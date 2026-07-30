package cachettlwheel

import (
	"testing"
	"time"
)

func TestWheel_ExpiresEntriesByTick(t *testing.T) {
	now := time.Unix(100, 0)
	w := New(time.Second, 4, now)
	w.Set("short", "a", 1500*time.Millisecond)
	w.Set("long", "b", 5*time.Second)

	now = now.Add(time.Second)
	if expired := w.Advance(now); len(expired) != 0 {
		t.Fatalf("过早淘汰: %v", expired)
	}
	now = now.Add(time.Second)
	if expired := w.Advance(now); len(expired) != 1 || expired[0] != "short" {
		t.Fatalf("expired = %v, want [short]", expired)
	}
	if value, ok := w.Get("long", now); !ok || value != "b" {
		t.Fatalf("长 TTL 条目丢失: %q, %v", value, ok)
	}
	now = now.Add(3 * time.Second)
	if expired := w.Advance(now); len(expired) != 1 || expired[0] != "long" {
		t.Fatalf("expired = %v, want [long]", expired)
	}
}

func TestWheel_OverwriteIgnoresOldSchedule(t *testing.T) {
	now := time.Unix(200, 0)
	w := New(time.Second, 4, now)
	w.Set("k", "old", time.Second)
	w.Set("k", "new", 3*time.Second)
	if expired := w.Advance(now.Add(time.Second)); len(expired) != 0 {
		t.Fatalf("旧调度不应删除新值: %v", expired)
	}
	if got, ok := w.Get("k", now.Add(time.Second)); !ok || got != "new" {
		t.Fatalf("Get = %q, %v", got, ok)
	}
}
