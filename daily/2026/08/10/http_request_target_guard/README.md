# HTTP 请求目标守卫

反向代理、签名校验和缓存键生成都依赖一致的请求目标。`Normalize` 拒绝绝对 URL、反斜杠、片段和路径穿越，并对查询参数排序，避免同一请求因表示不同而绕过策略或污染缓存。

```go
target, err := requesttargetguard.Normalize("/v1/orders?status=paid&page=2", 2048)
if err != nil {
	// 返回 400 Bad Request。
}
_ = target // /v1/orders?page=2&status=paid
```

守卫只接受 origin-form（以 `/` 开头的路径）。生产环境还应让网关和应用使用相同的长度限制，并在签名、鉴权和缓存查找之前完成规范化。
