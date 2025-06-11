package rds

import "time"

type Float struct {
	base
}

func NewFloat(key string, ops ...Option) *Float {
	return &Float{base: newBase(key, ops...)}
}

func (b *Float) Set(val float64, exp time.Duration) *BoolCmd {
	cmd := b.db().Set(ctx, b.key, val, exp)
	return &BoolCmd{cmd: cmd}
}

func (b *Float) Get() *FloatCmd {
	cmd := b.db().Get(ctx, b.key)
	return &FloatCmd{cmd: cmd}
}

func (b *Float) IncrBy(step float64) *FloatCmd {
	cmd := b.db().IncrByFloat(ctx, b.key, step)
	return &FloatCmd{cmd: cmd}
}
