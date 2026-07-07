package slowquerysampler

import (
	"errors"
	"hash/fnv"
	"time"
)

type Sampler struct {
	slowThreshold time.Duration
	sampleRate    uint32
}

func New(slowThreshold time.Duration, sampleRate uint32) (*Sampler, error) {
	if slowThreshold <= 0 {
		return nil, errors.New("slowThreshold must be positive")
	}
	if sampleRate == 0 || sampleRate > 100 {
		return nil, errors.New("sampleRate must be between 1 and 100")
	}
	return &Sampler{slowThreshold: slowThreshold, sampleRate: sampleRate}, nil
}

func (s *Sampler) ShouldLog(operation string, latency time.Duration, failed bool) bool {
	if failed {
		return true
	}
	if latency >= s.slowThreshold {
		return true
	}

	hash := fnv.New32a()
	_, _ = hash.Write([]byte(operation))
	score := hash.Sum32()%100 + 1
	return score <= s.sampleRate
}
