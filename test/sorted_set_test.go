package test

import (
	"cmp"
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

func TestSortedSet_ZIndex(t *testing.T) {
	cache := newSortedSet[string]()
	cache.ZAdd(zsetData)

	// 从小到大排
	v, err := cache.ZIndex("a", false).Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(0), v)

	v, err = cache.ZIndex("c", false).Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), v)

	v, err = cache.ZIndex("d", false).Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(3), v)

	// 从大到小排序
	v, err = cache.ZIndex("a", true).Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(3), v)

	cache.Del()
}

func TestSortedSet_ZMembersByScore(t *testing.T) {
	cache := newSortedSet[string]()
	cache.ZAdd(zsetData)

	// 正序
	v, err := cache.ZMembersByScore(2.2, 4.4, 0, 2).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, []string{"b", "c"}, v)

	// 正序偏移
	v, err = cache.ZMembersByScore(2.2, 4.4, 1, 1).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, []string{"c"}, v)

	// 倒序
	v, err = cache.ZMembersByScore(4.4, 2.2, 0, 2).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, []string{"d", "c"}, v)

	// 倒序偏移
	v, err = cache.ZMembersByScore(4.4, 2.2, 1, 1).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, []string{"c"}, v)

	cache.Del()
}

func TestSortedSet_ZMembersByIndex(t *testing.T) {
	cache := newSortedSet[string]()
	cache.ZAdd(zsetData)

	v, err := cache.ZMembersByIndex(0, 1, false).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, []string{"a", "b"}, v)

	v, err = cache.ZMembersByIndex(0, 1, true).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, []string{"d", "c"}, v)

	// 正序偏移
	v, err = cache.ZMembersByIndex(-1, -1, false).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, []string{"d"}, v)

	cache.Del()
}

func TestSortedSet_ZRangeByScore(t *testing.T) {
	cache := newSortedSet[string]()
	cache.ZAdd(zsetData)

	v, err := cache.ZRangeByScore(2.2, 3.3, 0, 2).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, []rds.Z[string]{{Member: "b", Score: 2.2}, {Member: "c", Score: 3.3}}, v)

	v, err = cache.ZRangeByScore(3.3, 2.2, 0, 2).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, []rds.Z[string]{{Member: "c", Score: 3.3}, {Member: "b", Score: 2.2}}, v)

	cache.Del()
}

func TestSortedSet_ZRangeByIndex(t *testing.T) {
	cache := newSortedSet[string]()
	cache.ZAdd(zsetData)

	v, err := cache.ZRangeByIndex(0, 1, false).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, []rds.Z[string]{{Member: "a", Score: 0.3}, {Member: "b", Score: 2.2}}, v)

	v, err = cache.ZRangeByIndex(0, 1, true).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, []rds.Z[string]{{Member: "d", Score: 4.4}, {Member: "c", Score: 3.3}}, v)

	cache.Del()
}

func TestSortedSet_ZRem(t *testing.T) {
	cache := newSortedSet[string]()
	cache.ZAdd(zsetData)

	v, err := cache.ZRem("a", "b", "x").Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), v)

	arr, err := cache.ZMembersByScore(0, 100, 0, 0).Result()
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

	arr, err := cache.ZMembersByScore(0, 100, 0, 0).Result()
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "d"}, arr)

	cache.Del()
}

func TestSortedSet_ZRemByIndex(t *testing.T) {
	cache := newSortedSet[string]()
	cache.ZAdd(zsetData)

	v, err := cache.ZRemByIndex(1, -1).Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(3), v)

	arr, err := cache.ZMembersByScore(0, 100, 0, 0).Result()
	assert.NoError(t, err)
	assert.Equal(t, []string{"a"}, arr)

	cache.Del()
}
