package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	defaultRefreshLeaseTTL     = 30 * time.Second
	defaultRefreshLeaseRenewal = 10 * time.Second
	defaultRedisIOTimeout      = time.Second
	acquireRetryFloor          = 10 * time.Millisecond
	acquireRetryCeiling        = 50 * time.Millisecond
)

// RefreshLeaseKeyBuilder returns the Redis key for a user ID. Keys must contain
// user identifiers only and must never include access or refresh token material.
type RefreshLeaseKeyBuilder func(userID string) string

type RedisRefreshCoordinatorOption func(*redisRefreshCoordinatorOptions) error

type redisRefreshCoordinatorOptions struct {
	ttl       time.Duration
	renewal   time.Duration
	ioTimeout time.Duration
}

// WithRefreshLeaseTTL sets the Redis lease duration.
func WithRefreshLeaseTTL(ttl time.Duration) RedisRefreshCoordinatorOption {
	return func(options *redisRefreshCoordinatorOptions) error {
		if ttl <= 0 {
			return ErrInvalidOption
		}
		options.ttl = ttl
		return nil
	}
}

// WithRefreshLeaseRenewal sets the automatic lease renewal interval.
func WithRefreshLeaseRenewal(renewal time.Duration) RedisRefreshCoordinatorOption {
	return func(options *redisRefreshCoordinatorOptions) error {
		if renewal <= 0 {
			return ErrInvalidOption
		}
		options.renewal = renewal
		return nil
	}
}

// WithRedisIOTimeout bounds lease renewal and release socket I/O. The
// coordinator uses a lease-scoped client that shares the supplied client's pool
// and hooks; it never closes or changes the supplied client.
func WithRedisIOTimeout(timeout time.Duration) RedisRefreshCoordinatorOption {
	return func(options *redisRefreshCoordinatorOptions) error {
		if timeout <= 0 {
			return ErrInvalidOption
		}
		options.ioTimeout = timeout
		return nil
	}
}

// RedisRefreshCoordinator coordinates one refresh lease per user through Redis.
type RedisRefreshCoordinator struct {
	client  *redis.Client
	key     RefreshLeaseKeyBuilder
	ttl     time.Duration
	renewal time.Duration
}

// NewRedisRefreshCoordinator derives a bounded lease client from client. The
// derived client shares the supplied client's pool and installed hooks, and is
// not closed by the coordinator.
func NewRedisRefreshCoordinator(
	client *redis.Client,
	key RefreshLeaseKeyBuilder,
	options ...RedisRefreshCoordinatorOption,
) (*RedisRefreshCoordinator, error) {
	if client == nil || key == nil {
		return nil, ErrInvalidOption
	}
	configuration := redisRefreshCoordinatorOptions{
		ttl:     defaultRefreshLeaseTTL,
		renewal: defaultRefreshLeaseRenewal,
	}
	for _, option := range options {
		if option == nil {
			return nil, ErrInvalidOption
		}
		if err := option(&configuration); err != nil {
			return nil, err
		}
	}
	if configuration.renewal > configuration.ttl/3 {
		return nil, ErrInvalidOption
	}
	if configuration.ioTimeout == 0 {
		configuration.ioTimeout = min(defaultRedisIOTimeout, configuration.ttl/3)
	}
	if configuration.ioTimeout <= 0 || configuration.ioTimeout > configuration.ttl/3 {
		return nil, ErrInvalidOption
	}
	leaseClient := client.WithTimeout(configuration.ioTimeout)
	leaseClient.Options().ContextTimeoutEnabled = true
	return &RedisRefreshCoordinator{
		client:  leaseClient,
		key:     key,
		ttl:     configuration.ttl,
		renewal: configuration.renewal,
	}, nil
}

func (coordinator *RedisRefreshCoordinator) Acquire(ctx context.Context, userID string) (RefreshLease, error) {
	if ctx == nil || userID == "" {
		return nil, ErrInvalidOption
	}
	key := coordinator.key(userID)
	if key == "" {
		return nil, ErrInvalidOption
	}
	owner, err := newRefreshLeaseOwner()
	if err != nil {
		return nil, fmt.Errorf("%w: generate lease owner: %w", ErrRefreshCoordinator, err)
	}
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		acquired, err := coordinator.client.SetNX(ctx, key, owner, coordinator.ttl).Result()
		if err != nil {
			return nil, fmt.Errorf("%w: acquire lease: %w", ErrRefreshCoordinator, err)
		}
		if acquired {
			return newRedisRefreshLease(ctx, coordinator, key, owner), nil
		}
		delay, err := refreshAcquireDelay()
		if err != nil {
			return nil, fmt.Errorf("%w: jitter acquisition delay: %w", ErrRefreshCoordinator, err)
		}
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return nil, ctx.Err()
		case <-timer.C:
		}
	}
}

func newRefreshLeaseOwner() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func refreshAcquireDelay() (time.Duration, error) {
	span := int64(acquireRetryCeiling - acquireRetryFloor)
	random, err := rand.Int(rand.Reader, big.NewInt(span+1))
	if err != nil {
		return 0, err
	}
	return acquireRetryFloor + time.Duration(random.Int64()), nil
}
