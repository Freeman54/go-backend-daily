# 缓存 Stale-if-error 兜底

缓存刚过期时直接删除，会让一次短暂的数据库或依赖故障立即暴露给用户。`Cache` 将生命周期拆为新鲜窗口
和旧值兜底窗口：新鲜值直接命中；过期后尝试加载；加载失败且旧值仍在兜底窗口内时返回旧值。

```go
cache, _ := New[User](time.Minute, 10*time.Minute)
user, err := cache.Get("user:42", time.Now(), func() (User, error) {
	return repository.FindUser(42)
})
```

成功刷新会替换旧值，超过兜底窗口后则原样返回加载错误。生产实现通常还要增加 singleflight、防止缓存
击穿，并记录“返回旧值”的指标；对余额、权限等强一致数据不应静默使用旧值。
