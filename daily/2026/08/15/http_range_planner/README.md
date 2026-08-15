# 用 Range planner 安全处理断点下载

文件下载接口不能把 `Range` 头直接切字符串后用于读取：越界、后缀范围、多段范围和空资源都需要明确语义。本例实现单段 `bytes` 范围规划器，支持 `2-5`、`7-` 和 `-3`，并区分格式错误与不可满足范围。

处理成功时可返回 `206 Partial Content`，使用 `Start`、`Length()` 读取文件，并通过 `ContentRange` 生成响应头。`ErrRangeNotSatisfy` 应映射为 `416`，多段范围则明确拒绝，避免在尚未实现 multipart 响应时悄悄返回错误内容。

```go
r, err := httprangeplanner.Plan("bytes=100-199", fileSize)
if err == nil {
    // reader.ReadAt(buf[:r.Length()], r.Start)
    contentRange := httprangeplanner.ContentRange(r, fileSize)
    _ = contentRange
}
```

运行测试：

```bash
go test ./daily/2026/08/15/http_range_planner
```
