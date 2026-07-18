package postcommithooks

import (
	"errors"
	"testing"
)

func TestCommitRunsHooksInOrder(t *testing.T) {
	var got []string
	var runner Runner
	runner.Add(func() error {
		got = append(got, "event")
		return nil
	})
	runner.Add(func() error {
		got = append(got, "cache")
		return nil
	})

	errs := runner.Commit()

	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
	if len(got) != 2 || got[0] != "event" || got[1] != "cache" {
		t.Fatalf("unexpected order: %v", got)
	}
}

func TestRollbackDropsPendingHooks(t *testing.T) {
	called := false
	var runner Runner
	runner.Add(func() error {
		called = true
		return nil
	})

	runner.Rollback()
	runner.Commit()

	if called {
		t.Fatal("hook should not run after rollback")
	}
}

func TestCommitCollectsHookErrorsAndClearsState(t *testing.T) {
	var runner Runner
	want := errors.New("publish failed")
	runner.Add(func() error { return want })

	errs := runner.Commit()
	again := runner.Commit()

	if len(errs) != 1 || !errors.Is(errs[0], want) {
		t.Fatalf("unexpected commit errors: %v", errs)
	}
	if len(again) != 0 {
		t.Fatalf("expected hooks to be cleared, got %v", again)
	}
}
