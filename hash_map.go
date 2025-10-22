package rds

import (
	"cmp"
	"context"

	"github.com/redis/go-redis/v9"
)

type HashMap[E cmp.Ordered] struct {
	base
}

// 需要自己管理key的过期时间
func NewHashMap[E cmp.Ordered](ctx context.Context, key string) *HashMap[E] {
	h := &HashMap[E]{base: NewBase(ctx, key)}
	return h
}

// hset 支持一次性设置多个字段，返回新增字段数（修改不算新增）
// hmset 已经从4.0视为被弃用 https://redis.io/docs/latest/commands/hmset/
func (h *HashMap[E]) HSet(kv map[string]E) Int64Cmd {
	if len(kv) == 0 {
		return newInt64Cmd(new(redis.IntCmd))
	}
	values := make([]any, 0, len(kv)*2)
	for k, v := range kv {
		values = append(values, k, v)
	}
	cmd := h.db().HSet(h.ctx, h.key, values...)
	h.done(cmd)
	return newInt64Cmd(cmd)
}

func (h *HashMap[E]) HSetNX(field string, value E) BoolCmd {
	cmd := h.db().HSetNX(h.ctx, h.key, field, value)
	h.done(cmd)
	return newBoolCmd(cmd)
}

func (h *HashMap[E]) HGet(field string) AnyCmd[E] {
	cmd := h.db().HGet(h.ctx, h.key, field)
	h.done(cmd)
	return newAnyCmd[E](cmd)
}

func (h *HashMap[E]) HMGet(fields ...string) MapCmd[E] {
	cmd := h.db().HMGet(h.ctx, h.key, fields...)
	h.done(cmd)
	return newMapCmd[E](cmd, fields)
}

func (h *HashMap[E]) HGetAll() MapCmd[E] {
	cmd := h.db().HGetAll(h.ctx, h.key)
	h.done(cmd)
	return newMapCmd[E](cmd, nil)
}

// 注意float有精度问题
func (h *HashMap[E]) HIncrBy(field string, incr E) AnyCmd[E] {
	cmd := h.db().Do(h.ctx, "hincrby", h.key)
	return newAnyCmd[E](cmd)
}

// 返回删除成功数
func (h *HashMap[E]) HDel(fields ...string) Int64Cmd {
	cmd := h.db().HDel(h.ctx, h.key, fields...)
	h.done(cmd)
	return newInt64Cmd(cmd)
}

// 返回field数量
func (h *HashMap[E]) HLen() Int64Cmd {
	cmd := h.db().HLen(h.ctx, h.key)
	h.done(cmd)
	return newInt64Cmd(cmd)
}

func (h *HashMap[E]) HExists(field string) BoolCmd {
	cmd := h.db().HExists(h.ctx, h.key, field)
	h.done(cmd)
	return newBoolCmd(cmd)
}
