# SQL 分页上限策略

直接信任客户端的 `limit` 会产生超大结果集，放大数据库、网络和序列化成本。`Policy` 统一默认值和硬上限，避免各接口各自处理边界。

```go
policy, err := sql_limit_policy.New(20, 100)
if err != nil { panic(err) }
limit := policy.Normalize(request.Limit)
rows, err := db.QueryContext(ctx, "SELECT id FROM users LIMIT ?", limit)
```

分页上限只是保护措施。大表查询仍应使用确定排序和游标分页，并为数据库调用设置超时；不要用字符串拼接把用户输入放入 SQL。
