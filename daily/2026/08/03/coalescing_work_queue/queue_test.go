package coalescing_work_queue

import "testing"

func TestQueue_CoalescesPendingKeys(t *testing.T) {
	q := New(2)
	if !q.Push("user:1") || q.Push("user:1") || !q.Push("user:2") {
		t.Fatal("pending keys should be unique")
	}
	if q.Push("user:3") {
		t.Fatal("full queue should reject new key")
	}
	key, ok := q.Pop()
	if !ok || key != "user:1" {
		t.Fatalf("Pop() = %q, %v", key, ok)
	}
	if !q.Push("user:1") {
		t.Fatal("popped key should be enqueueable again")
	}
}
