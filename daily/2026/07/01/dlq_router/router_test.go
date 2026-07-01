package dlqrouter

import "testing"

func TestDecideReturnsAckOnSuccess(t *testing.T) {
	t.Parallel()

	decision := (Classifier{}).Decide(Message{}, nil)
	if decision != DecisionAck {
		t.Fatalf("Decide() = %q, want ack", decision)
	}
}

func TestDecideRetriesTransientErrorWithinBudget(t *testing.T) {
	t.Parallel()

	decision := (Classifier{MaxAttempts: 5}).Decide(Message{Attempts: 2}, ErrTransient)
	if decision != DecisionRetry {
		t.Fatalf("Decide() = %q, want retry", decision)
	}
}

func TestDecideSendsPoisonMessageToDLQ(t *testing.T) {
	t.Parallel()

	decision := (Classifier{MaxAttempts: 5}).Decide(Message{Attempts: 1}, ErrPoison)
	if decision != DecisionDLQ {
		t.Fatalf("Decide() = %q, want dlq", decision)
	}
}

func TestDecideSendsExhaustedTransientErrorToDLQ(t *testing.T) {
	t.Parallel()

	decision := (Classifier{MaxAttempts: 3}).Decide(Message{Attempts: 3}, ErrTransient)
	if decision != DecisionDLQ {
		t.Fatalf("Decide() = %q, want dlq", decision)
	}
}
