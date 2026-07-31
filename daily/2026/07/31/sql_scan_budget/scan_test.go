package sqlscanbudget

import (
	"errors"
	"testing"
)

func TestBudget_TracksRowsAndBytes(t *testing.T) {
	b, err := New(3, 10)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if err := b.Consume(4); err != nil {
		t.Fatalf("first Consume() error = %v", err)
	}
	if err := b.Consume(6); err != nil {
		t.Fatalf("second Consume() error = %v", err)
	}
	if gotRows, gotBytes := b.Used(); gotRows != 2 || gotBytes != 10 {
		t.Fatalf("Used() = (%d, %d), want (2, 10)", gotRows, gotBytes)
	}
	if err := b.Consume(1); !errors.Is(err, ErrBudgetExceeded) {
		t.Fatalf("Consume() error = %v, want ErrBudgetExceeded", err)
	}
}

func TestBudget_RejectsInvalidLimitsAndNegativeRowSize(t *testing.T) {
	if _, err := New(0, 10); !errors.Is(err, ErrInvalidBudget) {
		t.Fatalf("New() error = %v, want ErrInvalidBudget", err)
	}
	b, _ := New(1, 1)
	if err := b.Consume(-1); !errors.Is(err, ErrInvalidRowSize) {
		t.Fatalf("Consume() error = %v, want ErrInvalidRowSize", err)
	}
}

func TestBudget_DoesNotCountRejectedRow(t *testing.T) {
	b, _ := New(1, 5)
	if err := b.Consume(6); !errors.Is(err, ErrBudgetExceeded) {
		t.Fatalf("Consume() error = %v", err)
	}
	if rows, bytes := b.Used(); rows != 0 || bytes != 0 {
		t.Fatalf("Used() after rejection = (%d, %d)", rows, bytes)
	}
}
