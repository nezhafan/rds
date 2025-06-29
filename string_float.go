package rds

import "time"

type StringFloat struct {
	base
}

func NewStringFloat(key string, ops ...Option) *StringFloat {
	return &StringFloat{base: newBase(key, ops...)}
}

func (s *StringFloat) Set(val float64, exp time.Duration) *BoolCmd {
	cmd := s.db().Set(ctx, s.key, val, exp)
	s.done(cmd)
	return &BoolCmd{cmd: cmd}
}

func (s *StringFloat) SetNX(val float64, exp time.Duration) *BoolCmd {
	cmd := s.db().SetNX(ctx, s.key, val, exp)
	s.done(cmd)
	return &BoolCmd{cmd: cmd}
}

func (s *StringFloat) Get() *StringCmd[float64] {
	cmd := s.db().Get(ctx, s.key)
	s.done(cmd)
	return &StringCmd[float64]{cmd: cmd}
}

func (s *StringFloat) IncrBy(step float64) *StringCmd[float64] {
	cmd := s.db().IncrByFloat(ctx, s.key, step)
	s.done(cmd)
	return &StringCmd[float64]{cmd: cmd}
}
