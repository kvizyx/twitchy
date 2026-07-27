package helix

import (
	"context"
	"time"
)

const (
	defaultRateLimitMaxWait = time.Minute
	minimumRateLimitMaxWait = time.Second
	maximumRateLimitMaxWait = 10 * time.Minute
)

type Clock interface {
	Now() time.Time
}

type Sleeper interface {
	Sleep(context.Context, time.Duration) error
}

type wallClock struct{}

func (wallClock) Now() time.Time { return time.Now() }

type timerSleeper struct{}

func (timerSleeper) Sleep(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (policy RateLimitPolicy) maxWait() time.Duration {
	if policy.MaxWait == 0 {
		return defaultRateLimitMaxWait
	}
	if policy.MaxWait < minimumRateLimitMaxWait {
		return minimumRateLimitMaxWait
	}
	if policy.MaxWait > maximumRateLimitMaxWait {
		return maximumRateLimitMaxWait
	}
	return policy.MaxWait
}
