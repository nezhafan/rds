package test

import (
	"testing"
	"time"

	"github.com/nezhafan/rds"
	"github.com/stretchr/testify/assert"
)

func newFloat64() *rds.Float64 {
	return rds.NewFloat64(ctx, "string_float64_test")
}

func TestFloat64_Set(t *testing.T) {
	cache := newFloat64()

	val, err := cache.Set(0, time.Second).Result()
	assert.NoError(t, err)
	assert.True(t, val)

	val, err = cache.Set(1, time.Second).Result()
	assert.NoError(t, err)
	assert.True(t, val)

	cache.Del()
}

func TestFloat64_SetNX(t *testing.T) {
	cache := newFloat64()

	v, err := cache.SetNX(0, time.Second).Result()
	assert.NoError(t, err)
	assert.True(t, v)

	v, err = cache.SetNX(1, time.Second).Result()
	assert.NoError(t, err)
	assert.False(t, v)

	cache.Del()
}

func TestFloat64_Get(t *testing.T) {
	cache := newFloat64()

	var zero float64
	var notZero float64 = 0.3

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

func TestFloat64_IncrByFloat(t *testing.T) {
	cache := newFloat64()

	v, err := cache.IncrByFloat(3.3).Result()
	assert.Nil(t, err)
	assert.Equal(t, float64(3.3), v)

	v, err = cache.IncrByFloat(0).Result()
	assert.Nil(t, err)
	assert.Equal(t, float64(3.3), v)

	v, err = cache.IncrByFloat(-3.3).Result()
	assert.Nil(t, err)
	assert.Equal(t, float64(0), v)

	cache.Del()
}
