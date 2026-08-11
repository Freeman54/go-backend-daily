# HTTP 方法策略

路由层散落的 method 判断容易产生不一致的 `405 Method Not Allowed` 响应。`Policy` 在启动时规范化并去重方法集合，被拒绝时稳定返回排序后的 `Allow` 响应头。

```go
policy, _ := methodpolicy.New("GET", "POST")
handler := policy.Middleware(apiHandler)
http.ListenAndServe(":8080", handler)
```

策略适合放在鉴权和业务处理之前。真实服务还应把 `OPTIONS`、CORS 和自动生成的路由元数据统一考虑，避免网关与应用声明不同。
