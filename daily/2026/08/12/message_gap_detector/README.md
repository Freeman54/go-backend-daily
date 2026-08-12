# 消息序列缺口检测

有序消费不能把“收到更大的 offset”误当成进度推进，否则缺失消息会被永久跳过。`Detector` 区分连续消息、重复投递和缺口，并在遇到缺口时保持期望位置不变。

```go
detector := gapdetector.New(checkpoint + 1)
result := detector.Observe(message.Offset)
if result.Kind == gapdetector.Gap {
    requestReplay(result.Expected, result.Offset-1)
}
```

生产环境通常按 topic/partition 各维护一个检测器，并把 checkpoint 持久化。它只检测顺序，不负责消息缓存、重放或跨分区全局排序。
