# SQL 复合游标谓词

深分页使用 `OFFSET` 会扫描并丢弃越来越多的行。Keyset pagination 通过“上一页最后一行”的排序键继续查询。`After` 为多个同向排序字段展开字典序谓词，并把游标值保留为参数。

```go
where, args, err := tuplecursor.After(
    []string{"created_at", "id"},
    []any{lastCreatedAt, lastID},
    true,
)
// (created_at < ?) OR (created_at = ? AND id < ?)
```

列名不能通过占位符传递，因此实现只接受安全标识符。实际查询还应使用匹配的复合索引，并保证最后一个排序字段能够稳定打破平局。
