# 用 `context.AfterFunc` 组织取消后的清理动作

请求被取消时，最怕资源只清了一半：span 没结束、连接没归还、临时文件没删除。Go 1.21 之后可以用 `context.AfterFunc` 把“取消后要做的事”挂到上下文上，避免每层业务都手动写重复的清理分支。

## 设计要点

- `Group.Add` 收集多个清理动作，适合把连接、trace、锁释放集中管理
- `Bind` 用 `context.AfterFunc` 绑定取消事件，拿到的 `stop` 可以在正常完成时主动解除回调
- `Cleanup` 保证只执行一次，并按后进先出顺序释放资源，减少依赖顺序出错

## 示例说明

示例里把连接归还和 span 结束放进同一个 `Group`。请求取消后会自动触发清理；如果请求正常完成，可以先调用 `stop()` 解除回调，再按显式路径收尾。

## 运行方式

```bash
go test ./daily/2026/07/14/context_afterfunc_cleanup
```
