package rds

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func Pipelined(ctx context.Context, fn func(redis.Pipeliner)) error {
	pipe := DB().Pipeline()
	fn(pipe)
	_, err := pipe.Exec(ctx)
	return err
}
