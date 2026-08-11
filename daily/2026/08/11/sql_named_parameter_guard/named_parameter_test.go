package namedparameter

import (
	"errors"
	"reflect"
	"testing"
)

func TestValidate(t *testing.T) {
	query := `SELECT * FROM orders WHERE tenant_id = :tenant AND id = :id AND note = ':ignored'`
	args := map[string]any{"id": 42, "tenant": "acme"}
	if err := Validate(query, args); err != nil {
		t.Fatal(err)
	}
	if got := Names(args); !reflect.DeepEqual(got, []string{"id", "tenant"}) {
		t.Fatalf("Names() = %v", got)
	}
}

func TestValidateRejectsMissingOrUnusedArguments(t *testing.T) {
	for _, args := range []map[string]any{
		{"tenant": "acme"},
		{"tenant": "acme", "id": 42, "unused": true},
	} {
		err := Validate("SELECT * FROM t WHERE tenant = :tenant AND id = :id", args)
		if !errors.Is(err, ErrArgumentMismatch) {
			t.Fatalf("error = %v", err)
		}
	}
}

func TestValidateAllowsRepeatedNamesAndIdentifiers(t *testing.T) {
	err := Validate("SELECT :value_2 + :value_2", map[string]any{"value_2": 2})
	if err != nil {
		t.Fatal(err)
	}
}
