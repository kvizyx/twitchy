package oauth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"
)

type reauthorizationRequiredError interface {
	error
	RequiresReauthorization() bool
}

func TestCoordinatedLifecycle_leaseLossBeforeHTTPDoesNotRefreshCommitOrActivate(t *testing.T) {
	clock := newTestClock(time.Unix(100_000, 0))
	server := newCoordinatedOAuthServer(t)
	coordinator := newMemoryRefreshCoordinator()
	store := newCoordinatedTestStore(TokenPair{
		AccessToken:  "unchanged-access",
		RefreshToken: "unchanged-refresh",
		ExpiresIn:    time.Hour,
	})
	source := newCoordinatedSource(t, clock, server, coordinator, store)
	drainCoordinatedRegistration(t, coordinator)
	loadStarted := make(chan struct{})
	store.mu.Lock()
	store.loadFn = func(ctx context.Context, _ string) (TokenPair, error) {
		close(loadStarted)
		<-ctx.Done()
		return TokenPair{}, ctx.Err()
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
	lease.lose(ErrRefreshLeaseLost)
	err := waitLifecycleResult(t, result)
	if !errors.Is(err, ErrRefreshLeaseLost) {
		t.Fatalf("refresh error = %v, want lease loss", err)
	}
	if got := server.refreshes.Load(); got != 0 {
		t.Fatalf("remote refreshes = %d, want 0", got)
	}
	if got := store.hooks.Load(); got != 0 {
		t.Fatalf("durable commits = %d, want 0", got)
	}
	requireNoActivation(t, source, "unchanged-access")
}

func TestCoordinatedLifecycle_joinsReloadAndReleaseErrors(t *testing.T) {
	clock := newTestClock(time.Unix(105_000, 0))
	server := newCoordinatedOAuthServer(t)
	coordinator := newMemoryRefreshCoordinator()
	store := newCoordinatedTestStore(TokenPair{
		AccessToken:  "unchanged-access",
		RefreshToken: "unchanged-refresh",
		ExpiresIn:    time.Hour,
	})
	source := newCoordinatedSource(t, clock, server, coordinator, store)
	drainCoordinatedRegistration(t, coordinator)
	loadStarted := make(chan struct{})
	releaseLoad := make(chan struct{})
	reloadErr := errors.New("reload unavailable")
	releaseErr := errors.New("release unavailable")
	store.mu.Lock()
	store.loadFn = func(ctx context.Context, _ string) (TokenPair, error) {
		close(loadStarted)
		select {
		case <-releaseLoad:
			return TokenPair{}, reloadErr
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
	lease.releaseErr = releaseErr
	lease.mu.Unlock()
	close(releaseLoad)
	err := waitLifecycleResult(t, result)
	for _, want := range []error{reloadErr, releaseErr} {
		if !errors.Is(err, want) {
			t.Fatalf("refresh error = %v, missing %v", err, want)
		}
	}
}

func TestCoordinatedLifecycle_leaseLossDuringHTTPReturnsSanitizedTerminalError(t *testing.T) {
	clock := newTestClock(time.Unix(110_000, 0))
	server := newCoordinatedOAuthServer(t)
	server.setRefreshFn(func(_ http.ResponseWriter, request *http.Request, _ int32) {
		<-request.Context().Done()
	})
	coordinator := newMemoryRefreshCoordinator()
	store := newCoordinatedTestStore(TokenPair{
		AccessToken:  "private-access-value",
		RefreshToken: "private-refresh-value",
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
	lease := coordinator.nextLease(t)
	select {
	case <-server.refreshStarted:
	case <-time.After(time.Second):
		t.Fatal("remote refresh did not start")
	}
	lease.lose(ErrRefreshLeaseLost)
	err := waitLifecycleResult(t, result)
	if !errors.Is(err, ErrRefreshLeaseLost) {
		t.Fatalf("refresh error = %v, want lease loss", err)
	}
	var terminal reauthorizationRequiredError
	if !errors.As(err, &terminal) || !terminal.RequiresReauthorization() {
		t.Fatalf("refresh error type = %T, want reauthorization terminal", err)
	}
	for _, secret := range []string{"private-access-value", "private-refresh-value"} {
		if strings.Contains(err.Error(), secret) {
			t.Fatalf("refresh error exposed credential material")
		}
	}
	if got := store.hooks.Load(); got != 0 {
		t.Fatalf("durable commits = %d, want 0", got)
	}
	requireNoActivation(t, source, "private-access-value")
	if _, tokenErr := source.Token(context.Background()); !errors.As(tokenErr, &terminal) {
		t.Fatalf("next Token error = %v, want terminal rotation error", tokenErr)
	}
}
