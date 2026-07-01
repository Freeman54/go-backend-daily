package latencywindow

import (
	"sort"
	"time"
)

type sample struct {
	at      time.Time
	latency time.Duration
}

type Window struct {
	span      time.Duration
	threshold time.Duration
	samples   []sample
}

func NewWindow(span, threshold time.Duration) *Window {
	return &Window{
		span:      span,
		threshold: threshold,
	}
}

func (w *Window) Record(now time.Time, latency time.Duration) {
	w.samples = append(w.samples, sample{at: now, latency: latency})
	w.evict(now)
}

func (w *Window) Count(now time.Time) int {
	w.evict(now)
	return len(w.samples)
}

func (w *Window) SuccessRatio(now time.Time) float64 {
	w.evict(now)
	if len(w.samples) == 0 {
		return 0
	}

	good := 0
	for _, item := range w.samples {
		if item.latency <= w.threshold {
			good++
		}
	}
	return float64(good) / float64(len(w.samples))
}

func (w *Window) P95(now time.Time) time.Duration {
	w.evict(now)
	if len(w.samples) == 0 {
		return 0
	}

	values := make([]time.Duration, 0, len(w.samples))
	for _, item := range w.samples {
		values = append(values, item.latency)
	}
	sort.Slice(values, func(i, j int) bool {
		return values[i] < values[j]
	})

	index := (len(values)*95 - 1) / 100
	if index < 0 {
		index = 0
	}
	return values[index]
}

func (w *Window) evict(now time.Time) {
	cutoff := now.Add(-w.span)
	drop := 0
	for drop < len(w.samples) && w.samples[drop].at.Before(cutoff) {
		drop++
	}
	if drop > 0 {
		w.samples = append([]sample(nil), w.samples[drop:]...)
	}
}
