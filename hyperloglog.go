package rds

import "context"

type hyperloglog struct {
	base
}

// 基数统计，存在0.81%误差。
// 最大仅占用 12k 左右内存，可以统计 2^64 个元素
func NewHyperLogLog(ctx context.Context, key string) hyperloglog {
	return hyperloglog{base: NewBase(ctx, key)}
}

// 添加，至少有一个添加成功返回true，否则返回false
func (h *hyperloglog) PFAdd(vals ...any) *BoolCmd {
	cmd := h.db().PFAdd(h.ctx, h.key, vals...)
	h.done(cmd)
	return &BoolCmd{cmd: cmd}
}

// 统计数量，是存在0.81%误差的
func (h *hyperloglog) PFCount() *IntCmd {
	cmd := h.db().PFCount(h.ctx, h.key)
	h.done(cmd)
	return &IntCmd{cmd: cmd}
}

// 合并其他的
func (h *hyperloglog) PFMerge(hyperloglogs ...hyperloglog) *BoolCmd {
	keys := make([]string, len(hyperloglogs))
	for i, hl := range hyperloglogs {
		keys[i] = hl.key
	}
	cmd := h.db().PFMerge(h.ctx, h.key, keys...)
	h.done(cmd)
	return &BoolCmd{cmd: cmd}
}
