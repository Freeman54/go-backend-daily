package context_cleanup_barrier

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestWait_CompletesAllCleanup(t *testing.T) {
	tasks := []Task{func() error { return nil }, func() error { return errors.New("flush failed") }}
	err := Wait(context.Background(), tasks...)
	if err == nil || err.Error() != "flush failed" {
		t.Fatalf("Wait() error = %v", err)
	}
}

func TestWait_StopsWaitingWhenContextExpires(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	err := Wait(ctx, func() error { time.Sleep(100 * time.Millisecond); return nil })
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Wait() error = %v", err)
	}
}
