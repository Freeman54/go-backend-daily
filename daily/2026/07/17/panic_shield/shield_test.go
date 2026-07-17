package panicshield

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestExecuteReturnsOriginalError(t *testing.T) {
	want := errors.New("db timeout")

	err := Execute(context.Background(), func(context.Context) error {
		return want
	})

	if !errors.Is(err, want) {
		t.Fatalf("expected original error, got %v", err)
	}
}

func TestExecuteWrapsPanic(t *testing.T) {
	err := Execute(context.Background(), func(context.Context) error {
		panic("worker crashed")
	})

	var panicErr *PanicError
	if !errors.As(err, &panicErr) {
		t.Fatalf("expected PanicError, got %T", err)
	}
	if panicErr.Value != "worker crashed" {
		t.Fatalf("unexpected panic value: %#v", panicErr.Value)
	}
	if !strings.Contains(string(panicErr.Stack), "TestExecuteWrapsPanic") {
		t.Fatalf("stack does not include test frame: %s", panicErr.Stack)
	}
}
