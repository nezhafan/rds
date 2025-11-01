package test

import (
	"math"
	"testing"

	"github.com/nezhafan/rds"
	"github.com/stretchr/testify/assert"
)

func newBit() *rds.Bit {
	return rds.NewBit(ctx, "bit_test")
}

func TestBit_Set(t *testing.T) {
	cache := newBit()
	// 首次设置为 0 , 返回 0
	v, err := cache.SetBit(0, 0).Result()
	assert.NoError(t, err)
	assert.False(t, v)
	// 再次设置为 1，返回 0
	v, err = cache.SetBit(0, 1).Result()
	assert.NoError(t, err)
	assert.False(t, v)

	// 再次设置为 1，返回 1
	v, err = cache.SetBit(0, 1).Result()
	assert.NoError(t, err)
	assert.True(t, v)

	cache.Del()
}

func TestBit_Get(t *testing.T) {
	cache := newBit()

	cache.SetBit(math.MaxUint32, 1)

	v, err := cache.GetBit(math.MaxUint32).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, v)

	v, err = cache.GetBit(1).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 0, v)

	cache.Del()
}

func TestBit_BitCount(t *testing.T) {
	cache := newBit()

	cache.SetBit(1, 0)
	cache.SetBit(2, 1)
	cache.SetBit(math.MaxUint32, 1)

	v, err := cache.BitCount(0, math.MaxUint32).Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), v)

	cache.Del()
}

func TestBit_BitPos(t *testing.T) {
	cache := newBit()

	cache.SetBit(1, 0)
	cache.SetBit(2, 1)
	cache.SetBit(math.MaxUint16, 1)

	v, err := cache.BitPos(1, 0, -1).Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), v)

	v, err = cache.BitPos(1, 0, 1).Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(-1), v)

	cache.Del()
}
