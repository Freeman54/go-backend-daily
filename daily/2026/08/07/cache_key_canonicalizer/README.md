# 缓存键规范化

同一个业务对象若因 map 遍历顺序不同而生成不同缓存键，会产生隐性穿透。`Build` 使用 `url.Values.Encode` 按名称排序并转义维度，确保等价输入落在同一键上。

```go
key, err := cachekeycanonicalizer.Build("product", map[string]string{
    "locale": "zh CN", "id": "42",
})
// product?id=42&locale=zh+CN
```

命名空间和维度名称不能为空。实际项目还应把租户、版本等隔离维度作为显式 part，避免不同数据域误共享缓存。
