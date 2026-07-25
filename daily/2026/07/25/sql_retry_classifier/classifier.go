package sql_retry_classifier

import "errors"

type sqlStateError interface {
	error
	SQLState() string
}

func IsRetryable(err error) bool {
	var stateErr sqlStateError
	if !errors.As(err, &stateErr) {
		return false
	}
	switch stateErr.SQLState() {
	case "40001", "40P01", "55P03":
		return true
	default:
		return false
	}
}
