package retryclassifier

import (
	"context"
	"errors"
	"testing"
)

type tempErr struct{}

func (tempErr) Error() string   { return "temporary" }
func (tempErr) Temporary() bool { return true }

func TestShouldRetryHTTPStatus(t *testing.T) {
	policy := Policy{MaxAttempts: 3}

	if !policy.ShouldRetry(1, StatusError{Code: 503, Err: errors.New("unavailable")}) {
		t.Fatal("503 should be retryable")
	}
	if !policy.ShouldRetry(1, StatusError{Code: 429, Err: errors.New("rate limited")}) {
		t.Fatal("429 should be retryable")
	}
	if policy.ShouldRetry(1, StatusError{Code: 400, Err: errors.New("bad request")}) {
		t.Fatal("400 should not be retryable")
	}
}

func TestShouldRetryPermanentAndCanceled(t *testing.T) {
	policy := Policy{MaxAttempts: 3}

	if policy.ShouldRetry(1, ErrPermanent) {
		t.Fatal("permanent error should not retry")
	}
	if policy.ShouldRetry(1, context.Canceled) {
		t.Fatal("canceled context should not retry")
	}
}

func TestShouldRetryTemporaryAndAttemptCap(t *testing.T) {
	policy := Policy{MaxAttempts: 2}

	if !policy.ShouldRetry(1, tempErr{}) {
		t.Fatal("temporary error should retry before cap")
	}
	if policy.ShouldRetry(2, tempErr{}) {
		t.Fatal("temporary error should stop at attempt cap")
	}
}

func TestShouldRetryDeadlineExceeded(t *testing.T) {
	policy := Policy{MaxAttempts: 4}

	if !policy.ShouldRetry(1, context.DeadlineExceeded) {
		t.Fatal("deadline exceeded should retry")
	}
}
