package contexttimeout

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRunWithWorkersCompletesAllJobs(t *testing.T) {
	jobs := []Job{
		func(context.Context) error { return nil },
		func(context.Context) error { return nil },
		func(context.Context) error { return nil },
	}

	result := RunWithWorkers(context.Background(), 2, jobs)

	if result.Err != nil {
		t.Fatalf("expected nil error, got %v", result.Err)
	}
	if result.Completed != 3 {
		t.Fatalf("expected 3 completed jobs, got %d", result.Completed)
	}
	if result.Canceled {
		t.Fatal("expected run to finish without cancellation")
	}
}

func TestRunWithWorkersCancelsOnFirstError(t *testing.T) {
	wantErr := errors.New("database timeout")
	jobs := []Job{
		func(context.Context) error { return wantErr },
		func(ctx context.Context) error {
			<-ctx.Done()
			return ctx.Err()
		},
	}

	result := RunWithWorkers(context.Background(), 2, jobs)

	if !errors.Is(result.Err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, result.Err)
	}
	if !result.Canceled {
		t.Fatal("expected cancellation after first job error")
	}
}

func TestRunWithWorkersRespectsParentTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	jobs := []Job{
		func(ctx context.Context) error {
			<-ctx.Done()
			return ctx.Err()
		},
	}

	result := RunWithWorkers(ctx, 1, jobs)

	if !errors.Is(result.Err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline exceeded, got %v", result.Err)
	}
	if !result.Canceled {
		t.Fatal("expected timeout to mark result as canceled")
	}
}

func TestRunWithWorkersRejectsInvalidWorkerCount(t *testing.T) {
	result := RunWithWorkers(context.Background(), 0, nil)

	if !errors.Is(result.Err, ErrInvalidWorkerCount) {
		t.Fatalf("expected invalid worker count error, got %v", result.Err)
	}
}
