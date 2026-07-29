package oauth

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestRedisRefreshCoordinator_releaseUsesCallerDeadlineWhenRedisDead(t *testing.T) {
	client := redisTestClient(t)
	key := fmt.Sprintf("oauth:dead-release:%d", time.Now().UnixNano())
	owner := "owner"
	if err := client.SetNX(context.Background(), key, owner, time.Second).Err(); err != nil {
		t.Fatal(err)
	}
	deadClient := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:1",
		Dialer: func(ctx context.Context, _, _ string) (net.Conn, error) {
			<-ctx.Done()
			return nil, ctx.Err()
		},
	})
	t.Cleanup(func() { _ = deadClient.Close() })
	coordinator, err := NewRedisRefreshCoordinator(
		deadClient,
		func(string) string { return key },
		WithRefreshLeaseTTL(300*time.Millisecond),
		WithRefreshLeaseRenewal(100*time.Millisecond),
	)
	if err != nil {
		t.Fatal(err)
	}
	lease := newRedisRefreshLease(context.Background(), coordinator, key, owner)
	releaseContext, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	result := make(chan error, 1)
	started := time.Now()
	go func() { result <- lease.Release(releaseContext) }()
	select {
	case err = <-result:
		if elapsed := time.Since(started); elapsed > 250*time.Millisecond {
			t.Fatalf("Release latency = %s, want <= 250ms", elapsed)
		}
	case <-time.After(250 * time.Millisecond):
		lease.cancel()
		<-result
		t.Fatal("Release did not cancel renewal before waiting")
	}
	if !errors.Is(err, ErrRefreshCoordinator) || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Release error = %v, want coordinator and deadline errors", err)
	}
	select {
	case <-lease.done:
	case <-time.After(time.Second):
		t.Fatal("renew goroutine did not exit")
	}
}

func TestRedisRefreshCoordinator_concurrentReleaseDuringRenewal(t *testing.T) {
	fixture := newRedisCoordinatorFixture(t)
	lease, err := fixture.coordinator.Acquire(context.Background(), "concurrent-release")
	if err != nil {
		t.Fatal(err)
	}
	timer := time.NewTimer(150 * time.Millisecond)
	defer timer.Stop()
	<-timer.C

	var group sync.WaitGroup
	errorsChannel := make(chan error, 16)
	for range 16 {
		group.Add(1)
		go func() {
			defer group.Done()
			context, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			errorsChannel <- lease.Release(context)
		}()
	}
	group.Wait()
	close(errorsChannel)
	for releaseErr := range errorsChannel {
		if releaseErr != nil {
			t.Fatalf("concurrent Release error = %v", releaseErr)
		}
	}
	exists, err := fixture.client.Exists(context.Background(), fixture.key("concurrent-release")).Result()
	if err != nil || exists != 0 {
		t.Fatalf("released key exists = %d, error = %v", exists, err)
	}
}
