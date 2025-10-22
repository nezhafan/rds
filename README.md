### 说明
- 本项目是对 `github.com/redis/go-redis/v9` 的再封装。
- 数据类型方法分离。
- 类型拆分。
  - `string` 类型拆分 `NewInt64`、`NewFloat64`、`NewBool`、`NewJSON[any]`
  - `string` 类型增加 `mutex` 分布式互斥锁封装，同样支持 `Lock`、`Unlock`、`TryLock`方法，`Lock` 阻塞时随机每10ms-20ms重试获取锁，直到 获取锁或上下文时间耗尽或锁过期。
  - `hash` 类型拆分 `NewStruct[any]` 存储对象 和 `newHashMap[any]` 存储同类泛型字段值。

- 空值显式判断。对于部分可能返回 `redis.Nil` 的方法进行处理，`error` 不再可能是 `redis.Nil` ，如果需要区分空值和不存在，使用 `R()` 方法，返回是否存在。


### 文档补充中...


### 用法
- 可以先参考 `test` 中的代码