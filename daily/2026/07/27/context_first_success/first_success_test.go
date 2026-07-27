package context_first_success

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRun_ReturnsFirstSuccessfulResultAndCancelsSlowerWork(t *testing.T) {
	slowCanceled := make(chan struct{})
	tasks := []Task[int]{
		func(context.Context) (int, error) { return 0, errors.New("temporary failure") },
		func(context.Context) (int, error) { return 42, nil },
		func(ctx context.Context) (int, error) {
			<-ctx.Done()
			close(slowCanceled)
			return 0, ctx.Err()
		},
	}

	got, err := Run(context.Background(), tasks)
	if err != nil || got != 42 {
		t.Fatalf("Run() = (%d, %v), want (42, nil)", got, err)
	}
	select {
	case <-slowCanceled:
	case <-time.After(time.Second):
		t.Fatal("slower task was not canceled")
	}
}

func TestRun_ReturnsJoinedErrorWhenAllTasksFail(t *testing.T) {
	errA := errors.New("a")
	errB := errors.New("b")
	_, err := Run(context.Background(), []Task[int]{
		func(context.Context) (int, error) { return 0, errA },
		func(context.Context) (int, error) { return 0, errB },
	})
	if !errors.Is(err, errA) || !errors.Is(err, errB) {
		t.Fatalf("Run() error = %v, want both task errors", err)
	}
}

func TestRun_RejectsEmptyTaskList(t *testing.T) {
	if _, err := Run[int](context.Background(), nil); !errors.Is(err, ErrNoTasks) {
		t.Fatalf("Run() error = %v, want ErrNoTasks", err)
	}
}
