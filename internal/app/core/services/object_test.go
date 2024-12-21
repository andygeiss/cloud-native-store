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
	logger := NewMockLogger[string, string]()
	port := NewMockPort[string, string]()
	service := services.NewObjectService(cfg).WithTransactionalLogger(logger).WithPort(port)

	encryptedValue := string(security.Encrypt([]byte(value), cfg.Key))
	port.data[key] = encryptedValue

	result, err := service.Get(ctx, key)

	assert.That(t, "err must be nil", err == nil, true)
	assert.That(t, "value must match original", result, value)
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
	encryptedValue := port.data[key]
	decryptedValue, _ := security.Decrypt([]byte(encryptedValue), cfg.Key)
	assert.That(t, "value must match original after decryption", string(decryptedValue), value)
	assert.That(t, "logger must record put event", len(logger.events), 1)
	assert.That(t, "event type must be put", logger.events[0].EventType, consistency.EventTypePut)
}

func TestObjectService_Setup(t *testing.T) {
	cfg := &config.Config{Key: security.GenerateKey()}
	logger := NewMockLogger[string, string]()
	port := NewMockPort[string, string]()
	service := services.NewObjectService(cfg).WithTransactionalLogger(logger).WithPort(port)

	logger.events = append(logger.events, consistency.Event[string, string]{
		EventType: consistency.EventTypePut,
		Key:       "test-key",
		Value:     string(security.Encrypt([]byte("test-value"), cfg.Key)),
	})

	err := service.Setup()

	assert.That(t, "err must be nil", err == nil, true)
	value, exists := port.data["test-key"]
	assert.That(t, "key must exist in port", exists, true)
	decryptedValue, _ := security.Decrypt([]byte(value), cfg.Key)
	assert.That(t, "value must match original", string(decryptedValue), "test-value")
}
