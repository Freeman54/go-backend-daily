# 自适应批量大小

固定批量在低负载时可能浪费吞吐，在数据库变慢时又可能放大尾延迟。`Sizer` 根据最近一次批处理耗时做
“加性增、加性减”：低于目标就增加一档，高于目标就减少一档，并始终限制在安全边界内。

```go
sizer, _ := New(100, 10, 500, 200)
batchSize := sizer.Observe(lastLatency.Milliseconds())
rows := queue.Take(batchSize)
```

真实系统可用滑动窗口或 EWMA 代替单次观测，并同时考虑错误率、锁等待和下游限流。`Sizer` 内部加锁，
允许多个 worker 共享，但批量上限仍应由数据库参数数量、消息大小等硬约束决定。
