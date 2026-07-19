# 用 trace log enricher 统一日志里的链路标签

链路追踪、业务日志和告警聚合往往分布在不同组件里。要是每个 handler、DAO、consumer 都手写一遍 `trace_id`、`span_id`、`tenant`，代码会快速重复，而且很容易漏字段或字段名不一致。

`trace log enricher` 的做法是把日志需要的上下文字段统一放进 `context`，再在写日志前提取成结构化属性。这样日志字段来源一致，查询也更稳定。

## 设计要点

- 用白名单约束可出现在日志中的上下文字段
- 统一字段名，减少不同模块各写各的情况
- 缺失字段时自动跳过，不污染日志

## 示例说明

示例实现提供：

- `WithTraceID`、`WithSpanID`、`WithTenant`、`WithRequestID`
- `AttrsFromContext(ctx)`：提取结构化 `slog.Attr`

适合 HTTP 中间件、异步消费者、数据库访问日志和审计日志场景。

## 运行方式

```bash
go test ./daily/2026/07/19/trace_log_enricher
```
