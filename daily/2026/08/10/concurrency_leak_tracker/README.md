# 并发任务泄漏追踪器

仅看 goroutine 总数很难定位谁没有退出。`Tracker` 为后台任务登记名称和启动时间，并生成按年龄降序排列的快照，让健康检查或诊断接口发现超时未结束的任务。

```go
tracker := leaktracker.New()
done := tracker.Start("refresh-cache")
defer done()

stuck := tracker.OlderThan(time.Now(), 30*time.Second)
```

返回的完成函数可重复调用，适合多个退出分支统一 `defer`。任务名称不应包含用户隐私或高基数字段；诊断接口也应限制访问权限。追踪器用于发现问题，不能替代 `context` 取消和可靠的关闭协议。
