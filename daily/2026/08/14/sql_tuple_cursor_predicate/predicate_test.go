package tuplecursor

import (
	"reflect"
	"testing"
)

func TestAfterCompositeAscending(t *testing.T) {
	query, args, err := After([]string{"created_at", "id"}, []any{"2026-08-14", 42}, false)
	if err != nil {
		t.Fatal(err)
	}
	wantQuery := "((created_at > ?) OR (created_at = ? AND id > ?))"
	if query != wantQuery {
		t.Fatalf("query = %q, want %q", query, wantQuery)
	}
	wantArgs := []any{"2026-08-14", "2026-08-14", 42}
	if !reflect.DeepEqual(args, wantArgs) {
		t.Fatalf("args = %#v, want %#v", args, wantArgs)
	}
}

func TestAfterDescending(t *testing.T) {
	query, args, err := After([]string{"score"}, []any{99}, true)
	if err != nil || query != "((score < ?))" || !reflect.DeepEqual(args, []any{99}) {
		t.Fatalf("query=%q args=%v err=%v", query, args, err)
	}
}

func TestAfterRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		columns []string
		values  []any
	}{
		{nil, nil},
		{[]string{"id"}, nil},
		{[]string{"id; DROP TABLE users"}, []any{1}},
		{[]string{"9id"}, []any{1}},
	}
	for _, tt := range tests {
		if _, _, err := After(tt.columns, tt.values, false); err == nil {
			t.Fatalf("After(%v, %v) should fail", tt.columns, tt.values)
		}
	}
}
