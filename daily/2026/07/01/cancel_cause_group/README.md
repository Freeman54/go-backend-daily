# 用 cancel cause 向并发任务传播失败原因

多个后台任务共享一个请求上下文时，普通 `cancel` 只能告诉其他协程“结束了”，但无法说明为什么结束。Go 1.20 引入的 `context.WithCancelCause` 可以把首个失败原因传给所有协程，便于日志和指标归因。

## 设计要点

- 首个失败任务触发 `cancel(err)`，其余任务通过 `context.Cause` 读取原因
- 成功路径保持 `nil`，避免把正常完成误判为错误
- 适合并发拉取、并行校验、批处理拆片等共享取消语义的场景

## 示例说明

实际项目里可以把这个模式放进 fan-out 查询、异步回写或多阶段校验流程。相比只看 `context.Canceled`，它能让日志里直接带出根因。

## 运行方式

```bash
go test ./daily/2026/07/01/cancel_cause_group
```
