# 安全的 SQL 动态过滤器

动态搜索接口常需要按若干字段添加条件。`Build` 只允许安全的列名，并把所有值放进参数列表，避免把用户
输入直接拼进 SQL；同时支持 PostgreSQL 的编号占位符和 MySQL 的问号占位符。

```go
query, args, err := Build(
	"SELECT id FROM orders",
	[]Filter{{Column: "tenant_id", Value: 7}},
	DialectPostgres,
)
// query: SELECT id FROM orders WHERE tenant_id = $1
```

列名不能通过普通参数绑定，因此必须来自服务端白名单；正则校验只是第二道防线。示例只处理等值条件，
实际项目可用枚举映射公开字段到固定数据库列，并为排序、分页和不同操作符分别设计受限结构。
