package rds

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type HashStruct[T any] struct {
	base
}

func NewHashStruct[T any](key string) *HashStruct[T] {
	return &HashStruct[T]{base: newBase(key)}
}

func (b *HashStruct[T]) Sub(subkey string) *HashStruct[T] {
	return NewHashStruct[T](b.key + ":" + subkey)
}

// 返回该字段是否为新增字段（修改不算新增）
func (b *HashStruct[T]) HSet(field string, value any) (bc BoolCmd) {
	bc.cmd = b.db().HSet(ctx, b.key, field, value)
	return
}

func (b *HashStruct[T]) HMSet(obj *T, exp time.Duration) (ec ErrCmd) {
	cmder := b.db()
	if _, ok := cmder.(redis.Pipeliner); ok {
		ec.cmd = cmder.HSet(ctx, b.key, obj)
		cmder.Expire(ctx, b.key, exp)
	} else {
		pipe := cmder.Pipeline()
		ec.cmd = pipe.HSet(ctx, b.key, obj)
		pipe.Expire(ctx, b.key, exp)
		pipe.Exec(ctx)
	}
	return
}

func (b *HashStruct[T]) HMGet(fields ...string) (sc StructCmd[T]) {
	if len(fields) == 0 {
		sc.cmd = new(redis.MapStringStringCmd)
		return
	}
	sc.fields = fields
	sc.cmd = b.db().HMGet(ctx, b.key, fields...)
	return
}

func (b *HashStruct[T]) HGetAll() (sc StructCmd[T]) {
	sc.cmd = b.db().HGetAll(ctx, b.key)
	return
}

// 返回删除成功数
func (b *HashStruct[T]) HDel(fields ...string) (ic IntCmd) {
	ic.cmd = b.db().HDel(ctx, b.key, fields...)
	return
}

func (b *HashStruct[T]) HExists(field string) (bc BoolCmd) {
	bc.cmd = b.db().HExists(ctx, b.key, field)
	return
}
