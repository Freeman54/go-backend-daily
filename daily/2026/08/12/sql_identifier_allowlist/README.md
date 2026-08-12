# SQL 标识符白名单

SQL 参数只能绑定值，不能安全绑定列名或表名。`Allowlist` 把 API 暴露的字段名映射为代码内可信的 SQL 标识符，拒绝任意用户输入直接进入查询结构。

```go
fields, _ := identifierallowlist.New(map[string]string{
    "createdAt": "orders.created_at",
})
column, err := fields.Resolve(request.SortBy)
query := "SELECT * FROM orders ORDER BY " + column
```

映射值在初始化时只允许字母、数字、下划线和点号。排序方向、`NULLS FIRST` 等语法应使用独立枚举处理；普通查询值仍必须使用占位符绑定。
