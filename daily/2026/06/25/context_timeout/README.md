# 用 context 控制并发任务的超时与取消

在后端服务里，一个请求往往会拆成多段下游调用：查数据库、访问缓存、调用 RPC、写消息队列。只要其中一段已经失败，继续执行剩余任务通常只会浪费资源，甚至放大故障。

Go 的 `context.Context` 很适合表达这种控制关系：上游决定生命周期，下游持续监听取消信号，任务编排层负责在错误或超时时尽快收敛。

## 核心思路

这个示例实现了一个小型 worker runner：

- 父级 `context` 超时后，所有 worker 都应该尽快停止。
- 任意任务返回错误后，runner 会主动取消其它任务。
- 每个任务都接收同一个派生 `context`，任务内部需要主动监听 `ctx.Done()`。
- runner 返回完成数量、失败数量、是否取消以及首个错误。

## 运行方式

在仓库根目录执行：

```bash
go test ./...
```

也可以只跑今天的示例：

```bash
go test ./daily/2026/06/25/context_timeout
```

## 关键代码

任务签名很简单：

```go
type Job func(context.Context) error
```

runner 会从父级 context 派生一个可取消的 context：

```go
ctx, cancel := context.WithCancel(ctx)
defer cancel()
```

当某个任务失败时，worker 记录错误并调用 `cancel()`，其它正在执行的任务会通过 `ctx.Done()` 感知取消：

```go
if err := job(ctx); err != nil {
    errCh <- err
    cancel()
}
```

## 后端实践建议

1. 不要在业务函数里凭空创建 `context.Background()`，优先接收上游传入的 `ctx`。
2. 数据库、HTTP、RPC、消息队列客户端调用都应该使用支持 context 的 API。
3. 任务内部如果包含循环、重试或等待，一定要监听 `ctx.Done()`，否则取消信号无法及时生效。
4. 错误触发取消时，只记录首个关键错误，避免并发错误把真正的根因淹没。
5. 超时应该由入口层或编排层统一设置，下游只负责遵守，不要各自随意设置互相冲突的超时。

## 小结

`context` 不是简单的参数传递工具，而是后端系统里的生命周期控制协议。把取消和超时放进并发编排中，可以让服务在失败时更快止损，在压力下更容易保持稳定。
