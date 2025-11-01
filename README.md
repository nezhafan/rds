### 说明
- 本项目是对 `github.com/redis/go-redis/v9` 的再封装。

### 改动
- 函数按照数据类型分类。结构体形式调用，提示友好，用法清晰。
- 类型拆分。新增了 分布式锁、整数、浮点数、布尔、JSON、结构体 存储。
- 泛型。不需要再手动 字符串 转 数字/布尔/JSON，存储和获取自动转换。
- 删除了 `redis.Nil`，增加 `R()` 方法显式返回是否存在。
- 复杂的方法增加了注释和传参转换，更易上手。基本避免了传参错误。
- 辅助方法。
  - `SetDebug`、`SetWriter` 把命令和结果打印出来，方便本地调试。
  - `SetPrefix`，为所有key加上前缀。
  - `SetErrorHook` 捕捉 `error` 自定义处理日志。

### 类型详情
> 也可参考 test 测试用例
#### `string` 
- `Mutex` 分布式锁
  - 包含 `Lock()`、`Unlock`、`TryLock` 方法
  - 若在 `Lock`时阻塞，使用从10ms开始的斐波那契数列间隔重试，直到上下文时间耗尽，或者达到缓存最大过期时间，第一把锁失效。
  - 示例代码
    ```go
    mu := rds.NewMutex(ctx, "order:"+uuid, 30)
    // 方式一：不阻塞。直接尝试获取
    if !mu.TryLock() {
      return errors.New("有其它请求在占用处理")
    }
    defer mu.Unlock()
    // 方式二：阻塞。尝试获取，拿不到重试，直到：拿到锁true/上下文超时false/上个锁到期true(避免这种情况)。 
    if !mu.Lock() {
      return errors.New("没拿到锁，上下文超时了")
    }
    defer mu.Unlock()
    ```

- `String` 存储字符串。
  - 包含 `Set`、`SetNX`、`Get` 方法。
  - 示例代码
    ```go
    cache := rds.NewString(ctx, "key_string")
    cache.Set("okk", time.Minute)
    // 方式一：只关心值. 拿不到就去数据库再查。
    v := cache.Get().Val()
    // 方式二：需要考虑redis错误，如果错误直接返回之类。
    v, err := cache.Get().Result()
    // 方式三：值可能是空字符串，显式判断
    exists, v, err := cache.Get().R() 
    ```
- `Bool` 存储布尔值（实际存1和0）。
  - 包含 `Set`、`SetNX`、`Get` 方法。
  - 示例代码
    ```go
    // 若只是为了防止某一时刻的并发，简单setnx互斥，不需要复杂的Mutex
    cache := rds.NewBool[map[string]any](ctx, "key_bool")
    if cache.SetNX(true, time.Second * 2).Val() {
      // do
    }
    ```
- `JSON[any]` 使用json形式存储 `map`、`slice`、`struct`等。 
  - 包含 `Set`、`SetNX`、`Get` 方法。
  - 示例代码
    ```go
    cache := rds.NewJSON[map[string]any](ctx, "key_json")
    cache.Set(&map[string]any{"a": 1}, time.Minute)
    v, err := cache.Get().Result()
    exists, v, err := cache.Get().R()
    ```
- `Int64` 存储整数。
  - 包含 `Set`、`SetNX`、`Get`、`IncrBy` 方法。
  - 示例代码
    ```go
    cache := rds.NewInt64[map[string]any](ctx, "key_int64")
    v := cache.IncrBy(10).Val()
    // 使用 Incr的方式不方便设置过期，所以可以自行判断首次增长值
    if v == 10 {
      cache.Expire(time.Minute)
    }
    ```
- `Float64` 存储浮点数。
  - 包含 `Set`、`SetNX`、`Get`、`IncrByFloat` 方法。
  - 示例代码
    ```go
    cache := rds.NewFloat64[map[string]any](ctx, "key_float64")
    v := cache.IncrByFloat(10.).Val()
    // 使用 Incr的方式不方便设置过期，所以可以自行判断首次增长值
    if v == 10. {
      cache.Expire(time.Minute)
    }
    ```
- `Bit` 即位图bitmap，实际是对字符串按位拆解。
  - 包含 `SetBit`、`GetBit`、`BitCount`、`BitPos` 方法。
  - 示例代码
    ```go
    // 例如存储用户某日登录情况
    cache := rds.NewBit(ctx, "key_bit:2025-01-01")
    // 登录时
    cache.SetBit(uid, 1)
    // 判断某用户是否登录
    if cache.GetBit(uid).Val() == 1 {

    }
    // 查看当日登录了多少人
    cache.BitCount().Val()
    // 
    ```

#### hash
- `HashStruct[any]` 存储任意结构体，返回数据自动转换为结构体。
  - 包含 `SubKey`、`HSetAll`、`HSet`、`HGet`、`HMGet`、`HGetAll`、`HIncrBy`、`HIncrByFloat`、`HDel`、`HExists` 方法。
  - 注意⚠️：必须设置 `redis:"xx"` 标签，才能存储。忘记设置标签会报错。
  - 示例代码
    ```go
    type User struct {
      Id int
      Name string `redis:"name"`
      Age int `redis:"age"`
    }
    u1 := User{Id: 3306, Name: "Alice", Age: 20}
    usercache := rds.NewHashStruct[User](ctx, "key_hash_struct")
    cache := usercache.SubKey("3306")
    cache.SetAll(&u1, time.Minute)
    // 获取整个
    if u2 := cache.HGetAll().Val(); u2 != nil {
       fmt.Println(u2.Name, u2.Age)
    }
    // 获取部分字段 （也返回整个结构体，但是只取需要的即可）
    if u3 := cache.HMGet("name").Val(); u3 != nil {
      fmt.Println(u3.Name, u3.Age)
    }
    // 修改字段
    cache.HSet(map[string]any{"age": 22}, time.Minute)
    // 增长年龄
    cache.IncrBy("age", 1)
    ```

- `HashMap[cmp.Ordered]` 存储一类键值对，返回`map[string]E`，`field`必须是字符串，`value`为泛型。包含 `HSet`、`HSetNX`、`HGet`、`HMGet`、`HGetAll`、`HIncrBy`、`HIncrByFloat`、`HDel`、`HExists`、`HLen`

  - `list` 
    - `List[any]` 存储队列元素。

  - `set` 
    - `Set[cmp.Ordered]` 存储去重元素。包含 `SAdd`、`SIsMember`、`SMembers`、`SRandMember`、`SCard`、`SPop`、`SRem`方法。

  - `sorted set` 
    - `SortedSet[cmp.Ordered]` 存储积分排序元素。包含 `ZAdd`、`ZIncrBy`、`ZCard`、`ZCountByScore`、`ZScore`、`ZIndex`、`ZMembersByScore`、`ZMembersByIndex`、`ZRangeByScore`、`ZRangeByIndex`、`ZRem`、`ZRemByIndex`、`ZRemByIndex`方法。

  - `geo`
    - `GEO` 经纬度坐标。 包含 
  - `hyperloglog`
    - `HyperLogLog` 存储基数统计。

