package retrybudget

import (
	"testing"
	"time"
)

func TestBudgetAllowsBaseRetriesWithoutSuccesses(t *testing.T) {
	now := time.Unix(0, 0)
	budget := New(0.2, 2, time.Minute)

	if !budget.AllowRetry("payments", now) {
		t.Fatal("first retry should be allowed")
	}
	if !budget.AllowRetry("payments", now.Add(time.Second)) {
		t.Fatal("second retry should be allowed")
	}
	if budget.AllowRetry("payments", now.Add(2*time.Second)) {
		t.Fatal("third retry should be blocked")
	}
}

func TestBudgetScalesWithSuccessVolume(t *testing.T) {
	now := time.Unix(0, 0)
	budget := New(0.5, 1, time.Minute)

	for range 4 {
		budget.RecordSuccess("search", now)
	}

	for i := 0; i < 3; i++ {
		if !budget.AllowRetry("search", now.Add(time.Duration(i)*time.Second)) {
			t.Fatalf("retry %d should be allowed", i+1)
		}
	}
	if budget.AllowRetry("search", now.Add(4*time.Second)) {
		t.Fatal("retry above budget should be blocked")
	}
}

func TestBudgetResetsAfterWindow(t *testing.T) {
	now := time.Unix(0, 0)
	budget := New(0, 1, 5*time.Second)

	if !budget.AllowRetry("inventory", now) {
		t.Fatal("retry should be allowed in first window")
	}
	if budget.AllowRetry("inventory", now.Add(time.Second)) {
		t.Fatal("second retry should be blocked in first window")
	}
	if !budget.AllowRetry("inventory", now.Add(6*time.Second)) {
		t.Fatal("retry should be allowed after window resets")
	}
}
