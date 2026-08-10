package requesttargetguard

import (
	"errors"
	"testing"
)

func TestNormalize(t *testing.T) {
	got, err := Normalize("/v1/orders?status=paid&page=2&tag=b&tag=a", 128)
	if err != nil {
		t.Fatal(err)
	}
	const want = "/v1/orders?page=2&status=paid&tag=b&tag=a"
	if got != want {
		t.Fatalf("Normalize() = %q, want %q", got, want)
	}
}

func TestNormalizeRejectsUnsafeTargets(t *testing.T) {
	tests := []struct {
		name   string
		target string
		max    int
		want   error
	}{
		{"absolute URL", "https://example.com/a", 128, ErrInvalid},
		{"traversal", "/a/../admin", 128, ErrInvalid},
		{"encoded traversal", "/a/%2e%2e/admin", 128, ErrInvalid},
		{"fragment", "/a#private", 128, ErrInvalid},
		{"backslash", `/a\b`, 128, ErrInvalid},
		{"too long", "/abcdef", 4, ErrTooLong},
		{"invalid budget", "/", 0, ErrTooLong},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Normalize(tt.target, tt.max)
			if !errors.Is(err, tt.want) {
				t.Fatalf("error = %v, want %v", err, tt.want)
			}
		})
	}
}
