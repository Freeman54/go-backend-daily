# 用 quiescence gate 实现静默关闭

优雅停机需要两个阶段：先拒绝新任务，再等待在途任务结束。单独使用 `WaitGroup` 很难安全地协调“停止 Add”和“开始 Wait”；本例把进入、离开和关闭状态放在同一把锁下，消除 `Add` 与 `Wait` 的竞态窗口。

`Enter` 成功后必须 `defer leave()`；`Close` 可重复调用，并返回一个只在所有在途任务完成后关闭的 channel。服务关闭时可同时等待该 channel 和总停机 deadline。

```go
leave, err := gate.Enter()
if err != nil { /* 返回 503 */ }
defer leave()

select {
case <-gate.Close():
case <-shutdownCtx.Done():
}
```

运行测试：

```bash
go test -race ./daily/2026/08/15/concurrency_quiescence_gate
```
