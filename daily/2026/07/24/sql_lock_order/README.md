# SQL 锁顺序规划：用确定性顺序降低死锁概率

两个事务以不同顺序锁定相同记录时，容易形成循环等待。`Plan` 对资源 ID 去重并升序排列；`ForUpdateQuery` 生成参数化 PostgreSQL 查询，让所有调用方遵循一致的加锁顺序。

```go
query, ids, err := sqllockorder.ForUpdateQuery("accounts", []int64{9, 2, 9})
rows, err := tx.QueryContext(ctx, query, ids[0], ids[1])
```

表名不能作为 SQL 参数，因此示例只接受安全标识符；实际项目更适合从代码内固定枚举表名。统一顺序能降低死锁概率，但仍应控制事务时长，并对数据库返回的死锁错误做有界重试。

运行：

```bash
go test ./daily/2026/07/24/sql_lock_order
```
