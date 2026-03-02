package rds

import (
	"cmp"

	"github.com/redis/go-redis/v9"
)

func NewPipeline() redis.Pipeliner {
	return GetDB().Pipeline()
}

func NewTxPipeline() redis.Pipeliner {
	return GetDB().TxPipeline()
}

func NewPipeString(p redis.Pipeliner, key string) *String {
	s := NewString(ctx, key)
	s.pipe = p
	return s
}

func NewPipeInt64(p redis.Pipeliner, key string) *Int64 {
	c := NewInt64(ctx, key)
	c.pipe = p
	return c
}

func NewPipeFloat64(p redis.Pipeliner, key string) *Float64 {
	c := NewFloat64(ctx, key)
	c.pipe = p
	return c
}

func NewPipeJSON[E any](p redis.Pipeliner, key string) *JSON[E] {
	c := NewJSON[E](ctx, key)
	c.pipe = p
	return c
}

func NewPipeBit(p redis.Pipeliner, key string) *Bit {
	c := NewBit(ctx, key)
	c.pipe = p
	return c
}

func NewPipeHashMap[E cmp.Ordered](p redis.Pipeliner, key string) *HashMap[E] {
	c := NewHashMap[E](ctx, key)
	c.pipe = p
	return c
}

func NewPipeHashStruct[E any](p redis.Pipeliner, key string) *HashStruct[E] {
	c := NewHashStruct[E](ctx, key)
	c.pipe = p
	return c
}

func NewPipeSet[E any](p redis.Pipeliner, key string) *Set[E] {
	c := NewSet[E](ctx, key)
	c.pipe = p
	return c
}

func NewPipeSortedSet[E cmp.Ordered](p redis.Pipeliner, key string) *SortedSet[E] {
	c := NewSortedSet[E](ctx, key)
	c.pipe = p
	return c
}

func NewPipeList[E any](p redis.Pipeliner, key string) *List[E] {
	c := NewList[E](ctx, key)
	c.pipe = p
	return c
}

func NewPipeHyperLogLog(p redis.Pipeliner, key string) *HyperLogLog {
	c := NewHyperLogLog(ctx, key)
	c.pipe = p
	return c
}

func NewPipeGeo(p redis.Pipeliner, key string) *Geo {
	c := NewGeo(ctx, key)
	c.pipe = p
	return c
}
