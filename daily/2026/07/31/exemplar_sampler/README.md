# 指标 Exemplar 稳定采样

Exemplar 把某次指标观测关联到 trace ID，帮助从延迟直方图跳转到具体链路，但每次都附加会增加存储和基数压力。`Sampler` 对 trace ID 做稳定哈希，以固定百分比采样：同一条链路在不同观测点得到一致决定。

```go
sampler := exemplarsampler.New(5)
if sampler.Sample(traceID) {
	histogram.ObserveWithExemplar(seconds, map[string]string{
		"trace_id": traceID,
	})
}
```

采样率只控制 exemplar 数量，不应作为请求 trace 的首采样策略。不要把用户 ID、订单号等高基数字段放入普通指标标签；哈希算法用于稳定分桶，不提供安全匿名化能力。
