package rds

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type base struct {
	key  string
	ctx  context.Context
	pipe redis.Pipeliner
}

func NewBase(ctx context.Context, key string) (b base) {
	if ctx == nil {
		ctx = context.Background()
	}
	b.ctx = ctx
	b.key = keyPrefix + key
	return
}

func (b *base) Key() string {
	return b.key
}

func (b *base) db() Cmdable {
	if b.pipe != nil {
		return b.pipe
	}
	return GetDB()
}

func (b *base) Expire(exp time.Duration) BoolCmd {
	var cmd *redis.BoolCmd
	if exp == KeepTTL {
		cmd = b.db().Persist(b.ctx, b.key)
	} else {
		cmd = b.db().Expire(b.ctx, b.key, exp)
	}
	b.done(cmd)
	return newBoolCmd(cmd)
}

func (b *base) ExpireAt(expAt time.Time) BoolCmd {
	cmd := b.db().ExpireAt(b.ctx, b.key, expAt)
	b.done(cmd)
	return newBoolCmd(cmd)
}

func (b *base) Exists() BoolCmd {
	cmd := b.db().Exists(b.ctx, b.key)
	b.done(cmd)
	return newBoolCmd(cmd)
}

func (b *base) Del() BoolCmd {
	cmd := b.db().Del(b.ctx, b.key)
	b.done(cmd)
	return newBoolCmd(cmd)
}

func (b *base) TTL() DurationCmd {
	cmd := b.db().TTL(b.ctx, b.key)
	b.done(cmd)
	return newDurationCmd(cmd)
}

func (b *base) done(cmd redis.Cmder) {
	// 开发模式打印命令和结果
	if isDebugMode {
		debugWriter.WriteString(cmd.String())
		debugWriter.WriteString("\n")
	}

	// cmd 钩子
	if cmdHook != nil {
		cmdHook(cmd)
	}

	// 错误钩子
	if errorHook != nil && cmd.Err() != nil && cmd.Err() != redis.Nil {
		errorHook(cmd.Err())
	}

}
