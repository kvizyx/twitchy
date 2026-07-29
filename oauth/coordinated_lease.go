package oauth

import (
	"context"
	"errors"
	"fmt"

	"github.com/kvizyx/twitchy/helix"
)

type coordinatedLeaseGuard struct {
	lease      RefreshLease
	workCtx    context.Context
	cancel     context.CancelFunc
	stopSource func() bool
}

func (source *RefreshingTokenSource) acquireCoordinatedLease(ctx context.Context) (*coordinatedLeaseGuard, error) {
	workCtx, cancel := context.WithCancel(ctx)
	stopSource := context.AfterFunc(source.ctx, cancel)
	coordination := source.coordination
	lease, err := coordination.coordinator.Acquire(workCtx, coordination.userID)
	if err != nil {
		stopSource()
		cancel()
		if closeErr := source.closedError(); closeErr != nil {
			return nil, errors.Join(closeErr, fmt.Errorf("acquire refresh lease: %w", err))
		}
		return nil, fmt.Errorf("acquire refresh lease: %w", err)
	}
	if isNilCoordinatorValue(lease) {
		stopSource()
		cancel()
		return nil, ErrInvalidOption
	}
	guard := &coordinatedLeaseGuard{
		lease:      lease,
		workCtx:    workCtx,
		cancel:     cancel,
		stopSource: stopSource,
	}
	if lease.Context() == nil {
		return nil, errors.Join(ErrInvalidOption, guard.release())
	}
	return guard, nil
}

func (guard *coordinatedLeaseGuard) release() error {
	releaseContext, cancel := context.WithTimeout(
		context.WithoutCancel(guard.workCtx),
		coordinatedReleaseTimeout,
	)
	defer cancel()
	releaseErr := guard.lease.Release(releaseContext)
	guard.stopSource()
	guard.cancel()
	return releaseErr
}

func (source *RefreshingTokenSource) closedError() error {
	source.mu.Lock()
	defer source.mu.Unlock()
	if source.closed {
		return helix.ErrSessionClosed
	}
	return nil
}

func (source *RefreshingTokenSource) coordinatedOwnershipError(lease RefreshLease) error {
	if err := source.closedError(); err != nil {
		return err
	}
	return lease.Err()
}

type coordinatedCandidate struct {
	pair     TokenPair
	snapshot helix.CredentialSnapshot
	clientID string
	userID   string
}

type coordinatedIdentity struct {
	clientID string
	userID   string
}

type coordinatedActivation struct {
	generation   uint64
	candidate    coordinatedCandidate
	clearPending bool
}

func (source *RefreshingTokenSource) candidateFromValidation(
	pair TokenPair,
	validation Validation,
	generation uint64,
) coordinatedCandidate {
	if len(validation.Scopes) > 0 {
		pair.Scopes = append([]helix.AuthorizationScope(nil), validation.Scopes...)
	}
	pair.ExpiresIn = validation.ExpiresIn
	return coordinatedCandidate{
		pair:     pair,
		clientID: validation.ClientID,
		userID:   validation.UserID,
		snapshot: source.coordinatedSnapshot(pair, generation, coordinatedIdentity{
			clientID: validation.ClientID,
			userID:   validation.UserID,
		}),
	}
}

func (source *RefreshingTokenSource) coordinatedSnapshot(
	pair TokenPair,
	generation uint64,
	identity coordinatedIdentity,
) helix.CredentialSnapshot {
	return helix.NewCredentialSnapshot(helix.Credential{
		AccessToken: pair.AccessToken,
		ClientID:    identity.clientID,
		TokenClass:  helix.TokenClassUser,
		UserID:      identity.userID,
		Scopes:      append([]helix.AuthorizationScope(nil), pair.Scopes...),
		ExpiresAt:   source.clock.Now().Add(pair.ExpiresIn),
		Refreshable: pair.RefreshToken != "",
		Generation:  generation,
	})
}

func (source *RefreshingTokenSource) activateCoordinatedCandidate(
	lease RefreshLease,
	activation coordinatedActivation,
) (helix.CredentialSnapshot, error) {
	if err := source.coordinatedOwnershipError(lease); err != nil {
		return helix.CredentialSnapshot{}, err
	}
	source.mu.Lock()
	defer source.mu.Unlock()
	if source.closed {
		return helix.CredentialSnapshot{}, helix.ErrSessionClosed
	}
	if err := lease.Err(); err != nil {
		return helix.CredentialSnapshot{}, err
	}
	if source.terminal != nil && !(activation.clearPending && errors.Is(source.terminal, helix.ErrCredentialCommit)) {
		return helix.CredentialSnapshot{}, source.terminal
	}
	if source.current.Generation() != activation.generation {
		return source.current, nil
	}
	source.current = activation.candidate.snapshot
	source.pair = cloneCoordinatedPair(activation.candidate.pair)
	source.clientID = activation.candidate.clientID
	source.userID = activation.candidate.userID
	if activation.clearPending {
		source.pending = nil
		source.terminal = nil
	}
	return activation.candidate.snapshot, nil
}

func (source *RefreshingTokenSource) coordinatedCommitFailure(pair TokenPair, generation uint64, cause error) error {
	commitErr := &coordinatedCommitError{cause: cause}
	source.mu.Lock()
	defer source.mu.Unlock()
	if source.closed {
		return errors.Join(helix.ErrSessionClosed, commitErr)
	}
	if source.commitFlight != nil {
		select {
		case <-source.commitFlight.done:
			source.commitFlight = nil
		default:
		}
	}
	source.pending = &pendingCommit{pair: cloneCoordinatedPair(pair), generation: generation}
	source.terminal = commitErr
	return commitErr
}

func (source *RefreshingTokenSource) coordinatedRotationFailure(cause error) error {
	rotationErr := &CredentialRotationError{cause: cause}
	source.mu.Lock()
	defer source.mu.Unlock()
	if source.closed {
		return errors.Join(helix.ErrSessionClosed, rotationErr)
	}
	source.terminal = rotationErr
	return rotationErr
}

func cloneCoordinatedPair(pair TokenPair) TokenPair {
	pair.Scopes = append([]helix.AuthorizationScope(nil), pair.Scopes...)
	return pair
}

func coordinatedPairsDiffer(first TokenPair, second TokenPair) bool {
	return first.AccessToken != second.AccessToken || first.RefreshToken != second.RefreshToken
}
