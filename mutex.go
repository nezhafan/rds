package redis

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	// 锁自动释放时间
	maxTimeout = strconv.Itoa(int(time.Second * 60))
	// 重试间隔
	retryTime = 10 * time.Millisecond
	// 初始时间
	initUnixNano = time.Now().UnixNano()
	// 随机数
	rd = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type mutex struct {
	ctx       context.Context
	key       string
	id        string
	retryTime time.Duration
}

func NewMutex(ctx context.Context, key string) *mutex {
	return &mutex{
		ctx:       ctx,
		key:       key,
		id:        strconv.Itoa(int(initUnixNano + rd.Int63())),
		retryTime: retryTime,
	}
}

func (m *mutex) WithRetryTime(d time.Duration) *mutex {
	m.retryTime = d
	return m
}

// 先判断是否为当前线程重复获取锁，如果是则返回OK。 (可重入锁)
// 否则尝试 SetNX 获取，成功都是返回 OK
const lockScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return "OK"
else
	return redis.call("SET", KEYS[1], ARGV[1], "NX", "PX", ARGV[2])
end`

func (m *mutex) Lock() error {
	keys := []string{m.key}
	for {
		resp, err := rdb.Eval(m.ctx, lockScript, keys, m.id, maxTimeout).Result()
		if err != nil && err != redis.Nil {
			return err
		}

		reply, ok := resp.(string)
		if ok && reply == "OK" {
			return nil
		}

		time.Sleep(m.retryTime)
	}
}

// 删除脚本。须匹配id，防止超时后另外线程获取到锁后误删。
const unLockScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
else
	return 0
end`

func (m *mutex) UnLock() {
	rdb.Eval(m.ctx, unLockScript, []string{m.key}, m.id)
}
