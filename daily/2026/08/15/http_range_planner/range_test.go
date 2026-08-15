package httprangeplanner

import (
	"errors"
	"testing"
)

func TestPlan(t *testing.T) {
	tests := []struct {
		header string
		size   int64
		want   Range
	}{
		{"bytes=2-5", 10, Range{2, 5}},
		{"bytes=7-", 10, Range{7, 9}},
		{"bytes=-3", 10, Range{7, 9}},
		{"bytes=-20", 10, Range{0, 9}},
		{"bytes=8-20", 10, Range{8, 9}},
	}
	for _, tt := range tests {
		got, err := Plan(tt.header, tt.size)
		if err != nil || got != tt.want {
			t.Fatalf("Plan(%q, %d) = %#v, %v", tt.header, tt.size, got, err)
		}
	}
}

func TestPlanRejectsInvalidAndUnsatisfiableRanges(t *testing.T) {
	for _, header := range []string{"items=0-1", "bytes=", "bytes=1-2,4-5", "bytes=-0", "bytes=4-2", "bytes=+1-2"} {
		if _, err := Plan(header, 10); !errors.Is(err, ErrInvalidRange) {
			t.Fatalf("Plan(%q) error = %v", header, err)
		}
	}
	if _, err := Plan("bytes=10-", 10); !errors.Is(err, ErrRangeNotSatisfy) {
		t.Fatalf("error = %v", err)
	}
	if _, err := Plan("bytes=0-", 0); !errors.Is(err, ErrRangeNotSatisfy) {
		t.Fatalf("empty resource error = %v", err)
	}
}

func TestRangeMetadata(t *testing.T) {
	r := Range{Start: 2, End: 5}
	if r.Length() != 4 || ContentRange(r, 10) != "bytes 2-5/10" {
		t.Fatalf("unexpected metadata")
	}
}
