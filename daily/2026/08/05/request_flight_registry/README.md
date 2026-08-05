# 可取消的请求合并注册表

热点请求同时到达时，可让相同 key 的调用共享一次后端计算。`Registry` 还记录等待者数量：只有最后一个等待者取消时，才取消底层任务，避免单个客户端断开误伤其他调用者。

```go
registry := request_flight_registry.New[string]()
value, err := registry.Do(ctx, "user:42", loadUser)
```

该示例适合合并幂等读请求。生产环境还应设置任务超时，避免计算函数忽略 `context` 后永久占用注册表。
