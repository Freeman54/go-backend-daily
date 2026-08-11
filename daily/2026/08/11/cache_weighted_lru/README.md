# 按权重限制的 LRU 缓存

只限制条目数会让一个超大对象挤爆进程内存。`Cache` 为每个值记录业务估算的权重，并在总权重超过容量时逐出最久未访问的条目。

```go
cache, _ := weightedlru.New[string, []byte](64 << 20)
payload := []byte("result")
_ = cache.Put("order:42", payload, int64(len(payload)))
value, ok := cache.Get("order:42")
_, _ = value, ok
```

示例使用调用方提供的近似字节数，避免依赖不可靠的对象深度测量。生产环境应额外暴露当前权重、逐出次数和超重拒绝指标。
