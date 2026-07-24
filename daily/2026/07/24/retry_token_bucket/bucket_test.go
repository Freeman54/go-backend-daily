package retrytokenbucket

import (
	"testing"
	"time"
)

func TestBucket_TakeConsumesCapacityAndRefills(t *testing.T) {
	now := time.Unix(100, 0)
	bucket := New(2, time.Second, now)
	if !bucket.Take(now) || !bucket.Take(now) {
		t.Fatal("initial tokens should be available")
	}
	if bucket.Take(now) {
		t.Fatal("third token should be rejected")
	}
	if !bucket.Take(now.Add(time.Second)) {
		t.Fatal("one token should refill after one interval")
	}
}

func TestNew_RejectsInvalidConfiguration(t *testing.T) {
	if _, err := NewChecked(0, time.Second, time.Now()); err == nil {
		t.Fatal("zero capacity should return an error")
	}
	if _, err := NewChecked(1, 0, time.Now()); err == nil {
		t.Fatal("zero refill interval should return an error")
	}
}
