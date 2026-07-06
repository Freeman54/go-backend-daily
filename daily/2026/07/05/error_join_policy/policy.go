package errorjoinpolicy

import "errors"

type Decision string

const (
	DecisionRetry  Decision = "retry"
	DecisionFail   Decision = "fail"
	DecisionIgnore Decision = "ignore"
)

type retryableError struct{ err error }

func (e retryableError) Error() string { return e.err.Error() }
func (e retryableError) Unwrap() error { return e.err }
func (retryableError) Retryable() bool { return true }

type fatalError struct{ err error }

func (e fatalError) Error() string { return e.err.Error() }
func (e fatalError) Unwrap() error { return e.err }
func (fatalError) Fatal() bool     { return true }

type ignoredError struct{ err error }

func (e ignoredError) Error() string { return e.err.Error() }
func (e ignoredError) Unwrap() error { return e.err }
func (ignoredError) Ignored() bool   { return true }

func Retryable(err error) error {
	if err == nil {
		return nil
	}
	return retryableError{err: err}
}

func Fatal(err error) error {
	if err == nil {
		return nil
	}
	return fatalError{err: err}
}

func Ignored(err error) error {
	if err == nil {
		return nil
	}
	return ignoredError{err: err}
}

func Decide(err error) Decision {
	if err == nil {
		return DecisionIgnore
	}

	queue := []error{err}
	seenRetryable := false

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == nil {
			continue
		}

		type unwrapper interface{ Unwrap() []error }
		if joined, ok := current.(unwrapper); ok {
			queue = append(queue, joined.Unwrap()...)
			continue
		}

		type fatalMarker interface{ Fatal() bool }
		if marked, ok := current.(fatalMarker); ok && marked.Fatal() {
			return DecisionFail
		}

		type retryableMarker interface{ Retryable() bool }
		if marked, ok := current.(retryableMarker); ok && marked.Retryable() {
			seenRetryable = true
			continue
		}

		type ignoredMarker interface{ Ignored() bool }
		if marked, ok := current.(ignoredMarker); ok && marked.Ignored() {
			continue
		}

		var target fatalMarker
		if errors.As(current, &target) && target.Fatal() {
			return DecisionFail
		}

		var retryTarget retryableMarker
		if errors.As(current, &retryTarget) && retryTarget.Retryable() {
			seenRetryable = true
			continue
		}
	}

	if seenRetryable {
		return DecisionRetry
	}
	return DecisionIgnore
}
