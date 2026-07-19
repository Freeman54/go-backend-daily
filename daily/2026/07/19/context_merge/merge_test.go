package contextmerge

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestMergeUsesEarliestDeadline(t *testing.T) {
	parentA, cancelA := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancelA()

	parentB, cancelB := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancelB()

	merged, cancel := Merge(parentA, parentB)
	defer cancel()

	deadline, ok := merged.Deadline()
	if !ok {
		t.Fatal("expected merged context to have deadline")
	}
	if remaining := time.Until(deadline); remaining > 120*time.Millisecond {
		t.Fatalf("expected earliest deadline, got remaining=%v", remaining)
	}
}

func TestMergeCancelsWhenAnyParentCancels(t *testing.T) {
	parentA, cancelA := context.WithCancelCause(context.Background())
	parentB, cancelB := context.WithCancel(context.Background())
	defer cancelB()

	merged, cancel := Merge(parentA, parentB)
	defer cancel()

	want := errors.New("scheduler stopped")
	cancelA(want)

	select {
	case <-merged.Done():
	case <-time.After(time.Second):
		t.Fatal("merged context was not canceled")
	}

	if got := context.Cause(merged); !errors.Is(got, want) {
		t.Fatalf("unexpected cause: %v", got)
	}
}

func TestMergeCancelFuncStopsContext(t *testing.T) {
	merged, cancel := Merge(context.Background())
	cancel()

	select {
	case <-merged.Done():
	default:
		t.Fatal("expected merged context to be canceled")
	}
}
