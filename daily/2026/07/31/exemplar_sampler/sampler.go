package exemplarsampler

import "hash/fnv"

// Sampler 按 trace ID 稳定采样指定百分比的指标 exemplar。
type Sampler struct {
	rate uint32
}

func New(rate int) Sampler {
	if rate < 0 {
		rate = 0
	}
	if rate > 100 {
		rate = 100
	}
	return Sampler{rate: uint32(rate)}
}

func (s Sampler) Rate() int {
	return int(s.rate)
}

func (s Sampler) Sample(traceID string) bool {
	if traceID == "" || s.rate == 0 {
		return false
	}
	if s.rate == 100 {
		return true
	}
	h := fnv.New32a()
	_, _ = h.Write([]byte(traceID))
	return h.Sum32()%100 < s.rate
}
