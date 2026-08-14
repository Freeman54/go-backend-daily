package attributeallowlist

import (
	"reflect"
	"testing"
)

func TestFilter(t *testing.T) {
	input := map[string]string{
		"traceparent": "00-abc-def-01",
		"tenant":      "t-1",
		"debug":       "internal detail",
	}
	got := Filter(input, []string{"traceparent", "tenant"})
	want := map[string]string{"traceparent": "00-abc-def-01", "tenant": "t-1"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Filter() = %#v, want %#v", got, want)
	}

	input["tenant"] = "changed"
	if got["tenant"] != "t-1" {
		t.Fatal("result should not alias the input map")
	}
}

func TestFilterEmptyInput(t *testing.T) {
	if got := Filter(nil, []string{"traceparent"}); len(got) != 0 {
		t.Fatalf("Filter(nil) = %#v", got)
	}
	if got := Filter(map[string]string{"traceparent": "value"}, nil); len(got) != 0 {
		t.Fatalf("Filter without allowlist = %#v", got)
	}
}
