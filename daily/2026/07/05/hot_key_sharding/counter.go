package hotkeysharding

import (
	"hash/fnv"
	"sync"
)

type shard struct {
	mu sync.Mutex
	m  map[string]int64
}

type Counter struct {
	shards []shard
}

func New(shardCount int) *Counter {
	if shardCount <= 0 {
		shardCount = 1
	}
	shards := make([]shard, shardCount)
	for i := range shards {
		shards[i].m = make(map[string]int64)
	}
	return &Counter{shards: shards}
}

func (c *Counter) Add(key string, delta int64) {
	s := &c.shards[index(key, len(c.shards))]
	s.mu.Lock()
	s.m[key] += delta
	s.mu.Unlock()
}

func (c *Counter) Get(key string) int64 {
	s := &c.shards[index(key, len(c.shards))]
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.m[key]
}

func (c *Counter) ShardLoads() []int {
	loads := make([]int, len(c.shards))
	for i := range c.shards {
		c.shards[i].mu.Lock()
		loads[i] = len(c.shards[i].m)
		c.shards[i].mu.Unlock()
	}
	return loads
}

func index(key string, size int) int {
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(key))
	return int(hash.Sum32() % uint32(size))
}
