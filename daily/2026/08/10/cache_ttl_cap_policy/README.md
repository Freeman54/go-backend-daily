# 缓存 TTL 上限策略

业务代码直接采用来源 TTL 时，异常配置可能制造长期脏数据。`Policy.Apply` 把期望 TTL 限制在统一上下界内，并为负缓存使用更严格的上限。

```go
policy := ttlcap.Policy{Min: time.Second, Max: time.Hour, NegativeMax: time.Minute}
ttl, err := policy.Apply(24*time.Hour, false)
// ttl == time.Hour
```

下限可避免接近零的 TTL 引发缓存穿透，上限限制陈旧窗口。负缓存应更短，以免“数据刚创建但仍显示不存在”。策略校验应在配置加载阶段完成。
