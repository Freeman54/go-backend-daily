# 消息属性白名单

消息中间件的 attributes/header 经常跨越多个服务。无条件透传会泄露内部调试信息，也可能让生产者伪造消费者信任的控制字段。`Filter` 只复制明确允许的属性，并返回独立 map。

```go
safe := attributeallowlist.Filter(message.Attributes, []string{
    "traceparent",
    "tenant",
})
```

白名单应由消费者契约定义，并限制字段长度和值域。鉴权身份、重试次数等可信字段应由网关或消息系统重新生成，不能直接信任外部生产者提供的值。
