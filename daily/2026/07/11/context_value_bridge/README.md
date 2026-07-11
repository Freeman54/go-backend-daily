# 用 context value bridge 安全启动后台任务

请求结束后再异步补日志、发审计事件、写埋点是很常见的需求，但直接把原始 `context.Context` 传给后台 goroutine 往往会马上被取消。另一方面，如果完全丢弃原 context，又会把 `request_id`、`trace_id` 这类排障关键字段一起丢掉。更稳妥的方式是只桥接需要的 value，同时给后台任务单独设置一个更短的生命周期。

## 设计要点

- 只复制显式声明的 context key，避免把大对象和无关状态带进后台任务
- 新 context 基于 `context.Background()` 创建，不继承上游取消信号
- 为后台任务设置独立 timeout，防止清理逻辑无限悬挂

## 示例说明

`Detach` 会从父 context 中摘取指定 key 的 value，再构造一个全新的 background context。测试覆盖了 value 精确复制、父请求取消后子任务继续执行，以及新 timeout 生效三类场景。

## 运行方式

```bash
go test ./daily/2026/07/11/context_value_bridge
```
