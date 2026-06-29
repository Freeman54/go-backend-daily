package gracefulshutdown

import (
	"context"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestShutdownCancelsWorkers(t *testing.T) {
	t.Parallel()

	group := NewGroup()
	var stopped atomic.Int32

	for range 3 {
		group.Go(func(ctx context.Context) error {
			<-ctx.Done()
			stopped.Add(1)
			return nil
		})
	}

	if err := group.Shutdown(100 * time.Millisecond); err != nil {
		t.Fatalf("Shutdown() error = %v", err)
	}
	if stopped.Load() != 3 {
		t.Fatalf("stopped workers = %d, want 3", stopped.Load())
	}
}

func TestShutdownReturnsTimeout(t *testing.T) {
	t.Parallel()

	group := NewGroup()
	release := make(chan struct{})
	group.Go(func(ctx context.Context) error {
		<-ctx.Done()
		<-release
		return nil
	})

	err := group.Shutdown(20 * time.Millisecond)
	close(release)

	if err == nil {
		t.Fatal("Shutdown() error = nil, want timeout")
	}
	if !strings.Contains(err.Error(), "timed out") {
		t.Fatalf("Shutdown() error = %v, want timeout text", err)
	}
}
