package leaktracker

import (
	"sync"
	"testing"
	"time"
)

func TestTrackerLifecycleAndOrdering(t *testing.T) {
	base := time.Unix(1_700_000_000, 0)
	current := base
	tracker := New()
	tracker.now = func() time.Time { return current }

	finishOld := tracker.Start("old")
	current = current.Add(5 * time.Second)
	finishNew := tracker.Start("new")

	got := tracker.OlderThan(current.Add(5*time.Second), 6*time.Second)
	if len(got) != 1 || got[0].Name != "old" {
		t.Fatalf("OlderThan() = %+v, want old task", got)
	}
	all := tracker.OlderThan(current, 0)
	if len(all) != 2 || all[0].Name != "old" || all[1].Name != "new" {
		t.Fatalf("OlderThan() ordering = %+v", all)
	}

	finishOld()
	finishOld()
	finishNew()
	if got := tracker.OlderThan(current, 0); len(got) != 0 {
		t.Fatalf("completed tasks remain: %+v", got)
	}
}

func TestTrackerConcurrentCompletion(t *testing.T) {
	tracker := New()
	const count = 50
	finishes := make([]func(), count)
	for i := range finishes {
		finishes[i] = tracker.Start("worker")
	}
	var wg sync.WaitGroup
	for _, finish := range finishes {
		wg.Add(1)
		go func(done func()) {
			defer wg.Done()
			done()
		}(finish)
	}
	wg.Wait()
	if got := tracker.OlderThan(time.Now(), 0); len(got) != 0 {
		t.Fatalf("completed tasks remain: %d", len(got))
	}
}
