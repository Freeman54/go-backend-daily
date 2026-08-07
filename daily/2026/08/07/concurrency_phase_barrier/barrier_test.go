package concurrencyphasebarrier

import (
	"context"
	"testing"
	"time"
)

func TestBarrier_ReleasesAllAtTarget(t *testing.T) {
	b := New(2)
	first := make(chan error, 1)
	go func() { first <- b.Arrive(context.Background()) }()
	select {
	case err := <-first:
		t.Fatalf("first arrival returned early: %v", err)
	case <-time.After(15 * time.Millisecond):
	}
	if err := b.Arrive(context.Background()); err != nil {
		t.Fatalf("second arrival error = %v", err)
	}
	select {
	case err := <-first:
		if err != nil {
			t.Fatalf("first arrival error = %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("first arrival was not released")
	}
}

func TestBarrier_ArriveRespectsCancellation(t *testing.T) {
	b := New(2)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := b.Arrive(ctx); err != context.Canceled {
		t.Fatalf("Arrive() error = %v, want %v", err, context.Canceled)
	}
}

func TestNew_PanicsForNonPositiveTarget(t *testing.T) {
	t.Parallel()
	defer func() {
		if recover() == nil {
			t.Fatal("New() did not panic")
		}
	}()
	New(0)
}
