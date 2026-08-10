# 消息重投分类器

消费者不应对所有错误统一重试。`Classifier` 按错误链把结果分为确认、延迟重投和死信，并阻止超过次数或消息年龄预算的无效重试。

```go
c := redelivery.Classifier{MaxAttempts: 5, MaxAge: 10 * time.Minute}
decision := c.Decide(context.DeadlineExceeded, 2, time.Minute)
// decision == redelivery.Retry
```

业务参数错误应包装 `PermanentError` 后进入死信；取消、超时等暂态故障可重投。`attempt` 表示本次投递序号，从 1 开始。生产环境还应记录分类原因，并为死信配置人工修复与重放流程。
