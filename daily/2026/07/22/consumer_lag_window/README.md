# 用滑动窗口观察消息消费者积压

单次采样的 consumer lag 容易因分区抖动产生误报。固定大小滑动窗口只保留最近 N 次采样，并同时输出最大值和平均值，可以区分瞬时尖峰与持续积压。

示例使用环形数组实现常量空间的 `Window`。`Add` 写入最新 lag，`Snapshot` 汇总当前有效样本。生产环境可按 topic/consumer group 分组维护窗口，并把最大值用于快速告警、平均值用于趋势判断；并发调用时需在外层串行化或增加互斥锁。

```bash
go test ./daily/2026/07/22/consumer_lag_window
```
