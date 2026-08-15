package concurrencyquiescencegate

import (
	"errors"
	"testing"
	"time"
)

func TestCloseWaitsForActiveWork(t *testing.T) {
	gate := New()
	leave, err := gate.Enter()
	if err != nil {
		t.Fatal(err)
	}
	done := gate.Close()
	select {
	case <-done:
		t.Fatal("closed before active work left")
	default:
	}
	if _, err := gate.Enter(); !errors.Is(err, ErrClosed) {
		t.Fatalf("Enter error = %v", err)
	}
	leave()
	leave()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("close did not finish")
	}
	if gate.Active() != 0 {
		t.Fatalf("active = %d", gate.Active())
	}
}

func TestCloseWithoutWorkIsImmediateAndIdempotent(t *testing.T) {
	gate := New()
	first := gate.Close()
	second := gate.Close()
	if first != second {
		t.Fatal("Close returned different channels")
	}
	select {
	case <-first:
	default:
		t.Fatal("empty gate did not close immediately")
	}
}
