package rds

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Set[E any] struct {
	base
}

// Set 去重
func NewSet[E any](ctx context.Context, key string) *Set[E] {
	return &Set[E]{base: NewBase(ctx, key)}
}

// 添加成员。 返回添加成功数
func (s *Set[E]) SAdd(members ...E) Int64Cmd {
	if len(members) == 0 {
		newInt64Cmd(new(redis.IntCmd))
	}
	args := sliceToAnys(members)
	cmd := s.db().SAdd(s.ctx, s.key, args...)
	s.done(cmd)
	return newInt64Cmd(cmd)
}

// 是否为成员
func (s *Set[E]) SIsMember(member E) BoolCmd {
	cmd := s.db().SIsMember(s.ctx, s.key, member)
	s.done(cmd)
	return newBoolCmd(cmd)
}

// 所有成员
func (s *Set[E]) SMembers() SliceCmd[E] {
	cmd := s.db().SMembers(s.ctx, s.key)
	s.done(cmd)
	return newSliceCmd[E](cmd)
}

// 随机返回count个成员（不删除）
func (s *Set[E]) SRandMember(count int64) SliceCmd[E] {
	cmd := s.db().SRandMemberN(s.ctx, s.key, count)
	s.done(cmd)
	return newSliceCmd[E](cmd)
}

// 成员数
func (s *Set[E]) SCard() Int64Cmd {
	cmd := s.db().SCard(s.ctx, s.key)
	s.done(cmd)
	return newInt64Cmd(cmd)
}

// 随机弹出count个成员（会删除）
func (s *Set[E]) SPop(count int64) SliceCmd[E] {
	cmd := s.db().SPopN(s.ctx, s.key, count)
	s.done(cmd)
	return newSliceCmd[E](cmd)
}

// 移除成员。 返回移除成功数
func (s *Set[E]) SRem(members ...E) Int64Cmd {
	args := sliceToAnys(members)
	cmd := s.db().SRem(s.ctx, s.key, args...)
	s.done(cmd)
	return newInt64Cmd(cmd)
}
