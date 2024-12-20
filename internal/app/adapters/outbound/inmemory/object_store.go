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
type ObjectStore struct {
	shards efficiency.Sharding[string, string]
}

// NewObjectStore initializes a new ObjectStore with the given number of shards.
// It returns an implementation of the ports.ObjectPort interface.
func NewObjectStore(numShards int) ports.ObjectPort[string, string] {
	return &ObjectStore{
		shards: efficiency.NewSharding[string, string](numShards),
	}
}

// Delete removes a key and its associated value from the store.
// If the key does not exist, it silently returns without an error.
func (a *ObjectStore) Delete(ctx context.Context, key string) (err error) {
	a.shards.Delete(key)
	return nil
}

// Get retrieves the value associated with the given key.
// If the key does not exist, it returns an error (ErrorKeyDoesNotExist).
func (a *ObjectStore) Get(ctx context.Context, key string) (value string, err error) {
	value, exists := a.shards.Get(key)
	if !exists {
		return value, ErrorKeyDoesNotExist
	}
	return value, nil
}

// Put inserts or updates the value associated with the given key in the store.
func (a *ObjectStore) Put(ctx context.Context, key, value string) (err error) {
	a.shards.Put(key, value)
	return nil
}
