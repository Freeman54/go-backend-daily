package deadlineclamp

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestClampUsesMaxBudgetWhenParentHasNoDeadline(t *testing.T) {
	ctx, cancel, budget, err := Clamp(context.Background(), 20*time.Millisecond, 80*time.Millisecond)
	if err != nil {
		t.Fatalf("Clamp() error = %v", err)
	}
	defer cancel()

	if budget != 80*time.Millisecond {
		t.Fatalf("budget = %v, want %v", budget, 80*time.Millisecond)
	}
	deadline, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected derived deadline")
	}
	if remaining := time.Until(deadline); remaining > 80*time.Millisecond || remaining < 40*time.Millisecond {
		t.Fatalf("remaining budget = %v, want around 80ms", remaining)
	}
}

func TestClampUsesParentRemainingWhenSmallerThanMax(t *testing.T) {
	parent, cancelParent := context.WithTimeout(context.Background(), 45*time.Millisecond)
	defer cancelParent()

	ctx, cancel, budget, err := Clamp(parent, 20*time.Millisecond, 200*time.Millisecond)
	if err != nil {
		t.Fatalf("Clamp() error = %v", err)
	}
	defer cancel()

	if budget > 45*time.Millisecond || budget < 10*time.Millisecond {
		t.Fatalf("budget = %v, want clamped to parent remainder", budget)
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected derived deadline")
	}
	if remaining := time.Until(deadline); remaining > 45*time.Millisecond || remaining < 10*time.Millisecond {
		t.Fatalf("remaining budget = %v, want clamped to parent remainder", remaining)
	}
}

func TestClampRejectsTooLittleParentBudget(t *testing.T) {
	parent, cancelParent := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancelParent()
	time.Sleep(5 * time.Millisecond)

	_, _, _, err := Clamp(parent, 20*time.Millisecond, 50*time.Millisecond)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Clamp() error = %v, want %v", err, context.DeadlineExceeded)
	}
}

func TestClampRejectsInvalidInputs(t *testing.T) {
	_, _, _, err := Clamp(context.Background(), 0, time.Second)
	if err == nil {
		t.Fatal("expected validation error for minBudget")
	}

	_, _, _, err = Clamp(context.Background(), time.Second, 500*time.Millisecond)
	if err == nil {
		t.Fatal("expected validation error for maxBudget")
	}
}
