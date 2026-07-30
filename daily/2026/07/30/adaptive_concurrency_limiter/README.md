# 自适应并发限制器

固定并发数无法同时适应低负载和下游退化。`Limiter` 按当前并发上限组成采样窗口：窗口内全部请求快速成功时加一；出现超时延或失败时按惩罚数收缩，并始终限制在配置范围内。

```go
l := adaptiveconcurrencylimiter.New(8, 2, 32, 100*time.Millisecond)
if !l.TryAcquire() {
	// 快速拒绝或排队
	return
}
started := time.Now()
err := callDependency()
l.Release(time.Since(started), err == nil)
```

获取与释放必须成对出现。生产系统还应加入滑动分位数、冷却时间和指标，避免短窗口导致振荡；限制器只保护并发量，不能替代超时、重试预算与熔断。
