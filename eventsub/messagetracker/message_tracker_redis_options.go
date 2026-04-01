package messagetracker

import (
	"time"
)

// RedisOption is an option for RedisMessageTracker.
type RedisOption func(*redisOptions)

type redisOptions struct {
	messageTTL time.Duration
	keyBuilder RedisKeyBuilder
}

func buildDefaultRedisOptions(opts ...RedisOption) redisOptions {
	o := redisOptions{
		messageTTL: SafeMessageTTL,
		keyBuilder: func(messageID string) string {
			return "twitch-eventsub-messages:" + messageID
		},
	}

	for _, applyOption := range opts {
		applyOption(&o)
	}

	return o
}

// RedisWithMessageTTL sets non-default TTL for tracked messages in the storage. You almost never want to use this
// option.
//
// Default value is SafeMessageTTL.
func RedisWithMessageTTL(ttl time.Duration) RedisOption {
	return func(opts *redisOptions) {
		opts.messageTTL = ttl
	}
}

// RedisWithKeyBuilder sets the key builder used to construct the Redis key for each tracked message.
//
// If not provided, default prefix will be used.
func RedisWithKeyBuilder(kb RedisKeyBuilder) RedisOption {
	return func(opts *redisOptions) {
		opts.keyBuilder = kb
	}
}
