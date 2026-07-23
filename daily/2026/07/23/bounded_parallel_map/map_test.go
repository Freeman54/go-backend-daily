package boundedparallelmap

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestMapPreservesOrderAndBoundsConcurrency(t *testing.T) {
	var running atomic.Int32
	var peak atomic.Int32
	got, err := Map(context.Background(), []int{1, 2, 3, 4}, 2, func(_ context.Context, value int) (int, error) {
		current := running.Add(1)
		defer running.Add(-1)
		for current > peak.Load() && !peak.CompareAndSwap(peak.Load(), current) {
		}
		time.Sleep(time.Millisecond)
		return value * value, nil
	})
	if err != nil {
		t.Fatalf("Map() error = %v", err)
	}
	want := []int{1, 4, 9, 16}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Map()[%d] = %d, want %d", i, got[i], want[i])
		}
	}
	if peak.Load() > 2 {
		t.Fatalf("peak concurrency = %d, want <= 2", peak.Load())
	}
}

func TestMapStopsSchedulingAfterError(t *testing.T) {
	sentinel := errors.New("boom")
	var calls atomic.Int32
	_, err := Map(context.Background(), []int{1, 2, 3}, 1, func(_ context.Context, value int) (int, error) {
		calls.Add(1)
		if value == 2 {
			return 0, sentinel
		}
		return value, nil
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("Map() error = %v, want %v", err, sentinel)
	}
	if calls.Load() != 2 {
		t.Fatalf("calls = %d, want 2", calls.Load())
	}
}

func TestMapRejectsInvalidLimit(t *testing.T) {
	if _, err := Map(context.Background(), []int{1}, 0, func(context.Context, int) (int, error) { return 0, nil }); err == nil {
		t.Fatal("Map() expected invalid limit error")
	}
}
