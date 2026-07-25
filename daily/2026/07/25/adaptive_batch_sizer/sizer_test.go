package adaptive_batch_sizer

import "testing"

func TestSizer_ObserveAdjustsTowardLatencyTarget(t *testing.T) {
	s, err := New(10, 2, 20, 100)
	if err != nil {
		t.Fatal(err)
	}
	if got := s.Observe(50); got != 12 {
		t.Fatalf("fast observation size = %d, want 12", got)
	}
	if got := s.Observe(150); got != 10 {
		t.Fatalf("slow observation size = %d, want 10", got)
	}
}

func TestNew_ValidatesBounds(t *testing.T) {
	if _, err := New(1, 0, 10, 100); err == nil {
		t.Fatal("step=0 应返回错误")
	}
	if _, err := New(11, 1, 10, 100); err == nil {
		t.Fatal("initial 超出上限应返回错误")
	}
}
