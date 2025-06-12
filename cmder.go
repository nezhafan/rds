package rds

import (
	"errors"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type MapCmd[E Ordered] struct {
	cmd    redis.Cmder
	fields []string
}

func (mc *MapCmd[E]) Val() map[string]E {
	v, _ := mc.Result()
	return v
}

func (mc *MapCmd[E]) Err() error {
	return mc.cmd.Err()
}

func (mc *MapCmd[E]) Result() (mp map[string]E, err error) {
	if err = mc.cmd.Err(); err != nil {
		return
	}

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

type StructCmd[E any] struct {
	cmd    redis.Cmder
	fields []string
}

func (c *StructCmd[E]) Val() *E {
	v, _ := c.Result()
	return v
}

func (c *StructCmd[E]) Err() error {
	return c.cmd.Err()
}

func (c *StructCmd[E]) Result() (obj *E, err error) {
	if err = c.cmd.Err(); err != nil {
		return
	}

	obj = new(E)
	switch cmd := c.cmd.(type) {
	case *redis.MapStringStringCmd:
		cmd.Scan(obj)
	case *redis.SliceCmd:
		cmd.Scan(obj)
	}
	return
}

type BoolCmd struct {
	cmd redis.Cmder
}

func (c *BoolCmd) Val() bool {
	v, _ := c.Result()
	return v
}

func (c *BoolCmd) Err() error {
	return c.cmd.Err()
}

func (c *BoolCmd) Result() (bool, error) {
	switch v := c.cmd.(type) {
	case *redis.BoolCmd:
		return v.Result()
	case *redis.IntCmd:
		return v.Val() == 1, nil
	case *redis.StatusCmd:
		return v.Val() == "OK", nil
	default:
		return false, nil
	}
}

type StringCmd struct {
	cmd redis.Cmder
}

func (c *StringCmd) Val() string {
	v, _, _ := c.Result()
	return v
}

func (c *StringCmd) Err() error {
	return c.cmd.Err()
}

func (c *StringCmd) Result() (val string, exists bool, err error) {
	if v, ok := c.cmd.(*redis.StringCmd); ok {
		val, err = v.Result()
		exists = err == nil
		if errors.Is(err, redis.Nil) {
			err = nil
		}
		return
	}

	return
}

type IntCmd struct {
	cmd redis.Cmder
}

func (c *IntCmd) Val() int64 {
	v, _ := c.Result()
	return v
}

func (c *IntCmd) Err() error {
	return c.cmd.Err()
}

func (c *IntCmd) Result() (int64, error) {
	switch cmd := c.cmd.(type) {
	case *redis.IntCmd:
		return cmd.Result()
	case *redis.StringCmd:
		var n int64
		s, err := cmd.Result()
		if err == nil {
			n, err = strconv.ParseInt(s, 10, 64)
		}
		return n, err
	}
	return 0, nil
}

type FloatCmd struct {
	cmd redis.Cmder
}

func (c *FloatCmd) Val() float64 {
	v, _ := c.Result()
	return v
}

func (c *FloatCmd) Err() error {
	return c.cmd.Err()
}

func (c *FloatCmd) Result() (float64, error) {
	switch cmd := c.cmd.(type) {
	case *redis.FloatCmd:
		return cmd.Result()
	case *redis.StringCmd:
		var n float64
		s, err := cmd.Result()
		if err == nil {
			n, err = strconv.ParseFloat(s, 64)
		}
		return n, err
	}
	return 0, nil
}

type SliceCmd[E Ordered] struct {
	cmd *redis.StringSliceCmd
}

func (c *SliceCmd[E]) Val() []E {
	v, _ := c.Result()
	return v
}

func (c *SliceCmd[E]) Err() error {
	return c.cmd.Err()
}

func (c *SliceCmd[E]) Result() (list []E, err error) {
	return stringsToSlice[E](c.cmd.Val()), c.cmd.Err()
}

type RedisCmd[E any] struct {
	cmd redis.Cmd
}

// func (c *RedisCmd[E]) Val() E {
// 	return c.cmd.Val().(E)
// }

// func (c *RedisCmd[E]) Err() error {
// 	return c.cmd.Err()
// }

// func (c *RedisCmd[E]) Result() (E, error) {
// 	return c.cmd.Val().(E), c.cmd.Err()
// }
