package boundedfanout

import (
	"context"
	"errors"
	"slices"
	"sync/atomic"
	"testing"
	"time"
)

func TestRunPreservesOrderAndConcurrencyBound(t *testing.T) {
	var running atomic.Int64
	var peak atomic.Int64

	tasks := []Task{
		func(context.Context) (string, error) {
			updatePeak(&running, &peak)
			time.Sleep(20 * time.Millisecond)
			running.Add(-1)
			return "a", nil
		},
		func(context.Context) (string, error) {
			updatePeak(&running, &peak)
			time.Sleep(20 * time.Millisecond)
			running.Add(-1)
			return "b", nil
		},
		func(context.Context) (string, error) {
			updatePeak(&running, &peak)
			time.Sleep(20 * time.Millisecond)
			running.Add(-1)
			return "c", nil
		},
	}

	got, err := Run(context.Background(), 2, tasks)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !slices.Equal(got, []string{"a", "b", "c"}) {
		t.Fatalf("Run() results = %v", got)
	}
	if peak.Load() > 2 {
		t.Fatalf("peak concurrency = %d, want <= 2", peak.Load())
	}
}

func TestRunCancelsSiblingsAfterFirstError(t *testing.T) {
	started := make(chan struct{}, 1)
	cancelObserved := make(chan struct{}, 1)
	wantErr := errors.New("downstream failed")

	tasks := []Task{
		func(context.Context) (string, error) {
			return "", wantErr
		},
		func(ctx context.Context) (string, error) {
			started <- struct{}{}
			<-ctx.Done()
			cancelObserved <- struct{}{}
			return "", ctx.Err()
		},
	}

	_, err := Run(context.Background(), 2, tasks)
	if !errors.Is(err, wantErr) {
		t.Fatalf("Run() error = %v, want %v", err, wantErr)
	}

	select {
	case <-started:
	default:
		t.Fatal("sibling task did not start")
	}
	select {
	case <-cancelObserved:
	default:
		t.Fatal("sibling task did not observe cancellation")
	}
}

func TestRunRejectsInvalidParallelism(t *testing.T) {
	_, err := Run(context.Background(), 0, nil)
	if err == nil {
		t.Fatal("Run() error = nil, want validation error")
	}
}

func updatePeak(running, peak *atomic.Int64) {
	current := running.Add(1)
	for {
		existing := peak.Load()
		if current <= existing {
			return
		}
		if peak.CompareAndSwap(existing, current) {
			return
		}
	}
}
