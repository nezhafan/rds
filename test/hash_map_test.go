package test

import (
	"testing"

	"github.com/nezhafan/rds"
	"github.com/stretchr/testify/assert"
)

var testHaspMap = map[string]float64{"a": 3.33, "b": -100, "c": 0}

func newHashMap() *rds.HashMap[float64] {
	hm := rds.NewHashMap[float64](ctx, "hash_map_test")
	return hm
}

func TestHashMap_HSet(t *testing.T) {
	hm := newHashMap()

	// 不执行
	v, err := hm.HSet(nil).Result()
	assert.Equal(t, int64(0), v)
	assert.NoError(t, err)

	v, err = hm.HSet(testHaspMap).Result()
	assert.Equal(t, int64(3), v)
	assert.NoError(t, err)

	hm.Del()
}

func TestHashMap_HSetNX(t *testing.T) {
	hm := newHashMap()
	hm.HSet(testHaspMap)

	v, err := hm.HSetNX("a", 1.11).Result()
	assert.False(t, v)
	assert.NoError(t, err)

	a, err := hm.HGet("a").Result()
	assert.Equal(t, 3.33, a)
	assert.NoError(t, err)

	v, err = hm.HSetNX("d", 1.11).Result()
	assert.True(t, v)
	assert.NoError(t, err)

	hm.Del()
}

func TestHashMap_HGet(t *testing.T) {
	hm := newHashMap()
	hm.HSet(testHaspMap)

	exists, v, err := hm.HGet("a").R()
	assert.True(t, exists)
	assert.Equal(t, 3.33, v)
	assert.NoError(t, err)

	exists, v, err = hm.HGet("b").R()
	assert.True(t, exists)
	assert.Equal(t, -100.0, v)
	assert.NoError(t, err)

	exists, v, err = hm.HGet("c").R()
	assert.True(t, exists)
	assert.Equal(t, 0.0, v)
	assert.NoError(t, err)

	exists, v, err = hm.HGet("d").R()
	assert.False(t, exists)
	assert.Equal(t, 0.0, v)
	assert.NoError(t, err)

	v, err = hm.HGet("d").Result()
	assert.Equal(t, 0.0, v)
	assert.NoError(t, err)

	hm.Del()
}

func TestHashMap_HMGet(t *testing.T) {
	hm := newHashMap()
	hm.HSet(testHaspMap)

	val, err := hm.HMGet("a", "b", "c", "d").Result()
	assert.Equal(t, 3.33, val["a"])
	assert.Equal(t, -100., val["b"])
	assert.Equal(t, 0.0, val["c"])
	assert.Equal(t, 0.0, val["d"])
	assert.NoError(t, err)

	hm.Del()
}

func TestHashMap_HMGetAll(t *testing.T) {
	hm := newHashMap()
	hm.HSet(testHaspMap)

	val, err := hm.HGetAll().Result()
	assert.Equal(t, 3, len(val))
	assert.Equal(t, 3.33, val["a"])
	assert.Equal(t, -100., val["b"])
	assert.Equal(t, 0.0, val["c"])
	assert.Equal(t, 0.0, val["d"])
	assert.NoError(t, err)

	hm.Del()
}

func TestHashMap_HLen(t *testing.T) {
	hm := newHashMap()
	hm.HSet(testHaspMap)

	val, err := hm.HLen().Result()
	assert.Equal(t, int64(3), val)
	assert.NoError(t, err)

	hm.Del()
}

// func TestHashMap_HIncrBy(t *testing.T) {
// 	hm := newHashMap()
// 	hm.HSet(testHaspMap)

// 	val, err := hm.HIncrBy("a", -2.22).Result()
// 	assert.Equal(t, 1.11, val)
// 	assert.NoError(t, err)

// 	hm.Del()
// }

func TestHashMap_HDel(t *testing.T) {
	hm := newHashMap()
	hm.HSet(testHaspMap)

	val, err := hm.HDel("a", "b", "d").Result()
	assert.Equal(t, int64(2), val)
	assert.NoError(t, err)

	a, err := hm.HGet("a").Result()
	assert.Equal(t, 0.0, a)
	assert.NoError(t, err)

	hm.Del()
}

func TestHashMap_HExists(t *testing.T) {
	hm := newHashMap()
	hm.HSet(testHaspMap)

	val, err := hm.HExists("a").Result()
	assert.True(t, val)
	assert.NoError(t, err)

	val, err = hm.HExists("d").Result()
	assert.False(t, val)
	assert.NoError(t, err)

	hm.Del()
}
