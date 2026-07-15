package cacheadmissionwindow

import (
	"testing"
	"time"
)

func TestRecordAdmitsAfterThresholdWithinWindow(t *testing.T) {
	current := time.Unix(0, 0)
	policy := New(3, time.Minute, func() time.Time { return current })

	if policy.Record("sku-1") {
		t.Fatal("first hit should not be admitted")
	}
	if policy.Record("sku-1") {
		t.Fatal("second hit should not be admitted")
	}
	if !policy.Record("sku-1") {
		t.Fatal("third hit should be admitted")
	}
}

func TestRecordResetsAfterWindowExpires(t *testing.T) {
	current := time.Unix(0, 0)
	policy := New(2, time.Minute, func() time.Time { return current })

	policy.Record("sku-1")
	current = current.Add(2 * time.Minute)

	if policy.Record("sku-1") {
		t.Fatal("count should reset after window expires")
	}
}

func TestRecordKeepsKeysIsolated(t *testing.T) {
	current := time.Unix(0, 0)
	policy := New(2, time.Minute, func() time.Time { return current })

	policy.Record("sku-1")
	if policy.Record("sku-2") {
		t.Fatal("different key should have an independent counter")
	}
}
