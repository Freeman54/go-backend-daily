# 缓存写入合并

热点数据更新时，多个 goroutine 可能同时把同一份结果写入远端缓存，增加网络和连接池压力。
`Coalescer` 只让同一 key 的首个调用执行写入，其余并发调用等待并共享错误结果。

```go
coalescer := New()
err := coalescer.Do("user:42", func() error {
	return redis.Set(ctx, "user:42", encoded, ttl).Err()
})
```

不同 key 不会互相阻塞。该模式只适合“同时发生且内容等价”的写操作；若不同调用携带不同版本，
应先比较版本号或使用 CAS，避免先到的旧数据覆盖新数据。写函数本身也应带超时。
