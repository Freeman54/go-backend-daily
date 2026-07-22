package consumerlagwindow

import "testing"

func TestWindowTracksRollingMaximumAndAverage(t *testing.T) {
	w, err := New(3)
	if err != nil {
		t.Fatalf("new window: %v", err)
	}
	for _, lag := range []int64{10, 20, 30, 5} {
		if err := w.Add(lag); err != nil {
			t.Fatalf("add: %v", err)
		}
	}
	s := w.Snapshot()
	if s.Count != 3 || s.Max != 30 || s.Average != 18 {
		t.Fatalf("snapshot = %#v", s)
	}
}

func TestWindowRejectsInvalidValues(t *testing.T) {
	if _, err := New(0); err == nil {
		t.Fatal("expected invalid size error")
	}
	w, _ := New(1)
	if err := w.Add(-1); err == nil {
		t.Fatal("expected negative lag error")
	}
	if got := w.Snapshot(); got != (Snapshot{}) {
		t.Fatalf("empty snapshot = %#v", got)
	}
}
