package sagacompensation

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"
)

func TestSagaCompensatesCompletedStepsInReverseOrder(t *testing.T) {
	var calls []string

	steps := []Step{
		{
			Name: "reserve_inventory",
			Do: func(context.Context) error {
				calls = append(calls, "do:inventory")
				return nil
			},
			Compensate: func(context.Context) error {
				calls = append(calls, "undo:inventory")
				return nil
			},
		},
		{
			Name: "charge_payment",
			Do: func(context.Context) error {
				calls = append(calls, "do:payment")
				return errors.New("bank timeout")
			},
		},
	}

	err := Execute(context.Background(), steps)
	if err == nil || !strings.Contains(err.Error(), "charge_payment failed") {
		t.Fatalf("expected step failure, got %v", err)
	}

	want := []string{"do:inventory", "do:payment", "undo:inventory"}
	if !slices.Equal(calls, want) {
		t.Fatalf("unexpected call order: got=%v want=%v", calls, want)
	}
}

func TestSagaJoinsCompensationErrors(t *testing.T) {
	steps := []Step{
		{
			Name: "write_order",
			Do:   func(context.Context) error { return nil },
			Compensate: func(context.Context) error {
				return errors.New("delete order failed")
			},
		},
		{
			Name: "publish_event",
			Do:   func(context.Context) error { return errors.New("broker unavailable") },
		},
	}

	err := Execute(context.Background(), steps)
	if err == nil {
		t.Fatal("expected saga error")
	}

	if !strings.Contains(err.Error(), "publish_event failed") {
		t.Fatalf("expected original failure, got %v", err)
	}

	if !strings.Contains(err.Error(), "compensate write_order failed") {
		t.Fatalf("expected compensation failure, got %v", err)
	}
}
