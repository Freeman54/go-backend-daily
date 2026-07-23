package contextbudgetreserve

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestWithReserveKeepsCleanupBudget(t *testing.T) {
	parentDeadline := time.Now().Add(time.Second)
	parent, cancel := context.WithDeadline(context.Background(), parentDeadline)
	defer cancel()

	child, childCancel, err := WithReserve(parent, 200*time.Millisecond)
	if err != nil {
		t.Fatalf("WithReserve() error = %v", err)
	}
	defer childCancel()
	got, ok := child.Deadline()
	if !ok {
		t.Fatal("child has no deadline")
	}
	want := parentDeadline.Add(-200 * time.Millisecond)
	if delta := got.Sub(want); delta < -time.Millisecond || delta > time.Millisecond {
		t.Fatalf("deadline delta = %v, want near zero", delta)
	}
}

func TestWithReserveRejectsMissingOrExhaustedBudget(t *testing.T) {
	if _, _, err := WithReserve(context.Background(), time.Second); !errors.Is(err, ErrNoDeadline) {
		t.Fatalf("missing deadline error = %v, want %v", err, ErrNoDeadline)
	}
	parent, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	if _, _, err := WithReserve(parent, 20*time.Millisecond); !errors.Is(err, ErrInsufficientBudget) {
		t.Fatalf("exhausted budget error = %v, want %v", err, ErrInsufficientBudget)
	}
}

func TestWithReserveRejectsNegativeReserve(t *testing.T) {
	parent, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if _, _, err := WithReserve(parent, -time.Millisecond); err == nil {
		t.Fatal("expected negative reserve error")
	}
}
