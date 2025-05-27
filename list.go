package rds

import (
	"context"
	"fmt"

	"golang.org/x/exp/constraints"
)

type List[E constraints.Ordered] struct {
	base
}

func NewList[E constraints.Ordered](ctx context.Context, key string) List[E] {
	return List[E]{base: newBase(ctx, key)}
}

// 左入
func (l *List[E]) LPush(vals ...E) (int64, error) {
	args := toAnys(vals)
	cmd := DB().LPush(l.ctx, l.key, args...)
	l.done(cmd)
	return cmd.Result()
}

// 右入
func (l *List[E]) RPush(vals ...E) (int64, error) {
	args := toAnys(vals)
	cmd := DB().RPush(l.ctx, l.key, args...)
	l.done(cmd)
	return cmd.Result()
}

// 左出
func (l *List[E]) LPop() (E, error) {
	cmd := DB().LPop(l.ctx, l.key)
	l.done(cmd)
	v, err := cmd.Result()
	return stringTo[E](v), err
}

// 右出
func (l *List[E]) RPop() (E, error) {
	cmd := DB().RPop(l.ctx, l.key)
	l.done(cmd)
	v, err := cmd.Result()
	return stringTo[E](v), err
}

// 遍历
func (l *List[E]) LRange(start, stop int) []E {
	cmd := DB().LRange(l.ctx, l.key, int64(start), int64(stop))
	l.done(cmd)
	vs := cmd.Val()
	fmt.Println("=====", vs)
	return stringsToSlice[E](vs)
}

// 长度
func (l *List[E]) LLen() int {
	cmd := DB().LLen(l.ctx, l.key)
	l.done(cmd)
	return int(cmd.Val())
}

// 从左开始移除count个
func (l *List[E]) LRem(value any, count int64) (int64, error) {
	cmd := DB().LRem(l.ctx, l.key, count, value)
	l.done(cmd)
	return cmd.Result()
}

// 从右开始移除count个
func (l *List[E]) RRem(value any, count int64) (int64, error) {
	cmd := DB().LRem(l.ctx, l.key, -count, value)
	l.done(cmd)
	return cmd.Result()
}

// 移除，返回被移除数量
func (l *List[E]) Rem(value any) (int64, error) {
	cmd := DB().LRem(l.ctx, l.key, 0, value)
	l.done(cmd)
	return cmd.Result()
}
