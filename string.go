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
	return &BoolCmd{cmd: cmd}
}

func (s *String) SetNX(val string, exp time.Duration) *BoolCmd {
	cmd := s.db().SetNX(ctx, s.key, val, exp)
	return &BoolCmd{cmd: cmd}
}

func (s *String) Get() *StringCmd {
	cmd := s.db().Get(ctx, s.key)
	return &StringCmd{cmd: cmd}
}
