package testkit

import (
	"context"
	"sync"
	"time"

	"github.com/kvizyx/twitchy/helix"
)

type FakeClock struct {
	mu  sync.RWMutex
	now time.Time
}

func NewFakeClock(now time.Time) *FakeClock { return &FakeClock{now: now} }

func (clock *FakeClock) Now() time.Time {
	clock.mu.RLock()
	defer clock.mu.RUnlock()
	return clock.now
}

func (clock *FakeClock) NewTimer(duration time.Duration) helix.Timer {
	return fakeTimer{timer: time.NewTimer(duration)}
}

type fakeTimer struct{ timer *time.Timer }

func (timer fakeTimer) C() <-chan time.Time { return timer.timer.C }
func (timer fakeTimer) Stop() bool          { return timer.timer.Stop() }

func (clock *FakeClock) Set(now time.Time) {
	clock.mu.Lock()
	clock.now = now
	clock.mu.Unlock()
}

func (clock *FakeClock) Advance(duration time.Duration) {
	clock.mu.Lock()
	clock.now = clock.now.Add(duration)
	clock.mu.Unlock()
}

type FakeSleeper struct {
	clock     *FakeClock
	mu        sync.Mutex
	durations []time.Duration
}

func NewFakeSleeper(clock *FakeClock) *FakeSleeper { return &FakeSleeper{clock: clock} }

func (sleeper *FakeSleeper) Sleep(ctx context.Context, duration time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	sleeper.mu.Lock()
	sleeper.durations = append(sleeper.durations, duration)
	sleeper.mu.Unlock()
	if sleeper.clock != nil {
		sleeper.clock.Advance(duration)
	}
	return nil
}

func (sleeper *FakeSleeper) Durations() []time.Duration {
	sleeper.mu.Lock()
	defer sleeper.mu.Unlock()
	return append([]time.Duration(nil), sleeper.durations...)
}

var _ helix.Clock = (*FakeClock)(nil)
var _ helix.Sleeper = (*FakeSleeper)(nil)
