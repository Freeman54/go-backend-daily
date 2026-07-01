package cancelcausegroup

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRunPropagatesFirstErrorAsCause(t *testing.T) {
	t.Parallel()

	expected := errors.New("db timeout")
	seenCause := make(chan error, 1)

	err := Run(context.Background(),
		func(context.Context) error {
			time.Sleep(10 * time.Millisecond)
			return expected
		},
		func(ctx context.Context) error {
			<-ctx.Done()
			seenCause <- context.Cause(ctx)
			return nil
		},
	)
	if !errors.Is(err, expected) {
		t.Fatalf("Run() error = %v, want %v", err, expected)
	}

	got := <-seenCause
	if !errors.Is(got, expected) {
		t.Fatalf("context cause = %v, want %v", got, expected)
	}
}

func TestRunReturnsNilWhenAllTasksSucceed(t *testing.T) {
	t.Parallel()

	err := Run(context.Background(),
		func(context.Context) error { return nil },
		func(context.Context) error { return nil },
	)
	if err != nil {
		t.Fatalf("Run() error = %v, want nil", err)
	}
}
