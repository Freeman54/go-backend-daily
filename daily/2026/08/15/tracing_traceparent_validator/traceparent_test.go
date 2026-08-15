package tracingtraceparentvalidator

import (
	"errors"
	"testing"
)

func TestParse(t *testing.T) {
	p, err := Parse("00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01")
	if err != nil {
		t.Fatal(err)
	}
	if !p.Sampled() || p.Version != 0 || p.TraceID[0] != 0x4b || p.SpanID[0] != 0x00 {
		t.Fatalf("parent = %#v", p)
	}
}

func TestParseRejectsMalformedParents(t *testing.T) {
	tests := []string{
		"", "00-short-00f067aa0ba902b7-01",
		"01-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
		"00-00000000000000000000000000000000-00f067aa0ba902b7-01",
		"00-4bf92f3577b34da6a3ce929d0e0e4736-0000000000000000-01",
		"00-4BF92F3577B34DA6A3CE929D0E0E4736-00f067aa0ba902b7-01",
		"00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-zz",
	}
	for _, value := range tests {
		if _, err := Parse(value); !errors.Is(err, ErrInvalidTraceparent) {
			t.Fatalf("Parse(%q) error = %v", value, err)
		}
	}
}

func TestSampledBit(t *testing.T) {
	p, err := Parse("00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-02")
	if err != nil || p.Sampled() {
		t.Fatalf("parent = %#v, error = %v", p, err)
	}
}
