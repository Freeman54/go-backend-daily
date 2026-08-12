package identifierallowlist

import (
	"errors"
	"testing"
)

func TestAllowlistResolve(t *testing.T) {
	a, err := New(map[string]string{"CreatedAt": "orders.created_at", "id": "order_id"})
	if err != nil {
		t.Fatal(err)
	}
	got, err := a.Resolve(" createdat ")
	if err != nil || got != "orders.created_at" {
		t.Fatalf("resolve = %q, %v", got, err)
	}
	if _, err := a.Resolve("created_at DESC"); !errors.Is(err, ErrIdentifierNotAllowed) {
		t.Fatalf("unknown error = %v", err)
	}
}

func TestAllowlistRejectsUnsafeMappings(t *testing.T) {
	bad := []string{"", "orders..id", "id DESC", "1column", "id;DROP"}
	for _, identifier := range bad {
		if _, err := New(map[string]string{"field": identifier}); !errors.Is(err, ErrIdentifierNotAllowed) {
			t.Fatalf("identifier %q error = %v", identifier, err)
		}
	}
	if _, err := New(map[string]string{" ": "id"}); !errors.Is(err, ErrIdentifierNotAllowed) {
		t.Fatalf("empty external name error = %v", err)
	}
}
