# 构造稳定的 HTTP Vary 缓存键

反向代理或应用内缓存不能只按 URL 命中：语言、编码等请求头也可能改变响应。`Build` 将方法、主机、路径、规范化查询参数和指定的 Vary 请求头组合成确定性缓存键。

```go
req, _ := http.NewRequest("GET", "https://api.example.com/items?b=2&a=1", nil)
req.Header.Set("Accept-Language", "zh-CN")
key, err := http_vary_cache_key.Build(req, []string{"Accept-Language"})
```

查询参数统一排序，头名称统一小写，并拒绝包含换行符的非法头名称，避免缓存键注入。生产环境可再对最终字符串做哈希，并限制参与 Vary 的头字段集合，防止缓存基数失控。
