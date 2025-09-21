package rds

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// 事务管道（统一发送命令，统一执行不可被别的命令插入中间）
func TxPipelined(ctx context.Context, fn func(redis.Pipeliner)) error {
	p := DB().TxPipeline()
	fn(p)
	_, err := p.Exec(ctx)
	return err
}

// 管道（统一发送命令，执行过程可能被别的命令插入）
func Pipelined(ctx context.Context, fn func(redis.Pipeliner)) error {
	p := DB().Pipeline()
	fn(p)
	_, err := p.Exec(ctx)
	return err
}
