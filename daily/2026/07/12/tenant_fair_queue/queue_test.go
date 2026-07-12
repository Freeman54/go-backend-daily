package tenantfairqueue

import "testing"

func TestQueueAlternatesTenants(t *testing.T) {
	queue := New[int]()
	queue.Push("a", 1)
	queue.Push("a", 2)
	queue.Push("a", 3)
	queue.Push("b", 10)
	queue.Push("b", 11)

	var got []Item[int]
	for i := 0; i < 5; i++ {
		item, ok := queue.Pop()
		if !ok {
			t.Fatalf("pop %d failed", i)
		}
		got = append(got, item)
	}

	wantTenants := []string{"a", "b", "a", "b", "a"}
	for i := range wantTenants {
		if got[i].Tenant != wantTenants[i] {
			t.Fatalf("pop %d tenant = %q want %q", i, got[i].Tenant, wantTenants[i])
		}
	}
}

func TestQueueRemovesEmptyTenant(t *testing.T) {
	queue := New[string]()
	queue.Push("api", "first")
	queue.Push("mq", "second")

	item, ok := queue.Pop()
	if !ok || item.Tenant != "api" {
		t.Fatalf("first pop = %#v, %v", item, ok)
	}

	item, ok = queue.Pop()
	if !ok || item.Tenant != "mq" {
		t.Fatalf("second pop = %#v, %v", item, ok)
	}

	if _, ok := queue.Pop(); ok {
		t.Fatal("queue should be empty")
	}
}

func TestQueueEmpty(t *testing.T) {
	queue := New[int]()
	if _, ok := queue.Pop(); ok {
		t.Fatal("empty queue should return ok=false")
	}
}
