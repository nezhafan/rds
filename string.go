package rds

import "time"

type String struct {
	base
}

func NewString(key string, ops ...Option) *String {
	return &String{base: newBase(key, ops...)}
}

func (s *String) Set(val string, exp time.Duration) (ec ErrCmd) {
	ec.cmd = s.db().Set(ctx, s.key, val, exp)
	return
}

func (s *String) SetNX(val string, exp time.Duration) (bc BoolCmd) {
	bc.cmd = s.db().SetNX(ctx, s.key, val, exp)
	return
}

func (s *String) Get() (sc StringCmd) {
	sc.cmd = s.db().Get(ctx, s.key)
	return
}
