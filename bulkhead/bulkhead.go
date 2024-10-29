// /observability/bulkhead/bulkhead.go
package bulkhead

import (
	"context"
	"time"

	"github.com/goletan/config"
	"github.com/goletan/observability/types"
)

// Bulkhead controls the number of concurrent operations allowed in a specific section of code to limit resource consumption.
type Bulkhead struct {
	capacity  int
	semaphore chan struct{}
	timeout   time.Duration
}

// NewBulkhead creates a new Bulkhead with a given capacity and timeout duration.
func NewBulkhead(capacity int, timeout time.Duration) *Bulkhead {
	return &Bulkhead{
		capacity:  capacity,
		semaphore: make(chan struct{}, capacity),
		timeout:   timeout,
	}
}

// InitBulkheadConfig initializes a Bulkhead based on configuration loaded from config library.
func InitBulkheadConfig() (*Bulkhead, error) {
	cfg := &types.ObservabilityConfig{}
	err := config.LoadConfig("Observability", cfg, nil)
	if err != nil {
		return nil, err
	}
	return NewBulkhead(cfg.Bulkhead.Capacity, cfg.Bulkhead.Timeout), nil
}

// Execute attempts to acquire a permit and run the given function within the bulkhead's capacity and timeout.
func (b *Bulkhead) Execute(ctx context.Context, fn func() error) error {
	select {
	case b.semaphore <- struct{}{}:
		// Permit acquired, execute the function
		defer func() { <-b.semaphore }()
		return fn()
	case <-time.After(b.timeout):
		// Timeout occurred, failed to acquire a permit
		return context.DeadlineExceeded
	case <-ctx.Done():
		// Context canceled before acquiring a permit
		return ctx.Err()
	}
}

// Capacity returns the current capacity of the bulkhead.
func (b *Bulkhead) Capacity() int {
	return b.capacity
}

// Usage returns the current number of occupied slots.
func (b *Bulkhead) Usage() int {
	return len(b.semaphore)
}
