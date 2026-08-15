# 在入口校验 W3C traceparent

链路追踪入口不应盲目信任外部 `traceparent`。畸形 ID、全零 ID、错误长度或非十六进制字符会污染追踪上下文，甚至制造难以关联的伪链路。本例解析 W3C 基础格式，校验 version 00、trace ID、parent ID 和 flags。

校验失败时通常应创建新的本地 trace，而不是让业务请求失败；同时用低基数指标记录无效头。示例刻意只接受 version 00，未来版本可能包含额外字段，生产实现应交给 OpenTelemetry SDK 等符合规范的库处理。

```go
parent, err := tracingtraceparentvalidator.Parse(header)
if err != nil {
    // 启动新 trace，并记录 invalid_traceparent。
}
_ = parent.Sampled()
```

运行测试：

```bash
go test ./daily/2026/08/15/tracing_traceparent_validator
```
