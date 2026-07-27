package result_order_buffer

import "testing"

func TestBuffer_AddReleasesContiguousResults(t *testing.T) {
	b := New[string](10)
	if got := b.Add(11, "eleven"); len(got) != 0 {
		t.Fatalf("out-of-order add released %v", got)
	}
	got := b.Add(10, "ten")
	if len(got) != 2 || got[0] != "ten" || got[1] != "eleven" {
		t.Fatalf("released %v, want [ten eleven]", got)
	}
}

func TestBuffer_AddIgnoresOldAndDuplicateSequence(t *testing.T) {
	b := New[string](2)
	b.Add(3, "first")
	b.Add(3, "duplicate")
	b.Add(1, "old")
	got := b.Add(2, "two")
	if len(got) != 2 || got[1] != "first" {
		t.Fatalf("released %v, want [two first]", got)
	}
}
