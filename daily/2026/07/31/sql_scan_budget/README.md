# SQL 结果集扫描预算

即使查询本身很快，无界结果集也可能在 Go 进程中消耗大量内存、CPU 和响应带宽。`Budget` 在逐行扫描时同时累计行数和估算字节数，超出任一上限便拒绝接收下一行，且被拒绝的行不计入用量。

```go
budget, _ := sqlscanbudget.New(1_000, 2<<20)
for rows.Next() {
	var payload []byte
	if err := rows.Scan(&payload); err != nil {
		return err
	}
	if err := budget.Consume(len(payload)); err != nil {
		return err
	}
}
```

真实代码应在超限时及时关闭 `Rows`，并优先用分页、`LIMIT` 和数据库侧成本控制减少数据传输。这里的字节数是应用层估算，不等于数据库页读取量，也不能替代查询超时和慢查询观测。
