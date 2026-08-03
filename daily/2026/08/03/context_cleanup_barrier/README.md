# Context 感知的清理栅栏

服务退出时通常要并行刷新缓冲区、提交 offset、关闭连接。`Wait` 并发执行清理函数，收集所有已完成错误，并在调用方的 context 到期时及时返回。

```go
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()
err := context_cleanup_barrier.Wait(ctx, flushLogs, commitOffsets, closePool)
```

带缓冲的结果通道确保超时返回后任务发送结果不会阻塞。注意：context 只限制“等待时间”，无法强制终止不接受 context 的清理函数；任务自身仍应实现取消或明确的超时。

