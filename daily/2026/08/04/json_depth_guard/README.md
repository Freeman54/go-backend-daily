# JSON 嵌套深度防护

攻击者可构造极深的 JSON，让递归式业务处理消耗大量栈或 CPU。`Validate` 使用标准库流式读取 token，在反序列化到业务结构前限制对象和数组的嵌套深度，同时校验 JSON 语法。

```go
body := []byte(`{"profile":{"name":"alice"}}`)
if err := json_depth_guard.Validate(body, 8); err != nil {
	// 返回 400；ErrTooDeep 可映射为明确的参数错误
}
```

深度从最外层容器算 1，标量值不增加深度。生产环境还应同时限制请求体字节数，因为深度限制不能防止超大的扁平数组。
