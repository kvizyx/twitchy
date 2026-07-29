package oauth

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kvizyx/twitchy/helix"
)

func TestCoordinatedRetryCommit_reacquiresReloadsAndPublishesPendingPair(t *testing.T) {
	clock := newTestClock(time.Unix(60_000, 0))
	server := newCoordinatedOAuthServer(t)
	coordinator := newMemoryRefreshCoordinator()
	store := newCoordinatedTestStore(TokenPair{
		AccessToken:  "old-access",
		RefreshToken: "old-refresh",
		ExpiresIn:    time.Hour,
	})
	commitErr := errors.New("commit unavailable")
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

	if err := source.RetryCommit(context.Background()); err != nil {
		t.Fatal(err)
	}
	snapshot, err := source.Token(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got := snapshot.AccessToken(); got != "rotated-access" {
		t.Fatalf("active access token = %q, want pending rotation", got)
	}
	if got := coordinator.acquireCalls.Load(); got != 3 {
		t.Fatalf("lease acquisitions = %d, want initial, refresh, and retry", got)
	}
	if got := store.hooks.Load(); got != 2 {
		t.Fatalf("commit attempts = %d, want 2", got)
	}
}

func TestCoordinatedRetryCommit_adoptsPairCommittedByAnotherProcess(t *testing.T) {
	clock := newTestClock(time.Unix(70_000, 0))
	server := newCoordinatedOAuthServer(t)
	coordinator := newMemoryRefreshCoordinator()
	store := newCoordinatedTestStore(TokenPair{
		AccessToken:  "old-access",
		RefreshToken: "old-refresh",
		ExpiresIn:    time.Hour,
		Scopes:       []helix.AuthorizationScope{"scope:old"},
	})
	commitErr := errors.New("commit unavailable")
	store.hookFn = func(context.Context, TokenPair) error { return commitErr }
	source := newCoordinatedSource(t, clock, server, coordinator, store)
	clock.Advance(time.Hour - defaultRefreshSkew)
	if _, err := source.Token(context.Background()); !errors.Is(err, helix.ErrCredentialCommit) {
		t.Fatalf("initial refresh error = %v, want credential commit failure", err)
	}
	pending := store.pendingPair()
	pending.ExpiresIn = 30 * time.Minute
	pending.Scopes = []helix.AuthorizationScope{"scope:committed"}
	store.setPair(pending)

	if err := source.RetryCommit(context.Background()); err != nil {
		t.Fatal(err)
	}
	snapshot, err := source.Token(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got := snapshot.AccessToken(); got != "rotated-access" {
		t.Fatalf("active access token = %q, want externally committed pair", got)
	}
	if got := store.hooks.Load(); got != 1 {
		t.Fatalf("commit attempts = %d, want no duplicate hook", got)
	}
}

func TestCoordinatedRetryCommit_concurrentCallersCoalesce(t *testing.T) {
	clock := newTestClock(time.Unix(80_000, 0))
	server := newCoordinatedOAuthServer(t)
	coordinator := newMemoryRefreshCoordinator()
	store := newCoordinatedTestStore(TokenPair{
		AccessToken:  "old-access",
		RefreshToken: "old-refresh",
		ExpiresIn:    time.Hour,
	})
	commitErr := errors.New("commit unavailable")
	commitStarted := make(chan struct{}, 1)
	releaseCommit := make(chan struct{})
	var attempts atomic.Int32
	store.hookFn = func(ctx context.Context, pair TokenPair) error {
		if attempts.Add(1) == 1 {
			return commitErr
		}
		select {
		case commitStarted <- struct{}{}:
		default:
		}
		select {
		case <-releaseCommit:
			store.setPair(pair)
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	source := newCoordinatedSource(t, clock, server, coordinator, store)
	clock.Advance(time.Hour - defaultRefreshSkew)
	if _, err := source.Token(context.Background()); !errors.Is(err, helix.ErrCredentialCommit) {
		t.Fatalf("initial refresh error = %v, want credential commit failure", err)
	}

	results := make(chan error, 20)
	var callers sync.WaitGroup
	callers.Add(20)
	for range 20 {
		go func() {
			defer callers.Done()
			results <- source.RetryCommit(context.Background())
		}()
	}
	select {
	case <-commitStarted:
	case <-time.After(time.Second):
		source.mu.Lock()
		pendingPresent := source.pending != nil
		flightPresent := source.commitFlight != nil
		terminal := source.terminal
		source.mu.Unlock()
		t.Fatalf(
			"retry commit did not start: attempts=%d hooks=%d acquisitions=%d pending=%t flight=%t terminal=%v",
			attempts.Load(),
			store.hooks.Load(),
			coordinator.acquireCalls.Load(),
			pendingPresent,
			flightPresent,
			terminal,
		)
	}
	close(releaseCommit)
	callers.Wait()
	close(results)
	for err := range results {
		if err != nil {
			t.Fatal(err)
		}
	}
	if got := store.hooks.Load(); got != 2 {
		t.Fatalf("commit attempts = %d, want one retry", got)
	}
	if got := coordinator.acquireCalls.Load(); got != 3 {
		t.Fatalf("lease acquisitions = %d, want one retry acquisition", got)
	}
}
