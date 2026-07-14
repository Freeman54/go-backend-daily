package slowstartlimiter

import "testing"

func TestLimiterRampsUpAfterStableSuccesses(t *testing.T) {
	limiter := New(4, 2)

	limiter.OnSuccess()
	if limiter.Limit() != 1 {
		t.Fatalf("limit = %d want 1", limiter.Limit())
	}

	limiter.OnSuccess()
	if limiter.Limit() != 2 {
		t.Fatalf("limit = %d want 2", limiter.Limit())
	}

	limiter.OnSuccess()
	limiter.OnSuccess()
	if limiter.Limit() != 3 {
		t.Fatalf("limit = %d want 3", limiter.Limit())
	}
}

func TestLimiterResetsAfterFailure(t *testing.T) {
	limiter := New(5, 1)
	limiter.OnSuccess()
	limiter.OnSuccess()
	if limiter.Limit() != 3 {
		t.Fatalf("limit = %d want 3", limiter.Limit())
	}

	limiter.OnFailure()
	if limiter.Limit() != 1 {
		t.Fatalf("limit = %d want 1", limiter.Limit())
	}
}
