package rds

import "github.com/redis/go-redis/v9"

type GeoLocation = redis.GeoLocation
type GeoPos = redis.GeoPos

type geo struct {
	base
}

func NewGeo(key string) geo {
	return geo{base: newBase(key)}
}

// 添加经纬度坐标。 https://redis.io/docs/latest/commands/geoadd/
// 有效经度范围为 -180 到 180 度。
// 有效纬度范围为 -85.05112878 到 85.05112878 度。
func (g *geo) GeoAdd(member string, longitude, latitude float64) error {
	location := &redis.GeoLocation{
		Name:      member,
		Longitude: longitude,
		Latitude:  latitude,
	}
	return rdb.GeoAdd(ctx, g.key, location).Err()
}

// 批量添加
func (g *geo) GeoBatchAdd(locations ...*GeoLocation) error {
	return rdb.GeoAdd(ctx, g.key, locations...).Err()
}

// 获取经纬度
func (g *geo) GeoPos(members ...string) []*GeoPos {
	return rdb.GeoPos(ctx, g.key, members...).Val()
}

// 计算距离（米) 如果其中一个成员不存在则返回0
func (g *geo) GeoDist(member1, member2 string) float64 {
	return rdb.GeoDist(ctx, g.key, member1, member2, "m").Val()
}

/* 6.2已弃用
func (g *geo) GeoRadius() {

}

func (g *geo) GeoRadiusByMember() {

}
*/

// 依据一个经纬度点搜索半径radius米范围内的点
func (g *geo) GeoSearchByCoord(longitude, latitude, radius float64, count int64, withCoord bool, withDist bool, ascSort bool) []redis.GeoLocation {
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
	return rdb.GeoSearchLocation(ctx, g.key, query).Val()
}

// 依据一个成员搜索半径radius米范围内的点 （结果包含其本身）
func (g *geo) GeoSearchByMember(member string, radius float64, count int64, withCoord bool, withDist bool, ascSort bool) []redis.GeoLocation {
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
	return rdb.GeoSearchLocation(ctx, g.key, query).Val()
}
