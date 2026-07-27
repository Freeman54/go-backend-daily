# 用 Context 竞争首个成功结果

向多个等价副本并发查询时，最先返回的不一定是成功结果。`Run` 会并发执行所有任务，忽略先到的错误，
返回第一个成功值，并立即通过派生 `context` 取消剩余任务；若全部失败，则用 `errors.Join` 保留所有错误。

```go
value, err := Run(ctx, []Task[string]{
	func(ctx context.Context) (string, error) { return queryReplica(ctx, "a") },
	func(ctx context.Context) (string, error) { return queryReplica(ctx, "b") },
})
```

结果通道容量等于任务数，因此即使函数已经返回，晚到的 goroutine 也不会因发送结果而阻塞。实际系统还应
限制并发副本数，并确保下游调用真正监听 `ctx.Done()`，否则取消信号无法节省资源。
