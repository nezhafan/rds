package test

import (
	"testing"

	"github.com/nezhafan/rds"
	"github.com/stretchr/testify/assert"
)

func newGeo() *rds.Geo {
	return rds.NewGeo(ctx, "geo_test")
}

func TestGeo_GeoAdd(t *testing.T) {
	g := newGeo()

	v, err := g.GeoAdd(map[string]rds.GeoPos{
		"Alice": {Longitude: 116.390, Latitude: 39.916},
	}).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, v)

	// 只返回新增数量。
	v, err = g.GeoAdd(map[string]rds.GeoPos{
		"Alice": {Longitude: 100.390, Latitude: 30.916},
		"1000":  {Longitude: 116.390, Latitude: 39.925},
	}).Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, v)

	// 6.2.0 版本以上模式
	if rds.IsReachVersion62 {
		// 仅新增
		v, err = g.GeoAdd(map[string]rds.GeoPos{
			"Alice": {Longitude: 101.390, Latitude: 30.916},
			"850":   {Longitude: 116.380, Latitude: 39.916},
		}, "nx").Result()
		assert.NoError(t, err)
		assert.EqualValues(t, 1, v)

		// 仅修改
		v, err = g.GeoAdd(map[string]rds.GeoPos{
			"Alice": {Longitude: 102.390, Latitude: 30.916},
			"None":  {Longitude: 116.380, Latitude: 39.916},
		}, "xx").Result()
		assert.NoError(t, err)
		assert.EqualValues(t, 0, v)
		assert.InDelta(t, 102.390, g.GeoPos("Alice").Val()["Alice"].Longitude, 0.01)
		assert.Nil(t, g.GeoPos("Alice").Val()["None"])

		// 新增或修改 (区别于默认，会返回新增+修改成功数量)
		v, err = g.GeoAdd(map[string]rds.GeoPos{
			"Alice": {Longitude: 103.390, Latitude: 30.916},
		}, "ch").Result()
		assert.NoError(t, err)
		assert.EqualValues(t, 1, v)
	}

	g.Del()
}

func TestGeo_GeoPos(t *testing.T) {
	g := newGeo()

	g.GeoAdd(map[string]rds.GeoPos{
		"Alice": {Longitude: 116.390, Latitude: 39.916},
		"Bob":   {Longitude: 116.390, Latitude: 39.925},
	})

	v, err := g.GeoPos("Alice", "Bob", "Cynthia").Result()
	assert.NoError(t, err)

	assert.InDelta(t, 116.390, v["Alice"].Longitude, 0.001)
	assert.InDelta(t, 39.916, v["Alice"].Latitude, 0.001)

	assert.InDelta(t, 116.390, v["Bob"].Longitude, 0.001)
	assert.InDelta(t, 39.925, v["Bob"].Latitude, 0.001)

	// 不存在为nil
	assert.Nil(t, v["Cynthia"])

	// 修改
	g.GeoAdd(map[string]rds.GeoPos{"Alice": {Longitude: 116.380, Latitude: 49.916}})

	v, err = g.GeoPos("Alice").Result()
	assert.NoError(t, err)
	assert.InDelta(t, 116.380, v["Alice"].Longitude, 0.001)
	assert.InDelta(t, 49.916, v["Alice"].Latitude, 0.001)

	g.Del()
}

func TestGeo_GeoDist(t *testing.T) {
	g := newGeo()

	g.GeoAdd(map[string]rds.GeoPos{
		"Alice": {Longitude: 116.390, Latitude: 39.916},
		"1000":  {Longitude: 116.390, Latitude: 39.925},
		"850":   {Longitude: 116.380, Latitude: 39.916},
	})

	// 正常计算距离 约1000米
	dist, err := g.GeoDist("Alice", "1000").Result()
	assert.NoError(t, err)
	assert.InDelta(t, 1000, dist, 5)

	// 正常计算距离 约850米
	dist, err = g.GeoDist("Alice", "850").Result()
	assert.NoError(t, err)
	assert.InDelta(t, 850, dist, 5)

	// 距离自己
	exists, dist, err := g.GeoDist("Alice", "Alice").R()
	assert.NoError(t, err)
	assert.EqualValues(t, 0., dist)
	assert.True(t, exists)

	// 距离不存在的人
	exists, dist, err = g.GeoDist("Alice", "Danny").R()
	assert.NoError(t, err)
	assert.EqualValues(t, 0., dist)
	assert.False(t, exists)

	g.Del()
}

func TestGeo_GeoSearchByCoord(t *testing.T) {
	g := newGeo()

	g.GeoAdd(map[string]rds.GeoPos{
		"Alice": {Longitude: 116.390, Latitude: 39.916},
		"1000":  {Longitude: 116.390, Latitude: 39.925},
		"850":   {Longitude: 116.380, Latitude: 39.916},
	})

	// 搜索990米内（1000米实际计算后大概是1001）
	const dist = 990.
	v, err := g.GeoSearchByCoord(116.390, 39.916, dist, 0, nil).Result()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(v))

	// 一个总是自己
	assert.EqualValues(t, "Alice", v[0].Name)
	// 可以找到850
	assert.EqualValues(t, "850", v[1].Name)
	assert.LessOrEqual(t, v[1].Dist, dist)

	// 不应该找到1000
	d := g.GeoDist("Alice", "1000").Val()
	assert.Greater(t, d, dist)

	g.Del()
}

func TestGeo_GeoSearchByMember(t *testing.T) {
	g := newGeo()

	g.GeoAdd(map[string]rds.GeoPos{
		"Alice": {Longitude: 116.390, Latitude: 39.916},
		"1000":  {Longitude: 116.390, Latitude: 39.925},
		"850":   {Longitude: 116.380, Latitude: 39.916},
	})

	// 搜索990米内（1000米实际计算后大概是1001）
	const dist = 990.
	v, err := g.GeoSearchByMember("Alice", dist, 0, nil).Result()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(v))

	// 一个总是自己
	assert.EqualValues(t, "Alice", v[0].Name)
	// 可以找到850
	assert.EqualValues(t, "850", v[1].Name)
	assert.LessOrEqual(t, v[1].Dist, dist)

	// 不应该找到1000
	d := g.GeoDist("Alice", "1000").Val()
	assert.Greater(t, d, dist)

	g.Del()
}

func TestGeo_GeoDel(t *testing.T) {
	g := newGeo()

	g.GeoAdd(map[string]rds.GeoPos{
		"Alice": {Longitude: 116.390, Latitude: 39.916},
		"1000":  {Longitude: 116.390, Latitude: 39.925},
		"850":   {Longitude: 116.380, Latitude: 39.916},
	})

	v, err := g.GeoDel("Alice", "1000").Result()
	assert.NoError(t, err)
	assert.EqualValues(t, 2, v)

	p, err := g.GeoPos("Alice", "1000", "850").Result()
	assert.NoError(t, err)
	assert.Nil(t, p["Alice"])
	assert.Nil(t, p["1000"])
	assert.NotNil(t, p["850"])

	g.Del()
}
