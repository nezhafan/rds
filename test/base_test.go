package test

import (
	"sync"
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

	// 正常过期时间
	v, err := cache.Expire(time.Minute).Result()
	assert.NoError(t, err)
	assert.True(t, v)

	n, err := cache.TTL().Result()
	assert.NoError(t, err)
	assert.EqualValues(t, time.Minute, n)

	// 永不过期
	assert.Equal(t, -1, rds.KeepTTL)
	v, err = cache.Expire(-1).Result()
	assert.NoError(t, err)
	assert.True(t, v)

	n, err = cache.TTL().Result()
	assert.NoError(t, err)
	assert.EqualValues(t, rds.KeepTTL, n)

	// 立即过期
	v, err = cache.Expire(-2).Result()
	assert.NoError(t, err)
	assert.True(t, v)

	b, err := cache.Exists().Result()
	assert.NoError(t, err)
	assert.False(t, b)

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

func TestOnceExpire(t *testing.T) {
	cache := rds.NewHashMap[int](ctx, "base_once_expire")
	cache.HSet(map[string]int{"a": 1})
	defer cache.Del()

	wg := new(sync.WaitGroup)
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cache.OnceExpire(time.Minute)
		}()
	}
	wg.Wait()

	v, err := cache.TTL().Result()
	assert.NoError(t, err)
	assert.InDelta(t, v, time.Minute, float64(time.Second))

}

func TestOnceExpireAt(t *testing.T) {
	cache := rds.NewHashMap[int](ctx, "base_once_expireat")
	cache.HSet(map[string]int{"a": 1})
	defer cache.Del()
	expAt := time.Now().Add(time.Minute)

	wg := new(sync.WaitGroup)
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cache.OnceExpireAt(expAt)
		}()
	}
	wg.Wait()

	v, err := cache.TTL().Result()
	assert.NoError(t, err)
	assert.InDelta(t, v, time.Minute, float64(time.Second))
}
