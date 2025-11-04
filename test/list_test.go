package test

import (
	"testing"

	"github.com/nezhafan/rds"
	"github.com/stretchr/testify/assert"
)

var (
	listUser1 = User{Name: "list1"}
	listUser2 = User{Name: "list2"}
)

func newList[E any]() *rds.List[E] {
	return rds.NewList[E](ctx, "list_test")
}

func TestListLPush(t *testing.T) {
	cache1 := newList[int]()

	v, err := cache1.LPush(6, 5, 4).Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(3), v)

	v, err = cache1.LPush(3, 2, 1).Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(6), v)

	v, err = cache1.LPush().Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(6), v)

	cache1.Del()

	cache2 := newList[User]()
	v, err = cache2.LPush(listUser1, listUser2).Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), v)

	cache2.Del()
}

func TestListRPush(t *testing.T) {
	cache1 := newList[int]()

	v, err := cache1.RPush(1, 2, 3).Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(3), v)

	v, err = cache1.RPush(4, 5, 6).Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(6), v)

	v, err = cache1.RPush().Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(6), v)

	cache1.Del()

	cache2 := newList[User]()
	v, err = cache2.LPush(listUser1, listUser2).Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), v)

	cache2.Del()
}

func TestListLPop(t *testing.T) {
	cache1 := newList[int]()
	cache1.RPush(1, 2)

	v1, err := cache1.LPop().Result()
	assert.NoError(t, err)
	assert.Equal(t, int(1), v1)

	v1, err = cache1.LPop().Result()
	assert.NoError(t, err)
	assert.Equal(t, int(2), v1)

	v1, err = cache1.LPop().Result()
	assert.NoError(t, err)
	assert.Equal(t, int(0), v1)

	cache1.Del()

	cache2 := newList[User]()
	cache2.RPush(listUser1, listUser2)
	v2, err := cache2.LPop().Result()
	assert.NoError(t, err)
	assert.Equal(t, listUser1.Name, v2.Name)
	cache2.Del()
}

func TestListRPop(t *testing.T) {
	cache1 := newList[int]()
	cache1.RPush(1, 2)

	v1, err := cache1.RPop().Result()
	assert.NoError(t, err)
	assert.Equal(t, int(2), v1)

	v1, err = cache1.RPop().Result()
	assert.NoError(t, err)
	assert.Equal(t, int(1), v1)

	v1, err = cache1.RPop().Result()
	assert.NoError(t, err)
	assert.Equal(t, int(0), v1)

	cache1.Del()

	cache2 := newList[User]()
	cache2.RPush(listUser1, listUser2)
	v2, err := cache2.RPop().Result()
	assert.NoError(t, err)
	assert.Equal(t, listUser2.Name, v2.Name)
	cache2.Del()
}

func TestListLIndex(t *testing.T) {
	cache := newList[int]()
	cache.RPush(1, 2)

	v1, err := cache.LIndex(0).Result()
	assert.NoError(t, err)
	assert.Equal(t, int(1), v1)

	v1, err = cache.LIndex(1).Result()
	assert.NoError(t, err)
	assert.Equal(t, int(2), v1)

	cache.Del()
}

func TestListLSet(t *testing.T) {
	cache := newList[float64]()
	cache.RPush(1)

	v, err := cache.LSet(0, 999.99).Result()
	assert.NoError(t, err)
	assert.True(t, v)

	v1, err := cache.LIndex(0).Result()
	assert.NoError(t, err)
	assert.Equal(t, float64(999.99), v1)

	cache.Del()
}

func TestListLLen(t *testing.T) {
	cache := newList[int]()
	cache.RPush(1, 3)

	v, err := cache.LLen().Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 2, v)

	cache.Del()
}

func TestListLRange(t *testing.T) {
	cache := newList[int]()
	cache.RPush(1, 2, 3)

	v, err := cache.LRange(0, -1).Result()
	assert.NoError(t, err)
	for i := range v {
		assert.EqualValues(t, i+1, v[i])
	}

	cache.Del()
}

func TestListLRem(t *testing.T) {
	cache := newList[string]()
	cache.RPush("1", "2", "3", "1", "2", "3")

	// 移除第一个1
	v, err := cache.LRem(1, "1").Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), v)

	// 移除倒数第一个2
	v, err = cache.LRem(-1, "2").Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), v)

	// 移除所有3
	v, err = cache.LRem(0, "3").Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), v)

	// 剩下1个，1个2
	v2, err := cache.LRange(0, -1).Result()
	assert.NoError(t, err)
	assert.Equal(t, []string{"2", "1"}, v2)

	cache.Del()
}

func TestListLTrim(t *testing.T) {
	cache := newList[string]()
	cache.RPush("a", "b", "c", "d", "e", "f")

	v, err := cache.LTrim(1, 3).Result()
	assert.NoError(t, err)
	assert.True(t, v)

	v2, err := cache.LRange(0, -1).Result()
	assert.NoError(t, err)
	assert.Equal(t, []string{"b", "c", "d"}, v2)

	cache.Del()
}
