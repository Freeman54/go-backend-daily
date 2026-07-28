package cache_probabilistic_refresh

import "time"

// ShouldRefresh 用随机样本把刷新请求分散到条目过期前。
func ShouldRefresh(now, expiresAt time.Time, ttl time.Duration, sample, earlyFraction float64) bool {
	if !now.Before(expiresAt) {
		return true
	}
	if ttl <= 0 || sample < 0 || sample > 1 || earlyFraction <= 0 || earlyFraction > 1 {
		return false
	}
	remainingRatio := float64(expiresAt.Sub(now)) / float64(ttl)
	return remainingRatio <= earlyFraction && sample >= remainingRatio
}
