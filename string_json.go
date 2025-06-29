package rds

import "time"

type StringJSON[E any] struct {
	base
}

func NewStringJSON[E any](key string, ops ...Option) *StringJSON[E] {
	return &StringJSON[E]{base: newBase(key, ops...)}
}

func (s *StringJSON[E]) Set(val *E, exp time.Duration) *BoolCmd {
	cmd := s.db().Set(ctx, s.key, toJSON(val), exp)
	s.done(cmd)
	return &BoolCmd{cmd: cmd}
}

func (s *StringJSON[E]) SetNX(val *E, exp time.Duration) *BoolCmd {
	cmd := s.db().SetNX(ctx, s.key, toJSON(val), exp)
	s.done(cmd)
	return &BoolCmd{cmd: cmd}
}

func (s *StringJSON[E]) Get() *StringJSONCmd[E] {
	cmd := s.db().Get(ctx, s.key)
	s.done(cmd)
	return &StringJSONCmd[E]{cmd: cmd}
}
