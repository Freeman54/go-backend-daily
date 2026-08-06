# 安全控制 SQL 的 NULL 排序

排序字段不能用参数占位符绑定，也不应直接拼接用户输入。`BuildOrderBy` 先把 API 暴露的排序键映射为可信列名，再限定方向与 `NULLS FIRST/LAST` 策略。

```go
clause, err := sql_null_sort_policy.BuildOrderBy(
    "updated", sql_null_sort_policy.Descending, sql_null_sort_policy.NullsLast,
    map[string]string{"updated": "updated_at"},
)
```

显式 NULL 策略能避免升降序切换或数据库迁移后分页顺序悄然变化；确实要沿用数据库默认值时可传 `NullsDefault`。列映射应由服务端静态配置提供；普通筛选值仍必须使用参数化查询。
