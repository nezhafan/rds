package test

import (
	"sync"
	"testing"

	"github.com/nezhafan/rds"
	"github.com/stretchr/testify/assert"
)

func TestPipe(t *testing.T) {
	rds.SetDebug(false)
	const key = "test_pipe"
	wg := new(sync.WaitGroup)

	const number = 9990

	// n个并发作为干扰项
	counter := rds.NewInt64(ctx, key)
	for i := 0; i < number; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			counter.IncrBy(1)
		}(i)
	}

	// 一个pipe，内部去增加100个值，若为管道，则不会受到干扰，最后一个和第一个一定相差100
	wg.Add(1)
	go func() {
		defer wg.Done()
		var cmdStart, cmdEnd rds.Int64Cmd
		err := rds.Pipeline(ctx, func(p *rds.Pipe) {
			// 使用pipe重新声明counter
			counter2 := rds.NewPipeInt64(p, key)
			for i := 1; i <= number; i++ {
				cmd := counter2.IncrBy(1)
				if i == 1 {
					cmdStart = cmd
				}
				if i == number {
					cmdEnd = cmd
				}
			}
		})
		assert.NoError(t, err)
		// 有插队，所以差值不对
		assert.NotEqualValues(t, cmdStart.Val()+number-1, cmdEnd.Val())
	}()

	wg.Wait()

	counter.Del()
}

func TestTxPipe(t *testing.T) {
	rds.SetDebug(false)
	const key = "test_pipe"
	wg := new(sync.WaitGroup)

	const number = 9990

	// n个并发作为干扰项
	counter := rds.NewInt64(ctx, key)
	for i := 0; i < number; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			counter.IncrBy(1)
		}(i)
	}

	// 一个pipe，内部去增加100个值，若为管道，则不会受到干扰，最后一个和第一个一定相差100
	wg.Add(1)
	go func() {
		defer wg.Done()
		var cmdStart, cmdEnd rds.Int64Cmd
		err := rds.TxPipeline(ctx, func(p *rds.Pipe) {
			// 使用pipe重新声明counter
			counter2 := rds.NewPipeInt64(p, key)
			for i := 1; i <= number; i++ {
				cmd := counter2.IncrBy(1)
				if i == 1 {
					cmdStart = cmd
				}
				if i == number {
					cmdEnd = cmd
				}
			}
		})
		assert.NoError(t, err)
		// 无插队，一定是这个差值
		assert.EqualValues(t, cmdStart.Val()+number-1, cmdEnd.Val())
	}()

	wg.Wait()

	counter.Del()
}
