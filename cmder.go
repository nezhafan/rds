package rds

import (
	"encoding/json"
	"reflect"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cmder[E any] interface {
	Val() E
	Err() error
	Result() (E, error)
}

type MapCmd[E Ordered] struct {
	cmd    redis.Cmder
	fields []string
}

func (mc *MapCmd[E]) Val() (mp map[string]E) {
	switch cmd := mc.cmd.(type) {
	case *redis.MapStringStringCmd:
		mp = make(map[string]E, len(cmd.Val()))
		for field, val := range cmd.Val() {
			mp[field] = stringTo[E](val)
		}
	case *redis.SliceCmd:
		mp = make(map[string]E, len(cmd.Val()))
		for i, val := range cmd.Val() {
			if s, ok := val.(string); ok {
				mp[mc.fields[i]] = stringTo[E](s)
			}
		}
	}
	return
}

func (c *MapCmd[E]) Err() error {
	return c.cmd.Err()
}

func (c *MapCmd[E]) Result() (mp map[string]E, err error) {
	return c.Val(), c.Err()
}

type StructCmd[E any] struct {
	cmd    redis.Cmder
	fields []string
}

func (c *StructCmd[E]) Val() *E {
	if c.cmd.Err() != nil {
		return nil
	}

	obj := new(E)
	switch cmd := c.cmd.(type) {
	case *redis.MapStringStringCmd:
		cmd.Scan(obj)
	case *redis.SliceCmd:
		cmd.Scan(obj)
	}
	return obj
}

func (c *StructCmd[E]) Err() error {
	return c.cmd.Err()
}

func (c *StructCmd[E]) Result() (obj *E, err error) {
	return c.Val(), c.Err()
}

type BoolCmd struct {
	cmd redis.Cmder
}

func (c *BoolCmd) Val() bool {
	switch v := c.cmd.(type) {
	case *redis.BoolCmd:
		return v.Val()
	case *redis.IntCmd:
		return v.Val() == 1
	case *redis.StatusCmd:
		return v.Val() == "OK"
	default:
		return false
	}
}

func (c *BoolCmd) Err() error {
	return c.cmd.Err()
}

func (c *BoolCmd) Result() (bool, error) {
	return c.Val(), c.Err()
}

type StringCmd[E Ordered] struct {
	cmd redis.Cmder
}

func (c *StringCmd[E]) Val() E {
	if v, ok := c.cmd.(*redis.StringCmd); ok {
		return stringTo[E](v.Val())
	}

	var zero E
	switch cmd := c.cmd.(type) {
	case *redis.IntCmd:
		reflect.ValueOf(&zero).Elem().SetInt(cmd.Val())
	case *redis.FloatCmd:
		reflect.ValueOf(&zero).Elem().SetFloat(cmd.Val())
	}

	return zero
}

func (c *StringCmd[E]) Err() error {
	if c.cmd.Err() == redis.Nil {
		return nil
	}
	return c.cmd.Err()
}

func (c *StringCmd[E]) Result() (val E, exists bool, err error) {
	val = c.Val()
	exists = c.cmd.Err() == nil
	err = c.Err()
	return
}

type StringJSONCmd[E any] struct {
	cmd redis.Cmder
}

func (c *StringJSONCmd[E]) Val() *E {
	if c.cmd.Err() != nil {
		return nil
	}

	cmd, ok := c.cmd.(*redis.StringCmd)
	if !ok {
		return nil
	}

	if len(cmd.Val()) == 0 || cmd.Val() == "null" {
		return nil
	}

	var data E
	err := json.Unmarshal([]byte(cmd.Val()), &data)
	if err != nil {
		return nil
	}

	return &data
}

func (c *StringJSONCmd[E]) Err() error {
	if c.cmd.Err() == redis.Nil {
		return nil
	}
	return c.cmd.Err()
}

func (c *StringJSONCmd[E]) Result() (data *E, exists bool, err error) {
	data = c.Val()
	exists = c.cmd.Err() == nil
	err = c.Err()
	return
}

type DurationCmd struct {
	cmd *redis.DurationCmd
}

func (c *DurationCmd) Val() time.Duration {
	return c.cmd.Val()
}

func (c *DurationCmd) Err() error {
	return c.cmd.Err()
}

func (c *DurationCmd) Result() (time.Duration, error) {
	return c.Val(), c.Err()
}

type IntCmd struct {
	cmd redis.Cmder
}

func (c *IntCmd) Val() int64 {
	switch cmd := c.cmd.(type) {
	case *redis.IntCmd:
		return cmd.Val()
	case *redis.StringCmd:
		s := cmd.Val()
		n, _ := strconv.ParseInt(s, 10, 64)
		return n
	case *redis.Cmd:
		n, _ := cmd.Int64()
		return n
	}
	return 0
}

func (c *IntCmd) Err() error {
	return c.cmd.Err()
}

func (c *IntCmd) Result() (int64, error) {
	return c.Val(), c.Err()
}

type FloatCmd struct {
	cmd redis.Cmder
}

func (c *FloatCmd) Val() float64 {
	switch cmd := c.cmd.(type) {
	case *redis.FloatCmd:
		return cmd.Val()
	case *redis.StringCmd:
		s := cmd.Val()
		n, _ := strconv.ParseFloat(s, 64)
		return n
	}
	return 0
}

func (c *FloatCmd) Err() error {
	return c.cmd.Err()
}

func (c *FloatCmd) Result() (float64, error) {
	return c.Val(), c.Err()
}

type SliceCmd[E Ordered] struct {
	cmd *redis.StringSliceCmd
}

func (c *SliceCmd[E]) Val() []E {
	return stringsToSlice[E](c.cmd.Val())
}

func (c *SliceCmd[E]) Err() error {
	return c.cmd.Err()
}

func (c *SliceCmd[E]) Result() (list []E, err error) {
	return c.Val(), c.cmd.Err()
}

type AnyCmd[E Ordered] struct {
	cmd *redis.StringCmd
}

func (c *AnyCmd[E]) Val() E {
	return stringTo[E](c.cmd.Val())
}

func (c *AnyCmd[E]) Err() error {
	return c.cmd.Err()
}

func (c *AnyCmd[E]) Result() (E, error) {
	return c.Val(), c.cmd.Err()
}

type ZSliceCmd[M Ordered] struct {
	cmd *redis.ZSliceCmd
}

func (c *ZSliceCmd[M]) Val() []Z[M] {
	list := make([]Z[M], 0, len(c.cmd.Val()))
	for _, v := range c.cmd.Val() {
		var member M
		if v, ok := v.Member.(string); ok {
			member = stringTo[M](v)
		}
		list = append(list, Z[M]{
			Score:  v.Score,
			Member: member,
		})
	}
	return list
}

func (c *ZSliceCmd[M]) Err() error {
	return c.cmd.Err()
}

func (c *ZSliceCmd[M]) Result() ([]Z[M], error) {
	return c.Val(), c.cmd.Err()
}
