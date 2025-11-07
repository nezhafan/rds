package test

import (
	"cmp"
	"fmt"
	"testing"

	"github.com/nezhafan/rds"
	"github.com/stretchr/testify/assert"
)

var zsetData = map[string]float64{"a": 0.3, "b": 2.2, "c": 3.3, "d": 4.4}

func newSortedSet[E cmp.Ordered]() *rds.SortedSet[E] {
	return rds.NewSortedSet[E](ctx, "sorted_set_test")
}

func TestSortedSet_ZAdd(t *testing.T) {
	cache := newSortedSet[int]()

	// 新增且更新。 新增3个，返回新增数
	v, err := cache.ZAdd(map[int]float64{100: 100, 200: 200, 300: 300}).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 3, v)

	// 新增且更新。 新增1个、更新1个、无更新1个，返回新增数
	v, err = cache.ZAdd(map[int]float64{100: 100.1, 200: 200, 400: 400}).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, v)

	// 新增且更新。 新增1个、更新1个、无更新1个，返回新增数+更新成功数
	v, err = cache.ZAdd(map[int]float64{100: 100.2, 200: 200, 500: 500}, "ch").Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 2, v)

	// nx 新增不更新。 新增1个、更新1个、无更新1个，返回新增数
	v, err = cache.ZAdd(map[int]float64{100: 100.3, 200: 200, 600: 600}, "nx").Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, v)
	assert.EqualValues(t, 100.2, cache.ZScore(100).Val()) // 更新不生效
	assert.EqualValues(t, 600, cache.ZScore(600).Val())

	// xx 更新不新增。 新增1个、更新1个、无更新1个，返回新增数（恒定0）
	v, err = cache.ZAdd(map[int]float64{100: 100.4, 200: 200, 700: 700}, "xx").Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 0, v)
	assert.EqualValues(t, 100.4, cache.ZScore(100).Val()) // 更新成功
	assert.EqualValues(t, 0, cache.ZScore(700).Val())     // 新增不生效

	// xx ch更新不新增。 新增1个、更新1个、无更新1个
	v, err = cache.ZAdd(map[int]float64{100: 100.5, 200: 200, 800: 800}, "xx", "ch").Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, v)                           // 可以返回更新成功数
	assert.EqualValues(t, 100.5, cache.ZScore(100).Val()) // 更新成功
	assert.EqualValues(t, 0, cache.ZScore(800).Val())     // 新增不生效

	// xx gt更新不新增。 新增1个、更新2个(一个变大，一个变小)
	v, err = cache.ZAdd(map[int]float64{100: 100.6, 200: 199, 900: 900}, "xx", "gt").Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 0, v)                           // 可以返回更新成功数
	assert.EqualValues(t, 100.6, cache.ZScore(100).Val()) // 更新成功
	assert.EqualValues(t, 200, cache.ZScore(200).Val())   // 更新不成功
	assert.EqualValues(t, 0, cache.ZScore(900).Val())     // 新增不生效

	cache.Del()
}

func TestSortedSet_ZIncrBy(t *testing.T) {
	cache := newSortedSet[string]()
	cache.ZAdd(zsetData)

	v, err := cache.ZIncrBy("a", 3).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 3.3, v)

	v, err = cache.ZIncrBy("c", -3).Result()
	assert.NoError(t, err)
	assert.NotEqualValues(t, 0.3, v)
	assert.InDelta(t, 0.3, v, 0.00000001) // 精度丢失

	v, err = cache.ZIncrBy("d", 0.1).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 4.5, v)

	cache.Del()
}

func TestSortedSet_ZCard(t *testing.T) {
	cache := newSortedSet[string]()
	cache.ZAdd(zsetData)

	v, err := cache.ZCard().Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(4), v)

	cache.Del()
}

func TestSortedSet_ZCountByScore(t *testing.T) {
	cache := newSortedSet[string]()
	cache.ZAdd(zsetData)

	v, err := cache.ZCountByScore(1.5, 3.5).Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), v)

	cache.Del()
}

func TestSortedSet_ZScore(t *testing.T) {
	cache := newSortedSet[string]()
	cache.ZAdd(zsetData)

	v, err := cache.ZScore("a").Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 0.3, v)

	v, err = cache.ZScore("d").Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 4.4, v)

	cache.Del()
}

func TestSortedSet_ZRank(t *testing.T) {
	cache := newSortedSet[string]()
	cache.ZAdd(zsetData)

	// 从小到大
	v, err := cache.ZRank(rds.ScoreASC, "a").Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 0, v)

	v, err = cache.ZRank(rds.ScoreASC, "c").Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 2, v)

	// 从大到小
	v, err = cache.ZRank(rds.ScoreDESC, "a").Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 3, v)

	cache.Del()
}

func TestSortedSet_ZMembersByScore(t *testing.T) {
	cache := newSortedSet[string]()
	cache.ZAdd(zsetData)

	// 从小到大
	v, err := cache.ZMembersByScore(rds.ScoreASC, 2.2, 4.4, 0, 2).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, []string{"b", "c"}, v)

	// 从大到小
	v, err = cache.ZMembersByScore(rds.ScoreDESC, 2.2, 4.4, 0, 2).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, []string{"d", "c"}, v)

	cache.Del()
}

func TestSortedSet_ZMembersByRank(t *testing.T) {
	cache := newSortedSet[string]()
	cache.ZAdd(zsetData)

	// 分数高-低 前2名
	v, err := cache.ZMembersByRank(rds.ScoreDESC, 0, 1).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, []string{"d", "c"}, v)

	// 分数高-低 后2名
	v, err = cache.ZMembersByRank(rds.ScoreDESC, -2, -1).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, []string{"b", "a"}, v)

	// 分数低-高 前两名 （即分数高-低的后2名）
	v, err = cache.ZMembersByRank(rds.ScoreASC, 0, 1).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, []string{"a", "b"}, v)

	cache.Del()
}

func TestSortedSet_ZRangeByScore(t *testing.T) {
	cache := newSortedSet[string]()
	cache.ZAdd(zsetData)

	v, err := cache.ZRangeByScore(rds.ScoreASC, 2.2, 3.3, 0, 2).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, []rds.Z[string]{{Member: "b", Score: 2.2}, {Member: "c", Score: 3.3}}, v)

	v, err = cache.ZRangeByScore(rds.ScoreDESC, 2.2, 3.3, 0, 2).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, []rds.Z[string]{{Member: "c", Score: 3.3}, {Member: "b", Score: 2.2}}, v)

	cache.Del()
}

func TestSortedSet_ZRangeByRank(t *testing.T) {
	cache := newSortedSet[string]()
	cache.ZAdd(zsetData)

	v, err := cache.ZRangeByRank(rds.ScoreDESC, 0, 1).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, []rds.Z[string]{{Member: "d", Score: 4.4}, {Member: "c", Score: 3.3}}, v)

	v, err = cache.ZRangeByRank(rds.ScoreDESC, -2, -1).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, []rds.Z[string]{{Member: "b", Score: 2.2}, {Member: "a", Score: 0.3}}, v)

	v, err = cache.ZRangeByRank(rds.ScoreASC, 0, 1).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, []rds.Z[string]{{Member: "a", Score: 0.3}, {Member: "b", Score: 2.2}}, v)

	cache.Del()
}

func TestSortedSet_ZRem(t *testing.T) {
	cache := newSortedSet[string]()
	cache.ZAdd(zsetData)

	v, err := cache.ZRem("a", "b", "x").Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), v)

	arr, err := cache.ZMembersByScore(rds.ScoreASC, 0, 100, 0, 0).Result()
	assert.NoError(t, err)
	assert.Equal(t, []string{"c", "d"}, arr)

	cache.Del()
}

func TestSortedSet_ZRemByScore(t *testing.T) {
	cache := newSortedSet[string]()
	cache.ZAdd(zsetData)

	v, err := cache.ZRemByScore(2.2, 3.3).Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), v)

	arr, err := cache.ZMembersByScore(rds.ScoreASC, 0, 100, 0, 0).Result()
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "d"}, arr)

	cache.Del()
}

func TestSortedSet_ZRemByRank(t *testing.T) {

	cache := newSortedSet[string]()

	// 从小到大移除前两个 a b
	cache.ZAdd(zsetData)
	v, err := cache.ZRemByRank(rds.ScoreASC, 0, 1).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 2, v)
	arr, err := cache.ZMembersByScore(rds.ScoreASC, 0, 100, 0, 0).Result()
	assert.NoError(t, err)
	assert.Equal(t, []string{"c", "d"}, arr)

	// 从小到大移除后两个 c d
	cache.ZAdd(zsetData)
	v, err = cache.ZRemByRank(rds.ScoreASC, -2, -1).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 2, v)
	arr, err = cache.ZMembersByScore(rds.ScoreASC, 0, 100, 0, 0).Result()
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "b"}, arr)

	// 从小到大 只保留前2个
	cache.ZAdd(zsetData)
	v, err = cache.ZRemByRank(rds.ScoreASC, 2, -1).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 2, v)
	arr, err = cache.ZMembersByScore(rds.ScoreASC, 0, 100, 0, 0).Result()
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "b"}, arr)

	// 从大到小移除前两个 d c
	cache.ZAdd(zsetData)
	v, err = cache.ZRemByRank(rds.ScoreDESC, 0, 1).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 2, v)
	arr, err = cache.ZMembersByScore(rds.ScoreASC, 0, 100, 0, 0).Result()
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "b"}, arr)

	// 从大到小移除后两个 b a
	cache.ZAdd(zsetData)
	v, err = cache.ZRemByRank(rds.ScoreDESC, -2, -1).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 2, v)
	arr, err = cache.ZMembersByScore(rds.ScoreASC, 0, 100, 0, 0).Result()
	assert.NoError(t, err)
	assert.Equal(t, []string{"c", "d"}, arr)

	// 从大到小 只保留前2个
	cache.ZAdd(zsetData)
	v, err = cache.ZRemByRank(rds.ScoreDESC, 2, -1).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 2, v)
	arr, err = cache.ZMembersByScore(rds.ScoreASC, 0, 100, 0, 0).Result()
	assert.NoError(t, err)
	assert.Equal(t, []string{"c", "d"}, arr)

	cache.Del()
}

func TestXxx(t *testing.T) {
	cache := rds.NewSortedSet[string](ctx, "key_sorted_set")
	defer cache.Del()
	// 插入/修改数据 （若仅新增/仅修改/需要查看新增数或修改数，参考说或 test/sorted_set_test.go）
	cache.ZAdd(map[string]float64{
		"a": 1.1,
		"b": 2.2,
		"c": 2.2,
		"d": 4.4,
	})
	// 获取 c 的分值
	v1, err := cache.ZScore("c").Result()
	fmt.Println(v1, err)
	// 增加 c 的分值
	v2, err := cache.ZIncrBy("c", 1).Result()
	fmt.Println(v2, err)
	// 获取 c 的分值
	v3, err := cache.ZScore("c").Result()
	fmt.Println(v3, err)
	// 获取 c 的排名 (分值从高到低)
	v4, err := cache.ZRank(rds.ScoreDESC, "c").Result()
	fmt.Println(v4, err)
	// 获取 c 的排名 (分值从低到高)
	v5, err := cache.ZRank(rds.ScoreASC, "c").Result()
	fmt.Println(v5, err)
	// 查询前两名 （分值从高到低）
	v6, err := cache.ZRangeByRank(rds.ScoreDESC, 0, 1).Result()
	fmt.Println(v6, err)
	// 查询分值在 2.2-3.3 的成员 （分值从高到低）
	v7, err := cache.ZRangeByScore(rds.ScoreDESC, 2.2, 3.3, 0, 0).Result()
	fmt.Println(v7, err)
	// 查询分值区间内有多少人
	v8, err := cache.ZCountByScore(2.2, 3.3).Result()
	fmt.Println(v8, err)
	// 移除 a
	cache.ZRem("a")
	// 移除分值最低的两个
	cache.ZRemByRank(rds.ScoreASC, 0, 1)
	// 移除分值在 2.2 - 3.3 的所有成员
	cache.ZRemByScore(2.2, 3.3)
}
