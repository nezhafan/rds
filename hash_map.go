package rds

type HashMap[E Ordered] struct {
	base
}

func NewHashMap[E Ordered](key string, ops ...Option) *HashMap[E] {
	h := &HashMap[E]{base: newBase(key, ops...)}
	return h
}

// 返回该字段是否为新增字段（修改不算新增）
func (b *HashMap[E]) HSet(field string, value any) *BoolCmd {
	cmd := b.db().HSet(ctx, b.key, field, value)
	return &BoolCmd{cmd: cmd}
}

func (b *HashMap[E]) HSetNX(field string, value any) *BoolCmd {
	cmd := b.db().HSetNX(ctx, b.key, field, value)
	return &BoolCmd{cmd: cmd}
}

func (b *HashMap[E]) HGet(field string) *StringCmd {
	cmd := b.db().HGet(ctx, b.key, field)
	return &StringCmd{cmd: cmd}
}

func (b *HashMap[E]) HMSet(obj map[string]any) *IntCmd {
	cmd := b.db().HSet(ctx, b.key, obj)
	return &IntCmd{cmd: cmd}
}

func (b *HashMap[E]) HMGet(fields ...string) *MapCmd[E] {
	cmd := b.db().HMGet(ctx, b.key, fields...)
	return &MapCmd[E]{cmd: cmd, fields: fields}
}

func (b *HashMap[E]) HGetAll() *MapCmd[E] {
	cmd := b.db().HGetAll(ctx, b.key)
	return &MapCmd[E]{cmd: cmd}
}

// 返回删除成功数
func (b *HashMap[E]) HDel(fields ...string) *IntCmd {
	cmd := b.db().HDel(ctx, b.key, fields...)
	return &IntCmd{cmd: cmd}
}

// 返回field数量
func (b *HashMap[E]) HLen() *IntCmd {
	cmd := b.db().HLen(ctx, b.key)
	return &IntCmd{cmd: cmd}
}

func (b *HashMap[E]) HExists(field string) *BoolCmd {
	cmd := b.db().HExists(ctx, b.key, field)
	return &BoolCmd{cmd: cmd}
}

func (b *HashMap[E]) HIncrBy(field string, incr int64) *MapCmd[int64] {
	cmd := b.db().HIncrBy(ctx, b.key, field, incr)
	return &MapCmd[int64]{cmd: cmd}
}

func (b *HashMap[E]) HIncrByFloat(field string, incr float64) *MapCmd[float64] {
	cmd := b.db().HIncrByFloat(ctx, b.key, field, incr)
	return &MapCmd[float64]{cmd: cmd}
}
