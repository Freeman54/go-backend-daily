package cachetombstone

import "time"

type Store struct {
	now   func() time.Time
	items map[string]entry
}

type entry struct {
	value          string
	version        int64
	tombstoneUntil time.Time
	deleted        bool
}

func New(now func() time.Time) *Store {
	if now == nil {
		now = time.Now
	}
	return &Store{
		now:   now,
		items: make(map[string]entry),
	}
}

func (s *Store) Put(key, value string, version int64) bool {
	current, ok := s.items[key]
	if ok && current.deleted && !s.now().Before(current.tombstoneUntil) {
		delete(s.items, key)
		ok = false
	}
	if ok && current.deleted && s.now().Before(current.tombstoneUntil) && version <= current.version {
		return false
	}
	if ok && version < current.version {
		return false
	}
	s.items[key] = entry{
		value:   value,
		version: version,
	}
	return true
}

func (s *Store) Delete(key string, version int64, ttl time.Duration) {
	current, ok := s.items[key]
	if ok && version < current.version {
		return
	}
	s.items[key] = entry{
		version:        version,
		deleted:        true,
		tombstoneUntil: s.now().Add(ttl),
	}
}

func (s *Store) Get(key string) (string, bool) {
	item, ok := s.items[key]
	if !ok {
		return "", false
	}
	if item.deleted && s.now().Before(item.tombstoneUntil) {
		return "", false
	}
	if item.deleted {
		delete(s.items, key)
		return "", false
	}
	return item.value, true
}
