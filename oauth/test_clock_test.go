package oauth

import (
	"sync"
	"time"

	"github.com/kvizyx/twitchy/helix"
)

type testClock struct {
	mu     sync.Mutex
	now    time.Time
	timers map[*testTimer]struct{}
}

func newTestClock(now time.Time) *testClock {
	return &testClock{now: now, timers: make(map[*testTimer]struct{})}
}

func (clock *testClock) Now() time.Time {
	clock.mu.Lock()
	defer clock.mu.Unlock()
	return clock.now
}

func (clock *testClock) NewTimer(duration time.Duration) helix.Timer {
	clock.mu.Lock()
	defer clock.mu.Unlock()
	timer := &testTimer{clock: clock, deadline: clock.now.Add(duration), channel: make(chan time.Time, 1)}
	clock.timers[timer] = struct{}{}
	return timer
}

func (clock *testClock) Advance(duration time.Duration) {
	clock.mu.Lock()
	clock.now = clock.now.Add(duration)
	for timer := range clock.timers {
		if !timer.stopped && !timer.fired && !clock.now.Before(timer.deadline) {
			timer.fired = true
			timer.channel <- clock.now
		}
	}
	clock.mu.Unlock()
}

type testTimer struct {
	clock    *testClock
	deadline time.Time
	channel  chan time.Time
	stopped  bool
	fired    bool
}

func (timer *testTimer) C() <-chan time.Time { return timer.channel }

func (timer *testTimer) Stop() bool {
	timer.clock.mu.Lock()
	defer timer.clock.mu.Unlock()
	wasActive := !timer.stopped && !timer.fired
	timer.stopped = true
	delete(timer.clock.timers, timer)
	return wasActive
}
