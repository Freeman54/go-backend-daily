# SQL 排序白名单

占位符只能绑定值，不能安全绑定列名和 `ASC`/`DESC`。把 API 的排序参数直接拼接到 SQL 会造成注入风险。`Build` 用业务字段到固定数据库列的白名单映射，并严格枚举排序方向。

```go
clause, err := orderallowlist.Build(request.Sort, map[string]string{
    "createdAt": "created_at", "id": "id",
})
query := "SELECT id, created_at FROM orders " + clause
```

白名单中的列名必须由开发者静态定义，不能来自请求。分页查询应追加唯一列作为稳定的最终排序键，避免相同值在翻页时漂移。
