// Package retryenvelope validates retry metadata crossing queue boundaries.
package retryenvelope

import (
	"errors"
	"time"
)

var (
	ErrInvalidAttempt = errors.New("invalid retry attempt")
	ErrClockSkew      = errors.New("retry timestamp is in the future")
)

type Envelope struct {
	Attempt     int
	FirstFailed time.Time
}

type Decision int

const (
	Retry Decision = iota
	DeadLetter
)

// Decide validates untrusted metadata and applies attempt and age budgets.
func Decide(now time.Time, env Envelope, maxAttempts int, maxAge, skew time.Duration) (Decision, error) {
	if maxAttempts <= 0 || env.Attempt < 1 {
		return DeadLetter, ErrInvalidAttempt
	}
	if env.FirstFailed.IsZero() || env.FirstFailed.After(now.Add(skew)) {
		return DeadLetter, ErrClockSkew
	}
	if env.Attempt >= maxAttempts || now.Sub(env.FirstFailed) >= maxAge {
		return DeadLetter, nil
	}
	return Retry, nil
}
