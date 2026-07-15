package spaneventlimiter

import "fmt"

type Limiter struct {
	limit   int
	counts  map[string]int
	dropped map[string]int
}

func New(limit int) *Limiter {
	return &Limiter{
		limit:   limit,
		counts:  make(map[string]int),
		dropped: make(map[string]int),
	}
}

func (l *Limiter) Allow(name string) bool {
	l.counts[name]++
	if l.counts[name] <= l.limit {
		return true
	}
	l.dropped[name]++
	return false
}

func (l *Limiter) FlushSummaries() []string {
	summaries := make([]string, 0, len(l.dropped))
	for name, dropped := range l.dropped {
		summaries = append(summaries, fmt.Sprintf("%s dropped %d duplicate events", name, dropped))
	}
	l.dropped = make(map[string]int)
	l.counts = make(map[string]int)
	return summaries
}
