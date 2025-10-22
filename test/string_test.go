package test

import (
	"testing"
	"time"

	"github.com/nezhafan/rds"
	"github.com/stretchr/testify/assert"
)

func newString() *rds.String {
	return rds.NewString(ctx, "string_test")
}

func TestString_Set(t *testing.T) {
	cache := newString()

	val, err := cache.Set("ok", time.Second).Result()
	assert.NoError(t, err)
	assert.True(t, val)

	val, err = cache.Set("", time.Second).Result()
	assert.NoError(t, err)
	assert.True(t, val)

	cache.Del()
}

func TestString_SetNX(t *testing.T) {
	cache := newString()

	v, err := cache.SetNX("", time.Second).Result()
	assert.NoError(t, err)
	assert.True(t, v)

	v, err = cache.SetNX("", time.Second).Result()
	assert.NoError(t, err)
	assert.False(t, v)

	cache.Del()
}

func TestString_Get(t *testing.T) {
	cache := newString()

	var zero string
	var notZero = "ok"

	// 测试无值
	v, err := cache.Get().Result()
	assert.Nil(t, err)
	assert.Equal(t, "", v)
	exists, v, err := cache.Get().R()
	assert.Nil(t, err)
	assert.Equal(t, "", v)
	assert.False(t, exists)

	// 测试有零值
	cache.Set(zero, time.Second)
	v, err = cache.Get().Result()
	assert.Nil(t, err)
	assert.Equal(t, "", v)
	exists, v, err = cache.Get().R()
	assert.Nil(t, err)
	assert.Equal(t, "", v)
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
