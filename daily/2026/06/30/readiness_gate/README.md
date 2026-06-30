# 用 readiness gate 控制实例摘流与上线

后端服务不仅在进程启动成功时才算可用，还要看数据库、缓存、消息队列等依赖是否就绪。发布摘流时，也需要把实例主动标记为不可接流。这个示例用一个简单的 gate 汇总依赖状态和摘流状态。

## 设计要点

- 每个依赖单独上报 `Ready` 状态
- 摘流时通过 `SetDraining(true)` 直接阻止新流量进入
- `Snapshot` 返回布尔结果和阻塞原因，便于健康检查接口直接输出

## 示例说明

真实服务里可以把 `Snapshot` 结果映射到 `/readyz` 接口。Kubernetes 或服务注册中心只看 `Ready`，排障页面再展示 `Reasons` 明细。

## 运行方式

```bash
go test ./daily/2026/06/30/readiness_gate
```
