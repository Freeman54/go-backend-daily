package messagesequenceguard

import (
	"errors"
	"testing"
)

func TestGuardAcceptsNextSequenceAndDetectsDuplicate(t *testing.T) {
	g := New()
	if result, err := g.Accept("order:7", 1); err != nil || result != Applied {
		t.Fatalf("first Accept() = %v, %v; want Applied", result, err)
	}
	if result, err := g.Accept("order:7", 1); err != nil || result != Duplicate {
		t.Fatalf("duplicate Accept() = %v, %v; want Duplicate", result, err)
	}
	if result, err := g.Accept("order:7", 2); err != nil || result != Applied {
		t.Fatalf("next Accept() = %v, %v; want Applied", result, err)
	}
}

func TestGuardRejectsSequenceGapWithoutAdvancing(t *testing.T) {
	g := New()
	if _, err := g.Accept("order:7", 2); !errors.Is(err, ErrSequenceGap) {
		t.Fatalf("gap error = %v, want %v", err, ErrSequenceGap)
	}
	if result, err := g.Accept("order:7", 1); err != nil || result != Applied {
		t.Fatalf("Accept() after gap = %v, %v; want Applied", result, err)
	}
}

func TestGuardTracksKeysIndependentlyAndRejectsInvalidSequence(t *testing.T) {
	g := New()
	if _, err := g.Accept("a", 0); err == nil {
		t.Fatal("expected invalid sequence error")
	}
	if result, err := g.Accept("a", 1); err != nil || result != Applied {
		t.Fatalf("key a = %v, %v", result, err)
	}
	if result, err := g.Accept("b", 1); err != nil || result != Applied {
		t.Fatalf("key b = %v, %v", result, err)
	}
}
