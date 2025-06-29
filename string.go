package rds

import "time"

type String struct {
	base
}

func NewString(key string, ops ...Option) *String {
	return &String{base: newBase(key, ops...)}
}

func (s *String) Set(val string, exp time.Duration) *BoolCmd {
	cmd := s.db().Set(ctx, s.key, val, exp)
	s.done(cmd)
	return &BoolCmd{cmd: cmd}
}

func (s *String) SetNX(val string, exp time.Duration) *BoolCmd {
	cmd := s.db().SetNX(ctx, s.key, val, exp)
	s.done(cmd)
	return &BoolCmd{cmd: cmd}
}

func (s *String) Get() *StringCmd[string] {
	cmd := s.db().Get(ctx, s.key)
	s.done(cmd)
	return &StringCmd[string]{cmd: cmd}
}
