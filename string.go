package redis

import (
	"time"
)

type str struct {
	base
}

func NewString(key string) str {
	return str{base{key}}
}

// 合并了Set和SetEX。 强制设置时间
func (s str) Set(val any, exp time.Duration) error {
	return rdb.Set(ctx, s.key, val, exp).Err()
}

func (s str) SetNX(val any, exp time.Duration) (bool, error) {
	return rdb.SetNX(ctx, s.key, val, exp).Result()
}

// 注意的是，如果不存在则error为redis.Nil
func (s str) Get() string {
	return rdb.Get(ctx, s.key).Val()
}

// 注意如果一个值已经是小数，则必须使用IncrByFloat
func (s str) IncrBy(incr int64) int64 {
	return rdb.IncrBy(ctx, s.key, incr).Val()
}

func (s str) IncrByFloat(incr float64) float64 {
	return rdb.IncrByFloat(ctx, s.key, incr).Val()
}
