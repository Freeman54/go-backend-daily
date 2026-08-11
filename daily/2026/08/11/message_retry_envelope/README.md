# 消息重试信封校验

重试次数和首次失败时间通常来自消息头，属于跨进程的不可信输入。`Decide` 同时限制最大尝试次数与总重试年龄，并容忍少量时钟偏差，防止损坏元数据造成无限重试。

```go
decision, err := retryenvelope.Decide(
	time.Now(), envelope, 5, 24*time.Hour, 10*time.Second,
)
if err != nil || decision == retryenvelope.DeadLetter {
	// 投递死信队列并记录原因。
}
```

消费者应在每次重新投递时递增 `Attempt`，但保留最初的 `FirstFailed`。死信处理仍需携带业务消息 ID，便于审计与人工重放。
