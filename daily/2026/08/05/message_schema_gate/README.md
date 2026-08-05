# 消息模式版本门禁

滚动升级期间，消费者应明确声明可处理的消息模式版本。`Gate` 区分可接受、过旧和过新的消息，便于分别丢弃历史垃圾与隔离尚未支持的新格式。

```go
gate, _ := message_schema_gate.New(2, 4)
decision := gate.Check(5) // Future
```

版本号只能表达兼容区间；字段兼容性仍应通过 Protobuf/Avro 规则和契约测试保证。
