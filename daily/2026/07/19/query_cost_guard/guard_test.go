package querycostguard

import (
	"context"
	"testing"
	"time"
)

func TestAcquireAndReleaseWeightedBudget(t *testing.T) {
	guard := New(3)

	releaseHeavy, err := guard.Acquire(context.Background(), 2)
	if err != nil {
		t.Fatalf("acquire heavy: %v", err)
	}
	releaseLight, err := guard.Acquire(context.Background(), 1)
	if err != nil {
		t.Fatalf("acquire light: %v", err)
	}

	done := make(chan error, 1)
	go func() {
		_, err := guard.Acquire(context.Background(), 1)
		done <- err
	}()

	select {
	case err := <-done:
		t.Fatalf("third query should block, got %v", err)
	case <-time.After(40 * time.Millisecond):
	}

	releaseLight()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("third query failed: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("third query did not resume after release")
	}

	releaseHeavy()
}

func TestAcquireRejectsOverBudgetQuery(t *testing.T) {
	guard := New(2)
	if _, err := guard.Acquire(context.Background(), 3); err == nil {
		t.Fatal("expected over-budget query to fail")
	}
}

func TestAcquireRespectsContextTimeout(t *testing.T) {
	guard := New(1)
	release, err := guard.Acquire(context.Background(), 1)
	if err != nil {
		t.Fatalf("first acquire: %v", err)
	}
	defer release()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	if _, err := guard.Acquire(ctx, 1); err == nil {
		t.Fatal("expected waiting query to time out")
	}
}
