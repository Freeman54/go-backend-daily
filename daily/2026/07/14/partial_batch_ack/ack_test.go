package partialbatchack

import "testing"

func TestBuildPlanSplitsBatchByDecision(t *testing.T) {
	plan := BuildPlan([]ItemResult{
		{ID: "1", Decision: DecisionAck},
		{ID: "2", Decision: DecisionRetry},
		{ID: "3", Decision: DecisionDrop},
		{ID: "4", Decision: DecisionAck},
	})

	if len(plan.AckIDs) != 2 {
		t.Fatalf("AckIDs len = %d want 2", len(plan.AckIDs))
	}
	if len(plan.RetryIDs) != 1 || plan.RetryIDs[0] != "2" {
		t.Fatalf("RetryIDs = %v want [2]", plan.RetryIDs)
	}
	if len(plan.DropIDs) != 1 || plan.DropIDs[0] != "3" {
		t.Fatalf("DropIDs = %v want [3]", plan.DropIDs)
	}
}

func TestRetryRatio(t *testing.T) {
	plan := BuildPlan([]ItemResult{
		{ID: "1", Decision: DecisionAck},
		{ID: "2", Decision: DecisionRetry},
		{ID: "3", Decision: DecisionRetry},
		{ID: "4", Decision: DecisionDrop},
	})

	if got := plan.RetryRatio(); got != 0.5 {
		t.Fatalf("RetryRatio = %v want 0.5", got)
	}
}
