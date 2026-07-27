# SQL 批量 VALUES 占位符

批量插入时应继续使用参数化 SQL，不能把业务值拼进语句。`Build` 根据列数、行数和起始编号生成
PostgreSQL 占位符，并返回下一个可用编号，方便前面已有查询参数的场景。

```go
clause, next, err := Build(3, 2, 2)
// clause: "($2,$3,$4),($5,$6,$7)"，next: 8
query := "INSERT INTO users(name,email,age) VALUES " + clause
_, err = db.ExecContext(ctx, query, args...)
```

占位符结构由可信的整数维度生成，业务数据仍通过 `args` 绑定。生产代码还要限制行数，使参数总数、
SQL 长度和事务时间不超过数据库及连接池的安全边界。
