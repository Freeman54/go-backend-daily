package cache_tag_index

import (
	"reflect"
	"sort"
	"sync"
	"testing"
)

func TestIndexInvalidate_ReturnsUniqueKeysAndClearsTag(t *testing.T) {
	index := New()
	index.Attach("post:1", "author:7", "feed")
	index.Attach("post:1", "author:7")
	index.Attach("post:2", "author:7")
	got := index.Invalidate("author:7")
	sort.Strings(got)
	if want := []string{"post:1", "post:2"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if again := index.Invalidate("author:7"); len(again) != 0 {
		t.Fatalf("二次失效返回 %v", again)
	}
}

func TestIndexAttach_IsConcurrentSafe(t *testing.T) {
	index := New()
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); index.Attach("shared", "tag") }()
	}
	wg.Wait()
	if got := index.Invalidate("tag"); len(got) != 1 || got[0] != "shared" {
		t.Fatalf("got %v", got)
	}
}
