# 用任务追踪器实现可靠的优雅下线

服务收到终止信号后，既要拒绝新任务，也要等待已接收任务完成。只依赖固定 `sleep` 无法适应任务时长，而裸 `WaitGroup` 又容易在 `Wait` 与 `Add` 并发时产生错误用法。

`Tracker` 用 `Begin` 登记任务并返回幂等完成函数；`Drain` 原子地关闭准入，再等待活动任务归零或 context 超时。HTTP 服务可先从负载均衡摘流，再调用 `Drain`，最后关闭数据库和消息客户端等共享资源。

```bash
go test ./daily/2026/07/22/graceful_drain_tracker
```
