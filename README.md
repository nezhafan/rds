### 说明
- 本项目是对 `github.com/redis/go-redis/v9` 的再封装。

### 特点
- 按数据类型分离结构体。使用时更加清晰和提示友好，同时也预防了类型冲突错误。
  - 所有类型都嵌入 `base` 结构体，继承 `Exists`、`Del`、`TTL`、`Expire`、`ExpireAt` 公共方法。
  - `string` 类型。
    - `String` 存储字符串，支持 `Get`、`Set`、`SetNX` 方法。
    - `StringInt` 存储整数，支持 `Get`、`Set`、`SetNX`、`IncrBy` 方法。
    - `StringFloat` 存储浮点数，支持 `Get`、`Set`、`SetNX`、`IncrBy` 方法。
    - `StringJSON[T]` 存储各种可`json`化的对象，包含`struct`、`slice`、`map`，支持 `Get`、`Set`、`SetNX` 方法。
    - `Mutex` 分布式锁，支持 `Lock`、`Unlock`、`TryLock` 方法。 `Lock`时会阻塞，每10ms-20ms重试获取锁，直到超过最大事件或上下文时间耗尽。
  - `hash` 类型。 
    - `HashStruct` 适用于存储对象信息，数据表的行。
    - `HashMap` 适用于存储一类键值对，如统一记录所有用户某个事的完成时间，需要同时取多个用户的该值。
  - `list` 类型。`List[T]` 队列，自定义元素类型。
  - `set` 类型。 `Set[T]` 去重，自定义元素类型。
  - `sorted set` 类型。 `SortedSet[T]` 做排序，自定义 `member` 类型。
  - `bitmap` 类型。 实际底层还是 `string` 类型，用做大批量数字的布尔值标记。
  - `geo` 类型。 `Geo` 做地理位置的标记和距离计算。
  - `hyperloglog` 类型。 `NewHyperLogLog` 做基数统计，非常节省存储，但存在0.81%误差。
- 参数/返回
  - 在第一次 new 结构体的时候传参 `key`，之后就不用每次传 `key`，很多类型可以直接声明在函数外。
  - 使用泛型对返回值进行自动转换。
  - 对于`string` ，从 `error` 中提取出来 `redis.Nil`，会返回3个值：val、exists、error 用以区分空字符串和不存在的两种情况。
- 辅助方法
  - 连接。可以使用这里的简易 `Connect`方法，也支持 `SetDB()` 自己实现连接后赋值进来。
  - DEBUG。支持`SetDebug`模式，会将执行的命令和返回内容清晰地打印，可在本地环境打开。
  - 事务/管道。 支持 `TxPipelined`事务 和 `TxPipelined` 普通管道
  

### 用法举例

##### 1. 连接 
```go
ctx := context.Background()

// 简单连接
if err := rds.Connect("127.0.0.1:6379", "password", 0); err != nil {
  panic(err)
}

// 自定义连接
options := &redis.Options{
  Addr:     "",
  Password: "",
  DB:       0,
}
if err := rds.ConnectByOption(options); err != nil {
  panic(err)
}

// 使用已有连接，把本包当作扩展。
rds.SetDB(db Cmdable)
```

##### 2. String 
- 存储值
```go
cache := rds.NewString(ctx, "string")
// 设定值，必须设置过期时间（永久有效可以使用rds.KeepTTL）
cache.Set("aaa", time.Minute)
// 当你防重复设置时，可以这样做
if !cache.SetNX("bbb", time.Minute).Val() {
  fmt.Println("重复设置")
}
// 获取值
val := cache.Get().Val()
fmt.Println(val)
// 如果需要区分空值还是不存在，或者需要判断error
val, exists, err := cache.Get().Result()
fmt.Println(val, exists, err)
```
- 存储整型
```go
cache := rds.NewStringInt(ctx, "string_int")
// 设定值，必须设置过期时间（永久有效可以使用rds.KeepTTL）
cache.Set(100, time.Minute)
// 当你防重复设置时，可以这样做
if !cache.SetNX(200, time.Minute).Val() {
  fmt.Println("重复设置")
}
// 获取值
val := cache.Get().Val()
fmt.Println(val)
// 如果需要区分空值还是不存在，或者需要判断error
val, exists, err := cache.Get().Result()
fmt.Println(val, exists, err)
// 数字类型，当然通常需要自增操作。 返回增加后的值
val = cache.IncrBy(2).Val()
fmt.Println(val)
```
- 存储浮点数
```go
cache := rds.NewStringFloat(ctx, "string_float")
// 设定值，必须设置过期时间（永久有效可以使用rds.KeepTTL）
cache.Set(0.3, time.Minute)
// 当你防重复设置时，可以这样做
if !cache.SetNX(0.3, time.Minute).Val() {
  fmt.Println("重复设置")
}
// 获取值
val := cache.Get().Val()
fmt.Println(val)
// 如果需要区分空值还是不存在，或者需要判断error
val, exists, err := cache.Get().Result()
fmt.Println(val, exists, err)
// 数字类型，当然通常需要自增操作。 返回增加后的值
val = cache.IncrBy(2).Val()
fmt.Println(val)
```
- 存储JSON
```go
type Config struct {
  App  string `json:"app"`
  Host string `json:"host"`
}
cache := rds.NewStringJSON[Config](ctx, "string_json")
// 设定值，必须设置过期时间。（永久有效可以使用rds.KeepTTL）
cache.Set(&Config{App: "app", Host: "https://cc.com"}, time.Minute)
// 当你防重复设置时，可以这样做
if !cache.SetNX(&Config{}, time.Minute).Val() {
  fmt.Println("重复设置")
}
// 获取值. 
// 无缓存时， val是nil， exists是false
// 有缓存但是缓存为空或null时， val是nil， exists是true
val, exists, err := cache.Get().Result()
fmt.Println(val,exists, err)
```

##### 3. Hash 
- 存储对象，数据库的单条记录
```go
// 注意hash存储和获取对象时，需要设置redis:"xxx"标签
type User struct {
  Id   int    `redis:"id"`
  Name string `redis:"name"`
}

// 缓存主key，可以声明为全局变量，说明所有的user都使用这个key前缀
var (
  userCache = NewHashStruct[User]("hash_struct")
)

func main() {
  user := User{
    Id:   123,
    Name: "张三",
  }
  SetUser(&user)
  u := GetUser(123)
  fmt.Println(u)
}

// 缓存用户
func SetUser(user *User) {
  // 函数内使用sub来分离区分是哪个用户
  cache := userCache.SubID(user.Id)
  // 设定值，必须设置过期时间（永久有效可以使用KeepTTL）
  cache.HMSet(user, time.Minute)
}

// 举例获取用户
func GetUser(id int, fields ...string) (u *User) {
  // 函数内使用sub来分离区分是哪个用户
  cache := userCache.SubID(strconv.Itoa(id))
  // fields不传时获取所有字段，否则获取指定字段
  u = cache.HMGet(fields...).Val()
  // 自己的缓存逻辑
  if u == nil {

  }
  return
}
```
- 存储一类键值对
```go
// 以存储数值为例.  
cache := NewHashMap[int64]("hash_map")
// 过期
cache.Expire(time.Minute)
// 设定单个值
cache.HSet("A", 0)
// 批量设定值 (不支持map[string]E ，所以自己确保any的值是E类型)
cache.HMSet(map[string]any{"A": 10, "B": 10})
// 数量
length := cache.HLen().Val()
fmt.Println(length)
// 获取单个值
val := cache.HGet("A").Val()
fmt.Println(val)
// 获取所有
all := cache.HGetAll().Val()
fmt.Println(all)
// 删除
cache.HDel("B", "D")
// 增长 (HIncrByInt 和 HIncrByFloat 分别对应 int64 和 float64 时，混用会error)
val = cache.HIncrByInt("A", 10).Val()
fmt.Println(val)
```

##### 4. Set 
```go
cache := NewSet[int64]("set")
// 过期
cache.Expire(time.Minute)
// 添加值，返回成功添加数
val := cache.SAdd(1, 2).Val()
fmt.Println(val)
val = cache.SAdd(1, 2, 3).Val()
fmt.Println(val)
// 判断是否存在
if cache.SIsMember(3).Val() {
  fmt.Println("存在")
}
// 删除
cache.SRem(3)
// 长度
length := cache.SCard().Val()
fmt.Println(length)
// 所有元素
all := cache.SMembers().Val()
fmt.Println(all)
```

##### 5. SortedSet
```go
// score只能是数字， member处理为数字或字符串
cache := rds.NewSortedSet[string](ctx, "sorted_set")
// 过期
cache.Expire(time.Minute)
// 添加。 member不可以重复，score可以重复
cache.ZAdd(map[string]float64{
  "A": 98.5,
  "B": 90,
  "C": 88.5,
})
// 元素数量
length := cache.ZCard().Val()
fmt.Println(length)
// 所有元素
length = cache.ZCount(90, 100).Val()
fmt.Println(length)
// 给A加1.5分到100
score := cache.ZIncrBy("A", 1.5).Val()
fmt.Println(score)
// 获取A的分数
score = cache.ZScore("A").Val()
fmt.Println(score)

// 积分区间内的成员 （从小到大）
members := cache.ZMembersByScore(90, 100, 0, 10, rds.ASC).Val()
fmt.Println(members)
// 积分区间内的成员 （从小到大）
members = cache.ZMembersByScore(90, 100, 0, 10, rds.DESC).Val()
fmt.Println(members)
// 排名区间内的成员 （从小到大）
members = cache.ZMembersByRank(0, 2, rds.ASC).Val()
fmt.Println(members)
// 排名区间内的成员 （从小到大）
members = cache.ZMembersByRank(0, 2, rds.DESC).Val()
fmt.Println(members)
// 积分区间内的成员，附带分数 （从小到大）
items := cache.ZRangeByScore(90, 100, 0, 10, rds.ASC).Val()
fmt.Println(items)
// 积分区间内的成员，附带分数 （从小到大）
items = cache.ZRangeByScore(90, 100, 0, 10, rds.DESC).Val()
fmt.Println(items)
// 排名区间内的成员，附带分数 （从小到大）
items = cache.ZRangeByRank(0, 2, rds.ASC).Val()
fmt.Println(items)
// 排名区间内的成员，附带分数 （从小到大）
items = cache.ZRangeByRank(0, 2, rds.DESC).Val()
fmt.Println(items)

// 移除排名区间内的成员 (移除排名第二的，即B)
cache.ZRemByRank(1, 1)
// 移除成员
cache.ZRem("D")
// 移除分数区间内的成员 (移除C)
cache.ZRemByScore(0, 89)
```

##### 6. List
```go
cache := rds.NewList[int](ctx, "list")
// 过期
cache.Expire(time.Minute)
// 插入顺序 1、2、3、4、5
cache.LPush(2, 1)
cache.RPush(3, 4, 5)
// 读取
all := cache.LRange(0, -1).Val()
fmt.Println(all)
// 当前长度
length := cache.LLen().Val()
fmt.Println(length)
// 左弹出
val := cache.LPop().Val()
fmt.Println(val)
// 右弹出
val = cache.RPop().Val()
fmt.Println(val)
// 修改第一个为333
cache.LSet(0, -1)
// 读取第一个
val = cache.LIndex(0).Val()
fmt.Println(val)
// 截断，左闭右闭
cache.LTrim(1, 1)
// 读取
all = cache.LRange(0, -1).Val()
fmt.Println(all)
```

##### 7. Bitmap
```go
var point uint32 = 999
cache := rds.NewBitmap(ctx, "bitmap")
// 过期
cache.Expire(time.Minute)
// 第一次设置 （返回未设置之前的状态false）
val := cache.SetBit(point, true).Val()
fmt.Println(val)
// 第二次设置（返回上次的状态true）
val = cache.SetBit(point, true).Val()
fmt.Println(val)
// 获取 （返回状态true）
ok := cache.GetBit(point).Val()
fmt.Println(ok)
// 统计区间内1的位数
n := cache.BitCount(0, int64(point)).Val()
fmt.Println(n)
```

#### 8.事务管道
```go
var cmd1 Cmder[map[string]int]
var cmd2 Cmder[time.Duration]

err := TxPipelined(context.Background(), func(p redis.Pipeliner) {
  // 需要手动 WithCmdable 
  cache := NewHashMap[int]("test_pipe").WithCmdable(p)
  // 此时命令不会真正去发送给redis-server
  cache.HSet("a", 1)
  cache.HSet("b", 2)

  // 错误方式 (在事务中无法直接拿到值赋值给变量)
  // result := cache.HGetAll().Val()
  // ttl :=  cache.TTL().Val()

  // 正确方式 - 第一步：拿到 cmd ，此时无值
  cmd1 = cache.HGetAll()
  cmd2 = cache.TTL()
})

if err != nil {
  fmt.Println(err)
  return
}

// 正确方式 - 第二步：从cmd中拿到值 （因为此时事务已经提交）
result := cmd1.Val()
ttl := cmd2.Val()
fmt.Println(result, ttl)
```

#### 9.分布式锁
```go
// 方式一：加锁。若失败则阻塞且不断重试
mu := rds.NewMutex(ctx, "1")
mu.Lock()
defer mu.Unlock()
// do...

// 方式二：加锁。若失败直接返回
mu := rds.NewMutex(ctx, "1")
if !mu.TryLock() {
  return ""
}
// do...
```