# 用分片 map 缓解热点 key 的锁竞争

计数器、限流器或本地缓存一旦承载高频热点 key，很容易因为一把全局锁把所有 goroutine 串起来。把数据按 key 哈希到多个 shard 上，是 Go 后端里很常见也很实用的第一层降争用手段。

## 设计要点

- 用稳定哈希把 key 分散到多个 shard，每个 shard 持有自己的锁和 map
- 热点 key 依然会落到单个 shard，但不同 key 之间的争用会被明显削弱
- shard 数量不是越多越好，需要结合内存占用和典型并发度权衡

## 示例说明

示例实现了一个分片计数器，提供 `Add`、`Get` 和 `ShardLoads` 三个接口。测试覆盖了并发累加正确性，以及不同 key 会被分散到多个 shard 的基本性质。

## 运行方式

```bash
go test ./daily/2026/07/05/hot_key_sharding
```
