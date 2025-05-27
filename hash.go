package rds

import (
	"context"
	"time"
)

type Hash struct {
	base
}

// https://redis.io/docs/latest/commands/hdel
func NewHash(ctx context.Context, key string) Hash {
	return Hash{base: newBase(ctx, key)}
}

func (h *Hash) SubKey(subkey string) *Hash {
	nh := &Hash{base: newBase(h.ctx, h.key)}
	nh.key += ":" + subkey
	return nh
}

// 返回该字段是否为新增（修改不算新增）
func (h *Hash) HSet(field string, val any) (isNew bool, err error) {
	cmd := DB().HSet(h.ctx, h.key, field, val)
	h.done(cmd)
	return cmd.Val() == 1, cmd.Err()
}

// 返回是否设置成功
func (h *Hash) HSetNX(field string, val any) (success bool, err error) {
	cmd := DB().HSetNX(h.ctx, h.key, field, val)
	h.done(cmd)
	return cmd.Result()
}

func (h *Hash) HGet(field string) (string, error) {
	cmd := DB().HGet(h.ctx, h.key, field)
	h.done(cmd)
	return cmd.Result()
}

// 存储结构体时需要定义标签
// 结构体须带redis标签。 例如 type User struct { Id int `json:"id" redis:"id"` }
func (h *Hash) HMSet(obj any, exp time.Duration) error {
	pipe := DB().Pipeline()
	pipe.HSet(h.ctx, h.key, obj)
	if exp != KeepTTL {
		pipe.Expire(h.ctx, h.key, exp)
	}
	cmds, err := pipe.Exec(h.ctx)
	h.done(cmds[0])
	return err
}

// 参数为结构体指针
// 结构体须带redis标签。 例如 type User struct { Id int `json:"id" redis:"id"` }
func (h *Hash) HMGet(obj any, fields ...string) (exists bool) {
	if len(fields) == 0 {
		return true
	}
	cmd := DB().HMGet(h.ctx, h.key, fields...)
	h.done(cmd)
	exists = cmd.Err() == nil && len(cmd.Val()) > 0
	if exists {
		cmd.Scan(obj)
	}
	return exists
}

// 参数为结构体指针
// 结构体须带redis标签。 例如 type User struct { Id int `json:"id" redis:"id"` }
func (h *Hash) HGetAll(obj any) (exists bool) {
	cmd := DB().HGetAll(h.ctx, h.key)
	h.done(cmd)
	exists = cmd.Err() == nil && len(cmd.Val()) > 0
	if exists {
		cmd.Scan(obj)
	}
	return exists
}

func (h *Hash) HIncrBy(field string, incr int64) (int64, error) {
	cmd := DB().HIncrBy(h.ctx, h.key, field, incr)
	h.done(cmd)
	return cmd.Result()
}

func (h *Hash) HIncrByFloat(field string, incr float64) (float64, error) {
	cmd := DB().HIncrByFloat(h.ctx, h.key, field, incr)
	h.done(cmd)
	return cmd.Result()
}

func (h *Hash) HDel(fields ...string) error {
	cmd := DB().HDel(h.ctx, h.key, fields...)
	h.done(cmd)
	return cmd.Err()
}

func (h *Hash) HExists(field string) bool {
	cmd := DB().HExists(h.ctx, h.key, field)
	h.done(cmd)
	return cmd.Val()
}

func (h *Hash) HLen(field string) int {
	cmd := DB().HLen(h.ctx, h.key)
	h.done(cmd)
	return int(cmd.Val())
}
