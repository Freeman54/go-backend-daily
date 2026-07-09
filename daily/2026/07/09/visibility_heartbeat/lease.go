package visibilityheartbeat

import "time"

type Lease struct {
	start        time.Time
	expiresAt    time.Time
	visibility   time.Duration
	renewBefore  time.Duration
	maxExtension time.Duration
}

func NewLease(start time.Time, visibility, renewBefore, maxExtension time.Duration) *Lease {
	if visibility <= 0 {
		visibility = 30 * time.Second
	}
	if renewBefore <= 0 || renewBefore >= visibility {
		renewBefore = visibility / 3
	}
	if maxExtension < visibility {
		maxExtension = visibility
	}

	return &Lease{
		start:        start,
		expiresAt:    start.Add(visibility),
		visibility:   visibility,
		renewBefore:  renewBefore,
		maxExtension: maxExtension,
	}
}

func (l *Lease) Due(now time.Time) bool {
	return !now.Before(l.expiresAt.Add(-l.renewBefore))
}

func (l *Lease) Heartbeat(now time.Time) (renewed bool, giveUp bool) {
	if !l.Due(now) {
		return false, false
	}

	nextExpiry := l.expiresAt.Add(l.visibility)
	if nextExpiry.Sub(l.start) > l.maxExtension {
		return false, true
	}

	l.expiresAt = nextExpiry
	return true, false
}

func (l *Lease) ExpiresAt() time.Time {
	return l.expiresAt
}
