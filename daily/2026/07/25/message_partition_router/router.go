package message_partition_router

import (
	"errors"
	"hash/fnv"
)

type Router struct {
	partitions uint32
}

func New(partitions int) (*Router, error) {
	if partitions <= 0 {
		return nil, errors.New("partitions must be positive")
	}
	return &Router{partitions: uint32(partitions)}, nil
}

func (r *Router) Partition(key string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return int(h.Sum32() % r.partitions)
}
