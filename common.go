package rds

import (
	"context"
	"io"
	"os"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/redis/go-redis/v9"
)

const (
	OK      = "OK"
	Nil     = redis.Nil
	KeepTTL = redis.KeepTTL
)

const (
	version62 = "6.2.0"
)

var (
	// 命令和结果打印
	cmdHook func(cmd redis.Cmder)
	// 错误打印
	errorHook func(err error)
	// 开发模式
	isDebugOpen atomic.Bool
	// 日志输出
	writer io.StringWriter = os.Stdout
	// 所有key加前缀
	keyPrefix string
	// 版本号
	versionCurrent string
	// 达到 6.2.0 版本
	IsAboveVersion62 bool
)

// DEBUG模式开启/关闭，打印请求和返回。（不要在生产环境开启）
func SetDebug(isOpen bool) {
	isDebugOpen.Store(isOpen)
}

func SetWriter(w io.StringWriter) {
	if w != nil {
		writer = w
	}
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
	IsAboveVersion62 = IsAboveVersion(version62)
}

// 当前连接redis是否达到目标版本
func IsAboveVersion(targetVerion string) bool {
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
