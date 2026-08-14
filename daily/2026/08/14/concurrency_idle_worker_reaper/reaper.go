package workerreaper

import (
	"fmt"
	"sort"
	"time"
)

// Worker 描述一个可回收 Worker 的最近活动时间。
type Worker struct {
	ID       string
	LastUsed time.Time
}

// Select 返回应回收的 Worker ID，同时至少保留 minWorkers 个。
// 最空闲的 Worker 优先被回收。
func Select(now time.Time, workers []Worker, idleTimeout time.Duration, minWorkers int) ([]string, error) {
	if idleTimeout <= 0 || minWorkers < 0 {
		return nil, fmt.Errorf("回收参数无效")
	}
	if minWorkers >= len(workers) {
		return nil, nil
	}

	candidates := append([]Worker(nil), workers...)
	sort.SliceStable(candidates, func(i, j int) bool {
		return candidates[i].LastUsed.Before(candidates[j].LastUsed)
	})
	limit := len(workers) - minWorkers
	result := make([]string, 0, limit)
	for _, worker := range candidates {
		if len(result) == limit || now.Sub(worker.LastUsed) < idleTimeout {
			break
		}
		result = append(result, worker.ID)
	}
	return result, nil
}
