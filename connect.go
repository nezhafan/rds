package rds

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// 文档 https://redis.uptrace.dev/zh/guide/
type Cmdable interface {
	redis.Cmdable
	Do(ctx context.Context, args ...any) *redis.Cmd
}

var (
	_ Cmdable = (redis.Pipeliner)(nil)
	_ Cmdable = (redis.UniversalClient)(nil)
	_ Cmdable = (*redis.Client)(nil)
	_ Cmdable = (*redis.ClusterClient)(nil)

	rdb Cmdable
)

func DB() Cmdable {
	return rdb
}

func SetDB(db Cmdable) {
	rdb = db
}

// Connect("127.0.0.1:6379", "", 0)
func Connect(addr string, auth string, db int) error {
	options := &redis.Options{
		Addr:     addr,
		Password: auth,
		DB:       db,
	}
	return ConnectByOption(options)
}

func ConnectByOption(option *redis.Options) error {
	rdb = redis.NewClient(option)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return err
	}
	return nil
}
