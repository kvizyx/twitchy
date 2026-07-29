package oauth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisCoordinatorFixture struct {
	client      *redis.Client
	coordinator *RedisRefreshCoordinator
	key         RefreshLeaseKeyBuilder
}

func newRedisCoordinatorFixture(t *testing.T) redisCoordinatorFixture {
	t.Helper()
	client := redisTestClient(t)
	prefix := fmt.Sprintf("oauth:test:%d:", time.Now().UnixNano())
	key := RefreshLeaseKeyBuilder(func(userID string) string { return prefix + userID })
	coordinator, err := NewRedisRefreshCoordinator(
		client,
		key,
		WithRefreshLeaseTTL(300*time.Millisecond),
		WithRefreshLeaseRenewal(100*time.Millisecond),
	)
	if err != nil {
		t.Fatal(err)
	}
	return redisCoordinatorFixture{client: client, coordinator: coordinator, key: key}
}

func TestRedisRefreshCoordinator(t *testing.T) {
	fixture := newRedisCoordinatorFixture(t)
	t.Run("owner-safe acquire and release", func(t *testing.T) { testOwnerSafeRedisLease(t, fixture) })
	t.Run("renews before expiry", func(t *testing.T) { testRedisLeaseRenewal(t, fixture) })
	t.Run("wrong owner retains key and cancels lease", func(t *testing.T) { testWrongOwnerRedisLease(t, fixture) })
	t.Run("renew loss cancels lease", func(t *testing.T) { testRedisLeaseLoss(t, fixture) })
}

func testOwnerSafeRedisLease(t *testing.T, fixture redisCoordinatorFixture) {
	t.Helper()
	first, err := fixture.coordinator.Acquire(context.Background(), "first")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = first.Release(context.Background()) })

	waitingContext, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	if _, err = fixture.coordinator.Acquire(waitingContext, "first"); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("contended Acquire error = %v, want context deadline exceeded", err)
	}
	if err = first.Release(context.Background()); err != nil {
		t.Fatal(err)
	}
	second, err := fixture.coordinator.Acquire(context.Background(), "first")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = second.Release(context.Background()) }()
}

func testRedisLeaseRenewal(t *testing.T, fixture redisCoordinatorFixture) {
	t.Helper()
	lease, err := fixture.coordinator.Acquire(context.Background(), "renewed")
	if err != nil {
		t.Fatal(err)
	}
	timer := time.NewTimer(500 * time.Millisecond)
	defer timer.Stop()
	select {
	case <-lease.Context().Done():
		t.Fatalf("lease ended before release: %v", lease.Err())
	case <-timer.C:
	}
	if err = lease.Release(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func testWrongOwnerRedisLease(t *testing.T, fixture redisCoordinatorFixture) {
	t.Helper()
	lease, err := fixture.coordinator.Acquire(context.Background(), "wrong-owner")
	if err != nil {
		t.Fatal(err)
	}
	if err = fixture.client.Set(
		context.Background(),
		fixture.key("wrong-owner"),
		"different-owner",
		time.Second,
	).Err(); err != nil {
		t.Fatal(err)
	}
	if err = lease.Release(context.Background()); !errors.Is(err, ErrRefreshLeaseLost) {
		t.Fatalf("wrong-owner Release error = %v, want ErrRefreshLeaseLost", err)
	}
	value, err := fixture.client.Get(context.Background(), fixture.key("wrong-owner")).Result()
	if err != nil {
		t.Fatal(err)
	}
	if value != "different-owner" {
		t.Fatalf("wrong-owner key value = %q", value)
	}
	select {
	case <-lease.Context().Done():
	case <-time.After(time.Second):
		t.Fatal("lease context was not canceled after ownership loss")
	}
}

func testRedisLeaseLoss(t *testing.T, fixture redisCoordinatorFixture) {
	t.Helper()
	lease, err := fixture.coordinator.Acquire(context.Background(), "renew-loss")
	if err != nil {
		t.Fatal(err)
	}
	if err = fixture.client.Set(
		context.Background(),
		fixture.key("renew-loss"),
		"different-owner",
		time.Second,
	).Err(); err != nil {
		t.Fatal(err)
	}
	select {
	case <-lease.Context().Done():
	case <-time.After(time.Second):
		t.Fatal("lease context was not canceled after renewal lost ownership")
	}
	if !errors.Is(lease.Err(), ErrRefreshLeaseLost) {
		t.Fatalf("lease loss error = %v, want ErrRefreshLeaseLost", lease.Err())
	}
	redisLease := lease.(*redisRefreshLease)
	select {
	case <-redisLease.done:
	case <-time.After(time.Second):
		t.Fatal("renew goroutine did not exit after ownership loss")
	}
}

func TestRedisRefreshCoordinatorContracts(t *testing.T) {
	key := func(string) string { return "lease" }
	if _, err := NewRedisRefreshCoordinator(nil, key); !errors.Is(err, ErrInvalidOption) {
		t.Fatalf("nil client error = %v, want ErrInvalidOption", err)
	}
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:1", DialTimeout: 10 * time.Millisecond})
	t.Cleanup(func() { _ = client.Close() })
	if _, err := NewRedisRefreshCoordinator(client, nil); !errors.Is(err, ErrInvalidOption) {
		t.Fatalf("nil key builder error = %v, want ErrInvalidOption", err)
	}
	for _, options := range [][]RedisRefreshCoordinatorOption{
		{WithRefreshLeaseTTL(0)},
		{WithRefreshLeaseRenewal(0)},
		{WithRefreshLeaseTTL(time.Second), WithRefreshLeaseRenewal(time.Second)},
		{nil},
	} {
		if _, err := NewRedisRefreshCoordinator(client, key, options...); !errors.Is(err, ErrInvalidOption) {
			t.Fatalf("invalid options error = %v, want ErrInvalidOption", err)
		}
	}

	coordinator, err := NewRedisRefreshCoordinator(client, func(string) string { return "" })
	if err != nil {
		t.Fatal(err)
	}
	if _, err = coordinator.Acquire(context.Background(), "empty-key"); !errors.Is(err, ErrInvalidOption) {
		t.Fatalf("empty key error = %v, want ErrInvalidOption", err)
	}
	coordinator, err = NewRedisRefreshCoordinator(client, key)
	if err != nil {
		t.Fatal(err)
	}
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err = coordinator.Acquire(canceled, "canceled"); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled Acquire error = %v, want context canceled", err)
	}
	var nilContext context.Context
	if _, err = coordinator.Acquire(nilContext, "nil-context"); !errors.Is(err, ErrInvalidOption) {
		t.Fatalf("nil context error = %v, want ErrInvalidOption", err)
	}
	if _, err = coordinator.Acquire(context.Background(), "redis-error"); !errors.Is(err, ErrRefreshCoordinator) {
		t.Fatalf("Redis failure error = %v, want ErrRefreshCoordinator", err)
	}
}

func redisTestClient(t *testing.T) *redis.Client {
	t.Helper()
	address := os.Getenv("TWITCHY_TEST_REDIS_ADDR")
	if address == "" {
		t.Skip("TWITCHY_TEST_REDIS_ADDR is not set")
	}
	client := redis.NewClient(&redis.Options{Addr: address})
	t.Cleanup(func() { _ = client.Close() })
	if err := client.Ping(context.Background()).Err(); err != nil {
		t.Fatalf("ping Redis: %v", err)
	}
	return client
}
