// Package latency_adaptive_sampler preserves important traces while sampling routine traffic.
package latency_adaptive_sampler

import "time"

type Sampler struct {
	slowThreshold time.Duration
	fastRate      uint64
}

func New(slowThreshold time.Duration, fastRate uint64) Sampler {
	if fastRate == 0 {
		fastRate = 1
	}
	return Sampler{slowThreshold: slowThreshold, fastRate: fastRate}
}

func (s Sampler) Keep(latency time.Duration, failed bool, stableHash uint64) bool {
	return failed || latency >= s.slowThreshold || stableHash%s.fastRate == 0
}
