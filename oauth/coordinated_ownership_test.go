package oauth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kvizyx/twitchy/helix"
)

func TestCoordinatedLifecycle_leaseLossDuringHookRetainsPendingPair(t *testing.T) {
	clock := newTestClock(time.Unix(120_000, 0))
	server := newCoordinatedOAuthServer(t)
	coordinator := newMemoryRefreshCoordinator()
	store := newCoordinatedTestStore(TokenPair{
		AccessToken:  "old-access",
		RefreshToken: "old-refresh",
		ExpiresIn:    time.Hour,
	})
	source := newCoordinatedSource(t, clock, server, coordinator, store)
	drainCoordinatedRegistration(t, coordinator)
	hookStarted := make(chan struct{})
	store.mu.Lock()
	store.hookFn = func(ctx context.Context, _ TokenPair) error {
		close(hookStarted)
		<-ctx.Done()
		return ctx.Err()
	}
	store.mu.Unlock()
	clock.Advance(time.Hour - defaultRefreshSkew)

	result := make(chan error, 1)
	go func() {
		_, err := source.Token(context.Background())
		result <- err
	}()
	waitAcquireAttempt(t, coordinator)
	lease := coordinator.nextLease(t)
	select {
	case <-hookStarted:
	case <-time.After(time.Second):
		t.Fatal("credential hook did not start")
	}
	lease.lose(ErrRefreshLeaseLost)
	err := waitLifecycleResult(t, result)
	if !errors.Is(err, helix.ErrCredentialCommit) || !errors.Is(err, ErrRefreshLeaseLost) {
		t.Fatalf("refresh error = %v, want commit and lease loss", err)
	}
	requireNoActivation(t, source, "old-access")
	if got := store.pendingPair().AccessToken; got != "rotated-access" {
		t.Fatalf("pending access token = %q, want rotated pair retained", got)
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
		t.Fatalf("RetryCommit error = %v", err)
	}
}

func TestCoordinatedLifecycle_leaseLossBeforeHookSkipsCommitAndActivation(t *testing.T) {
	clock := newTestClock(time.Unix(125_000, 0))
	server := newCoordinatedOAuthServer(t)
	coordinator := newMemoryRefreshCoordinator()
	store := newCoordinatedTestStore(TokenPair{
		AccessToken:  "old-access",
		RefreshToken: "old-refresh",
		ExpiresIn:    time.Hour,
	})
	source := newCoordinatedSource(t, clock, server, coordinator, store)
	drainCoordinatedRegistration(t, coordinator)
	loadStarted := make(chan struct{})
	releaseLoad := make(chan struct{})
	store.mu.Lock()
	store.loadFn = func(ctx context.Context, _ string) (TokenPair, error) {
		close(loadStarted)
		select {
		case <-releaseLoad:
			return store.currentPair(), nil
		case <-ctx.Done():
			return TokenPair{}, ctx.Err()
		}
	}
	store.mu.Unlock()
	clock.Advance(time.Hour - defaultRefreshSkew)

	result := make(chan error, 1)
	go func() {
		_, err := source.Token(context.Background())
		result <- err
	}()
	waitAcquireAttempt(t, coordinator)
	lease := coordinator.nextLease(t)
	select {
	case <-loadStarted:
	case <-time.After(time.Second):
		t.Fatal("durable reload did not start")
	}
	lease.mu.Lock()
	lease.loseOnCall = 3
	lease.mu.Unlock()
	close(releaseLoad)
	err := waitLifecycleResult(t, result)
	if !errors.Is(err, helix.ErrCredentialCommit) || !errors.Is(err, ErrRefreshLeaseLost) {
		t.Fatalf("refresh error = %v, want commit and lease loss", err)
	}
	if got := store.hooks.Load(); got != 0 {
		t.Fatalf("durable commits = %d, want 0", got)
	}
	requireNoActivation(t, source, "old-access")
}
