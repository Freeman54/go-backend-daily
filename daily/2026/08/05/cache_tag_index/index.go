package cache_tag_index

import "sync"

type Index struct {
	mu   sync.Mutex
	tags map[string]map[string]struct{}
}

func New() *Index {
	return &Index{tags: make(map[string]map[string]struct{})}
}

func (i *Index) Attach(key string, tags ...string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	for _, tag := range tags {
		if i.tags[tag] == nil {
			i.tags[tag] = make(map[string]struct{})
		}
		i.tags[tag][key] = struct{}{}
	}
}

func (i *Index) Invalidate(tag string) []string {
	i.mu.Lock()
	defer i.mu.Unlock()
	keys := i.tags[tag]
	delete(i.tags, tag)
	result := make([]string, 0, len(keys))
	for key := range keys {
		result = append(result, key)
	}
	return result
}
