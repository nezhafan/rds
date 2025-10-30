package test

import (
	"testing"
	"time"

	"github.com/nezhafan/rds"
	"github.com/stretchr/testify/assert"
)

var testJsonUser = User{Name: "Alice"}

func newJSON[E any]() *rds.JSON[E] {
	return rds.NewJSON[E](ctx, "json_test")
}

func TestJSON_Set(t *testing.T) {
	cache1 := newJSON[[]int]()
	val, err := cache1.Set(&[]int{1, 2}, time.Second).Result()
	assert.NoError(t, err)
	assert.True(t, val)
	cache1.Del()

	cache2 := newJSON[map[string]int]()
	val, err = cache2.Set(&map[string]int{"a": 1}, time.Second).Result()
	assert.NoError(t, err)
	assert.True(t, val)
	cache2.Del()

	cache3 := newJSON[User]()
	val, err = cache3.Set(&testJsonUser, time.Second).Result()
	assert.NoError(t, err)
	assert.True(t, val)
	cache3.Del()

	cache4 := newJSON[User]()
	val, err = cache4.Set(nil, time.Second).Result()
	assert.NoError(t, err)
	assert.True(t, val)
	cache4.Del()
}

func TestJSON_SetNX(t *testing.T) {
	cache := newJSON[User]()
	v, err := cache.SetNX(nil, time.Second).Result()
	assert.NoError(t, err)
	assert.True(t, v)

	v, err = cache.SetNX(nil, time.Second).Result()
	assert.NoError(t, err)
	assert.False(t, v)

	cache.Del()
}

func TestJSON_Get(t *testing.T) {
	// === 测试结构体 ===
	cache1 := newJSON[User]()
	// 测试无缓存
	exists, v1, err := cache1.Get().R()
	assert.Nil(t, err)
	assert.EqualValues(t, User{}, v1)
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
	cache2 := newJSON[[]int]()
	cache2.Set(&[]int{1, 2}, time.Second)
	v2, err := cache2.Get().Result()
	assert.NoError(t, err)
	assert.EqualValues(t, []int{1, 2}, v2)
	cache2.Del()

	// === 测试map ===
	cache3 := newJSON[map[string]any]()
	cache3.Set(&map[string]any{"a": 1, "b": "2"}, time.Second)
	v3, err := cache3.Get().Result()
	assert.NoError(t, err)
	assert.EqualValues(t, map[string]any{"a": 1.0, "b": "2"}, v3)
	cache3.Del()
}
