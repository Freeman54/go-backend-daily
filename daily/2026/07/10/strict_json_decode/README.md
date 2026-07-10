# 用 strict JSON decode 收紧 Go HTTP 接口输入边界

后端接口最怕“客户端悄悄多传字段、服务端悄悄忽略”。这类问题会把拼写错误、版本漂移和灰度兼容风险都藏起来。更稳妥的做法是在入口直接拒绝未知字段、空 body 和尾随多段 JSON。

## 设计要点

- `json.Decoder.DisallowUnknownFields()` 让字段拼写错误尽早暴露
- 二次 `Decode` 检查是否存在尾随 JSON，避免多个对象拼接绕过校验
- 结构化解码后再补业务字段校验，区分协议错误和业务错误

## 示例说明

示例里的 `DecodeRequest` 模拟 HTTP handler 对请求体做严格解码。测试覆盖正常输入、未知字段、尾随 payload，以及业务字段不合法等情况。

## 运行方式

```bash
go test ./daily/2026/07/10/strict_json_decode
```
