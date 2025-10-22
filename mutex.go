package rds

import (
	"context"
	"math/rand"
	"strconv"
	"time"
)

var (
	// 随机数
	rd = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type Mutex struct {
	base
	id  string
	exp time.Duration
}

// ctx 不要传context.Background()，需要传一个带超时的context
// exp 传该锁的过期时间，建议 10s-60s
// exp 如果锁过期而上一个业务没处理完，则会产生问题，所以要对业务有一定预估
func NewMutex(ctx context.Context, key string, exp time.Duration) *Mutex {
	return &Mutex{
		base: NewBase(ctx, key),
		id:   strconv.Itoa(int(rd.Int31())),
		exp:  exp,
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
	resp, err := DB().Eval(m.ctx, lockScript, keys, m.id, m.exp/time.Second).Result()
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
