### 说明
- 本项目是对 `github.com/redis/go-redis/v9` 的再封装。
- 将原来混合的函数按照数据类型封装为结构体区分函数调用。提示友好，用法清晰。



### 类型详情
#### `string` 
- `String` 存储字符串。
  - 包含 `Set`、`SetNX`、`Get` 方法。
  - 示例代码
    ```go
    cache := rds.NewString(ctx, "key_string")
    cache.Set("okk")
    v, err := cache.Get().Result()
    exists, v, err := cache.Get().R() 
    ```
- `JSON[any]` 使用json形式存储 `map`、`slice`、`struct`等。 
  - 包含 `Set`、`SetNX`、`Get` 方法。
  - 示例代码
    ```go
    cache := rds.NewJSON[map[string]any](ctx, "key_json")
    // cache := rds.NewJSON[User](ctx, "key_json")
    cache.Set(&map[string]any{"a": 1})
    v, err := cache.Get().Result()
    exists, v, err := cache.Get().R() 
    ```
- `Int64` 存储整数。包含 `Set`、`SetNX`、`Get`、`IncrBy` 方法。
- `Float64` 存储浮点数。包含 `Set`、`SetNX`、`Get`、`IncrByFloat` 方法。
- `Bool` 存储布尔值（实际存1和0）。包含 `Set`、`SetNX`、`Get` 方法。
- `Bit` 即位图bitmap，实际是对字符串按位拆解。包含 `SetBit`、`GetBit`、`BitCount`、`BitPos` 方法。
- `Mutex` 额外单独封装了分布式锁。 有 `Lock`、`Unlock`、`TryLock`方法。在 `Lock`时也会发生阻塞，每 10ms-20ms随机重试，直到 上下文timeout 或者达到最大限定时间。

#### hash
- `HashStruct[any]` 存储任意结构体，返回数据自动转换为结构体。包含 `SubKey`、`HSetAll`、`HSet`、`HGet`、`HMGet`、`HGetAll`、`HIncrBy`、`HIncrByFloat`、`HDel`、`HExists` 方法。
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

