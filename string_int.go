package rds

import (
	"context"
	"time"
)

type StringInt struct {
	base
}

func NewStringInt(ctx context.Context, key string) *StringInt {
	return &StringInt{base: NewBase(ctx, key)}
}

func (s *StringInt) Set(val int64, exp time.Duration) *BoolCmd {
	cmd := s.db().Set(s.ctx, s.key, val, exp)
	s.done(cmd)
	return &BoolCmd{cmd}
}

func (s *StringInt) SetNX(val int64, exp time.Duration) *BoolCmd {
	cmd := s.db().SetNX(s.ctx, s.key, val, exp)
	s.done(cmd)
	return &BoolCmd{cmd}
}

func (s *StringInt) Get() *StringCmd[int64] {
	cmd := s.db().Get(s.ctx, s.key)
	s.done(cmd)
	return &StringCmd[int64]{cmd: cmd}
}

func (s *StringInt) IncrBy(step int64) *StringCmd[int64] {
	cmd := s.db().IncrBy(s.ctx, s.key, step)
	s.done(cmd)
	return &StringCmd[int64]{cmd: cmd}
}

func (s *StringInt) WithCmdable(cmdable Cmdable) *StringInt {
	b := s.base
	b.cmdable = cmdable
	return &StringInt{base: b}
}
