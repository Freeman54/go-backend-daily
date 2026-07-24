# 重试令牌桶：给失败流量单独设预算

下游故障时，无限制重试会把一次请求放大成多次调用，进一步压垮依赖。`Bucket` 把“是否允许重试”变成独立的并发安全预算：令牌耗尽时应快速放弃重试，固定间隔后再逐步恢复。

```go
bucket, err := retrytokenbucket.NewChecked(20, 100*time.Millisecond, time.Now())
if err != nil {
	log.Fatal(err)
}
if bucket.Take(time.Now()) {
	// 执行一次带退避和 context 超时的重试
}
```

令牌桶应按下游或操作类型隔离，并配合最大尝试次数、指数退避和 jitter。示例接收显式时间，便于确定性测试；生产环境可在调用点传入 `time.Now()`。

运行：

```bash
go test ./daily/2026/07/24/retry_token_bucket
```
