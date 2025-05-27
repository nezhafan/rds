package rds

import (
	"context"
	"strconv"
	"time"
)

type Float struct {
	base
}

func NewFloat(ctx context.Context, key string) Float {
	return Float{base: newBase(ctx, key)}
}

func (s *Float) Set(val int, exp time.Duration) error {
	cmd := DB().Set(s.ctx, s.key, val, exp)
	s.done(cmd)
	return cmd.Err()
}

func (s *Float) Get() (val float64, err error) {
	cmd := DB().Get(s.ctx, s.key)
	s.done(cmd)
	err = cmd.Err()
	if err == nil {
		val, _ = strconv.ParseFloat(cmd.Val(), 64)
	}
	return
}

func (s *Float) IncrBy(incr float64) (val float64, err error) {
	cmd := DB().IncrByFloat(s.ctx, s.key, incr)
	s.done(cmd)
	err = cmd.Err()
	if err == nil {
		val = cmd.Val()
	}
	return
}
