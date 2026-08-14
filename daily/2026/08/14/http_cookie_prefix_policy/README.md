# HTTP Cookie 前缀策略

浏览器为 `__Secure-` 和 `__Host-` Cookie 定义了更严格的安全约束。`Validate` 在服务端写响应前检查这些约束，避免配置错误导致 Cookie 被浏览器拒绝，或意外扩大 Cookie 的作用域。

```go
cookie := &http.Cookie{Name: "__Host-session", Secure: true, Path: "/"}
if err := cookieprefix.Validate(cookie); err == nil {
    http.SetCookie(w, cookie)
}
```

`__Secure-` 要求 `Secure`；`__Host-` 还要求 `Path=/` 且不能设置 `Domain`。前缀区分大小写，普通 Cookie 不受这组额外规则影响。
