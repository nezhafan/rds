package rds

import (
	"cmp"
	"context"

	"github.com/redis/go-redis/v9"
)

type Pipe struct {
	ctx  context.Context
	pipe redis.Pipeliner
}

// 管道，内部命令统一发送，不保证这组命令的中间有插队。内部使用 NewPipeXXX
func Pipeline(ctx context.Context, fn func(*Pipe)) error {
	p := GetDB().Pipeline()
	fn(&Pipe{ctx: ctx, pipe: p})
	_, err := p.Exec(ctx)
	return err
}

// 事务管道，内部命令统一发送，保证这组命令执行无插队。内部使用 NewPipeXXX
func TxPipeline(ctx context.Context, fn func(*Pipe)) error {
	p := GetDB().TxPipeline()
	fn(&Pipe{ctx: ctx, pipe: p})
	_, err := p.Exec(ctx)
	return err
}

func NewPipeString(p *Pipe, key string) *String {
	s := NewString(p.ctx, key)
	s.pipe = p.pipe
	return s
}

func NewPipeInt64(p *Pipe, key string) *Int64 {
	s := NewInt64(p.ctx, key)
	s.pipe = p.pipe
	return s
}

func NewPipeFloat64(p *Pipe, key string) *Float64 {
	s := NewFloat64(p.ctx, key)
	s.pipe = p.pipe
	return s
}

func NewPipeBool(p *Pipe, key string) *Bool {
	s := NewBool(p.ctx, key)
	s.pipe = p.pipe
	return s
}

func NewPipeJSON[E any](p *Pipe, key string) *JSON[E] {
	s := NewJSON[E](p.ctx, key)
	s.pipe = p.pipe
	return s
}

func NewPipeBit(p *Pipe, key string) *Bit {
	s := NewBit(p.ctx, key)
	s.pipe = p.pipe
	return s
}

func NewPipeHashMap[E cmp.Ordered](p *Pipe, key string) *HashMap[E] {
	s := NewHashMap[E](p.ctx, key)
	s.pipe = p.pipe
	return s
}

func NewPipeHashStruct[E any](p *Pipe, key string) *HashStruct[E] {
	s := NewHashStruct[E](p.ctx, key)
	s.pipe = p.pipe
	return s
}

func NewPipeSet[E any](p *Pipe, key string) *Set[E] {
	s := NewSet[E](p.ctx, key)
	s.pipe = p.pipe
	return s
}

func NewPipeSortedSet[E cmp.Ordered](p *Pipe, key string) *SortedSet[E] {
	s := NewSortedSet[E](p.ctx, key)
	s.pipe = p.pipe
	return s
}

func NewPipeList[E any](p *Pipe, key string) *List[E] {
	s := NewList[E](p.ctx, key)
	s.pipe = p.pipe
	return s
}

func NewPipeHyperLogLog(p *Pipe, key string) *HyperLogLog {
	s := NewHyperLogLog(p.ctx, key)
	s.pipe = p.pipe
	return s
}

func NewPipeGeo(p *Pipe, key string) *Geo {
	s := NewGeo(p.ctx, key)
	s.pipe = p.pipe
	return s
}
