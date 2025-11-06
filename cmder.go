package rds

import (
	"cmp"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	_ convert[int64]   = toInt64
	_ convert[float64] = toFloat64
	_ convert[[]int]   = toE[[]int]
)

type convert[E any] func(redis.Cmder) E

type cmd[E any] struct {
	cmder    redis.Cmder
	convert  func(redis.Cmder) E
	isExists bool
}

func newCmd[E any](cmder redis.Cmder, cv convert[E]) cmd[E] {
	isExists := cmder.Err() == nil
	if cmder.Err() == redis.Nil {
		cmder.SetErr(nil)
	}
	c := cmd[E]{cmder: cmder, convert: cv, isExists: isExists}
	return c
}

func (c cmd[E]) Err() error {
	return c.cmder.Err()
}

func (c cmd[E]) Val() E {
	if !c.isExists {
		var e E
		return e
	}
	return c.convert(c.cmder)
}

func (c cmd[E]) Result() (E, error) {
	return c.Val(), c.Err()
}

type cmdR[E any] struct {
	cmd[E]
}

func newCmdR[E any](cmd redis.Cmder, cv convert[E]) cmdR[E] {
	return cmdR[E]{newCmd(cmd, cv)}
}

func (c cmdR[E]) R() (exists bool, val E, err error) {
	return c.isExists, c.Val(), c.Err()
}

type StringCmdR = cmdR[string]

func newStringCmdR(cmd *redis.StringCmd) StringCmdR {
	return newCmdR(cmd, toString)
}

type Int64CmdR = cmdR[int64]

type Int64Cmd = cmd[int64]

func newInt64Cmd(cmd redis.Cmder) Int64Cmd {
	return newCmd(cmd, toInt64)
}

func newInt64CmdR(cmd redis.Cmder) Int64CmdR {
	return newCmdR(cmd, toInt64)
}

type Float64Cmd = cmd[float64]

type Float64CmdR = cmdR[float64]

func newFloat64Cmd(cmd redis.Cmder) Float64Cmd {
	return newCmd(cmd, toFloat64)
}
func newFloat64CmdR(cmd redis.Cmder) Float64CmdR {
	return newCmdR(cmd, toFloat64)
}

type BoolCmd = cmd[bool]

type BoolCmdR = cmdR[bool]

func newBoolCmd(cmd redis.Cmder) BoolCmd {
	return newCmd(cmd, toBool)
}

func newBoolCmdR(cmd redis.Cmder) BoolCmdR {
	return newCmdR(cmd, toBool)
}

type DurationCmd = cmd[time.Duration]

func newDurationCmd(cmd redis.Cmder) DurationCmd {
	return newCmd(cmd, toDuration)
}

type JSONCmd[E any] struct {
	cmdR[*E]
}

func newJSONCmd[E any](cmd *redis.StringCmd) JSONCmd[E] {
	return JSONCmd[E]{newCmdR(cmd, toStructPtr[E])}
}

type AnyCmd[E any] struct {
	cmdR[E]
}

func newAnyCmd[E any](cmd *redis.StringCmd) AnyCmd[E] {
	return AnyCmd[E]{newCmdR(cmd, toE[E])}
}

type SliceCmd[E any] struct {
	cmd[[]E]
}

func newSliceCmd[E any](cmd redis.Cmder) SliceCmd[E] {
	return SliceCmd[E]{newCmd(cmd, toEs[E])}
}

type ZSliceCmd[E cmp.Ordered] struct {
	cmd[[]Z[E]]
}

func newZSliceCmd[E cmp.Ordered](cmd *redis.ZSliceCmd) ZSliceCmd[E] {
	return ZSliceCmd[E]{newCmd(cmd, toZSlice[E])}
}

type MapCmd[E cmp.Ordered] struct {
	cmd    redis.Cmder
	fields []string
}

func newMapCmd[E cmp.Ordered](cmd redis.Cmder, fields []string) MapCmd[E] {
	return MapCmd[E]{cmd: cmd, fields: fields}
}

func (c MapCmd[E]) Val() (mp map[string]E) {
	return toMap[E](c.cmd, c.fields)
}

func (c MapCmd[E]) Err() error {
	if c.cmd == nil {
		return nil
	}
	return c.cmd.Err()
}

func (c MapCmd[E]) Result() (map[string]E, error) {
	return c.Val(), c.Err()
}

type StructCmd[E any] struct {
	cmd    redis.Cmder
	fields []string
}

func newStructCmd[E any](cmd redis.Cmder, fields []string) StructCmd[E] {
	return StructCmd[E]{cmd: cmd, fields: fields}
}

func (c StructCmd[E]) Val() *E {
	_, v := toStruct[E](c.cmd, c.fields)
	return v
}

func (c StructCmd[E]) Err() error {
	return c.cmd.Err()
}

func (c StructCmd[E]) Result() (obj *E, err error) {
	return c.Val(), c.Err()
}

func (c StructCmd[E]) R() (exists bool, obj *E, err error) {
	exists, v := toStruct[E](c.cmd, c.fields)
	return exists, v, c.Err()
}

type GeoPos = redis.GeoPos

type GeoPosCmd struct {
	cmder   *redis.GeoPosCmd
	members []string
}

func newGeoPosCmd(cmder *redis.GeoPosCmd, members []string) GeoPosCmd {
	return GeoPosCmd{cmder: cmder, members: members}
}

func (g GeoPosCmd) Val() map[string]*GeoPos {
	mp := make(map[string]*GeoPos, len(g.members))
	for i, v := range g.cmder.Val() {
		mp[g.members[i]] = v
	}
	return mp
}

func (g GeoPosCmd) Err() error {
	return g.cmder.Err()
}

func (g GeoPosCmd) Result() (map[string]*GeoPos, error) {
	return g.Val(), g.Err()
}

type GeoLocation = redis.GeoLocation

type GeoLocationCmd struct {
	cmder *redis.GeoSearchLocationCmd
}

func newGeoLocationCmd(cmder *redis.GeoSearchLocationCmd) GeoLocationCmd {
	return GeoLocationCmd{cmder: cmder}
}

func (g GeoLocationCmd) Val() []GeoLocation {
	return g.cmder.Val()
}

func (g GeoLocationCmd) Err() error {
	return g.cmder.Err()
}

func (g GeoLocationCmd) Result() ([]GeoLocation, error) {
	return g.Val(), g.Err()
}
