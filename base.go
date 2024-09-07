package rds

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// 文档 https://redis.uptrace.dev/zh/guide/

var (
	rdb *redis.Client
	ctx = context.Background()
)

const (
	OK  = "OK"
	Nil = redis.Nil
)

func Connect(addr string, auth string, db int) error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: auth,
		DB:       db,
	})

	return rdb.Ping(ctx).Err()
}

// func Get() *redis.Client {
// 	return rdb
// }

func Do(args ...any) (any, error) {
	return rdb.Do(ctx, args...).Result()
}

type base struct {
	key string
}

func newBase(key string) base {
	return base{key: key}
}

func (b *base) Expire(exp time.Duration) (bool, error) {
	return rdb.Expire(ctx, b.key, exp).Result()
}

func (b *base) Exists() (bool, error) {
	i, err := rdb.Exists(ctx, b.key).Result()
	return i == 1, err
}

func (b *base) Del() bool {
	return rdb.Del(ctx, b.key).Val() == 1
}

func (b *base) TTL() time.Duration {
	return rdb.TTL(ctx, b.key).Val()
}
