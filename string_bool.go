package rds

import "time"

type Bool struct {
	base
}

func NewBool(key string, ops ...Option) *Bool {
	return &Bool{base: newBase(key, ops...)}
}

func (b *Bool) Set(val bool, exp time.Duration) *BoolCmd {
	var ok string
	if val {
		ok = OK
	}
	cmd := b.db().Set(ctx, b.key, ok, exp)
	return &BoolCmd{cmd: cmd}
}

func (b *Bool) Get() *BoolCmd {
	cmd := b.db().Get(ctx, b.key)
	return &BoolCmd{cmd: cmd}
}
