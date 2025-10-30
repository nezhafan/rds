package rds

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

var (
	// 随机数
	rd = rand.New(rand.NewSource(time.Now().UnixNano()))
	// 重试间隔(毫秒)，斐波那契数列
	retryIntervals = []int{10, 10, 20, 30, 50, 80, 130, 210, 340, 550}
)

type Mutex struct {
	base
	id           string
	maxExpSecond int64
}

/*
分布式锁
- 可重入锁，不同用户请求应该没次都New而不是获取同一个Mutex，因为一个Mutex对象Lock总是成功的，同一个上下文请求可以重复加锁。
- 重试机制：使用毫秒级的斐波那契数列，10ms，10ms，20ms... 最大550ms的间隔一直重试
- ctx 不要传context.Background()，需要传一个带超时的context
- maxExpSecond 传该锁的最大过期秒数，建议 10s-60s。 如果锁过期而上一个业务没处理完，则会产生问题，所以要对业务有一定预估
*/
func NewMutex(ctx context.Context, key string, maxExpSecond int64) *Mutex {
	return &Mutex{
		base:         NewBase(ctx, key),
		id:           strconv.Itoa(int(rd.Int31())),
		maxExpSecond: maxExpSecond,
	}
}

// 先判断是否为当前线程重复获取锁，如果是则返回OK。 (可重入锁)
// 否则尝试 SetNX 获取，成功都是返回 OK
const lockScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return "OK"
else
	return redis.call("SET", KEYS[1], ARGV[1], "NX", "EX", ARGV[2])
end
`

// 尝试加锁
func (m *Mutex) TryLock() bool {
	keys := []string{m.key}
	cmd := DB().Eval(m.ctx, lockScript, keys, m.id, m.maxExpSecond)
	m.done(cmd)
	resp, err := cmd.Result()
	return err == nil && resp.(string) == "OK"
}

// 加锁。 每10-20ms重试一次
func (m *Mutex) Lock() bool {
	var retry int
	for {
		select {
		case <-m.ctx.Done():
			fmt.Println("加锁超时")
			return false
		default:
			if ok := m.TryLock(); ok {
				fmt.Println("加锁成功")
				return true
			}
			milli := retryIntervals[min(retry, len(retryIntervals)-1)]
			time.Sleep(time.Duration(milli) * time.Millisecond)
			fmt.Println("等待", milli, "毫秒")
			retry++
		}
	}
}

// 删除脚本。须匹配id，防止超时后另外线程获取到锁后误删。
const unLockScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
else
	return 0
end
`

func (m *Mutex) Unlock() {
	cmd := DB().Eval(m.ctx, unLockScript, []string{m.key}, m.id)
	m.done(cmd)
	fmt.Println("解锁成功")
}
