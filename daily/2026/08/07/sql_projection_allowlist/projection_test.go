package sqlprojectionallowlist

import "testing"

func TestSelect_BuildsQuotedProjection(t *testing.T) {
	t.Parallel()
	got, err := Select([]string{"id", "created_at"}, map[string]struct{}{"id": {}, "created_at": {}})
	if err != nil {
		t.Fatalf("Select() error = %v", err)
	}
	if want := "\"id\", \"created_at\""; got != want {
		t.Errorf("Select() = %q, want %q", got, want)
	}
}

func TestSelect_RejectsUnknownColumn(t *testing.T) {
	t.Parallel()
	if _, err := Select([]string{"id", "password_hash"}, map[string]struct{}{"id": {}}); err == nil {
		t.Fatal("Select() error = nil, want allowlist error")
	}
}

func TestSelect_RejectsDuplicateColumn(t *testing.T) {
	t.Parallel()
	if _, err := Select([]string{"id", "id"}, map[string]struct{}{"id": {}}); err == nil {
		t.Fatal("Select() error = nil, want duplicate error")
	}
}
