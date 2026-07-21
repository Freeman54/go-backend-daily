package errorsignalmap

import (
	"errors"
	"testing"
)

func TestMapRecognizesWrappedErrors(t *testing.T) {
	err := Wrapf(ErrTimeout, "fetch profile for user %d", 42)

	got := Map(err)
	want := Signal{
		Code:      "timeout",
		Level:     "error",
		Retryable: true,
		Metric:    "dependency_timeout_total",
	}

	if got != want {
		t.Fatalf("signal = %#v, want %#v", got, want)
	}
}

func TestMapHandlesKnownClasses(t *testing.T) {
	cases := []struct {
		name string
		err  error
		code string
	}{
		{name: "nil", err: nil, code: "ok"},
		{name: "invalid", err: ErrInvalid, code: "invalid_argument"},
		{name: "conflict", err: ErrConflict, code: "conflict"},
		{name: "unavailable", err: ErrUnavailable, code: "unavailable"},
		{name: "internal", err: errors.New("boom"), code: "internal"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Map(tc.err); got.Code != tc.code {
				t.Fatalf("code = %s, want %s", got.Code, tc.code)
			}
		})
	}
}

func TestWrapfReturnsNilForNilError(t *testing.T) {
	if err := Wrapf(nil, "ignored"); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}
