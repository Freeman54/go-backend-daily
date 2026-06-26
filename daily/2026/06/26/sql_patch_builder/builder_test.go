package sqlpatchbuilder

import "testing"

func TestBuildUpdateIncludesOnlyPresentFields(t *testing.T) {
	t.Parallel()

	builder := NewBuilder("accounts").
		Allow("display_name", "DisplayName").
		Allow("email", "Email").
		Allow("active", "Active")

	patch := Patch{
		"DisplayName": "alice",
		"Active":      true,
	}

	query, args, err := builder.BuildUpdate(patch, "id = ?", 7)
	if err != nil {
		t.Fatalf("BuildUpdate() returned error: %v", err)
	}

	wantQuery := "UPDATE accounts SET active = ?, display_name = ? WHERE id = ?"
	if query != wantQuery {
		t.Fatalf("query = %q, want %q", query, wantQuery)
	}

	wantArgs := []any{true, "alice", 7}
	if len(args) != len(wantArgs) {
		t.Fatalf("args len = %d, want %d", len(args), len(wantArgs))
	}
	for i := range args {
		if args[i] != wantArgs[i] {
			t.Fatalf("args[%d] = %#v, want %#v", i, args[i], wantArgs[i])
		}
	}
}

func TestBuildUpdateRejectsUnknownField(t *testing.T) {
	t.Parallel()

	builder := NewBuilder("accounts").Allow("display_name", "DisplayName")

	_, _, err := builder.BuildUpdate(Patch{"Role": "admin"}, "id = ?", 1)
	if err == nil {
		t.Fatal("BuildUpdate() expected error for unknown field")
	}
}
