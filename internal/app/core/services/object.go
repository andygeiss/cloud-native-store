package services

import (
	"context"
	"log"

	"github.com/andygeiss/cloud-native-store/internal/app/core/ports"
	"github.com/andygeiss/cloud-native-utils/consistency"
)

// ObjectService is a generic service for managing objects with transactional logging.
// It uses a transactional logger and an ObjectPort to interact with objects.
type ObjectService[K comparable, V any] struct {
	tx   consistency.Logger[K, V] // Transactional logger for recording operations.
	port ports.ObjectPort[K, V]   // Port interface for object interactions (e.g., CRUD operations).
}

// NewObjectService creates a new instance of ObjectService without any dependencies.
func NewObjectService[K comparable, V any]() *ObjectService[K, V] {
	return &ObjectService[K, V]{}
}

// Delete removes an object identified by the key from the port and logs the operation.
func (a *ObjectService[K, V]) Delete(ctx context.Context, key K) (err error) {
	if err = a.port.Delete(ctx, key); err != nil {
		return
	}
	a.tx.WriteDelete(key)
	return nil
}

// Get retrieves an object identified by the key from the port.
func (a *ObjectService[K, V]) Get(ctx context.Context, key K) (value V, err error) {
	value, err = a.port.Get(ctx, key)
	if err != nil {
		return
	}
	return value, nil
}

// Put adds or updates an object identified by the key and logs the operation.
func (a *ObjectService[K, V]) Put(ctx context.Context, key K, value V) (err error) {
	if err = a.port.Put(ctx, key, value); err != nil {
		return
	}
	a.tx.WritePut(key, value)
	return nil
}

// Setup initializes the service by processing pending events from the transactional logger.
func (a *ObjectService[K, V]) Setup() (err error) {
	// If no transactional logger is present, return early.
	if a.tx == nil {
		return nil
	}

	// Read pending events from the logger.
	eventCh, _ := a.tx.ReadEvents()
	ctx := context.Background()

	// Process each event in the channel.
	for event := range eventCh {
		switch event.EventType {
		case consistency.EventTypeDelete:
			log.Printf("event delete: %v", event)
			_ = a.Delete(ctx, event.Key)
		case consistency.EventTypePut:
			_ = a.Put(ctx, event.Key, event.Value)
			log.Printf("event put: %v", event)
		}
	}

	return nil
}

// Teardown cleans up any resources used by the ObjectService.
func (a *ObjectService[K, V]) Teardown() {
	if err := a.tx.Close(); err != nil {
		log.Fatalf("error during close: %v", err)
	}
}

// WithTransactionalLogger sets the transactional logger for the service and returns the updated service.
func (a *ObjectService[K, V]) WithTransactionalLogger(logger consistency.Logger[K, V]) *ObjectService[K, V] {
	a.tx = logger
	return a
}

// WithPort sets the ObjectPort for the service and returns the updated service.
func (a *ObjectService[K, V]) WithPort(port ports.ObjectPort[K, V]) *ObjectService[K, V] {
	a.port = port
	return a
}
