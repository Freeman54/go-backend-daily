# 用 span event limiter 控制热点路径的观测噪音

排查线上问题时，大家都希望 trace 和日志足够详细；但在重试循环、批量消费或大规模 fanout 里，如果每一次失败都往 span 里塞事件，很快就会把单条 trace 撑爆，既增加存储成本，也让真正关键信息被噪音淹没。

## 设计要点

- 每类事件只保留前 N 次，避免同一错误在热点循环里重复刷屏
- 超过限制的事件不再逐条输出，而是累计成 dropped 计数
- 在 span 结束前统一输出摘要，兼顾可读性和观测成本

## 示例说明

示例实现了一个事件限流器，适合包在 tracing instrumentation 外层。调用方可以先问 `Allow(name)`，允许时再真正写入 span；结束前再把 `FlushSummaries()` 的结果作为摘要事件补进去。测试验证：

- 超过上限的重复事件会被丢弃
- `FlushSummaries` 会返回丢弃摘要，并重置内部状态

## 运行方式

```bash
go test ./daily/2026/07/15/span_event_limiter
```
