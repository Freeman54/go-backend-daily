package keyedsemaphore

import (
	"context"
	"testing"
	"time"
)

func TestLimiterLimitsEachKeyIndependently(t *testing.T) {
	limiter := New(1)

	if err := limiter.Acquire(context.Background(), "tenant-a"); err != nil {
		t.Fatalf("acquire tenant-a: %v", err)
	}
	if err := limiter.Acquire(context.Background(), "tenant-b"); err != nil {
		t.Fatalf("acquire tenant-b: %v", err)
	}
}

func TestLimiterBlocksSameKeyUntilRelease(t *testing.T) {
	limiter := New(1)
	if err := limiter.Acquire(context.Background(), "tenant-a"); err != nil {
		t.Fatalf("first acquire: %v", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- limiter.Acquire(context.Background(), "tenant-a")
	}()

	select {
	case err := <-done:
		t.Fatalf("second acquire returned too early: %v", err)
	case <-time.After(40 * time.Millisecond):
	}

	if err := limiter.Release("tenant-a"); err != nil {
		t.Fatalf("release: %v", err)
	}

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("second acquire failed: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("second acquire did not unblock")
	}
}

func TestLimiterAcquireRespectsContextCancellation(t *testing.T) {
	limiter := New(1)
	if err := limiter.Acquire(context.Background(), "tenant-a"); err != nil {
		t.Fatalf("first acquire: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	if err := limiter.Acquire(ctx, "tenant-a"); err == nil {
		t.Fatal("expected acquire to fail after timeout")
	}
}

func TestLimiterRejectsInvalidRelease(t *testing.T) {
	limiter := New(1)
	if err := limiter.Release("missing"); err == nil {
		t.Fatal("expected missing key release to fail")
	}
}
