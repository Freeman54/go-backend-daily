package spaneventlimiter

import "testing"

func TestAllowCapsDuplicateEvents(t *testing.T) {
	limiter := New(2)

	if !limiter.Allow("db.retry") {
		t.Fatal("first event should pass")
	}
	if !limiter.Allow("db.retry") {
		t.Fatal("second event should pass")
	}
	if limiter.Allow("db.retry") {
		t.Fatal("third event should be dropped")
	}
}

func TestFlushSummariesReturnsDroppedCountsAndResets(t *testing.T) {
	limiter := New(1)
	limiter.Allow("rpc.timeout")
	limiter.Allow("rpc.timeout")

	summaries := limiter.FlushSummaries()
	if len(summaries) != 1 {
		t.Fatalf("summaries len = %d want 1", len(summaries))
	}

	if !limiter.Allow("rpc.timeout") {
		t.Fatal("limiter should reset after flush")
	}
}
