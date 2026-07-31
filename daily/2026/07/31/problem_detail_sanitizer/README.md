# Problem Details 输出净化

错误对象常混有 SQL、内部主机名、堆栈或凭证，直接序列化给客户端会泄露实现细节。`Sanitize` 只保留稳定的 `type`、`title` 和合法 HTTP 状态码；`detail` 只提取命中的公开白名单短语，实例路径默认移除。

```go
public := problemdetailsanitizer.Sanitize(problem, []string{
	"order has already been paid",
	"requested quantity is unavailable",
})
_ = json.NewEncoder(w).Encode(public)
```

白名单应使用固定业务文案，不要包含原始错误字符串。完整错误、请求 ID 和堆栈应写入受控日志，客户端通过不含敏感信息的关联 ID 发起排查；净化逻辑应位于统一错误响应边界。
