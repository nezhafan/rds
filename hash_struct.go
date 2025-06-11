package rds

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type HashStruct[E any] struct {
	base
}

func NewHashStruct[E any](key string) *HashStruct[E] {
	return &HashStruct[E]{base: newBase(key)}
}

func (b *HashStruct[E]) Sub(subkey string) *HashStruct[E] {
	return NewHashStruct[E](b.key + ":" + subkey)
}

// 返回该字段是否为新增字段（修改不算新增）
func (b *HashStruct[E]) HSet(field string, value any) *BoolCmd {
	cmd := b.db().HSet(ctx, b.key, field, value)
	return &BoolCmd{cmd: cmd}
}

func (b *HashStruct[E]) HMSet(obj *E, exp time.Duration) *BoolCmd {
	cmder := b.db()
	var cmd redis.Cmder
	if _, ok := cmder.(redis.Pipeliner); ok {
		cmd = cmder.HSet(ctx, b.key, obj)
		cmder.Expire(ctx, b.key, exp)
	} else {
		pipe := cmder.Pipeline()
		cmd = pipe.HSet(ctx, b.key, obj)
		pipe.Expire(ctx, b.key, exp)
		pipe.Exec(ctx)
	}
	return &BoolCmd{cmd: cmd}
}

func (b *HashStruct[E]) HMGet(fields ...string) *StructCmd[E] {
	cmd := b.db().HMGet(ctx, b.key, fields...)
	return &StructCmd[E]{cmd: cmd, fields: fields}
}

func (b *HashStruct[E]) HGetAll() *StructCmd[E] {
	cmd := b.db().HGetAll(ctx, b.key)
	return &StructCmd[E]{cmd: cmd}
}

// 返回删除成功数
func (b *HashStruct[E]) HDel(fields ...string) *IntCmd {
	cmd := b.db().HDel(ctx, b.key, fields...)
	return &IntCmd{cmd: cmd}
}

func (b *HashStruct[E]) HExists(field string) *BoolCmd {
	cmd := b.db().HExists(ctx, b.key, field)
	return &BoolCmd{cmd: cmd}
}
