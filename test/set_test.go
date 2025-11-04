package test

import (
	"testing"

	"github.com/nezhafan/rds"
	"github.com/stretchr/testify/assert"
)

func newSet[E any]() *rds.Set[E] {
	return rds.NewSet[E](ctx, "set_test")
}

func TestSet_SAdd(t *testing.T) {
	cache := newSet[int64]()
	v, err := cache.SAdd(1, 2, 3, 1).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 3, v)

	cache.Del()
}

func TestSet_SIsMember(t *testing.T) {
	cache := newSet[int64]()
	cache.SAdd(1, 2, 3, 1)

	v, err := cache.SIsMember(1).Result()
	assert.NoError(t, err)
	assert.True(t, v)

	v, err = cache.SIsMember(4).Result()
	assert.NoError(t, err)
	assert.False(t, v)

	cache.Del()
}

func TestSet_SMembers(t *testing.T) {
	cache := newSet[int64]()
	cache.SAdd(1, 2, 3, 1)

	v, err := cache.SMembers().Result()
	assert.NoError(t, err)
	assert.EqualValues(t, []int64{1, 2, 3}, v)

	cache.Del()

}

func TestSet_SRandMember(t *testing.T) {
	cache := newSet[int64]()
	cache.SAdd(1, 2, 3, 1)

	v, err := cache.SRandMember(2).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 2, len(v))

	cache.Del()
}

func TestSet_SCard(t *testing.T) {
	cache := newSet[int64]()
	cache.SAdd(1, 2, 3, 1)

	v, err := cache.SCard().Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 3, v)

	cache.Del()
}

func TestSet_SPop(t *testing.T) {
	cache := newSet[int64]()
	cache.SAdd(1, 2, 3, 1)

	v, err := cache.SPop(2).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 2, len(v))

	// 测试重复弹出
	n, err := cache.SCard().Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, n)

	cache.Del()
}

func TestSet_SRem(t *testing.T) {
	cache := newSet[int64]()
	cache.SAdd(1, 2, 3, 1)

	v, err := cache.SRem(1, 2).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 2, v)

	v, err = cache.SCard().Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, v)

	cache.Del()
}

func BenchmarkSet_SMembers(b *testing.B) {
	cache := newSet[User]()
	items := make([]User, 0, b.N)
	for i := 0; i < b.N; i++ {
		items = append(items, testHashUser)
	}
	cache.SAdd(items...)

	_ = cache.SMembers().Val()
	cache.Del()
}
