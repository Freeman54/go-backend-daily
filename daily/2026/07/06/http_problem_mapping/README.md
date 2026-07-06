# 用 HTTP Problem Details 稳定输出接口错误

后端接口如果直接把内部错误字符串暴露给客户端，既不稳定，也不利于前端和调用方按类型处理。更好的方式是把领域错误、校验错误和冲突错误映射成稳定的 Problem Details 响应结构，让状态码、标题和错误类型都具备可预期语义。

## 设计要点

- 在 handler 边界统一做错误映射，不把业务层绑定到 HTTP 细节
- 优先识别 `ValidationError` 和包装过的哨兵错误，避免状态码漂移
- 对未知错误返回稳定的 500 响应，隐藏内部实现细节

## 示例说明

示例实现了 `Map` 函数，把 `ValidationError`、`ErrUnauthorized`、`ErrNotFound` 和 `ErrConflict` 映射到不同的 Problem Details；未知错误统一回退到 500。测试覆盖了校验错误、包装错误和未知错误三类场景。

## 运行方式

```bash
go test ./daily/2026/07/06/http_problem_mapping
```
