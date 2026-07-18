package priorityadmission

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestLowPriorityRespectsReservedSlots(t *testing.T) {
	limiter := New(3, 1, 8)

	releaseA, err := limiter.Acquire(context.Background(), Low)
	if err != nil {
		t.Fatalf("acquire low A: %v", err)
	}
	defer releaseA()
	releaseB, err := limiter.Acquire(context.Background(), Low)
	if err != nil {
		t.Fatalf("acquire low B: %v", err)
	}
	defer releaseB()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	if _, err := limiter.Acquire(ctx, Low); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected low request to wait, got %v", err)
	}

	releaseHigh, err := limiter.Acquire(context.Background(), High)
	if err != nil {
		t.Fatalf("acquire high: %v", err)
	}
	releaseHigh()
}

func TestHighPriorityWaiterBeatsLowPriorityWaiter(t *testing.T) {
	limiter := New(2, 1, 8)

	releaseLowA, _ := limiter.Acquire(context.Background(), Low)
	defer releaseLowA()
	releaseHighA, _ := limiter.Acquire(context.Background(), High)
	defer releaseHighA()

	highReady := make(chan struct{})
	lowReady := make(chan struct{})
	holdHigh := make(chan struct{})

	go func() {
		release, err := limiter.Acquire(context.Background(), High)
		if err != nil {
			t.Errorf("high waiter failed: %v", err)
			return
		}
		close(highReady)
		<-holdHigh
		release()
	}()

	go func() {
		release, err := limiter.Acquire(context.Background(), Low)
		if err != nil {
			t.Errorf("low waiter failed: %v", err)
			return
		}
		close(lowReady)
		release()
	}()

	time.Sleep(20 * time.Millisecond)
	releaseLowA()

	select {
	case <-highReady:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("expected high waiter to be admitted first")
	}

	select {
	case <-lowReady:
		t.Fatal("low waiter should still wait until high queue drains")
	case <-time.After(30 * time.Millisecond):
	}

	close(holdHigh)
}

func TestQueueFull(t *testing.T) {
	limiter := New(1, 0, 1)
	release, _ := limiter.Acquire(context.Background(), Low)
	defer release()

	ctx := context.Background()
	go func() {
		_, _ = limiter.Acquire(ctx, Low)
	}()
	time.Sleep(10 * time.Millisecond)

	if _, err := limiter.Acquire(context.Background(), High); !errors.Is(err, ErrQueueFull) {
		t.Fatalf("expected queue full, got %v", err)
	}
}
