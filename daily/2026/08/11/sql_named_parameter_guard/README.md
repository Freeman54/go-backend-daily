# SQL 命名参数守卫

动态拼装查询时，遗漏参数或多传无效参数往往直到运行时才暴露。`Validate` 提取 `:name` 占位符，并要求它与参数映射的键完全一致；字符串字面量中的冒号会被忽略。

```go
args := map[string]any{"tenant": "acme", "id": 42}
if err := namedparameter.Validate(query, args); err != nil {
	return err
}
```

守卫只校验绑定契约，不负责拼接 SQL，也不替代驱动的参数化查询。示例解析器面向常见 SQL；生产环境应针对所用数据库补齐转义字符串、注释和类型转换语法。
