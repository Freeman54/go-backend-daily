# 用复合游标实现稳定的 SQL 翻页

偏移量分页在深页会越来越慢，而且并发插入会造成重复或遗漏。游标分页可用 `(created_at, id)` 作为稳定排序键，下一页查询写成 `WHERE (created_at, id) < (?, ?) ORDER BY created_at DESC, id DESC LIMIT ?`。

示例把时间和 ID 编码成 URL 安全的不透明 token，并在解码时严格校验字段。真实 API 若担心客户端篡改，可在 payload 外增加 HMAC；排序字段和比较方向必须与 SQL 完全一致。

```bash
go test ./daily/2026/07/22/sql_cursor_codec
```
