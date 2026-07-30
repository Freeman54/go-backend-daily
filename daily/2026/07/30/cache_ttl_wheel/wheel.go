package cachettlwheel

import (
	"sort"
	"sync"
	"time"
)

type entry struct {
	value     string
	expiresAt time.Time
	version   uint64
}

type scheduled struct {
	key     string
	version uint64
}

// Wheel 将到期检查分散到固定数量的时间槽。
type Wheel struct {
	mu      sync.Mutex
	tick    time.Duration
	now     time.Time
	cursor  int
	version uint64
	entries map[string]entry
	slots   [][]scheduled
}

func New(tick time.Duration, slots int, now time.Time) *Wheel {
	if tick <= 0 || slots < 1 {
		panic("tick and slots must be positive")
	}
	return &Wheel{tick: tick, now: now, entries: make(map[string]entry), slots: make([][]scheduled, slots)}
}

func (w *Wheel) Set(key, value string, ttl time.Duration) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.version++
	expiresAt := w.now.Add(ttl)
	w.entries[key] = entry{value: value, expiresAt: expiresAt, version: w.version}
	steps := int((ttl + w.tick - 1) / w.tick)
	if steps < 1 {
		steps = 1
	}
	slot := (w.cursor + steps) % len(w.slots)
	w.slots[slot] = append(w.slots[slot], scheduled{key: key, version: w.version})
}

func (w *Wheel) Advance(now time.Time) []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	var expired []string
	for !w.now.Add(w.tick).After(now) {
		w.now = w.now.Add(w.tick)
		w.cursor = (w.cursor + 1) % len(w.slots)
		bucket := w.slots[w.cursor]
		w.slots[w.cursor] = nil
		for _, item := range bucket {
			e, ok := w.entries[item.key]
			if !ok || e.version != item.version {
				continue
			}
			if !e.expiresAt.After(w.now) {
				delete(w.entries, item.key)
				expired = append(expired, item.key)
				continue
			}
			steps := int((e.expiresAt.Sub(w.now) + w.tick - 1) / w.tick)
			slot := (w.cursor + steps) % len(w.slots)
			w.slots[slot] = append(w.slots[slot], item)
		}
	}
	sort.Strings(expired)
	return expired
}

func (w *Wheel) Get(key string, now time.Time) (string, bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	e, ok := w.entries[key]
	if !ok || !e.expiresAt.After(now) {
		return "", false
	}
	return e.value, true
}
