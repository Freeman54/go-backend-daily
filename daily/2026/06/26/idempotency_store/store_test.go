package idempotencystore

import (
	"testing"
	"time"
)

func TestStoreRememberReturnsCachedResponseBeforeExpiry(t *testing.T) {
	t.Parallel()

	store := NewStore()
	key := "order:create:42"
	want := Response{StatusCode: 201, Body: `{"id":42}`}

	store.Remember(key, want, time.Minute)

	got, ok := store.Lookup(key, time.Unix(30, 0))
	if !ok {
		t.Fatal("Lookup() reported missing response")
	}
	if got != want {
		t.Fatalf("Lookup() = %#v, want %#v", got, want)
	}
}

func TestStoreLookupRemovesExpiredEntry(t *testing.T) {
	t.Parallel()

	store := NewStore()
	key := "order:create:43"
	store.RememberAt(key, Response{StatusCode: 202, Body: "accepted"}, time.Minute, time.Unix(0, 0))

	_, ok := store.Lookup(key, time.Unix(120, 0))
	if ok {
		t.Fatal("Lookup() should evict expired entry")
	}
}
