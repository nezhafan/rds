
本项目是对 `github.com/redis/go-redis/v9` 的再封装。

1. 分离了各类型的操作方法，不再混在一起。
2. 使代码操作更加清晰和简便。 减少了参数和函数调用。
3. 增加对几个不常用类型`bitmap、hyperloglog、geo`的简单举例
4. 增加了分布式锁的封装`mutex.go`
5. 具体的使用查看 `z_test.go`

```go
// go-redis
const key1 = "key1"
rdb.Set(ctx, key1, "value", time.Hour)
rdb.Get(ctx, key1).Val()

// 封装后，减少了参数传入
key1 := redis.NewString("key1")
key1.Set("value", time.Hour)
key1.Get()
```

