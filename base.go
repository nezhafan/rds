package rds

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type base struct {
	key     string
	ctx     context.Context
	cmdable Cmdable
}

func NewBase(ctx context.Context, key string) (b base) {
	if ctx == nil {
		ctx = context.Background()
	}
	b.ctx = ctx
	b.key = key
	return
}

func (b *base) Key() string {
	return b.key
}

func (b *base) db() Cmdable {
	if b.cmdable != nil {
		return b.cmdable
	}
	return DB()
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

func (b *base) done(cmd redis.Cmder) {
	// 打印命令
	if debugMode == ModeCommand {
		fmt.Println(cmd.Args()...)
	} else if debugMode == ModeFull {
		fmt.Println(cmd.String())
	}
}
