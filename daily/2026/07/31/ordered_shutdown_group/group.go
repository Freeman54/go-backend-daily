package orderedshutdowngroup

import (
	"context"
	"errors"
	"sort"
)

type Hook func(context.Context) error

type entry struct {
	phase int
	order int
	name  string
	hook  Hook
}

// Group 按阶段从高到低执行停机钩子，同阶段保持注册顺序。
type Group struct {
	entries []entry
}

func New() *Group {
	return &Group{}
}

func (g *Group) Add(phase int, name string, hook Hook) {
	g.entries = append(g.entries, entry{phase: phase, order: len(g.entries), name: name, hook: hook})
}

func (g *Group) Shutdown(ctx context.Context) error {
	entries := append([]entry(nil), g.entries...)
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].phase > entries[j].phase
	})

	var errs []error
	for _, item := range entries {
		if err := ctx.Err(); err != nil {
			errs = append(errs, err)
			break
		}
		if item.hook == nil {
			continue
		}
		if err := item.hook(ctx); err != nil {
			errs = append(errs, errors.New(item.name+": "+err.Error()))
		}
	}
	return errors.Join(errs...)
}
