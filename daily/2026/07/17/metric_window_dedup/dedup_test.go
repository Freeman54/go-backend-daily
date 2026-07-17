package metricwindowdedup

import (
	"testing"
	"time"
)

func TestShouldEmitSuppressesDuplicatesWithinWindow(t *testing.T) {
	deduper := New(30 * time.Second)
	now := time.Unix(100, 0)

	if !deduper.ShouldEmit("db-primary-down", now) {
		t.Fatal("first emit should pass")
	}
	if deduper.ShouldEmit("db-primary-down", now.Add(10*time.Second)) {
		t.Fatal("duplicate emit should be suppressed")
	}
	if !deduper.ShouldEmit("db-primary-down", now.Add(31*time.Second)) {
		t.Fatal("emit after window should pass")
	}
}

func TestShouldEmitTracksEachKeyIndependently(t *testing.T) {
	deduper := New(time.Minute)
	now := time.Unix(100, 0)

	if !deduper.ShouldEmit("cache-miss-spike", now) {
		t.Fatal("first key should pass")
	}
	if !deduper.ShouldEmit("mq-backlog", now.Add(5*time.Second)) {
		t.Fatal("different key should pass")
	}
}
