# SQL 参数预算

批量写入和动态 `IN` 查询容易超过数据库或驱动的参数数量上限。`Planner` 根据固定参数、每行参数数和总预算计算安全批次，避免请求到达数据库后才失败。

```go
p := parametbudget.Planner{MaxParameters: 65535}
sizes, err := p.BatchSizes(10000, 8, 3)
// sizes 描述每批应包含的行数。
```

`fixed` 用于租户 ID、时间范围等每条 SQL 只出现一次的参数。调用方应逐批执行并明确事务边界；如果必须原子写入全部数据，应考虑临时表或数据库原生批量导入能力。
