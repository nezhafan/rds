package rds

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type HashStruct[E any] struct {
	base
}

// 泛型请传入结构体而非结构体指针，例：NewHashStruct[User]
func NewHashStruct[E any](ctx context.Context, key string) *HashStruct[E] {
	return &HashStruct[E]{base: NewBase(ctx, key)}
}

func (h *HashStruct[E]) SubKey(ctx context.Context, subkey string) *HashStruct[E] {
	return NewHashStruct[E](ctx, h.key+":"+subkey)
}

// 对象替换，强制设定过期时间。 //TODO
func (h *HashStruct[E]) HSetAll(obj *E, exp time.Duration) BoolCmd {
	// 处理缓存nil的情况
	var values []any
	if obj == nil {
		values = []any{emptyField, true}
	} else {
		values = structToAnys(obj)
	}
	// 使用管道。 需要先删除，防止出现字段不对齐的情况
	pipe := h.db().Pipeline()
	pipe.Del(h.ctx, h.key)
	cmd := pipe.HSet(h.ctx, h.key, values)
	pipe.Expire(h.ctx, h.key, exp)
	_, err := pipe.Exec(h.ctx)
	cmd.SetErr(err)
	h.done(cmd)
	return newBoolCmd(cmd)
}

// hset 支持一次性设置多个字段
// hmset 已经从4.0视为被弃用 https://redis.io/docs/latest/commands/hmset/
func (h *HashStruct[E]) HSet(update map[string]any, exp time.Duration) Int64Cmd {
	if len(update) == 0 {
		return newInt64Cmd(new(redis.IntCmd))
	}
	values := mapToAnys(update)
	pipe := h.db().Pipeline()
	cmd := pipe.HSet(h.ctx, h.key, values)
	pipe.Expire(h.ctx, h.key, exp)
	cmds, err := pipe.Exec(h.ctx)
	cmd.SetErr(err)
	for _, c := range cmds {
		h.done(c)
	}
	return newInt64Cmd(cmd)
}

func (h *HashStruct[E]) HGet(field string) StringCmdR {
	cmd := h.db().HGet(h.ctx, h.key, field)
	h.done(cmd)
	return newStringCmdR(cmd)
}

func (h *HashStruct[E]) HMGet(fields ...string) StructCmd[E] {
	cmd := h.db().HMGet(h.ctx, h.key, fields...)
	h.done(cmd)
	return newStructCmd[E](cmd, fields)
}

func (h *HashStruct[E]) HGetAll() StructCmd[E] {
	cmd := h.db().HGetAll(h.ctx, h.key)
	h.done(cmd)
	return newStructCmd[E](cmd, nil)
}

// 注意 hincrby 去增长一个浮点数字段会报错
func (h *HashStruct[E]) HIncrBy(field string, incr int64) Int64Cmd {
	cmd := h.db().HIncrBy(h.ctx, h.key, field, incr)
	h.done(cmd)
	return newInt64Cmd(cmd)
}

// 注意 hincrbyfloat 可能出现精度问题
func (h *HashStruct[E]) HIncrByFloat(field string, incr float64) Float64Cmd {
	cmd := h.db().HIncrByFloat(h.ctx, h.key, field, incr)
	h.done(cmd)
	return newFloat64Cmd(cmd)
}

// 返回删除成功数
func (h *HashStruct[E]) HDel(fields ...string) Int64Cmd {
	cmd := h.db().HDel(h.ctx, h.key, fields...)
	h.done(cmd)
	return newInt64Cmd(cmd)
}

func (h *HashStruct[E]) HExists(field string) BoolCmd {
	cmd := h.db().HExists(h.ctx, h.key, field)
	h.done(cmd)
	return newBoolCmd(cmd)
}
