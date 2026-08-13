# 可取消并发屏障

并发任务组遇到首个错误时应尽快通知兄弟任务停止，但调用方仍需等待所有 goroutine 完成清理后再返回。`Run` 用带原因的 context 广播取消，用 `sync.Once` 保留首个业务错误，并用 `WaitGroup` 形成退出屏障。

```go
err := cancelbarrier.Run(ctx,
    func(ctx context.Context) error { return fetchProfile(ctx) },
    func(ctx context.Context) error { return fetchOrders(ctx) },
)
```

任务函数必须监听 `ctx.Done()`，否则取消只能发出信号，无法强制终止 goroutine。返回前等待清理可避免资源仍在后台使用、测试泄漏或服务关闭时遗留工作。
