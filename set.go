package rds

import (
	"golang.org/x/exp/constraints"
)

type Set[E constraints.Ordered] struct {
	base
}

// Set 去重
func NewSet[E constraints.Ordered](key string, ops ...Option) Set[E] {
	return Set[E]{base: newBase(key, ops...)}
}

// 添加成员，返回添加成功数
func (s *Set[E]) SAdd(members ...any) (c IntCmd) {
	c.cmd = s.db().SAdd(ctx, s.key, members...)
	return
}

// 获取所有成员
func (s *Set[E]) SMembers() (c SliceCmd[E]) {
	c.cmd = s.db().SMembers(ctx, s.key)
	return
}

// 获取成员数
func (s *Set[E]) SCard() (c IntCmd) {
	c.cmd = s.db().SCard(ctx, s.key)
	return
}

// 是否为成员
func (s *Set[E]) SIsMember(member E) (c BoolCmd) {
	c.cmd = s.db().SIsMember(ctx, s.key, member)
	return
}

// 移除成员，返回移除成功数
func (s *Set[E]) SRem(members ...E) (c IntCmd) {
	args := toAnys(members)
	c.cmd = s.db().SRem(ctx, s.key, args...)
	return
}
