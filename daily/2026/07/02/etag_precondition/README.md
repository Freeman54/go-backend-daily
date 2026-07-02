# 用 ETag 条件更新保护接口并发写

开放编辑接口时，两个客户端先后读取同一份资源、再分别提交更新，是最常见的并发写覆盖来源。这个示例用 `ETag` 和 `If-Match` 做条件更新，只允许基于最新版本的修改落库。

## 设计要点

- 读取资源时返回版本号派生的 `ETag`
- 更新时强制客户端回传 `If-Match`
- 当前版本与 `If-Match` 不一致时返回前置条件失败，要求客户端先刷新

## 示例说明

HTTP 场景里常见做法是把 `ErrPreconditionFailed` 映射为 `412 Precondition Failed`。这样比“最后写入覆盖前者”更安全，也能让前端显式处理编辑冲突。

## 运行方式

```bash
go test ./daily/2026/07/02/etag_precondition
```
