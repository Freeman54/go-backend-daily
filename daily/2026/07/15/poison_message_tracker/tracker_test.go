package poisonmessagetracker

import "testing"

func TestFailRetriesUntilThreshold(t *testing.T) {
	tracker := New(3, nil)

	if got := tracker.Fail("msg-1", "timeout"); got != DecisionRetry {
		t.Fatalf("decision = %s want retry", got)
	}
	if got := tracker.Fail("msg-1", "timeout"); got != DecisionRetry {
		t.Fatalf("decision = %s want retry", got)
	}
	if got := tracker.Fail("msg-1", "timeout"); got != DecisionQuarantine {
		t.Fatalf("decision = %s want quarantine", got)
	}
}

func TestFailQuarantinesKnownPoisonCodeImmediately(t *testing.T) {
	tracker := New(5, []string{"schema_invalid"})

	if got := tracker.Fail("msg-2", "schema_invalid"); got != DecisionQuarantine {
		t.Fatalf("decision = %s want quarantine", got)
	}
	if got := tracker.Attempts("msg-2"); got != 0 {
		t.Fatalf("attempts = %d want 0 after quarantine", got)
	}
}

func TestSucceedClearsAttempts(t *testing.T) {
	tracker := New(3, nil)
	tracker.Fail("msg-3", "timeout")
	tracker.Fail("msg-3", "timeout")

	if got := tracker.Succeed("msg-3"); got != DecisionAck {
		t.Fatalf("decision = %s want ack", got)
	}
	if got := tracker.Attempts("msg-3"); got != 0 {
		t.Fatalf("attempts = %d want 0 after success", got)
	}
}
