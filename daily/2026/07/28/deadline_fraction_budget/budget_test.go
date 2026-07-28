package deadline_fraction_budget

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestWithFraction_UsesFractionOfRemainingBudget(t *testing.T) {
	now := time.Unix(100, 0)
	parent, cancel := context.WithDeadline(context.Background(), now.Add(10*time.Second))
	defer cancel()

	ctx, childCancel, err := WithFraction(parent, now, 0.25, 500*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	defer childCancel()
	got, ok := ctx.Deadline()
	if !ok || !got.Equal(now.Add(2500*time.Millisecond)) {
		t.Fatalf("deadline = %v, want %v", got, now.Add(2500*time.Millisecond))
	}
}

func TestWithFraction_AppliesMinimumBudget(t *testing.T) {
	now := time.Now()
	parent, cancel := context.WithDeadline(context.Background(), now.Add(10*time.Second))
	defer cancel()
	ctx, childCancel, err := WithFraction(parent, now, 0.01, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer childCancel()
	got, _ := ctx.Deadline()
	if !got.Equal(now.Add(time.Second)) {
		t.Fatalf("deadline = %v, want %v", got, now.Add(time.Second))
	}
}

func TestWithFraction_RejectsParentWithoutDeadline(t *testing.T) {
	_, _, err := WithFraction(context.Background(), time.Now(), 0.5, time.Second)
	if !errors.Is(err, ErrNoDeadline) {
		t.Fatalf("error = %v, want ErrNoDeadline", err)
	}
}
