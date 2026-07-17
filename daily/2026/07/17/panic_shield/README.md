# 用 panic shield 把后台任务崩溃转成可观测错误

后端里常见的异步任务、消费者 handler、批处理 worker 往往跑在 goroutine 中，一旦内部 `panic`，如果没有统一兜底，轻则丢一批任务，重则直接把进程打挂。`panic shield` 的做法是在执行边界统一 `recover`，把崩溃转换成普通错误，再交给重试、告警和日志链路处理。

## 设计要点

- 把 `panic` 恢复逻辑放在 worker 执行入口，而不是散落在业务代码里
- 返回带有 `panic value` 和 `stack` 的错误对象，方便告警和排查
- 正常错误直接透传，不改变业务层的错误语义

## 示例说明

示例提供 `Execute(ctx, fn)`：

- `fn` 正常返回错误时原样透传
- `fn` 内部发生 `panic` 时返回 `PanicError`
- `PanicError` 中保留栈信息，便于日志系统和报警平台消费

适合消息消费、后台补偿任务、批量同步作业等需要“单任务失败可控”的场景。

## 运行方式

```bash
go test ./daily/2026/07/17/panic_shield
```
