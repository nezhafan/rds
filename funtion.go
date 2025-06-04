package rds

import (
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
