package rds

import (
	"strconv"

	"github.com/redis/go-redis/v9"
)

const (
	Inf    = "+inf"
	NegInf = "-inf"
)

type SortedSet[M Ordered, S Number] struct {
	base
}

type Z[M Ordered, S Number] struct {
	Score  S
	Member M
}

func NewSortedSet[M Ordered, S Number](key string, ops ...Option) *SortedSet[M, S] {
	return &SortedSet[M, S]{base: newBase(key, ops...)}
}

// 添加
func (s *SortedSet[M, S]) ZAdd(zs ...Z[M, S]) *IntCmd {
	args := make([]any, 0, len(zs)*2+3)
	args = append(args, "ZADD", s.key, "CH")
	for i := range zs {
		args = append(args, zs[i].Score, zs[i].Member)
	}
	cmd := s.db().Do(ctx, args...)
	return &IntCmd{cmd: cmd}
}

// 移除
func (s *SortedSet[M, S]) ZRem(members ...M) *IntCmd {
	args := toAnys(members)
	cmd := s.db().ZRem(ctx, s.key, args...)
	return &IntCmd{cmd: cmd}
}

// 成员数
func (s *SortedSet[M, S]) ZCard() *IntCmd {
	cmd := s.db().ZCard(ctx, s.key)
	return &IntCmd{cmd: cmd}
}

// 成员数
func (s *SortedSet[M, S]) ZCount(min, max S) *IntCmd {
	minS := strconv.Itoa(int(min))
	maxS := strconv.Itoa(int(max))
	cmd := s.db().ZCount(ctx, s.key, minS, maxS)
	return &IntCmd{cmd: cmd}
}

// 成员。 积分从低到高
func (s *SortedSet[M, S]) ZRangeByScore(min, max, offset, limit int64) *SliceCmd[M] {
	by := &redis.ZRangeBy{
		Min:    strconv.Itoa(int(min)),
		Max:    strconv.Itoa(int(max)),
		Offset: 0,
		Count:  limit,
	}
	cmd := s.db().ZRangeByScore(ctx, s.key, by)
	return &SliceCmd[M]{cmd: cmd}
}

// 成员。 积分从高到低
func (s *SortedSet[M, S]) ZRevRangeByScore(min, max, offset, limit int64) *SliceCmd[M] {
	by := &redis.ZRangeBy{
		Min:    strconv.Itoa(int(min)),
		Max:    strconv.Itoa(int(max)),
		Offset: 0,
		Count:  limit,
	}
	cmd := s.db().ZRevRangeByScore(ctx, s.key, by)

	return &SliceCmd[M]{cmd: cmd}
}

func (s *SortedSet[M, S]) ZRangeByScoreWithScores(start, stop, offset, limit int64) []Z[M, S] {
	by := &redis.ZRangeBy{
		Min:    strconv.Itoa(int(start)),
		Max:    strconv.Itoa(int(stop)),
		Offset: 0,
	}
	cmd := s.db().ZRangeByScoreWithScores(ctx, s.key, by)
	list := make([]Z[M, S], 0, len(cmd.Val()))
	for _, v := range cmd.Val() {
		var member M
		if v, ok := v.Member.(string); ok {
			member = stringTo[M](v)
		}
		list = append(list, Z[M, S]{
			Score:  S(v.Score),
			Member: member,
		})
	}

	return list
}
