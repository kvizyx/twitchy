package oauth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kvizyx/twitchy/helix"
)

type CredentialHook func(context.Context, TokenPair) error

type SourceOption func(*sourceOptions) error

type sourceOptions struct {
	refreshSkew time.Duration
}

func WithRefreshSkew(skew time.Duration) SourceOption {
	return func(options *sourceOptions) error {
		if skew < 0 || skew > 15*time.Minute {
			return ErrInvalidOption
		}
		options.refreshSkew = skew
		return nil
	}
}

type refreshFlight struct {
	generation uint64
	done       chan struct{}
	snapshot   helix.CredentialSnapshot
	err        error
}

type pendingCommit struct {
	pair       TokenPair
	generation uint64
}

type commitFlight struct {
	done chan struct{}
	err  error
}

func (source *RefreshingTokenSource) refreshGeneration(ctx context.Context, generation uint64, reason helix.RefreshReason) (helix.CredentialSnapshot, error) {
	source.mu.Lock()
	if err := source.lifecycleErrorLocked(); err != nil {
		source.mu.Unlock()
		return helix.CredentialSnapshot{}, err
	}
	if source.current.Generation() > generation {
		snapshot := source.current
		source.mu.Unlock()
		return snapshot, nil
	}
	if source.flight != nil {
		flight := source.flight
		source.mu.Unlock()
		select {
		case <-flight.done:
			return flight.snapshot, flight.err
		case <-ctx.Done():
			return helix.CredentialSnapshot{}, ctx.Err()
		}
	}
	flight := &refreshFlight{generation: generation, done: make(chan struct{})}
	source.flight = flight
	source.mu.Unlock()

	snapshot, rotated, err := source.performRefresh(reason)
	source.mu.Lock()
	if source.closed {
		snapshot = helix.CredentialSnapshot{}
		err = helix.ErrSessionClosed
	} else if err == nil {
		if source.terminal != nil {
			snapshot = helix.CredentialSnapshot{}
			err = source.terminal
		} else if source.current.Generation() != generation {
			snapshot = source.current
		} else {
			source.current = snapshot
			source.pair = rotated
		}
	}
	flight.snapshot = snapshot
	flight.err = err
	deleteFlight := source.flight == flight
	if deleteFlight {
		source.flight = nil
	}
	close(flight.done)
	source.mu.Unlock()
	return snapshot, err
}

func (source *RefreshingTokenSource) performRefresh(reason helix.RefreshReason) (helix.CredentialSnapshot, TokenPair, error) {
	source.mu.Lock()
	current := source.current
	pair := source.pair
	clientID := source.clientID
	source.mu.Unlock()
	if current.TokenClass() != helix.TokenClassUser || !current.Refreshable() || pair.RefreshToken == "" {
		return helix.CredentialSnapshot{}, TokenPair{}, helix.ErrNotRefreshable
	}
	if reason != helix.RefreshReasonExpired && reason != helix.RefreshReasonUnauthorized {
		return helix.CredentialSnapshot{}, TokenPair{}, helix.ErrNotRefreshable
	}
	if clientID == "" {
		validation, err := source.client.Validate(source.ctx, ValidateRequest{AccessToken: current.AccessToken()})
		if err != nil {
			return helix.CredentialSnapshot{}, TokenPair{}, err
		}
		clientID = validation.ClientID
		source.applyValidation(*validation)
	}
	rotated, err := source.client.Refresh(source.ctx, RefreshRequest{ClientID: clientID, RefreshToken: pair.RefreshToken})
	if err != nil {
		return helix.CredentialSnapshot{}, TokenPair{}, err
	}
	rotatedPair := *rotated
	if len(rotatedPair.Scopes) == 0 {
		rotatedPair.Scopes = append([]helix.AuthorizationScope(nil), current.Scopes()...)
	}
	if source.hook != nil {
		if err := source.hook(source.ctx, rotatedPair); err != nil {
			source.mu.Lock()
			if !source.closed && source.terminal == nil {
				source.terminal = fmt.Errorf("%w: %w", helix.ErrCredentialCommit, err)
				source.pending = &pendingCommit{pair: rotatedPair, generation: current.Generation()}
			}
			terminal := source.terminal
			source.mu.Unlock()
			if terminal == nil {
				return helix.CredentialSnapshot{}, TokenPair{}, helix.ErrSessionClosed
			}
			return helix.CredentialSnapshot{}, TokenPair{}, terminal
		}
	}
	return source.snapshotForPair(rotatedPair, current.Generation()+1), rotatedPair, nil
}

func (source *RefreshingTokenSource) RetryCommit(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	source.mu.Lock()
	if err := source.lifecycleErrorLocked(); err != nil && !errors.Is(err, helix.ErrCredentialCommit) {
		source.mu.Unlock()
		return err
	}
	if source.commitFlight != nil {
		flight := source.commitFlight
		source.mu.Unlock()
		select {
		case <-flight.done:
			return flight.err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	if source.pending == nil {
		source.mu.Unlock()
		return helix.ErrCredentialCommit
	}
	pending := *source.pending
	flight := &commitFlight{done: make(chan struct{})}
	source.commitFlight = flight
	source.mu.Unlock()

	err := error(nil)
	if source.hook != nil {
		err = source.hook(source.ctx, pending.pair)
	}
	source.mu.Lock()
	if source.closed {
		err = helix.ErrSessionClosed
	} else if err != nil {
		err = fmt.Errorf("%w: %w", helix.ErrCredentialCommit, err)
	} else if source.terminal != nil && !errors.Is(source.terminal, helix.ErrCredentialCommit) {
		err = source.terminal
	} else {
		source.current = source.snapshotForPair(pending.pair, pending.generation+1)
		source.pair = pending.pair
		source.pending = nil
		source.terminal = nil
	}
	flight.err = err
	source.commitFlight = nil
	close(flight.done)
	source.mu.Unlock()
	return err
}
