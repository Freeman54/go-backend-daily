package contexthandoff

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestPlanAllocatesBudgetByWeight(t *testing.T) {
	now := time.Date(2026, 7, 5, 10, 0, 0, 0, time.UTC)
	ctx, cancel := context.WithDeadline(context.Background(), now.Add(100*time.Millisecond))
	defer cancel()

	slices, err := Plan(ctx, now, 10*time.Millisecond, 10*time.Millisecond, []Stage{
		{Name: "cache", Weight: 1},
		{Name: "db", Weight: 3},
		{Name: "cleanup", Weight: 1},
	})
	if err != nil {
		t.Fatalf("Plan returned error: %v", err)
	}

	want := []time.Duration{22 * time.Millisecond, 46 * time.Millisecond, 22 * time.Millisecond}
	for i, slice := range slices {
		if slice.Duration != want[i] {
			t.Fatalf("slice %d = %v, want %v", i, slice.Duration, want[i])
		}
	}
}

func TestPlanRejectsMissingDeadline(t *testing.T) {
	_, err := Plan(context.Background(), time.Now(), 0, 0, []Stage{{Name: "rpc"}})
	if !errors.Is(err, ErrNoDeadline) {
		t.Fatalf("expected ErrNoDeadline, got %v", err)
	}
}

func TestPlanRejectsInsufficientBudget(t *testing.T) {
	now := time.Date(2026, 7, 5, 10, 0, 0, 0, time.UTC)
	ctx, cancel := context.WithDeadline(context.Background(), now.Add(20*time.Millisecond))
	defer cancel()

	_, err := Plan(ctx, now, 5*time.Millisecond, 10*time.Millisecond, []Stage{
		{Name: "a"},
		{Name: "b"},
	})
	if !errors.Is(err, ErrInsufficientBudget) {
		t.Fatalf("expected ErrInsufficientBudget, got %v", err)
	}
}
