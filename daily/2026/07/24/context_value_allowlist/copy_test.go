package contextvalueallowlist

import (
	"context"
	"testing"
)

type key string

func TestCopy_CopiesOnlyAllowedValues(t *testing.T) {
	source := context.WithValue(context.Background(), key("trace"), "t-1")
	source = context.WithValue(source, key("secret"), "password")
	got := Copy(context.Background(), source, []any{key("trace")})
	if got.Value(key("trace")) != "t-1" {
		t.Fatal("allowed trace value was not copied")
	}
	if got.Value(key("secret")) != nil {
		t.Fatal("secret value must not be copied")
	}
}

func TestCopy_PreservesDestinationValues(t *testing.T) {
	destination := context.WithValue(context.Background(), key("tenant"), "acme")
	got := Copy(destination, context.Background(), nil)
	if got.Value(key("tenant")) != "acme" {
		t.Fatal("destination value was lost")
	}
}
