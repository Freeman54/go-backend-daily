# 消息分区键校验

分区键决定同一业务实体的事件是否能进入同一分区并维持顺序。`Validate` 约束键非空、长度上限及跨语言稳定的 ASCII 字符集（字母、数字、`-_.:`）。

```go
if err := messagepartitionkey.Validate("tenant-42:order-99", 64); err != nil {
    // 拒绝写入，避免不可预期的分区或键爆炸。
}
```

不要把随机 UUID、完整 JSON 或用户自由文本直接作为键；应使用稳定的租户与实体标识组合，并结合实际 Broker 的键长度限制配置上限。
