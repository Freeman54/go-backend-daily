package quorum_response_collector

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestCollect_ReturnsAfterQuorumAndCancelsRemainingTasks(t *testing.T) {
	canceled := make(chan struct{})
	tasks := []Task[int]{
		func(context.Context) (int, error) { return 1, nil },
		func(context.Context) (int, error) { return 2, nil },
		func(ctx context.Context) (int, error) {
			<-ctx.Done()
			close(canceled)
			return 0, ctx.Err()
		},
	}

	got, err := Collect(context.Background(), 2, tasks)
	if err != nil || len(got) != 2 {
		t.Fatalf("Collect() = (%v, %v), want two values", got, err)
	}
	select {
	case <-canceled:
	case <-time.After(time.Second):
		t.Fatal("remaining task was not canceled")
	}
}

func TestCollect_ReturnsJoinedErrorsWhenQuorumIsImpossible(t *testing.T) {
	errA := errors.New("a")
	errB := errors.New("b")
	_, err := Collect(context.Background(), 2, []Task[int]{
		func(context.Context) (int, error) { return 1, nil },
		func(context.Context) (int, error) { return 0, errA },
		func(context.Context) (int, error) { return 0, errB },
	})
	if !errors.Is(err, ErrQuorumUnavailable) || !errors.Is(err, errA) || !errors.Is(err, errB) {
		t.Fatalf("Collect() error = %v, want quorum and task errors", err)
	}
}

func TestCollect_RejectsInvalidQuorum(t *testing.T) {
	if _, err := Collect[int](context.Background(), 0, nil); !errors.Is(err, ErrInvalidQuorum) {
		t.Fatalf("Collect() error = %v, want ErrInvalidQuorum", err)
	}
}
