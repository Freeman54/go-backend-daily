package withoutcancelcleanup

import (
	"context"
	"testing"
	"time"
)

type key string

func TestDetachKeepsValuesButIgnoresParentCancel(t *testing.T) {
	parent, cancelParent := context.WithCancel(context.WithValue(context.Background(), key("request_id"), "req-9"))
	detached, cancel := Detach(parent, 0)
	defer cancel()

	cancelParent()

	if got := detached.Value(key("request_id")); got != "req-9" {
		t.Fatalf("expected request_id to be preserved, got %v", got)
	}
	select {
	case <-detached.Done():
		t.Fatal("detached context should ignore parent cancellation")
	default:
	}
}

func TestDetachAddsOwnTimeout(t *testing.T) {
	parent, cancelParent := context.WithCancel(context.Background())
	defer cancelParent()

	detached, cancel := Detach(parent, 20*time.Millisecond)
	defer cancel()

	select {
	case <-detached.Done():
		t.Fatal("timeout should not fire immediately")
	case <-time.After(5 * time.Millisecond):
	}

	select {
	case <-detached.Done():
	case <-time.After(100 * time.Millisecond):
		t.Fatal("expected detached context to timeout")
	}
}
