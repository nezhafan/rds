package rds

import "context"

type HyperLogLog struct {
	base
}

// 基数统计，存在0.81%误差。
// 最大仅占用 12k 左右内存，可以统计 2^64 个元素
func NewHyperLogLog(ctx context.Context, key string) *HyperLogLog {
	return &HyperLogLog{base: NewBase(ctx, key)}
}

// 添加，至少有一个添加成功返回true，否则返回false
func (h *HyperLogLog) PFAdd(vals ...any) BoolCmd {
	cmd := h.db().PFAdd(h.ctx, h.key, vals...)
	h.done(cmd)
	return newBoolCmd(cmd)
}

// 统计数量，是存在0.81%误差的
func (h *HyperLogLog) PFCount() Int64Cmd {
	cmd := h.db().PFCount(h.ctx, h.key)
	h.done(cmd)
	return newInt64Cmd(cmd)
}

// 合并其他的
func (h *HyperLogLog) PFMerge(Hyperloglogs ...HyperLogLog) BoolCmd {
	keys := make([]string, len(Hyperloglogs))
	for i, hl := range Hyperloglogs {
		keys[i] = hl.key
	}
	cmd := h.db().PFMerge(h.ctx, h.key, keys...)
	h.done(cmd)
	return newBoolCmd(cmd)
}
