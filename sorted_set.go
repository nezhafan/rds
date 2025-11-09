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

type sortedType int

const (
	ScoreASC  sortedType = 1
	ScoreDESC sortedType = 2
)

/*
有序集合
  - member不可重复，score可重复
  - score相同时，会按member字典排序，而不是按先后添加顺序
*/
func NewSortedSet[E cmp.Ordered](ctx context.Context, key string) *SortedSet[E] {
	return &SortedSet[E]{base: newBase(ctx, key)}
}

/*
添加/修改 https://redis.io/docs/latest/commands/zadd/
zadd key [nx|xx] [gt|lt] [ch] score member
  - 不传参新增或更新；nx仅新增不更新； xx仅更新不新增
  - gt仅新分值大于原分值时更新；lt仅新分值小于原分值时更新；不影响新增；6.2.0版本有效
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

// 增加分值 (注意可能会丢失精度)
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

// 分值区间内的成员数
func (s *SortedSet[E]) ZCountByScore(min, max float64) Int64Cmd {
	minS := any2String(min)
	maxS := any2String(max)
	cmd := s.db().ZCount(s.ctx, s.key, minS, maxS)
	s.done(cmd)
	return newInt64Cmd(cmd)
}

// 获取分值。 member => score
func (s *SortedSet[E]) ZScore(member E) Float64Cmd {
	cmd := s.db().ZScore(s.ctx, s.key, any2String(member))
	s.done(cmd)
	return newFloat64Cmd(cmd)
}

/*
获取排名 member => rank
  - 参数 scoreSort 分值排序方式 rds.ScoreDESC从高到低 rds.ScoreASC从低到高
  - 参数 member 成员
  - 返回 按照分值排序方向返回 排名 （从0开始）
*/
func (s *SortedSet[E]) ZRank(scoreSort sortedType, member E) Int64Cmd {
	var cmd *redis.IntCmd
	if scoreSort == ScoreASC {
		cmd = s.db().ZRank(s.ctx, s.key, any2String(member))
	} else {
		cmd = s.db().ZRevRank(s.ctx, s.key, any2String(member))
	}
	s.done(cmd)
	return newInt64Cmd(cmd)
}

/*
分值范围内的成员
  - 参数 scoreSort 分值排序方式 rds.ScoreDESC从高到低 rds.ScoreASC从低到高
  - 参数 min max [左闭右闭]
  - 参数 因为分值区间的元素数量不确定，所以限制offset、limit，若返回全部都传0即可
  - 返回 按照分值排序方向返回 []{member1, member2, ...}
*/
func (s *SortedSet[E]) ZMembersByScore(scoreSort sortedType, min, max float64, offset, limit int64) SliceCmd[E] {
	by := &redis.ZRangeBy{
		Min:    any2String(min),
		Max:    any2String(max),
		Offset: offset,
		Count:  limit,
	}
	var cmd *redis.StringSliceCmd
	if scoreSort == ScoreASC {
		cmd = s.db().ZRangeByScore(s.ctx, s.key, by)
	} else {
		cmd = s.db().ZRevRangeByScore(s.ctx, s.key, by)
	}
	s.done(cmd)
	return newSliceCmd[E](cmd)
}

/*
排序范围内的成员
  - 参数 scoreSort 分值排序方式 rds.ScoreDESC从高到低 rds.ScoreASC从低到高
  - 参数 start stop [左闭右闭]
  - 返回 按照分值排序方向返回 []{member1, member2, ...}
  - 举例 考试成绩 A-94分  B-96分 C-98分 D-100分 （成绩排名是从分值降序排列看）
  - 考试成绩最高的前2名 ZMembersByRank(rds.ScoreDESC, 0, 1) 返回 D、C
  - 考试成绩最低的后2名 ZMembersByRank(rds.ScoreDESC, -2, -1) 返回 B、A
  - 考试成绩最低的后2名 ZMembersByRank(rds.ScoreASC, 0, 1) 返回 A、B，这是从分值升序的视角返回，相同的元素，排序方式不同
*/
func (s *SortedSet[E]) ZMembersByRank(scoreSort sortedType, start, stop int64) SliceCmd[E] {
	var cmd *redis.StringSliceCmd
	if scoreSort == ScoreASC {
		cmd = s.db().ZRange(s.ctx, s.key, start, stop)
	} else {
		cmd = s.db().ZRevRange(s.ctx, s.key, start, stop)
	}
	s.done(cmd)
	return newSliceCmd[E](cmd)
}

/*
分值范围内的成员和其分值。
  - 参数 scoreSort 分值排序方式 rds.ScoreDESC从高到低 rds.ScoreASC从低到高
  - 参数 min max [左闭右闭] [min, max]
  - 参数 是否按照从大到小排序返回
  - 参数 因为分值区间的元素数量不确定，所以限制offset、limit，若返回全部都传0即可
  - 返回 按照分值排序方向返回 []{(member1,score1), (member2,score2), ...}
*/
func (s *SortedSet[E]) ZRangeByScore(scoreSort sortedType, min, max float64, offset, limit int64) ZSliceCmd[E] {
	by := &redis.ZRangeBy{
		Min:    any2String(min),
		Max:    any2String(max),
		Offset: offset,
		Count:  limit,
	}
	var cmd *redis.ZSliceCmd
	if scoreSort == ScoreASC {
		cmd = s.db().ZRangeByScoreWithScores(s.ctx, s.key, by)
	} else {
		cmd = s.db().ZRevRangeByScoreWithScores(s.ctx, s.key, by)
	}
	s.done(cmd)
	return newZSliceCmd[E](cmd)
}

/*
排序范围内的成员和其分值
  - 参数 scoreSort 分值排序方式 rds.ScoreDESC从高到低 rds.ScoreASC从低到高
  - 参数 start stop [左闭右闭]
  - 返回 按照分值排序方向返回 []{member1, member2, ...}
  - 举例 考试成绩 A-94分  B-96分 C-98分 D-100分 （成绩排名是从分值降序排列看）
  - 考试成绩最高的前2名 ZMembersByRank(rds.ScoreDESC, 0, 1) 返回 D-100、C-98
  - 考试成绩最低的后2名 ZMembersByRank(rds.ScoreDESC, -2, -1) 返回 B-96、A-94
  - 考试成绩最低的后2名 ZMembersByRank(rds.ScoreASC, 0, 1) 返回 A-94、B-96，这是从分值升序的视角返回，相同的元素，排序方式不同
*/
func (s *SortedSet[E]) ZRangeByRank(scoreSort sortedType, start, stop int64) ZSliceCmd[E] {
	var cmd *redis.ZSliceCmd
	if scoreSort == ScoreASC {
		cmd = s.db().ZRangeWithScores(s.ctx, s.key, start, stop)
	} else {
		cmd = s.db().ZRevRangeWithScores(s.ctx, s.key, start, stop)
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

/*
移除分值区间内的成员
  - 参数 min max [左闭右闭]
*/
func (s *SortedSet[E]) ZRemByScore(min, max float64) Int64Cmd {
	minS := any2String(min)
	maxS := any2String(max)
	cmd := s.db().ZRemRangeByScore(s.ctx, s.key, minS, maxS)
	s.done(cmd)
	return newInt64Cmd(cmd)
}

/*
移除排名区间内的成员。
  - 参数 scoreSort 分值排序方式 rds.ScoreDESC从高到低 rds.ScoreASC从低到高
  - 参数 start stop [左闭右闭]
  - 返回 移除数量
  - 举例 考试成绩 A-94分  B-96分 C-98分 D-100分 （成绩排名是从分值降序排列看）
  - 只保留前2名 即删除第三名之后所有 ZRemByRank(rds.ScoreDESC, 2, -1)
  - 移除最后2名 ZRemByRank(rds.ScoreDESC, -2, -1) 或 ZRemByRank(rds.ScoreASC, 0, 1)
*/
func (s *SortedSet[E]) ZRemByRank(scoreSort sortedType, start, stop int64) Int64Cmd {
	var cmd *redis.IntCmd
	if scoreSort == ScoreASC {
		cmd = s.db().ZRemRangeByRank(s.ctx, s.key, start, stop)
	} else {
		start, stop = 0-stop-1, 0-start-1
		cmd = s.db().ZRemRangeByRank(s.ctx, s.key, start, stop)
	}
	s.done(cmd)
	return newInt64Cmd(cmd)
}
