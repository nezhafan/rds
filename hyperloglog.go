package rds

type hyperloglog struct {
	base
}

// 基数统计，存在0.81%误差。
// 最大 12k 左右内存，可以统计 2^64 个元素
// 绝大多数公司都用不上，可以使用bitmap
func NewHyperLogLog(key string) hyperloglog {
	return hyperloglog{base: newBase(key)}
}

// 添加，至少有一个添加成功返回1，否则返回0
func (h *hyperloglog) PFAdd(vals ...any) (bool, error) {
	n, err := rdb.PFAdd(ctx, h.key, vals...).Result()
	return n == 1, err
}

// 统计数量，是存在0.81%误差的
func (h *hyperloglog) PFCount() int64 {
	return rdb.PFCount(ctx, h.key).Val()
}

// 合并其他的
func (h *hyperloglog) PFMerge(hyperloglogs ...hyperloglog) (bool, error) {
	keys := make([]string, len(hyperloglogs))
	for i, hl := range hyperloglogs {
		keys[i] = hl.key
	}
	r, err := rdb.PFMerge(ctx, h.key, keys...).Result()
	return r == OK, err
}
