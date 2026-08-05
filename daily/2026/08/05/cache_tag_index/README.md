# 缓存标签索引

当一个实体变化会影响多个缓存 key 时，逐个硬编码删除容易遗漏。`Index` 维护标签到 key 的反向索引，`Invalidate` 原子取出并清空某标签关联的 key 集合。

```go
index := cache_tag_index.New()
index.Attach("post:7", "author:2", "feed:hot")
keys := index.Invalidate("author:2")
```

示例聚焦并发安全的进程内索引。分布式部署应使用 Redis Set 与 Lua 脚本，并为索引设置生命周期，防止已过期 key 留在标签集合中。
