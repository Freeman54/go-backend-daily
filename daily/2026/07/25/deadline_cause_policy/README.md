# 带原因的超时策略

普通 `context.WithTimeout` 只能让调用方看到 `DeadlineExceeded`，难以判断是哪一层预算耗尽。Go 1.21 的
`context.WithTimeoutCause` 可以保留稳定的业务原因，日志和指标可通过 `context.Cause` 分类。

示例中的 `NewTimeoutCause` 校验超时时间和原因，适合作为服务内部的统一入口：

```go
ctx, cancel, err := NewTimeoutCause(parent, 200*time.Millisecond, errors.New("库存查询超时"))
if err != nil { return err }
defer cancel()

if err := query(ctx); err != nil {
    log.Printf("query stopped: %v", context.Cause(ctx))
}
```

注意：原因应是低基数、可归类的错误，不要把用户数据拼进错误文本；仍需及时调用 `cancel` 释放定时器。
