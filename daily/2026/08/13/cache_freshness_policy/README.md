# 缓存新鲜度策略

缓存 TTL 之外常需要一个短暂的“可陈旧服务”窗口：新鲜数据直接返回，陈旧数据可以先返回并异步刷新，彻底过期的数据则必须回源。`Policy` 把这一决策集中成三个明确状态。

```go
switch policy.Classify(entry.StoredAt, time.Now()) {
case freshnesspolicy.Fresh: return entry.Value
case freshnesspolicy.Stale: go refresh(key); return entry.Value
case freshnesspolicy.Expired: return loadSynchronously(key)
}
```

时间戳应来自可信时钟；多机缓存还需考虑时钟偏差。异步刷新应配合请求合并，避免陈旧窗口触发回源风暴。
