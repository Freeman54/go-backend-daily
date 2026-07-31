# 有序停机任务组

后端进程停机时存在依赖顺序：应先停止接收消息，再等待业务请求，最后刷新指标和关闭连接。`Group` 为钩子标记阶段，停机时按阶段从高到低执行，同一阶段保持注册顺序，并用 `errors.Join` 汇总失败。

```go
g := orderedshutdowngroup.New()
g.Add(20, "consumer", stopConsumer)
g.Add(10, "http", shutdownHTTP)
g.Add(0, "metrics", flushMetrics)
if err := g.Shutdown(ctx); err != nil {
	log.Printf("shutdown: %v", err)
}
```

每个钩子都应尊重传入的 `context` 并保证幂等。示例串行执行以表达确定的依赖顺序；独立资源可注册到一个钩子内并发关闭，但仍要受统一的停机截止时间约束。
