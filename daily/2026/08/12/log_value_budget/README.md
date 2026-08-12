# 日志字段值预算

结构化日志中的错误堆栈、请求片段或第三方响应可能异常巨大，抬高采集成本并挤占正常日志。`Truncate` 按字节限制字段值，同时避免切断 UTF-8 字符。

```go
logger.Info("dependency failed", valuebudget.Fields(map[string]string{
    "response": body,
    "error":    err.Error(),
}, 1024))
```

字节预算比字符预算更贴近日志管道的传输与存储成本。敏感字段仍应先脱敏，再做截断；否则前缀中仍可能泄露凭据或个人信息。
