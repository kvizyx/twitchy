package oauth

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kvizyx/twitchy/helix"
)

func TestCoordinatedRetryCommit_reloadFailureRetainsPendingPair(t *testing.T) {
	clock := newTestClock(time.Unix(90_000, 0))
	server := newCoordinatedOAuthServer(t)
	coordinator := newMemoryRefreshCoordinator()
	store := newCoordinatedTestStore(TokenPair{
		AccessToken:  "old-access",
		RefreshToken: "old-refresh",
		ExpiresIn:    time.Hour,
	})
	commitErr := errors.New("commit unavailable")
	reloadErr := errors.New("reload unavailable")
	var attempts atomic.Int32
	store.hookFn = func(ctx context.Context, pair TokenPair) error {
		if attempts.Add(1) == 1 {
			return commitErr
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		store.setPair(pair)
		return nil
	}
	source := newCoordinatedSource(t, clock, server, coordinator, store)
	clock.Advance(time.Hour - defaultRefreshSkew)
	if _, err := source.Token(context.Background()); !errors.Is(err, helix.ErrCredentialCommit) {
		t.Fatalf("initial refresh error = %v, want credential commit failure", err)
	}
	store.mu.Lock()
	store.loadErr = reloadErr
	store.mu.Unlock()

	if err := source.RetryCommit(context.Background()); !errors.Is(err, reloadErr) {
		t.Fatalf("RetryCommit error = %v, want reload failure", err)
	}
	store.mu.Lock()
	store.loadErr = nil
	store.mu.Unlock()
	if err := source.RetryCommit(context.Background()); err != nil {
		t.Fatalf("second RetryCommit error = %v", err)
	}
	if got := store.hooks.Load(); got != 2 {
		t.Fatalf("commit attempts = %d, want pending pair retained", got)
	}
}

func TestCoordinatedRetryCommit_completedFlightDoesNotMaskLaterPendingPair(t *testing.T) {
	clock := newTestClock(time.Unix(95_000, 0))
	server := newCoordinatedOAuthServer(t)
	coordinator := newMemoryRefreshCoordinator()
	store := newCoordinatedTestStore(TokenPair{
		AccessToken:  "old-access",
		RefreshToken: "old-refresh",
		ExpiresIn:    time.Hour,
	})
	firstCommitErr := errors.New("first commit unavailable")
	var attempts atomic.Int32
	store.hookFn = func(ctx context.Context, pair TokenPair) error {
		if attempts.Add(1) == 1 {
			return firstCommitErr
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		store.setPair(pair)
		return nil
	}
	source := newCoordinatedSource(t, clock, server, coordinator, store)
	clock.Advance(time.Hour - defaultRefreshSkew)
	if _, err := source.Token(context.Background()); !errors.Is(err, helix.ErrCredentialCommit) {
		t.Fatalf("first refresh error = %v, want commit failure", err)
	}
	if err := source.RetryCommit(context.Background()); err != nil {
		t.Fatalf("first RetryCommit error = %v", err)
	}
	if err := source.RetryCommit(context.Background()); err != nil {
		t.Fatalf("idempotent RetryCommit error = %v", err)
	}
	snapshot, err := source.Token(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	secondCommitErr := errors.New("second commit unavailable")
	store.mu.Lock()
	store.hookFn = func(context.Context, TokenPair) error { return secondCommitErr }
	store.mu.Unlock()
	_, err = source.Refresh(context.Background(), snapshot, helix.RefreshReasonUnauthorized)
	if !errors.Is(err, secondCommitErr) {
		t.Fatalf("second refresh error = %v, want new commit failure", err)
	}
	store.mu.Lock()
	store.hookFn = func(ctx context.Context, pair TokenPair) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		store.setPair(pair)
		return nil
	}
	store.mu.Unlock()
	if err = source.RetryCommit(context.Background()); err != nil {
		t.Fatalf("second RetryCommit error = %v", err)
	}
	if got := store.hooks.Load(); got != 4 {
		t.Fatalf("commit attempts = %d, want both pending pairs retried", got)
	}
}
