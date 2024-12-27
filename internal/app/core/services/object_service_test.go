package services_test

import (
	"context"
	"testing"

	"github.com/andygeiss/cloud-native-store/internal/app/adapters/outbound/inmemory"
	"github.com/andygeiss/cloud-native-store/internal/app/config"
	"github.com/andygeiss/cloud-native-store/internal/app/core/services"
	"github.com/andygeiss/cloud-native-utils/assert"
)

func TestObjectService_Delete_With_No_Key(t *testing.T) {
	cfg := &config.Config{}
	ctx := context.Background()
	port := inmemory.NewObjectStore(1)
	service := services.NewObjectService(cfg).WithPort(port)
	err := service.Delete(ctx, "foo")
	assert.That(t, "err must be nil", err, nil)
}

func TestObjectService_Delete_With_Key(t *testing.T) {
	cfg := &config.Config{}
	ctx := context.Background()
	port := inmemory.NewObjectStore(1)
	service := services.NewObjectService(cfg).WithPort(port)
	service.Put(ctx, "foo", "bar")
	err := service.Delete(ctx, "foo")
	assert.That(t, "err must be nil", err, nil)
}

func TestObjectService_Get_Key_Does_Not_Exist(t *testing.T) {
	cfg := &config.Config{}
	ctx := context.Background()
	port := inmemory.NewObjectStore(1)
	service := services.NewObjectService(cfg).WithPort(port)
	_, err := service.Get(ctx, "foo")
	assert.That(t, "err must be correct", err.Error(), "key does not exist")
}

func TestObjectService_Get_Key_Exist(t *testing.T) {
	cfg := &config.Config{}
	ctx := context.Background()
	port := inmemory.NewObjectStore(1)
	service := services.NewObjectService(cfg).WithPort(port)
	_ = service.Put(ctx, "foo", "bar")
	value, err := service.Get(ctx, "foo")
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "value must be correct", value, "bar")
}

func TestObjectService_Put(t *testing.T) {
	cfg := &config.Config{}
	ctx := context.Background()
	port := inmemory.NewObjectStore(1)
	service := services.NewObjectService(cfg).WithPort(port)
	err := service.Put(ctx, "foo", "bar")
	assert.That(t, "err must be nil", err, nil)
}

func TestObjectService_Put_With_Overwrite(t *testing.T) {
	cfg := &config.Config{}
	ctx := context.Background()
	port := inmemory.NewObjectStore(1)
	service := services.NewObjectService(cfg).WithPort(port)
	service.Put(ctx, "foo", "bar")
	err := service.Put(ctx, "foo", "baz")
	value, err := service.Get(ctx, "foo")
	assert.That(t, "err must be nil", err, nil)
	assert.That(t, "value must be correct", value, "baz")
}
