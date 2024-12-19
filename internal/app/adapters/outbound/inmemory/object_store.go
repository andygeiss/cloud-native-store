package inmemory

import (
	"context"
	"errors"

	"github.com/andygeiss/cloud-native-store/internal/app/core/ports"
	"github.com/andygeiss/cloud-native-utils/efficiency"
)

var (
	// ErrorKeyDoesNotExist is returned when a key is not found in the object store.
	ErrorKeyDoesNotExist = errors.New("key does not exist")
)

// ObjectStore is an in-memory object storage that uses sharding for efficient access.
// K represents the type of the keys (comparable) and V represents the type of the values (any).
type ObjectStore[K comparable, V any] struct {
	shards efficiency.Sharding[K, V]
}

// NewObjectStore initializes a new ObjectStore with the given number of shards.
// It returns an implementation of the ports.ObjectPort interface.
func NewObjectStore[K comparable, V any](numShards int) ports.ObjectPort[K, V] {
	return &ObjectStore[K, V]{
		shards: efficiency.NewSharding[K, V](numShards),
	}
}

// Delete removes a key and its associated value from the store.
// If the key does not exist, it silently returns without an error.
func (a *ObjectStore[K, V]) Delete(ctx context.Context, key K) (err error) {
	a.shards.Delete(key)
	return nil
}

// Get retrieves the value associated with the given key.
// If the key does not exist, it returns an error (ErrorKeyDoesNotExist).
func (a *ObjectStore[K, V]) Get(ctx context.Context, key K) (value V, err error) {
	value, exists := a.shards.Get(key)
	if !exists {
		return value, ErrorKeyDoesNotExist
	}
	return value, nil
}

// Put inserts or updates the value associated with the given key in the store.
func (a *ObjectStore[K, V]) Put(ctx context.Context, key K, value V) (err error) {
	a.shards.Put(key, value)
	return nil
}
