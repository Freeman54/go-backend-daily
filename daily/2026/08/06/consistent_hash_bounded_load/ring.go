// Package consistent_hash_bounded_load combines consistent hashing with load bounds.
package consistent_hash_bounded_load

import (
	"errors"
	"hash/fnv"
	"math"
	"sort"
	"strconv"
)

type point struct {
	hash uint64
	node string
}

type Ring struct {
	points []point
	nodes  []string
	factor float64
}

func New(nodes []string, replicas int, maxLoadFactor float64) (*Ring, error) {
	if len(nodes) == 0 || replicas <= 0 || maxLoadFactor < 1 {
		return nil, errors.New("nodes, positive replicas, and factor >= 1 are required")
	}
	r := &Ring{nodes: append([]string(nil), nodes...), factor: maxLoadFactor}
	for _, node := range nodes {
		if node == "" {
			return nil, errors.New("node name cannot be empty")
		}
		for replica := 0; replica < replicas; replica++ {
			r.points = append(r.points, point{hash: hash(node + "#" + strconv.Itoa(replica)), node: node})
		}
	}
	sort.Slice(r.points, func(i, j int) bool { return r.points[i].hash < r.points[j].hash })
	return r, nil
}

// Pick returns the first clockwise node within the bounded-load threshold.
func (r *Ring) Pick(key string, loads map[string]int) (string, error) {
	if len(r.points) == 0 {
		return "", errors.New("empty ring")
	}
	start := sort.Search(len(r.points), func(i int) bool { return r.points[i].hash >= hash(key) }) % len(r.points)
	if loads == nil {
		return r.points[start].node, nil
	}
	total := 0
	for _, node := range r.nodes {
		total += loads[node]
	}
	limit := int(math.Ceil(float64(total) / float64(len(r.nodes)) * r.factor))
	seen := make(map[string]struct{}, len(r.nodes))
	for offset := 0; offset < len(r.points); offset++ {
		node := r.points[(start+offset)%len(r.points)].node
		if _, ok := seen[node]; ok {
			continue
		}
		seen[node] = struct{}{}
		if loads[node] <= limit {
			return node, nil
		}
	}
	return "", errors.New("all nodes exceed the load bound")
}

func hash(value string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(value))
	return h.Sum64()
}
