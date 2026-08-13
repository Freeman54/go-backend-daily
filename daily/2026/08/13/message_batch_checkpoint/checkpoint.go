// Package batchcheckpoint advances a checkpoint only through contiguous acknowledgements.
package batchcheckpoint

import "sync"

type Tracker struct {
	mu   sync.Mutex
	next int64
	done map[int64]struct{}
}

func New(next int64) *Tracker { return &Tracker{next: next, done: make(map[int64]struct{})} }

// Ack records completed work and returns the exclusive checkpoint safe to persist.
func (t *Tracker) Ack(offset int64) int64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	if offset >= t.next {
		t.done[offset] = struct{}{}
	}
	for {
		if _, ok := t.done[t.next]; !ok {
			break
		}
		delete(t.done, t.next)
		t.next++
	}
	return t.next
}
func (t *Tracker) Checkpoint() int64 { t.mu.Lock(); defer t.mu.Unlock(); return t.next }
