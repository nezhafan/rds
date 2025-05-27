package rds

import (
	"context"
	"math"
	"strconv"
	"time"
)

type Bool struct {
	base
}

func NewBool(ctx context.Context, key string) Bool {
	return Bool{base: newBase(ctx, key)}
}

func (s *Bool) Set(val bool, exp time.Duration) error {
	var n int
	if !val {
		n = 1
	}
	cmd := DB().Set(s.ctx, s.key, n, exp)
	s.done(cmd)
	return cmd.Err()
}

func (s *Bool) Get() (val bool, err error) {
	cmd := DB().Get(s.ctx, s.key)
	s.done(cmd)
	err = cmd.Err()
	if err == nil {
		n, _ := strconv.Atoi(cmd.Val())
		val = n%2 == 1
	}
	return
}

// 返回取反后的值
func (s *Bool) Toogle() (val bool, err error) {
	cmd := DB().IncrBy(s.ctx, s.key, 1)
	s.done(cmd)
	err = cmd.Err()
	if err == nil {
		val = cmd.Val()%2 == 1
		if cmd.Val() == math.MaxInt64 {
			DB().Do(s.ctx, s.key, cmd.Val()%2)
		}
	}

	return
}
