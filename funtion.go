package rds

import (
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type Mode int

const (
	OK      = "OK"
	Nil     = redis.Nil
	KeepTTL = redis.KeepTTL

	// DEBUG 模式
	ModeClose   Mode = 1 // 不器用
	ModeCommand Mode = 2 // 打印执行的命令
	ModeFull    Mode = 3 // 打印执行的命令和返回

	// version74 = "7.4.0"
	// version80 = "8.0.0"
)

var (
	// 所有key前缀
	allKeyPrefix = ""
	// debug模式
	debugMode = ModeClose
)

// 设置所有key的前缀
func SetPrefix(prefix string) {
	allKeyPrefix = prefix
}

// 设置DEBUG模式
func SetDebug(mode Mode) {
	debugMode = mode
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

func toJSON(data any) string {
	if data == nil {
		return "null"
	}
	b, _ := json.Marshal(data)
	return string(b)
}

func toString(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case int8:
		return strconv.Itoa(int(x))
	case int16:
		return strconv.Itoa(int(x))
	case int32:
		return strconv.Itoa(int(x))
	case int64:
		return strconv.Itoa(int(x))
	case uint8:
		return strconv.Itoa(int(x))
	case uint16:
		return strconv.Itoa(int(x))
	case uint32:
		return strconv.Itoa(int(x))
	case uint64:
		return strconv.Itoa(int(x))
	case float32:
		return strconv.FormatFloat(float64(x), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(x)
	default:
		return toJSON(x)
	}
}
