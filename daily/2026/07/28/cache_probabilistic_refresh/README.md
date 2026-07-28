# 缓存概率提前刷新

大量请求同时看到缓存过期会形成击穿。`ShouldRefresh` 在 TTL 尾部开启提前刷新窗口，并依据剩余 TTL 比例与随机样本决定是否刷新：越接近过期，触发概率越高，从而把刷新压力分散到一段时间内。

```go
if ShouldRefresh(time.Now(), entry.ExpiresAt, entry.TTL, rand.Float64(), 0.2) {
	go refreshWithSingleflight(key)
}
```

概率刷新不能替代请求合并；应配合 singleflight、分布式锁或后台刷新，避免多个实例仍然同时回源。随机源应复用而非每次重新播种，过期数据是否可返回则应由一致性要求决定。
