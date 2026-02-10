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
	// 默认锁过期时间(秒)
	defaultExpireSecond = 60
	// 重试间隔(毫秒)，斐波那契数列
	retryIntervals = []int{5, 5, 10, 20, 30, 50, 80, 130, 210}
)

type Mutex struct {
	base
	id        string
	expSecond int
}

/*
分布式锁
- 可重入锁：不同用户请求应该每次都New而不是使用同一个Mutex实例
- 重试机制：使用斐波那契数列 5ms,5ms,10ms,20ms...210ms + 随机0到自身随机值 的间隔一直重试，直到拿到锁或上下文时间耗尽
- 锁默认过期时间为60秒，不主动续约，需要自己估计业务时常是否会超过这个值，自己去开定时器利用可重入特性再次Lock
*/
func NewMutex(ctx context.Context, key string) *Mutex {
	return &Mutex{
		base:      newBase(ctx, key),
		id:        strconv.Itoa(int(rd.Int31())),
		expSecond: defaultExpireSecond,
	}
}

func (m *Mutex) WithExpire(exp time.Duration) *Mutex {
	expSecond := int(exp.Seconds())
	if expSecond <= 0 {
		expSecond = defaultExpireSecond
	}
	m.expSecond = expSecond
	return m
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
	cmd := GetDB().Eval(m.ctx, lockScript, []string{m.Key()}, m.id, m.expSecond)
	resp, err := cmd.Result()
	ok := err == nil && resp.(string) == "OK"
	if isDebugMode {
		if ok {
			_, _ = debugWriter.WriteString(m.Key() + " " + m.id + " 加锁成功\n")
		} else {
			_, _ = debugWriter.WriteString(m.Key() + " " + m.id + " 加锁失败\n")
		}
	}
	return ok
}

// 加锁。 阻塞，定时重试。
func (m *Mutex) Lock() bool {
	var retry int
	for {
		select {
		case <-m.ctx.Done():
			return false
		default:
			if ok := m.TryLock(); ok {
				return true
			}
			milli := retryIntervals[min(retry, len(retryIntervals)-1)]
			milli += rd.Intn(milli) // 加上0到自身的随机值
			if isDebugMode {
				debugWriter.WriteString(m.Key() + " " + m.id + " " + strconv.Itoa(milli) + "毫秒后重试 \n")
			}
			time.Sleep(time.Duration(milli) * time.Millisecond)
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

// 解锁，不会误解其它实例的锁
func (m *Mutex) Unlock() {
	GetDB().Eval(m.ctx, unLockScript, []string{m.Key()}, m.id)
	if isDebugMode {
		debugWriter.WriteString(m.Key() + " " + m.id + " 解锁\n")
	}
}
