package optionalfieldpatch

import (
	"encoding/json"
	"testing"
)

func TestFieldDistinguishesMissingNullAndConcreteValue(t *testing.T) {
	var patch Patch
	if err := json.Unmarshal([]byte(`{"nickname":null,"age":18}`), &patch); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if !patch.Nickname.Present() || !patch.Nickname.Null() {
		t.Fatal("nickname should be present and null")
	}
	if !patch.Age.Present() || patch.Age.Null() {
		t.Fatal("age should be present with concrete value")
	}
	value, err := patch.Age.Value()
	if err != nil || value != 18 {
		t.Fatalf("Age.Value() = (%d, %v), want (18, nil)", value, err)
	}
}

func TestApplyOnlyTouchesPresentFields(t *testing.T) {
	user := User{Nickname: "alice", Age: 20}

	var patch Patch
	if err := json.Unmarshal([]byte(`{"nickname":"bob"}`), &patch); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	updated := Apply(user, patch)
	if updated.Nickname != "bob" {
		t.Fatalf("Nickname = %q, want %q", updated.Nickname, "bob")
	}
	if updated.Age != 20 {
		t.Fatalf("Age = %d, want 20", updated.Age)
	}
}

func TestApplyClearsFieldOnNull(t *testing.T) {
	user := User{Nickname: "alice", Age: 20}

	var patch Patch
	if err := json.Unmarshal([]byte(`{"nickname":null,"age":null}`), &patch); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	updated := Apply(user, patch)
	if updated.Nickname != "" || updated.Age != 0 {
		t.Fatalf("Apply() = %+v, want zeroed fields", updated)
	}
}
