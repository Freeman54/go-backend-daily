package cachesingleflight

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestGetDeduplicatesConcurrentLoads(t *testing.T) {
	t.Parallel()

	cache := NewCache[string]()
	var loads atomic.Int32

	loader := func(context.Context) (string, time.Duration, error) {
		loads.Add(1)
		time.Sleep(20 * time.Millisecond)
		return "value", time.Minute, nil
	}

	results := make(chan string, 4)
	for range 4 {
		go func() {
			value, err := cache.Get(context.Background(), "article:1", loader)
			if err != nil {
				t.Errorf("Get() returned error: %v", err)
				return
			}
			results <- value
		}()
	}

	for range 4 {
		select {
		case got := <-results:
			if got != "value" {
				t.Fatalf("Get() = %q, want %q", got, "value")
			}
		case <-time.After(time.Second):
			t.Fatal("timed out waiting for cache result")
		}
	}

	if loads.Load() != 1 {
		t.Fatalf("loader called %d times, want 1", loads.Load())
	}
}

func TestGetReloadsAfterExpiry(t *testing.T) {
	t.Parallel()

	cache := NewCache[string]()
	now := time.Unix(0, 0)
	cache.now = func() time.Time { return now }

	var version atomic.Int32
	loader := func(context.Context) (string, time.Duration, error) {
		next := version.Add(1)
		return string(rune('0' + next)), time.Second, nil
	}

	first, err := cache.Get(context.Background(), "profile:1", loader)
	if err != nil {
		t.Fatalf("first Get() error: %v", err)
	}

	now = now.Add(2 * time.Second)

	second, err := cache.Get(context.Background(), "profile:1", loader)
	if err != nil {
		t.Fatalf("second Get() error: %v", err)
	}

	if first == second {
		t.Fatalf("expected reload after expiry, got same value %q", second)
	}
}
