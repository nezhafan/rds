package rdb

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	// 锁自动释放时间(秒)
	maxTimeout = "60"
	// 随机数
	rd = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type mutex struct {
	ctx context.Context
	key string
	id  string
}

func NewMutex(key string) *mutex {
	return &mutex{
		ctx: context.Background(),
		key: key,
		id:  strconv.Itoa(rd.Int()),
	}
}

func (m *mutex) WithContext(ctx context.Context) *mutex {
	return &mutex{
		ctx: ctx,
		key: m.key,
		id:  m.id,
	}
}

// 先判断是否为当前线程重复获取锁，如果是则返回OK。 (可重入锁)
// 否则尝试 SetNX 获取，成功都是返回 OK
const lockScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return "OK"
else
	return redis.call("SET", KEYS[1], ARGV[1], "NX", "EX", ARGV[2])
end`

func (m *mutex) TryLock() bool {
	keys := []string{m.key}
	resp, err := rdb.Eval(m.ctx, lockScript, keys, m.id, maxTimeout).Result()
	if err != nil && err != redis.Nil {
		return false
	}

	reply, ok := resp.(string)
	return ok && reply == "OK"
}

func (m *mutex) Lock() bool {
	for {
		if m.TryLock() {
			return true
		}
		select {
		case <-m.ctx.Done():
			return false
		default:
			retryTime := time.Duration(rd.Intn(30)+10) * time.Millisecond
			fmt.Println("重试", retryTime)
			time.Sleep(retryTime)
		}
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
