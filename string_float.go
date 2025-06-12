package rds

import "time"

type StringFloat struct {
	base
}

func NewStringFloat(key string, ops ...Option) *StringFloat {
	return &StringFloat{base: newBase(key, ops...)}
}

func (b *StringFloat) Set(val float64, exp time.Duration) *BoolCmd {
	cmd := b.db().Set(ctx, b.key, val, exp)
	return &BoolCmd{cmd: cmd}
}

func (b *StringFloat) Get() *FloatCmd {
	cmd := b.db().Get(ctx, b.key)
	return &FloatCmd{cmd: cmd}
}

func (b *StringFloat) IncrBy(step float64) *FloatCmd {
	cmd := b.db().IncrByFloat(ctx, b.key, step)
	return &FloatCmd{cmd: cmd}
}
