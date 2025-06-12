package rds

import "time"

type StringInt struct {
	base
}

func NewStringInt(key string, ops ...Option) *StringInt {
	return &StringInt{base: newBase(key, ops...)}
}

func (b *StringInt) Set(val int64, exp time.Duration) *BoolCmd {
	cmd := b.db().Set(ctx, b.key, val, exp)
	return &BoolCmd{cmd}
}

func (b *StringInt) Get() *IntCmd {
	cmd := b.db().Get(ctx, b.key)
	return &IntCmd{cmd: cmd}
}

func (b *StringInt) IncrBy(step int64) *IntCmd {
	cmd := b.db().IncrBy(ctx, b.key, step)
	return &IntCmd{cmd: cmd}
}
