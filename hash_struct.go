package rds

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type HashStruct[E any] struct {
	base
}

// 泛型为value元素
func NewHashStruct[E any](ctx context.Context, key string) *HashStruct[E] {
	return &HashStruct[E]{base: NewBase(ctx, key)}
}

func (h *HashStruct[E]) SubKey(ctx context.Context, subkey string) *HashStruct[E] {
	return NewHashStruct[E](ctx, h.key+":"+subkey)
}

// 返回该字段是否为新增字段（修改不算新增）
func (h *HashStruct[E]) HSet(field string, value any) *BoolCmd {
	cmd := h.db().HSet(h.ctx, h.key, field, value)
	h.done(cmd)
	return &BoolCmd{cmd: cmd}
}

// 强制设置过期时间
func (h *HashStruct[E]) HMSet(obj *E, exp time.Duration) *BoolCmd {
	cmder := h.db()
	var cmd redis.Cmder
	if _, ok := cmder.(redis.Pipeliner); ok {
		cmd = cmder.HSet(h.ctx, h.key, obj)
		cmder.Expire(h.ctx, h.key, exp)
	} else {
		pipe := cmder.Pipeline()
		cmd = pipe.HSet(h.ctx, h.key, obj)
		pipe.Expire(h.ctx, h.key, exp)
		pipe.Exec(h.ctx)
	}
	h.done(cmd)
	return &BoolCmd{cmd: cmd}
}

func (h *HashStruct[E]) HMGet(fields ...string) *StructCmd[E] {
	if len(fields) == 0 {
		return h.HGetAll()
	}
	cmd := h.db().HMGet(h.ctx, h.key, fields...)
	h.done(cmd)
	return &StructCmd[E]{cmd: cmd, fields: fields}
}

func (h *HashStruct[E]) HGetAll() *StructCmd[E] {
	cmd := h.db().HGetAll(h.ctx, h.key)
	h.done(cmd)
	return &StructCmd[E]{cmd: cmd}
}

// 返回删除成功数
func (h *HashStruct[E]) HDel(fields ...string) *IntCmd {
	cmd := h.db().HDel(h.ctx, h.key, fields...)
	h.done(cmd)
	return &IntCmd{cmd: cmd}
}

func (h *HashStruct[E]) HExists(field string) *BoolCmd {
	cmd := h.db().HExists(h.ctx, h.key, field)
	h.done(cmd)
	return &BoolCmd{cmd: cmd}
}

func (h *HashStruct[E]) HIncrBy(field string, incr int64) *IntCmd {
	cmd := h.db().HIncrBy(h.ctx, h.key, field, incr)
	h.done(cmd)
	return &IntCmd{cmd: cmd}
}

func (h *HashStruct[E]) HIncrByFloat(field string, incr float64) *IntCmd {
	cmd := h.db().HIncrByFloat(h.ctx, h.key, field, incr)
	h.done(cmd)
	return &IntCmd{cmd: cmd}
}

func (h *HashStruct[E]) WithCmdable(cmdable Cmdable) *HashStruct[E] {
	b := h.base
	h.cmdable = cmdable
	return &HashStruct[E]{base: b}
}
