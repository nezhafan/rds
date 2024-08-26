package rdb

import "strconv"

type bitmap struct {
	base
}

// https://redis.io/docs/latest/commands/setbit/
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

// 获取范围内1的个数
func (b bitmap) BitCount(start, end int) int {
	args := []any{"BITCOUNT", b.key, start, end, "BIT"}
	i, _ := rdb.Do(ctx, args...).Int()
	return i
}

// 返回第一个0或1的位置
func (b bitmap) BitPos(search bool, start, end int) int {
	var n int
	if search {
		n = 1
	}
	args := []any{"BITPOS", b.key, n, start, end, "BIT"}
	i, _ := rdb.Do(ctx, args...).Int()
	return i
}

// 合并别的bitmap
// op: AND OR XOR NOT
// 如果是临时统计，请给key加上过期时间
func (b bitmap) BitOP(op string, srcKeys ...any) {
	commands := make([]any, 0, len(srcKeys)+3)
	commands = append(commands, "BITOP", op, b.key)
	commands = append(commands, srcKeys...)
	rdb.Do(ctx, commands...)
}

type bitfield struct {
	base
}

// https://redis.io/docs/latest/commands/bitfield/
func NewBitField(key string) bitfield {
	return bitfield{
		base: base{key},
	}
}

func (b bitfield) Set(typ string, offset uint32, value uint32) (uint32, error) {
	slice, err := rdb.Do(ctx, "BITFIELD", b.key, "OVERFLOW", "SAT", "SET", typ, offset, value).Slice()
	if err != nil {
		return 0, err
	}
	return uint32(slice[0].(int64)), nil
}

func (b bitfield) IncrBy(typ string, offset uint32, value uint32) (uint32, error) {
	slice, err := rdb.Do(ctx, "BITFIELD", b.key, "OVERFLOW", "SAT", "INCRBY", typ, offset, value).Slice()
	if err != nil {
		return 0, err
	}
	return uint32(slice[0].(int64)), nil
}

func (b bitfield) Get(typ string, offset uint32) (uint32, error) {
	slice, err := rdb.Do(ctx, "BITFIELD_RO", b.key, "GET", typ, offset).Slice()
	if err != nil {
		return 0, err
	}
	return uint32(slice[0].(int64)), nil
}

type autobitfield struct {
	base
	bits []uint8
}

// bit位可以是1-32位
// 如果存IP、时间戳使用32位，如果存用户ID使用24位就能存1600万
func NewAutoBitField(key string, bits ...uint8) autobitfield {
	if len(bits) == 0 {
		panic("至少需要一个参数")
	}
	for _, b := range bits {
		if b > 32 {
			panic("限制最大32位")
		}
		if b == 0 {
			panic("禁止为0")
		}
	}
	return autobitfield{
		base: base{key},
		bits: bits,
	}
}

// 返回原值。不会溢出。
func (b autobitfield) AutoSet(values ...uint32) ([]uint32, error) {
	if len(values) != len(b.bits) {
		panic("参数值数量必须与New时一一对应")
	}
	commands := make([]any, 0, len(b.bits)*6+2)
	commands = append(commands, "BITFIELD", b.key)
	var offset int
	for i, bit := range b.bits {
		commands = append(commands, "OVERFLOW", "SAT", "SET", "u"+strconv.Itoa(int(bit)), offset, values[i])
		offset += int(bit) + 1
	}

	return b.autodo(commands)
}

// 返回增长后的值。不会溢出。
func (b autobitfield) AutoIncrBy(values ...uint32) ([]uint32, error) {
	if len(values) != len(b.bits) {
		panic("参数值数量必须与New时一一对应")
	}
	commands := make([]any, 0, len(values)*6+2)
	commands = append(commands, "BITFIELD", b.key)
	var offset int
	for i, bit := range b.bits {
		commands = append(commands, "OVERFLOW", "SAT", "INCRBY", "u"+strconv.Itoa(int(bit)), offset, values[i])
		offset += int(bit) + 1
	}

	return b.autodo(commands)
}

func (b autobitfield) AutoGet() ([]uint32, error) {
	commands := make([]any, 0, len(b.bits)*3+2)
	commands = append(commands, "BITFIELD_RO", b.key)
	var offset int
	for _, bit := range b.bits {
		commands = append(commands, "GET", "u"+strconv.Itoa(int(bit)), offset)
		offset += int(bit) + 1
	}

	return b.autodo(commands)
}

func (b autobitfield) autodo(commands []any) ([]uint32, error) {
	slice, err := rdb.Do(ctx, commands...).Int64Slice()
	if err != nil {
		return nil, err
	}
	r := make([]uint32, 0, len(slice))
	for _, n := range slice {
		r = append(r, uint32(n))
	}
	return r, nil
}
