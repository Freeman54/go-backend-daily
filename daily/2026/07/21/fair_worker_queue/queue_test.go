package fairworkerqueue

import "testing"

func TestQueueDrainsInTenantRoundRobinOrder(t *testing.T) {
	q := New()
	mustEnqueue(t, q, "tenant-a", "a1")
	mustEnqueue(t, q, "tenant-a", "a2")
	mustEnqueue(t, q, "tenant-b", "b1")
	mustEnqueue(t, q, "tenant-c", "c1")
	mustEnqueue(t, q, "tenant-b", "b2")

	var got []Item
	for {
		item, ok := q.Next()
		if !ok {
			break
		}
		got = append(got, item)
	}

	want := []Item{
		{Tenant: "tenant-a", Value: "a1"},
		{Tenant: "tenant-b", Value: "b1"},
		{Tenant: "tenant-c", Value: "c1"},
		{Tenant: "tenant-a", Value: "a2"},
		{Tenant: "tenant-b", Value: "b2"},
	}

	if len(got) != len(want) {
		t.Fatalf("got %d items, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("item %d = %#v, want %#v", i, got[i], want[i])
		}
	}
}

func TestQueueRejectsEmptyTenant(t *testing.T) {
	q := New()
	if err := q.Enqueue("", "job"); err == nil {
		t.Fatal("expected empty tenant to fail")
	}
}

func TestQueueLenTracksRemainingItems(t *testing.T) {
	q := New()
	mustEnqueue(t, q, "tenant-a", "a1")
	mustEnqueue(t, q, "tenant-b", "b1")

	if q.Len() != 2 {
		t.Fatalf("len = %d, want 2", q.Len())
	}

	if _, ok := q.Next(); !ok {
		t.Fatal("expected item")
	}
	if q.Len() != 1 {
		t.Fatalf("len = %d, want 1", q.Len())
	}
}

func mustEnqueue(t *testing.T, q *Queue, tenant string, value string) {
	t.Helper()
	if err := q.Enqueue(tenant, value); err != nil {
		t.Fatalf("enqueue %s/%s: %v", tenant, value, err)
	}
}
