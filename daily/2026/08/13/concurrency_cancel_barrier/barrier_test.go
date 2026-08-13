package cancelbarrier

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRunCancelsSiblingsAndWaits(t *testing.T) {
	want := errors.New("failed")
	cleaned := make(chan struct{})
	err := Run(context.Background(), func(context.Context) error { return want }, func(ctx context.Context) error { <-ctx.Done(); close(cleaned); return nil })
	if !errors.Is(err, want) {
		t.Fatalf("error = %v", err)
	}
	select {
	case <-cleaned:
	default:
		t.Fatal("returned before sibling cleanup")
	}
}
func TestRunSuccessAndParentCancellation(t *testing.T) {
	if err := Run(context.Background(), func(context.Context) error { return nil }); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	err := Run(ctx, func(ctx context.Context) error { <-ctx.Done(); return nil })
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("error = %v", err)
	}
}
