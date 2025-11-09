package rds

import (
	"context"
	"time"
)

type Float64 struct {
	base
}

func NewFloat64(ctx context.Context, key string) *Float64 {
	return &Float64{base: newBase(ctx, key)}
}

func (s *Float64) Set(val float64, exp time.Duration) BoolCmd {
	cmd := s.db().Set(s.ctx, s.key, val, exp)
	s.done(cmd)
	return newBoolCmd(cmd)
}

func (s *Float64) SetNX(val float64, exp time.Duration) BoolCmd {
	cmd := s.db().SetNX(s.ctx, s.key, val, exp)
	s.done(cmd)
	return newBoolCmd(cmd)
}

func (s *Float64) Get() Float64CmdR {
	cmd := s.db().Get(s.ctx, s.key)
	s.done(cmd)
	return newFloat64CmdR(cmd)
}

// 建议：判断 返回值 == 增长值 时，设置一下过期时间
func (s *Float64) IncrByFloat(increment float64) Float64Cmd {
	cmd := s.db().IncrByFloat(s.ctx, s.key, increment)
	s.done(cmd)
	return newFloat64Cmd(cmd)
}
