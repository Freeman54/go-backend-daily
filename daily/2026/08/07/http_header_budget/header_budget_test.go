package httpheaderbudget

import (
	"net/http"
	"testing"
)

func TestValidate_WithinBudget(t *testing.T) {
	t.Parallel()
	headers := http.Header{"Accept": {"application/json"}, "X-Request-ID": {"abc"}}
	if err := Validate(headers, 2, 64); err != nil {
		t.Fatalf("Validate() error = %v, want nil", err)
	}
}

func TestValidate_RejectsTooManyFields(t *testing.T) {
	t.Parallel()
	headers := http.Header{"A": {"1"}, "B": {"2"}}
	if err := Validate(headers, 1, 64); err == nil {
		t.Fatal("Validate() error = nil, want field budget error")
	}
}

func TestValidate_RejectsByteBudget(t *testing.T) {
	t.Parallel()
	headers := http.Header{"X-Trace": {"12345"}}
	if err := Validate(headers, 3, 10); err == nil {
		t.Fatal("Validate() error = nil, want byte budget error")
	}
}
