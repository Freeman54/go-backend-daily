package fallbackbudget

import (
	"time"
)

type Budget struct {
	Limit       int
	Window      time.Duration
	used        int
	windowStart time.Time
}

func New(limit int, window time.Duration) *Budget {
	return &Budget{
		Limit:  limit,
		Window: window,
	}
}

func (b *Budget) Allow(now time.Time) bool {
	b.rollWindow(now)
	if b.used >= b.Limit {
		return false
	}
	b.used++
	return true
}

func (b *Budget) Remaining(now time.Time) int {
	b.rollWindow(now)
	if b.used >= b.Limit {
		return 0
	}
	return b.Limit - b.used
}

func (b *Budget) rollWindow(now time.Time) {
	if b.windowStart.IsZero() {
		b.windowStart = now
		return
	}
	if now.Sub(b.windowStart) >= b.Window {
		b.windowStart = now
		b.used = 0
	}
}
