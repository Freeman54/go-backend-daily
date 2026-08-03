# SQL IN 参数分块

数据库和驱动通常限制单条语句的绑定参数数量。`Plan` 按上限切分 ID，调用方可对每块执行一次参数化查询，避免拼接 SQL 和超出参数限制。

```go
chunks, err := sql_in_clause_chunks.Plan(userIDs, 500)
for _, ids := range chunks {
    // SELECT ... WHERE id IN (?, ?, ...)
}
```

返回的块共享原始切片底层数组，并用 full slice expression 限制容量，防止向某块 append 时覆盖相邻块。实际批量查询还应考虑事务一致性、结果合并顺序和数据库总往返成本。

