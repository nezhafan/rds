package rds

import (
	"cmp"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	redisTag   = "redis"
	null       = "null"
	emptyField = "__empty__"
)

func bytes2String(b []byte) string {
	return string(b)
}

func string2Bytes(s string) []byte {
	return []byte(s)
}

func structToAnys(obj any) []any {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}
	result := make([]any, 0, v.NumField()*2)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		tag := t.Field(i).Tag.Get(redisTag)
		if tag == "" || tag == "-" {
			continue
		}
		name, _, _ := strings.Cut(tag, ",")
		if name == "" {
			continue
		}
		result = append(result, name, any2String(v.Field(i).Interface()))
	}
	return result
}

func mapToAnys(vals map[string]any) []any {
	anys := make([]any, 0, len(vals)*2)
	for key, val := range vals {
		anys = append(anys, key, val)
	}
	return anys
}

func slice2Any[E any](cmder redis.Cmder) (output []E) {
	switch c := cmder.(type) {
	case *redis.StringSliceCmd:
		val := c.Val()
		if len(val) == 0 {
			return nil
		}
		output = make([]E, 0, len(val))
		for i := range val {
			output = append(output, string2Any[E](val[i]))
		}
	case *redis.SliceCmd:
		val := c.Val()
		if len(val) == 0 {
			return nil
		}
		output = make([]E, 0, len(val))
		for i := range val {
			s, ok := val[i].(string)
			if !ok {
				return nil
			}
			output = append(output, string2Any[E](s))
		}
	}
	return output
}

func sliceToAnys[E any](vals []E) []any {
	anys := make([]any, len(vals))
	for i := range vals {
		anys[i] = vals[i]
		switch x := anys[i].(type) {
		case string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float64, float32, bool:
		case json.Marshaler:
			b, _ := x.MarshalJSON()
			anys[i] = bytes2String(b)
		default:
			b, _ := json.Marshal(x)
			anys[i] = bytes2String(b)
		}
	}
	return anys
}

func string2Any[E any](s string) E {
	var e E
	if s == "" {
		return e
	}

	switch x := any(e).(type) {
	case string:
		return any(s).(E)
	case int:
		n, _ := strconv.ParseInt(s, 10, 64)
		return any(int(n)).(E)
	case int8:
		n, _ := strconv.ParseInt(s, 10, 64)
		return any(int8(n)).(E)
	case int16:
		n, _ := strconv.ParseInt(s, 10, 64)
		return any(int16(n)).(E)
	case int32:
		n, _ := strconv.ParseInt(s, 10, 64)
		return any(int32(n)).(E)
	case int64:
		n, _ := strconv.ParseInt(s, 10, 64)
		return any(n).(E)
	case uint:
		n, _ := strconv.ParseUint(s, 10, 64)
		return any(uint(n)).(E)
	case uint8:
		n, _ := strconv.ParseUint(s, 10, 64)
		return any(uint8(n)).(E)
	case uint16:
		n, _ := strconv.ParseUint(s, 10, 64)
		return any(uint16(n)).(E)
	case uint32:
		n, _ := strconv.ParseUint(s, 10, 64)
		return any(uint32(n)).(E)
	case uint64:
		n, _ := strconv.ParseUint(s, 10, 64)
		return any(n).(E)
	case float64:
		n, _ := strconv.ParseFloat(s, 64)
		return any(n).(E)
	case float32:
		n, _ := strconv.ParseFloat(s, 32)
		return any(float32(n)).(E)
	case bool:
		n, _ := strconv.ParseBool(s)
		return any(n).(E)
	case json.Unmarshaler:
		x.UnmarshalJSON(string2Bytes(s))
	default:
		// 判断是否是指针
		if reflect.ValueOf(e).Kind() == reflect.Ptr {
			e = reflect.New(reflect.ValueOf(e).Type().Elem()).Interface().(E)
		}
		json.Unmarshal(string2Bytes(s), &e)
		return e
	}
	return e
}

func any2String(v any) string {
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
	case nil:
		return null
	case json.Marshaler:
		b, _ := x.MarshalJSON()
		return bytes2String(b)
	default:
		b, _ := json.Marshal(x)
		return bytes2String(b)
	}
}

func toInt64(cmder redis.Cmder) int64 {
	switch c := cmder.(type) {
	case *redis.IntCmd:
		return c.Val()
	case *redis.StringCmd:
		n, _ := strconv.Atoi(c.Val())
		return int64(n)
	case *redis.Cmd:
		n, _ := c.Int64()
		return n
	default:
		return 0
	}
}

func toFloat64(cmder redis.Cmder) float64 {
	switch c := cmder.(type) {
	case *redis.FloatCmd:
		return c.Val()
	case *redis.StringCmd:
		n, _ := strconv.ParseFloat(c.Val(), 64)
		return n
	case *redis.Cmd:
		n, _ := c.Float64()
		return n
	default:
		return 0
	}
}

func toString(cmder redis.Cmder) string {
	switch c := cmder.(type) {
	case *redis.StringCmd:
		return c.Val()
	case *redis.Cmd:
		s, _ := c.Text()
		return s
	default:
		return ""
	}
}

func toDuration(cmder redis.Cmder) time.Duration {
	c, ok := cmder.(*redis.DurationCmd)
	if ok {
		return c.Val()
	}
	return 0
}

func toBool(cmder redis.Cmder) bool {
	switch c := cmder.(type) {
	case *redis.BoolCmd:
		return c.Val()
	case *redis.StatusCmd:
		return c.Val() == OK
	case *redis.IntCmd:
		return c.Val() > 0
	case *redis.StringCmd:
		return c.Val() == "1"
	default:
		return false
	}
}

func toAny[E any](cmder redis.Cmder) E {
	c, ok := cmder.(*redis.StringCmd)
	if !ok {
		var e E
		return e
	}
	return string2Any[E](c.Val())
}

func toMap[E cmp.Ordered](cmder redis.Cmder, fields []string) (mp map[string]E) {
	if cmder == nil || cmder.Err() != nil {
		return
	}
	switch c := cmder.(type) {
	case *redis.MapStringStringCmd:
		val := c.Val()
		mp = make(map[string]E, len(val))
		for field, v := range val {
			mp[field] = string2Any[E](v)
		}
	case *redis.SliceCmd:
		val := c.Val()
		mp = make(map[string]E, len(val))
		for i, v := range val {
			if s, ok := v.(string); ok {
				mp[fields[i]] = string2Any[E](s)
			}
		}
	}
	return
}

func toStruct[E any](cmder redis.Cmder, fields []string) *E {
	if cmder == nil {
		return nil
	}
	var mp map[string]string
	switch c := cmder.(type) {
	case *redis.MapStringStringCmd:
		if c == nil || len(c.Val()) == 0 {
			return nil
		}
		// 判断缓存空值
		if len(c.Val()) == 1 {
			if _, ok := c.Val()[emptyField]; ok {
				return nil
			}
		}
		mp = c.Val()
	case *redis.SliceCmd:
		if c == nil || len(c.Val()) == 0 {
			return nil
		}
		mp = make(map[string]string, len(c.Val()))
		for i, val := range c.Val() {
			if s, ok := val.(string); ok {
				mp[fields[i]] = s
			}
		}
	default:
		return nil
	}
	obj := new(E)
	if len(mp) == 0 {
		return obj
	}
	v := reflect.ValueOf(obj).Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.CanSet() {
			continue
		}
		// 通过标签获取key
		tag := t.Field(i).Tag.Get(redisTag)
		if tag == "" || tag == "-" {
			continue
		}
		key, _, _ := strings.Cut(tag, ",")
		if key == "" {
			continue
		}
		// 是否取值
		val, ok := mp[key]
		if !ok || val == "" {
			continue
		}

		kind := field.Kind()
		if kind == reflect.Ptr {
			if field.IsNil() {
				field.Set(reflect.New(field.Type().Elem()))
			}
			field = field.Elem()
			kind = field.Kind()
		}
		switch kind {
		case reflect.String:
			field.SetString(val)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, _ := strconv.ParseInt(val, 10, 64)
			field.Set(reflect.ValueOf(n).Convert(field.Type()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			n, _ := strconv.ParseUint(val, 10, 64)
			field.Set(reflect.ValueOf(n).Convert(field.Type()))
		case reflect.Float32, reflect.Float64:
			n, _ := strconv.ParseFloat(val, field.Type().Bits())
			field.Set(reflect.ValueOf(n).Convert(field.Type()))
		case reflect.Bool:
			b, _ := strconv.ParseBool(val)
			field.Set(reflect.ValueOf(b).Convert(field.Type()))
		default:
			if field.CanAddr() {
				json.Unmarshal([]byte(val), field.Addr().Interface())
			}
		}
	}
	return obj
}

func toZSlice[E cmp.Ordered](cmder redis.Cmder) []Z[E] {
	c, ok := cmder.(*redis.ZSliceCmd)
	if !ok {
		return nil
	}
	list := make([]Z[E], 0, len(c.Val()))
	for _, v := range c.Val() {
		var member E
		if v, ok := v.Member.(string); ok {
			member = string2Any[E](v)
		}
		list = append(list, Z[E]{
			Score:  v.Score,
			Member: member,
		})
	}
	return list
}
