// Package redelivery classifies message processing outcomes.
package redelivery

import (
	"errors"
	"time"
)

type Decision uint8

const (
	Ack Decision = iota
	Retry
	DeadLetter
)

type PermanentError struct{ Err error }

func (e *PermanentError) Error() string { return "permanent message error: " + e.Err.Error() }
func (e *PermanentError) Unwrap() error { return e.Err }

type Classifier struct {
	MaxAttempts int
	MaxAge      time.Duration
}

// Decide classifies an outcome. attempt is one-based and age is time since publication.
func (c Classifier) Decide(err error, attempt int, age time.Duration) Decision {
	if err == nil {
		return Ack
	}
	var permanent *PermanentError
	if errors.As(err, &permanent) {
		return DeadLetter
	}
	if c.MaxAttempts <= 0 || c.MaxAge <= 0 || attempt <= 0 || attempt >= c.MaxAttempts || age < 0 || age >= c.MaxAge {
		return DeadLetter
	}
	return Retry
}
