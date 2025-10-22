package rds

import (
	"context"
	"slices"

	"github.com/redis/go-redis/v9"
)

type GeoLocation = redis.GeoLocation
type GeoPos = redis.GeoPos

type Geo struct {
	base
}

// 坐标。底层使用 zset 存储
func NewGeo(ctx context.Context, key string) *Geo {
	return &Geo{base: NewBase(ctx, key)}
}

// 添加或修改经纬度坐标 https://redis.io/docs/latest/commands/geoadd/
// 参数先经度longitude 后纬度 latitude
// 若大于 6.2.0 版本，可以设置模式: xx仅更新不新增; nx仅新增不更新; ch新增或更新;
func (g *Geo) GeoAdd(longitude, latitude float64, member string, mode string) Int64Cmd {
	args := make([]any, 0, 4)
	args = append(args, "geoadd", g.key)
	if IsReachVersion62 && slices.Contains(zsetModes, mode) {
		args = append(args, mode)
	}
	args = append(args, longitude, latitude, member)
	cmd := g.db().Do(g.ctx, args...)
	g.done(cmd)
	return newInt64Cmd(cmd)
}

// 批量添加
func (g *Geo) GeoBatchAdd(locations map[string]GeoPos, mode string) Int64Cmd {
	args := make([]any, 0, 3*len(locations)+3)
	args = append(args, "geoadd", g.key)
	if IsReachVersion62 && slices.Contains(zsetModes, mode) {
		args = append(args, mode)
	}
	for member, pos := range locations {
		args = append(args, pos.Longitude, pos.Latitude, member)
	}
	cmd := g.db().Do(g.ctx, args...)
	g.done(cmd)
	return newInt64Cmd(cmd)
}

// 获取经纬度
func (g *Geo) GeoPos(members ...string) GeoPosCmd {
	cmd := g.db().GeoPos(g.ctx, g.key, members...)
	g.done(cmd)
	return cmd
}

// 计算距离（米) 如果其中一个成员不存在则返回0
func (g *Geo) GeoDist(member1, member2 string) Float64CmdR {
	cmd := g.db().GeoDist(g.ctx, g.key, member1, member2, "m")
	g.done(cmd)
	return newFloat64CmdR(cmd)
}

/* 6.2已弃用
func (g *geo) GeoRadius() {

}

func (g *geo) GeoRadiusByMember() {

}
*/

// 依据一个经纬度点搜索半径radius米范围内的点
// func (g *Geo) GeoSearchByCoord(longitude, latitude, radius float64, count int64, withCoord bool, withDist bool, ascSort bool) *redis.GeoSearchLocationCmd {
// 	var sort string
// 	if ascSort {
// 		sort = "ASC"
// 	}
// 	query := &redis.GeoSearchLocationQuery{
// 		GeoSearchQuery: redis.GeoSearchQuery{
// 			Longitude:  longitude,
// 			Latitude:   latitude,
// 			Radius:     radius,     // 半径
// 			RadiusUnit: "m",        // 米
// 			Sort:       sort,       // 排序方式
// 			CountAny:   !ascSort,   // 是否找到count个匹配项后立即返回。这意味着返回的结果可能不是最接近指定点，但是效率高
// 			Count:      int(count), // 返回数量
// 		},
// 		WithCoord: withCoord, // 返回经纬度
// 		WithDist:  withDist,  // 计算距离
// 		WithHash:  false,
// 	}
// 	return g.db().GeoSearchLocation(g.ctx, g.key, query)
// }

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
	return g.db().GeoSearchLocation(g.ctx, g.key, query)
}

func (g *Geo) GeoDel(members ...string) Int64Cmd {
	args := sliceToAnys(members)
	cmd := g.db().ZRem(g.ctx, g.key, args...)
	g.done(cmd)
	return newInt64Cmd(cmd)
}
