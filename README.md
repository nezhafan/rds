### 说明
- 本项目是对 `github.com/redis/go-redis/v9` 的再封装。

### 特性
- 用结构体分离了各数据类型。使用时更加清晰和提示友好，同时也预防了类型冲突错误。
- 类型拆分
  - 将 `Hash` 拆分为 `HashMap` 和 `HashStruct` 两种类型。 `HashMap` 使用于存储一类KV，其中V为`string、int64、float64`的泛型，取出时为 `map[string]E`;`HashStruct` 适用于存储对象，取出时为 `struct`。(`HashStruct` 要设定 `redis:"xxx"`标签)
  - 从 `String` 类型中拆分出了 `StringInt` 和 `StringFloat` 类型，用来做数字存储 `incr` 相关操作。
  - 从 `String` 类型中拆分出了 `StringJSON` 类型做对象存储。(优先推荐 `HashStruct` 类型做对象存储)
  - 从 `String` 类型中拆分出了 `Mutex` 分布式锁
- 泛型支持
  - `HashStruct` 和 `StringJSON` 支持自定义结构体，在存储和取出时自动转换。
  - `Set`、`List`、`SortedSet` 支持区分存储字符串还是数值，不会混淆。
  - `Bitmap` 位图，将0和1转为true和false
- 简化参数
  - 在第一次 new 结构体的时候传参 `key`，之后就不用每次传 `key`，很多类型可以直接声明在函数外。
  - 省略了 `context.Context` 参数。 1通过可选参数`WithContext`设置，2默认使用连接时设置的超时
  - 可选参数`WithPipe`设置管道，使用方式见下例。
- 优化返回。对 `redis.Nil` 处理，用以区分空字符串和不存在，涉及这个的会返回3个值：val、exists、error
- 全局参数
  - 支持debug模式`SetDebug`，会将执行的命令和返回内容清晰的打印
  - 增加了自动设置前缀方式 `SetPrefix()`

### 用法举例
##### 1. String 
- 存储值
```go
cache := rds.NewString("string")
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
cache := rds.NewStringInt("string_int")
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
cache := rds.NewStringFloat("string_float")
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
cache := rds.NewStringJSON[Config]("string_json")
// 设定值，必须设置过期时间。（永久有效可以使用rds.KeepTTL）
cache.Set(&Config{App: "app", Host: "https://cc.com"}, time.Minute)
// 当你防重复设置时，可以这样做
if !cache.SetNX(&Config{}, time.Minute).Val() {
  fmt.Println("重复设置")
}
// 获取值. （不存在可以直接判断指针是否为nil）
val, err := cache.Get().Result()
fmt.Println(val, err)
```

##### 2. Hash 
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

##### 3. Set 
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

##### 4. SortedSet
```go
// score只能是数字， member处理为数字或字符串
cache := rds.NewSortedSet[string]("sorted_set")
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
items := cache.ZItemsByScore(90, 100, 0, 10, rds.ASC).Val()
fmt.Println(items)
// 积分区间内的成员，附带分数 （从小到大）
items = cache.ZItemsByScore(90, 100, 0, 10, rds.DESC).Val()
fmt.Println(items)
// 排名区间内的成员，附带分数 （从小到大）
items = cache.ZItemsByRank(0, 2, rds.ASC).Val()
fmt.Println(items)
// 排名区间内的成员，附带分数 （从小到大）
items = cache.ZItemsByRank(0, 2, rds.DESC).Val()
fmt.Println(items)

// 移除排名区间内的成员 (移除排名第二的，即B)
cache.ZRemByRank(1, 1)
// 移除成员
cache.ZRem("D")
// 移除分数区间内的成员 (移除C)
cache.ZRemByScore(0, 89)
```

##### 5. List
```go
cache := rds.NewList[int]("list")
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

##### 6. Bitmap
```go
var point uint32 = 999
cache := rds.NewBitmap("bitmap")
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

#### 7.事务管道
```go
var cmd1 Cmder[map[string]int]
var cmd2 Cmder[time.Duration]

err := TxPipelined(context.Background(), func(p redis.Pipeliner) {
  // 需要手动 withpipe
  cache := NewHashMap[int]("test_pipe", WithPipe(p))
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