package rds

type HashMap[E string | int64 | float64] struct {
	base
}

func NewHashMap[E string | int64 | float64](key string, ops ...Option) *HashMap[E] {
	h := &HashMap[E]{base: newBase(key, ops...)}
	return h
}

// 返回该字段是否为新增字段（修改不算新增）
func (b *HashMap[E]) HSet(field string, value E) *BoolCmd {
	cmd := b.db().HSet(ctx, b.key, field, value)
	b.done(cmd)
	return &BoolCmd{cmd: cmd}
}

func (b *HashMap[E]) HSetNX(field string, value E) *BoolCmd {
	cmd := b.db().HSetNX(ctx, b.key, field, value)
	b.done(cmd)
	return &BoolCmd{cmd: cmd}
}

func (b *HashMap[E]) HGet(field string) *StringCmd[E] {
	cmd := b.db().HGet(ctx, b.key, field)
	b.done(cmd)
	return &StringCmd[E]{cmd: cmd}
}

// 不支持map[string]E ，所以自己确保any的值是E类型
func (b *HashMap[E]) HMSet(obj map[string]any) *IntCmd {
	cmd := b.db().HSet(ctx, b.key, obj)
	b.done(cmd)
	return &IntCmd{cmd: cmd}
}

func (b *HashMap[E]) HMGet(fields ...string) *MapCmd[E] {
	cmd := b.db().HMGet(ctx, b.key, fields...)
	b.done(cmd)
	return &MapCmd[E]{cmd: cmd, fields: fields}
}

func (b *HashMap[E]) HGetAll() *MapCmd[E] {
	cmd := b.db().HGetAll(ctx, b.key)
	b.done(cmd)
	return &MapCmd[E]{cmd: cmd}
}

// 返回删除成功数
func (b *HashMap[E]) HDel(fields ...string) *IntCmd {
	cmd := b.db().HDel(ctx, b.key, fields...)
	b.done(cmd)
	return &IntCmd{cmd: cmd}
}

// 返回field数量
func (b *HashMap[E]) HLen() *IntCmd {
	cmd := b.db().HLen(ctx, b.key)
	b.done(cmd)
	return &IntCmd{cmd: cmd}
}

func (b *HashMap[E]) HExists(field string) *BoolCmd {
	cmd := b.db().HExists(ctx, b.key, field)
	b.done(cmd)
	return &BoolCmd{cmd: cmd}
}

func (b *HashMap[E]) HIncrByInt(field string, incr int64) *StringCmd[E] {
	cmd := b.db().HIncrBy(ctx, b.key, field, incr)
	b.done(cmd)
	return &StringCmd[E]{cmd: cmd}
}

func (b *HashMap[E]) HIncrByFloat(field string, incr float64) *StringCmd[E] {
	cmd := b.db().HIncrByFloat(ctx, b.key, field, incr)
	b.done(cmd)
	return &StringCmd[E]{cmd: cmd}
}
