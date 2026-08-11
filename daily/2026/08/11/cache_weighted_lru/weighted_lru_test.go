package weightedlru

import (
	"errors"
	"testing"
)

func TestCacheEvictsLeastRecentlyUsedByWeight(t *testing.T) {
	cache, err := New[string, string](5)
	if err != nil {
		t.Fatal(err)
	}
	_ = cache.Put("a", "A", 2)
	_ = cache.Put("b", "B", 2)
	if _, ok := cache.Get("a"); !ok {
		t.Fatal("a should exist")
	}
	_ = cache.Put("c", "C", 2)
	if _, ok := cache.Get("b"); ok {
		t.Fatal("b should be evicted")
	}
	if cache.Weight() != 4 {
		t.Fatalf("weight = %d", cache.Weight())
	}
}

func TestCacheUpdatesAndValidatesWeight(t *testing.T) {
	if _, err := New[string, int](0); err == nil {
		t.Fatal("New() expected error")
	}
	cache, _ := New[string, int](3)
	_ = cache.Put("a", 1, 2)
	_ = cache.Put("a", 2, 1)
	if got, _ := cache.Get("a"); got != 2 || cache.Weight() != 1 {
		t.Fatalf("value/weight = %d/%d", got, cache.Weight())
	}
	if err := cache.Put("x", 1, 0); err == nil {
		t.Fatal("zero weight expected error")
	}
	if err := cache.Put("x", 1, 4); !errors.Is(err, ErrTooHeavy) {
		t.Fatalf("error = %v", err)
	}
}
