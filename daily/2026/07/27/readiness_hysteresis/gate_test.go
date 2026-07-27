package readiness_hysteresis

import "testing"

func TestGate_RequiresConsecutiveSuccessesToRecover(t *testing.T) {
	gate, err := New(2, 3)
	if err != nil {
		t.Fatal(err)
	}
	if !gate.Observe(false) {
		t.Fatal("one failure should not close gate")
	}
	if gate.Observe(false) {
		t.Fatal("two consecutive failures should close gate")
	}
	if gate.Observe(true) || gate.Observe(true) {
		t.Fatal("gate should remain closed before recovery threshold")
	}
	if !gate.Observe(true) {
		t.Fatal("three consecutive successes should reopen gate")
	}
}

func TestGate_ResetsStreakWhenObservationChanges(t *testing.T) {
	gate, _ := New(2, 2)
	gate.Observe(false)
	if !gate.Observe(true) {
		t.Fatal("success should reset failure streak")
	}
	if !gate.Observe(false) {
		t.Fatal("failure streak should restart at one")
	}
}

func TestNew_RejectsNonPositiveThresholds(t *testing.T) {
	if _, err := New(0, 1); err == nil {
		t.Fatal("zero failure threshold should fail")
	}
	if _, err := New(1, 0); err == nil {
		t.Fatal("zero recovery threshold should fail")
	}
}
