package tenantfairqueue

import "sync"

type Item[T any] struct {
	Tenant string
	Value  T
}

type Queue[T any] struct {
	mu      sync.Mutex
	order   []string
	buckets map[string][]T
	next    int
}

func New[T any]() *Queue[T] {
	return &Queue[T]{
		buckets: make(map[string][]T),
	}
}

func (q *Queue[T]) Push(tenant string, value T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.buckets[tenant]) == 0 {
		q.order = append(q.order, tenant)
	}

	q.buckets[tenant] = append(q.buckets[tenant], value)
}

func (q *Queue[T]) Pop() (Item[T], bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for len(q.order) > 0 {
		if q.next >= len(q.order) {
			q.next = 0
		}

		idx := q.next
		tenant := q.order[idx]
		bucket := q.buckets[tenant]
		if len(bucket) == 0 {
			q.removeTenant(idx, tenant)
			continue
		}

		value := bucket[0]
		if len(bucket) == 1 {
			delete(q.buckets, tenant)
			q.removeTenant(idx, tenant)
		} else {
			q.buckets[tenant] = bucket[1:]
			q.next = (idx + 1) % len(q.order)
		}

		return Item[T]{Tenant: tenant, Value: value}, true
	}

	var zero Item[T]
	return zero, false
}

func (q *Queue[T]) removeTenant(idx int, tenant string) {
	delete(q.buckets, tenant)
	q.order = append(q.order[:idx], q.order[idx+1:]...)
	if len(q.order) == 0 {
		q.next = 0
		return
	}
	if idx < q.next || q.next >= len(q.order) {
		q.next = idx % len(q.order)
	}
}
