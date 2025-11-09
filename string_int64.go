package rds

import (
	"context"
	"time"
)

type Int64 struct {
	base
}

func NewInt64(ctx context.Context, key string) *Int64 {
	return &Int64{base: newBase(ctx, key)}
}

func (s *Int64) Set(val int64, exp time.Duration) BoolCmd {
	cmd := s.db().Set(s.ctx, s.key, val, exp)
	s.done(cmd)
	return newBoolCmd(cmd)
}

func (s *Int64) SetNX(val int64, exp time.Duration) BoolCmd {
	cmd := s.db().SetNX(s.ctx, s.key, val, exp)
	s.done(cmd)
	return newBoolCmd(cmd)
}

func (s *Int64) Get() Int64CmdR {
	cmd := s.db().Get(s.ctx, s.key)
	s.done(cmd)
	return newInt64CmdR(cmd)
}

// 建议：判断 返回值 == 增长值 时，设置一下过期时间
func (s *Int64) IncrBy(increment int64) Int64Cmd {
	cmd := s.db().IncrBy(s.ctx, s.key, increment)
	s.done(cmd)
	return newInt64Cmd(cmd)
}
