# 消息批次提交点

并发消费时，offset 完成顺序通常与读取顺序不同。若收到较大 offset 的成功结果就直接提交，进程崩溃后可能永久跳过尚未完成的较小 offset。`Tracker` 暂存乱序确认，只把提交点推进到连续完成区间之后。

```go
checkpoint := tracker.Ack(message.Offset)
consumer.Commit(checkpoint) // exclusive：表示 checkpoint 之前都已完成
```

示例使用内存状态，重启后由 broker 的已提交位置重新构建。真实系统还需限定在单个 partition 内使用，并将业务副作用设计为幂等。
