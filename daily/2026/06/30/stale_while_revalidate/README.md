# 用 stale-while-revalidate 保持缓存可用

缓存过期时，如果所有请求都同步回源，容易把数据库或下游接口打爆。`stale-while-revalidate` 的思路是先返回短时间内允许的旧值，同时只触发一次刷新，让用户体验和后端稳定性都更平衡。

## 设计要点

- `ExpiresAt` 之前返回新值
- 进入 `StaleUntil` 窗口后返回旧值，但只放行一次刷新
- 刷新失败后清理刷新标记，让后续请求能再次尝试回源

## 示例说明

真实业务里通常把 `Refresh` 放到 goroutine 或后台任务里异步执行，这里为了便于测试，显式暴露成一个方法，演示“读旧值”和“单次刷新”的核心机制。

## 运行方式

```bash
go test ./daily/2026/06/30/stale_while_revalidate
```
