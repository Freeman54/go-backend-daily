# 延迟自适应采样

固定比例采样可能漏掉低频错误和慢请求。`Sampler` 强制保留失败或超过延迟阈值的请求，仅对快速成功请求按稳定哈希降采样。

```go
sampler := latency_adaptive_sampler.New(300*time.Millisecond, 100)
if sampler.Keep(elapsed, err != nil, hashRequestID(requestID)) {
    exportTrace()
}
```

稳定哈希让同一请求在不同节点得到一致决定，`fastRate=100` 表示大约保留百分之一的普通流量。阈值和比例应由 SLO、流量与存储预算驱动，并监控采样后的估算偏差。

