package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/andygeiss/cloud-native-store/internal/app/adapters/outbound/inmemory"
	"github.com/andygeiss/cloud-native-store/internal/app/config"
	"github.com/andygeiss/cloud-native-store/internal/app/core/services"
	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/consistency"
)

// ----------------------------------------------------------------------------
// 1) Test Setup() and Teardown() without a logger
// ----------------------------------------------------------------------------

func TestObjectService_Setup_WithoutLogger(t *testing.T) {
	cfg := &config.Config{}
	svc := services.NewObjectService(cfg) // no .WithTransactionalLogger()

	err := svc.Setup()
	assert.That(t, "err must be nil", err, nil)
}

func TestObjectService_Teardown_WithoutLogger(t *testing.T) {
	cfg := &config.Config{}
	svc := services.NewObjectService(cfg) // no .WithTransactionalLogger()

	// Should not cause any panic or error:
	svc.Teardown()

	// If you want to verify something after teardown, you could do:
	// assert.That(t, "teardown must have completed", someCondition, true)
}

// ----------------------------------------------------------------------------
// 2) Test Setup() and Teardown() with a mock/fake logger
// ----------------------------------------------------------------------------

// fakeLogger simulates a transactional logger for testing.
type fakeLogger struct {
	closed      bool
	events      chan consistency.Event[string, string]
	errors      chan error
	wrotePut    map[string]string
	wroteDelete []string
}

// NewFakeLogger creates a new fake logger with optional initial events/errors.
func NewFakeLogger(initialEvents []consistency.Event[string, string], initialErrors []error) *fakeLogger {
	fl := &fakeLogger{
		events:      make(chan consistency.Event[string, string], 10),
		errors:      make(chan error, 10),
		wrotePut:    make(map[string]string),
		wroteDelete: []string{},
	}
	// Populate channels with initial events/errors
	for _, e := range initialEvents {
		fl.events <- e
	}
	for _, err := range initialErrors {
		fl.errors <- err
	}
	// Close them to simulate “end of stream”
	close(fl.events)
	close(fl.errors)
	return fl
}

func (f *fakeLogger) ReadEvents() (<-chan consistency.Event[string, string], <-chan error) {
	return f.events, f.errors
}
func (f *fakeLogger) WritePut(key, value string) {
	f.wrotePut[key] = value
}
func (f *fakeLogger) WriteDelete(key string) {
	f.wroteDelete = append(f.wroteDelete, key)
}
func (f *fakeLogger) Close() error {
	f.closed = true
	return nil
}

// Setup test: ensures events are correctly applied to the port.
func TestObjectService_Setup_WithLogger_Events(t *testing.T) {
	initialEvents := []consistency.Event[string, string]{
		{EventType: consistency.EventTypePut, Key: "k1", Value: "v1"},
		{EventType: consistency.EventTypeDelete, Key: "k2"},
	}
	logger := NewFakeLogger(initialEvents, nil)

	cfg := &config.Config{}
	ctx := context.Background()
	port := inmemory.NewObjectStore(10)

	svc := services.NewObjectService(cfg).
		WithPort(port).
		WithTransactionalLogger(logger)

	err := svc.Setup()
	assert.That(t, "err must be nil", err, nil)

	// Check that "k1" was put with "v1"
	v, getErr := port.Get(ctx, "k1")
	assert.That(t, "getErr must be nil", getErr, nil)
	assert.That(t, "value of k1 must be 'v1'", v, "v1")

	// Check that "k2" was deleted (does not exist)
	_, getErr2 := port.Get(ctx, "k2")
	// We expect "key does not exist" from the inmemory store
	assert.That(t, "err message must be 'key does not exist'", getErr2.Error(), "key does not exist")
}

// Setup test: ensures errors on the logger channel are returned.
func TestObjectService_Setup_WithLogger_Errors(t *testing.T) {
	initialEvents := []consistency.Event[string, string]{
		{EventType: consistency.EventTypePut, Key: "k1", Value: "v1"},
	}
	initialErrors := []error{
		errors.New("failed to read log file"),
	}
	logger := NewFakeLogger(initialEvents, initialErrors)

	cfg := &config.Config{}
	port := inmemory.NewObjectStore(10)

	svc := services.NewObjectService(cfg).
		WithPort(port).
		WithTransactionalLogger(logger)

	err := svc.Setup()
	assert.That(t, "err message must be 'failed to read log file'", err.Error(), "failed to read log file")
}

func TestObjectService_Teardown_WithLogger(t *testing.T) {
	logger := NewFakeLogger(nil, nil)
	cfg := &config.Config{}
	svc := services.NewObjectService(cfg).WithTransactionalLogger(logger)

	svc.Teardown()
	// Check the logger is closed
	assert.That(t, "logger must be closed", logger.closed, true)
}

// ----------------------------------------------------------------------------
// 3) Test WithTransactionalLogger and WithPort chaining
// ----------------------------------------------------------------------------

func TestObjectService_WithTransactionalLogger(t *testing.T) {
	logger := NewFakeLogger(nil, nil)
	svc := services.NewObjectService(&config.Config{}).WithTransactionalLogger(logger)

	// Verify it’s not nil
	assert.That(t, "service must not be nil", svc == nil, false)
}

func TestObjectService_WithPort(t *testing.T) {
	p := inmemory.NewObjectStore(10)
	svc := services.NewObjectService(&config.Config{}).WithPort(p)

	// Verify it’s not nil
	assert.That(t, "service must not be nil", svc == nil, false)
}

// ----------------------------------------------------------------------------
// 4) Test that the logger actually logs Put / Delete
// ----------------------------------------------------------------------------

func TestObjectService_Logger_WritePut_And_Delete(t *testing.T) {
	logger := NewFakeLogger(nil, nil)
	cfg := &config.Config{}
	ctx := context.Background()
	port := inmemory.NewObjectStore(10)

	svc := services.NewObjectService(cfg).
		WithPort(port).
		WithTransactionalLogger(logger)

	// Put something
	err := svc.Put(ctx, "foo", "bar")
	assert.That(t, "err must be nil", err, nil)

	// Confirm the logger recorded the Put for "foo"
	loggedVal, ok := logger.wrotePut["foo"]
	assert.That(t, "logger must have recorded WritePut for 'foo'", ok, true)
	// Because of encryption/base64, it may not be exactly "bar". At least confirm it’s not empty:
	assert.That(t, "loggedVal must not be empty", loggedVal == "", false)

	// Now Delete it
	err = svc.Delete(ctx, "foo")
	assert.That(t, "err must be nil", err, nil)

	// Confirm the logger recorded the Delete
	assert.That(t, "must have exactly 1 delete entry", len(logger.wroteDelete), 1)
	assert.That(t, "must have deleted 'foo'", logger.wroteDelete[0], "foo")
}

// ----------------------------------------------------------------------------
// 5) Test failures from the port (to see retries, circuit breaker, etc.)
// ----------------------------------------------------------------------------

type failingPort struct {
	failsRemaining int
	store          map[string]string
}

func newFailingPort(fails int) *failingPort {
	return &failingPort{
		failsRemaining: fails,
		store:          map[string]string{},
	}
}

func (f *failingPort) failIfNeeded() error {
	if f.failsRemaining > 0 {
		f.failsRemaining--
		return errors.New("simulated port failure")
	}
	return nil
}

func (f *failingPort) Get(ctx context.Context, key string) (string, error) {
	if err := f.failIfNeeded(); err != nil {
		return "", err
	}
	val, ok := f.store[key]
	if !ok {
		return "", errors.New("key does not exist")
	}
	return val, nil
}

func (f *failingPort) Put(ctx context.Context, key, value string) error {
	if err := f.failIfNeeded(); err != nil {
		return err
	}
	f.store[key] = value
	return nil
}

func (f *failingPort) Delete(ctx context.Context, key string) error {
	if err := f.failIfNeeded(); err != nil {
		return err
	}
	delete(f.store, key)
	return nil
}

// Test that retries eventually succeed if failures < retry limit.
func TestObjectService_Put_WithFailingPort(t *testing.T) {
	// This port fails the first 2 calls, then succeeds on the 3rd.
	p := newFailingPort(2)
	cfg := &config.Config{}
	ctx := context.Background()

	svc := services.NewObjectService(cfg).WithPort(p)

	err := svc.Put(ctx, "foo", "bar")
	// With 3 retries (and only 2 initial failures), it should eventually succeed:
	assert.That(t, "err must be nil after retries", err, nil)

	// Confirm that we can Get the value back
	val, getErr := svc.Get(ctx, "foo")
	assert.That(t, "getErr must be nil", getErr, nil)
	assert.That(t, "value must be 'bar'", val, "bar")
}

// Test that one failure is handled by retry
func TestObjectService_Delete_WithFailingPort(t *testing.T) {
	// This port fails the first time, then succeeds on the second attempt
	p := newFailingPort(1)
	cfg := &config.Config{}
	ctx := context.Background()
	// Put an item in the store directly
	p.store["foo"] = "bar"

	svc := services.NewObjectService(cfg).WithPort(p)

	err := svc.Delete(ctx, "foo")
	assert.That(t, "err must be nil after one retry", err, nil)

	// Confirm "foo" is deleted
	_, getErr := p.Get(ctx, "foo")
	assert.That(t, "err must be 'key does not exist'", getErr.Error(), "key does not exist")
}

// Test that all retries fail and open the circuit breaker.
func TestObjectService_Get_WithFailingPort(t *testing.T) {
	// This port fails all 3 attempts.
	p := newFailingPort(4)
	cfg := &config.Config{}
	ctx := context.Background()
	p.store["foo"] = "bar"

	// This will open the circuit breaker after 3 failures.
	svc := services.NewObjectService(cfg).WithPort(p)

	_, err := svc.Get(ctx, "foo")

	// We expect an error because all 3 attempts fail and open the breaker.
	assert.That(t, "err must not be nil", err == nil, false)

	// Max retries reached, so we expect the error message from the port.
	assert.That(t, "error must be correct", err.Error(), "simulated port failure")
}
