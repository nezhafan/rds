package rds

import (
	"context"
	"time"
)

type StringJSON[E any] struct {
	base
}

func NewStringJSON[E any](ctx context.Context, key string) *StringJSON[E] {
	return &StringJSON[E]{base: NewBase(ctx, key)}
}

func (s *StringJSON[E]) Set(val *E, exp time.Duration) BoolCmd {
	cmd := s.db().Set(s.ctx, s.key, any2String(val), exp)
	s.done(cmd)
	return newBoolCmd(cmd)
}

func (s *StringJSON[E]) SetNX(val *E, exp time.Duration) BoolCmd {
	cmd := s.db().SetNX(s.ctx, s.key, any2String(val), exp)
	s.done(cmd)
	return newBoolCmd(cmd)
}

func (s *StringJSON[E]) Get() JSONCmd[E] {
	cmd := s.db().Get(s.ctx, s.key)
	s.done(cmd)
	return newJSONCmd[E](cmd)
}
