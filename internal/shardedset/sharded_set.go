package shardedset

import (
	"sync"
	"time"
)

const shardCount uint64 = 64

// ShardedSet is a concurrent set split into fixed shards with TTL-based eviction.
type ShardedSet[K comparable] struct {
	shards   [shardCount]*shard[K]
	hasher   Hasher[K]
	entryTTL time.Duration
	stop     chan struct{}
}

type (
	expirationTimestamp = time.Time

	shard[K comparable] struct {
		mu      sync.Mutex
		entries map[K]expirationTimestamp
	}
)

func newShard[K comparable]() *shard[K] {
	return &shard[K]{
		entries: make(map[K]expirationTimestamp),
	}
}

func (s *shard[K]) evict() {
	now := time.Now()

	s.mu.Lock()
	for key, expiresAt := range s.entries {
		if now.After(expiresAt) {
			delete(s.entries, key)
		}
	}
	s.mu.Unlock()
}

func newShardedSet[K comparable](hasher Hasher[K], entryTTL time.Duration) ShardedSet[K] {
	var shards [shardCount]*shard[K]

	for i := range shardCount {
		shards[i] = newShard[K]()
	}

	shardedSet := ShardedSet[K]{
		shards:   shards,
		hasher:   hasher,
		entryTTL: entryTTL,
		stop:     make(chan struct{}),
	}

	if entryTTL > 0 {
		go shardedSet.startEviction()
	}

	return shardedSet
}

// startEviction runs a single goroutine that sweeps all shards on a fixed interval.
func (ss *ShardedSet[K]) startEviction() {
	ticker := time.NewTicker(ss.entryTTL)
	defer ticker.Stop()

	for {
		select {
		case <-ss.stop:
			return
		case <-ticker.C:
			for _, s := range ss.shards {
				s.evict()
			}
		}
	}
}

// Stop stops the background TTL eviction.
//
// The ShardedSet must not be used after Stop is called.
func (ss *ShardedSet[K]) Stop() {
	select {
	case <-ss.stop:
		// Already stopped.
	default:
		close(ss.stop)
	}
}

// NewString creates a ShardedSet for string keys with the given entry TTL.
func NewString(entryTTL time.Duration) ShardedSet[string] {
	return newShardedSet[string](newFNVHasher(), entryTTL)
}

// SetIfAbsent records the key and reports whether it was already present (a duplicate).
// An expired entry is treated as absent.
func (ss *ShardedSet[K]) SetIfAbsent(key K) bool {
	s := ss.getShard(key)

	s.mu.Lock()
	defer s.mu.Unlock()

	if expiresAt, exists := s.entries[key]; exists {
		if ss.entryTTL > 0 && time.Now().After(expiresAt) {
			// Entry with given key is expired so we can delete it.
			delete(s.entries, key)
			return false
		}
		return true
	}

	s.entries[key] = time.Now().Add(ss.entryTTL)
	return false
}

func (ss *ShardedSet[K]) getShard(key K) *shard[K] {
	shardIndex := ss.hasher(key) & (shardCount - 1)
	return ss.shards[shardIndex]
}
