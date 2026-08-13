# HTTP 内容协商

服务端不应只用字符串包含判断解析 `Accept`，因为客户端可能携带权重、通配符和多个候选类型。`Choose` 使用标准库解析媒体类型，忽略 `q=0` 的禁用项，并在权重相同时保持服务端偏好顺序。

```go
mediaType, ok := contentnegotiator.Choose(r.Header.Get("Accept"), []string{"application/json", "text/plain"})
if !ok { http.Error(w, "not acceptable", http.StatusNotAcceptable); return }
w.Header().Set("Content-Type", mediaType)
```

生产环境还应设置 `Vary: Accept`，避免共享缓存把一种表示错误复用于其他客户端。
