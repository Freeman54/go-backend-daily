package deadlineslackguard

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestCheck_RejectsMissingOrInsufficientDeadline(t *testing.T) {
	now := time.Date(2026, 7, 31, 10, 0, 0, 0, time.UTC)
	if err := Check(context.Background(), now, time.Second); !errors.Is(err, ErrNoDeadline) {
		t.Fatalf("Check() error = %v, want ErrNoDeadline", err)
	}

	ctx, cancel := context.WithDeadline(context.Background(), now.Add(500*time.Millisecond))
	defer cancel()
	if err := Check(ctx, now, time.Second); !errors.Is(err, ErrInsufficientSlack) {
		t.Fatalf("Check() error = %v, want ErrInsufficientSlack", err)
	}
}

func TestCheck_AcceptsEnoughSlackAndValidatesMinimum(t *testing.T) {
	now := time.Date(2026, 7, 31, 10, 0, 0, 0, time.UTC)
	ctx, cancel := context.WithDeadline(context.Background(), now.Add(2*time.Second))
	defer cancel()
	if err := Check(ctx, now, time.Second); err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	if err := Check(ctx, now, 0); !errors.Is(err, ErrInvalidMinimum) {
		t.Fatalf("Check() error = %v, want ErrInvalidMinimum", err)
	}
}
