package oauth

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	renewRefreshLease = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
  return redis.call("PEXPIRE", KEYS[1], ARGV[2])
end
return 0`)
	releaseRefreshLease = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
  return redis.call("DEL", KEYS[1])
end
return 0`)
)

type redisRefreshLease struct {
	coordinator *RedisRefreshCoordinator
	key         string
	owner       string
	ctx         context.Context
	cancel      context.CancelFunc
	stop        chan struct{}
	done        chan struct{}

	mu         sync.Mutex
	lost       bool
	lossErr    error
	releaseErr error
	release    sync.Once
}

func newRedisRefreshLease(
	ctx context.Context,
	coordinator *RedisRefreshCoordinator,
	key string,
	owner string,
) *redisRefreshLease {
	leaseContext, cancel := context.WithCancel(ctx)
	lease := &redisRefreshLease{
		coordinator: coordinator,
		key:         key,
		owner:       owner,
		ctx:         leaseContext,
		cancel:      cancel,
		stop:        make(chan struct{}),
		done:        make(chan struct{}),
	}
	go lease.renew()
	return lease
}

func (lease *redisRefreshLease) Context() context.Context {
	return lease.ctx
}

func (lease *redisRefreshLease) Err() error {
	lease.mu.Lock()
	defer lease.mu.Unlock()
	if lease.lossErr != nil {
		return lease.lossErr
	}
	return lease.ctx.Err()
}

func (lease *redisRefreshLease) AssertOwnership(ctx context.Context) error {
	if ctx == nil {
		return ErrInvalidOption
	}
	if err := lease.Err(); err != nil {
		return err
	}
	owner, err := lease.coordinator.client.Get(ctx, lease.key).Result()
	if errors.Is(err, redis.Nil) {
		lease.markLost(ErrRefreshLeaseLost)
		return ErrRefreshLeaseLost
	}
	if err != nil {
		return fmt.Errorf("%w: assert lease ownership: %w", ErrRefreshCoordinator, err)
	}
	if owner != lease.owner {
		lease.markLost(ErrRefreshLeaseLost)
		return ErrRefreshLeaseLost
	}
	return nil
}

func (lease *redisRefreshLease) Release(ctx context.Context) error {
	if ctx == nil {
		return ErrInvalidOption
	}
	lease.release.Do(func() {
		close(lease.stop)
		lease.cancel()
		select {
		case <-lease.done:
		case <-ctx.Done():
			lease.releaseErr = errors.Join(
				fmt.Errorf("%w: wait for renewal: %w", ErrRefreshCoordinator, ctx.Err()),
				lease.Err(),
			)
			return
		}
		lease.mu.Lock()
		lost := lease.lost
		lossErr := lease.lossErr
		lease.mu.Unlock()
		if lost {
			lease.releaseErr = lossErr
			return
		}
		deleted, err := releaseRefreshLease.Run(ctx, lease.coordinator.client, []string{lease.key}, lease.owner).Int64()
		if err != nil {
			lease.releaseErr = fmt.Errorf("%w: release lease: %w", ErrRefreshCoordinator, err)
			lease.markLost(lease.releaseErr)
			return
		}
		if deleted != 1 {
			lease.releaseErr = ErrRefreshLeaseLost
			lease.markLost(lease.releaseErr)
			return
		}
	})
	return lease.releaseErr
}

func (lease *redisRefreshLease) renew() {
	ticker := time.NewTicker(lease.coordinator.renewal)
	defer ticker.Stop()
	defer close(lease.done)
	for {
		select {
		case <-lease.stop:
			return
		case <-lease.ctx.Done():
			return
		case <-ticker.C:
			renewed, err := renewRefreshLease.Run(
				lease.ctx,
				lease.coordinator.client,
				[]string{lease.key},
				lease.owner,
				lease.coordinator.ttl.Milliseconds(),
			).Int64()
			if err != nil {
				if lease.ctx.Err() == nil {
					lease.markLost(fmt.Errorf("%w: renew lease: %w", ErrRefreshCoordinator, err))
				}
				return
			}
			if renewed != 1 {
				lease.markLost(ErrRefreshLeaseLost)
				return
			}
		}
	}
}

func (lease *redisRefreshLease) markLost(err error) {
	lease.mu.Lock()
	if !lease.lost {
		lease.lost = true
		lease.lossErr = err
	}
	lease.mu.Unlock()
	lease.cancel()
}
