package rds

import (
	"cmp"
	"context"

	"github.com/redis/go-redis/v9"
)

type SortedSet[E cmp.Ordered] struct {
	base
}

type Z[E cmp.Ordered] struct {
	Score  float64
	Member E
}

// 有序集合
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

// 增加分数
func (s *SortedSet[E]) ZIncrBy(member E, incr float64) *FloatCmd {
	cmd := s.db().ZIncrBy(s.ctx, s.key, float64(incr), toString(member))
	s.done(cmd)
	return &FloatCmd{cmd: cmd}
}

// 所有成员数
func (s *SortedSet[E]) ZCard() *IntCmd {
	cmd := s.db().ZCard(s.ctx, s.key)
	s.done(cmd)
	return &IntCmd{cmd: cmd}
}

// 积分区间内的成员数
func (s *SortedSet[E]) ZCountByScore(minScore, maxScore float64) *IntCmd {
	minS := toString(minScore)
	maxS := toString(maxScore)
	cmd := s.db().ZCount(s.ctx, s.key, minS, maxS)
	s.done(cmd)
	return &IntCmd{cmd: cmd}
}

// 获取分数。 member => score
func (s *SortedSet[E]) ZScore(member E) *FloatCmd {
	cmd := s.db().ZScore(s.ctx, s.key, toString(member))
	s.done(cmd)
	return &FloatCmd{cmd: cmd}
}

// 获取位置。 (从0开始) member => index
func (s *SortedSet[E]) ZIndex(member E, scoreDesc bool) *IntCmd {
	var cmd *redis.IntCmd
	if scoreDesc {
		cmd = s.db().ZRevRank(s.ctx, s.key, toString(member))
	} else {
		cmd = s.db().ZRank(s.ctx, s.key, toString(member))
	}
	s.done(cmd)
	return &IntCmd{cmd: cmd}
}

// 积分范围内的成员。 [左闭右闭] [minScore, maxScore] => []{member1, member2, ...}
// 若不需要偏移offset和数量限制count，参数传0 即可
func (s *SortedSet[E]) ZMembersByScore(minScore, maxScore float64, desc bool, offset, limit int64) *SliceCmd[E] {
	by := &redis.ZRangeBy{
		Min:    toString(minScore),
		Max:    toString(maxScore),
		Offset: offset,
		Count:  limit,
	}
	var cmd *redis.StringSliceCmd
	if desc {
		cmd = s.db().ZRevRangeByScore(s.ctx, s.key, by)
	} else {
		cmd = s.db().ZRangeByScore(s.ctx, s.key, by)
	}
	s.done(cmd)
	return &SliceCmd[E]{cmd: cmd}
}

// 排序范围内的成员。 (从0开始) [左闭右闭] [start, stop] => []{member1, member2, ...}
func (s *SortedSet[E]) ZMembersByIndex(start, stop int64, scoreDesc bool) *SliceCmd[E] {
	var cmd *redis.StringSliceCmd
	if scoreDesc {
		cmd = s.db().ZRevRange(s.ctx, s.key, start, stop)
	} else {
		cmd = s.db().ZRange(s.ctx, s.key, start, stop)
	}
	s.done(cmd)
	return &SliceCmd[E]{cmd: cmd}
}

// 积分范围内的成员和其积分。  [左闭右闭] [minScore, maxScore] => []{(member1,score1), (member2,score2), ...}
// 若不需要偏移offset和数量限制count，参数传0 即可
func (s *SortedSet[E]) ZRangeByScore(minScore, maxScore float64, desc bool, offset, limit int64) *ZSliceCmd[E] {
	by := &redis.ZRangeBy{
		Min:    toString(minScore),
		Max:    toString(maxScore),
		Offset: offset,
		Count:  limit,
	}
	var cmd *redis.ZSliceCmd
	if desc {
		cmd = s.db().ZRevRangeByScoreWithScores(s.ctx, s.key, by)
	} else {
		cmd = s.db().ZRangeByScoreWithScores(s.ctx, s.key, by)
	}
	s.done(cmd)
	return &ZSliceCmd[E]{cmd: cmd}
}

// 排名范围内的成员和其积分。 (从0开始) [左闭右闭] [start, stop] => []{(member1,score1), (member2,score2), ...}
func (s *SortedSet[E]) ZRangeByIndex(start, stop int64, scoreDesc bool) *ZSliceCmd[E] {
	var cmd *redis.ZSliceCmd
	if scoreDesc {
		cmd = s.db().ZRevRangeWithScores(s.ctx, s.key, start, stop)
	} else {
		cmd = s.db().ZRangeWithScores(s.ctx, s.key, start, stop)
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
func (s *SortedSet[E]) ZRemByScore(minScore, maxScore float64) *IntCmd {
	minS := toString(minScore)
	maxS := toString(maxScore)
	cmd := s.db().ZRemRangeByScore(s.ctx, s.key, minS, maxS)
	s.done(cmd)
	return &IntCmd{cmd: cmd}
}

// 移除排名区间内的成员。 [左闭右闭] (按积分降序删除使用负数偏移值)
func (s *SortedSet[E]) ZRemByIndex(start, stop int64) *IntCmd {
	cmd := s.db().ZRemRangeByRank(s.ctx, s.key, start, stop)
	s.done(cmd)
	return &IntCmd{cmd: cmd}
}

func (s *SortedSet[E]) WithCmdable(cmdable Cmdable) *SortedSet[E] {
	b := s.base
	b.cmdable = cmdable
	return &SortedSet[E]{base: b}
}
