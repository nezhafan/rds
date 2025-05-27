package rds

import (
	"context"
	"strconv"
	"time"
)

type Int struct {
	base
}

func NewInt(ctx context.Context, key string) Int {
	return Int{base: newBase(ctx, key)}
}

func (s *Int) Set(val int, exp time.Duration) error {
	cmd := DB().Set(s.ctx, s.key, val, exp)
	s.done(cmd)
	return cmd.Err()
}

func (s *Int) Get() (val int, err error) {
	cmd := DB().Get(s.ctx, s.key)
	s.done(cmd)
	err = cmd.Err()
	if err == nil {
		val, _ = strconv.Atoi(cmd.Val())
	}
	return
}

func (s *Int) IncrBy(incr int) (val int, err error) {
	cmd := DB().IncrBy(s.ctx, s.key, int64(incr))
	s.done(cmd)
	err = cmd.Err()
	if err == nil {
		val = int(cmd.Val())
	}
	return
}
