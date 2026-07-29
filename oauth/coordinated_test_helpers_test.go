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

type memoryRefreshCoordinator struct {
	mu           sync.Mutex
	held         map[string]bool
	released     map[string]chan struct{}
	acquired     chan *memoryRefreshLease
	attempted    chan string
	acquireCalls atomic.Int32
}

func newMemoryRefreshCoordinator() *memoryRefreshCoordinator {
	return &memoryRefreshCoordinator{
		held:      make(map[string]bool),
		released:  make(map[string]chan struct{}),
		acquired:  make(chan *memoryRefreshLease, 64),
		attempted: make(chan string, 64),
	}
}

func (coordinator *memoryRefreshCoordinator) Acquire(ctx context.Context, userID string) (RefreshLease, error) {
	coordinator.acquireCalls.Add(1)
	coordinator.attempted <- userID
	for {
		coordinator.mu.Lock()
		if !coordinator.held[userID] {
			coordinator.held[userID] = true
			released := make(chan struct{})
			coordinator.released[userID] = released
			leaseContext, cancel := context.WithCancel(ctx)
			lease := &memoryRefreshLease{
				ctx:    leaseContext,
				cancel: cancel,
				releaseOwner: func() {
					coordinator.mu.Lock()
					if coordinator.held[userID] && coordinator.released[userID] == released {
						coordinator.held[userID] = false
						close(released)
					}
					coordinator.mu.Unlock()
				},
			}
			coordinator.mu.Unlock()
			coordinator.acquired <- lease
			return lease, nil
		}
		released := coordinator.released[userID]
		coordinator.mu.Unlock()
		select {
		case <-released:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (coordinator *memoryRefreshCoordinator) nextLease(t *testing.T) *memoryRefreshLease {
	t.Helper()
	select {
	case lease := <-coordinator.acquired:
		return lease
	case <-time.After(time.Second):
		t.Fatal("coordinated lease was not acquired")
		return nil
	}
}

type memoryRefreshLease struct {
	ctx          context.Context
	cancel       context.CancelFunc
	releaseOwner func()

	mu          sync.Mutex
	lossErr     error
	releaseErr  error
	errCalls    int
	loseOnCall  int
	releaseOnce sync.Once
}

func (lease *memoryRefreshLease) Context() context.Context { return lease.ctx }

func (lease *memoryRefreshLease) AssertOwnership(context.Context) error { return lease.Err() }

func (lease *memoryRefreshLease) Err() error {
	lease.mu.Lock()
	defer lease.mu.Unlock()
	lease.errCalls++
	if lease.loseOnCall > 0 && lease.errCalls == lease.loseOnCall && lease.lossErr == nil {
		lease.lossErr = ErrRefreshLeaseLost
		lease.cancel()
	}
	if lease.lossErr != nil {
		return lease.lossErr
	}
	return lease.ctx.Err()
}

func (lease *memoryRefreshLease) Release(context.Context) error {
	lease.releaseOnce.Do(func() {
		lease.releaseOwner()
		lease.cancel()
	})
	lease.mu.Lock()
	defer lease.mu.Unlock()
	return errors.Join(lease.lossErr, lease.releaseErr)
}

func (lease *memoryRefreshLease) lose(err error) {
	lease.mu.Lock()
	if lease.lossErr == nil {
		lease.lossErr = err
	}
	lease.mu.Unlock()
	lease.cancel()
}

type coordinatedTestStore struct {
	mu       sync.Mutex
	pair     TokenPair
	loadErr  error
	loadFn   func(context.Context, string) (TokenPair, error)
	hookFn   func(context.Context, TokenPair) error
	lastHook TokenPair
	loads    atomic.Int32
	hooks    atomic.Int32
}

func newCoordinatedTestStore(pair TokenPair) *coordinatedTestStore {
	return &coordinatedTestStore{pair: cloneTokenPair(pair)}
}

func (store *coordinatedTestStore) load(ctx context.Context, userID string) (TokenPair, error) {
	store.loads.Add(1)
	if err := ctx.Err(); err != nil {
		return TokenPair{}, err
	}
	store.mu.Lock()
	loadFn := store.loadFn
	pair := cloneTokenPair(store.pair)
	loadErr := store.loadErr
	store.mu.Unlock()
	if loadFn != nil {
		return loadFn(ctx, userID)
	}
	return pair, loadErr
}

func (store *coordinatedTestStore) hook(ctx context.Context, pair TokenPair) error {
	store.hooks.Add(1)
	store.mu.Lock()
	store.lastHook = cloneTokenPair(pair)
	hookFn := store.hookFn
	store.mu.Unlock()
	if hookFn != nil {
		return hookFn(ctx, cloneTokenPair(pair))
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	store.setPair(pair)
	return nil
}

func (store *coordinatedTestStore) setPair(pair TokenPair) {
	store.mu.Lock()
	store.pair = cloneTokenPair(pair)
	store.mu.Unlock()
}

func (store *coordinatedTestStore) currentPair() TokenPair {
	store.mu.Lock()
	defer store.mu.Unlock()
	return cloneTokenPair(store.pair)
}

func (store *coordinatedTestStore) pendingPair() TokenPair {
	store.mu.Lock()
	defer store.mu.Unlock()
	return cloneTokenPair(store.lastHook)
}

func cloneTokenPair(pair TokenPair) TokenPair {
	pair.Scopes = append([]helix.AuthorizationScope(nil), pair.Scopes...)
	return pair
}
