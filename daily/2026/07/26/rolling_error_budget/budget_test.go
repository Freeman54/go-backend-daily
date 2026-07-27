package rolling_error_budget

import (
	"testing"
	"time"
)

func TestBudget_AllowsWhenRollingErrorRateIsWithinLimit(t *testing.T) {
	b, err := New(4, 0.5)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Unix(100, 0)
	b.Record(now, true)
	b.Record(now, false)
	if !b.Allow(now) {
		t.Fatal("50% errors should be allowed at a 50% limit")
	}
	b.Record(now, false)
	if b.Allow(now) {
		t.Fatal("error rate above limit should be rejected")
	}
}

func TestBudget_ExpiresOldBuckets(t *testing.T) {
	b, _ := New(2, 0)
	now := time.Unix(200, 0)
	b.Record(now, false)
	if b.Allow(now) {
		t.Fatal("current error should exhaust zero error budget")
	}
	if !b.Allow(now.Add(2 * time.Second)) {
		t.Fatal("expired error should not affect current window")
	}
}

func TestNew_ValidatesConfiguration(t *testing.T) {
	if _, err := New(0, 0.1); err == nil {
		t.Fatal("zero window should fail")
	}
	if _, err := New(1, 1.1); err == nil {
		t.Fatal("limit above one should fail")
	}
}
