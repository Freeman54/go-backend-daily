package insertbatchplan

import "testing"

func TestSplitUsesStricterParamLimit(t *testing.T) {
	got, err := Split(10, 3, 5, 12)
	if err != nil {
		t.Fatalf("split: %v", err)
	}

	want := []Plan{
		{Start: 0, End: 4},
		{Start: 4, End: 8},
		{Start: 8, End: 10},
	}

	if len(got) != len(want) {
		t.Fatalf("got %d plans, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("plan %d = %#v, want %#v", i, got[i], want[i])
		}
	}
}

func TestSplitReturnsNilForNoRows(t *testing.T) {
	got, err := Split(0, 3, 100, 300)
	if err != nil {
		t.Fatalf("split: %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil plan, got %#v", got)
	}
}

func TestSplitRejectsInvalidInputs(t *testing.T) {
	cases := []struct {
		totalRows int
		columns   int
		maxRows   int
		maxParams int
	}{
		{-1, 3, 100, 300},
		{1, 0, 100, 300},
		{1, 3, 0, 300},
		{1, 3, 100, 2},
	}

	for _, tc := range cases {
		if _, err := Split(tc.totalRows, tc.columns, tc.maxRows, tc.maxParams); err == nil {
			t.Fatalf("expected error for %#v", tc)
		}
	}
}
