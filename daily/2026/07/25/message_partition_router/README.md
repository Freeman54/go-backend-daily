# 消息分区路由

有序消息通常按业务 key（如租户、订单）固定路由到同一分区。`Router` 使用标准库 FNV-1a 哈希，
相同 key 在分区数不变时会得到稳定且有界的分区编号。

```go
router, _ := New(32)
partition := router.Partition(order.TenantID)
producer.Send(partition, order)
```

扩缩分区会改变取模结果，因此不能把它误认为一致性哈希；需要扩容稳定性时应使用虚拟节点或显式路由表。
空 key 会集中到一个分区，调用方应保证路由 key 具备足够基数，并监控分区流量倾斜。
