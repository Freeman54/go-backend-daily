package aftercommithooks

import (
	"context"
	"errors"
	"slices"
	"testing"
)

func TestCommitRunsHooksInOrder(t *testing.T) {
	var unit Unit
	var order []string

	for _, step := range []string{"outbox", "cache", "audit"} {
		step := step
		if err := unit.OnCommit(func(context.Context) error {
			order = append(order, step)
			return nil
		}); err != nil {
			t.Fatalf("OnCommit returned error: %v", err)
		}
	}

	if err := unit.Commit(context.Background()); err != nil {
		t.Fatalf("Commit returned error: %v", err)
	}

	if !slices.Equal(order, []string{"outbox", "cache", "audit"}) {
		t.Fatalf("hook order = %v", order)
	}
}

func TestCommitJoinsHookErrors(t *testing.T) {
	var unit Unit
	errA := errors.New("publish failed")
	errB := errors.New("metric flush failed")

	_ = unit.OnCommit(func(context.Context) error { return errA })
	_ = unit.OnCommit(func(context.Context) error { return errB })

	err := unit.Commit(context.Background())
	if !errors.Is(err, errA) || !errors.Is(err, errB) {
		t.Fatalf("joined error = %v", err)
	}
}

func TestOnCommitRejectsAfterCommit(t *testing.T) {
	var unit Unit
	if err := unit.Commit(context.Background()); err != nil {
		t.Fatalf("first Commit returned error: %v", err)
	}
	if err := unit.OnCommit(func(context.Context) error { return nil }); !errors.Is(err, ErrAlreadyCommitted) {
		t.Fatalf("OnCommit error = %v, want ErrAlreadyCommitted", err)
	}
}
