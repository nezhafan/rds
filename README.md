### 说明
- 本项目是对 `github.com/redis/go-redis/v9` 的再封装。

### 特性
- 用结构体分离了各数据类型。使用更加清晰和提示友好。
- 减少参数
  - 省略了 `context.Context` 参数
  - 在第一次 new 结构体的时候传参 `key`， 之后就不用每次传 `key`
- 优化返回
  - 对 `redis.Nil` 处理，替换为返回 `bool` 表示是否存在。
    - `String`: `GET()` 
    - `Hash`: `HGET()`、`HGetAll()` (特殊处理，结构体无法通过len判断有没有)
    - `List`: `LPOP()`、`RPOP()`、`LINDEX()`
    - `Set`: `SPOP()`
    - `ZSet`: `ZRANK()`、`ZREVRANK()`
- 使用 `SETNX` 增加了分布式锁的封装 `mutex.go`
- 泛型支持


### 用法举例
- String 
  - 存储值
  ```go
  
  ```
- Hash 
  - 存储对象，数据库的单条记录
  ```go
  
  ```

