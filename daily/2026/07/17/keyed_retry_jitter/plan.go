package keyedretryjitter

import (
	"hash/fnv"
	"math"
	"time"
)

func Delay(key string, attempt int, base time.Duration, max time.Duration, jitterRatio float64) time.Duration {
	if attempt < 1 {
		attempt = 1
	}
	if base <= 0 {
		base = time.Second
	}
	if max <= 0 {
		max = base
	}
	if jitterRatio < 0 {
		jitterRatio = 0
	}
	if jitterRatio > 1 {
		jitterRatio = 1
	}

	delay := base * time.Duration(1<<(attempt-1))
	if delay > max {
		delay = max
	}

	jitter := time.Duration(float64(delay) * jitterRatio)
	if jitter == 0 {
		return delay
	}

	spread := stableSpread(key, attempt)
	offset := time.Duration(float64(jitter) * spread)
	adjusted := delay - jitter/2 + offset
	if adjusted < 0 {
		return 0
	}
	if adjusted > max {
		return max
	}
	return adjusted
}

func stableSpread(key string, attempt int) float64 {
	hasher := fnv.New32a()
	_, _ = hasher.Write([]byte(key))
	_, _ = hasher.Write([]byte{byte(attempt), byte(attempt >> 8), byte(attempt >> 16), byte(attempt >> 24)})
	return float64(hasher.Sum32()) / float64(math.MaxUint32)
}
