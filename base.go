package rds

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

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

func (b *base) db() Cmdable {
	if b.pipe != nil {
		return b.pipe
	}
	return DB()
}

func (b *base) done(cmd redis.Cmder) {
	// 若给定了过期时间，则在操作后设置。若key不存在，则不置为0
	// 这有个问题是，如果一直是是操作不存在的key
	// 那么这里就会一直额外开销，但是相对于忘记设置过期这是可接受的
	if b.exp > 0 && b.db().Expire(ctx, b.key, b.exp).Val() {
		b.exp = 0
	}

	if debugMode == ModeCommand {
		fmt.Println(cmd.Args()...)
	} else if debugMode == ModeFull {
		fmt.Println(cmd.String())
	}
}

func (b *base) Expire(exp time.Duration) bool {
	cmd := b.db().Expire(ctx, b.key, exp)
	b.done(cmd)
	return cmd.Val()
}

func (b *base) Exists() bool {
	cmd := b.db().Exists(ctx, b.key)
	b.done(cmd)
	return cmd.Val() == 1
}

func (b *base) Del() bool {
	cmd := b.db().Del(ctx, b.key)
	b.done(cmd)
	return cmd.Val() == 1
}

func (b *base) TTL() time.Duration {
	cmd := b.db().TTL(ctx, b.key)
	b.done(cmd)
	return cmd.Val()
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
