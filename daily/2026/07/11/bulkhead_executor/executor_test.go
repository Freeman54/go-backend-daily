package bulkheadexecutor

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestExecutorIsolatesLanes(t *testing.T) {
	exec, err := New(map[string]int{
		"db":    1,
		"email": 1,
	})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	release := make(chan struct{})
	errCh := make(chan error, 2)

	go func() {
		errCh <- exec.Do(context.Background(), "db", func(context.Context) error {
			<-release
			return nil
		})
	}()

	time.Sleep(20 * time.Millisecond)

	if _, limit, ok := exec.Occupancy("db"); !ok || limit != 1 {
		t.Fatalf("unexpected db occupancy metadata: ok=%v limit=%d", ok, limit)
	}

	if err := exec.Do(context.Background(), "email", func(context.Context) error { return nil }); err != nil {
		t.Fatalf("email lane should not be blocked by db lane: %v", err)
	}

	close(release)
	if err := <-errCh; err != nil {
		t.Fatalf("db lane should finish cleanly: %v", err)
	}
}

func TestExecutorReturnsBusyWhenContextExpires(t *testing.T) {
	exec, err := New(map[string]int{"db": 1})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	block := make(chan struct{})
	go func() {
		_ = exec.Do(context.Background(), "db", func(context.Context) error {
			<-block
			return nil
		})
	}()

	time.Sleep(20 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	err = exec.Do(ctx, "db", func(context.Context) error { return nil })
	if !errors.Is(err, ErrLaneBusy) {
		t.Fatalf("expected ErrLaneBusy, got %v", err)
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline exceeded, got %v", err)
	}

	close(block)
}
