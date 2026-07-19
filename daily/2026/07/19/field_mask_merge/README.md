# 用 field mask merge 实现安全的 PATCH 更新

很多 HTTP PATCH 或 RPC Update 接口都要同时处理三件事：客户端想改哪些字段、哪些字段允许改、字段缺失和字段清空该怎么区分。直接把整个请求结构体覆盖到数据库模型上，通常会把“没传”误判成“要置空”。

`field mask merge` 通过显式字段列表驱动合并，只更新 mask 指定的字段，并允许把字段设置为空值或删除。

## 设计要点

- `mask` 明确声明客户端要修改的字段
- 只允许更新白名单中的字段，避免越权写入
- 未出现在 `mask` 中的字段保持原值不变

## 示例说明

示例实现提供 `Apply(base, updates, mask, allowed)`：

- `base` 是当前持久化状态
- `updates` 是待更新值，`nil` 表示清空该字段
- `mask` 指定本次真正生效的字段

适合用户资料修改、配置更新、后台管理接口等场景。

## 运行方式

```bash
go test ./daily/2026/07/19/field_mask_merge
```
