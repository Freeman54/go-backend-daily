# 按比例分配 Context 时间预算

多阶段请求若把父请求的完整截止时间传给每一层，前置阶段可能耗尽全部时间。`WithFraction` 根据父上下文的剩余时间，为当前阶段分配固定比例，同时支持最小预算且绝不越过父截止时间。

```go
stageCtx, cancel, err := WithFraction(ctx, time.Now(), 0.3, 50*time.Millisecond)
if err != nil {
	return err
}
defer cancel()
return queryDatabase(stageCtx)
```

调用者必须执行 `cancel` 及时释放定时器。示例显式传入 `now`，让时间计算可以稳定测试。真实链路还应为序列阶段预留收尾时间，并在父上下文没有截止时间时明确决定策略，而不是悄悄制造无限等待。
