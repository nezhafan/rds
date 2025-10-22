package test

import (
	"testing"
	"time"

	"github.com/nezhafan/rds"
	"github.com/stretchr/testify/assert"
)

func newBool() *rds.Bool {
	return rds.NewBool(ctx, "string_bool_test")
}

func TestBool_Set(t *testing.T) {
	cache := newBool()

	val, err := cache.Set(false, time.Second).Result()
	assert.NoError(t, err)
	assert.True(t, val)

	val, err = cache.Set(true, time.Second).Result()
	assert.NoError(t, err)
	assert.True(t, val)

	cache.Del()
}

func TestBool_SetNX(t *testing.T) {
	cache := newBool()

	v, err := cache.SetNX(false, time.Second).Result()
	assert.NoError(t, err)
	assert.True(t, v)

	v, err = cache.SetNX(false, time.Second).Result()
	assert.NoError(t, err)
	assert.False(t, v)

	cache.Del()
}

func TestBool_Get(t *testing.T) {
	cache := newBool()

	// 测试无值
	v, err := cache.Get().Result()
	assert.Nil(t, err)
	assert.Equal(t, false, v)
	exists, v, err := cache.Get().R()
	assert.Nil(t, err)
	assert.Equal(t, false, v)
	assert.False(t, exists)

	// 测试有非值
	cache.Set(true, time.Second)
	v, err = cache.Get().Result()
	assert.Nil(t, err)
	assert.Equal(t, true, v)
	exists, v, err = cache.Get().R()
	assert.Nil(t, err)
	assert.Equal(t, true, v)
	assert.True(t, exists)

	cache.Del()
}
