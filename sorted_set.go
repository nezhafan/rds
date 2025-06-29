package rds

import (
	"github.com/redis/go-redis/v9"
)

type Order string

const (
	Inf          = "+inf"
	NegInf       = "-inf"
	ASC    Order = "ASC"
	DESC   Order = "DESC"
)

type SortedSet[M Ordered] struct {
	base
}

type Z[M Ordered] struct {
	Score  float64
	Member M
}

func NewSortedSet[M Ordered](key string, ops ...Option) *SortedSet[M] {
	return &SortedSet[M]{base: newBase(key, ops...)}
}

// 添加
func (s *SortedSet[M]) ZAdd(zs map[M]float64) *IntCmd {
	args := make([]any, 0, len(zs)*2+3)
	args = append(args, "zadd", s.key, "ch")
	for member, score := range zs {
		args = append(args, score, member)
	}
	cmd := s.db().Do(ctx, args...)
	s.done(cmd)
	return &IntCmd{cmd: cmd}
}

// 增长积分
// func (s *SortedSet[M, S]) ZIncrBy(member M, increment S) *FloatCmd[S] {
// 	cmd := s.db().ZIncrBy(ctx, s.key, float64(increment), member)
// 	return &FloatCmd[S]{cmd: cmd}
// }

// 所有成员数
func (s *SortedSet[M]) ZCard() *IntCmd {
	cmd := s.db().ZCard(ctx, s.key)
	s.done(cmd)
	return &IntCmd{cmd: cmd}
}

// 积分区间内的成员数
func (s *SortedSet[M]) ZCount(min, max float64) *IntCmd {
	minS := toString(min)
	maxS := toString(max)
	cmd := s.db().ZCount(ctx, s.key, minS, maxS)
	s.done(cmd)
	return &IntCmd{cmd: cmd}
}

// 获取分数
func (s *SortedSet[M]) ZScore(member M) *FloatCmd {
	cmd := s.db().ZScore(ctx, s.key, toString(member))
	s.done(cmd)
	return &FloatCmd{cmd: cmd}
}

// 增加分数
func (s *SortedSet[M]) ZIncrBy(member M, incr float64) *FloatCmd {
	cmd := s.db().ZIncrBy(ctx, s.key, float64(incr), toString(member))
	s.done(cmd)
	return &FloatCmd{cmd: cmd}
}

// 按照积分获取：成员
func (s *SortedSet[M]) ZMembersByScore(min, max float64, offset, limit int64, order Order) *SliceCmd[M] {
	by := &redis.ZRangeBy{
		Min:    toString(min),
		Max:    toString(max),
		Offset: offset,
		Count:  limit,
	}
	var cmd *redis.StringSliceCmd
	if order == ASC {
		cmd = s.db().ZRangeByScore(ctx, s.key, by)
	} else if order == DESC {
		cmd = s.db().ZRevRangeByScore(ctx, s.key, by)
	}
	s.done(cmd)
	return &SliceCmd[M]{cmd: cmd}
}

// 按照积分获取：成员
func (s *SortedSet[M]) ZMembersByRank(start, stop int64, order Order) *SliceCmd[M] {
	var cmd *redis.StringSliceCmd
	if order == ASC {
		cmd = s.db().ZRange(ctx, s.key, start, stop)
	} else if order == DESC {
		cmd = s.db().ZRevRange(ctx, s.key, start, stop)
	}
	s.done(cmd)
	return &SliceCmd[M]{cmd: cmd}
}

// 按照积分获取：成员+积分
func (s *SortedSet[M]) ZItemsByScore(min, max float64, offset, limit int64, order Order) *ZSliceCmd[M] {
	by := &redis.ZRangeBy{
		Min:    toString(min),
		Max:    toString(max),
		Offset: offset,
		Count:  limit,
	}
	var cmd *redis.ZSliceCmd
	if order == ASC {
		cmd = s.db().ZRangeByScoreWithScores(ctx, s.key, by)
	} else if order == DESC {
		cmd = s.db().ZRevRangeByScoreWithScores(ctx, s.key, by)
	}
	s.done(cmd)
	return &ZSliceCmd[M]{cmd: cmd}
}

// 按照积分获取：成员+积分
func (s *SortedSet[M]) ZItemsByRank(start, stop int64, order Order) *ZSliceCmd[M] {
	var cmd *redis.ZSliceCmd
	if order == ASC {
		cmd = s.db().ZRangeWithScores(ctx, s.key, start, stop)
	} else if order == DESC {
		cmd = s.db().ZRevRangeWithScores(ctx, s.key, start, stop)
	}
	s.done(cmd)
	return &ZSliceCmd[M]{cmd: cmd}
}

// 迭代
// func (s *SortedSet[M]) ZScan(cursor uint64, match string, count int64)

// 移除成员
func (s *SortedSet[M]) ZRem(members ...M) *IntCmd {
	args := toAnys(members)
	cmd := s.db().ZRem(ctx, s.key, args...)
	s.done(cmd)
	return &IntCmd{cmd: cmd}
}

// 移除积分区间内的成员
func (s *SortedSet[M]) ZRemByScore(min, max float64) *IntCmd {
	minS := toString(min)
	maxS := toString(max)
	cmd := s.db().ZRemRangeByScore(ctx, s.key, minS, maxS)
	s.done(cmd)
	return &IntCmd{cmd: cmd}
}

// 移除排名区间内的成员
func (s *SortedSet[M]) ZRemByRank(start, stop int64) *IntCmd {
	cmd := s.db().ZRemRangeByRank(ctx, s.key, start, stop)
	s.done(cmd)
	return &IntCmd{cmd: cmd}
}
