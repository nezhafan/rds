package rds

import (
	"context"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/redis/go-redis/v9"
)

type Mode int

const (
	OK      = "OK"
	Nil     = redis.Nil
	KeepTTL = redis.KeepTTL
)

const (
	version74 = "7.4.0"
	version62 = "6.2.0"
)

var (
	// 命令和结果打印
	cmdHook func(cmd redis.Cmder)
	// 错误打印
	errorHook func(err error)
	// 开发模式
	isDebugMode atomic.Bool
	// 所有key加前缀
	keyPrefix string
	// 版本号
	versionCurrent string
	// 达到 6.2.0 版本
	IsReachVersion62 bool
	// 达到 7.4.0 版本
	IsReachVersion74 bool
)

// 开启DEBUG模式，打印请求和返回。（不要在生产环境开启）
func OpenDebug() {
	isDebugMode.Store(true)
}

// 自定义额外处理cmd
func SetCmdHook(fn func(cmd redis.Cmder)) {
	cmdHook = fn
}

// 自定义错误处理 (用于自定义日志打印和消息通知)
func SetErrorHook(fn func(err error)) {
	errorHook = fn
}

// 设置所有的key的前缀
func SetPrefix(prefix string) {
	keyPrefix = prefix
}

// 获取版本号
func Version() string {
	return versionCurrent
}

func initInfo() {
	info := DB().Info(context.Background(), "Server").Val()
	for _, line := range strings.Split(info, "\r\n") {
		if strings.HasPrefix(line, "redis_version:") {
			versionCurrent = strings.TrimPrefix(line, "redis_version:")
			break
		}
	}
	IsReachVersion62 = IsReachVersion(version62)
	IsReachVersion74 = IsReachVersion(version74)
}

// 当前连接redis是否达到目标版本
func IsReachVersion(targetVerion string) bool {
	cv := Version()
	if cv == targetVerion {
		return true
	}

	current := strings.Split(Version(), ".")
	target := strings.Split(targetVerion, ".")
	for i, v := range current {
		if i < len(target) && v != target[i] {
			n1, _ := strconv.Atoi(v)
			n2, _ := strconv.Atoi(target[i])
			return n1 > n2
		}
	}

	return false
}
