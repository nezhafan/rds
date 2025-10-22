package test

import (
	"testing"
	"time"

	"github.com/nezhafan/rds"
	"github.com/stretchr/testify/assert"
)

func newInt64() *rds.Int64 {
	return rds.NewInt64(ctx, "string_int64_test")
}

func TestInt64_Set(t *testing.T) {
	cache := newInt64()

	val, err := cache.Set(0, time.Second).Result()
	assert.NoError(t, err)
	assert.True(t, val)

	val, err = cache.Set(1, time.Second).Result()
	assert.NoError(t, err)
	assert.True(t, val)

	cache.Del()
}

func TestInt64_SetNX(t *testing.T) {
	cache := newInt64()

	v, err := cache.SetNX(0, time.Second).Result()
	assert.NoError(t, err)
	assert.True(t, v)

	v, err = cache.SetNX(1, time.Second).Result()
	assert.NoError(t, err)
	assert.False(t, v)

	cache.Del()
}

func TestInt64_Get(t *testing.T) {
	cache := newInt64()

	var zero int64
	var notZero int64 = 333

	// 测试无值
	v, err := cache.Get().Result()
	assert.Nil(t, err)
	assert.Equal(t, zero, v)
	exists, v, err := cache.Get().R()
	assert.Nil(t, err)
	assert.Equal(t, zero, v)
	assert.False(t, exists)

	// 测试有零值

	cache.Set(zero, time.Second)
	v, err = cache.Get().Result()
	assert.Nil(t, err)
	assert.Equal(t, zero, v)
	exists, v, err = cache.Get().R()
	assert.Nil(t, err)
	assert.Equal(t, zero, v)
	assert.True(t, exists)

	// 测试有非零值
	cache.Set(notZero, time.Second)
	v, err = cache.Get().Result()
	assert.Nil(t, err)
	assert.Equal(t, notZero, v)
	exists, v, err = cache.Get().R()
	assert.Nil(t, err)
	assert.Equal(t, notZero, v)
	assert.True(t, exists)

	cache.Del()
}

func TestInt64_IncrBy(t *testing.T) {
	cache := newInt64()

	v, err := cache.IncrBy(10000).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(10000), v)

	v, err = cache.IncrBy(0).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(10000), v)

	v, err = cache.IncrBy(-10000).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(0), v)

	cache.Del()
}
