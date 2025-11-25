### 一、项目说明
本项目是对 github.com/redis/go-redis/v9 的轻量封装，核心优化如下：
- 按 Redis 数据类型分类封装函数，以结构体方式调用，提示友好、用法清晰；
- 扩展支持分布式锁、布尔 / 整数 / 浮点数 / JSON / 结构体等专属存储类型，通过 `rds.NewXXX` 快速使用；
- 基于泛型实现存储 / 获取的自动类型转换，无需手动序列化；
- 简化错误处理：移除 `error` 可能为 `redis.Nil` 的场景，新增 `R()` 方法显式返回值是否存在；
- 增强易用性：复杂方法补充注释和参数转换，减少传参错误；
- 支持 Key 前缀、调试日志、错误钩子、管道 / 事务管道等实用能力。
- 一些方法
  - `Connect` 、`ConnectByOption` 连接
  - `SetDB` 使用自定义连接，`GetDB` 获取连接。
  - `SetDebug` 把命令和结果打印出来，方便本地调试，线上不建议开启。
  - `SetPrefix`，为所有key加上前缀。
  - `SetErrorHook` 捕捉 `error` 后进行自定义处理日志。
  - `Pipeline` 管道、`TxPipeline` 事务管道。将各类型串联起来。在函数内部使用 `rds.NewPipeXXX` 来引用封装结构。
- 有用的话，欢迎 star。


### 二、核心用法示例


#### 1. 哈希
```go
/* HashStruct[any] 存储任意结构体，返回数据自动转换为结构体 */

// 以User对象举例(注意：结构体字段必须有redis标签，否则不会存储)
type User struct {
  Id   int    `json:"id"` // id没redis标签不会存储
  Name string `redis:"name" json:"name"`
  Age  int    `redis:"age" json:"age"`
}
u1 := User{Id: 3306, Name: "Alice", Age: 20}
cache := rds.NewHashStruct[User](ctx, "key_hash_struct")
defer cache.Del()
// 设置对象
cache.HSet(&u1, time.Minute)
// 查询
exists, u2, err := cache.HGetAll().R()
fmt.Println(exists, u2, err)
```

```go
/* HashMap[E] 存储一类键值对，具有统一的value类型 */

// 以 value 类型为 int 举例
cache := rds.NewHashMap[int](ctx, "key_hash_map")
defer cache.Del()
// 设置值（不设置过期时间，需要自己去管理）
cache.HSet(map[string]int{"A": 1, "B": 2, "C": 3})
// 获取所有值
v1 := cache.HGetAll().Val()
fmt.Println(v1)
// 获取部分字段 （也返回整个结构体，但是只取需要的即可）
v2 := cache.HMGet("A", "C").Val()
fmt.Println(v2)
```

#### 2. 字符串
```go
/* JSON类型泛型，自动序列化和反序列化 */
cache := rds.NewJSON[map[string]any](ctx, "key_json")
cache.Set(&map[string]any{"a": 1}, time.Minute)
v, err := cache.Get().Result()
exists, v, err := cache.Get().R()
fmt.Println(v, exists, err)
```
```go
/* int64类型 */
cache := rds.NewInt64(ctx, "key_int64")
v := cache.IncrBy(1).Val()
fmt.Println(v)
```
```go
/* string类型 */
rds.NewString(ctx, "key_string")
/* float64类型 */
rds.NewFloat64(ctx, "key_string")
/* bool类型 */
rds.NewBool(ctx, "key_bool")
```

#### 3.分布式锁
> 使用 string setnx + lua 脚本
```go
/* 方式一：不阻塞。直接尝试获取 */
mu := rds.NewMutex(ctx, "lock:"+uuid)
if !mu.TryLock() {
  return errors.New("有其它请求在占用")
}
defer mu.Unlock()
// 业务逻辑...
```
```go
/* 方式二：阻塞。尝试获取，拿不到重试，直到：拿到锁true/上下文超时false/上个锁到期true */
mu := rds.NewMutex(ctx, "lock:"+uuid) 
if !mu.Lock() {
  return errors.New("没拿到锁，上下文超时了")
}
defer mu.Unlock()
// 业务逻辑...
```

#### 更多类型和用法见 [demo.md](demo.md)