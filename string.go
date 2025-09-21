package rds

import (
	"context"
	"time"
)

type String struct {
	base
}

func NewString(ctx context.Context, key string) *String {
	return &String{base: NewBase(ctx, key)}
}

func (s *String) Set(val string, exp time.Duration) *BoolCmd {
	cmd := s.db().Set(s.ctx, s.key, val, exp)
	s.done(cmd)
	return &BoolCmd{cmd: cmd}
}

func (s *String) SetNX(val string, exp time.Duration) *BoolCmd {
	cmd := s.db().SetNX(s.ctx, s.key, val, exp)
	s.done(cmd)
	return &BoolCmd{cmd: cmd}
}

func (s *String) Get() *StringCmd[string] {
	cmd := s.db().Get(s.ctx, s.key)
	s.done(cmd)
	return &StringCmd[string]{cmd: cmd}
}

func (s *String) WithCmdable(cmdable Cmdable) *String {
	b := s.base
	b.cmdable = cmdable
	return &String{base: b}
}
