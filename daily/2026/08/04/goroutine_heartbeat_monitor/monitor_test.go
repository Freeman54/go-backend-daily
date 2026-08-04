package goroutine_heartbeat_monitor

import (
	"testing"
	"time"
)

func TestMonitor_ReportsOnlyStaleWorkers(t *testing.T) {
	now := time.Unix(100, 0)
	m := New(5 * time.Second)
	m.Beat("healthy", now.Add(-4*time.Second))
	m.Beat("stale", now.Add(-6*time.Second))
	got := m.Stale(now)
	if len(got) != 1 || got[0] != "stale" {
		t.Fatalf("Stale() = %v, want [stale]", got)
	}
}

func TestMonitor_BeatRefreshesWorkerAndRemoveForgetsIt(t *testing.T) {
	now := time.Unix(100, 0)
	m := New(time.Second)
	m.Beat("worker", now.Add(-2*time.Second))
	m.Beat("worker", now)
	if got := m.Stale(now); len(got) != 0 {
		t.Fatalf("Stale() = %v, want empty", got)
	}
	m.Remove("worker")
	if got := m.Count(); got != 0 {
		t.Fatalf("Count() = %d, want 0", got)
	}
}

func TestNew_RejectsNonPositiveTimeout(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("New() did not panic")
		}
	}()
	New(0)
}
