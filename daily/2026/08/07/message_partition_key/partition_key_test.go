package messagepartitionkey

import "testing"

func TestValidate_AcceptsStableKey(t *testing.T) {
	t.Parallel()
	if err := Validate("tenant-42:order-99", 64); err != nil {
		t.Fatalf("Validate() error = %v, want nil", err)
	}
}

func TestValidate_RejectsWhitespace(t *testing.T) {
	t.Parallel()
	if err := Validate("tenant 42", 64); err == nil {
		t.Fatal("Validate() error = nil, want invalid character error")
	}
}

func TestValidate_RejectsTooLongKey(t *testing.T) {
	t.Parallel()
	if err := Validate("abcdefgh", 7); err == nil {
		t.Fatal("Validate() error = nil, want length error")
	}
}
