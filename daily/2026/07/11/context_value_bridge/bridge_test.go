package contextvaluebridge

import (
	"context"
	"testing"
	"time"
)

type ctxKey string

func TestDetachCopiesSelectedValues(t *testing.T) {
	parent := context.Background()
	parent = context.WithValue(parent, ctxKey("request_id"), "req-1")
	parent = context.WithValue(parent, ctxKey("user_id"), "u-42")

	child, cancel := Detach(parent, time.Second, ctxKey("request_id"))
	defer cancel()

	if got := child.Value(ctxKey("request_id")); got != "req-1" {
		t.Fatalf("request_id = %v, want req-1", got)
	}
	if got := child.Value(ctxKey("user_id")); got != nil {
		t.Fatalf("user_id should not be copied, got %v", got)
	}
}

func TestDetachDoesNotInheritParentCancellation(t *testing.T) {
	parent, parentCancel := context.WithCancel(context.Background())
	parent = context.WithValue(parent, ctxKey("trace_id"), "trace-9")

	child, cancel := Detach(parent, time.Second, ctxKey("trace_id"))
	defer cancel()

	parentCancel()

	select {
	case <-child.Done():
		t.Fatal("child context should stay alive after parent cancellation")
	default:
	}

	if got := child.Value(ctxKey("trace_id")); got != "trace-9" {
		t.Fatalf("trace_id = %v, want trace-9", got)
	}
}

func TestDetachAppliesFreshTimeout(t *testing.T) {
	child, cancel := Detach(context.Background(), 30*time.Millisecond)
	defer cancel()

	select {
	case <-child.Done():
	case <-time.After(200 * time.Millisecond):
		t.Fatal("child context did not time out")
	}

	if err := child.Err(); err != context.DeadlineExceeded {
		t.Fatalf("child err = %v, want deadline exceeded", err)
	}
}
