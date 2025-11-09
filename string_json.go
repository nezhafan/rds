package rds

import (
	"context"
	"time"
)

type JSON[E any] struct {
	base
}

func NewJSON[E any](ctx context.Context, key string) *JSON[E] {
	return &JSON[E]{base: newBase(ctx, key)}
}

func (s *JSON[E]) Set(val *E, exp time.Duration) BoolCmd {
	cmd := s.db().Set(s.ctx, s.key, any2String(val), exp)
	s.done(cmd)
	return newBoolCmd(cmd)
}

func (s *JSON[E]) SetNX(val *E, exp time.Duration) BoolCmd {
	cmd := s.db().SetNX(s.ctx, s.key, any2String(val), exp)
	s.done(cmd)
	return newBoolCmd(cmd)
}

func (s *JSON[E]) Get() JSONCmd[E] {
	cmd := s.db().Get(s.ctx, s.key)
	s.done(cmd)
	return newJSONCmd[E](cmd)
}
