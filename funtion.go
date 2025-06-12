package rds

import (
	"reflect"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var (
	// 所有key前缀
	allKeyPrefix = ""
)

const (
	OK      = "OK"
	Nil     = redis.Nil
	KeepTTL = redis.KeepTTL

	// version74 = "7.4.0"
	// version80 = "8.0.0"
)

// 设置所有key的前缀
func SetPrefix(prefix string) {
	allKeyPrefix = prefix
}

func toAnys[E any](vals []E) []any {
	ans := make([]any, len(vals))
	for i, v := range vals {
		ans[i] = v
	}
	return ans
}

func stringTo[E Ordered](input string) E {
	var zero E
	rt := reflect.TypeOf(zero)
	switch rt.Kind() {
	case reflect.String:
		return any(input).(E)
	case reflect.Int:
		if n, err := strconv.ParseInt(input, 10, rt.Bits()); err == nil {
			return any(int(n)).(E)
		}
	case reflect.Int8:
		if n, err := strconv.ParseInt(input, 10, rt.Bits()); err == nil {
			return any(int8(n)).(E)
		}
	case reflect.Int16:
		if n, err := strconv.ParseInt(input, 10, rt.Bits()); err == nil {
			return any(int16(n)).(E)
		}
	case reflect.Int32:
		if n, err := strconv.ParseInt(input, 10, rt.Bits()); err == nil {
			return any(int32(n)).(E)
		}
	case reflect.Int64:
		if n, err := strconv.ParseInt(input, 10, rt.Bits()); err == nil {
			return any(n).(E)
		}
	case reflect.Uint:
		if n, err := strconv.ParseUint(input, 10, rt.Bits()); err == nil {
			return any(uint(n)).(E)
		}
	case reflect.Uint8:
		if n, err := strconv.ParseUint(input, 10, rt.Bits()); err == nil {
			return any(uint8(n)).(E)
		}
	case reflect.Uint16:
		if n, err := strconv.ParseUint(input, 10, rt.Bits()); err == nil {
			return any(uint16(n)).(E)
		}
	case reflect.Uint32:
		if n, err := strconv.ParseUint(input, 10, rt.Bits()); err == nil {
			return any(uint32(n)).(E)
		}
	case reflect.Uint64:
		if n, err := strconv.ParseUint(input, 10, rt.Bits()); err == nil {
			return any(n).(E)
		}
	case reflect.Float32:
		if n, err := strconv.ParseFloat(input, rt.Bits()); err == nil {
			return any(float32(n)).(E)
		}
	case reflect.Float64:
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

func stringsToSlice[E Ordered](input []string) []E {
	if len(input) == 0 {
		return nil
	}
	output := make([]E, 0, len(input))

	for _, s := range input {
		output = append(output, stringTo[E](s))
	}
	return output
}
