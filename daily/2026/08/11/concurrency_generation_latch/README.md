# 并发代际闩锁

配置热更新、服务发现和缓存失效需要反复广播状态变化。一次性关闭 channel 只能通知一轮；`Latch` 每次推进代际都会关闭旧 channel 并创建新 channel，让当前等待者全部唤醒，同时避免检查状态与开始等待之间丢失通知。

```go
observed := latch.Current()
generation, err := latch.Wait(ctx, observed)
if err == nil {
	// 读取 generation 对应的新快照。
}
```

调用方应先保存已处理的代际，再进入等待。闩锁只负责变化通知；具体状态应存放在独立、可原子读取的快照中。
