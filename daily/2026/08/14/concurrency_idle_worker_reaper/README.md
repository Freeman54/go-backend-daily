# 并发空闲 Worker 回收

弹性 Worker 池在流量高峰后需要缩容，但不能一次回收过多导致下一波请求排队。`Select` 按最后活动时间选择空闲 Worker，并保证至少保留指定数量。

```go
ids, err := workerreaper.Select(time.Now(), workers, 30*time.Second, 2)
for _, id := range ids {
    // 通知对应 Worker 在完成当前任务后退出
}
```

选择和真正退出应分成两个阶段：先标记 Worker 为 draining，停止分配新任务，再等待当前任务完成。生产实现还需要用锁或单线程控制循环保护 Worker 状态。
