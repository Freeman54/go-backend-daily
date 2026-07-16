package baggageallowlist

import "testing"

func TestFilterKeepsAllowedEntries(t *testing.T) {
	allowed := map[string]struct{}{
		"tenant": {},
		"trace":  {},
	}

	got := Filter("tenant=acme,user=42,trace=req-1", allowed, 64)
	want := "tenant=acme,trace=req-1"
	if got != want {
		t.Fatalf("unexpected filtered baggage: %q", got)
	}
}

func TestFilterRespectsByteBudget(t *testing.T) {
	allowed := map[string]struct{}{
		"tenant": {},
		"trace":  {},
	}

	got := Filter("tenant=acme,trace=request-12345", allowed, len("tenant=acme"))
	if got != "tenant=acme" {
		t.Fatalf("unexpected budgeted baggage: %q", got)
	}
}

func TestFilterSkipsMalformedParts(t *testing.T) {
	allowed := map[string]struct{}{
		"tenant": {},
	}

	got := Filter("tenant=acme,broken-part,=missing-key", allowed, 64)
	if got != "tenant=acme" {
		t.Fatalf("unexpected filtered baggage: %q", got)
	}
}
