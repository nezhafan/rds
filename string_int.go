package rds

import "time"

type Int struct {
	base
}

func NewInt(key string, ops ...Option) *Int {
	return &Int{base: newBase(key, ops...)}
}

func (b *Int) Set(val int64, exp time.Duration) *BoolCmd {
	cmd := b.db().Set(ctx, b.key, val, exp)
	return &BoolCmd{cmd}
}

func (b *Int) Get() *IntCmd {
	cmd := b.db().Get(ctx, b.key)
	return &IntCmd{cmd: cmd}
}

func (b *Int) IncrBy(step int64) *IntCmd {
	cmd := b.db().IncrBy(ctx, b.key, step)
	return &IntCmd{cmd: cmd}
}
