# Goroutine 心跳监测

goroutine 没有崩溃并不代表仍在推进：它可能卡在外部调用、锁或错误的循环中。`Monitor` 记录工作单元最后一次完成进度的时间，定期找出超过阈值未更新的 worker。

```go
monitor := goroutine_heartbeat_monitor.New(30 * time.Second)
monitor.Beat("partition-3", time.Now())
for _, worker := range monitor.Stale(time.Now()) {
	log.Printf("worker %s has made no progress", worker)
}
```

心跳应代表“完成了有效进度”，而非仅代表循环仍在运行。生产环境可把 stale 数量导出为指标并触发告警，但自动重启前要确认任务具备幂等性，避免重复副作用。
