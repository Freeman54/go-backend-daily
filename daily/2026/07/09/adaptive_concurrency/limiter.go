package adaptiveconcurrency

import "time"

type Sample struct {
	Latency time.Duration
	Failed  bool
}

type Limiter struct {
	min              int
	max              int
	current          int
	targetLatency    time.Duration
	successThreshold int
	successStreak    int
}

func New(min, max int, targetLatency time.Duration, successThreshold int) *Limiter {
	if min <= 0 {
		min = 1
	}
	if max < min {
		max = min
	}
	if targetLatency <= 0 {
		targetLatency = 50 * time.Millisecond
	}
	if successThreshold <= 0 {
		successThreshold = 3
	}

	return &Limiter{
		min:              min,
		max:              max,
		current:          min,
		targetLatency:    targetLatency,
		successThreshold: successThreshold,
	}
}

func (l *Limiter) Limit() int {
	return l.current
}

func (l *Limiter) Observe(sample Sample) int {
	if sample.Failed || sample.Latency > l.targetLatency {
		l.successStreak = 0
		l.current = maxInt(l.min, l.current/2)
		return l.current
	}

	l.successStreak++
	if sample.Latency <= l.targetLatency/2 && l.successStreak >= l.successThreshold && l.current < l.max {
		l.current++
		l.successStreak = 0
	}

	return l.current
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
