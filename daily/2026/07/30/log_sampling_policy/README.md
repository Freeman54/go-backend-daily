# 确定性日志采样策略

高流量服务若记录每个成功请求，日志成本和噪声都会迅速增长。`Policy` 无条件保留强制、服务端错误和慢请求，其余事件按请求 ID 做确定性哈希采样。

```go
policy, _ := logsamplingpolicy.New(0.01)
if policy.Keep(logsamplingpolicy.Event{
	RequestID: requestID, StatusCode: status, DurationMillis: elapsed.Milliseconds(),
}) {
	logger.Info("request completed")
}
```

确定性决策让同一请求在重试或多处记录时保持一致，更利于串联排查。请求 ID 必须足够均匀；空值会让所有无 ID 事件作出同一决策。采样不能用于审计日志，并应把采样率作为字段写入指标或日志以支持流量估算。
