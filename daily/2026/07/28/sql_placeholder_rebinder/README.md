# SQL 占位符安全重绑定

不同数据库驱动使用不同占位符。`Rebind` 把通用的 `?` 依次转换为 PostgreSQL 的 `$1`、`$2`，同时跳过单引号字符串中的问号，并识别 SQL 的双单引号转义。

```go
query, count, err := Rebind("SELECT * FROM users WHERE id = ? AND state = ?")
// query: SELECT * FROM users WHERE id = $1 AND state = $2
```

这只是一个用于讲解状态机的最小实现，不是完整 SQL 解析器；它没有处理注释、美元引用字符串和方言特性。生产项目应优先使用成熟驱动或查询构造器。无论占位符形式如何，都必须通过参数绑定传值，不能字符串拼接用户输入。
