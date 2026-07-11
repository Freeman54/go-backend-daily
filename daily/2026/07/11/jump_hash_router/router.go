package jumphashrouter

import "hash/fnv"

type Router struct {
	partitions int32
}

func New(partitions int) Router {
	return Router{partitions: int32(partitions)}
}

func (r Router) Route(key string) int {
	if r.partitions <= 0 {
		return -1
	}
	return int(jumpHash(hashKey(key), r.partitions))
}

func MovedRatio(keys []string, fromPartitions, toPartitions int) float64 {
	if len(keys) == 0 || fromPartitions <= 0 || toPartitions <= 0 {
		return 0
	}

	before := New(fromPartitions)
	after := New(toPartitions)

	var moved int
	for _, key := range keys {
		if before.Route(key) != after.Route(key) {
			moved++
		}
	}
	return float64(moved) / float64(len(keys))
}

func hashKey(key string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(key))
	return h.Sum64()
}

func jumpHash(key uint64, buckets int32) int32 {
	var b int64 = -1
	var j int64

	for j < int64(buckets) {
		b = j
		key = key*2862933555777941757 + 1
		j = int64(float64(b+1) * (float64(1<<31) / float64((key>>33)+1)))
	}

	return int32(b)
}
