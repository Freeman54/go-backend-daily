# 结构化日志字段脱敏

结构化日志常混入令牌、密码和邮箱。`Redactor` 按不区分大小写的字段名删除秘密，并对邮箱保留域名以兼顾排障和隐私。

```go
r := log_field_redactor.New("password", "authorization")
safe := r.Redact(map[string]any{"email": "alice@example.com"})
```

返回值是新 map，不修改调用方数据。生产系统还应限制嵌套结构、请求体和异常文本中的敏感信息。
