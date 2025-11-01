package test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/nezhafan/rds"
	"github.com/stretchr/testify/assert"
)

func TestMutex(t *testing.T) {
	// 用sleep模拟业务耗时
	const taskCost = time.Millisecond * 12

	wg := new(sync.WaitGroup)

	start := time.Now()

	const n = 3

	// 模拟并发请求。
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(ctx, time.Second*2)
			defer cancel()
			mu := rds.NewMutex(ctx, "mutex_test")
			ok := mu.Lock()
			assert.True(t, ok)
			defer mu.Unlock()
			time.Sleep(taskCost)
		}()
	}

	wg.Wait()

	cost := time.Since(start).Milliseconds()
	assert.LessOrEqual(t, cost, taskCost*n)
}

func BenchmarkMutex(b *testing.B) {
	taskCost := time.Millisecond
	start := time.Now()
	wg := new(sync.WaitGroup)
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(ctx, time.Second*20)
			defer cancel()
			mu := rds.NewMutex(ctx, "mutex_benchmark")
			ok := mu.Lock()
			assert.True(b, ok)
			defer mu.Unlock()
			time.Sleep(taskCost)
		}()
	}
	wg.Wait()
	cost := time.Since(start).Milliseconds()
	assert.LessOrEqual(b, cost, taskCost*time.Duration(b.N))
}
