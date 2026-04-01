package messagetracker

import (
	"time"
)

// InMemoryOption is an option for InMemoryMessageTracker.
type InMemoryOption func(*inMemoryOptions)

type inMemoryOptions struct {
	messageTTL time.Duration
}

func buildDefaultInMemoryOptions(opts ...InMemoryOption) inMemoryOptions {
	o := inMemoryOptions{
		messageTTL: SafeMessageTTL,
	}

	for _, applyOption := range opts {
		applyOption(&o)
	}

	return o
}

// InMemoryWithMessageTTL sets non-default TTL for tracked messages in the storage. You almost never want to use this
// option.
//
// Default value is SafeMessageTTL.
func InMemoryWithMessageTTL(ttl time.Duration) InMemoryOption {
	return func(opts *inMemoryOptions) {
		opts.messageTTL = ttl
	}
}
