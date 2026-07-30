# 缓存 TTL 时间轮

逐条创建定时器会带来大量定时器对象和调度开销。`Wheel` 把过期任务放进固定时间槽，每次 `Advance` 只检查当前槽；超过一轮的长 TTL 在槽触发时继续调度。

```go
w := cachettlwheel.New(time.Second, 60, time.Now())
w.Set("session:42", "payload", 90*time.Second)
expiredKeys := w.Advance(time.Now())
```

每次覆盖写都会生成版本号，因此旧槽记录不会误删新值。时间轮以一个 tick 为精度，适合允许近似过期的本地缓存；严格过期仍由 `Get` 的截止时间检查保证。示例使用字符串值，实际项目可改为泛型并由单独 goroutine 驱动 `Advance`。
