package rdb

type set struct {
	base
}

func NewSet(key string) set {
	return set{base{key}}
}

// 添加成员，返回添加后成员数
func (s set) SAdd(members ...any) (int64, error) {
	return rdb.SAdd(ctx, s.key, members...).Result()
}

// 获取成员数
func (s set) SCard() int64 {
	return rdb.SCard(ctx, s.key).Val()
}

// 是否为成员
func (s set) SIsMember(member any) bool {
	return rdb.SIsMember(ctx, s.key, member).Val()
}

// 移除成员
func (s set) SRem(members ...any) (int64, error) {
	return rdb.SRem(ctx, s.key, members...).Result()
}

// 获取所有成员
func (s set) SMembers() []string {
	return rdb.SMembers(ctx, s.key).Val()
}
