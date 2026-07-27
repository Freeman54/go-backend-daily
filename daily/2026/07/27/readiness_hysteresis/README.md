# 就绪探针迟滞门

依赖偶发超时会让实例的就绪状态在成功与失败之间快速切换，造成流量频繁迁移。`Gate` 使用连续失败阈值
关闭门，再用独立的连续成功阈值恢复门，形成迟滞，过滤短暂毛刺。

```go
gate, _ := New(3, 5)
ready := gate.Observe(checkDatabase() == nil)
if !ready {
	http.Error(w, "not ready", http.StatusServiceUnavailable)
}
```

任一成功会清空失败连续计数，任一失败也会清空恢复连续计数。阈值应结合探测周期和可接受摘流时间设置；
就绪探针适合控制是否接收新流量，进程存活探针不应因普通依赖故障而重启服务。
