package helix

import "time"

type RateLimitPolicy struct {
	Wait    bool
	MaxWait time.Duration
}
