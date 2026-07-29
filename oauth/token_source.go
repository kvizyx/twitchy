package oauth

import (
	"context"
	"sync"
	"time"

	"github.com/kvizyx/twitchy/helix"
)

const defaultRefreshSkew = time.Minute

type RefreshingTokenSource struct {
	client       *Client
	hook         CredentialHook
	clock        helix.Clock
	skew         time.Duration
	coordination *sourceCoordination

	mu           sync.Mutex
	current      helix.CredentialSnapshot
	pair         TokenPair
	clientID     string
	userID       string
	flight       *refreshFlight
	commitFlight *commitFlight
	pending      *pendingCommit
	terminal     error
	closed       bool
	ctx          context.Context
	cancel       context.CancelFunc
	done         chan struct{}
}

func NewRefreshingTokenSource(client *Client, pair TokenPair, hook CredentialHook, options ...SourceOption) (*RefreshingTokenSource, error) {
	if err := client.validClient(); err != nil {
		return nil, err
	}
	configuration := sourceOptions{refreshSkew: defaultRefreshSkew}
	for _, option := range options {
		if option == nil {
			return nil, ErrInvalidOption
		}
		if err := option(&configuration); err != nil {
			return nil, err
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	source := &RefreshingTokenSource{
		client:       client,
		hook:         hook,
		clock:        client.clock,
		skew:         configuration.refreshSkew,
		coordination: configuration.coordination,
		pair:         pair,
		ctx:          ctx,
		cancel:       cancel,
		done:         make(chan struct{}),
	}
	source.current = source.snapshotForPair(pair, 0)
	return source, nil
}

func (source *RefreshingTokenSource) Token(ctx context.Context) (helix.CredentialSnapshot, error) {
	if err := ctx.Err(); err != nil {
		return helix.CredentialSnapshot{}, err
	}
	source.mu.Lock()
	if err := source.lifecycleErrorLocked(); err != nil {
		source.mu.Unlock()
		return helix.CredentialSnapshot{}, err
	}
	snapshot := source.current
	flight := source.flight
	due := !snapshot.ExpiresAt().After(source.clock.Now().Add(source.skew))
	refreshable := snapshot.TokenClass() == helix.TokenClassUser && snapshot.Refreshable()
	source.mu.Unlock()
	if flight != nil && flight.generation == snapshot.Generation() {
		select {
		case <-flight.done:
			return flight.snapshot, flight.err
		case <-ctx.Done():
			return helix.CredentialSnapshot{}, ctx.Err()
		}
	}
	if !due || !refreshable {
		return snapshot, nil
	}
	return source.refreshGeneration(ctx, snapshot.Generation(), helix.RefreshReasonExpired)
}

func (source *RefreshingTokenSource) Refresh(ctx context.Context, credential helix.CredentialSnapshot, reason helix.RefreshReason) (helix.CredentialSnapshot, error) {
	if err := ctx.Err(); err != nil {
		return helix.CredentialSnapshot{}, err
	}
	source.mu.Lock()
	if err := source.lifecycleErrorLocked(); err != nil {
		source.mu.Unlock()
		return helix.CredentialSnapshot{}, err
	}
	if source.current.Generation() > credential.Generation() {
		snapshot := source.current
		source.mu.Unlock()
		return snapshot, nil
	}
	if source.current.TokenClass() != helix.TokenClassUser || !source.current.Refreshable() {
		source.mu.Unlock()
		return helix.CredentialSnapshot{}, helix.ErrNotRefreshable
	}
	generation := source.current.Generation()
	source.mu.Unlock()
	return source.refreshGeneration(ctx, generation, reason)
}

func (source *RefreshingTokenSource) Close() error {
	source.mu.Lock()
	if source.closed {
		source.mu.Unlock()
		return nil
	}
	source.closed = true
	source.cancel()
	close(source.done)
	source.mu.Unlock()
	return nil
}

func (source *RefreshingTokenSource) lifecycleErrorLocked() error {
	if source.closed {
		return helix.ErrSessionClosed
	}
	return source.terminal
}

func (source *RefreshingTokenSource) snapshotForPair(pair TokenPair, generation uint64) helix.CredentialSnapshot {
	tokenClass := helix.TokenClassApp
	refreshable := pair.RefreshToken != ""
	if refreshable {
		tokenClass = helix.TokenClassUser
	}
	return helix.NewCredentialSnapshot(helix.Credential{
		AccessToken: pair.AccessToken,
		ClientID:    source.clientID,
		TokenClass:  tokenClass,
		UserID:      source.userID,
		Scopes:      append([]helix.AuthorizationScope(nil), pair.Scopes...),
		ExpiresAt:   source.clock.Now().Add(pair.ExpiresIn),
		Refreshable: refreshable,
		Generation:  generation,
	})
}

func (source *RefreshingTokenSource) applyValidation(validation Validation) {
	source.mu.Lock()
	defer source.mu.Unlock()
	if source.closed || source.terminal != nil {
		return
	}
	source.clientID = validation.ClientID
	source.userID = validation.UserID
	if len(validation.Scopes) > 0 {
		source.pair.Scopes = append([]helix.AuthorizationScope(nil), validation.Scopes...)
	}
	credential := source.current
	credential = helix.NewCredentialSnapshot(helix.Credential{
		AccessToken: credential.AccessToken(),
		ClientID:    source.clientID,
		TokenClass:  helix.TokenClassUser,
		UserID:      source.userID,
		Scopes:      source.pair.Scopes,
		ExpiresAt:   source.clock.Now().Add(validation.ExpiresIn),
		Refreshable: credential.Refreshable(),
		Generation:  credential.Generation(),
	})
	source.current = credential
}
