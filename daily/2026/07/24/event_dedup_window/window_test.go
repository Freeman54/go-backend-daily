package eventdedupwindow

import (
	"testing"
	"time"
)

func TestWindow_FirstSeenThenDuplicateUntilExpiry(t *testing.T) {
	now := time.Unix(100, 0)
	window, err := New(2 * time.Second)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if !window.Accept("evt-1", now) {
		t.Fatal("first event should be accepted")
	}
	if window.Accept("evt-1", now.Add(time.Second)) {
		t.Fatal("duplicate inside window should be rejected")
	}
	if !window.Accept("evt-1", now.Add(2*time.Second)) {
		t.Fatal("event at expiry should be accepted again")
	}
}

func TestWindow_RejectsEmptyIDAndInvalidTTL(t *testing.T) {
	if _, err := New(0); err == nil {
		t.Fatal("zero TTL should return an error")
	}
	window, _ := New(time.Second)
	if window.Accept("", time.Now()) {
		t.Fatal("empty event ID should be rejected")
	}
}
