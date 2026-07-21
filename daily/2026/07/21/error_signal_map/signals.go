package errorsignalmap

import (
	"errors"
	"fmt"
)

var (
	ErrTimeout     = errors.New("dependency timeout")
	ErrConflict    = errors.New("write conflict")
	ErrInvalid     = errors.New("invalid request")
	ErrUnavailable = errors.New("dependency unavailable")
)

type Signal struct {
	Code      string
	Level     string
	Retryable bool
	Metric    string
}

func Map(err error) Signal {
	switch {
	case err == nil:
		return Signal{Code: "ok", Level: "info", Retryable: false, Metric: "request_ok_total"}
	case errors.Is(err, ErrInvalid):
		return Signal{Code: "invalid_argument", Level: "warn", Retryable: false, Metric: "request_invalid_total"}
	case errors.Is(err, ErrConflict):
		return Signal{Code: "conflict", Level: "warn", Retryable: true, Metric: "request_conflict_total"}
	case errors.Is(err, ErrTimeout):
		return Signal{Code: "timeout", Level: "error", Retryable: true, Metric: "dependency_timeout_total"}
	case errors.Is(err, ErrUnavailable):
		return Signal{Code: "unavailable", Level: "error", Retryable: true, Metric: "dependency_unavailable_total"}
	default:
		return Signal{Code: "internal", Level: "error", Retryable: false, Metric: "request_internal_total"}
	}
}

func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf(format+": %w", append(args, err)...)
}
