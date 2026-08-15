# 用 correlation tracker 管理消息 RPC 响应

通过消息队列实现请求-响应时，调用方需要用 correlation ID 把异步响应交还给等待者。本例用互斥锁保护等待表，用容量为 1 的 channel 投递结果；响应到达后先从表中删除，保证重复响应只有一个成功。

必须先 `Register` 再发布消息。调用方超时或 `context` 取消后应执行 `Cancel`，否则等待表会持续增长。真实系统还要给 correlation ID 加实例前缀，并记录未知、重复响应指标。

```go
resultCh, err := tracker.Register(requestID)
// 注册成功后发布消息，再 select resultCh 与 ctx.Done()。
```

运行测试（并发检查建议加 `-race`）：

```bash
go test -race ./daily/2026/08/15/message_correlation_tracker
```
