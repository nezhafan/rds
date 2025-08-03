package rds

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type base struct {
	key  string
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

func (b *base) Key() string {
	return b.key
}

func (b *base) Expire(exp time.Duration) *BoolCmd {
	cmd := b.db().Expire(b.ctx, b.key, exp)
	b.done(cmd)
	return &BoolCmd{cmd: cmd}
}

func (b *base) ExpireAt(expAt time.Time) *BoolCmd {
	cmd := b.db().ExpireAt(b.ctx, b.key, expAt)
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

func (b *base) done(cmd redis.Cmder) {
	// 打印命令
	if debugMode == ModeCommand {
		fmt.Println(cmd.Args()...)
	} else if debugMode == ModeFull {
		fmt.Println(cmd.String())
	}

}
