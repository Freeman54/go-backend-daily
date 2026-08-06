# 给 Tracing Baggage 设置传播预算

Baggage 会跨服务随请求传播，未经限制的键值既会放大网络开销，也可能把调试信息或敏感字段扩散到整条调用链。`Apply` 只按优先级保留允许的键，并同时限制条目数和 UTF-8 字节数。

```go
kept, dropped := tracing_baggage_budget.Apply(
    incoming, []string{"tenant", "region", "experiment"}, 8, 512,
)
```

优先级由调用方显式给出，使预算不足时的取舍可预测。实际接入 OpenTelemetry 时，还应校验键值语法、记录丢弃计数，并禁止凭证及个人信息进入 Baggage。
