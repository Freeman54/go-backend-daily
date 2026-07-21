package stagebudget

import (
	"context"
	"testing"
	"time"
)

func TestAllocateSplitsDeadlineByWeight(t *testing.T) {
	now := time.Date(2026, 7, 21, 10, 0, 0, 0, time.UTC)
	ctx, cancel := context.WithDeadline(context.Background(), now.Add(10*time.Second))
	defer cancel()

	planner := NewPlanner(time.Second)
	mustAddStage(t, planner, "validate", 1)
	mustAddStage(t, planner, "query", 2)
	mustAddStage(t, planner, "write", 1)

	got, err := planner.Allocate(ctx, now)
	if err != nil {
		t.Fatalf("allocate: %v", err)
	}

	if got["validate"] != 2250*time.Millisecond {
		t.Fatalf("validate = %v, want 2.25s", got["validate"])
	}
	if got["query"] != 4500*time.Millisecond {
		t.Fatalf("query = %v, want 4.5s", got["query"])
	}
	if got["write"] != 2250*time.Millisecond {
		t.Fatalf("write = %v, want 2.25s", got["write"])
	}
}

func TestAllocateFailsWithoutUsefulBudget(t *testing.T) {
	now := time.Date(2026, 7, 21, 10, 0, 0, 0, time.UTC)
	ctx, cancel := context.WithDeadline(context.Background(), now.Add(500*time.Millisecond))
	defer cancel()

	planner := NewPlanner(time.Second)
	mustAddStage(t, planner, "query", 1)

	if _, err := planner.Allocate(ctx, now); err == nil {
		t.Fatal("expected exhausted budget error")
	}
}

func TestAddRejectsInvalidStage(t *testing.T) {
	planner := NewPlanner(0)
	if err := planner.Add("", 1); err == nil {
		t.Fatal("expected empty name to fail")
	}
	if err := planner.Add("query", 0); err == nil {
		t.Fatal("expected non-positive weight to fail")
	}
}

func mustAddStage(t *testing.T, planner *Planner, name string, weight int) {
	t.Helper()
	if err := planner.Add(name, weight); err != nil {
		t.Fatalf("add %s: %v", name, err)
	}
}
