package shardedset

import (
	"hash/fnv"
)

// Hasher is a function that returns a consistent hash for the given key.
type Hasher[K comparable] func(K) uint64

// newFNVHasher returns a Hasher for string keys based on FNV-1a.
func newFNVHasher() Hasher[string] {
	return func(key string) uint64 {
		h := fnv.New64a()
		_, _ = h.Write([]byte(key))
		return h.Sum64()
	}
}
