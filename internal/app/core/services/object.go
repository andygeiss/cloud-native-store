package services

import (
	"context"
	"log"
	"time"

	"github.com/andygeiss/cloud-native-store/internal/app/core/ports"
	"github.com/andygeiss/cloud-native-utils/consistency"
	"github.com/andygeiss/cloud-native-utils/service"
	"github.com/andygeiss/cloud-native-utils/stability"
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

	fn := func() service.Function[K, V] {
		return func(context.Context, K) (value V, err error) {
			err = a.port.Delete(ctx, key)
			return value, err
		}
	}()

	fn = stability.Timeout(fn, 5*time.Second)
	fn = stability.Debounce(fn, time.Second/time.Duration(10))
	fn = stability.Retry(fn, 3, 5*time.Second)
	fn = stability.Breaker(fn, 3)

	_, err = fn(ctx, key)
	if err != nil {
		return
	}

	a.tx.WriteDelete(key)
	return nil
}

// Get retrieves an object identified by the key from the port.
func (a *ObjectService[K, V]) Get(ctx context.Context, key K) (value V, err error) {

	fn := func() service.Function[K, V] {
		return func(context.Context, K) (V, error) {
			return a.port.Get(ctx, key)
		}
	}()

	fn = stability.Timeout(fn, 5*time.Second)
	fn = stability.Debounce(fn, time.Second/time.Duration(10))
	fn = stability.Retry(fn, 3, 5*time.Second)
	fn = stability.Breaker(fn, 3)

	value, err = fn(ctx, key)
	if err != nil {
		return
	}

	return value, nil
}

// Put adds or updates an object identified by the key and logs the operation.
func (a *ObjectService[K, V]) Put(ctx context.Context, key K, value V) (err error) {

	fn := func() service.Function[K, V] {
		return func(context.Context, K) (V, error) {
			err = a.port.Put(ctx, key, value)
			return value, err
		}
	}()

	fn = stability.Timeout(fn, 5*time.Second)
	fn = stability.Debounce(fn, time.Second/time.Duration(10))
	fn = stability.Retry(fn, 3, 5*time.Second)
	fn = stability.Breaker(fn, 3)

	value, err = fn(ctx, key)
	if err != nil {
		return
	}

	a.tx.WritePut(key, value)
	return nil
}

// Setup initializes the ObjectService by processing pending events
// from the transactional logger and applying them to the data store.
func (a *ObjectService[K, V]) Setup() (err error) {
	// Start reading events and errors from the transactional logger.
	eventCh, errCh := a.tx.ReadEvents()

	// Create a context with cancellation to manage the lifecycle of operations.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure the context is cancelled when the function exits.

	for {
		select {
		case event, ok := <-eventCh:
			if !ok {
				// The event channel has been closed, signaling no more events.
				return nil
			}
			// Handle the specific type of event received.
			switch event.EventType {
			case consistency.EventTypeDelete:
				// If the event is a delete operation, attempt to delete the key from the data store.
				if err := a.port.Delete(ctx, event.Key); err != nil {
					return err // Return the error if the delete operation fails.
				}
			case consistency.EventTypePut:
				// If the event is a put operation, attempt to update the key-value pair in the data store.
				if err := a.port.Put(ctx, event.Key, event.Value); err != nil {
					return err // Return the error if the put operation fails.
				}
			}
		case err, ok := <-errCh:
			// Handle errors reported by the error channel.
			if ok && err != nil {
				return err // Return the error if any occurred during processing.
			}
		}
	}
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
