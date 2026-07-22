package gracefuldraintracker

import (
	"context"
	"testing"
	"time"
)

func TestTrackerWaitsUntilAllTasksFinish(t *testing.T) {
	tracker := New()
	doneA, err := tracker.Begin()
	if err != nil {
		t.Fatalf("begin: %v", err)
	}
	doneB, _ := tracker.Begin()
	doneA()
	doneB()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := tracker.Drain(ctx); err != nil {
		t.Fatalf("drain: %v", err)
	}
	if _, err := tracker.Begin(); err == nil {
		t.Fatal("expected begin rejection after drain")
	}
}

func TestTrackerDrainHonorsContext(t *testing.T) {
	tracker := New()
	_, _ = tracker.Begin()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := tracker.Drain(ctx); err == nil {
		t.Fatal("expected context error")
	}
}
