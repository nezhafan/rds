### 一、说明
- 本项目是对 `github.com/redis/go-redis/v9` 的再封装。
- 函数按照数据类型分类。结构体形式调用，提示友好，用法清晰。
- 类型拆分。新增了 分布式锁、整数、浮点数、布尔、JSON、结构体 存储。
- 泛型。存储和获取时自动转换。
- 删除了 `error` 可能是 `redis.Nil` 的场景，增加 `R()` 方法显式返回是否存在。
- 复杂的方法增加了注释和传参转换，更易上手。基本避免了传参错误。
- 一些方法。
  - `Connect` 、`ConnectByOption` 连接
  - `SetDB` 使用自定义连接，`GetDB` 获取连接。
  - `SetDebug` 把命令和结果打印出来，方便本地调试，线上不建议开启。
  - `SetPrefix`，为所有key加上前缀。
  - `SetErrorHook` 捕捉 `error` 后进行自定义处理日志。
  - `Pipeline` 管道、`TxPipeline` 事务管道。将各类型串联起来。
- 有用的话，欢迎 star。

### 二、类型和用法
> 也可参考 test 测试用例
#### 1.string
- `Mutex` 分布式锁（可重入）
  - 包含 `Lock()`、`Unlock`、`TryLock` 方法
  - 示例代码
    > 方式一：不阻塞。直接尝试获取
    ```go
    mu := rds.NewMutex(ctx, "order:"+uuid)
    if !mu.TryLock() {
      return errors.New("有其它请求在占用")
    }
    defer mu.Unlock()
    // 业务逻辑...
    ```
    > 方式二：阻塞。尝试获取，拿不到重试，直到：拿到锁true/上下文超时false/上个锁到期true(避免这种情况)。
    ```go
    mu := rds.NewMutex(ctx, "order:"+uuid) 
    if !mu.Lock() {
      return errors.New("没拿到锁，上下文超时了")
    }
    defer mu.Unlock()
    // 业务逻辑...
    ```

- `String` 存储字符串。
  - 包含 `Set`、`SetNX`、`Get` 方法。
  - 示例代码
    ```go
    cache := rds.NewString(ctx, "key_string")
    cache.Set("okk", time.Minute)
    // 方式一：只关心值. 拿不到就去数据库再查。
    v := cache.Get().Val()
    fmt.Println(v)
    // 方式二：需要考虑redis错误，如果错误直接返回之类。
    v, err := cache.Get().Result()
    fmt.Println(v, err)
    // 方式三：值可能是空字符串，显式判断
    exists, v, err := cache.Get().R()
    fmt.Println(exists, v, err)
    ```
- `Bool` 存储布尔值（实际存1和0）。
  - 包含 `Set`、`SetNX`、`Get` 方法。
  - 示例代码
    ```go
    // 若只是为了防止某一时刻的并发，简单setnx互斥，不需要复杂的Mutex
    cache := rds.NewBool(ctx, "key_bool")
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
    cache := rds.NewInt64(ctx, "key_int64")
    v := cache.IncrBy(10).Val()
    // 使用 Incr的方式不方便设置过期，所以可以自行判断首次增长值
    if v == 10 {
      cache.Expire(time.Minute)
    }
    ```
- `Float64` 存储浮点数，存储浮点数要注意精度问题。
  - 包含 `Set`、`SetNX`、`Get`、`IncrByFloat` 方法。
  - 示例代码
    ```go
    cache := rds.NewFloat64(ctx, "key_float64")
    defer cache.Del()
    v := cache.IncrByFloat(0.3).Val()
    fmt.Println(v)
    // 使用 Incr的方式不方便设置过期，所以可以自行判断首次增长值
    if v == 0.3 {
      cache.Expire(time.Minute)
    }
    ```
- `Bit` 即位图bitmap，实际是对字符串按位拆解。
  - 包含 `SetBit`、`GetBit`、`BitCount`、`BitPos` 方法。
  - 示例代码
    ```go
    cache := rds.NewBit(ctx, "key_bit")
    defer cache.Del()
    // 举例某个用户的任务达成
    cache.SetBit(uid, 1)
    // 判断某用户任务是否达成
    if cache.GetBit(uid).Val() == 1 {

    }
    // 统计达成人数
    cache.BitCount().Val()
    ```

#### 2.hash
- `HashStruct[any]` 存储任意结构体，返回数据自动转换为结构体。
  - 包含 `HSet`、`HGet`、`HMGet`、`HGetAll` 方法。 
  - 注意⚠️：必须设置 `redis:"xx"` 标签，才能存储。忘记设置标签会报错。
  - 示例代码
    ```go
    // 以User对象举例
    type User struct {
      Id   int    `json:"id"` // id没redis标签不会存储
      Name string `redis:"name" json:"name"`
      Age  int    `redis:"age" json:"age"`
    }
    u := User{Id: 3306, Name: "Alice", Age: 20}
    cache := rds.NewHashStruct[User](ctx, "key_hash_struct")
    defer cache.Del()
    // 获取不存在的对象 返回nil，告知无缓存
    exists, u1, err := cache.HGetAll().R()
    fmt.Println(exists, u1, err)
    // 缓存nil对象，返回nil，告知有缓存
    cache.HSet(nil, time.Minute)
    exists, u2, err := cache.HGetAll().R()
    fmt.Println(exists, u2, err)
    // 设置有效对象
    cache.HSet(&u, time.Minute)
    exists, u3, err := cache.HGetAll().R()
    fmt.Println(exists, u3, err)
    // 获取单个字段 （也返回整个结构体，但是只取需要的即可）
    age := cache.HGet("age").Val()
    fmt.Println(age)
    // 获取部分字段 （也返回整个结构体，但是只取需要的即可）
    exists, u4, err := cache.HMGet("age").R()
    fmt.Println(exists, u4, err)
    ```

- `HashMap[cmp.Ordered]` 存储一类键值对，具有统一的value类型。
  - 包含 `HSet`、`HSetNX`、`HGet`、`HMGet`、`HGetAll`、`HIncrBy`、`HIncrByFloat`、`HDel`、`HExists`、`HLen` 方法。
  - 示例代码
    ```go
    // 以 value 类型为 int 举例
    cache := rds.NewHashMap[int](ctx, "key_hash_map")
    defer cache.Del()
    // 设置值（不设置过期时间，需要自己去管理）
    cache.HSet(map[string]int{"A": 1, "B": 2, "C": 3})
    // 增长 (注意类型为int值时，不可使用IncryByFloat)
    cache.HIncrBy("B", 1)
    // 删除
    cache.HDel("C")
    // 获取所有值
    v1 := cache.HGetAll().Val()
    fmt.Println(v1)
    // 获取部分字段 （也返回整个结构体，但是只取需要的即可）
    v2 := cache.HMGet("A", "C").Val()
    fmt.Println(v2)
    ```

#### 3.list
- `List[any]` 存储队列元素。
  - 包含 `LPush`、`RPush`、`LPop`、`RPop`、`LSet`、`LIndex`、`LRange`、`LLen`、`LRem`、`LTrim`
  - 示例代码LLen
    ```go
    cache := rds.NewList[string](ctx, "key_list")
    defer cache.Del()
    // 加入队列
    cache.RPush("a", "b", "c")
    // 取出队列
    exists, v, err := cache.LPop().R()
    fmt.Println(exists, v, err)
    // 查看队列数量
    size := cache.LLen().Val()
    fmt.Println(size)
    // 遍历队列
    vals := cache.LRange(0, -1).Val()
    fmt.Println(vals)
    ```
    ```go
    // 支持可JSON化的类型
    cache := rds.NewList[[]string](ctx, "key_list_2")
    defer cache.Del()
    // 加入队列
    cache.RPush([]string{"a"}, []string{"b", "c"}, []string{"d"})
    // 取出队列
    exists, v, err := cache.LPop().R()
    fmt.Println(exists, v, err)
    // 查看队列数量
    size := cache.LLen().Val()
    fmt.Println(size)
    // 遍历队列
    vals2 := cache.LRange(0, -1).Val()
    fmt.Println(vals2)
    ```

#### 4.set
- `Set[cmp.Ordered]` 存储去重元素
  - 包含 `SAdd`、`SIsMember`、`SMembers`、`SScan`、`SCard`、`SPop`、`SRem`方法。
  - 示例代码
    ```go
    // 以存储string类型举例
    cache := rds.NewSet[string](ctx, "key_set")
    defer cache.Del()
    // 添加
    cache.SAdd("a", "b", "c", "d", "e", "f", "g")
    // 元素个数
    size := cache.SCard().Val()
    fmt.Println(size)
    // 是否存在
    exists := cache.SIsMember("a").Val()
    fmt.Println(exists)
    // 随机删除并返回n个成员
    v1 := cache.SPop(2).Val()
    fmt.Println(v1)
    // 返回所有成员(注意控制返回的数量不要巨大)
    v2 := cache.SMembers().Val()
    fmt.Println(v2)
    // 数量巨大可以使用游标，多次查看
    err := cache.SScan("", 10, func(vals []string) error {
      fmt.Println(vals)
      return nil
    })
    fmt.Println(err)
    ```

#### 5.sorted set
- `SortedSet[cmp.Ordered]` 存储分值排序元素。
  - 包含 `ZAdd`、`ZIncrBy`、`ZCard`、`ZCountByScore`、`ZScore`、`ZRank`、`ZMembersByScore`、`ZMembersByRank`、`ZRangeByScore`、`ZRangeByRank`、`ZRem`、`ZRemByRank`、`ZRemByRank`方法。
  - 示例代码
    ```go
    cache := rds.NewSortedSet[string](ctx, "key_sorted_set")
    defer cache.Del()
    // 插入/修改数据 （若仅新增/仅修改/需要查看新增数或修改数，参考说或 test/sorted_set_test.go）
    cache.ZAdd(map[string]float64{
      "a": 1.1,
      "b": 2.2,
      "c": 2.2,
      "d": 4.4,
    })
    // 获取 c 的分值
    v1, err := cache.ZScore("c").Result()
    fmt.Println(v1, err)
    // 增加 c 的分值
    v2, err := cache.ZIncrBy("c", 1).Result()
    fmt.Println(v2, err)
    // 获取 c 的分值
    v3, err := cache.ZScore("c").Result()
    fmt.Println(v3, err)
    // 获取 c 的排名 (分值从高到低)
    v4, err := cache.ZRank(rds.ScoreDESC, "c").Result()
    fmt.Println(v4, err)
    // 获取 c 的排名 (分值从低到高)
    v5, err := cache.ZRank(rds.ScoreASC, "c").Result()
    fmt.Println(v5, err)
    // 查询前两名 （分值从高到低）
    v6, err := cache.ZRangeByRank(rds.ScoreDESC, 0, 1).Result()
    fmt.Println(v6, err)
    // 查询分值在 2.2-3.3 的成员 （分值从高到低）
    v7, err := cache.ZRangeByScore(rds.ScoreDESC, 2.2, 3.3, 0, 0).Result()
    fmt.Println(v7, err)
    // 查询分值区间内有多少人
    v8, err := cache.ZCountByScore(2.2, 3.3).Result()
    fmt.Println(v8, err)
    // 移除 a
    cache.ZRem("a")
    // 移除分值最低的两个
    cache.ZRemByRank(rds.ScoreASC, 0, 1)
    // 移除分值在 2.2 - 3.3 的所有成员
    cache.ZRemByScore(2.2, 3.3)
    ```
- `GEO` 存储经纬度坐标。 
  - 包含 `GeoAdd`、`GeoPos`、`GeoDist`、`GeoSearchByCoord`、`GeoSearchByMember`、`GeoDel` 方法
  - 参考 `test/geo_test.go`


#### 6.hyperloglog
- `HyperLogLog` 存储基数统计。
  - 包含 `PFAdd`、`PFCount`、`PFMerge` 方法
  - 使用 `type` 查看会说它是 `string` ，但是实际并非SDS，是借用了`string`类型便于处理一些其它逻辑
  - 参考 `test/hyperloglog_test.go`

