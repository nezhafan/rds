package rds

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/constraints"
)

type Z = redis.Z

const (
	MinInf = "-inf"
	MaxInf = "+inf"
)

type Z2[M constraints.Ordered, S int | float32] struct {
	Member M
	Score  S
}

type ZSet[M constraints.Ordered, S int | float32] struct {
	base
}

func NewZSet[M constraints.Ordered, S int | float32](ctx context.Context, key string) ZSet[M, S] {
	return ZSet[M, S]{base: newBase(ctx, key)}
}

// 添加成员. 返回新增成员数（修改不算新增）
func (z *ZSet[M, S]) ZAdd(args ...Z2[M, S]) (int64, error) {
	slice := make([]any, len(args)*2+2)
	slice = append(slice, "ZADD", z.key)
	for i := range args {
		slice = append(slice, args[i].Score)
		slice = append(slice, args[i].Member)
	}

	cmd := DB().Do(z.ctx, slice...)
	z.done(cmd)
	return cmd.Int64()
}

// 移除成员
func (z *ZSet[M, S]) ZRem(members ...M) (int64, error) {
	vs := toAnys(members)
	return DB().ZRem(z.ctx, z.key, vs...).Result()
}

// 获取成员数
func (z *ZSet[M, S]) ZCard() int64 {
	return DB().ZCard(z.ctx, z.key).Val()
}

// 获取member的score
func (z *ZSet[M, S]) ZScore(member string) (S, bool) {
	v, err := DB().ZScore(z.ctx, z.key, member).Result()
	return S(v), err == nil
}

// 从小到大获取member
func (z *ZSet[M, S]) ZRange(start, stop int64) []M {
	vs := DB().ZRange(z.ctx, z.key, start, stop).Val()
	return stringsToSlice[M](vs)
}

// 从小到大获取member和score
func (z *ZSet[M, S]) ZRangeWithScores(start, stop int64) []redis.Z {
	args := []any{"zrange", z.key, "withscores", start, stop}
	slices, err := DB().Do(z.ctx, args...).StringSlice()
	fmt.Println("slices", slices, err)
	return DB().ZRangeWithScores(z.ctx, z.key, start, stop).Val()
}

// 从大到小获取member
func (z *ZSet[M, S]) ZRevRange(start, stop int64) []string {
	return DB().ZRevRange(z.ctx, z.key, start, stop).Val()
}

// 从大到小获取member和score
func (z *ZSet[M, S]) ZRevRangeWithScores(start, stop int64) []Z {
	return DB().ZRevRangeWithScores(z.ctx, z.key, start, stop).Val()
}

// 获取积分区间(闭区间)内的member和score.
// min和max可以传 MinInf和MaxInf , limit 为0时返回全部
func (z *ZSet[M, S]) ZRangeByScore(min, max string, limit int64) []Z {
	by := &redis.ZRangeBy{Min: min, Max: max, Count: limit}
	return DB().ZRangeByScoreWithScores(z.ctx, z.key, by).Val()
}

// 移除指定score区间内的所有成员
func (z *ZSet[M, S]) ZRemRangeByScore(min, max string) (int64, error) {
	return DB().ZRemRangeByScore(z.ctx, z.key, min, max).Result()
}

// member区间内成员数，可以传入 -、+  或者 [A、[B
// https://www.runoob.com/redis/sorted-sets-zlexcount.html
func (z *ZSet[M, S]) ZLexCount(minMember, maxMember string) int64 {
	return DB().ZLexCount(z.ctx, z.key, minMember, maxMember).Val()
}

// memberA和memberB中间的 member
func (z *ZSet[M, S]) ZRangeByLex(memberA, memberB string, limit int64) []string {
	by := &redis.ZRangeBy{Min: memberA, Max: memberB, Count: limit}
	return DB().ZRangeByLex(z.ctx, z.key, by).Val()
}

// 获取成员从小到大的排名
func (z *ZSet[M, S]) ZRank(member string) int64 {
	return DB().ZRank(z.ctx, z.key, member).Val()
}

// 获取成员从大到小的排名
func (z *ZSet[M, S]) ZRevRank(member string) int64 {
	return DB().ZRevRank(z.ctx, z.key, member).Val()
}

// 删除从小到大的排名区间
func (z *ZSet[M, S]) ZRemRangeByRank(start, stop int64) (int64, error) {
	return DB().ZRemRangeByRank(z.ctx, z.key, start, stop).Result()
}
