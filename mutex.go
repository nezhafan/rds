package rds

import (
	"context"
	"math/rand"
	"strconv"
	"time"
)

var (
	// 锁自动释放时间(秒)
	mutexTimeout = "30"
	// 随机数
	rd = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type Mutex struct {
	ctx context.Context
	key string
	id  string
}

func NewMutex(ctx context.Context, key string) *Mutex {
	if len(allKeyPrefix) > 0 {
		key = allKeyPrefix + ":mutex:" + key
	} else {
		key = "mutex:" + key
	}
	return &Mutex{
		ctx: ctx,
		key: key,
		id:  strconv.Itoa(rd.Int()),
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

// 尝试加锁
func (m *Mutex) TryLock() bool {
	keys := []string{m.key}
	resp, err := DB().Eval(m.ctx, lockScript, keys, m.id, mutexTimeout).Result()
	return err == nil && resp.(string) == "OK"
}

// 加锁。 每10-20ms重试一次
func (m *Mutex) Lock() error {
	for {
		if m.TryLock() {
			return nil
		}
		select {
		case <-m.ctx.Done():
			return context.DeadlineExceeded
		default:
			retryTime := time.Duration(rd.Intn(10)+10) * time.Millisecond
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

func (m *Mutex) Unlock() {
	DB().Eval(m.ctx, unLockScript, []string{m.key}, m.id)
}
