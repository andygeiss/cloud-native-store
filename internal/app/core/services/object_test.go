package services_test

import (
	"context"
	"testing"

	"github.com/andygeiss/cloud-native-store/internal/app/config"
	"github.com/andygeiss/cloud-native-store/internal/app/core/services"
	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/consistency"
	"github.com/andygeiss/cloud-native-utils/security"
)

func TestObjectService_Delete(t *testing.T) {
	cfg := &config.Config{Key: security.GenerateKey()}
	ctx := context.Background()
	key := "test-key"
	logger := NewMockLogger[string, string]()
	port := NewMockPort[string, string]()
	service := services.NewObjectService(cfg).WithTransactionalLogger(logger).WithPort(port)

	port.data[key] = "value"
	err := service.Delete(ctx, key)

	assert.That(t, "err must be nil", err == nil, true)
	_, exists := port.data[key]
	assert.That(t, "key must be deleted", exists, false)
	assert.That(t, "logger must record delete event", len(logger.events), 1)
	assert.That(t, "event type must be delete", logger.events[0].EventType, consistency.EventTypeDelete)
}

func TestObjectService_Get(t *testing.T) {
	cfg := &config.Config{Key: security.GenerateKey()}
	ctx := context.Background()
	key := "test-key"
	value := "test-value"
	port := NewMockPort[string, string]()
	service := services.NewObjectService(cfg).WithPort(port)

	port.data[key] = value

	result, err := service.Get(ctx, key)

	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "value must be correct", result, value)
}

func TestObjectService_Put(t *testing.T) {
	cfg := &config.Config{Key: security.GenerateKey()}
	ctx := context.Background()
	key := "test-key"
	value := "test-value"
	logger := NewMockLogger[string, string]()
	port := NewMockPort[string, string]()
	service := services.NewObjectService(cfg).WithTransactionalLogger(logger).WithPort(port)

	err := service.Put(ctx, key, value)

	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "value must be stored in port", port.data[key], value)
	assert.That(t, "logger must record put event", len(logger.events), 1)
	assert.That(t, "event type must be put", logger.events[0].EventType, consistency.EventTypePut)
	assert.That(t, "event value must match", logger.events[0].Value, value)
}

func TestObjectService_FunctionPatterns(t *testing.T) {
	cfg := &config.Config{Key: security.GenerateKey()}
	ctx := context.Background()
	key := "test-key"
	value := "test-value"
	logger := NewMockLogger[string, string]()
	port := NewMockPort[string, string]()
	service := services.NewObjectService(cfg).WithTransactionalLogger(logger).WithPort(port)

	// Test Delete pattern
	port.data[key] = value
	err := service.Delete(ctx, key)
	assert.That(t, "err must be nil for delete", err == nil, true)
	_, exists := port.data[key]
	assert.That(t, "key must be deleted", exists, false)

	// Test Get pattern
	port.data[key] = value
	result, err := service.Get(ctx, key)
	assert.That(t, "err must be nil for get", err == nil, true)
	assert.That(t, "value must be correct for get", result, value)

	// Test Put pattern
	err = service.Put(ctx, key, value)
	assert.That(t, "err must be nil for put", err == nil, true)
	assert.That(t, "value must be stored for put", port.data[key], value)
}

func TestObjectService_Setup(t *testing.T) {
	cfg := &config.Config{Key: security.GenerateKey()}
	eventCh := make(chan consistency.Event[string, string], 2)
	errCh := make(chan error, 1)

	eventCh <- consistency.Event[string, string]{
		EventType: consistency.EventTypePut,
		Key:       "test-key",
		Value:     "test-value",
	}
	eventCh <- consistency.Event[string, string]{
		EventType: consistency.EventTypeDelete,
		Key:       "test-key",
	}
	close(eventCh)
	close(errCh)

	logger := NewMockLogger[string, string]()
	port := NewMockPort[string, string]()
	service := services.NewObjectService(cfg).WithTransactionalLogger(logger).WithPort(port)

	logger.events = append(logger.events,
		consistency.Event[string, string]{EventType: consistency.EventTypePut, Key: "test-key", Value: "test-value"},
		consistency.Event[string, string]{EventType: consistency.EventTypeDelete, Key: "test-key"},
	)

	err := service.Setup()
	assert.That(t, "err must be nil for setup", err == nil, true)
	_, exists := port.data["test-key"]
	assert.That(t, "key must be deleted after setup", exists, false)
}
