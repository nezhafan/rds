package rds

import (
	"context"
	"time"
)

type String struct {
	base
}

func NewString(ctx context.Context, key string) *String {
	return &String{base: newBase(ctx, key)}
}

func (s *String) Set(val string, exp time.Duration) BoolCmd {
	cmd := s.db().Set(s.ctx, s.key, val, exp)
	s.done(cmd)
	return newBoolCmd(cmd)
}

func (s *String) SetNX(val string, exp time.Duration) BoolCmd {
	cmd := s.db().SetNX(s.ctx, s.key, val, exp)
	s.done(cmd)
	return newBoolCmd(cmd)
}

func (s *String) Get() StringCmdR {
	cmd := s.db().Get(s.ctx, s.key)
	s.done(cmd)
	return newStringCmdR(cmd)
}
