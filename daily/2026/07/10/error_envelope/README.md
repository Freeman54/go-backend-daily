# 用 error envelope 统一后端错误输出

后端错误链条通常很长：存储错误、领域错误、HTTP 错误会一层层包起来。如果直接把 `err.Error()` 暴露给客户端，既不稳定也不安全；如果完全丢掉底层原因，排障又会很慢。更合适的做法是构建一个对外稳定、对内可追踪的错误 envelope。

## 设计要点

- 顶层错误决定客户端可见的 `code` 与 `message`
- 底层 `cause` 保留在 envelope 中，便于日志、审计和告警关联
- 对重复错误文案做去重，避免多层 wrap 后的噪音

## 示例说明

示例里的 `CodedError` 让服务层能够在保留 `Unwrap` 链的同时给错误打上业务 code。`Build` 最终生成统一错误 envelope，测试覆盖显式业务错误、普通内部错误和重复 cause 去重。

## 运行方式

```bash
go test ./daily/2026/07/10/error_envelope
```
