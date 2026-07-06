package errorjoinpolicy

import (
	"errors"
	"testing"
)

func TestDecidePrefersFatalOverRetryable(t *testing.T) {
	err := errors.Join(
		Retryable(errors.New("db timeout")),
		Fatal(errors.New("bad payload")),
	)

	if got := Decide(err); got != DecisionFail {
		t.Fatalf("Decide() = %q, want %q", got, DecisionFail)
	}
}

func TestDecideReturnsRetryForJoinedRetryableErrors(t *testing.T) {
	err := errors.Join(
		Ignored(errors.New("metrics flush failed")),
		Retryable(errors.New("redis timeout")),
	)

	if got := Decide(err); got != DecisionRetry {
		t.Fatalf("Decide() = %q, want %q", got, DecisionRetry)
	}
}

func TestDecideReturnsIgnoreForNilOrIgnored(t *testing.T) {
	if got := Decide(nil); got != DecisionIgnore {
		t.Fatalf("Decide(nil) = %q, want %q", got, DecisionIgnore)
	}

	if got := Decide(Ignored(errors.New("noop"))); got != DecisionIgnore {
		t.Fatalf("Decide(ignored) = %q, want %q", got, DecisionIgnore)
	}
}
