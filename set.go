package rds

import (
	"context"

	"golang.org/x/exp/constraints"
)

type Set[E constraints.Ordered] struct {
	base
}

// Set 去重
func NewSet[E constraints.Ordered](ctx context.Context, key string) Set[E] {
	return Set[E]{base: newBase(ctx, key)}
}

// 添加成员，返回添加成功数
func (s *Set[E]) SAdd(members ...any) (success int, ok bool) {
	cmd := DB().SAdd(s.ctx, s.key, members...)
	s.done(cmd)
	success, ok = int(cmd.Val()), cmd.Err() == nil
	return
}

// 获取所有成员
func (s *Set[E]) SMembers() (members []E, ok bool) {
	cmd := DB().SMembers(s.ctx, s.key)
	s.done(cmd)
	members, ok = stringsToSlice[E](cmd.Val()), cmd.Err() == nil
	return
}

// 获取成员数
func (s *Set[E]) SCard() int {
	cmd := DB().SCard(s.ctx, s.key)
	s.done(cmd)
	return int(cmd.Val())
}

// 是否为成员
func (s *Set[E]) SIsMember(member E) (exists bool, err error) {
	cmd := DB().SIsMember(s.ctx, s.key, member)
	s.done(cmd)
	return cmd.Result()
}

// 移除成员，返回移除成功数
func (s *Set[E]) SRem(members ...E) (n int64, err error) {
	args := toAnys(members)
	cmd := DB().SRem(s.ctx, s.key, args...)
	s.done(cmd)
	return cmd.Result()
}
