package queryshapehash

import "testing"

func TestNormalizeMasksLiterals(t *testing.T) {
	queryA := "SELECT * FROM orders WHERE user_id = 42 AND state = 'paid'"
	queryB := " select *   from orders where user_id = 99 and state = 'failed' "

	gotA := Normalize(queryA)
	gotB := Normalize(queryB)

	if gotA != gotB {
		t.Fatalf("normalized query mismatch:\n%s\n%s", gotA, gotB)
	}
}

func TestHashChangesWhenShapeChanges(t *testing.T) {
	queryA := "SELECT * FROM orders WHERE user_id = 42"
	queryB := "SELECT * FROM orders WHERE user_id = 42 ORDER BY created_at DESC"

	if Hash(queryA) == Hash(queryB) {
		t.Fatal("expected different hashes for different query shapes")
	}
}
