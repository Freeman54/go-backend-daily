package json_depth_guard

import "testing"

func TestValidate_AllowsValuesWithinDepth(t *testing.T) {
	t.Parallel()
	if err := Validate([]byte(`{"items":[{"id":1}]}`), 3); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestValidate_RejectsValuesBeyondDepth(t *testing.T) {
	t.Parallel()
	if err := Validate([]byte(`{"items":[{"meta":{"id":1}}]}`), 3); err != ErrTooDeep {
		t.Fatalf("Validate() error = %v, want %v", err, ErrTooDeep)
	}
}

func TestValidate_RejectsInvalidInput(t *testing.T) {
	t.Parallel()
	if err := Validate([]byte(`{"id":`), 3); err == nil {
		t.Fatal("Validate() error = nil, want syntax error")
	}
}

func TestValidate_RejectsInvalidLimit(t *testing.T) {
	t.Parallel()
	if err := Validate([]byte(`{}`), 0); err != ErrInvalidMaxDepth {
		t.Fatalf("Validate() error = %v, want %v", err, ErrInvalidMaxDepth)
	}
}
