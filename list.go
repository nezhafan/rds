package rds

type List[E Ordered] struct {
	base
}

func NewList[E Ordered](key string, ops ...Option) *List[E] {
	return &List[E]{base: newBase(key, ops...)}
}

// 左入。 返回list新长度
func (b *List[E]) LPush(vals ...E) *IntCmd {
	args := toAnys(vals)
	cmd := b.db().LPush(ctx, b.key, args...)
	b.done(cmd)
	return &IntCmd{cmd: cmd}
}

// 右入。 返回list新长度
func (b *List[E]) RPush(vals ...E) *IntCmd {
	args := toAnys(vals)
	cmd := b.db().RPush(ctx, b.key, args...)
	b.done(cmd)
	return &IntCmd{cmd: cmd}
}

// 左出
func (b *List[E]) LPop() *AnyCmd[E] {
	cmd := b.db().LPop(ctx, b.key)
	b.done(cmd)
	return &AnyCmd[E]{cmd: cmd}
}

// 右出
func (b *List[E]) RPop() *AnyCmd[E] {
	cmd := b.db().RPop(ctx, b.key)
	b.done(cmd)
	return &AnyCmd[E]{cmd: cmd}
}

// 选取
func (b *List[E]) LIndex(index int64) *AnyCmd[E] {
	cmd := b.db().LIndex(ctx, b.key, index)
	b.done(cmd)
	return &AnyCmd[E]{cmd: cmd}
}

// 设置
func (b *List[E]) LSet(index int64, val E) *BoolCmd {
	cmd := b.db().LSet(ctx, b.key, index, val)
	b.done(cmd)
	return &BoolCmd{cmd: cmd}
}

// 长度
func (b *List[E]) LLen() *IntCmd {
	cmd := b.db().LLen(ctx, b.key)
	b.done(cmd)
	return &IntCmd{cmd: cmd}
}

// 遍历 （左闭右闭）
func (b *List[E]) LRange(start, stop int64) *SliceCmd[E] {
	cmd := b.db().LRange(ctx, b.key, start, stop)
	b.done(cmd)
	return &SliceCmd[E]{cmd: cmd}
}

// 移除
func (b *List[E]) LRem(count int64, val E) *IntCmd {
	cmd := b.db().LRem(ctx, b.key, count, val)
	b.done(cmd)
	return &IntCmd{cmd: cmd}
}

// 截取 （左闭右闭）
func (b *List[E]) LTrim(start, stop int64) *BoolCmd {
	cmd := b.db().LTrim(ctx, b.key, start, stop)
	b.done(cmd)
	return &BoolCmd{cmd: cmd}
}
