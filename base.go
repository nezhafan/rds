package rds

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type CmdableDo interface {
	redis.Cmdable
	Do(ctx context.Context, args ...interface{}) *redis.Cmd
}

func Do(ctx context.Context, cmd, key string, args ...any) (any, error) {
	cmds := append([]any{cmd, key}, args...)
	return DB().Do(ctx, cmds...).Result()
}

var (
	timeout = time.Second * 2
)

func SetTimeout(t time.Duration) {
	timeout = t
}

type base struct {
	key    string
	ctx    context.Context
	cancel context.CancelFunc
}

func newBase(ctx context.Context, key string) (b base) {
	if ctx == nil {
		b.ctx, b.cancel = context.WithTimeout(context.Background(), timeout)
	} else {
		b.ctx = ctx
	}

	if len(allKeyPrefix) == 0 {
		b.key = key
	} else {
		b.key = allKeyPrefix + ":" + key
	}
	return
}

var (
	hooks []Hook
)

type Hook func(cmd redis.Cmder)

func (b *base) done(cmd redis.Cmder) {
	if b.cancel != nil {
		b.cancel()
	}

	if len(hooks) > 0 {
		for _, h := range hooks {
			h(cmd)
		}
	}

}

func (b *base) WithContext(ctx context.Context) *base {
	b.ctx = ctx
	return b
}

func (b *base) Expire(exp time.Duration) bool {
	return DB().Expire(b.ctx, b.key, exp).Val()
}

func (b *base) Exists() bool {
	i := DB().Exists(b.ctx, b.key).Val()
	return i == 1
}

func (b *base) Del() bool {
	return DB().Del(b.ctx, b.key).Val() == 1
}

func (b *base) TTL() time.Duration {
	return DB().TTL(b.ctx, b.key).Val()
}
