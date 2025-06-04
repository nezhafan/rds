package rds

import "time"

type Float struct {
	base
}

func NewFloat(key string, ops ...Option) *Float {
	return &Float{base: newBase(key, ops...)}
}

func (b *Float) Set(val float64, exp time.Duration) (c ErrCmd) {
	c.cmd = b.db().Set(ctx, b.key, val, exp)
	return
}

func (b *Float) Get() (c FloatCmd) {
	c.cmd = b.db().Get(ctx, b.key)
	return
}

func (b *Float) IncrBy(step float64) (c FloatCmd) {
	c.cmd = b.db().IncrByFloat(ctx, b.key, step)
	return
}
