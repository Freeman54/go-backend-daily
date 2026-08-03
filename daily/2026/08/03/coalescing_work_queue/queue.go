// Package coalescing_work_queue implements a bounded FIFO that keeps one pending item per key.
package coalescing_work_queue

import "sync"

type Queue struct {
	mu      sync.Mutex
	limit   int
	items   []string
	pending map[string]struct{}
}

func New(limit int) *Queue {
	if limit < 0 {
		limit = 0
	}
	return &Queue{limit: limit, pending: make(map[string]struct{})}
}

func (q *Queue) Push(key string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if _, exists := q.pending[key]; exists || len(q.items) >= q.limit {
		return false
	}
	q.items = append(q.items, key)
	q.pending[key] = struct{}{}
	return true
}

func (q *Queue) Pop() (string, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return "", false
	}
	key := q.items[0]
	q.items = q.items[1:]
	delete(q.pending, key)
	return key, true
}
