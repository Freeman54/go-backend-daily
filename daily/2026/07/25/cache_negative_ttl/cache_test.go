package cache_negative_ttl

import (
	"testing"
	"time"
)

func TestCache_UsesShorterTTLForMissingValues(t *testing.T) {
	now := time.Unix(100, 0)
	c := New[string](time.Minute, 5*time.Second, func() time.Time { return now })
	c.Set("hit", "value", true)
	c.Set("miss", "", false)

	now = now.Add(6 * time.Second)
	if _, ok, found := c.Get("miss"); found || ok {
		t.Fatal("负缓存应已过期")
	}
	if value, ok, found := c.Get("hit"); !found || !ok || value != "value" {
		t.Fatalf("正缓存读取异常: value=%q ok=%v found=%v", value, ok, found)
	}
}

func TestNew_RejectsInvalidTTL(t *testing.T) {
	if c := New[string](0, time.Second, time.Now); c != nil {
		t.Fatal("非正数 TTL 应返回 nil")
	}
}
