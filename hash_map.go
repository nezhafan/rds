package rds

import (
	"github.com/redis/go-redis/v9"
)

type HashMap struct {
	base
}

func NewHashMap(key string, ops ...Option) *HashMap {
	h := &HashMap{base: newBase(key, ops...)}
	return h
}

// 返回该字段是否为新增字段（修改不算新增）
func (b *HashMap) HSet(field string, value any) (bc BoolCmd) {
	bc.cmd = b.db().HSet(ctx, b.key, field, value)
	return
}

func (b *HashMap) HSetNX(field string, value any) (bc BoolCmd) {
	bc.cmd = b.db().HSetNX(ctx, b.key, field, value)
	return
}

func (b *HashMap) HGet(field string) (sc StringCmd) {
	sc.cmd = b.db().HGet(ctx, b.key, field)
	return
}

func (b *HashMap) HMSet(obj map[string]any) (ec ErrCmd) {
	ec.cmd = b.db().HSet(ctx, b.key, obj)
	return
}

func (b *HashMap) HMGet(fields ...string) (mc MapCmd) {
	if len(fields) == 0 {
		mc.cmd = new(redis.MapStringStringCmd)
		return
	}
	mc.fields = fields
	mc.cmd = b.db().HMGet(ctx, b.key, fields...)
	return
}

func (b *HashMap) HGetAll() (mc MapCmd) {
	mc.cmd = b.db().HGetAll(ctx, b.key)
	return
}

// 返回删除成功数
func (b *HashMap) HDel(fields ...string) (ic IntCmd) {
	ic.cmd = b.db().HDel(ctx, b.key, fields...)
	return
}

// 返回field数量
func (b *HashMap) HLen() (ic IntCmd) {
	ic.cmd = b.db().HLen(ctx, b.key)
	return
}

func (b *HashMap) HExists(field string) (bc BoolCmd) {
	bc.cmd = b.db().HExists(ctx, b.key, field)
	return
}
