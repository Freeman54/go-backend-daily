package aftercommithooks

import (
	"context"
	"errors"
	"sync"
)

var ErrAlreadyCommitted = errors.New("transaction already committed")

type Hook func(context.Context) error

type Unit struct {
	mu        sync.Mutex
	hooks     []Hook
	committed bool
}

func (u *Unit) OnCommit(h Hook) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.committed {
		return ErrAlreadyCommitted
	}
	u.hooks = append(u.hooks, h)
	return nil
}

func (u *Unit) Commit(ctx context.Context) error {
	u.mu.Lock()
	if u.committed {
		u.mu.Unlock()
		return ErrAlreadyCommitted
	}
	u.committed = true
	hooks := append([]Hook(nil), u.hooks...)
	u.mu.Unlock()

	var joined error
	for _, hook := range hooks {
		if err := hook(ctx); err != nil {
			joined = errors.Join(joined, err)
		}
	}
	return joined
}
