package fallbackbudget

import (
	"testing"
	"time"
)

func TestAllowWithinBudget(t *testing.T) {
	budget := New(2, time.Minute)
	now := time.Unix(1000, 0)

	if !budget.Allow(now) {
		t.Fatal("first fallback should be allowed")
	}
	if !budget.Allow(now.Add(10 * time.Second)) {
		t.Fatal("second fallback should be allowed")
	}
	if budget.Allow(now.Add(20 * time.Second)) {
		t.Fatal("third fallback should be blocked")
	}
}

func TestWindowReset(t *testing.T) {
	budget := New(1, time.Minute)
	now := time.Unix(2000, 0)

	if !budget.Allow(now) {
		t.Fatal("first fallback should be allowed")
	}
	if budget.Allow(now.Add(30 * time.Second)) {
		t.Fatal("budget should still be exhausted in same window")
	}
	if !budget.Allow(now.Add(61 * time.Second)) {
		t.Fatal("budget should reset after window")
	}
}

func TestRemainingTracksUsage(t *testing.T) {
	budget := New(3, time.Minute)
	now := time.Unix(3000, 0)

	if remaining := budget.Remaining(now); remaining != 3 {
		t.Fatalf("remaining = %d want 3", remaining)
	}
	_ = budget.Allow(now)
	if remaining := budget.Remaining(now.Add(time.Second)); remaining != 2 {
		t.Fatalf("remaining = %d want 2", remaining)
	}
}
