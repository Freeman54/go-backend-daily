package logsamplingpolicy

import (
	"errors"
	"hash/fnv"
)

type Event struct {
	RequestID      string
	StatusCode     int
	DurationMillis int64
	Force          bool
}

type Policy struct{ threshold uint64 }

func New(rate float64) (Policy, error) {
	if rate < 0 || rate > 1 {
		return Policy{}, errors.New("sample rate must be between zero and one")
	}
	return Policy{threshold: uint64(rate * 1_000_000)}, nil
}

func (p Policy) Keep(e Event) bool {
	if e.Force || e.StatusCode >= 500 || e.DurationMillis >= 1000 {
		return true
	}
	h := fnv.New64a()
	_, _ = h.Write([]byte(e.RequestID))
	return h.Sum64()%1_000_000 < p.threshold
}
