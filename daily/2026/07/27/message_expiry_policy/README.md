# 消息过期时间传播策略

消息经过重试、延迟队列或跨服务转发后，如果每一跳都重新设置完整 TTL，过期业务可能被无限延长。
`Policy` 根据消息最初创建时间和允许的最大年龄计算剩余 TTL，到期后明确拒绝继续投递。

```go
policy, _ := New(10 * time.Minute)
ttl, err := policy.RemainingTTL(message.CreatedAt, time.Now())
if errors.Is(err, ErrExpired) {
	return sendToDeadLetter(message)
}
return broker.Publish(message, ttl)
```

剩余 TTL 应写入消息代理原生的过期配置，而 `CreatedAt` 应在首次生产时生成并跨重试保持不变。跨机器时钟
可能存在偏差，因此示例拒绝“当前时间早于创建时间”的输入；生产系统可配合时钟同步和有限偏差容忍策略。
