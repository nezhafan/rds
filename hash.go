package redis

type hash struct {
	base
}

func NewHash(key string) hash {
	return hash{base{key}}
}

func (h hash) SubKey(subkey string) hash {
	return hash{base{h.key + ":" + subkey}}
}

func (h hash) HSet(field string, val any) error {
	return rdb.HSet(ctx, h.key, field, val).Err()
}

// 存储结构体时需要使用redis标签。 Id int `json:"id" redis:"id"`
func (h hash) HMSet(obj any) error {
	return rdb.HSet(ctx, h.key, obj).Err()
}

func (h hash) HGet(field string) (string, error) {
	return rdb.HGet(ctx, h.key, field).Result()
}

func (h hash) HMGet(to any, fields ...string) error {
	if len(fields) == 0 {
		return nil
	}
	return rdb.HMGet(ctx, h.key, fields...).Scan(to)
}

func (h hash) HGetAll(to any) error {
	return rdb.HGetAll(ctx, h.key).Scan(to)
}

func (h hash) HIncrBy(field string, incr int64) (int64, error) {
	return rdb.HIncrBy(ctx, h.key, field, incr).Result()
}

func (h hash) HIncrByFloat(field string, incr float64) (float64, error) {
	return rdb.HIncrByFloat(ctx, h.key, field, incr).Result()
}

func (h hash) HDel(fields ...string) error {
	return rdb.HDel(ctx, h.key, fields...).Err()
}

func (h hash) HExists(field string) bool {
	return rdb.HExists(ctx, h.key, field).Val()
}
