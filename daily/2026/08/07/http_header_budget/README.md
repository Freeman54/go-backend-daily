# HTTP 请求头预算

网关转发请求时，请求头也是不受信任的输入：超大的 `Cookie`、追踪字段或大量自定义头都会放大代理、日志和下游服务的内存成本。`Validate` 同时限制字段数和字段名/值的总字节数，超过预算立即返回可用 `errors.Is` 判断的错误。

```go
headers := http.Header{"X-Request-ID": {"req-42"}}
if err := httpheaderbudget.Validate(headers, 16, 4<<10); err != nil {
	// 映射为 431 Request Header Fields Too Large。
}
```

计数不包含 HTTP 的分隔符，因而策略不会因序列化实现不同而变化。预算应由入口网关统一配置，并与反向代理的最大头部限制保持一致。
