package helix

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/kvizyx/twitchy/helix/internal/manifest"
)

const maxHTTPAttempts = 3

type retryCause uint8

const (
	retryUnauthorized retryCause = iota
	retryRateLimit
	retryUnavailable
)

type retryState struct {
	seen map[retryCause]bool
}

func (state *retryState) canReplay(cause retryCause, attempt int) bool {
	if attempt >= maxHTTPAttempts || state.seen[cause] {
		return false
	}
	state.seen[cause] = true
	return true
}

func requestReplaySafe(request *http.Request, operation manifest.Operation) bool {
	if !operation.Replay.Replayable || !operation.Request.BodyReconstructible {
		return false
	}
	return request.Body == nil || request.Body == http.NoBody || request.GetBody != nil
}

func replayRequest(ctx context.Context, request *http.Request, attempt int, credential CredentialSnapshot) (*http.Request, error) {
	clone := request.Clone(ctx)
	if attempt > 1 && request.Body != nil && request.Body != http.NoBody {
		if request.GetBody == nil {
			return nil, fmt.Errorf("helix: request body cannot be replayed")
		}
		body, err := request.GetBody()
		if err != nil {
			return nil, fmt.Errorf("helix: reconstruct request body: %w", err)
		}
		clone.Body = body
	}
	if credential.AccessToken() != "" {
		clone.Header.Set("Authorization", "Bearer "+credential.AccessToken())
	}
	if credential.ClientID() != "" {
		clone.Header.Set("Client-Id", credential.ClientID())
	}
	return clone, nil
}

type refreshFlight struct {
	done     chan struct{}
	snapshot CredentialSnapshot
	err      error
}

type refreshCoordinator struct {
	mu      sync.Mutex
	flights map[uint64]*refreshFlight
	latest  map[uint64]CredentialSnapshot
}

func newRefreshCoordinator() *refreshCoordinator {
	return &refreshCoordinator{
		flights: make(map[uint64]*refreshFlight),
		latest:  make(map[uint64]CredentialSnapshot),
	}
}

func (coordinator *refreshCoordinator) refresh(ctx context.Context, source RefreshableTokenSource, snapshot CredentialSnapshot) (CredentialSnapshot, error) {
	coordinator.mu.Lock()
	if refreshed, exists := coordinator.latest[snapshot.Generation()]; exists {
		coordinator.mu.Unlock()
		return refreshed, nil
	}
	flight, exists := coordinator.flights[snapshot.Generation()]
	if !exists {
		flight = &refreshFlight{done: make(chan struct{})}
		coordinator.flights[snapshot.Generation()] = flight
	}
	coordinator.mu.Unlock()
	if exists {
		select {
		case <-flight.done:
			return flight.snapshot, flight.err
		case <-ctx.Done():
			return CredentialSnapshot{}, ctx.Err()
		}
	}

	refreshed, err := source.Refresh(ctx, snapshot, RefreshReasonUnauthorized)
	coordinator.mu.Lock()
	flight.snapshot = refreshed
	flight.err = err
	if err == nil {
		coordinator.latest[snapshot.Generation()] = refreshed
	}
	delete(coordinator.flights, snapshot.Generation())
	close(flight.done)
	coordinator.mu.Unlock()
	return refreshed, err
}
