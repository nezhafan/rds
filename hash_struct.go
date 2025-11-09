package rds

import (
	"context"
	"time"
)

type HashStruct[E any] struct {
	base
}

// 泛型请传入结构体而非结构体指针，例：NewHashStruct[User]
func NewHashStruct[E any](ctx context.Context, key string) *HashStruct[E] {
	return &HashStruct[E]{base: newBase(ctx, key)}
}

// 缓存对象，强制设定过期时间。 ⚠️ 注意字段必须定义 `redis:"xx"` 标签才会存储，无redis标签会报错
func (h *HashStruct[E]) HSet(obj *E, exp time.Duration) BoolCmd {
	// 处理缓存nil的情况
	var values []any
	if obj == nil {
		values = []any{emptyField, true}
	} else {
		values = struct2Anys(obj)
	}
	// 使用管道同时设置过期时间
	pipe := h.db().Pipeline()
	cmd1 := pipe.HSet(h.ctx, h.key, values)
	cmd2 := pipe.Expire(h.ctx, h.key, exp)
	_, err := pipe.Exec(h.ctx)
	cmd1.SetErr(err)
	h.done(cmd1)
	h.done(cmd2)
	return newBoolCmd(cmd1)
}

func (h *HashStruct[E]) HGet(field string) StringCmdR {
	cmd := h.db().HGet(h.ctx, h.key, field)
	h.done(cmd)
	return newStringCmdR(cmd)
}

func (h *HashStruct[E]) HMGet(fields ...string) StructCmd[E] {
	if len(fields) == 0 {
		return h.HGetAll()
	}
	cmd := h.db().HMGet(h.ctx, h.key, fields...)
	h.done(cmd)
	return newStructCmd[E](cmd, fields)
}

func (h *HashStruct[E]) HGetAll() StructCmd[E] {
	cmd := h.db().HGetAll(h.ctx, h.key)
	h.done(cmd)
	return newStructCmd[E](cmd, nil)
}
