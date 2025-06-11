package rds

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// type CmdableDo interface {
// 	redis.Cmdable
// 	Do(ctx context.Context, args ...interface{}) *redis.Cmd
// }

// func Do(ctx context.Context, cmd, key string, args ...any) (any, error) {
// 	cmds := append([]any{cmd, key}, args...)
// 	return DB().Do(ctx, cmds...).Result()
// }

var (
	ctx = context.Background()
)

type base struct {
	key  string
	exp  time.Duration
	pipe redis.Pipeliner
}

func newBase(key string, ops ...Option) (b base) {
	if len(allKeyPrefix) == 0 {
		b.key = key
	} else {
		b.key = allKeyPrefix + ":" + key
	}
	for _, op := range ops {
		op(&b)
	}
	return
}

func (b *base) db() redis.Cmdable {
	if b.pipe != nil {
		return b.pipe
	}
	return DB()
}

func (b *base) done() {
	if b.exp > 0 {
		b.exp = 0
		b.Expire(b.exp)
	}
}

func (b *base) Expire(exp time.Duration) bool {
	return DB().Expire(ctx, b.key, exp).Val()
}

func (b *base) Exists() bool {
	i := DB().Exists(ctx, b.key).Val()
	return i == 1
}

func (b *base) Del() bool {
	return DB().Del(ctx, b.key).Val() == 1
}

func (b *base) TTL() time.Duration {
	return DB().TTL(ctx, b.key).Val()
}

type Option func(b *base)

func WithPipe(pipe redis.Pipeliner) Option {
	return func(b *base) {
		b.pipe = pipe
	}
}

func WithExpire(exp time.Duration) Option {
	return func(b *base) {
		b.exp = exp
	}
}
