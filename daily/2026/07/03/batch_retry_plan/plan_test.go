package batchretryplan

import (
	"reflect"
	"testing"
)

func TestBatchRetryPlanRoutesItems(t *testing.T) {
	results := []ItemResult{
		{ID: "a", Outcome: OutcomeSuccess, Attempt: 0},
		{ID: "b", Outcome: OutcomeRetryable, Attempt: 0},
		{ID: "c", Outcome: OutcomePermanent, Attempt: 2},
		{ID: "d", Outcome: OutcomeRetryable, Attempt: 2},
	}

	plan := Build(results, 3)

	if !reflect.DeepEqual(plan.AckIDs, []string{"a"}) {
		t.Fatalf("unexpected ack ids: %#v", plan.AckIDs)
	}

	if !reflect.DeepEqual(plan.RetryIDs, []string{"b"}) {
		t.Fatalf("unexpected retry ids: %#v", plan.RetryIDs)
	}

	if !reflect.DeepEqual(plan.DLQIDs, []string{"c", "d"}) {
		t.Fatalf("unexpected dlq ids: %#v", plan.DLQIDs)
	}
}

func TestBatchRetryPlanDropsRetryWhenBudgetExhausted(t *testing.T) {
	plan := Build([]ItemResult{
		{ID: "x", Outcome: OutcomeRetryable, Attempt: 0},
	}, 1)

	if len(plan.RetryIDs) != 0 {
		t.Fatalf("expected no retry ids, got %#v", plan.RetryIDs)
	}

	if !reflect.DeepEqual(plan.DLQIDs, []string{"x"}) {
		t.Fatalf("unexpected dlq ids: %#v", plan.DLQIDs)
	}
}
