# HTTP 响应提交守卫

并发 fan-out 或异步回调中，成功、超时和错误路径可能竞争写入同一个 `http.ResponseWriter`。
一旦响应头提交，后续写入无法更改状态码，还可能把不同响应体拼接在一起。

```go
guard := New(w)
go func() {
	if err := callDependency(ctx); err != nil {
		guard.Commit(http.StatusBadGateway, []byte("dependency failed"))
		return
	}
	guard.Commit(http.StatusOK, []byte("ok"))
}()
```

`Guard` 用互斥锁保证只有首个 `Commit` 生效，返回值让失败的竞争者停止后续动作。实际 handler
仍应等待 goroutine 退出，并在提交前设置 Content-Type 等响应头；该结构不负责取消后台任务。
