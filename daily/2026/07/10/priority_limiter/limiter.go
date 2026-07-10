package prioritylimiter

type Priority string

const (
	High Priority = "high"
	Low  Priority = "low"
)

type Limiter struct {
	total        int
	highReserved int
	highInUse    int
	lowInUse     int
}

func New(total, highReserved int) *Limiter {
	if total <= 0 {
		total = 1
	}
	if highReserved < 0 {
		highReserved = 0
	}
	if highReserved > total {
		highReserved = total
	}
	return &Limiter{
		total:        total,
		highReserved: highReserved,
	}
}

func (l *Limiter) Acquire(priority Priority) bool {
	switch priority {
	case High:
		if l.highInUse+l.lowInUse >= l.total {
			return false
		}
		l.highInUse++
		return true
	case Low:
		lowCap := l.total - l.highReserved
		if l.lowInUse >= lowCap {
			return false
		}
		if l.highInUse+l.lowInUse >= l.total {
			return false
		}
		l.lowInUse++
		return true
	default:
		return false
	}
}

func (l *Limiter) Release(priority Priority) {
	switch priority {
	case High:
		if l.highInUse > 0 {
			l.highInUse--
		}
	case Low:
		if l.lowInUse > 0 {
			l.lowInUse--
		}
	}
}

func (l *Limiter) Snapshot() (highInUse, lowInUse int) {
	return l.highInUse, l.lowInUse
}
