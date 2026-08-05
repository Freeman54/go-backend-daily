# SQL 事务隔离级别选择器

隔离级别应由一致性需求驱动。`Select` 根据是否只读、是否允许不可重复读、是否必须防止写偏差，返回 Go `database/sql` 的事务选项。

```go
options := sql_isolation_selector.Select(sql_isolation_selector.Requirement{ReadOnly: true})
```

数据库对隔离级别的实现并不完全一致，上线前需结合所用数据库验证锁行为与重试策略。
