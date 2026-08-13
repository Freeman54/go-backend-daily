// Package freshnesspolicy classifies cache entries by age.
package freshnesspolicy

import "time"

type State uint8

const (
	Fresh State = iota
	Stale
	Expired
)

type Policy struct{ FreshFor, ServeStaleFor time.Duration }

func (p Policy) Classify(storedAt, now time.Time) State {
	age := now.Sub(storedAt)
	if age < 0 {
		age = 0
	}
	if age <= max(p.FreshFor, 0) {
		return Fresh
	}
	if age <= max(p.FreshFor, 0)+max(p.ServeStaleFor, 0) {
		return Stale
	}
	return Expired
}

func max(v time.Duration, floor time.Duration) time.Duration {
	if v < floor {
		return floor
	}
	return v
}
