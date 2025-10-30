package test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/nezhafan/rds"
	"github.com/stretchr/testify/assert"
)

func TestMutex_Lock(t *testing.T) {
	wg := new(sync.WaitGroup)
	wg.Add(2)

	var counter atomic.Int64
	mu := rds.NewMutex(ctx, "mutex_test", time.Millisecond*100)
	start := time.Now()

	go func() {
		if err := mu.Lock(); err != nil {
			counter.Add(1)
			return
		}
		defer mu.Unlock()
		time.Sleep(time.Millisecond * 10) // 模拟业务处理5毫秒
	}()

	go func() {
		if err := mu.Lock(); err != nil {
			counter.Add(1)
			return
		}
		defer mu.Unlock()
		time.Sleep(time.Millisecond * 5) // 模拟业务处理5毫秒
	}()

	wg.Wait()

	cost := time.Since(start).Milliseconds()
	assert.LessOrEqual(t, cost, time.Millisecond*16)
	assert.GreaterOrEqual(t, cost, time.Millisecond*15)
	assert.EqualValues(t, 1, counter.Load())
}
