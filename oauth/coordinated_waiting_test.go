package oauth

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/kvizyx/twitchy/helix"
)

func TestCoordinatedLifecycle_callerCancellationWhileWaitingSkipsSecondRefresh(t *testing.T) {
	testCoordinatedWaitingCancellation(t, false)
}

func TestCoordinatedLifecycle_closeWhileWaitingSkipsSecondRefresh(t *testing.T) {
	testCoordinatedWaitingCancellation(t, true)
}

func testCoordinatedWaitingCancellation(t *testing.T, closeSource bool) {
	t.Helper()
	clock := newTestClock(time.Unix(130_000, 0))
	server := newCoordinatedOAuthServer(t)
	releaseFirst := make(chan struct{})
	server.setRefreshFn(func(writer http.ResponseWriter, request *http.Request, count int32) {
		if count == 1 {
			select {
			case <-releaseFirst:
			case <-request.Context().Done():
				return
			}
		}
		_, _ = io.WriteString(writer, coordinatedRotationResponse)
	})
	coordinator := newMemoryRefreshCoordinator()
	store := newCoordinatedTestStore(TokenPair{
		AccessToken:  "old-access",
		RefreshToken: "old-refresh",
		ExpiresIn:    time.Hour,
	})
	first := newCoordinatedSource(t, clock, server, coordinator, store)
	second := newCoordinatedSource(t, clock, server, coordinator, store)
	drainCoordinatedRegistration(t, coordinator)
	drainCoordinatedRegistration(t, coordinator)
	clock.Advance(time.Hour - defaultRefreshSkew)

	firstResult := make(chan error, 1)
	go func() {
		_, err := first.Token(context.Background())
		firstResult <- err
	}()
	waitAcquireAttempt(t, coordinator)
	_ = coordinator.nextLease(t)
	select {
	case <-server.refreshStarted:
	case <-time.After(time.Second):
		t.Fatal("first remote refresh did not start")
	}

	waitContext, cancel := context.WithCancel(context.Background())
	defer cancel()
	secondResult := make(chan error, 1)
	go func() {
		_, err := second.Token(waitContext)
		secondResult <- err
	}()
	select {
	case <-coordinator.attempted:
	case err := <-secondResult:
		t.Fatalf("waiting refresh returned before lease acquisition: %v", err)
	case <-time.After(time.Second):
		t.Fatal("waiting refresh did not attempt lease acquisition")
	}
	if closeSource {
		if err := second.Close(); err != nil {
			t.Fatal(err)
		}
	} else {
		cancel()
	}
	err := waitLifecycleResult(t, secondResult)
	if closeSource {
		if !errors.Is(err, helix.ErrSessionClosed) {
			t.Fatalf("waiting refresh error = %v, want session closed", err)
		}
	} else if !errors.Is(err, context.Canceled) {
		t.Fatalf("waiting refresh error = %v, want caller cancellation", err)
	}
	close(releaseFirst)
	if err = waitLifecycleResult(t, firstResult); err != nil {
		t.Fatalf("first refresh error = %v", err)
	}
	if got := server.refreshes.Load(); got != 1 {
		t.Fatalf("remote refreshes = %d, want 1", got)
	}
}

func drainCoordinatedRegistration(t *testing.T, coordinator *memoryRefreshCoordinator) {
	t.Helper()
	waitAcquireAttempt(t, coordinator)
	_ = coordinator.nextLease(t)
}

func waitAcquireAttempt(t *testing.T, coordinator *memoryRefreshCoordinator) {
	t.Helper()
	select {
	case <-coordinator.attempted:
	case <-time.After(time.Second):
		t.Fatal("coordinated lease acquisition was not attempted")
	}
}

func waitLifecycleResult(t *testing.T, result <-chan error) error {
	t.Helper()
	select {
	case err := <-result:
		return err
	case <-time.After(time.Second):
		t.Fatal("coordinated operation did not finish")
		return nil
	}
}
