package refreshbackoff

import (
	"testing"
	"time"
)

func TestPolicyNextCapsExponentialDelay(t *testing.T) {
	now := time.Unix(0, 0)
	policy := Policy{Base: time.Second, Max: 5 * time.Second}
	tests := []struct {
		attempt int
		want    time.Duration
	}{{1, time.Second}, {2, 2 * time.Second}, {3, 4 * time.Second}, {4, 5 * time.Second}, {100, 5 * time.Second}}
	for _, tt := range tests {
		got, err := policy.Next(now, tt.attempt)
		if err != nil {
			t.Fatal(err)
		}
		if delay := got.Sub(now); delay != tt.want {
			t.Fatalf("attempt %d delay = %v, want %v", tt.attempt, delay, tt.want)
		}
	}
}

func TestPolicyValidationAndReady(t *testing.T) {
	now := time.Now()
	if _, err := (Policy{}).Next(now, 1); err == nil {
		t.Fatal("zero policy should fail")
	}
	if _, err := (Policy{Base: 2, Max: 1}).Next(now, 1); err == nil {
		t.Fatal("max below base should fail")
	}
	if _, err := (Policy{Base: 1, Max: 2}).Next(now, 0); err == nil {
		t.Fatal("zero attempt should fail")
	}
	if Ready(now, now.Add(time.Second)) {
		t.Fatal("retry should not be ready early")
	}
	if !Ready(now, now) {
		t.Fatal("retry should be ready at deadline")
	}
}
