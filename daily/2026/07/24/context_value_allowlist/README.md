# Context 值白名单复制：跨生命周期传递最少元数据

后台任务有时不能继承请求的取消信号，却仍需要 trace ID、租户 ID 等元数据。直接把整个请求 `context` 保存到异步任务会错误继承 deadline，也可能泄漏认证信息。`Copy` 以新 context 为生命周期，只复制显式白名单中的值。

```go
background := contextvalueallowlist.Copy(
	context.Background(),
	requestContext,
	[]any{traceKey, tenantKey},
)
```

白名单应只包含小型、不可变、请求范围的数据。不要把数据库连接、大对象或权限凭证塞入 context；业务必需参数更适合显式函数参数。本例不复制取消和截止时间，这是有意的生命周期切割。

运行：

```bash
go test ./daily/2026/07/24/context_value_allowlist
```
