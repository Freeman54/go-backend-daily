package contexthandoff

import (
	"context"
	"errors"
	"time"
)

var ErrNoDeadline = errors.New("context handoff: missing deadline")
var ErrInsufficientBudget = errors.New("context handoff: insufficient budget")

type Stage struct {
	Name   string
	Weight int
}

type Slice struct {
	Name     string
	Duration time.Duration
}

func Plan(ctx context.Context, now time.Time, reserve time.Duration, minEach time.Duration, stages []Stage) ([]Slice, error) {
	deadline, ok := ctx.Deadline()
	if !ok {
		return nil, ErrNoDeadline
	}
	if len(stages) == 0 {
		return nil, nil
	}
	if reserve < 0 {
		reserve = 0
	}
	if minEach < 0 {
		minEach = 0
	}

	available := deadline.Sub(now) - reserve
	if available <= 0 {
		return nil, ErrInsufficientBudget
	}
	required := minEach * time.Duration(len(stages))
	if available < required {
		return nil, ErrInsufficientBudget
	}

	totalWeight := 0
	for _, stage := range stages {
		if stage.Weight <= 0 {
			totalWeight++
			continue
		}
		totalWeight += stage.Weight
	}

	base := available - required
	slices := make([]Slice, 0, len(stages))
	used := time.Duration(0)

	for i, stage := range stages {
		weight := stage.Weight
		if weight <= 0 {
			weight = 1
		}
		duration := minEach
		if base > 0 {
			if i == len(stages)-1 {
				duration += base - used
			} else {
				extra := time.Duration(int64(base) * int64(weight) / int64(totalWeight))
				duration += extra
				used += extra
			}
		}
		slices = append(slices, Slice{Name: stage.Name, Duration: duration})
	}

	return slices, nil
}
