package rds

import (
	"cmp"
	"context"
)

type List[E cmp.Ordered] struct {
	base
}

func NewList[E cmp.Ordered](ctx context.Context, key string) *List[E] {
	return &List[E]{base: NewBase(ctx, key)}
}

// 左入。 返回list新长度
func (l *List[E]) LPush(vals ...E) *IntCmd {
	args := toAnys(vals)
	cmd := l.db().LPush(l.ctx, l.key, args...)
	l.done(cmd)
	return &IntCmd{cmd: cmd}
}

// 右入。 返回list新长度
func (l *List[E]) RPush(vals ...E) *IntCmd {
	args := toAnys(vals)
	cmd := l.db().RPush(l.ctx, l.key, args...)
	l.done(cmd)
	return &IntCmd{cmd: cmd}
}

// 左出
func (l *List[E]) LPop() *AnyCmd[E] {
	cmd := l.db().LPop(l.ctx, l.key)
	l.done(cmd)
	return &AnyCmd[E]{cmd: cmd}
}

// 右出
func (l *List[E]) RPop() *AnyCmd[E] {
	cmd := l.db().RPop(l.ctx, l.key)
	l.done(cmd)
	return &AnyCmd[E]{cmd: cmd}
}

// 选取
func (l *List[E]) LIndex(index int64) *AnyCmd[E] {
	cmd := l.db().LIndex(l.ctx, l.key, index)
	l.done(cmd)
	return &AnyCmd[E]{cmd: cmd}
}

// 设置
func (l *List[E]) LSet(index int64, val E) *BoolCmd {
	cmd := l.db().LSet(l.ctx, l.key, index, val)
	l.done(cmd)
	return &BoolCmd{cmd: cmd}
}

// 长度
func (l *List[E]) LLen() *IntCmd {
	cmd := l.db().LLen(l.ctx, l.key)
	l.done(cmd)
	return &IntCmd{cmd: cmd}
}

// 遍历 （左闭右闭）
func (l *List[E]) LRange(start, stop int64) *SliceCmd[E] {
	cmd := l.db().LRange(l.ctx, l.key, start, stop)
	l.done(cmd)
	return &SliceCmd[E]{cmd: cmd}
}

// 移除
func (l *List[E]) LRem(count int64, val E) *IntCmd {
	cmd := l.db().LRem(l.ctx, l.key, count, val)
	l.done(cmd)
	return &IntCmd{cmd: cmd}
}

// 截取 （左闭右闭）
func (l *List[E]) LTrim(start, stop int64) *BoolCmd {
	cmd := l.db().LTrim(l.ctx, l.key, start, stop)
	l.done(cmd)
	return &BoolCmd{cmd: cmd}
}

func (l *List[E]) WithCmdable(cmdable Cmdable) *List[E] {
	b := l.base
	b.cmdable = cmdable
	return &List[E]{base: b}
}
