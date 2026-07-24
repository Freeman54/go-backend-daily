# 事件去重窗口：抑制消息系统的短期重复投递

消息队列通常提供“至少一次”投递，消费者必须预期同一事件重复到达。`Window` 记录事件 ID 的过期时间：首次到达返回 `true`，窗口内重复到达返回 `false`，到期后可以再次接受。

```go
window, err := eventdedupwindow.New(5 * time.Minute)
if err != nil {
	log.Fatal(err)
}
if window.Accept(message.ID, time.Now()) {
	handle(message)
}
```

这是单进程教学实现，进程重启会丢失状态，也不会跨实例共享。需要强一致副作用时，应使用数据库唯一键、幂等表或 Redis 原子操作，并让去重记录与业务提交保持一致。

运行：

```bash
go test ./daily/2026/07/24/event_dedup_window
```
