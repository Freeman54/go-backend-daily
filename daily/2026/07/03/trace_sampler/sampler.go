package tracesampler

import (
	"hash/fnv"
	"time"
)

type Meta struct {
	TraceID  string
	Duration time.Duration
	HasError bool
	Forced   bool
}

type Decision struct {
	Sample bool
	Reason string
}

type Sampler struct {
	basePermille  int
	slowThreshold time.Duration
}

func New(basePermille int, slowThreshold time.Duration) Sampler {
	if basePermille < 0 {
		basePermille = 0
	}
	if basePermille > 1000 {
		basePermille = 1000
	}
	return Sampler{
		basePermille:  basePermille,
		slowThreshold: slowThreshold,
	}
}

func (s Sampler) Decide(meta Meta) Decision {
	switch {
	case meta.Forced:
		return Decision{Sample: true, Reason: "forced"}
	case meta.HasError:
		return Decision{Sample: true, Reason: "error"}
	case meta.Duration >= s.slowThreshold:
		return Decision{Sample: true, Reason: "slow"}
	case s.basePermille == 0:
		return Decision{Sample: false, Reason: "base-drop"}
	default:
		if bucket(meta.TraceID) < uint32(s.basePermille) {
			return Decision{Sample: true, Reason: "base-hit"}
		}
		return Decision{Sample: false, Reason: "base-drop"}
	}
}

func bucket(traceID string) uint32 {
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(traceID))
	return hash.Sum32() % 1000
}
