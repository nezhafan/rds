package rds

import (
	"github.com/redis/go-redis/v9"
)

func NewPipeline() redis.Pipeliner {
	return GetDB().Pipeline()
}

func NewTxPipeline() redis.Pipeliner {
	return GetDB().TxPipeline()
}
