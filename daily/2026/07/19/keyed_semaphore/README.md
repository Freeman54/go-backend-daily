# 用 keyed semaphore 控制热点资源的并发度

很多后端场景不是要限制“全局并发”，而是要限制“同一个 key 的并发”。例如同一个租户的批处理、同一个订单的补偿任务、同一条配置的刷新操作。如果只做全局限流，热点 key 仍然可能把下游打爆。

`keyed semaphore` 的思路是：每个 key 拥有独立的并发配额，互不影响。冷门 key 不会被热门 key 挤压，热点 key 也不能无限膨胀。

## 设计要点

- 每个 key 独立计数，适合订单号、租户 ID、分片 ID 等维度
- `Acquire` 支持 `context.Context`，超时或取消时能及时放弃排队
- `Release` 归还当前 key 的令牌，避免热点资源长期占住配额

## 示例说明

示例实现提供：

- `New(limit)`：创建每个 key 共用同一上限的 limiter
- `Acquire(ctx, key)`：申请某个 key 的一个并发槽位
- `Release(key)`：释放槽位

适合热点 key 保护、分区消费限流、单租户隔离等场景。

## 运行方式

```bash
go test ./daily/2026/07/19/keyed_semaphore
```
