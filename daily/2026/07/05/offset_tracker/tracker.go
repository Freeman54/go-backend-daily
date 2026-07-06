package offsettracker

type Tracker struct {
	next    int64
	pending map[int64]struct{}
}

func New(startOffset int64) *Tracker {
	return &Tracker{
		next:    startOffset,
		pending: make(map[int64]struct{}),
	}
}

func (t *Tracker) Ack(offset int64) (commitTo int64, advanced bool) {
	if offset < t.next {
		return t.next - 1, false
	}

	previousNext := t.next
	t.pending[offset] = struct{}{}
	commit := t.next - 1
	for {
		if _, ok := t.pending[t.next]; !ok {
			break
		}
		delete(t.pending, t.next)
		commit = t.next
		t.next++
	}

	return commit, t.next > previousNext
}

func (t *Tracker) Next() int64 {
	return t.next
}
