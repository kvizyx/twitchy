package messagetracker

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type (
	// RedisKeyBuilder builds and returns key for Redis to store event with provided event identifier.
	RedisKeyBuilder func(messageID string) string

	// RedisAnyClient is an alias for redis.Cmdable that represents any Redis client that can be used by
	// RedisMessageTracker.
	RedisAnyClient = redis.Cmdable

	// RedisMessageTracker is a Redis implementation of MessageTracker based on official Redis client.
	// It accepts any redis.Cmdable, so it works with *redis.Client, *redis.ClusterClient, *redis.Ring, and others.
	RedisMessageTracker struct {
		client RedisAnyClient
		opts   redisOptions
	}
)

var _ MessageTracker = (*RedisMessageTracker)(nil)

// NewRedisMessageTracker creates a new RedisMessageTracker with the provided options.
//
// The client parameter accepts any RedisAnyClient implementation (*redis.Client, *redis.ClusterClient, etc.).
func NewRedisMessageTracker(client RedisAnyClient, options ...RedisOption) RedisMessageTracker {
	return RedisMessageTracker{
		client: client,
		opts:   buildDefaultRedisOptions(options...),
	}
}

func (rmt RedisMessageTracker) Track(ctx context.Context, messageID string) (bool, error) {
	key := rmt.opts.keyBuilder(messageID)

	firstTimeSeen, err := rmt.client.SetNX(ctx, key, 1, rmt.opts.messageTTL).Result()
	if err != nil {
		return false, fmt.Errorf("set nx: %w", err)
	}

	isDuplicate := !firstTimeSeen
	return isDuplicate, nil
}
