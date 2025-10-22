package test

import (
	"testing"
	"time"

	"github.com/nezhafan/rds"
	"github.com/stretchr/testify/assert"
)

func TestDel(t *testing.T) {
	cache := rds.NewString(ctx, "base_test")

	v, err := cache.Del().Result()
	assert.NoError(t, err)
	assert.False(t, v)

	cache.Set("ok", time.Second)
	v, err = cache.Del().Result()
	assert.NoError(t, err)
	assert.True(t, v)

	cache.Del()
}

func TestExists(t *testing.T) {
	cache := rds.NewString(ctx, "base_test")

	v, err := cache.Exists().Result()
	assert.NoError(t, err)
	assert.False(t, v)

	cache.Set("ok", time.Second)
	v, err = cache.Exists().Result()
	assert.NoError(t, err)
	assert.True(t, v)

	cache.Del()
}

func TestTTL(t *testing.T) {
	cache := rds.NewString(ctx, "base_test")

	v, err := cache.TTL().Result()
	assert.NoError(t, err)
	assert.EqualValues(t, -2, v)

	cache.Set("ok", rds.KeepTTL)
	v, err = cache.TTL().Result()
	assert.NoError(t, err)
	assert.EqualValues(t, -1, v)

	cache.Set("ok", time.Minute)
	v, err = cache.TTL().Result()
	assert.NoError(t, err)
	assert.EqualValues(t, time.Minute, v)

	cache.Del()
}

func TestExpire(t *testing.T) {
	cache := rds.NewString(ctx, "base_test")

	cache.Set("ok", time.Second)

	v, err := cache.Expire(time.Minute).Result()
	assert.NoError(t, err)
	assert.True(t, v)

	n, err := cache.TTL().Result()
	assert.NoError(t, err)
	assert.EqualValues(t, time.Minute, n)

	v, err = cache.Expire(rds.KeepTTL).Result()
	assert.NoError(t, err)
	assert.True(t, v)

	n, err = cache.TTL().Result()
	assert.NoError(t, err)
	assert.EqualValues(t, rds.KeepTTL, n)

	cache.Del()
}

func TestExpireAt(t *testing.T) {
	cache := rds.NewString(ctx, "base_test")

	cache.Set("ok", time.Second)

	now := time.Now()
	v, err := cache.ExpireAt(now.Add(time.Hour)).Result()
	assert.NoError(t, err)
	assert.True(t, v)

	n, err := cache.TTL().Result()
	assert.NoError(t, err)
	assert.InDelta(t, n, time.Hour, float64(time.Second))

	cache.Del()
}
