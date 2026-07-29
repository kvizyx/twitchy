package oauth

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kvizyx/twitchy/helix"
)

func TestCoordinatedLifecycle_closeDuringHTTPReleasesWithoutCommitOrActivation(t *testing.T) {
	clock := newTestClock(time.Unix(140_000, 0))
	server := newCoordinatedOAuthServer(t)
	server.setRefreshFn(func(_ http.ResponseWriter, request *http.Request, _ int32) {
		<-request.Context().Done()
	})
	coordinator := newMemoryRefreshCoordinator()
	store := newCoordinatedTestStore(TokenPair{
		AccessToken:  "unchanged-access",
		RefreshToken: "unchanged-refresh",
		ExpiresIn:    time.Hour,
	})
	source := newCoordinatedSource(t, clock, server, coordinator, store)
	drainCoordinatedRegistration(t, coordinator)
	clock.Advance(time.Hour - defaultRefreshSkew)

	result := make(chan error, 1)
	go func() {
		_, err := source.Token(context.Background())
		result <- err
	}()
	waitAcquireAttempt(t, coordinator)
	_ = coordinator.nextLease(t)
	select {
	case <-server.refreshStarted:
	case <-time.After(time.Second):
		t.Fatal("remote refresh did not start")
	}
	if err := source.Close(); err != nil {
		t.Fatal(err)
	}
	if err := waitLifecycleResult(t, result); !errors.Is(err, helix.ErrSessionClosed) {
		t.Fatalf("refresh error = %v, want session closed", err)
	}
	if got := store.hooks.Load(); got != 0 {
		t.Fatalf("durable commits = %d, want 0", got)
	}
	requireNoActivation(t, source, "unchanged-access")
}

func TestCoordinatedLifecycle_closeDuringHookPreventsActivation(t *testing.T) {
	clock := newTestClock(time.Unix(150_000, 0))
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
	_ = coordinator.nextLease(t)
	select {
	case <-hookStarted:
	case <-time.After(time.Second):
		t.Fatal("credential hook did not start")
	}
	if err := source.Close(); err != nil {
		t.Fatal(err)
	}
	if err := waitLifecycleResult(t, result); !errors.Is(err, helix.ErrSessionClosed) {
		t.Fatalf("refresh error = %v, want session closed", err)
	}
	requireNoActivation(t, source, "old-access")
}

func TestCoordinatedLifecycle_closeDuringRetryCommitPreventsActivation(t *testing.T) {
	clock := newTestClock(time.Unix(160_000, 0))
	server := newCoordinatedOAuthServer(t)
	coordinator := newMemoryRefreshCoordinator()
	store := newCoordinatedTestStore(TokenPair{
		AccessToken:  "old-access",
		RefreshToken: "old-refresh",
		ExpiresIn:    time.Hour,
	})
	commitErr := errors.New("commit unavailable")
	var attempts atomic.Int32
	store.hookFn = func(context.Context, TokenPair) error {
		if attempts.Add(1) == 1 {
			return commitErr
		}
		return nil
	}
	source := newCoordinatedSource(t, clock, server, coordinator, store)
	drainCoordinatedRegistration(t, coordinator)
	clock.Advance(time.Hour - defaultRefreshSkew)
	if _, err := source.Token(context.Background()); !errors.Is(err, helix.ErrCredentialCommit) {
		t.Fatalf("initial refresh error = %v, want commit failure", err)
	}
	waitAcquireAttempt(t, coordinator)
	_ = coordinator.nextLease(t)

	hookStarted := make(chan struct{})
	store.mu.Lock()
	store.hookFn = func(ctx context.Context, _ TokenPair) error {
		close(hookStarted)
		<-ctx.Done()
		return ctx.Err()
	}
	store.mu.Unlock()
	result := make(chan error, 1)
	go func() { result <- source.RetryCommit(context.Background()) }()
	waitAcquireAttempt(t, coordinator)
	_ = coordinator.nextLease(t)
	select {
	case <-hookStarted:
	case <-time.After(time.Second):
		t.Fatal("retry hook did not start")
	}
	if err := source.Close(); err != nil {
		t.Fatal(err)
	}
	if err := waitLifecycleResult(t, result); !errors.Is(err, helix.ErrSessionClosed) {
		t.Fatalf("RetryCommit error = %v, want session closed", err)
	}
	requireNoActivation(t, source, "old-access")
}
