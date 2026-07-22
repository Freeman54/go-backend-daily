package cacherefreshjitter

import (
	"hash/fnv"
	"time"
)

func RefreshAfter(key string, base time.Duration, percent int) time.Duration {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	if base <= 0 || percent == 0 {
		return base
	}
	h := fnv.New64a()
	_, _ = h.Write([]byte(key))
	unit := float64(h.Sum64()) / float64(^uint64(0))
	factor := 1 - float64(percent)/100 + 2*float64(percent)/100*unit
	return time.Duration(float64(base) * factor)
}
