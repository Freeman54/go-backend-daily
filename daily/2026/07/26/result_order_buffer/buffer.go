package result_order_buffer

import "sync"

// Buffer 暂存乱序结果，并在序号连续时按顺序释放。
type Buffer[T any] struct {
	mu      sync.Mutex
	next    uint64
	pending map[uint64]T
}

func New[T any](next uint64) *Buffer[T] {
	return &Buffer[T]{next: next, pending: make(map[uint64]T)}
}

func (b *Buffer[T]) Add(sequence uint64, value T) []T {
	b.mu.Lock()
	defer b.mu.Unlock()
	if sequence < b.next {
		return nil
	}
	if _, exists := b.pending[sequence]; exists {
		return nil
	}
	b.pending[sequence] = value
	var ready []T
	for {
		value, exists := b.pending[b.next]
		if !exists {
			return ready
		}
		ready = append(ready, value)
		delete(b.pending, b.next)
		b.next++
	}
}
