package test

import (
	"testing"
	"time"

	"github.com/nezhafan/rds"
	"github.com/stretchr/testify/assert"
)

var testJsonUser = User{Name: "Alice"}

func newStringJSON[E any]() *rds.StringJSON[E] {
	return rds.NewStringJSON[E](ctx, "string_json_test")
}

func TestStringJSON_Set(t *testing.T) {
	cache1 := newStringJSON[[]int]()
	val, err := cache1.Set(&[]int{1, 2}, time.Second).Result()
	assert.NoError(t, err)
	assert.True(t, val)
	cache1.Del()

	cache2 := newStringJSON[map[string]int]()
	val, err = cache2.Set(&(map[string]int{"a": 1}), time.Second).Result()
	assert.NoError(t, err)
	assert.True(t, val)

	cache3 := newStringJSON[User]()
	val, err = cache3.Set(&testJsonUser, time.Second).Result()
	assert.NoError(t, err)
	assert.True(t, val)

	cache4 := newStringJSON[*int]()
	val, err = cache4.Set(nil, time.Second).Result()
	assert.NoError(t, err)
	assert.True(t, val)
	cache2.Del()
}

func TestStringJSON_SetNX(t *testing.T) {
	cache := newStringJSON[User]()
	v, err := cache.SetNX(nil, time.Second).Result()
	assert.NoError(t, err)
	assert.True(t, v)

	v, err = cache.SetNX(nil, time.Second).Result()
	assert.NoError(t, err)
	assert.False(t, v)

	cache.Del()
}

func TestStringJSON_Get(t *testing.T) {
	// === 测试结构体 ===
	cache1 := newStringJSON[User]()
	// 测试无缓存
	exists, v1, err := cache1.Get().R()
	assert.Nil(t, err)
	assert.False(t, exists)

	// 测试缓存nil
	cache1.Set(nil, time.Second)
	v1, err = cache1.Get().Result()
	assert.Nil(t, err)
	assert.Equal(t, "", v1.Name)
	exists, v1, err = cache1.Get().R()
	assert.Nil(t, err)
	assert.Equal(t, "", v1.Name)
	assert.True(t, exists)

	// 测试有值结构体
	cache1.Set(&testJsonUser, time.Second)
	v1, err = cache1.Get().Result()
	assert.NoError(t, err)
	assert.EqualValues(t, testJsonUser.Name, v1.Name)
	exists, v1, err = cache1.Get().R()
	assert.NoError(t, err)
	assert.EqualValues(t, testJsonUser.Name, v1.Name)
	assert.True(t, exists)

	cache1.Del()

	// === 测试切片 ===
	cache2 := newStringJSON[[]int]()
	cache2.Set(&[]int{1, 2}, time.Second)
	v2, err := cache2.Get().Result()
	assert.NoError(t, err)
	assert.EqualValues(t, []int{1, 2}, v2)
	cache2.Del()

	// === 测试map ===
	cache3 := newStringJSON[map[string]any]()
	cache3.Set(&map[string]any{"a": 1, "b": "2"}, time.Second)
	v3, err := cache3.Get().Result()
	assert.NoError(t, err)
	assert.EqualValues(t, map[string]any{"a": 1.0, "b": "2"}, v3)
	cache3.Del()
}
