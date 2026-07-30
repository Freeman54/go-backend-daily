package adaptiveconcurrencylimiter

import (
	"testing"
	"time"
)

func TestLimiter_AdjustsLimitFromLatencySamples(t *testing.T) {
	l := New(2, 1, 4, 100*time.Millisecond)
	if !l.TryAcquire() || !l.TryAcquire() || l.TryAcquire() {
		t.Fatal("初始并发上限应为 2")
	}
	l.Release(50*time.Millisecond, true)
	l.Release(50*time.Millisecond, true)
	if got := l.Limit(); got != 3 {
		t.Fatalf("快速成功后 limit = %d, want 3", got)
	}
	if !l.TryAcquire() || !l.TryAcquire() || !l.TryAcquire() || l.TryAcquire() {
		t.Fatal("并发上限提升后应允许 3 个请求")
	}
	l.Release(200*time.Millisecond, true)
	l.Release(20*time.Millisecond, false)
	l.Release(20*time.Millisecond, false)
	if got := l.Limit(); got != 1 {
		t.Fatalf("慢请求和失败后 limit = %d, want 1", got)
	}
}

func TestNew_RejectsInvalidConfiguration(t *testing.T) {
	tests := []struct {
		initial, min, max int
		target            time.Duration
	}{
		{0, 1, 2, time.Second},
		{3, 1, 2, time.Second},
		{1, 2, 3, time.Second},
		{1, 1, 2, 0},
	}
	for _, tt := range tests {
		if _, err := NewChecked(tt.initial, tt.min, tt.max, tt.target); err == nil {
			t.Fatalf("NewChecked(%d, %d, %d, %s) 应失败", tt.initial, tt.min, tt.max, tt.target)
		}
	}
}
