# 消息确认租约

耗时消费者必须在可见性超时前续租，否则消息可能在处理未完成时被再次投递。`Lease` 在截止时间减去安全边界后提示续租，成功续租后更新本地截止时间，确认完成后永久停止续租。

```go
lease, _ := messageacklease.New(time.Now(), 30*time.Second, 5*time.Second)
if lease.ShouldExtend(time.Now()) {
	if err := broker.Extend(messageID, 30*time.Second); err == nil {
		lease.Extended(time.Now(), 30*time.Second)
	}
}
// 业务提交并成功 ack 后
lease.Ack()
```

只有 broker 确认续租成功后才能调用 `Extended`。真实消费者还需要限制最大处理时长、记录续租失败指标，并保证业务幂等；本地租约状态无法消除网络分区下的重复投递。
