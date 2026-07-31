# 截止时间余量守卫

请求虽然尚未超时，但剩余几十毫秒时继续访问下游，通常只会制造无效负载。`Check` 在启动昂贵阶段之前验证 `context` 是否携带截止时间，以及剩余余量是否达到最低要求。

```go
if err := deadlineslackguard.Check(ctx, time.Now(), 200*time.Millisecond); err != nil {
	return fmt.Errorf("skip database query: %w", err)
}
rows, err := db.QueryContext(ctx, query)
```

最低余量应根据下游延迟分位数、网络抖动和响应序列化成本设定。传入 `now` 而非在函数内部读取时钟，能够让规则易于测试；守卫不能替代下游调用自身的超时配置。
