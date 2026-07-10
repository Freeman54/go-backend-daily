# 用 outbox batch claim 控制消息中继的并发领取

事务 outbox 只解决“先落库再投递”，并没有解决 relay 进程如何稳定领取消息。若多个 relay 同时扫表，缺少 claim 租约会导致重复发送；若 claim 规则不稳定，又容易让旧消息一直饥饿。

## 设计要点

- 只领取 `AvailableAt <= now` 且租约已过期的消息
- 领取时写入短租约 `ClaimedUntil`，避免并发 relay 重复处理
- 先按可用时间、再按尝试次数排序，让旧消息和低重试消息优先出队

## 示例说明

示例里的 `Claim` 模拟 relay 从 outbox 表中批量挑选可发送消息，并在内存层把租约和尝试次数更新出来。测试覆盖已被占用消息跳过，以及 claim 顺序与租约设置。

## 运行方式

```bash
go test ./daily/2026/07/10/outbox_batch_claim
```
