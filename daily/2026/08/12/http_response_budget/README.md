# HTTP 响应体预算

后端接口可能因错误分页、异常聚合或序列化膨胀返回巨大响应。`Writer` 在字节写入前检查累计大小，超限时返回 `ErrResponseTooLarge`，避免继续消耗带宽。

```go
limited := &responsebudget.Writer{ResponseWriter: w, Limit: 1 << 20}
if err := json.NewEncoder(limited).Encode(result); err != nil {
    // 记录超限指标；若响应尚未提交，可返回受控错误。
}
```

流式写入一旦已经提交响应头便无法改写状态码，因此生产环境更适合先写入有上限的缓冲区，再统一提交。预算还应与网关限制、压缩策略和分页上限保持一致。
