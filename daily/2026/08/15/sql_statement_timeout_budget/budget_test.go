package sqlstatementtimeoutbudget

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestBudgetUsesSmallerLimit(t *testing.T) {
	now := time.Unix(100, 0)
	ctx, cancel := context.WithDeadline(context.Background(), now.Add(8500*time.Millisecond))
	defer cancel()
	got, err := Budget(ctx, now, 10*time.Second, 1500*time.Millisecond)
	if err != nil || got != 7*time.Second {
		t.Fatalf("Budget = %v, %v", got, err)
	}
	got, err = Budget(context.Background(), now, 3*time.Second, time.Second)
	if err != nil || got != 3*time.Second {
		t.Fatalf("maximum Budget = %v, %v", got, err)
	}
}

func TestBudgetRejectsExhaustedOrInvalidLimits(t *testing.T) {
	now := time.Unix(100, 0)
	ctx, cancel := context.WithDeadline(context.Background(), now.Add(500*time.Millisecond))
	defer cancel()
	for _, tc := range []struct {
		ctx     context.Context
		max     time.Duration
		reserve time.Duration
	}{
		{ctx, time.Second, time.Second},
		{context.Background(), 0, 0},
		{context.Background(), time.Second, -1},
	} {
		if _, err := Budget(tc.ctx, now, tc.max, tc.reserve); !errors.Is(err, ErrNoBudget) {
			t.Fatalf("error = %v", err)
		}
	}
}

func TestSetLocalSQL(t *testing.T) {
	query, args, err := SetLocalSQL(2500 * time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	if query != "SELECT set_config('statement_timeout', $1, true)" || !reflect.DeepEqual(args, []any{"2500ms"}) {
		t.Fatalf("query = %q, args = %#v", query, args)
	}
	if _, _, err := SetLocalSQL(time.Microsecond); !errors.Is(err, ErrNoBudget) {
		t.Fatalf("error = %v", err)
	}
}
