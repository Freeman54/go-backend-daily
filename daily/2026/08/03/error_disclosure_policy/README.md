# 错误披露策略

把数据库地址、SQL 或内部调用链直接返回给客户端会泄漏系统细节。此示例用专用错误类型标记可公开消息；普通错误统一映射为稳定的兜底文案，同时内部仍可通过 `%w` 保留完整错误链供日志记录。

```go
err := error_disclosure_policy.WrapForOperation(
    "reserve inventory",
    error_disclosure_policy.Safe("库存暂不可用"),
)
message := error_disclosure_policy.PublicMessage(err)
```

生产系统可进一步为安全错误增加机器码和 HTTP 状态码。安全标记应只在受控边界创建，不能用用户输入或下游原始报错直接构造。

