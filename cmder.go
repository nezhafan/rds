package rds

import (
	"strconv"

	"github.com/redis/go-redis/v9"
)

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
	if c.Err() != nil {
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

type StringCmd struct {
	cmd redis.Cmder
}

func (c *StringCmd) Val() string {
	if v, ok := c.cmd.(*redis.StringCmd); ok {
		return v.Val()
	}

	return ""
}

func (c *StringCmd) Err() error {
	if c.cmd.Err() == redis.Nil {
		return nil
	}
	return c.cmd.Err()
}

func (c *StringCmd) Result() (val string, exists bool, err error) {
	val = c.Val()
	exists = c.cmd.Err() == nil
	err = c.Err()
	return
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

// type RedisCmd[E any] struct {
// 	cmd redis.Cmd
// }

// func (c *RedisCmd[E]) Val() E {
// 	return c.cmd.Val().(E)
// }

// func (c *RedisCmd[E]) Err() error {
// 	return c.cmd.Err()
// }

// func (c *RedisCmd[E]) Result() (E, error) {
// 	return c.cmd.Val().(E), c.cmd.Err()
// }
