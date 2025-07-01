package rds

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// 文档 https://redis.uptrace.dev/zh/guide/
var (
	rdb redis.UniversalClient
)

func DB() redis.UniversalClient {
	return rdb
}

func SetDB(db redis.UniversalClient) {
	rdb = db
}

func Connect(addr string, auth string, db int) error {
	timeout := time.Second * 3
	options := &redis.Options{
		Addr:         addr,
		Password:     auth,
		DB:           db,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		DialTimeout:  timeout,
	}
	return ConnectByOption(options)
}

func ConnectByOption(option *redis.Options) error {
	rdb = redis.NewClient(option)
	if err := DB().Ping(context.Background()).Err(); err != nil {
		return err
	}
	return nil
}
