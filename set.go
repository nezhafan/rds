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
	args := slice2Anys(members)
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

/*
游标扫描获取，自动管理游标位置
  - match 使用*匹配部分，可传空查找全部。
  - count 并不一定返回要求的个数，在数量少的时候没作用可能返回全部，数量多的时候会返回近似个元素。
  - fn 处理函数，随时return false停止继续获取
*/
func (s *Set[E]) SScan(match string, count int64, fn func(vals []E) error) (err error) {
	var cursor uint64
	var vals []string
	var result []E
	for {
		cmd := s.db().SScan(s.ctx, s.key, cursor, match, count)
		s.done(cmd)

		vals, cursor, err = cmd.Result()
		if err != nil {
			return
		}

		for i := range vals {
			result = append(result, string2E[E](vals[i]))
		}

		if len(result) > 0 {
			if err = fn(result); err != nil {
				return
			}
		}

		if cursor == 0 {
			return nil
		}
		result = result[:0]
	}
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
	args := slice2Anys(members)
	cmd := s.db().SRem(s.ctx, s.key, args...)
	s.done(cmd)
	return newInt64Cmd(cmd)
}
