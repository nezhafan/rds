package rds

import (
	"github.com/redis/go-redis/v9"
)

type GeoLocation = redis.GeoLocation
type GeoPos = redis.GeoPos

type Geo struct {
	base
}

func NewGeo(key string, ops ...Option) *Geo {
	return &Geo{base: newBase(key, ops...)}
}

// 添加经纬度坐标。 https://redis.io/docs/latest/commands/geoadd/
// 有效经度范围为 -180 到 180 度。
// 有效纬度范围为 -85.05112878 到 85.05112878 度。
func (g *Geo) GeoAdd(member string, longitude, latitude float64) *IntCmd {
	location := &redis.GeoLocation{
		Name:      member,
		Longitude: longitude,
		Latitude:  latitude,
	}
	cmd := g.db().GeoAdd(ctx, g.key, location)
	g.done(cmd)
	return &IntCmd{cmd: cmd}
}

// 批量添加
func (g *Geo) GeoBatchAdd(locations ...*GeoLocation) *IntCmd {
	cmd := g.db().GeoAdd(ctx, g.key, locations...)
	g.done(cmd)
	return &IntCmd{cmd: cmd}
}

// 获取经纬度
func (g *Geo) GeoPos(members ...string) *redis.GeoPosCmd {
	return g.db().GeoPos(ctx, g.key, members...)
}

// 计算距离（米) 如果其中一个成员不存在则返回0
func (g *Geo) GeoDist(member1, member2 string) *FloatCmd {
	cmd := g.db().GeoDist(ctx, g.key, member1, member2, "m")
	return &FloatCmd{cmd: cmd}
}

/* 6.2已弃用
func (g *geo) GeoRadius() {

}

func (g *geo) GeoRadiusByMember() {

}
*/

// 依据一个经纬度点搜索半径radius米范围内的点
func (g *Geo) GeoSearchByCoord(longitude, latitude, radius float64, count int64, withCoord bool, withDist bool, ascSort bool) *redis.GeoSearchLocationCmd {
	var sort string
	if ascSort {
		sort = "ASC"
	}
	query := &redis.GeoSearchLocationQuery{
		GeoSearchQuery: redis.GeoSearchQuery{
			Longitude:  longitude,
			Latitude:   latitude,
			Radius:     radius,     // 半径
			RadiusUnit: "m",        // 米
			Sort:       sort,       // 排序方式
			CountAny:   !ascSort,   // 是否找到count个匹配项后立即返回。这意味着返回的结果可能不是最接近指定点，但是效率高
			Count:      int(count), // 返回数量
		},
		WithCoord: withCoord, // 返回经纬度
		WithDist:  withDist,  // 计算距离
		WithHash:  false,
	}
	return g.db().GeoSearchLocation(ctx, g.key, query)
}

// 依据一个成员搜索半径radius米范围内的点 （结果包含其本身）
func (g *Geo) GeoSearchByMember(member string, radius float64, count int64, withCoord bool, withDist bool, ascSort bool) *redis.GeoSearchLocationCmd {
	var sort string
	if ascSort {
		sort = "ASC"
	}
	query := &redis.GeoSearchLocationQuery{
		GeoSearchQuery: redis.GeoSearchQuery{
			Member:     member,
			Radius:     radius,     // 半径
			RadiusUnit: "m",        // 米
			Sort:       sort,       // 排序方式
			CountAny:   !ascSort,   // 是否找到count个匹配项后立即返回。这意味着返回的结果可能不是最接近指定点，但是效率高
			Count:      int(count), // 返回数量
		},
		WithCoord: withCoord, // 返回经纬度
		WithDist:  withDist,  // 计算距离
		WithHash:  false,
	}
	return g.db().GeoSearchLocation(ctx, g.key, query)
}
