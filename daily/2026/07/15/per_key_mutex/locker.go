package perkeymutex

import "sync"

type Locker struct {
	mu    sync.Mutex
	locks map[string]*entry
}

type entry struct {
	mu   sync.Mutex
	refs int
}

func New() *Locker {
	return &Locker{
		locks: make(map[string]*entry),
	}
}

func (l *Locker) Lock(key string) func() {
	l.mu.Lock()
	e, ok := l.locks[key]
	if !ok {
		e = &entry{}
		l.locks[key] = e
	}
	e.refs++
	l.mu.Unlock()

	e.mu.Lock()

	return func() {
		e.mu.Unlock()

		l.mu.Lock()
		defer l.mu.Unlock()

		e.refs--
		if e.refs == 0 {
			delete(l.locks, key)
		}
	}
}

func (l *Locker) ActiveKeys() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.locks)
}
