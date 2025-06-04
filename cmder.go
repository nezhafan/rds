package rds

import (
	"errors"
	"reflect"
	"strconv"

	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/constraints"
)

// type cmder[T any] struct {
//   redis.Cmder
// }

// func NewCmder[T any](cmder redis.Cmder) *Cmder[T] {
//  return &cmder[T]{cmder: cmder}
// }

// func (c *cmder[T]) Result() (T, error) {
//  return c.Cmder.(T), c.Cmder.Err()
// }

// func (c *cmder[T]) Val() T {
//  return c.Cmder.(T)
// }

// func (c *cmder[T]) Err() error {
//  return c.Cmder.Err()
// }

type MapCmd struct {
	cmd    redis.Cmder
	fields []string
}

func (mc MapCmd) Result() (mp map[string]string, err error) {
	if err = mc.cmd.Err(); err != nil {
		return
	}

	switch cmd := mc.cmd.(type) {
	case *redis.MapStringStringCmd:
		mp = cmd.Val()
	case *redis.SliceCmd:
		mp = make(map[string]string, len(cmd.Val()))
		for i, val := range cmd.Val() {
			if s, ok := val.(string); ok {
				mp[mc.fields[i]] = s
			}
		}
	}
	return
}

func (mc MapCmd) Val() map[string]string {
	v, _ := mc.Result()
	return v
}

func (mc MapCmd) Err() error {
	return mc.cmd.Err()
}

type StructCmd[T any] MapCmd

func (sc StructCmd[T]) Result() (obj *T, err error) {
	if err = sc.cmd.Err(); err != nil {
		return
	}

	obj = new(T)
	switch cmd := sc.cmd.(type) {
	case *redis.MapStringStringCmd:
		cmd.Scan(obj)
	case *redis.SliceCmd:
		cmd.Scan(obj)
	}
	return
}

func (sc StructCmd[T]) Val() *T {
	v, _ := sc.Result()
	return v
}

func (sc StructCmd[T]) Err() error {
	return sc.cmd.Err()
}

type BoolCmd struct {
	cmd redis.Cmder
}

func (c *BoolCmd) Result() (bool, error) {
	switch v := c.cmd.(type) {
	case *redis.BoolCmd:
		return v.Result()
	case *redis.IntCmd:
		return v.Val() == 1, nil
	case *redis.StringCmd:
		return v.Val() == "OK", nil
	default:
		return false, nil
	}
}

func (c *BoolCmd) Val() bool {
	v, _ := c.Result()
	return v
}

func (c *BoolCmd) Err() error {
	return c.cmd.Err()
}

type StringCmd struct {
	cmd redis.Cmder
}

func (c StringCmd) Result() (val string, exists bool, err error) {
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

func (c StringCmd) Val() string {
	v, _, _ := c.Result()
	return v
}

func (c StringCmd) Err() error {
	return c.cmd.Err()
}

type IntCmd struct {
	cmd redis.Cmder
}

func (c IntCmd) Result() (int64, error) {
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

func (c IntCmd) Val() int64 {
	v, _ := c.Result()
	return v
}

func (c IntCmd) Err() error {
	return c.cmd.Err()
}

type FloatCmd struct {
	cmd redis.Cmder
}

func (c FloatCmd) Result() (float64, error) {
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

func (c FloatCmd) Val() float64 {
	v, _ := c.Result()
	return v
}

func (c FloatCmd) Err() error {
	return c.cmd.Err()
}

type ErrCmd struct {
	cmd redis.Cmder
}

func (c ErrCmd) Err() error {
	return c.cmd.Err()
}

type SliceCmd[T constraints.Ordered] struct {
	cmd *redis.StringSliceCmd
}

func (c SliceCmd[T]) Result() (list []T, err error) {
	err = c.cmd.Err()
	if err != nil {
		list = stringsToSlice[T](c.cmd.Val())
	}
	return
}

func (c SliceCmd[T]) Val() []T {
	v, _ := c.Result()
	return v
}

func (c SliceCmd[T]) Err() error {
	return c.cmd.Err()
}

func stringTo[E constraints.Ordered](input string) E {
	var zero E
	rt := reflect.TypeOf(zero)
	switch rt.Kind() {
	case reflect.String:
		return any(input).(E)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if n, err := strconv.ParseInt(input, 10, rt.Bits()); err == nil {
			return any(n).(E)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if n, err := strconv.ParseUint(input, 10, rt.Bits()); err == nil {
			return any(n).(E)
		}
	case reflect.Float32, reflect.Float64:
		if n, err := strconv.ParseFloat(input, rt.Bits()); err == nil {
			return any(n).(E)
		}
	case reflect.Bool:
		if b, err := strconv.ParseBool(input); err == nil {
			return any(b).(E)
		}
	}
	return zero
}

func stringsToSlice[E constraints.Ordered](input []string) []E {
	if len(input) == 0 {
		return nil
	}
	output := make([]E, 0, len(input))

	for _, s := range input {
		output = append(output, stringTo[E](s))

	}
	return output
}
