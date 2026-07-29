package oauth

import (
	"context"
	"errors"

	"github.com/kvizyx/twitchy/helix"
)

type coordinatedRefreshState struct {
	current  helix.CredentialSnapshot
	pair     TokenPair
	identity coordinatedIdentity
}

type coordinatedRemoteRefresh struct {
	durable      TokenPair
	generation   uint64
	identity     coordinatedIdentity
	fallback     []helix.AuthorizationScope
	clearPending bool
}

func (source *RefreshingTokenSource) performCoordinatedRefresh(
	ctx context.Context,
	generation uint64,
	reason helix.RefreshReason,
) (snapshot helix.CredentialSnapshot, returnErr error) {
	state, err := source.coordinatedRefreshState(generation, reason)
	if err != nil {
		return helix.CredentialSnapshot{}, err
	}
	if state.current.Generation() != generation {
		return state.current, nil
	}
	guard, err := source.acquireCoordinatedLease(ctx)
	if err != nil {
		return helix.CredentialSnapshot{}, err
	}
	defer func() { returnErr = errors.Join(returnErr, guard.release()) }()

	durable, loadErr := source.coordination.loader(guard.lease.Context(), source.coordination.userID)
	if ownershipErr := source.coordinatedOwnershipError(guard.lease); loadErr != nil || ownershipErr != nil {
		return helix.CredentialSnapshot{}, errors.Join(wrapCoordinatedLoadError(loadErr), ownershipErr)
	}
	if coordinatedPairsDiffer(durable, state.pair) && durable.ExpiresIn > 0 {
		if len(durable.Scopes) == 0 {
			durable.Scopes = append([]helix.AuthorizationScope(nil), state.current.Scopes()...)
		}
		validation, validationErr := source.client.Validate(
			guard.lease.Context(),
			ValidateRequest{AccessToken: durable.AccessToken},
		)
		if ownershipErr := source.coordinatedOwnershipError(guard.lease); validationErr != nil || ownershipErr != nil {
			return helix.CredentialSnapshot{}, errors.Join(validationErr, ownershipErr)
		}
		candidate := source.candidateFromValidation(durable, *validation, generation+1)
		return source.activateCoordinatedCandidate(guard.lease, coordinatedActivation{
			generation: generation,
			candidate:  candidate,
		})
	}
	return source.refreshCoordinatedDurable(guard, coordinatedRemoteRefresh{
		durable:    durable,
		generation: generation,
		identity:   state.identity,
		fallback:   state.current.Scopes(),
	})
}

func (source *RefreshingTokenSource) coordinatedRefreshState(
	generation uint64,
	reason helix.RefreshReason,
) (coordinatedRefreshState, error) {
	source.mu.Lock()
	defer source.mu.Unlock()
	if err := source.lifecycleErrorLocked(); err != nil {
		return coordinatedRefreshState{}, err
	}
	if source.current.Generation() != generation {
		return coordinatedRefreshState{current: source.current}, nil
	}
	if source.current.TokenClass() != helix.TokenClassUser || !source.current.Refreshable() {
		return coordinatedRefreshState{}, helix.ErrNotRefreshable
	}
	if reason != helix.RefreshReasonExpired && reason != helix.RefreshReasonUnauthorized {
		return coordinatedRefreshState{}, helix.ErrNotRefreshable
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

func (source *RefreshingTokenSource) refreshCoordinatedDurable(
	guard *coordinatedLeaseGuard,
	request coordinatedRemoteRefresh,
) (helix.CredentialSnapshot, error) {
	if request.durable.RefreshToken == "" {
		return helix.CredentialSnapshot{}, helix.ErrNotRefreshable
	}
	if request.identity.clientID == "" {
		validation, err := source.client.Validate(
			guard.lease.Context(),
			ValidateRequest{AccessToken: request.durable.AccessToken},
		)
		if err != nil {
			return helix.CredentialSnapshot{}, err
		}
		request.identity = coordinatedIdentity{clientID: validation.ClientID, userID: validation.UserID}
		if len(validation.Scopes) > 0 {
			request.durable.Scopes = append([]helix.AuthorizationScope(nil), validation.Scopes...)
		}
	}
	if err := source.coordinatedOwnershipError(guard.lease); err != nil {
		return helix.CredentialSnapshot{}, err
	}

	// OAuth rotation has an unavoidable crash window: process death after the
	// remote rotation and before the durable hook requires user reauthorization.
	rotated, err := source.client.Refresh(guard.lease.Context(), RefreshRequest{
		ClientID:     request.identity.clientID,
		RefreshToken: request.durable.RefreshToken,
	})
	if err != nil {
		return helix.CredentialSnapshot{}, source.coordinatedRefreshFailure(guard.workCtx, guard.lease, err)
	}
	rotatedPair := cloneCoordinatedPair(*rotated)
	if len(rotatedPair.Scopes) == 0 {
		rotatedPair.Scopes = append([]helix.AuthorizationScope(nil), request.durable.Scopes...)
		if len(rotatedPair.Scopes) == 0 {
			rotatedPair.Scopes = append([]helix.AuthorizationScope(nil), request.fallback...)
		}
	}
	if err := source.coordinatedOwnershipError(guard.lease); err != nil {
		return helix.CredentialSnapshot{}, source.coordinatedCommitFailure(
			rotatedPair,
			request.generation,
			err,
		)
	}
	if err := source.hook(guard.lease.Context(), rotatedPair); err != nil {
		return helix.CredentialSnapshot{}, source.coordinatedCommitFailure(
			rotatedPair,
			request.generation,
			errors.Join(err, guard.lease.Err()),
		)
	}
	if err := source.coordinatedOwnershipError(guard.lease); err != nil {
		return helix.CredentialSnapshot{}, source.coordinatedCommitFailure(
			rotatedPair,
			request.generation,
			err,
		)
	}
	candidate := coordinatedCandidate{
		pair:     rotatedPair,
		clientID: request.identity.clientID,
		userID:   request.identity.userID,
		snapshot: source.coordinatedSnapshot(rotatedPair, request.generation+1, request.identity),
	}
	return source.activateCoordinatedCandidate(guard.lease, coordinatedActivation{
		generation:   request.generation,
		candidate:    candidate,
		clearPending: request.clearPending,
	})
}
