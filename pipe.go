package rds

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type pipe struct {
	ctx context.Context
	p   redis.Pipeliner
}

// func NewPipe(ctx context.Context) pipe {
// 	p := DB().Pipeline()
// 	return pipe{ctx: ctx, p: p}
// }

func (p *pipe) NewString(key string) *String {
	c := NewString(p.ctx, key)
	return c
}

func (p *pipe) NewInt64(key string) *Int64 {
	s := NewInt64(p.ctx, key)
	return s
}

func (p *pipe) NewFloat64(key string) *Float64 {
	s := NewFloat64(p.ctx, key)
	return s
}

func Pipeline(ctx context.Context, fn func(pipe)) error {
	p := DB().Pipeline()
	fn(pipe{ctx: ctx, p: p})
	_, err := p.Exec(ctx)
	return err
}

// func init() {
// 	err := Pipeline(context.Background(), func(p pipe) {
// 		p.NewString("xxx").SetNX("1", time.Second)
// 		p.NewInt64("ddd").Set(111, time.Second)

// 	})
// 	fmt.Println(err)
// }
