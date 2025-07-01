package rds

type HashMap[E Ordered] struct {
	base
}

func NewHashMap[E Ordered](key string, ops ...Option) *HashMap[E] {
	h := &HashMap[E]{base: newBase(key, ops...)}
	return h
}

// 返回该字段是否为新增字段（修改不算新增）
func (b *HashMap[E]) HSet(field string, value E) *BoolCmd {
	cmd := b.db().HSet(b.ctx, b.key, field, value)
	b.done(cmd)
	return &BoolCmd{cmd: cmd}
}

func (b *HashMap[E]) HSetNX(field string, value E) *BoolCmd {
	cmd := b.db().HSetNX(b.ctx, b.key, field, value)
	b.done(cmd)
	return &BoolCmd{cmd: cmd}
}

func (b *HashMap[E]) HGet(field string) *StringCmd[E] {
	cmd := b.db().HGet(b.ctx, b.key, field)
	b.done(cmd)
	return &StringCmd[E]{cmd: cmd}
}

func (b *HashMap[E]) HMSet(obj map[string]E) *BoolCmd {
	args := make([]any, 0, len(obj)*2)
	for k, v := range obj {
		args = append(args, k, v)
	}
	cmd := b.db().HMSet(b.ctx, b.key, args...)
	b.done(cmd)
	return &BoolCmd{cmd: cmd}
}

func (b *HashMap[E]) HMGet(fields ...string) *MapCmd[E] {
	cmd := b.db().HMGet(b.ctx, b.key, fields...)
	b.done(cmd)
	return &MapCmd[E]{cmd: cmd, fields: fields}
}

func (b *HashMap[E]) HGetAll() *MapCmd[E] {
	cmd := b.db().HGetAll(b.ctx, b.key)
	b.done(cmd)
	return &MapCmd[E]{cmd: cmd}
}

// 返回删除成功数
func (b *HashMap[E]) HDel(fields ...string) *IntCmd {
	cmd := b.db().HDel(b.ctx, b.key, fields...)
	b.done(cmd)
	return &IntCmd{cmd: cmd}
}

// 返回field数量
func (b *HashMap[E]) HLen() *IntCmd {
	cmd := b.db().HLen(b.ctx, b.key)
	b.done(cmd)
	return &IntCmd{cmd: cmd}
}

func (b *HashMap[E]) HExists(field string) *BoolCmd {
	cmd := b.db().HExists(b.ctx, b.key, field)
	b.done(cmd)
	return &BoolCmd{cmd: cmd}
}

func (b *HashMap[E]) HIncrByInt(field string, incr int64) *StringCmd[E] {
	cmd := b.db().HIncrBy(b.ctx, b.key, field, incr)
	b.done(cmd)
	return &StringCmd[E]{cmd: cmd}
}

func (b *HashMap[E]) HIncrByFloat(field string, incr float64) *StringCmd[E] {
	cmd := b.db().HIncrByFloat(b.ctx, b.key, field, incr)
	b.done(cmd)
	return &StringCmd[E]{cmd: cmd}
}
