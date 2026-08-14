# 缓存刷新失败退避

热点缓存过期后，如果上游持续失败，每个请求都立即触发刷新会进一步放大故障。`Policy` 按连续失败次数计算有上限的指数退避时间，`Ready` 用于判断当前请求是否可以再次尝试刷新。

```go
policy := refreshbackoff.Policy{Base: time.Second, Max: time.Minute}
retryAt, err := policy.Next(time.Now(), consecutiveFailures)
if err == nil && refreshbackoff.Ready(time.Now(), retryAt) {
    // 尝试刷新
}
```

生产实现通常还会叠加随机抖动，并在退避期间返回可接受的旧值。成功刷新后必须把连续失败次数归零。
