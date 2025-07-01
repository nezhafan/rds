package rds

import "time"

type StringInt struct {
	base
}

func NewStringInt(key string, ops ...Option) *StringInt {
	return &StringInt{base: newBase(key, ops...)}
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
