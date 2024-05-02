package redis

import "github.com/redis/go-redis/v9"

type Z = redis.Z

const (
	MinInf = "-inf"
	MaxInf = "+inf"
)

type zSet struct {
	base
}

func NewZSet(key string) zSet {
	return zSet{base{key}}
}

// 添加成员
func (z zSet) ZAdd(args ...Z) (int64, error) {
	return rdb.ZAdd(ctx, z.key, args...).Result()
}

// 移除成员
func (z zSet) ZRem(members ...any) (int64, error) {
	return rdb.ZRem(ctx, z.key, members...).Result()
}

// 获取成员数
func (z zSet) ZCard() int64 {
	return rdb.ZCard(ctx, z.key).Val()
}

// 获取member的score
func (z zSet) ZScore(member string) float64 {
	return rdb.ZScore(ctx, z.key, member).Val()
}

// 从小到大获取member
func (z zSet) ZRange(start, stop int64) []string {
	return rdb.ZRange(ctx, z.key, start, stop).Val()
}

// 从小到大获取member和score
func (z zSet) ZRangeWithScores(start, stop int64) []Z {
	return rdb.ZRangeWithScores(ctx, z.key, start, stop).Val()
}

// 从大到小获取member
func (z zSet) ZRevRange(start, stop int64) []string {
	return rdb.ZRevRange(ctx, z.key, start, stop).Val()
}

// 从大到小获取member和score
func (z zSet) ZRevRangeWithScores(start, stop int64) []Z {
	return rdb.ZRevRangeWithScores(ctx, z.key, start, stop).Val()
}

// 获取积分区间(闭区间)内的member和score.
// min和max可以传 MinInf和MaxInf , limit 为0时返回全部
func (z zSet) ZRangeByScore(min, max string, limit int64) []Z {
	by := &redis.ZRangeBy{Min: min, Max: max, Count: limit}
	return rdb.ZRangeByScoreWithScores(ctx, z.key, by).Val()
}

// 移除指定score区间内的所有成员
func (z zSet) ZRemRangeByScore(min, max string) (int64, error) {
	return rdb.ZRemRangeByScore(ctx, z.key, min, max).Result()
}

// member区间内成员数，可以传入 -、+  或者 [A、[B
// https://www.runoob.com/redis/sorted-sets-zlexcount.html
func (z zSet) ZLexCount(minMember, maxMember string) int64 {
	return rdb.ZLexCount(ctx, z.key, minMember, maxMember).Val()
}

// memberA和memberB中间的 member
func (z zSet) ZRangeByLex(memberA, memberB string, limit int64) []string {
	by := &redis.ZRangeBy{Min: memberA, Max: memberB, Count: limit}
	return rdb.ZRangeByLex(ctx, z.key, by).Val()
}

// 获取成员从小到大的排名
func (z zSet) ZRank(member string) int64 {
	return rdb.ZRank(ctx, z.key, member).Val()
}

// 获取成员从大到小的排名
func (z zSet) ZRevRank(member string) int64 {
	return rdb.ZRevRank(ctx, z.key, member).Val()
}

// 删除从小到大的排名区间
func (z zSet) ZRemRangeByRank(start, stop int64) (int64, error) {
	return rdb.ZRemRangeByRank(ctx, z.key, start, stop).Result()
}
