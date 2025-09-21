package rds

import (
	"cmp"
	"context"
)

type HashMap[E cmp.Ordered] struct {
	base
}

func NewHashMap[E cmp.Ordered](ctx context.Context, key string) *HashMap[E] {
	h := &HashMap[E]{base: NewBase(ctx, key)}
	return h
}

// 返回该字段是否为新增字段（修改不算新增）
func (h *HashMap[E]) HSet(field string, value E) *BoolCmd {
	cmd := h.db().HSet(h.ctx, h.key, field, value)
	h.done(cmd)
	return &BoolCmd{cmd: cmd}
}

func (h *HashMap[E]) HSetNX(field string, value E) *BoolCmd {
	cmd := h.db().HSetNX(h.ctx, h.key, field, value)
	h.done(cmd)
	return &BoolCmd{cmd: cmd}
}

func (h *HashMap[E]) HGet(field string) *StringCmd[E] {
	cmd := h.db().HGet(h.ctx, h.key, field)
	h.done(cmd)
	return &StringCmd[E]{cmd: cmd}
}

func (h *HashMap[E]) HMSet(obj map[string]E) *BoolCmd {
	args := make([]any, 0, len(obj)*2)
	for k, v := range obj {
		args = append(args, k, v)
	}
	cmd := h.db().HMSet(h.ctx, h.key, args...)
	h.done(cmd)
	return &BoolCmd{cmd: cmd}
}

func (h *HashMap[E]) HMGet(fields ...string) *MapCmd[E] {
	cmd := h.db().HMGet(h.ctx, h.key, fields...)
	h.done(cmd)
	return &MapCmd[E]{cmd: cmd, fields: fields}
}

func (h *HashMap[E]) HGetAll() *MapCmd[E] {
	cmd := h.db().HGetAll(h.ctx, h.key)
	h.done(cmd)
	return &MapCmd[E]{cmd: cmd}
}

// 返回删除成功数
func (h *HashMap[E]) HDel(fields ...string) *IntCmd {
	cmd := h.db().HDel(h.ctx, h.key, fields...)
	h.done(cmd)
	return &IntCmd{cmd: cmd}
}

// 返回field数量
func (h *HashMap[E]) HLen() *IntCmd {
	cmd := h.db().HLen(h.ctx, h.key)
	h.done(cmd)
	return &IntCmd{cmd: cmd}
}

func (h *HashMap[E]) HExists(field string) *BoolCmd {
	cmd := h.db().HExists(h.ctx, h.key, field)
	h.done(cmd)
	return &BoolCmd{cmd: cmd}
}

func (h *HashMap[E]) HIncrByInt(field string, incr int64) *StringCmd[E] {
	cmd := h.db().HIncrBy(h.ctx, h.key, field, incr)
	h.done(cmd)
	return &StringCmd[E]{cmd: cmd}
}

func (h *HashMap[E]) HIncrByFloat(field string, incr float64) *StringCmd[E] {
	cmd := h.db().HIncrByFloat(h.ctx, h.key, field, incr)
	h.done(cmd)
	return &StringCmd[E]{cmd: cmd}
}

func (h *HashMap[E]) WithCmdable(cmdable Cmdable) *HashMap[E] {
	h.cmdable = cmdable
	return h
}
