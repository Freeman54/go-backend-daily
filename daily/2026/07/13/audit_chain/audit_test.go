package auditchain

import (
	"errors"
	"testing"
	"time"
)

func TestVerifyPassesForUntouchedChain(t *testing.T) {
	var trail Trail
	trail.Append(Event{
		ID:       "1",
		Actor:    "alice",
		Action:   "grant",
		Resource: "role/admin",
		At:       time.Unix(100, 0),
	})
	trail.Append(Event{
		ID:       "2",
		Actor:    "bob",
		Action:   "revoke",
		Resource: "role/admin",
		At:       time.Unix(110, 0),
	})

	if err := trail.Verify(); err != nil {
		t.Fatalf("Verify failed: %v", err)
	}
}

func TestVerifyDetectsTamper(t *testing.T) {
	var trail Trail
	trail.Append(Event{
		ID:       "1",
		Actor:    "alice",
		Action:   "grant",
		Resource: "wallet/42",
		At:       time.Unix(200, 0),
	})
	trail.Append(Event{
		ID:       "2",
		Actor:    "alice",
		Action:   "transfer",
		Resource: "wallet/42",
		At:       time.Unix(210, 0),
	})

	events := trail.Events()
	events[1].Action = "delete"

	if err := VerifyEvents(events); !errors.Is(err, ErrBrokenChain) {
		t.Fatalf("VerifyEvents error = %v want ErrBrokenChain", err)
	}
}
