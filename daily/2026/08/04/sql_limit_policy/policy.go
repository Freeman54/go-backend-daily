package sql_limit_policy

import "errors"

var ErrInvalidPolicy = errors.New("default limit must be positive and no greater than maximum")

type Policy struct{ defaultLimit, maxLimit int }

func New(defaultLimit, maxLimit int) (Policy, error) {
	if defaultLimit <= 0 || maxLimit <= 0 || defaultLimit > maxLimit {
		return Policy{}, ErrInvalidPolicy
	}
	return Policy{defaultLimit: defaultLimit, maxLimit: maxLimit}, nil
}

func (p Policy) Normalize(requested int) int {
	if requested <= 0 {
		return p.defaultLimit
	}
	if requested > p.maxLimit {
		return p.maxLimit
	}
	return requested
}
