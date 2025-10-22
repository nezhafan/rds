package rds

import (
	"cmp"
	"context"

	"github.com/redis/go-redis/v9"
)

var zsetModes = []string{"xx", "nx", "ch"}

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

// 添加，返回新增+修改成功数 https://redis.io/docs/latest/commands/zadd/
func (s *SortedSet[E]) ZAdd(zs map[E]float64) Int64Cmd {
	args := make([]any, 0, len(zs)*2+3)
	args = append(args, "zadd", s.key, "ch")
	for member, score := range zs {
		args = append(args, score, member)
	}
	cmd := s.db().Do(s.ctx, args...)
	s.done(cmd)
	return newInt64Cmd(cmd)
}

// 增加分数 (注意如果操作小数部分，可能会丢失精度)
func (s *SortedSet[E]) ZIncrBy(member E, incr float64) Float64Cmd {
	cmd := s.db().ZIncrBy(s.ctx, s.key, float64(incr), any2String(member))
	s.done(cmd)
	return newFloat64Cmd(cmd)
}

// 所有成员数
func (s *SortedSet[E]) ZCard() Int64Cmd {
	cmd := s.db().ZCard(s.ctx, s.key)
	s.done(cmd)
	return newInt64Cmd(cmd)
}

// 积分区间内的成员数
func (s *SortedSet[E]) ZCountByScore(minScore, maxScore float64) Int64Cmd {
	minS := any2String(minScore)
	maxS := any2String(maxScore)
	cmd := s.db().ZCount(s.ctx, s.key, minS, maxS)
	s.done(cmd)
	return newInt64Cmd(cmd)
}

// 获取分数。 member => score
func (s *SortedSet[E]) ZScore(member E) Float64Cmd {
	cmd := s.db().ZScore(s.ctx, s.key, any2String(member))
	s.done(cmd)
	return newFloat64Cmd(cmd)
}

// 获取位置。 (从0开始) member => index
func (s *SortedSet[E]) ZIndex(member E, scoreDesc bool) Int64Cmd {
	var cmd *redis.IntCmd
	if scoreDesc {
		cmd = s.db().ZRevRank(s.ctx, s.key, any2String(member))
	} else {
		cmd = s.db().ZRank(s.ctx, s.key, any2String(member))
	}
	s.done(cmd)
	return newInt64Cmd(cmd)
}

// 积分范围内的成员。 [左闭右闭] [startScore, stopScore] => []{member1, member2, ...}
// 以 startScore 和 stopScore 的大小关系确定排序方向
// 因为分数区间的数量不确定，所以限制offset、limit。若返回全部都传0即可
func (s *SortedSet[E]) ZMembersByScore(startScore, stopScore float64, offset, limit int64) SliceCmd[E] {
	var desc bool
	if startScore > stopScore {
		desc = true
		startScore, stopScore = stopScore, startScore
	}
	by := &redis.ZRangeBy{
		Min:    any2String(startScore),
		Max:    any2String(stopScore),
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
	return newSliceCmd[E](cmd)
}

// 排序范围内的成员。 (从0开始) [左闭右闭] [start, stop] => []{member1, member2, ...}
func (s *SortedSet[E]) ZMembersByIndex(start, stop int64, scoreDesc bool) SliceCmd[E] {
	var cmd *redis.StringSliceCmd
	if scoreDesc {
		cmd = s.db().ZRevRange(s.ctx, s.key, start, stop)
	} else {
		cmd = s.db().ZRange(s.ctx, s.key, start, stop)
	}
	s.done(cmd)
	return newSliceCmd[E](cmd)
}

// 积分范围内的成员和其积分。  [左闭右闭] [minScore, maxScore] => []{(member1,score1), (member2,score2), ...}
// 若不需要偏移offset和数量限制count，参数传0 即可
func (s *SortedSet[E]) ZRangeByScore(startScore, stopScore float64, offset, limit int64) ZSliceCmd[E] {
	var desc bool
	if startScore > stopScore {
		desc = true
		startScore, stopScore = stopScore, startScore
	}
	by := &redis.ZRangeBy{
		Min:    any2String(startScore),
		Max:    any2String(stopScore),
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
	return newZSliceCmd[E](cmd)
}

// 排名范围内的成员和其积分。 (从0开始) [左闭右闭] [start, stop] => []{(member1,score1), (member2,score2), ...}
func (s *SortedSet[E]) ZRangeByIndex(start, stop int64, scoreDesc bool) ZSliceCmd[E] {
	var cmd *redis.ZSliceCmd
	if scoreDesc {
		cmd = s.db().ZRevRangeWithScores(s.ctx, s.key, start, stop)
	} else {
		cmd = s.db().ZRangeWithScores(s.ctx, s.key, start, stop)
	}
	s.done(cmd)
	return newZSliceCmd[E](cmd)
}

// 移除成员
func (s *SortedSet[E]) ZRem(members ...E) Int64Cmd {
	args := sliceToAnys(members)
	cmd := s.db().ZRem(s.ctx, s.key, args...)
	s.done(cmd)
	return newInt64Cmd(cmd)
}

// 移除积分区间内的成员
func (s *SortedSet[E]) ZRemByScore(minScore, maxScore float64) Int64Cmd {
	minS := any2String(minScore)
	maxS := any2String(maxScore)
	cmd := s.db().ZRemRangeByScore(s.ctx, s.key, minS, maxS)
	s.done(cmd)
	return newInt64Cmd(cmd)
}

// 移除排名区间内的成员。 [左闭右闭] (按积分降序删除使用负数偏移值)
func (s *SortedSet[E]) ZRemByIndex(start, stop int64) Int64Cmd {
	cmd := s.db().ZRemRangeByRank(s.ctx, s.key, start, stop)
	s.done(cmd)
	return newInt64Cmd(cmd)
}
