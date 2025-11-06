package test

import (
	"testing"
	"time"

	"github.com/nezhafan/rds"
	"github.com/stretchr/testify/assert"
)

var testHashUser = User{
	Name:    "Alice",
	Age:     20,
	Likes:   []string{"电影", "历史"},
	Pet:     &[]string{"狗狗"}[0],
	Guns:    []gun{{Name: "AK47", Price: 10000}},
	Nothing: "this will be ignored",
}

func newHashStruct() *rds.HashStruct[User] {
	hm := rds.NewHashStruct[User](ctx, "hash_struct_test")
	return hm
}

func TestHashStruct_HSet(t *testing.T) {
	hm := newHashStruct()

	val, err := hm.HSet(&testHashUser, time.Minute).Result()
	assert.NoError(t, err)
	assert.True(t, val)

	ttl, err := hm.TTL().Result()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, time.Minute, ttl)
	assert.LessOrEqual(t, time.Minute-time.Second*2, ttl)

	hm.Del()
}

func TestHashStruct_HGet(t *testing.T) {
	hm := newHashStruct()
	hm.HSet(&testHashUser, time.Minute)

	exists, val, err := hm.HGet("name").R()
	assert.NoError(t, err)
	assert.Equal(t, testHashUser.Name, val)
	assert.True(t, exists)
	exists, val, err = hm.HGet("nothing").R()
	assert.NoError(t, err)
	assert.Equal(t, "", val)
	assert.False(t, exists)

	hm.Del()
}

func TestHashStruct_HGetAll(t *testing.T) {
	cache := newHashStruct()

	// 无缓存
	exists, val, err := cache.HGetAll().R()
	assert.NoError(t, err)
	assert.Nil(t, val)
	assert.False(t, exists)

	// 缓存nil
	cache.HSet(nil, time.Minute)
	exists, val, err = cache.HGetAll().R()
	assert.NoError(t, err)
	assert.Nil(t, val)
	assert.True(t, exists)
	cache.Del()

	// 缓存有效值
	cache.HSet(&testHashUser, time.Minute)
	val, err = cache.HGetAll().Result()
	assert.NoError(t, err)
	assert.Equal(t, "Alice", val.Name)
	assert.Equal(t, age(20), val.Age)
	pet := val.Pet
	assert.Equal(t, "狗狗", *pet)
	assert.Equal(t, []string{"电影", "历史"}, val.Likes)
	assert.Equal(t, 1, len(val.Guns))
	assert.Empty(t, val.Nothing)

	cache.Del()
}

func TestHashStruct_HMGet(t *testing.T) {
	hm := newHashStruct()
	hm.HSet(&testHashUser, time.Minute)

	val, err := hm.HMGet("name", "age", "xxx").Result()
	assert.NoError(t, err)
	assert.Equal(t, testHashUser.Name, val.Name)
	assert.Nil(t, val.Pet)
	assert.Nil(t, val.Likes)
	assert.Nil(t, val.Guns)
	assert.Equal(t, testHashUser.Age, val.Age)

	hm.Del()
}
