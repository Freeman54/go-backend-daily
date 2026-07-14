package deadlinepartition

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestPartitionSplitsRemainingTimeByWeight(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Second))
	defer cancel()

	budgets, err := Partition(ctx, 100*time.Millisecond, []Step{
		{Name: "cache", Weight: 1},
		{Name: "db", Weight: 2},
		{Name: "fallback", Weight: 1},
	})
	if err != nil {
		t.Fatalf("Partition error = %v", err)
	}

	if len(budgets) != 3 {
		t.Fatalf("budgets len = %d want 3", len(budgets))
	}

	total := budgets[0].Timeout + budgets[1].Timeout + budgets[2].Timeout
	if total < 750*time.Millisecond || total > 920*time.Millisecond {
		t.Fatalf("total timeout = %v want around 900ms", total)
	}
	if budgets[1].Timeout <= budgets[0].Timeout {
		t.Fatalf("db timeout = %v should be larger than cache timeout = %v", budgets[1].Timeout, budgets[0].Timeout)
	}
}

func TestPartitionRequiresDeadline(t *testing.T) {
	_, err := Partition(context.Background(), 0, []Step{{Name: "db", Weight: 1}})
	if !errors.Is(err, ErrNoDeadline) {
		t.Fatalf("Partition error = %v want ErrNoDeadline", err)
	}
}
