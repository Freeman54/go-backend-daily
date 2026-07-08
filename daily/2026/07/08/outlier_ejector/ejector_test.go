package outlierejector

import (
	"testing"
	"time"
)

func TestNodeGetsEjectedAfterConsecutiveFailures(t *testing.T) {
	now := time.Unix(0, 0)
	ejector := New(3, 10*time.Second)

	ejector.Record("node-a", false, now)
	ejector.Record("node-a", false, now.Add(time.Second))
	if !ejector.Allow("node-a", now.Add(2*time.Second)) {
		t.Fatal("node should still be available before threshold")
	}

	ejector.Record("node-a", false, now.Add(2*time.Second))
	if ejector.Allow("node-a", now.Add(3*time.Second)) {
		t.Fatal("node should be ejected after threshold failures")
	}
}

func TestSuccessClearsFailureStreak(t *testing.T) {
	now := time.Unix(0, 0)
	ejector := New(2, 5*time.Second)

	ejector.Record("node-b", false, now)
	ejector.Record("node-b", true, now.Add(time.Second))
	ejector.Record("node-b", false, now.Add(2*time.Second))

	if !ejector.Allow("node-b", now.Add(3*time.Second)) {
		t.Fatal("node should not be ejected because success reset the streak")
	}
}

func TestNodeReturnsAfterCooldown(t *testing.T) {
	now := time.Unix(0, 0)
	ejector := New(1, 3*time.Second)

	ejector.Record("node-c", false, now)
	if ejector.Allow("node-c", now.Add(time.Second)) {
		t.Fatal("node should be in cooldown")
	}
	if !ejector.Allow("node-c", now.Add(4*time.Second)) {
		t.Fatal("node should return after cooldown")
	}
}
