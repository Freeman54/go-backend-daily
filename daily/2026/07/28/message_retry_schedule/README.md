# 消息重试时间表

消费者立即重试会持续冲击故障中的依赖。`Policy.Delay` 使用带上限的指数退避，并加入对称抖动，让不同消费者的重试时刻错开。随机样本作为参数传入，便于确定性测试。

```go
delay, err := policy.Delay(message.Attempt, rand.Float64())
if err != nil {
	return err
}
retryQueue.PublishAt(message, time.Now().Add(delay))
```

重试次数必须持久化，达到上限后转入死信队列。业务错误与临时错误应分类处理；参数校验失败、权限拒绝等永久错误通常不应重试。调度时间也要写入指标，避免队列积压被隐藏。
