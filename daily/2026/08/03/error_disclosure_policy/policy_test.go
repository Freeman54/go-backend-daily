package error_disclosure_policy

import (
	"errors"
	"testing"
)

func TestPublicMessage_DisclosesOnlySafeErrors(t *testing.T) {
	safe := Safe("quota exceeded")
	if got := PublicMessage(safe); got != "quota exceeded" {
		t.Fatalf("PublicMessage(safe) = %q", got)
	}
	wrapped := errors.New("dial tcp 10.0.0.8:5432: refused")
	if got := PublicMessage(wrapped); got != "internal server error" {
		t.Fatalf("PublicMessage(private) = %q", got)
	}
}

func TestPublicMessage_FindsWrappedSafeError(t *testing.T) {
	err := WrapForOperation("create order", Safe("inventory unavailable"))
	if got := PublicMessage(err); got != "inventory unavailable" {
		t.Fatalf("PublicMessage() = %q", got)
	}
}
