package rds

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type base struct {
	key  string
	exp  time.Duration
	ctx  context.Context
	pipe redis.Pipeliner
}

func newBase(key string, ops ...Option) (b base) {
	b.ctx = context.Background()
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
	if b.exp > 0 && b.db().Expire(b.ctx, b.key, b.exp).Val() {
		b.exp = 0
	}

	if debugMode == ModeCommand {
		fmt.Println(cmd.Args()...)
	} else if debugMode == ModeFull {
		fmt.Println(cmd.String())
	}
}

func (b *base) Expire(exp time.Duration) *BoolCmd {
	cmd := b.db().Expire(b.ctx, b.key, exp)
	b.done(cmd)
	return &BoolCmd{cmd: cmd}
}

func (b *base) Exists() *BoolCmd {
	cmd := b.db().Exists(b.ctx, b.key)
	b.done(cmd)
	return &BoolCmd{cmd: cmd}
}

func (b *base) Del() *BoolCmd {
	cmd := b.db().Del(b.ctx, b.key)
	b.done(cmd)
	return &BoolCmd{cmd: cmd}
}

func (b *base) TTL() *DurationCmd {
	cmd := b.db().TTL(b.ctx, b.key)
	b.done(cmd)
	return &DurationCmd{cmd: cmd}
}

type Option func(b *base)

func WithContext(ctx context.Context) Option {
	return func(b *base) {
		b.ctx = ctx
	}
}

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
