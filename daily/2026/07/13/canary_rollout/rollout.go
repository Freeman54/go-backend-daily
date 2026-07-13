package canaryrollout

import "hash/fnv"

type Rule struct {
	Key       string
	Percent   uint8
	Allowlist map[string]struct{}
	Denylist  map[string]struct{}
}

func Enabled(rule Rule, actorID string) bool {
	if _, ok := rule.Denylist[actorID]; ok {
		return false
	}
	if _, ok := rule.Allowlist[actorID]; ok {
		return true
	}
	if rule.Percent == 0 {
		return false
	}
	if rule.Percent >= 100 {
		return true
	}
	return bucket(rule.Key, actorID) < uint32(rule.Percent)
}

func bucket(key string, actorID string) uint32 {
	hasher := fnv.New32a()
	_, _ = hasher.Write([]byte(key))
	_, _ = hasher.Write([]byte{':'})
	_, _ = hasher.Write([]byte(actorID))
	return hasher.Sum32() % 100
}
