package consistent_hash_bounded_load

import "testing"

func TestRing_PickSkipsOverloadedPrimary(t *testing.T) {
	r, err := New([]string{"a", "b", "c"}, 64, 1.25)
	if err != nil {
		t.Fatal(err)
	}
	primary, err := r.Pick("order-42", nil)
	if err != nil {
		t.Fatal(err)
	}
	loads := map[string]int{"a": 1, "b": 1, "c": 1}
	loads[primary] = 100
	got, err := r.Pick("order-42", loads)
	if err != nil {
		t.Fatal(err)
	}
	if got == primary {
		t.Fatalf("Pick() returned overloaded primary %q", primary)
	}
}

func TestNew_RejectsInvalidConfiguration(t *testing.T) {
	if _, err := New(nil, 10, 1.25); err == nil {
		t.Fatal("New() should reject an empty node set")
	}
}
