package rds

type list struct {
	base
}

func NewList(key string) list {
	return list{base: newBase(key)}
}

// 左入
func (l *list) LPush(vals ...any) (int64, error) {
	return rdb.LPush(ctx, l.key, vals...).Result()
}

// 右入
func (l *list) RPush(vals ...any) (int64, error) {
	return rdb.RPush(ctx, l.key, vals...).Result()
}

// 左出
func (l *list) LPop() (string, error) {
	return rdb.LPop(ctx, l.key).Result()
}

// 右出
func (l *list) RPop() (string, error) {
	return rdb.RPop(ctx, l.key).Result()
}

// 遍历
func (l *list) LRange(start, stop int) []string {
	return rdb.LRange(ctx, l.key, int64(start), int64(stop)).Val()
}

// 长度
func (l *list) LLen() int64 {
	return rdb.LLen(ctx, l.key).Val()
}

// 从左开始移除count个
func (l *list) LRem(value any, count int64) (int64, error) {
	return rdb.LRem(ctx, l.key, count, value).Result()
}

// 从右开始移除count个
func (l *list) RRem(value any, count int64) (int64, error) {
	return rdb.LRem(ctx, l.key, -count, value).Result()
}

// 移除，返回被移除数量
func (l *list) Rem(value any) (int64, error) {
	return rdb.LRem(ctx, l.key, 0, value).Result()
}
