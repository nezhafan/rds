package redis

type bitmap struct {
	base
}

// 1. 每个用户的某个状态位，offset直接使用id
// 2. 单个用户的多个状态位，例如每日登录，划分区域366*10年，这仍然能存储100万人。 offset = id*366+每年的第几天

func NewBitmap(key string) bitmap {
	return bitmap{base{key}}
}

// 设置
func (b bitmap) SetBit(offset uint32, ok bool) error {
	var v int
	if ok {
		v = 1
	}
	return rdb.SetBit(ctx, b.key, int64(offset), v).Err()
}

// 获取
func (b bitmap) GetBit(offset uint32) bool {
	v := rdb.GetBit(ctx, b.key, int64(offset)).Val()
	return v == 1
}

// 获取范围内1的个数，带BIT参数为>7.0版本支持。
// 不带BIT的BITCOUNT为按照bytes统计，如果一个bytes被分割为两个区域，难以统计
func (b bitmap) BitCount(start, end uint32) int64 {
	args := []any{"BITCOUNT", b.key, start, end, "BIT"}
	v := rdb.Do(ctx, args...).Val()
	if n, ok := v.(int64); ok {
		return n
	}
	return 0
}
