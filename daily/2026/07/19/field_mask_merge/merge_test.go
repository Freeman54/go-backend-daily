package fieldmaskmerge

import "testing"

func TestApplyOnlyTouchesMaskedFields(t *testing.T) {
	base := map[string]string{
		"name":   "alice",
		"email":  "a@example.com",
		"region": "hz",
	}
	email := "new@example.com"

	got, err := Apply(base, map[string]*string{
		"email": &email,
	}, []string{"email"}, allowed("name", "email", "region"))
	if err != nil {
		t.Fatalf("apply: %v", err)
	}

	if got["email"] != "new@example.com" {
		t.Fatalf("expected email to change, got %q", got["email"])
	}
	if got["name"] != "alice" || got["region"] != "hz" {
		t.Fatalf("unexpected side effect: %#v", got)
	}
}

func TestApplyDeletesFieldWhenUpdateIsNil(t *testing.T) {
	base := map[string]string{"nickname": "ops"}

	got, err := Apply(base, map[string]*string{
		"nickname": nil,
	}, []string{"nickname"}, allowed("nickname"))
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	if _, ok := got["nickname"]; ok {
		t.Fatalf("expected nickname to be deleted, got %#v", got)
	}
}

func TestApplyRejectsUnknownField(t *testing.T) {
	_, err := Apply(nil, map[string]*string{}, []string{"role"}, allowed("name"))
	if err == nil {
		t.Fatal("expected unknown field to fail")
	}
}

func TestApplyRejectsMissingUpdatePayload(t *testing.T) {
	_, err := Apply(nil, map[string]*string{}, []string{"name"}, allowed("name"))
	if err == nil {
		t.Fatal("expected missing payload to fail")
	}
}

func allowed(keys ...string) map[string]struct{} {
	set := make(map[string]struct{}, len(keys))
	for _, key := range keys {
		set[key] = struct{}{}
	}
	return set
}
