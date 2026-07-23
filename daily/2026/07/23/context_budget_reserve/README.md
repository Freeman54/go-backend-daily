# 为收尾动作预留 Context 时间预算

一个请求的总 deadline 若全部交给主要下游调用，超时后往往已没有时间记录审计、回滚临时资源或返回结构化错误。`WithReserve` 从父 context 的截止时间中扣除一段保留时间，生成更早结束的子 context。

业务调用使用子 context，调用结束后仍可在父 context 剩余的 reserve 内执行必要收尾。没有截止时间、保留时间为负数或剩余预算不足都会明确报错，避免创建一个已经过期的 context。保留时间应依据收尾操作的延迟分位数配置，而不是随意写死。

```go
workCtx, cancel, err := contextbudgetreserve.WithReserve(requestCtx, 100*time.Millisecond)
defer cancel()
```

运行示例：

```bash
go test ./daily/2026/07/23/context_budget_reserve
```
