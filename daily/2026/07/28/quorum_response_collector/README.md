# 并发收集法定人数响应

读多个副本时，不一定要等待全部请求。`Collect` 并发执行任务，获得指定数量的成功结果后立即返回，并取消尚未结束的任务；当剩余任务已不可能凑齐法定人数时提前失败，并用 `errors.Join` 保留失败原因。

```go
values, err := Collect(ctx, 2, []Task[string]{
	func(ctx context.Context) (string, error) { return readReplica(ctx, "a") },
	func(ctx context.Context) (string, error) { return readReplica(ctx, "b") },
	func(ctx context.Context) (string, error) { return readReplica(ctx, "c") },
})
```

带缓冲结果通道可避免函数返回后发送方阻塞。生产环境还应限制副本并发量，并让下游 I/O 真正响应取消信号。法定人数只保证拿到足够结果；若需要值的一致性，还应增加版本比较或多数值选择。
