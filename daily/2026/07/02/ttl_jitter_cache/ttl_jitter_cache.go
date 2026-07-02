package ttljittercache

import (
	"math/rand"
	"time"
)

func JitterTTL(base time.Duration, maxJitter time.Duration, rnd *rand.Rand) time.Duration {
	if rnd == nil || maxJitter <= 0 {
		return base
	}
	return base + time.Duration(rnd.Int63n(int64(maxJitter)+1))
}

func StaggerExpirations(now time.Time, keys []string, baseTTL time.Duration, maxJitter time.Duration, rnd *rand.Rand) map[string]time.Time {
	expiresAt := make(map[string]time.Time, len(keys))
	for _, key := range keys {
		expiresAt[key] = now.Add(JitterTTL(baseTTL, maxJitter, rnd))
	}
	return expiresAt
}
