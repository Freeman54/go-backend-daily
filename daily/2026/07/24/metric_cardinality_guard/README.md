# 指标标签基数守卫：限制时序数量膨胀

把用户 ID、URL 或原始错误文本直接作为指标标签，会持续创建新时间序列，增加内存和查询成本。`Guard` 允许有限数量的新标签值；超过上限后统一映射为低基数兜底值，同时继续识别已经登记的值。

```go
guard, err := metriccardinalityguard.New(20, "other")
if err != nil {
	log.Fatal(err)
}
methodLabel := guard.Normalize(operation)
requests.WithLabelValues(methodLabel).Inc()
```

更稳妥的做法是预先定义有限枚举，守卫适合保护暂时无法完全控制的输入。生产环境还应记录降级次数，并按指标名称和标签键分别设置预算。

运行：

```bash
go test ./daily/2026/07/24/metric_cardinality_guard
```
