// Package ttlcap constrains cache TTL values to operationally safe bounds.
package ttlcap

import (
	"errors"
	"time"
)

var ErrInvalidPolicy = errors.New("invalid TTL cap policy")

type Policy struct {
	Min         time.Duration
	Max         time.Duration
	NegativeMax time.Duration
}

// Apply clamps a desired TTL. Non-positive desired values use the applicable maximum.
func (p Policy) Apply(desired time.Duration, negative bool) (time.Duration, error) {
	if p.Min <= 0 || p.Max < p.Min || p.NegativeMax < p.Min || p.NegativeMax > p.Max {
		return 0, ErrInvalidPolicy
	}
	upper := p.Max
	if negative {
		upper = p.NegativeMax
	}
	if desired <= 0 || desired > upper {
		return upper, nil
	}
	if desired < p.Min {
		return p.Min, nil
	}
	return desired, nil
}
