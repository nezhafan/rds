package rds

import (
	"context"
)

type Bit struct {
	base
}

// 位操作，适用于表达二元情况
// 不支持设置负数，最大uint32
// 占用内存取决于设置的最大值，1000万占用1.2MB，math.MaxUint32最大占用512M内存
// https://redis.io/docs/latest/commands/setbit/
func NewBit(ctx context.Context, key string) *Bit {
	return &Bit{base: newBase(ctx, key)}
}

// 设置位，返回设置之前该位的状态
func (b *Bit) SetBit(offset uint32, val uint8) BoolCmd {
	if val != 0 {
		val = 1
	}
	cmd := b.db().SetBit(b.ctx, b.key, int64(offset), int(val))
	b.done(cmd)
	return newBoolCmd(cmd)
}

// 获取位状态
func (b *Bit) GetBit(offset uint32) Int64Cmd {
	cmd := b.db().GetBit(b.ctx, b.key, int64(offset))
	b.done(cmd)
	return newInt64Cmd(cmd)
}

// 获取范围内1的个数
func (b *Bit) BitCount(start, end int64) Int64Cmd {
	args := []any{"bitcount", b.key, start, end, "bit"}
	cmd := b.db().Do(b.ctx, args...)
	b.done(cmd)
	return newInt64Cmd(cmd)
}

// 搜索第一个0或1的位置，若不存在返回-1
func (b *Bit) BitPos(search uint8, start, end int64) Int64Cmd {
	if search != 0 {
		search = 1
	}
	args := []any{"BITPOS", b.key, search, start, end, "BIT"}
	cmd := b.db().Do(b.ctx, args...)
	b.done(cmd)
	return newInt64Cmd(cmd)
}
