package rds

import "fmt"

type Set[E Ordered] struct {
	base
}

// Set 去重
func NewSet[E Ordered](key string, ops ...Option) *Set[E] {
	return &Set[E]{base: newBase(key, ops...)}
}

// 添加成员，返回添加成功数
func (s *Set[E]) SAdd(members ...E) *IntCmd {
	args := toAnys(members)
	cmd := s.db().SAdd(ctx, s.key, args...)
	return &IntCmd{cmd: cmd}
}

// 获取所有成员
func (s *Set[E]) SMembers() *SliceCmd[E] {
	cmd := s.db().SMembers(ctx, s.key)
	fmt.Println(cmd.Val())
	return &SliceCmd[E]{cmd}
}

// 获取成员数
func (s *Set[E]) SCard() *IntCmd {
	cmd := s.db().SCard(ctx, s.key)
	return &IntCmd{cmd: cmd}
}

// 是否为成员
func (s *Set[E]) SIsMember(member E) *BoolCmd {
	cmd := s.db().SIsMember(ctx, s.key, member)
	return &BoolCmd{cmd: cmd}
}

// 移除成员，返回移除成功数
func (s *Set[E]) SRem(members ...E) *IntCmd {
	args := toAnys(members)
	cmd := s.db().SRem(ctx, s.key, args...)
	return &IntCmd{cmd: cmd}
}
