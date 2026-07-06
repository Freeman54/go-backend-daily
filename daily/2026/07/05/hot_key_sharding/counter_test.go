package hotkeysharding

import (
	"fmt"
	"sync"
	"testing"
)

func TestCounterSupportsConcurrentAdd(t *testing.T) {
	counter := New(8)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				counter.Add("tenant-a", 1)
			}
		}()
	}

	wg.Wait()

	if got := counter.Get("tenant-a"); got != 10000 {
		t.Fatalf("Get() = %d, want 10000", got)
	}
}

func TestCounterDistributesDifferentKeysAcrossShards(t *testing.T) {
	counter := New(16)
	for i := 0; i < 128; i++ {
		counter.Add(fmt.Sprintf("key-%d", i), 1)
	}

	loads := counter.ShardLoads()
	used := 0
	for _, load := range loads {
		if load > 0 {
			used++
		}
	}

	if used < 8 {
		t.Fatalf("used shards = %d, want at least 8", used)
	}
}
