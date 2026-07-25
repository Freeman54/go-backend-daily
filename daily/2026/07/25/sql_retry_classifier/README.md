# SQL 重试分类

数据库错误不能一律重试。序列化失败（`40001`）、死锁（`40P01`）和锁不可用（`55P03`）通常是瞬态错误；
唯一约束冲突等确定性错误继续重试只会增加负载。

```go
if err := tx.Commit(); err != nil {
    if IsRetryable(err) {
        return retryWithBackoff(ctx)
    }
    return err
}
```

示例依赖驱动错误实现 `SQLState() string`，并通过 `errors.As` 支持包装错误。重试必须设置次数、总时长、
指数退避和随机抖动；事务函数还必须具备幂等性，避免已生效的外部副作用被重复执行。
