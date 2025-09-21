package rds

import (
	"cmp"
	"context"
)

type Set[E cmp.Ordered] struct {
	base
}

// Set 去重
func NewSet[E cmp.Ordered](ctx context.Context, key string) *Set[E] {
	return &Set[E]{base: NewBase(ctx, key)}
}

// 添加成员。 返回添加成功数
func (s *Set[E]) SAdd(members ...E) *IntCmd {
	args := toAnys(members)
	cmd := s.db().SAdd(s.ctx, s.key, args...)
	s.done(cmd)
	return &IntCmd{cmd: cmd}
}

// 所有成员
func (s *Set[E]) SMembers() *SliceCmd[E] {
	cmd := s.db().SMembers(s.ctx, s.key)
	s.done(cmd)
	return &SliceCmd[E]{cmd}
}

// 成员数
func (s *Set[E]) SCard() *IntCmd {
	cmd := s.db().SCard(s.ctx, s.key)
	s.done(cmd)
	return &IntCmd{cmd: cmd}
}

// 是否为成员
func (s *Set[E]) SIsMember(member E) *BoolCmd {
	cmd := s.db().SIsMember(s.ctx, s.key, member)
	s.done(cmd)
	return &BoolCmd{cmd: cmd}
}

// 移除成员。 返回移除成功数
func (s *Set[E]) SRem(members ...E) *IntCmd {
	args := toAnys(members)
	cmd := s.db().SRem(s.ctx, s.key, args...)
	s.done(cmd)
	return &IntCmd{cmd: cmd}
}

func (s *Set[E]) WithCmdable(cmdable Cmdable) *Set[E] {
	b := s.base
	b.cmdable = cmdable
	return &Set[E]{base: b}
}
