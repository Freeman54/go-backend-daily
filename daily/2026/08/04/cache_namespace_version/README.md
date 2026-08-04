# 缓存命名空间版本

批量删除某类缓存通常昂贵且容易漏删。`Namespace` 把版本号写入 key；配置或数据模型发生整体变化时只需提升版本，旧 key 会自然过期，不再进入读路径。

```go
ns := cache_namespace_version.New("user-profile")
key := ns.Key("42") // user-profile:v1:42
ns.Bump()
key = ns.Key("42")  // user-profile:v2:42
```

示例用原子变量演示进程内并发安全。多实例服务应把版本保存在 Redis 或配置中心，并通过原子自增和订阅机制同步；旧版本数据仍需设置 TTL 控制存储占用。
