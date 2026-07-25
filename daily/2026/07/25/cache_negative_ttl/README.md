# 负缓存 TTL

当不存在的 key 被频繁查询时，只缓存命中结果会让所有请求穿透到数据库。负缓存把“不存在”也短暂保存，
但使用比正常值更短的 TTL，减少新数据写入后仍返回旧空结果的时间。

```go
cache := New[User](time.Minute, 5*time.Second, time.Now)
cache.Set("user:42", User{}, false)
value, exists, found := cache.Get("user:42")
// found=true 表示缓存命中；exists=false 表示源数据不存在。
```

实现用三态返回值区分“缓存未命中”和“已缓存不存在”。生产环境还应在写入成功后主动删除负缓存，并为
TTL 加随机抖动，避免大量 key 同时过期。
