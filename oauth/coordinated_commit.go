package oauth

import (
	"context"
	"errors"

	"github.com/kvizyx/twitchy/helix"
)

func (source *RefreshingTokenSource) performCoordinatedRetryCommit(
	ctx context.Context,
	pending pendingCommit,
) (returnErr error) {
	guard, err := source.acquireCoordinatedLease(ctx)
	if err != nil {
		return err
	}
	defer func() { returnErr = errors.Join(returnErr, guard.release()) }()

	durable, loadErr := source.coordination.loader(guard.lease.Context(), source.coordination.userID)
	if ownershipErr := source.coordinatedOwnershipError(guard.lease); loadErr != nil || ownershipErr != nil {
		return errors.Join(wrapCoordinatedLoadError(loadErr), ownershipErr)
	}
	state, err := source.coordinatedPendingState(pending.generation)
	if err != nil {
		return err
	}
	if coordinatedPairsDiffer(durable, state.pair) {
		if durable.ExpiresIn > 0 {
			if len(durable.Scopes) == 0 {
				durable.Scopes = append([]helix.AuthorizationScope(nil), state.current.Scopes()...)
			}
			validation, validationErr := source.client.Validate(
				guard.lease.Context(),
				ValidateRequest{AccessToken: durable.AccessToken},
			)
			if ownershipErr := source.coordinatedOwnershipError(guard.lease); validationErr != nil || ownershipErr != nil {
				return errors.Join(validationErr, ownershipErr)
			}
			candidate := source.candidateFromValidation(durable, *validation, pending.generation+1)
			_, err = source.activateCoordinatedCandidate(guard.lease, coordinatedActivation{
				generation:   pending.generation,
				candidate:    candidate,
				clearPending: true,
			})
			return err
		}
		_, err = source.refreshCoordinatedDurable(guard, coordinatedRemoteRefresh{
			durable:      durable,
			generation:   pending.generation,
			identity:     state.identity,
			fallback:     state.current.Scopes(),
			clearPending: true,
		})
		return err
	}

	if err = source.coordinatedOwnershipError(guard.lease); err != nil {
		return &coordinatedCommitError{cause: err}
	}
	if err = source.hook(guard.lease.Context(), pending.pair); err != nil {
		return source.coordinatedCommitFailure(
			pending.pair,
			pending.generation,
			errors.Join(err, guard.lease.Err()),
		)
	}
	if err = source.coordinatedOwnershipError(guard.lease); err != nil {
		return source.coordinatedCommitFailure(pending.pair, pending.generation, err)
	}
	candidate := coordinatedCandidate{
		pair:     cloneCoordinatedPair(pending.pair),
		clientID: state.identity.clientID,
		userID:   state.identity.userID,
		snapshot: source.coordinatedSnapshot(pending.pair, pending.generation+1, state.identity),
	}
	_, err = source.activateCoordinatedCandidate(guard.lease, coordinatedActivation{
		generation:   pending.generation,
		candidate:    candidate,
		clearPending: true,
	})
	return err
}

func (source *RefreshingTokenSource) coordinatedPendingState(generation uint64) (coordinatedRefreshState, error) {
	source.mu.Lock()
	defer source.mu.Unlock()
	if source.closed {
		return coordinatedRefreshState{}, helix.ErrSessionClosed
	}
	if source.terminal != nil && !errors.Is(source.terminal, helix.ErrCredentialCommit) {
		return coordinatedRefreshState{}, source.terminal
	}
	if source.current.Generation() != generation {
		return coordinatedRefreshState{
			current: source.current,
			pair:    cloneCoordinatedPair(source.pair),
			identity: coordinatedIdentity{
				clientID: source.clientID,
				userID:   source.userID,
			},
		}, nil
	}
	return coordinatedRefreshState{
		current: source.current,
		pair:    cloneCoordinatedPair(source.pair),
		identity: coordinatedIdentity{
			clientID: source.clientID,
			userID:   source.userID,
		},
	}, nil
}
