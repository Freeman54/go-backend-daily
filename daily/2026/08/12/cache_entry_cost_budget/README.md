# 缓存条目成本准入

只限制缓存条目数量会让少量大对象挤占全部内存。`Budget` 用调用方定义的成本（通常是序列化字节数，也可以加入重建耗时权重）原子地执行准入与释放。

```go
budget := entrycostbudget.New(64 << 20)
if budget.Reserve(int64(len(encoded))) {
    cache.Set(key, encoded)
}
// 淘汰时调用 budget.Release(cost)
```

真实缓存需要把成本和条目一起保存，确保覆盖、淘汰和删除时只释放一次。该预算解决容量安全，不替代 LRU/LFU 等淘汰策略。
