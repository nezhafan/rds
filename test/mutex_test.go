package test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/nezhafan/rds"
	"github.com/stretchr/testify/assert"
)

func TestMutex_Lock(t *testing.T) {
	// 用sleep模拟业务耗时
	const taskCost = time.Millisecond * 20

	wg := new(sync.WaitGroup)

	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	start := time.Now()

	// 模拟并发请求。
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu := rds.NewMutex(ctx, "mutex_test", 10)
			ok := mu.Lock()
			assert.True(t, ok)
			defer mu.Unlock()
			time.Sleep(taskCost)
		}()
	}

	wg.Wait()

	cost := time.Since(start).Milliseconds()
	assert.LessOrEqual(t, cost, taskCost*2)
}
