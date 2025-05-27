package rds

import (
	"context"
	"time"
)

// type str struct {
// 	base
// }

type String struct {
	base
}

func NewString(ctx context.Context, key string) String {
	return String{base: newBase(ctx, key)}
}

// 必须强制设置时间，若希望永久，则使用 rds.KeepTTL
func (s *String) Set(val any, exp time.Duration) error {
	cmd := DB().Set(s.ctx, s.key, val, exp)
	s.done(cmd)
	return cmd.Err()
}

func (s *String) SetNX(val any, exp time.Duration) (success bool, err error) {
	cmd := DB().SetNX(s.ctx, s.key, val, exp)
	s.done(cmd)
	return cmd.Result()
}

func (s *String) Get() (val string, ok bool) {
	cmd := DB().Get(s.ctx, s.key)
	s.done(cmd)
	val, ok = cmd.Val(), cmd.Err() == nil
	return
}
