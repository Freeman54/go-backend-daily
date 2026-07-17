package metricwindowdedup

import "time"

type Deduper struct {
	window time.Duration
	last   map[string]time.Time
}

func New(window time.Duration) *Deduper {
	if window < 0 {
		window = 0
	}
	return &Deduper{
		window: window,
		last:   make(map[string]time.Time),
	}
}

func (d *Deduper) ShouldEmit(key string, now time.Time) bool {
	last, ok := d.last[key]
	if ok && now.Sub(last) < d.window {
		return false
	}
	d.last[key] = now
	return true
}
