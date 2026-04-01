package messagetracker

import (
	"context"

	"github.com/kvizyx/twitchy/internal/shardedset"
)

// InMemoryMessageTracker is an in-memory concurrent-safe implementation of MessageTracker based on
// shardedset.ShardedSet.
//
// This implementation is suitable for simple or test use-cases when you're running only one instance of your
// application.
type InMemoryMessageTracker struct {
	messages shardedset.ShardedSet[string]
}

var _ MessageTracker = (*InMemoryMessageTracker)(nil)

// NewInMemoryMessageTracker creates a new InMemoryMessageTracker with the provided options.
func NewInMemoryMessageTracker(options ...InMemoryOption) *InMemoryMessageTracker {
	opts := buildDefaultInMemoryOptions(options...)

	return &InMemoryMessageTracker{
		messages: shardedset.NewString(opts.messageTTL),
	}
}

func (imt *InMemoryMessageTracker) Track(_ context.Context, messageID string) (bool, error) {
	isDuplicate := imt.messages.SetIfAbsent(messageID)
	return isDuplicate, nil
}

// Stop stops the background eviction of expired messages in storage.
//
// The InMemoryMessageTracker must not be used after Stop is called.
func (imt *InMemoryMessageTracker) Stop() {
	imt.messages.Stop()
}
