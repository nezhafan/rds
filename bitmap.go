package rds

import "context"

type Bitmap struct {
	base
}

// 位操作，适用于表达二元情况
// 不支持设置负数，最大uint32
// 占用内存取决于设置的最大值，1-1千万占用1.2MB，math.MaxUint32最大占用512M内存
// https://redis.io/docs/latest/commands/setbit/
func NewBitmap(ctx context.Context, key string) Bitmap {
	return Bitmap{base: NewBase(ctx, key)}
}

// 设置位，返回设置之前该位的布尔值
func (b *Bitmap) SetBit(offset uint32, ok bool) *BoolCmd {
	var v int
	if ok {
		v = 1
	}
	cmd := b.db().SetBit(b.ctx, b.key, int64(offset), v)
	b.done(cmd)
	return &BoolCmd{cmd: cmd}
}

// 获取位状态
func (b *Bitmap) GetBit(offset uint32) *BoolCmd {
	cmd := b.db().GetBit(b.ctx, b.key, int64(offset))
	b.done(cmd)
	return &BoolCmd{cmd: cmd}
}

// 获取范围内1的个数
func (b *Bitmap) BitCount(start, end int64) *IntCmd {
	args := []any{"bitcount", b.key, start, end, "bit"}
	cmd := b.db().Do(b.ctx, args...)
	b.done(cmd)
	return &IntCmd{cmd: cmd}
}

// // 返回第一个0或1的位置
// func (b *Bitmap) BitPos(search int, start, end int64) int64 {
// 	if search > 1 {
// 		search = 1
// 	}
// 	b.ctx, cancel := b.context()
// 	args := []any{"BITPOS", b.key, search, start, end, "BIT"}
// 	i, _ := DB().Do(b.ctx, args...).Int64()
// 	return i
// }

// 合并别的bitmap
// op: AND OR XOR NOT
// 如果是临时统计，请给key加上过期时间
// func (b *Bitmap) BitOp(op string, srcKeys ...any) error {
// 	commands := make([]any, 0, len(srcKeys)+3)
// 	commands = append(commands, "BITOP", op, b.key)
// 	commands = append(commands, srcKeys...)
// 	b.ctx, cancel := b.context()
// 	cmd := DB().Do(b.ctx, commands...)
// 	b.done(cmd, cancel)
// 	return cmd.Err()
// }

// type bitfield struct {
// 	base
// }

// bitfield是对bitmap的分段切割，如果用不好，使用NewAutoBitField
// 用一个bitmap表示多个作用，bitmap如果是对多用户区分二元，bitfield更像对单用户记录多个数字类型字段。
// https://redis.io/docs/latest/commands/bitfield/
// func NewBitField(key string) bitfield {
// 	return bitfield{
// 		base: NewBase(key),
// 	}
// }

// func (b *bitfield) Set(typ string, offset uint32, value uint32) (uint32, error) {
// 	b.ctx, cancel := b.context()
// 	slice, err := DB().Do(b.ctx, "BITFIELD", b.key, "OVERFLOW", "SAT", "SET", typ, offset, value).Slice()
// 	if err != nil {
// 		return 0, err
// 	}
// 	return uint32(slice[0].(int64)), nil
// }

// func (b *bitfield) IncrBy(typ string, offset uint32, value uint32) (uint32, error) {
// 	b.ctx, cancel := b.context()
// 	slice, err := DB().Do(b.ctx, "BITFIELD", b.key, "OVERFLOW", "SAT", "INCRBY", typ, offset, value).Slice()
// 	if err != nil {
// 		return 0, err
// 	}
// 	return uint32(slice[0].(int64)), nil
// }

// func (b *bitfield) Get(typ string, offset uint32) (uint32, error) {
// 	b.ctx, cancel := b.context()
// 	slice, err := DB().Do(b.ctx, "BITFIELD_RO", b.key, "GET", typ, offset).Slice()
// 	if err != nil {
// 		return 0, err
// 	}
// 	return uint32(slice[0].(int64)), nil
// }

// type autobitfield[E Unsigned] struct {
// 	base
// 	bits []uint8
// }

// 对bitfield的自动切割，也是bit位操作
// 例如使用 32，32 记录 登录IP、登录时间戳。
// bit位的大小不必为8的倍数（但是实际内存会对齐，剩余部分可以预留）
// 在考虑数字最大值的情况下节约，如果设置的值超过范围，会保持在最大值，不会溢出。
// 自动处理都是无符号类型，如果需要存负数，要么使用bitfield，要么用1位表示正负，代码再判断拼接。
// func NewAutoBitField[E Unsigned](key string, bits ...uint8) autobitfield[E] {
// 	if len(bits) == 0 {
// 		panic("至少需要一个参数")
// 	}
// 	for _, b := range bits {
// 		if b > 32 {
// 			panic("限制最大32位")
// 		}
// 		if b == 0 {
// 			panic("禁止为0")
// 		}
// 	}
// 	return autobitfield[E]{
// 		base: NewBase(key),
// 		bits: bits,
// 	}
// }

// // 返回原值。不会溢出。
// func (b *autobitfield[E]) AutoSet(values ...uint32) ([]uint32, error) {
// 	if len(values) != len(b.bits) {
// 		panic("参数值数量必须与New时一一对应")
// 	}
// 	commands := make([]any, 0, len(b.bits)*6+2)
// 	commands = append(commands, "BITFIELD", b.key)
// 	var offset int
// 	for i, bit := range b.bits {
// 		commands = append(commands, "OVERFLOW", "SAT", "SET", "u"+strconv.Itoa(int(bit)), offset, values[i])
// 		offset += int(bit) + 1
// 	}

// 	return b.autodo(commands)
// }

// // 返回增长后的值。不会溢出。
// func (b *autobitfield[E]) AutoIncrBy(values ...uint32) ([]uint32, error) {
// 	if len(values) != len(b.bits) {
// 		panic("参数值数量必须与New时一一对应")
// 	}
// 	commands := make([]any, 0, len(values)*6+2)
// 	commands = append(commands, "BITFIELD", b.key)
// 	var offset int
// 	for i, bit := range b.bits {
// 		commands = append(commands, "OVERFLOW", "SAT", "INCRBY", "u"+strconv.Itoa(int(bit)), offset, values[i])
// 		offset += int(bit) + 1
// 	}

// 	return b.autodo(commands)
// }

// func (b *autobitfield[E]) AutoGet() ([]uint32, error) {
// 	commands := make([]any, 0, len(b.bits)*3+2)
// 	commands = append(commands, "BITFIELD_RO", b.key)
// 	var offset int
// 	for _, bit := range b.bits {
// 		commands = append(commands, "GET", "u"+strconv.Itoa(int(bit)), offset)
// 		offset += int(bit) + 1
// 	}

// 	return b.autodo(commands)
// }

// func (b *autobitfield[E]) autodo(commands []any) ([]uint32, error) {
// 	slice, err := DB().Do(b.ctx, commands...).Int64Slice()
// 	if err != nil {
// 		return nil, err
// 	}
// 	r := make([]uint32, 0, len(slice))
// 	for _, n := range slice {
// 		r = append(r, uint32(n))
// 	}
// 	return r, nil
// }

func (b *Bitmap) WithCmdable(cmdable Cmdable) *Bitmap {
	bc := b.base
	bc.cmdable = cmdable
	return &Bitmap{base: bc}
}
