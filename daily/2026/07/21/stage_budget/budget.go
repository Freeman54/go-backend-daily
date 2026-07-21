package stagebudget

import (
	"context"
	"fmt"
	"time"
)

type Planner struct {
	stages  []stage
	reserve time.Duration
}

type stage struct {
	name   string
	weight int
}

func NewPlanner(reserve time.Duration) *Planner {
	return &Planner{reserve: reserve}
}

func (p *Planner) Add(name string, weight int) error {
	if name == "" {
		return fmt.Errorf("stage name must not be empty")
	}
	if weight <= 0 {
		return fmt.Errorf("stage weight must be positive")
	}
	p.stages = append(p.stages, stage{name: name, weight: weight})
	return nil
}

func (p *Planner) Allocate(ctx context.Context, now time.Time) (map[string]time.Duration, error) {
	if len(p.stages) == 0 {
		return nil, fmt.Errorf("no stages configured")
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return nil, fmt.Errorf("context has no deadline")
	}

	remaining := deadline.Sub(now) - p.reserve
	if remaining <= 0 {
		return nil, fmt.Errorf("deadline budget exhausted")
	}

	totalWeight := 0
	for _, s := range p.stages {
		totalWeight += s.weight
	}

	allocations := make(map[string]time.Duration, len(p.stages))
	var assigned time.Duration
	for i, s := range p.stages {
		if i == len(p.stages)-1 {
			allocations[s.name] = remaining - assigned
			break
		}

		slice := remaining * time.Duration(s.weight) / time.Duration(totalWeight)
		allocations[s.name] = slice
		assigned += slice
	}
	return allocations, nil
}
