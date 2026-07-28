package sql_placeholder_rebinder

import (
	"errors"
	"testing"
)

func TestRebind_ReplacesOnlyPlaceholdersOutsideQuotedText(t *testing.T) {
	query := "SELECT '?' AS literal, name FROM users WHERE id = ? AND note = 'it''s ?'"
	got, count, err := Rebind(query)
	if err != nil {
		t.Fatal(err)
	}
	want := "SELECT '?' AS literal, name FROM users WHERE id = $1 AND note = 'it''s ?'"
	if got != want || count != 1 {
		t.Fatalf("Rebind() = (%q, %d), want (%q, 1)", got, count, want)
	}
}

func TestRebind_ReplacesSequentialPlaceholders(t *testing.T) {
	got, count, err := Rebind("WHERE a = ? OR b = ?")
	if err != nil || got != "WHERE a = $1 OR b = $2" || count != 2 {
		t.Fatalf("Rebind() = (%q, %d, %v)", got, count, err)
	}
}

func TestRebind_RejectsUnterminatedQuote(t *testing.T) {
	_, _, err := Rebind("SELECT 'broken ?")
	if !errors.Is(err, ErrUnterminatedQuote) {
		t.Fatalf("error = %v, want ErrUnterminatedQuote", err)
	}
}
