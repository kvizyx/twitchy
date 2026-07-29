package oauth

import (
	"context"
	"errors"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"
)

func TestCoordinatedRefresh_realRedisSerializesSources(t *testing.T) {
	fixture := newRedisCoordinatorFixture(t)
	clock := newTestClock(time.Unix(170_000, 0))
	server := newCoordinatedOAuthServer(t)
	store := newCoordinatedTestStore(TokenPair{
		AccessToken:  "old-access",
		RefreshToken: "old-refresh",
		ExpiresIn:    time.Hour,
	})
	first := newCoordinatedSource(t, clock, server, fixture.coordinator, store)
	second := newCoordinatedSource(t, clock, server, fixture.coordinator, store)
	clock.Advance(time.Hour - defaultRefreshSkew)

	results := make(chan error, 2)
	var callers sync.WaitGroup
	callers.Add(2)
	for _, source := range []*RefreshingTokenSource{first, second} {
		go func() {
			defer callers.Done()
			_, err := source.Token(context.Background())
			results <- err
		}()
	}
	callers.Wait()
	close(results)
	for err := range results {
		if err != nil {
			t.Fatal(err)
		}
	}
	if got := server.refreshes.Load(); got != 1 {
		t.Fatalf("remote refreshes = %d, want 1", got)
	}
	if got := store.hooks.Load(); got != 1 {
		t.Fatalf("durable commits = %d, want 1", got)
	}
	t.Logf("REMOTE_REFRESHES=%d DURABLE_COMMITS=%d ACTIVE_SOURCES=2", server.refreshes.Load(), store.hooks.Load())
}

func TestCoordinatedLifecycle_wrongOwnerReplacementInHookPreventsActivation(t *testing.T) {
	fixture := newRedisCoordinatorFixture(t)
	clock := newTestClock(time.Unix(190_000, 0))
	server := newCoordinatedOAuthServer(t)
	store := newCoordinatedTestStore(TokenPair{
		AccessToken:  "old-access",
		RefreshToken: "old-refresh",
		ExpiresIn:    time.Hour,
	})
	key := fixture.key("42")
	source := newCoordinatedSource(t, clock, server, fixture.coordinator, store)
	store.mu.Lock()
	store.hookFn = func(ctx context.Context, pair TokenPair) error {
		if err := fixture.client.Set(context.Background(), key, "replacement-owner", time.Second).Err(); err != nil {
			return err
		}
		store.setPair(pair)
		return nil
	}
	store.mu.Unlock()
	t.Cleanup(func() { _ = fixture.client.Del(context.Background(), key).Err() })
	clock.Advance(time.Hour - defaultRefreshSkew)

	_, err := source.Token(context.Background())
	if !errors.Is(err, ErrRefreshLeaseLost) {
		t.Fatalf("refresh error = %v, want lease loss", err)
	}
	requireNoActivation(t, source, "old-access")
	if got := store.hooks.Load(); got != 1 {
		t.Fatalf("durable commits = %d, want 1", got)
	}
	t.Logf("REMOTE_REFRESHES=%d DURABLE_COMMITS=%d ACTIVATIONS=0", server.refreshes.Load(), store.hooks.Load())
}

func TestCoordinatedLifecycle_realRedisLossDuringHTTPPreventsCommitAndActivation(t *testing.T) {
	fixture := newRedisCoordinatorFixture(t)
	clock := newTestClock(time.Unix(180_000, 0))
	server := newCoordinatedOAuthServer(t)
	server.setRefreshFn(func(_ http.ResponseWriter, request *http.Request, _ int32) {
		<-request.Context().Done()
	})
	store := newCoordinatedTestStore(TokenPair{
		AccessToken:  "unchanged-access",
		RefreshToken: "unchanged-refresh",
		ExpiresIn:    time.Hour,
	})
	source := newCoordinatedSource(t, clock, server, fixture.coordinator, store)
	clock.Advance(time.Hour - defaultRefreshSkew)

	result := make(chan error, 1)
	go func() {
		_, err := source.Token(context.Background())
		result <- err
	}()
	select {
	case <-server.refreshStarted:
	case <-time.After(time.Second):
		t.Fatal("remote refresh did not start")
	}
	key := fixture.key("42")
	externalLoss := os.Getenv("TWITCHY_TEST_EXTERNAL_REDIS_LOSS") == "1"
	if !externalLoss {
		if err := fixture.client.Set(context.Background(), key, "replacement-owner", time.Second).Err(); err != nil {
			t.Fatal(err)
		}
	}
	t.Cleanup(func() { _ = fixture.client.Del(context.Background(), key).Err() })
	err := waitLifecycleResult(t, result)
	if externalLoss {
		if !errors.Is(err, ErrRefreshCoordinator) {
			t.Fatalf("refresh error = %v, want coordinator loss", err)
		}
	} else if !errors.Is(err, ErrRefreshLeaseLost) {
		t.Fatalf("refresh error = %v, want ownership loss", err)
	}
	var terminal reauthorizationRequiredError
	if !errors.As(err, &terminal) || !terminal.RequiresReauthorization() {
		t.Fatalf("refresh error type = %T, want reauthorization terminal", err)
	}
	if got := store.hooks.Load(); got != 0 {
		t.Fatalf("durable commits = %d, want 0", got)
	}
	requireNoActivation(t, source, "unchanged-access")
	t.Logf("REMOTE_REFRESHES=%d DURABLE_COMMITS=%d ACTIVATIONS=0", server.refreshes.Load(), store.hooks.Load())
}
