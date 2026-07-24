package metriccardinalityguard

import "testing"

func TestGuard_AdmitsKnownValuesAndCapsNewValues(t *testing.T) {
	guard, err := New(2, "other")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if got := guard.Normalize("read"); got != "read" {
		t.Fatalf("Normalize(read) = %q", got)
	}
	if got := guard.Normalize("write"); got != "write" {
		t.Fatalf("Normalize(write) = %q", got)
	}
	if got := guard.Normalize("delete"); got != "other" {
		t.Fatalf("Normalize(delete) = %q, want other", got)
	}
	if got := guard.Normalize("read"); got != "read" {
		t.Fatalf("known value changed to %q", got)
	}
}

func TestNew_RejectsInvalidConfiguration(t *testing.T) {
	if _, err := New(0, "other"); err == nil {
		t.Fatal("zero limit should return an error")
	}
	if _, err := New(1, ""); err == nil {
		t.Fatal("empty fallback should return an error")
	}
}
