package retryclassifier

import (
	"context"
	"errors"
)

var ErrPermanent = errors.New("permanent error")

type temporary interface {
	Temporary() bool
}

type StatusError struct {
	Code int
	Err  error
}

func (e StatusError) Error() string {
	if e.Err == nil {
		return "status error"
	}
	return e.Err.Error()
}

func (e StatusError) Unwrap() error {
	return e.Err
}

type Policy struct {
	MaxAttempts int
}

func (p Policy) ShouldRetry(attempt int, err error) bool {
	if err == nil {
		return false
	}
	if p.MaxAttempts > 0 && attempt >= p.MaxAttempts {
		return false
	}
	if errors.Is(err, ErrPermanent) || errors.Is(err, context.Canceled) {
		return false
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	var statusErr StatusError
	if errors.As(err, &statusErr) && (statusErr.Code == 429 || statusErr.Code >= 500) {
		return true
	}

	var temp temporary
	return errors.As(err, &temp) && temp.Temporary()
}
