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
	cache1 := newSortedSet[string]()

	v, err := cache1.ZAdd(zsetData).Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(4), v)

	cache2 := newSortedSet[int]()
	v, err = cache2.ZAdd(map[int]float64{100: 100.1, 200: 200.2, 300: 300.3}).Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(3), v)

	// 修改成功1个值
	v, err = cache2.ZAdd(map[int]float64{100: 100.1, 200: 200.3}).Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), v)

	cache1.Del()
	cache2.Del()
}

func TestSortedSet_ZIncrBy(t *testing.T) {
	cache := newSortedSet[string]()
	cache.ZAdd(zsetData)

	v, err := cache.ZIncrBy("a", 3).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 3.3, v)

	v, err = cache.ZIncrBy("c", -3).Result()
	assert.NoError(t, err)
	assert.Equal(t, 0.2999999999999998, v) // 精度丢失

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
