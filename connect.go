package rds

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// 文档 https://redis.uptrace.dev/zh/guide/
var (
	rdb *redis.Client
)

func Connect(addr string, auth string, db int) error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: auth,
		DB:       db,
	})

	if err := DB().Ping(context.Background()).Err(); err != nil {
		return err
	}

	return nil
}

func DB() redis.UniversalClient {
	return rdb
}
