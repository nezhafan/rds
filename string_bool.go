package rds

import (
	"context"
	"time"
)

type Bool struct {
	base
}

func NewBool(ctx context.Context, key string) *Bool {
	return &Bool{base: newBase(ctx, key)}
}

func (s *Bool) Set(val bool, exp time.Duration) BoolCmd {
	var v string
	if val {
		v = "1"
	} else {
		v = "0"
	}
	cmd := s.db().Set(s.ctx, s.key, v, exp)
	s.done(cmd)
	return newBoolCmd(cmd)
}

func (s *Bool) SetNX(val bool, exp time.Duration) BoolCmd {
	cmd := s.db().SetNX(s.ctx, s.key, val, exp)
	s.done(cmd)
	return newBoolCmd(cmd)
}

func (s *Bool) Get() BoolCmdR {
	cmd := s.db().Get(s.ctx, s.key)
	s.done(cmd)
	return newBoolCmdR(cmd)
}
