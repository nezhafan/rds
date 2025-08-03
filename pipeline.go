package rds

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// 事务
func TxPipelined(ctx context.Context, fn func(redis.Pipeliner)) error {
	pipe := DB().TxPipeline()
	fn(pipe)
	time.Sleep(time.Second * 2)
	_, err := pipe.Exec(ctx)
	return err
}

func Pipelined(ctx context.Context, fn func(redis.Pipeliner)) error {
	pipe := DB().Pipeline()
	fn(pipe)
	_, err := pipe.Exec(ctx)
	return err
}
