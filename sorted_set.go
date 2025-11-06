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

/*
有序集合
  - member不可重复，score可重复
  - score相同时，会按member字典排序，而不是先后添加顺序
*/
func NewSortedSet[E cmp.Ordered](ctx context.Context, key string) *SortedSet[E] {
	return &SortedSet[E]{base: NewBase(ctx, key)}
}

/*
添加/修改 https://redis.io/docs/latest/commands/zadd/
zadd key [nx|xx] [gt|lt] [ch] score member
  - 不传参新增或更新；nx仅新增不更新； xx仅更新不新增
  - gt仅新分数大于原分数时更新；lt仅新分数小于原分数时更新；不影响新增；6.2.0版本有效
  - ch返回结果是否要加上更新成功数，默认仅返回新增数。
*/
func (s *SortedSet[E]) ZAdd(zs map[E]float64, params ...string) Int64Cmd {
	args := make([]any, 0, len(zs)*2+5)
	args = append(args, "zadd", s.key)
	if len(params) > 0 {
		args = append(args, params[0])
	}
	if len(params) > 1 {
		args = append(args, params[1])
	}
	if len(params) > 2 {
		args = append(args, params[2])
	}
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
func (s *SortedSet[E]) ZCountByScore(min, max float64) Int64Cmd {
	if min > max {
		min, max = max, min
	}
	minS := any2String(min)
	maxS := any2String(max)
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

// 获取排名。  member => rank (从0开始)
func (s *SortedSet[E]) ZRank(member E, scoreDesc bool) Int64Cmd {
	var cmd *redis.IntCmd
	if scoreDesc {
		cmd = s.db().ZRevRank(s.ctx, s.key, any2String(member))
	} else {
		cmd = s.db().ZRank(s.ctx, s.key, any2String(member))
	}
	s.done(cmd)
	return newInt64Cmd(cmd)
}

/*
积分范围内的成员
  - 参数 [左闭右闭] [min, max]
  - 参数 是否按照从大到小排序
  - 参数 因为分数区间的元素数量不确定，所以限制offset、limit。若返回全部都传0即可
  - 返回 []{member1, member2, ...}
*/
func (s *SortedSet[E]) ZMembersByScore(min, max float64, scoreDesc bool, offset, limit int64) SliceCmd[E] {
	if min > max {
		min, max = max, min
	}
	by := &redis.ZRangeBy{
		Min:    any2String(min),
		Max:    any2String(max),
		Offset: offset,
		Count:  limit,
	}
	var cmd *redis.StringSliceCmd
	if scoreDesc {
		cmd = s.db().ZRevRangeByScore(s.ctx, s.key, by)
	} else {
		cmd = s.db().ZRangeByScore(s.ctx, s.key, by)
	}
	s.done(cmd)
	return newSliceCmd[E](cmd)
}

/*
排序范围内的成员。
  - 参数 [左闭右闭] [start, stop] (从0开始)
  - 返回 []{member1, member2, ...}
*/
func (s *SortedSet[E]) ZMembersByRank(start, stop int64, scoreDesc bool) SliceCmd[E] {
	var cmd *redis.StringSliceCmd
	if scoreDesc {
		cmd = s.db().ZRevRange(s.ctx, s.key, start, stop)
	} else {
		cmd = s.db().ZRange(s.ctx, s.key, start, stop)
	}
	s.done(cmd)
	return newSliceCmd[E](cmd)
}

/*
积分范围内的成员和其积分。
  - 参数 [左闭右闭] [min, max]
  - 参数 是否按照从大到小排序
  - 参数 因为分数区间的元素数量不确定，所以限制offset、limit。若返回全部都传0即可
  - 返回 []{(member1,score1), (member2,score2), ...}
*/
func (s *SortedSet[E]) ZRangeByScore(min, max float64, scoreDesc bool, offset, limit int64) ZSliceCmd[E] {
	if min > max {
		min, max = max, min
	}
	by := &redis.ZRangeBy{
		Min:    any2String(min),
		Max:    any2String(max),
		Offset: offset,
		Count:  limit,
	}
	var cmd *redis.ZSliceCmd
	if scoreDesc {
		cmd = s.db().ZRevRangeByScoreWithScores(s.ctx, s.key, by)
	} else {
		cmd = s.db().ZRangeByScoreWithScores(s.ctx, s.key, by)
	}
	s.done(cmd)
	return newZSliceCmd[E](cmd)
}

/*
排名范围内的成员和其积分。
  - 参数 [左闭右闭] [start, stop] (从0开始)
  - 参数 是否按照从大到小排序
  - 返回 []{(member1,score1), (member2,score2), ...}
*/
func (s *SortedSet[E]) ZRangeByRank(start, stop int64, scoreDesc bool) ZSliceCmd[E] {
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
	args := slice2Anys(members)
	cmd := s.db().ZRem(s.ctx, s.key, args...)
	s.done(cmd)
	return newInt64Cmd(cmd)
}

// 移除积分区间内的成员 [左闭右闭]
func (s *SortedSet[E]) ZRemByScore(min, max float64) Int64Cmd {
	if min > max {
		min, max = max, min
	}
	minS := any2String(min)
	maxS := any2String(max)
	cmd := s.db().ZRemRangeByScore(s.ctx, s.key, minS, maxS)
	s.done(cmd)
	return newInt64Cmd(cmd)
}

/*
移除排名区间内的成员。
  - 参数 [左闭右闭] 都是正数则从低到高排序移除（从0开始 ），都是负数则从高到低排序移除（-1开始）
  - 举例 移除分数最低的两个成员 ZRemByRank(0, 1)
  - 举例 移除分数最高的两个成员 ZRemByRank(-1, -2)
*/
func (s *SortedSet[E]) ZRemByRank(start, stop int64) Int64Cmd {
	cmd := s.db().ZRemRangeByRank(s.ctx, s.key, start, stop)
	s.done(cmd)
	return newInt64Cmd(cmd)
}
