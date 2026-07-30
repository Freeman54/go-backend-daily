# SQL 保存点规划器

长事务中的批量操作可以用保存点只回滚局部步骤。保存点名称通常不能作为参数绑定，因此必须严格验证标识符，不能直接拼接外部输入。

```go
p, err := sqlsavepointplan.New("import_batch_1")
if err != nil {
	return err
}
_, _ = tx.ExecContext(ctx, p.Begin())
// 执行可恢复步骤；失败时执行 p.Rollback()
_, _ = tx.ExecContext(ctx, p.Release())
```

示例只接受 ASCII 字母、数字和下划线，且限制为 63 字符，再使用双引号引用。不同数据库对保存点语法和事务错误状态的处理不同，接入前应核对驱动文档；保存点也不能缩短外层事务持锁时间。
