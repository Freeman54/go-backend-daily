# 用 keyed retry jitter 打散消费者重试风暴

消息消费失败后如果所有实例都按固定指数退避重试，同一批失败消息常常会在完全一致的时间点再次涌入，形成周期性流量尖峰。`keyed retry jitter` 的思路是在指数退避基础上叠加稳定抖动：同一条消息重试节奏可重复，不同消息之间又能自然打散。

## 设计要点

- 延迟按 `base * 2^(attempt-1)` 增长，并设置最大上限
- 抖动不是纯随机，而是根据 `key + attempt` 计算稳定偏移
- 同一个 key 在重启后仍能得到一致延迟，便于排查和回放

## 示例说明

示例提供 `Delay(key, attempt, base, max, jitterRatio)`，适合：

- MQ 消费失败后的重试调度
- outbox 投递任务的退避重放
- webhook 重发时分散对下游的冲击

测试覆盖指数增长封顶、同 key 稳定性和抖动区间约束。

## 运行方式

```bash
go test ./daily/2026/07/17/keyed_retry_jitter
```
