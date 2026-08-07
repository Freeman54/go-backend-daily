# SQL 投影白名单

绑定参数只能保护值，不能安全地绑定列名。`Select` 接收调用方请求的列和服务端白名单，拒绝未知列与重复列后再引用标识符，从而可安全拼进 `SELECT` 子句。

```go
columns, err := sqlprojectionallowlist.Select(
    []string{"id", "created_at"},
    map[string]struct{}{"id": {}, "created_at": {}},
)
// "id", "created_at"
```

白名单应由代码维护，不要由客户端提交。表名、排序方向、过滤操作符同样需要独立的枚举/白名单策略。
