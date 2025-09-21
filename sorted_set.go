package rds

import (
	"cmp"
	"context"

	"github.com/redis/go-redis/v9"
)

const (
	Inf    = "+inf"
	NegInf = "-inf"
)

type SortedSet[E cmp.Ordered] struct {
	base
}

type Z[E cmp.Ordered] struct {
	Score  float64
	Member E
}

func NewSortedSet[E cmp.Ordered](ctx context.Context, key string) *SortedSet[E] {
	return &SortedSet[E]{base: NewBase(ctx, key)}
}

// 添加
func (s *SortedSet[E]) ZAdd(zs map[E]float64) *IntCmd {
	args := make([]any, 0, len(zs)*2+3)
	args = append(args, "zadd", s.key, "ch")
	for member, score := range zs {
		args = append(args, score, member)
	}
	cmd := s.db().Do(s.ctx, args...)
	s.done(cmd)
	return &IntCmd{cmd: cmd}
}

// 所有成员数
func (s *SortedSet[E]) ZCard() *IntCmd {
	cmd := s.db().ZCard(s.ctx, s.key)
	s.done(cmd)
	return &IntCmd{cmd: cmd}
}

// 积分区间内的成员数
func (s *SortedSet[E]) ZCount(min, max float64) *IntCmd {
	minS := toString(min)
	maxS := toString(max)
	cmd := s.db().ZCount(s.ctx, s.key, minS, maxS)
	s.done(cmd)
	return &IntCmd{cmd: cmd}
}

// 获取分数
func (s *SortedSet[E]) ZScore(member E) *FloatCmd {
	cmd := s.db().ZScore(s.ctx, s.key, toString(member))
	s.done(cmd)
	return &FloatCmd{cmd: cmd}
}

// 增加分数
func (s *SortedSet[E]) ZIncrBy(member E, incr float64) *FloatCmd {
	cmd := s.db().ZIncrBy(s.ctx, s.key, float64(incr), toString(member))
	s.done(cmd)
	return &FloatCmd{cmd: cmd}
}

// 按照积分获取：成员
func (s *SortedSet[E]) ZMembersByScore(min, max float64, offset, limit int64, asc bool) *SliceCmd[E] {
	by := &redis.ZRangeBy{
		Min:    toString(min),
		Max:    toString(max),
		Offset: offset,
		Count:  limit,
	}
	var cmd *redis.StringSliceCmd
	if asc {
		cmd = s.db().ZRangeByScore(s.ctx, s.key, by)
	} else {
		cmd = s.db().ZRevRangeByScore(s.ctx, s.key, by)
	}
	s.done(cmd)
	return &SliceCmd[E]{cmd: cmd}
}

// 按照排名获取：成员
func (s *SortedSet[E]) ZMembersByRank(start, stop int64, asc bool) *SliceCmd[E] {
	var cmd *redis.StringSliceCmd
	if asc {
		cmd = s.db().ZRange(s.ctx, s.key, start, stop)
	} else {
		cmd = s.db().ZRevRange(s.ctx, s.key, start, stop)
	}
	s.done(cmd)
	return &SliceCmd[E]{cmd: cmd}
}

// 按照积分获取：成员和积分
func (s *SortedSet[E]) ZRangeByScore(min, max float64, offset, limit int64, asc bool) *ZSliceCmd[E] {
	by := &redis.ZRangeBy{
		Min:    toString(min),
		Max:    toString(max),
		Offset: offset,
		Count:  limit,
	}
	var cmd *redis.ZSliceCmd
	if asc {
		cmd = s.db().ZRangeByScoreWithScores(s.ctx, s.key, by)
	} else {
		cmd = s.db().ZRevRangeByScoreWithScores(s.ctx, s.key, by)
	}
	s.done(cmd)
	return &ZSliceCmd[E]{cmd: cmd}
}

// 按照排名获取：成员和积分
func (s *SortedSet[E]) ZRangeByRank(start, stop int64, asc bool) *ZSliceCmd[E] {
	var cmd *redis.ZSliceCmd
	if asc {
		cmd = s.db().ZRangeWithScores(s.ctx, s.key, start, stop)
	} else {
		cmd = s.db().ZRevRangeWithScores(s.ctx, s.key, start, stop)
	}
	s.done(cmd)
	return &ZSliceCmd[E]{cmd: cmd}
}

// 移除成员
func (s *SortedSet[E]) ZRem(members ...E) *IntCmd {
	args := toAnys(members)
	cmd := s.db().ZRem(s.ctx, s.key, args...)
	s.done(cmd)
	return &IntCmd{cmd: cmd}
}

// 移除积分区间内的成员
func (s *SortedSet[E]) ZRemByScore(min, max float64) *IntCmd {
	minS := toString(min)
	maxS := toString(max)
	cmd := s.db().ZRemRangeByScore(s.ctx, s.key, minS, maxS)
	s.done(cmd)
	return &IntCmd{cmd: cmd}
}

// 移除排名区间内的成员
func (s *SortedSet[E]) ZRemByRank(start, stop int64) *IntCmd {
	cmd := s.db().ZRemRangeByRank(s.ctx, s.key, start, stop)
	s.done(cmd)
	return &IntCmd{cmd: cmd}
}

func (s *SortedSet[E]) WithCmdable(cmdable Cmdable) *SortedSet[E] {
	b := s.base
	b.cmdable = cmdable
	return &SortedSet[E]{base: b}
}
