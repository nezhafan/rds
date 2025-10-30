package rds

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Geo struct {
	base
}

// 坐标。底层使用 zset 存储
func NewGeo(ctx context.Context, key string) *Geo {
	return &Geo{base: NewBase(ctx, key)}
}

/*
添加/修改经纬度坐标 https://redis.io/docs/latest/commands/geoadd/
注意因为存储浮点数，经常会出现显示的精度问题，计算的距离也有米级别的误差。
参数先经度longitude 后纬度 latitude
zadd key [nx|xx] [ch] score member (可选参数仅6.2.0以上支持)
  - 不传参新增或更新；nx仅新增不更新； xx仅更新不新增
  - ch返回结果是否要加上更新成功数，默认仅返回新增数。
*/
func (g *Geo) GeoAdd(locations map[string]GeoPos, params ...string) Int64Cmd {
	args := make([]any, 0, len(locations)*3+4)
	args = append(args, "geoadd", g.key)
	if IsReachVersion62 {
		if len(params) > 0 {
			args = append(args, params[0])
		}
		if len(params) > 1 {
			args = append(args, params[1])
		}
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
	return newGeoPosCmd(cmd, members)
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

type GeoQuery struct {
	WithDist  bool // 是否返回距离。默认 true
	WithCoord bool // 是否返回经纬度。默认 false
	WithHash  bool // 是否返回hash。默认 false
	CountAny  bool // 是否不严格按照距离返回。开启时，返回的结果可能不是最接近指定点，但是效率高。默认 false
}

// 依据一个经纬度点搜索其半径radius米范围内的点，count为0时返回所有。（结果包含其本身）
func (g *Geo) GeoSearchByCoord(longitude, latitude, radius float64, count int64, q *GeoQuery) GeoLocationCmd {
	withDist, withCoord, withHash, countAny := true, false, false, false
	if q != nil {
		withDist = q.WithDist
		withCoord = q.WithCoord
		withHash = q.WithHash
		countAny = q.CountAny
	}
	query := &redis.GeoSearchLocationQuery{
		GeoSearchQuery: redis.GeoSearchQuery{
			Longitude:  longitude,
			Latitude:   latitude,
			Radius:     radius,     // 半径
			RadiusUnit: "m",        // 米
			Sort:       "ASC",      // 排序方式
			CountAny:   countAny,   // 是否找到count个匹配项后立即返回。
			Count:      int(count), // 返回数量
		},
		WithCoord: withCoord, // 返回经纬度
		WithDist:  withDist,  // 计算距离
		WithHash:  withHash,
	}
	cmd := g.db().GeoSearchLocation(g.ctx, g.key, query)
	g.done(cmd)
	return newGeoLocationCmd(cmd)
}

// 依据一个成员搜索半径radius米范围内的点，count为0时返回所有。（结果包含其本身）
func (g *Geo) GeoSearchByMember(member string, radius float64, count int64, q *GeoQuery) *redis.GeoSearchLocationCmd {
	withDist, withCoord, withHash, countAny := true, false, false, false
	if q != nil {
		withDist = q.WithDist
		withCoord = q.WithCoord
		withHash = q.WithHash
		countAny = q.CountAny
	}
	query := &redis.GeoSearchLocationQuery{
		GeoSearchQuery: redis.GeoSearchQuery{
			Member:     member,
			Radius:     radius,     // 半径
			RadiusUnit: "m",        // 米
			Sort:       "ASC",      // 排序方式
			CountAny:   countAny,   // 是否找到count个匹配项后立即返回。
			Count:      int(count), // 返回数量
		},
		WithCoord: withCoord, // 返回经纬度
		WithDist:  withDist,  // 计算距离
		WithHash:  withHash,
	}
	cmd := g.db().GeoSearchLocation(g.ctx, g.key, query)
	g.done(cmd)
	return cmd
}

func (g *Geo) GeoDel(members ...string) Int64Cmd {
	args := sliceToAnys(members)
	cmd := g.db().ZRem(g.ctx, g.key, args...)
	g.done(cmd)
	return newInt64Cmd(cmd)
}
