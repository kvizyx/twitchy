package oauth

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/kvizyx/twitchy/helix"
)

func TestCoordinatedRefresh_refreshesAndCommitsOnceAcrossSources(t *testing.T) {
	clock := newTestClock(time.Unix(10_000, 0))
	server := newCoordinatedOAuthServer(t)
	coordinator := newMemoryRefreshCoordinator()
	store := newCoordinatedTestStore(TokenPair{
		AccessToken:  "initial-access",
		RefreshToken: "initial-refresh",
		ExpiresIn:    time.Hour,
		Scopes:       []helix.AuthorizationScope{"scope:initial"},
		TokenType:    "bearer",
	})
	first := newCoordinatedSource(t, clock, server, coordinator, store)
	second := newCoordinatedSource(t, clock, server, coordinator, store)
	clock.Advance(time.Hour - defaultRefreshSkew)

	type result struct {
		snapshot helix.CredentialSnapshot
		err      error
	}
	results := make(chan result, 2)
	var workers sync.WaitGroup
	workers.Add(2)
	for _, source := range []*RefreshingTokenSource{first, second} {
		go func() {
			defer workers.Done()
			snapshot, err := source.Token(context.Background())
			results <- result{snapshot: snapshot, err: err}
		}()
	}
	workers.Wait()
	close(results)

	for refreshResult := range results {
		if refreshResult.err != nil {
			t.Fatal(refreshResult.err)
		}
		if got := refreshResult.snapshot.AccessToken(); got != "rotated-access" {
			t.Fatalf("active access token = %q, want durable rotation", got)
		}
	}
	if got := server.refreshes.Load(); got != 1 {
		t.Fatalf("remote refreshes = %d, want 1", got)
	}
	if got := store.hooks.Load(); got != 1 {
		t.Fatalf("durable commits = %d, want 1", got)
	}
	if got := coordinator.acquireCalls.Load(); got != 4 {
		t.Fatalf("lease acquisitions = %d, want 4", got)
	}
}

func TestCoordinatedRefresh_adoptsValidDurablePairWithoutRemoteRefresh(t *testing.T) {
	clock := newTestClock(time.Unix(20_000, 0))
	server := newCoordinatedOAuthServer(t)
	coordinator := newMemoryRefreshCoordinator()
	store := newCoordinatedTestStore(TokenPair{
		AccessToken:  "local-access",
		RefreshToken: "local-refresh",
		ExpiresIn:    time.Hour,
		Scopes:       []helix.AuthorizationScope{"scope:local"},
	})
	source := newCoordinatedSource(t, clock, server, coordinator, store)
	store.setPair(TokenPair{
		AccessToken:  "durable-access",
		RefreshToken: "durable-refresh",
		ExpiresIn:    30 * time.Minute,
		Scopes:       []helix.AuthorizationScope{"scope:durable"},
		TokenType:    "bearer",
	})
	clock.Advance(time.Hour - defaultRefreshSkew)

	snapshot, err := source.Token(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got := snapshot.AccessToken(); got != "durable-access" {
		t.Fatalf("active access token = %q, want durable pair", got)
	}
	if got := snapshot.Scopes(); len(got) != 1 || got[0] != "scope:durable" {
		t.Fatalf("adopted scopes = %v, want durable scopes", got)
	}
	if got := server.refreshes.Load(); got != 0 {
		t.Fatalf("remote refreshes = %d, want 0", got)
	}
	if got := store.hooks.Load(); got != 0 {
		t.Fatalf("durable commits = %d, want 0", got)
	}
}

func TestCoordinatedRefresh_adoptionPreservesCurrentScopesWhenDurableScopesOmitted(t *testing.T) {
	clock := newTestClock(time.Unix(25_000, 0))
	server := newCoordinatedOAuthServer(t)
	coordinator := newMemoryRefreshCoordinator()
	store := newCoordinatedTestStore(TokenPair{
		AccessToken:  "local-access",
		RefreshToken: "local-refresh",
		ExpiresIn:    time.Hour,
		Scopes:       []helix.AuthorizationScope{"scope:current"},
	})
	source := newCoordinatedSource(t, clock, server, coordinator, store)
	store.setPair(TokenPair{
		AccessToken:  "durable-access",
		RefreshToken: "durable-refresh",
		ExpiresIn:    30 * time.Minute,
		TokenType:    "bearer",
	})
	clock.Advance(time.Hour - defaultRefreshSkew)

	snapshot, err := source.Token(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got := snapshot.Scopes(); len(got) != 1 || got[0] != "scope:current" {
		t.Fatalf("adopted scopes = %v, want current scopes preserved", got)
	}
	if got := server.refreshes.Load(); got != 0 {
		t.Fatalf("remote refreshes = %d, want 0", got)
	}
}

func TestCoordinatedRefresh_usesExpiredDurableRefreshTokenAndPreservesScopes(t *testing.T) {
	clock := newTestClock(time.Unix(30_000, 0))
	server := newCoordinatedOAuthServer(t)
	coordinator := newMemoryRefreshCoordinator()
	store := newCoordinatedTestStore(TokenPair{
		AccessToken:  "local-access",
		RefreshToken: "local-refresh",
		ExpiresIn:    time.Hour,
		Scopes:       []helix.AuthorizationScope{"scope:local"},
	})
	source := newCoordinatedSource(t, clock, server, coordinator, store)
	store.setPair(TokenPair{
		AccessToken:  "durable-expired-access",
		RefreshToken: "durable-expired-refresh",
		ExpiresIn:    0,
		Scopes:       []helix.AuthorizationScope{"scope:durable"},
		TokenType:    "bearer",
	})
	clock.Advance(time.Hour - defaultRefreshSkew)

	snapshot, err := source.Token(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	select {
	case refreshToken := <-server.refreshInputs:
		if refreshToken != "durable-expired-refresh" {
			t.Fatalf("refresh input = %q, want durable refresh token", refreshToken)
		}
	default:
		t.Fatal("remote refresh input was not captured")
	}
	if got := snapshot.Scopes(); len(got) != 1 || got[0] != "scope:durable" {
		t.Fatalf("rotated scopes = %v, want durable scopes", got)
	}
	if got := store.currentPair().RefreshToken; got != "rotated-refresh" {
		t.Fatalf("committed refresh token = %q, want rotated refresh token", got)
	}
}

func TestCoordinatedRefresh_reloadFailureDoesNotRefreshCommitOrActivate(t *testing.T) {
	clock := newTestClock(time.Unix(40_000, 0))
	server := newCoordinatedOAuthServer(t)
	coordinator := newMemoryRefreshCoordinator()
	store := newCoordinatedTestStore(TokenPair{
		AccessToken:  "unchanged-access",
		RefreshToken: "unchanged-refresh",
		ExpiresIn:    time.Hour,
	})
	source := newCoordinatedSource(t, clock, server, coordinator, store)
	reloadErr := errors.New("durable store unavailable")
	store.mu.Lock()
	store.loadErr = reloadErr
	store.mu.Unlock()
	clock.Advance(time.Hour - defaultRefreshSkew)

	_, err := source.Token(context.Background())
	if !errors.Is(err, reloadErr) {
		t.Fatalf("refresh error = %v, want durable reload failure", err)
	}
	if got := server.refreshes.Load(); got != 0 {
		t.Fatalf("remote refreshes = %d, want 0", got)
	}
	if got := store.hooks.Load(); got != 0 {
		t.Fatalf("durable commits = %d, want 0", got)
	}
	requireNoActivation(t, source, "unchanged-access")
}

func TestCoordinatedRefresh_doesNotRetryOldRefreshTokenAfterCommitFailure(t *testing.T) {
	clock := newTestClock(time.Unix(50_000, 0))
	server := newCoordinatedOAuthServer(t)
	coordinator := newMemoryRefreshCoordinator()
	store := newCoordinatedTestStore(TokenPair{
		AccessToken:  "old-access",
		RefreshToken: "old-refresh",
		ExpiresIn:    time.Hour,
	})
	commitErr := errors.New("commit unavailable")
	store.hookFn = func(context.Context, TokenPair) error { return commitErr }
	source := newCoordinatedSource(t, clock, server, coordinator, store)
	clock.Advance(time.Hour - defaultRefreshSkew)

	_, err := source.Token(context.Background())
	if !errors.Is(err, helix.ErrCredentialCommit) || !errors.Is(err, commitErr) {
		t.Fatalf("refresh error = %v, want credential commit failure", err)
	}
	if _, err = source.Token(context.Background()); !errors.Is(err, helix.ErrCredentialCommit) {
		t.Fatalf("terminal Token error = %v, want credential commit failure", err)
	}
	if got := server.refreshes.Load(); got != 1 {
		t.Fatalf("remote refreshes = %d, want no old-token retry", got)
	}
	requireNoActivation(t, source, "old-access")
}
