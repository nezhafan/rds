package rds

import (
	"context"

	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/constraints"
)

type Signed constraints.Signed

type Unsigned constraints.Unsigned

type Integer constraints.Integer

type Float constraints.Float

type Number interface {
	Integer | Float
}

type Ordered constraints.Ordered

type Cmdable interface {
	redis.Cmdable
	Do(ctx context.Context, args ...any) *redis.Cmd
}
