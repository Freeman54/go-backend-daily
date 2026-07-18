# 用 context snapshot 把必要元数据安全交给异步任务

`context.Context` 适合表达取消、超时和请求范围内的少量元数据，但并不适合整条链路原样传给后台 goroutine。直接复用原始 `context`，容易把取消信号、无关 value、甚至很大的对象一起带过去。更稳妥的做法是提取白名单字段，形成一个轻量快照，再在异步任务里按需挂回去。

## 设计要点

- 只抓取有限且明确的字段，例如 `trace_id`、`user_id`、`tenant`
- 避免把原请求的取消语义直接泄漏到后台任务
- 异步任务拿到的是可序列化、可测试、可审查的上下文快照

## 示例说明

示例提供 `Capture(ctx)` 和 `Snapshot.Attach(ctx)`：

- `Capture` 从原始 `context` 中提取白名单字段
- `Attach` 把快照重新挂到一个新的 `context` 上
- 未在白名单中的值不会被复制

适合审计日志、异步通知、延迟补偿、任务投递前的上下文瘦身。

## 运行方式

```bash
go test ./daily/2026/07/18/context_snapshot
```
