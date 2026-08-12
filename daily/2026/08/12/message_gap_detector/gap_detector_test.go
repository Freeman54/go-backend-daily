package gapdetector

import "testing"

func TestDetectorDoesNotAdvanceAcrossGap(t *testing.T) {
	d := New(10)
	if got := d.Observe(12); got.Kind != Gap || got.Missing != 2 || got.Expected != 10 {
		t.Fatalf("gap = %+v", got)
	}
	if got := d.Observe(10); got.Kind != Accepted {
		t.Fatalf("expected offset = %+v", got)
	}
	if got := d.Observe(10); got.Kind != Duplicate || got.Expected != 11 {
		t.Fatalf("duplicate = %+v", got)
	}
	if got := d.Observe(11); got.Kind != Accepted {
		t.Fatalf("next offset = %+v", got)
	}
}

func TestDetectorAcceptsContiguousOffsets(t *testing.T) {
	d := New(0)
	for i := int64(0); i < 3; i++ {
		if got := d.Observe(i); got.Kind != Accepted || got.Expected != i {
			t.Fatalf("offset %d = %+v", i, got)
		}
	}
}
