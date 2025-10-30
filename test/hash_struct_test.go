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

func TestHashStruct_SubKey(t *testing.T) {
	hm := newHashStruct()
	subhm := hm.SubKey(ctx, "subkey")
	assert.Equal(t, "hash_struct_test", hm.Key())
	assert.Equal(t, hm.Key()+":subkey", subhm.Key())
}

func TestHashStruct_HSetAll(t *testing.T) {
	hm := newHashStruct()

	val, err := hm.HSetAll(&testHashUser, time.Minute).Result()
	assert.NoError(t, err)
	assert.True(t, val)

	ttl, err := hm.TTL().Result()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, time.Minute, ttl)
	assert.LessOrEqual(t, time.Minute-time.Second*2, ttl)

	hm.Del()
}

func TestHashStruct_HSet(t *testing.T) {
	hm := newHashStruct()

	// 不执行
	val, err := hm.HSet(nil, time.Minute).Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(0), val)

	val, err = hm.HSet(map[string]any{"name": "Bob", "age": 21, "": ""}, time.Minute).Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(3), val)

	hm.Del()
}

func TestHashStruct_HGet(t *testing.T) {
	hm := newHashStruct()
	hm.HSetAll(&testHashUser, time.Minute)

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

func TestHashStruct_HMGet(t *testing.T) {
	hm := newHashStruct()
	hm.HSetAll(&testHashUser, time.Minute)
	hm.HSet(map[string]any{"name": "", "likes": "", "pet": "", "guns": nil}, time.Minute)

	val, err := hm.HMGet("name", "age", "xxx").Result()
	assert.NoError(t, err)
	assert.Equal(t, "", val.Name)
	assert.Nil(t, val.Pet)
	assert.Nil(t, val.Likes)
	assert.Nil(t, val.Guns)
	assert.Equal(t, testHashUser.Age, val.Age)

	hm.Del()
}

func TestHashStruct_HGetAll(t *testing.T) {
	hm := newHashStruct()

	// 无缓存
	val, err := hm.HGetAll().Result()
	assert.NoError(t, err)
	assert.Nil(t, val)

	hm.HSetAll(&testHashUser, time.Minute)

	val, err = hm.HGetAll().Result()
	assert.NoError(t, err)
	assert.Equal(t, "Alice", val.Name)
	assert.Equal(t, age(20), val.Age)
	// assert.Equal(t, "狗狗", *(val.Pet))
	assert.Equal(t, []string{"电影", "历史"}, val.Likes)
	assert.Equal(t, 1, len(val.Guns))
	assert.Empty(t, val.Nothing)

	// 缓存nil
	hm.HSetAll(nil, time.Minute)
	val, err = hm.HGetAll().Result()
	assert.NoError(t, err)
	assert.Nil(t, val)

	hm.Del()
}

func TestHashStruct_HIncrBy(t *testing.T) {
	hm := newHashStruct()
	hm.HSetAll(&testHashUser, time.Minute)

	val, err := hm.HIncrBy("money", 100).Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(100), val)

	val2, err := hm.HIncrByFloat("money", 33.3).Result()
	assert.NoError(t, err)
	assert.Equal(t, float64(133.3), val2)

	hm.Del()
}

func TestHashStruct_HDel(t *testing.T) {
	hm := newHashStruct()
	hm.HSetAll(&testHashUser, time.Minute)

	val, err := hm.HDel("likes", "age", "nothing").Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), val)

	hm.Del()
}

func TestHashStruct_HExists(t *testing.T) {
	hm := newHashStruct()
	hm.HSetAll(&testHashUser, time.Minute)

	val, err := hm.HExists("name").Result()
	assert.NoError(t, err)
	assert.True(t, val)

	val, err = hm.HExists("nothing").Result()
	assert.NoError(t, err)
	assert.False(t, val)

	hm.Del()
}
