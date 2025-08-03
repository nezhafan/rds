package rds

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type HashStruct[E any] struct {
	base
}

func NewHashStruct[E any](key string, ops ...Option) *HashStruct[E] {
	return &HashStruct[E]{base: newBase(key, ops...)}
}

func (b *HashStruct[E]) SubKey(subkey string) *HashStruct[E] {
	return NewHashStruct[E](b.key + ":" + subkey)
}

func (b *HashStruct[E]) SubID(subid string) *HashStruct[E] {
	return NewHashStruct[E](b.key + ":" + subid)
}

// 返回该字段是否为新增字段（修改不算新增）
func (b *HashStruct[E]) HSet(field string, value any) *BoolCmd {
	cmd := b.db().HSet(b.ctx, b.key, field, value)
	b.done(cmd)
	return &BoolCmd{cmd: cmd}
}

// 强制设置过期时间
func (b *HashStruct[E]) HMSet(obj *E, exp time.Duration) *BoolCmd {
	cmder := b.db()
	var cmd redis.Cmder
	if _, ok := cmder.(redis.Pipeliner); ok {
		cmd = cmder.HSet(b.ctx, b.key, obj)
		cmder.Expire(b.ctx, b.key, exp)
	} else {
		pipe := cmder.Pipeline()
		cmd = pipe.HSet(b.ctx, b.key, obj)
		pipe.Expire(b.ctx, b.key, exp)
		pipe.Exec(b.ctx)
	}
	b.done(cmd)
	return &BoolCmd{cmd: cmd}
}

func (b *HashStruct[E]) HMGet(fields ...string) *StructCmd[E] {
	if len(fields) == 0 {
		return b.HGetAll()
	}
	cmd := b.db().HMGet(b.ctx, b.key, fields...)
	b.done(cmd)
	return &StructCmd[E]{cmd: cmd, fields: fields}
}

func (b *HashStruct[E]) HGetAll() *StructCmd[E] {
	cmd := b.db().HGetAll(b.ctx, b.key)
	b.done(cmd)
	return &StructCmd[E]{cmd: cmd}
}

// 返回删除成功数
func (b *HashStruct[E]) HDel(fields ...string) *IntCmd {
	cmd := b.db().HDel(b.ctx, b.key, fields...)
	b.done(cmd)
	return &IntCmd{cmd: cmd}
}

func (b *HashStruct[E]) HExists(field string) *BoolCmd {
	cmd := b.db().HExists(b.ctx, b.key, field)
	b.done(cmd)
	return &BoolCmd{cmd: cmd}
}
