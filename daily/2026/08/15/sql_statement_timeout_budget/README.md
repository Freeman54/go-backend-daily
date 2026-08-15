# 用请求预算设置事务级 statement_timeout

只有 Go `context` 超时并不足以保护数据库：客户端取消到达数据库前可能有延迟，连接异常时语句也可能继续消耗资源。本例从请求 deadline 中扣除提交、回滚和响应编码的预留时间，再与服务端允许的单语句上限取较小值。

`SetLocalSQL` 使用 PostgreSQL `set_config(..., true)`，让 `statement_timeout` 仅在当前事务生效，避免连接归还池后污染下一个请求。设置语句仍采用参数绑定，不拼接外部输入。

```go
budget, err := sqlstatementtimeoutbudget.Budget(ctx, time.Now(), 3*time.Second, 200*time.Millisecond)
if err == nil {
    query, args, _ := sqlstatementtimeoutbudget.SetLocalSQL(budget)
    _, err = tx.ExecContext(ctx, query, args...)
}
```

数据库超时应被归类为可观测的 deadline 错误；是否重试还要结合事务幂等性和剩余总预算判断。

运行测试：

```bash
go test ./daily/2026/08/15/sql_statement_timeout_budget
```
