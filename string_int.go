package rds

import "time"

type Int struct {
	base
}

func NewInt(key string, ops ...Option) *Int {
	return &Int{base: newBase(key, ops...)}
}

func (b *Int) Set(val int64, exp time.Duration) (ec ErrCmd) {
	ec.cmd = b.db().Set(ctx, b.key, val, exp)
	return
}

func (b *Int) Get() (ic IntCmd) {
	ic.cmd = b.db().Get(ctx, b.key)
	return
}

func (b *Int) IncrBy(step int64) (ic IntCmd) {
	ic.cmd = b.db().IncrBy(ctx, b.key, step)
	return
}
