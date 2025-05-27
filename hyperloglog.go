package rds

import "context"

type hyperloglog struct {
	base
}

// 基数统计，存在0.81%误差。
// 最大仅占用 12k 左右内存，可以统计 2^64 个元素
func NewHyperLogLog(ctx context.Context, key string) hyperloglog {
	return hyperloglog{base: newBase(ctx, key)}
}

// 添加，至少有一个添加成功返回true，否则返回false
func (h *hyperloglog) PFAdd(vals ...any) (bool, error) {
	cmd := DB().PFAdd(h.ctx, h.key, vals...)
	h.done(cmd)
	return cmd.Val() == 1, cmd.Err()
}

// 统计数量，是存在0.81%误差的
func (h *hyperloglog) PFCount() int64 {
	cmd := DB().PFCount(h.ctx, h.key)
	h.done(cmd)
	return cmd.Val()
}

// 合并其他的
func (h *hyperloglog) PFMerge(hyperloglogs ...hyperloglog) (bool, error) {
	keys := make([]string, len(hyperloglogs))
	for i, hl := range hyperloglogs {
		keys[i] = hl.key
	}
	cmd := DB().PFMerge(h.ctx, h.key, keys...)
	h.done(cmd)
	return cmd.Val() == OK, cmd.Err()
}
