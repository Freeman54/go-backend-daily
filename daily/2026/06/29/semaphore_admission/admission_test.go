package semaphoreadmission

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestGateLimitsConcurrency(t *testing.T) {
	t.Parallel()

	gate := NewGate(2)
	var active atomic.Int32
	var maxActive atomic.Int32

	var wg sync.WaitGroup
	for range 6 {
		wg.Add(1)
		go func() {
			defer wg.Done()

			err := gate.Do(context.Background(), func(context.Context) error {
				current := active.Add(1)
				defer active.Add(-1)

				for {
					seen := maxActive.Load()
					if current <= seen || maxActive.CompareAndSwap(seen, current) {
						break
					}
				}

				time.Sleep(15 * time.Millisecond)
				return nil
			})
			if err != nil {
				t.Errorf("Do() error = %v", err)
			}
		}()
	}

	wg.Wait()

	if maxActive.Load() != 2 {
		t.Fatalf("max concurrency = %d, want 2", maxActive.Load())
	}
	if gate.InFlight() != 0 {
		t.Fatalf("inflight = %d, want 0", gate.InFlight())
	}
}

func TestGateAcquireRespectsContext(t *testing.T) {
	t.Parallel()

	gate := NewGate(1)
	if err := gate.Acquire(context.Background()); err != nil {
		t.Fatalf("Acquire() error = %v", err)
	}
	defer gate.Release()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	start := time.Now()
	err := gate.Acquire(ctx)
	if err == nil {
		t.Fatal("Acquire() error = nil, want deadline exceeded")
	}
	if time.Since(start) < 15*time.Millisecond {
		t.Fatal("Acquire() returned too early without waiting for context")
	}
}
