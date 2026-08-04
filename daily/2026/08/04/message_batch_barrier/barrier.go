package message_batch_barrier

import "sync"

// Barrier tracks out-of-order acknowledgements and exposes the highest contiguous offset.
type Barrier struct {
	mu        sync.Mutex
	next      int64
	committed int64
	acked     map[int64]struct{}
}

func New(firstOffset int64) *Barrier {
	return &Barrier{next: firstOffset, committed: firstOffset - 1, acked: make(map[int64]struct{})}
}

func (b *Barrier) Ack(offset int64) int64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	if offset < b.next {
		return b.committed
	}
	b.acked[offset] = struct{}{}
	for {
		if _, ok := b.acked[b.next]; !ok {
			break
		}
		delete(b.acked, b.next)
		b.committed = b.next
		b.next++
	}
	return b.committed
}

func (b *Barrier) Pending() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.acked)
}
