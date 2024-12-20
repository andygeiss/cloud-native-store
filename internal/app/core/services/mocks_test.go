package services_test

import (
	"context"
	"errors"
	"sync"

	"github.com/andygeiss/cloud-native-utils/consistency"
)

type MockLogger[K comparable, V any] struct {
	events []consistency.Event[K, V]
	mutex  sync.Mutex
	err    error
}

func NewMockLogger[K comparable, V any]() *MockLogger[K, V] {
	return &MockLogger[K, V]{
		events: make([]consistency.Event[K, V], 0),
	}
}

func (m *MockLogger[K, V]) WriteDelete(key K) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.events = append(m.events, consistency.Event[K, V]{EventType: consistency.EventTypeDelete, Key: key})
}

func (m *MockLogger[K, V]) WritePut(key K, value V) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.events = append(m.events, consistency.Event[K, V]{EventType: consistency.EventTypePut, Key: key, Value: value})
}

func (m *MockLogger[K, V]) ReadEvents() (<-chan consistency.Event[K, V], <-chan error) {
	eventCh := make(chan consistency.Event[K, V], len(m.events))
	errCh := make(chan error, 1)

	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, e := range m.events {
		eventCh <- e
	}
	close(eventCh)

	if m.err != nil {
		errCh <- m.err
	}
	close(errCh)

	return eventCh, errCh
}

func (m *MockLogger[K, V]) Close() error {
	return m.err
}

// MockPort is a mock implementation of the ObjectPort interface.
type MockPort[K comparable, V any] struct {
	data map[K]V
	err  error
}

func NewMockPort[K comparable, V any]() *MockPort[K, V] {
	return &MockPort[K, V]{
		data: make(map[K]V),
	}
}

func (m *MockPort[K, V]) Delete(ctx context.Context, key K) error {
	if m.err != nil {
		return m.err
	}
	delete(m.data, key)
	return nil
}

func (m *MockPort[K, V]) Get(ctx context.Context, key K) (V, error) {
	value, ok := m.data[key]
	if !ok {
		return value, errors.New("key not found")
	}
	return value, nil
}

func (m *MockPort[K, V]) Put(ctx context.Context, key K, value V) error {
	if m.err != nil {
		return m.err
	}
	m.data[key] = value
	return nil
}
