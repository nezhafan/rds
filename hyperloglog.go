package rds

type hyperloglog struct {
	base
}

// 基数统计，存在0.81%误差。
// 最大仅占用 12k 左右内存，可以统计 2^64 个元素
func NewHyperLogLog(key string, ops ...Option) hyperloglog {
	return hyperloglog{base: newBase(key, ops...)}
}

// 添加，至少有一个添加成功返回true，否则返回false
func (h *hyperloglog) PFAdd(vals ...any) *BoolCmd {
	cmd := h.db().PFAdd(ctx, h.key, vals...)
	return &BoolCmd{cmd: cmd}
}

// 统计数量，是存在0.81%误差的
func (h *hyperloglog) PFCount() *IntCmd {
	cmd := h.db().PFCount(ctx, h.key)
	return &IntCmd{cmd: cmd}
}

// 合并其他的
func (h *hyperloglog) PFMerge(hyperloglogs ...hyperloglog) *BoolCmd {
	keys := make([]string, len(hyperloglogs))
	for i, hl := range hyperloglogs {
		keys[i] = hl.key
	}
	cmd := h.db().PFMerge(ctx, h.key, keys...)
	return &BoolCmd{cmd: cmd}
}
