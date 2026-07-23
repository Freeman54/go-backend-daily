# 有界并发 Map：控制扇出并保持结果顺序

批量调用下游接口时，为每条数据直接启动 goroutine 会让瞬时并发随输入规模增长，可能耗尽连接池或压垮依赖。`Map` 只维持固定数量的在途任务，并按输入下标写回结果，因此执行顺序可以不同，返回顺序仍稳定。

实现采用“完成一个再补一个”的调度方式。任一任务失败后立即取消派生 context，不再调度新任务；已启动的函数应监听收到的 context，及时释放资源。生产环境还应结合超时、指标和下游并发预算。

```go
values, err := boundedparallelmap.Map(ctx, ids, 8, loadUser)
```

运行示例：

```bash
go test ./daily/2026/07/23/bounded_parallel_map
```
