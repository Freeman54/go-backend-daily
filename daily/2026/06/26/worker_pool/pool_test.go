package workerpool

import (
	"context"
	"errors"
	"slices"
	"testing"
	"time"
)

func TestRunProcessesAllJobsAndKeepsInputOrder(t *testing.T) {
	t.Parallel()

	pool := New(3)
	jobs := []int{1, 2, 3, 4}

	got, err := pool.Run(context.Background(), jobs, func(ctx context.Context, job int) (int, error) {
		time.Sleep(time.Duration(5-job) * time.Millisecond)
		return job * job, nil
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	want := []int{1, 4, 9, 16}
	if !slices.Equal(got, want) {
		t.Fatalf("Run() = %v, want %v", got, want)
	}
}

func TestRunStopsSchedulingWhenContextCancelled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool := New(2)
	_, err := pool.Run(ctx, []int{1, 2, 3}, func(ctx context.Context, job int) (int, error) {
		if job == 1 {
			cancel()
		}
		<-ctx.Done()
		return 0, ctx.Err()
	})

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Run() error = %v, want context.Canceled", err)
	}
}
