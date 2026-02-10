package rds

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type base struct {
	key  string
	ctx  context.Context
	pipe redis.Pipeliner
}

func newBase(ctx context.Context, key string) (b base) {
	if ctx == nil {
		ctx = context.Background()
	}
	b.ctx = ctx
	b.key = key
	return
}

func (b *base) Key() string {
	if keyPrefix == "" {
		return b.key
	}
	return keyPrefix + b.key
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
		cmd = b.db().Persist(b.ctx, b.Key())
	} else {
		cmd = b.db().Expire(b.ctx, b.Key(), exp)
	}
	b.done(cmd)
	return newBoolCmd(cmd)
}

func (b *base) ExpireAt(expAt time.Time) BoolCmd {
	cmd := b.db().ExpireAt(b.ctx, b.Key(), expAt)
	b.done(cmd)
	return newBoolCmd(cmd)
}

func (b *base) Exists() BoolCmd {
	cmd := b.db().Exists(b.ctx, b.Key())
	b.done(cmd)
	return newBoolCmd(cmd)
}

func (b *base) Del() BoolCmd {
	cmd := b.db().Del(b.ctx, b.Key())
	b.done(cmd)
	return newBoolCmd(cmd)
}

func (b *base) TTL() DurationCmd {
	cmd := b.db().TTL(b.ctx, b.Key())
	b.done(cmd)
	return newDurationCmd(cmd)
}

var (
	onceExpire sync.Map
)

// 设置一次过期时间
func (b *base) OnceExpire(exp time.Duration) BoolCmd {
	if _, ok := onceExpire.LoadOrStore(b.Key(), struct{}{}); ok {
		return newBoolCmd(new(redis.BoolCmd))
	}
	cmd := b.Expire(exp)
	if !cmd.Val() {
		onceExpire.Delete(b.Key())
	}
	return cmd
}

func (b *base) OnceExpireAt(expAt time.Time) BoolCmd {
	if _, ok := onceExpire.LoadOrStore(b.Key(), struct{}{}); ok {
		return newBoolCmd(new(redis.BoolCmd))
	}
	cmd := b.ExpireAt(expAt)
	if !cmd.Val() {
		onceExpire.Delete(b.Key())
	}
	return cmd
}

func (b *base) done(cmd redis.Cmder) {
	// 开发模式打印命令和结果()
	if isDebugMode {
		_, _ = debugWriter.WriteString(cmd.String() + "\n")
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
