package rds

import (
	"reflect"
	"strconv"

	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/constraints"
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
		// var val E
		// v := reflect.ValueOf(&val).Elem()

		// switch rt.Kind() {
		// case reflect.String:
		// 	v.Set(reflect.ValueOf(s).Convert(rt))
		// 	output = append(output, val)
		// case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// 	if n, err := strconv.ParseInt(s, 10, rt.Bits()); err == nil {
		// 		v.SetInt(n)
		// 		output = append(output, val)
		// 	}

		// case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// 	if n, err := strconv.ParseUint(s, 10, rt.Bits()); err == nil {
		// 		v.SetUint(n)
		// 		output = append(output, val)
		// 	}

		// case reflect.Float32, reflect.Float64:
		// 	if n, err := strconv.ParseFloat(s, rt.Bits()); err == nil {
		// 		v.SetFloat(n)
		// 		output = append(output, val)
		// 	}

		// case reflect.Bool:
		// 	if b, err := strconv.ParseBool(s); err == nil {
		// 		v.SetBool(b)
		// 		output = append(output, val)
		// 	}
		// }
	}
	return output
}
