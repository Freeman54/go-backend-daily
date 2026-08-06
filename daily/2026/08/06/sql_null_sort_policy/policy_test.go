package sql_null_sort_policy

import "testing"

func TestBuildOrderBy_UsesWhitelistedColumn(t *testing.T) {
	got, err := BuildOrderBy("updated", Descending, NullsLast, map[string]string{"updated": "updated_at"})
	if err != nil {
		t.Fatal(err)
	}
	if want := "ORDER BY updated_at DESC NULLS LAST"; got != want {
		t.Fatalf("BuildOrderBy() = %q, want %q", got, want)
	}
}

func TestBuildOrderBy_RejectsUnknownColumn(t *testing.T) {
	if _, err := BuildOrderBy("updated_at; DROP TABLE users", Ascending, NullsFirst, map[string]string{"updated": "updated_at"}); err == nil {
		t.Fatal("BuildOrderBy() should reject unknown columns")
	}
}

func TestBuildOrderBy_OmitsDatabaseDefaultNullPlacement(t *testing.T) {
	got, err := BuildOrderBy("name", Ascending, NullsDefault, map[string]string{"name": "display_name"})
	if err != nil {
		t.Fatal(err)
	}
	if want := "ORDER BY display_name ASC"; got != want {
		t.Fatalf("BuildOrderBy() = %q, want %q", got, want)
	}
}
