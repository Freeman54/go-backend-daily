# 消息批次提交屏障

并发消费消息时，较大的 offset 可能先处理完成。若直接提交它，进程崩溃后尚未完成的较小 offset 会被跳过。`Barrier` 暂存乱序确认，只返回最高的连续已完成 offset。

```go
barrier := message_batch_barrier.New(100)
barrier.Ack(102) // 返回 99，仍有缺口
barrier.Ack(100) // 返回 100
commit := barrier.Ack(101) // 返回 102，可安全提交
```

`Ack` 支持并发调用并忽略重复、过期确认。实际消费者还应限制最大 in-flight 数，避免某条消息长期阻塞时暂存集合无限增长。
