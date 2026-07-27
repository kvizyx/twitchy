package testkit

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

type Barrier struct {
	mu        sync.Mutex
	remaining int
	ready     chan struct{}
	once      sync.Once
}

func NewBarrier(participants int) (*Barrier, error) {
	if participants < 1 {
		return nil, fmt.Errorf("testkit: barrier participants must be positive")
	}
	return &Barrier{remaining: participants, ready: make(chan struct{})}, nil
}

func (barrier *Barrier) Wait(ctx context.Context) error {
	barrier.mu.Lock()
	if barrier.remaining > 0 {
		barrier.remaining--
		if barrier.remaining == 0 {
			barrier.once.Do(func() { close(barrier.ready) })
		}
	}
	barrier.mu.Unlock()
	select {
	case <-barrier.ready:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

type AtomicCounter struct{ value atomic.Int64 }

func (counter *AtomicCounter) Add(delta int64) int64 { return counter.value.Add(delta) }
func (counter *AtomicCounter) Load() int64           { return counter.value.Load() }
func (counter *AtomicCounter) Reset()                { counter.value.Store(0) }

type GoroutineSnapshot struct{ Count int }

func SnapshotGoroutines() GoroutineSnapshot { return GoroutineSnapshot{Count: runtime.NumGoroutine()} }

func StabilizedGoroutineCount() int {
	previous := runtime.NumGoroutine()
	for range 100 {
		runtime.Gosched()
		current := runtime.NumGoroutine()
		if current == previous {
			return current
		}
		previous = current
	}
	return previous
}

func CheckGoroutineLeak(before GoroutineSnapshot) error {
	after := StabilizedGoroutineCount()
	if after > before.Count {
		return fmt.Errorf("testkit: goroutine count increased from %d to %d", before.Count, after)
	}
	return nil
}
