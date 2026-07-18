package deliveryretrybudget

import "testing"

func TestTransientErrorCanRetryBeforeBudgetRunsOut(t *testing.T) {
	budget := Budget{MaxAttempts: 3}

	if got := budget.Decide(2, true); got != Retry {
		t.Fatalf("expected retry, got %v", got)
	}
}

func TestBudgetExhaustionMovesMessageToQuarantine(t *testing.T) {
	budget := Budget{MaxAttempts: 3}

	if got := budget.Decide(3, true); got != Quarantine {
		t.Fatalf("expected quarantine, got %v", got)
	}
}

func TestNonTransientErrorSkipsRetry(t *testing.T) {
	budget := Budget{MaxAttempts: 5}

	if got := budget.Decide(1, false); got != Quarantine {
		t.Fatalf("expected quarantine, got %v", got)
	}
}
