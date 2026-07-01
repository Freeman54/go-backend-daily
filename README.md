# go-backend-daily
用 Go 记录后端研发的每日练习与技术思考，涵盖并发、工程设计、数据库、缓存、消息队列、性能优化、可观测性与架构实践。每篇内容包含可运行代码、测试用例和中文说明，作为持续学习、复盘和技术分享的沉淀。

## 目录

- [2026-06-25 用 context 控制并发任务的超时与取消](daily/2026/06/25/context_timeout)
- [2026-06-25 后端错误分类与 HTTP 状态码映射](daily/2026/06/25/error_taxonomy)
- [2026-06-26 用有界 worker pool 控制并发任务](daily/2026/06/26/worker_pool)
- [2026-06-26 用幂等键缓存重复请求响应](daily/2026/06/26/idempotency_store)
- [2026-06-26 用白名单构建局部更新 SQL](daily/2026/06/26/sql_patch_builder)
- [2026-06-26 用 singleflight 思路避免缓存击穿](daily/2026/06/26/cache_singleflight)
- [2026-06-26 用断路器隔离不稳定下游](daily/2026/06/26/circuit_breaker)
- [2026-06-29 用信号量做接口并发准入控制](daily/2026/06/29/semaphore_admission)
- [2026-06-29 用可重试错误和退避策略保护下游](daily/2026/06/29/retry_backoff)
- [2026-06-29 用游标分页稳定返回大列表](daily/2026/06/29/cursor_pagination)
- [2026-06-29 用事务 outbox 保证事件最终投递](daily/2026/06/29/transaction_outbox)
- [2026-06-29 用优雅停机收敛后台任务](daily/2026/06/29/graceful_shutdown)
- [2026-06-30 用令牌桶限制突发流量](daily/2026/06/30/token_bucket)
- [2026-06-30 用超时预算拆分下游调用时间](daily/2026/06/30/timeout_budget)
- [2026-06-30 用 stale-while-revalidate 保持缓存可用](daily/2026/06/30/stale_while_revalidate)
- [2026-06-30 用 readiness gate 控制实例摘流与上线](daily/2026/06/30/readiness_gate)
- [2026-06-30 用事务重试兜住串行化冲突](daily/2026/06/30/tx_retry)
- [2026-07-01 用 quorum 读取降低副本抖动影响](daily/2026/07/01/fanout_quorum)
- [2026-07-01 用 cancel cause 向并发任务传播失败原因](daily/2026/07/01/cancel_cause_group)
- [2026-07-01 用乐观锁避免并发写覆盖](daily/2026/07/01/optimistic_lock)
- [2026-07-01 用重试与死信路由保护消息消费](daily/2026/07/01/dlq_router)
- [2026-07-01 用滑动窗口统计延迟 SLO](daily/2026/07/01/latency_window)
