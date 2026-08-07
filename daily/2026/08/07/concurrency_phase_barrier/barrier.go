// Package concurrencyphasebarrier provides a one-shot synchronization barrier.
package concurrencyphasebarrier

import (
	"context"
	"sync"
)

// Barrier releases every arrived caller once target callers have arrived.
type Barrier struct {
	mu      sync.Mutex
	target  int
	arrived int
	done    chan struct{}
}

// New constructs a barrier. A non-positive target is a programming error.
func New(target int) *Barrier {
	if target < 1 {
		panic("barrier target must be positive")
	}
	return &Barrier{target: target, done: make(chan struct{})}
}

// Arrive records this caller and waits for the group or context cancellation.
// An arrival remains counted if its caller is later cancelled.
func (b *Barrier) Arrive(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	b.mu.Lock()
	b.arrived++
	if b.arrived == b.target {
		close(b.done)
	}
	done := b.done
	b.mu.Unlock()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
