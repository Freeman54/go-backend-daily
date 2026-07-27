# 并发结果顺序缓冲

并行处理能提高吞吐，但任务完成顺序通常不稳定。若下游协议要求严格递增的序号，可以让 worker 把
`sequence + result` 交给 `Buffer`：缺口之前的结果暂存，缺口补齐后一次释放连续结果。

```go
buffer := New[string](100)
buffer.Add(101, "second")          // 暂存，无输出
ready := buffer.Add(100, "first") // ["first", "second"]
for _, result := range ready {
	_ = result // 按序写入下游
}
```

实现会忽略已消费序号和重复序号，并用互斥锁保护状态，可由多个 worker 调用。真实服务还应限制
`pending` 大小并设置缺口超时，否则一个永久缺失的结果会导致内存持续增长。
