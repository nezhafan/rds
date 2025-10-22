package rds

import (
	"context"
)

type List[E any] struct {
	base
}

func NewList[E any](ctx context.Context, key string) *List[E] {
	return &List[E]{base: NewBase(ctx, key)}
}

// 左入。 返回list新长度
func (l *List[E]) LPush(vals ...E) Int64Cmd {
	if len(vals) == 0 {
		return l.LLen()
	}
	args := sliceToAnys(vals)
	cmd := l.db().LPush(l.ctx, l.key, args...)
	l.done(cmd)
	return newInt64Cmd(cmd)
}

// 右入。 返回list新长度
func (l *List[E]) RPush(vals ...E) Int64Cmd {
	if len(vals) == 0 {
		return l.LLen()
	}
	args := sliceToAnys(vals)
	cmd := l.db().RPush(l.ctx, l.key, args...)
	l.done(cmd)
	return newInt64Cmd(cmd)
}

// 左出
func (l *List[E]) LPop() AnyCmd[E] {
	cmd := l.db().LPop(l.ctx, l.key)
	l.done(cmd)
	return newAnyCmd[E](cmd)
}

// 右出
func (l *List[E]) RPop() AnyCmd[E] {
	cmd := l.db().RPop(l.ctx, l.key)
	l.done(cmd)
	return newAnyCmd[E](cmd)
}

// 设置指定索引位置的元素值。 key必须存在，索引超出范围会报错
func (l *List[E]) LSet(index int64, val E) BoolCmd {
	cmd := l.db().LSet(l.ctx, l.key, index, val)
	l.done(cmd)
	return newBoolCmd(cmd)
}

// 选取
func (l *List[E]) LIndex(index int64) AnyCmd[E] {
	cmd := l.db().LIndex(l.ctx, l.key, index)
	l.done(cmd)
	return newAnyCmd[E](cmd)
}

// 长度
func (l *List[E]) LLen() Int64Cmd {
	cmd := l.db().LLen(l.ctx, l.key)
	l.done(cmd)
	return newInt64Cmd(cmd)
}

// 遍历 [左闭右闭]
func (l *List[E]) LRange(start, stop int64) SliceCmd[E] {
	cmd := l.db().LRange(l.ctx, l.key, start, stop)
	l.done(cmd)
	return newSliceCmd[E](cmd)
}

// 移除 指定值的元素。 count>0 从表头开始移除； count<0 从表尾开始移除； count=0 移除所有
func (l *List[E]) LRem(count int64, val E) Int64Cmd {
	cmd := l.db().LRem(l.ctx, l.key, count, val)
	l.done(cmd)
	return newInt64Cmd(cmd)
}

// 截取保留段，其余删除 [左闭右闭]
func (l *List[E]) LTrim(start, stop int64) BoolCmd {
	cmd := l.db().LTrim(l.ctx, l.key, start, stop)
	l.done(cmd)
	return newBoolCmd(cmd)
}
