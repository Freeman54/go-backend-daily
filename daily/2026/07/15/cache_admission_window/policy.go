package cacheadmissionwindow

import "time"

type Policy struct {
	threshold int
	window    time.Duration
	now       func() time.Time
	entries   map[string]entry
}

type entry struct {
	count int
	start time.Time
}

func New(threshold int, window time.Duration, now func() time.Time) *Policy {
	if now == nil {
		now = time.Now
	}
	return &Policy{
		threshold: threshold,
		window:    window,
		now:       now,
		entries:   make(map[string]entry),
	}
}

func (p *Policy) Record(key string) bool {
	current := p.now()
	state, ok := p.entries[key]
	if !ok || current.Sub(state.start) >= p.window {
		state = entry{start: current}
	}

	state.count++
	p.entries[key] = state

	return state.count >= p.threshold
}

func (p *Policy) SeenCount(key string) int {
	return p.entries[key].count
}
