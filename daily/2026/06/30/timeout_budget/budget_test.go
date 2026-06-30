package timeoutbudget

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestSplitAllocatesByWeight(t *testing.T) {
	t.Parallel()

	base := time.Unix(0, 0)
	ctx, cancel := context.WithDeadline(context.Background(), base.Add(110*time.Millisecond))
	defer cancel()

	got, err := Split(ctx, base, []Step{
		{Name: "cache", Weight: 1, MinShare: 10 * time.Millisecond},
		{Name: "db", Weight: 3, MinShare: 20 * time.Millisecond},
	})
	if err != nil {
		t.Fatalf("Split() error = %v", err)
	}

	if got[0].Timeout != 30*time.Millisecond {
		t.Fatalf("cache timeout = %v, want 30ms", got[0].Timeout)
	}
	if got[1].Timeout != 80*time.Millisecond {
		t.Fatalf("db timeout = %v, want 80ms", got[1].Timeout)
	}
}

func TestSplitRequiresDeadline(t *testing.T) {
	t.Parallel()

	_, err := Split(context.Background(), time.Unix(0, 0), []Step{{Name: "db", Weight: 1}})
	if !errors.Is(err, ErrNoDeadline) {
		t.Fatalf("Split() error = %v, want ErrNoDeadline", err)
	}
}

func TestSplitFailsWhenBudgetTooSmall(t *testing.T) {
	t.Parallel()

	base := time.Unix(0, 0)
	ctx, cancel := context.WithDeadline(context.Background(), base.Add(30*time.Millisecond))
	defer cancel()

	_, err := Split(ctx, base, []Step{
		{Name: "cache", Weight: 1, MinShare: 20 * time.Millisecond},
		{Name: "db", Weight: 1, MinShare: 20 * time.Millisecond},
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Split() error = %v, want deadline exceeded", err)
	}
}
