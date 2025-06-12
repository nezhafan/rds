package rds

import (
	"fmt"
	"testing"
)

type User struct {
	Id   int    `redis:"id"`
	Name string `redis:"name"`
}

func init() {
	if err := Connect("127.0.0.1:6379", "123456", 0); err != nil {
		panic(err)
	}
}

func TestString(t *testing.T) {
	// str := NewString(nil, "test-string")
	// err := str.Set("123", time.Hour)
	// fmt.Println("set 123:", err)
	// b, err := str.SetNX("234", time.Hour)
	// fmt.Println("setnx 234:", b, err)
	// v, has := str.Get()
	// fmt.Println("get:", v, has)
	// ttl := str.TTL()
	// fmt.Println("ttl:", ttl)
	// ok := str.Expire(time.Second)
	// fmt.Println("expire:", ok)
	// ttl = str.TTL()
	// fmt.Println("ttl:", ttl)
	// ok = str.Del()
	// fmt.Println("del:", ok)
	// v, has = str.Get()
	// fmt.Println("get:", `"`+v+`"`, has)
	// ii := NewInt(nil, "test-int")
	// ii.IncrBy(1)

	ff := NewHashMap[string]("test-float")

	fmt.Println(ff.HMGet().Val())

	// aa := newBase("tes")
	// c := aa.db().Set(ctx, "tes", 1, time.Second)
	// fmt.Println(c.Result())
}

// func BenchmarkString(b *testing.B) {
// 	str := NewString(nil, "test-string")
// 	defer str.Del()

// 	wg := &sync.WaitGroup{}
// 	wg.Add(b.N)
// 	fmt.Println("测试次数：", b.N)
// 	for i := 0; i < b.N; i++ {
// 		go func() {
// 			defer wg.Done()
// 			ok, err := str.SetNX("1", time.Minute)
// 			if err != nil {
// 				panic(err)
// 			}
// 			if ok {
// 				fmt.Println("设置成功")
// 			}
// 		}()
// 	}
// 	wg.Wait()
// }

// func TestBitmap(t *testing.T) {
// 	bf1 := NewBitmap(nil, "test-bitmap1")
// 	defer bf1.Del()
// 	bf1.SetBit(10000, true)
// 	bf1.SetBit(2, true)
// 	fmt.Println(bf1.GetBit(10000) == true)
// 	fmt.Println(bf1.BitCount(0, -1) == 2)
// 	fmt.Println(bf1.BitPos(1, 0, -1) == 2)
// 	fmt.Println(bf1.BitPos(0, 0, -1) == 0)

// 	bf2 := NewBitmap(nil, "test-bitmap2")
// 	defer bf2.Del()
// 	bf2.SetBit(1, true)
// 	bfmerge := NewBitmap(nil, "test-bitmap-merge")
// 	defer bfmerge.Del()
// 	bfmerge.BitOp("OR", bf1.key, bf2.key)
// 	fmt.Println(bfmerge.BitCount(0, -1) == 3)
// 	bfmerge.BitOp("AND", bf1.key, bf2.key)
// 	fmt.Println(bfmerge.BitCount(0, -1) == 0)

// 	bitfield := NewBitField(nil, "test-bitfield")

// 	defer bitfield.Del()
// 	bitfield.Set("i8", 0, 17)
// 	bitfield.Set("i8", 9, 32)
// 	bitfield.IncrBy("i8", 9, 1)

// 	n, err := bitfield.Get("i8", 0)
// 	fmt.Println(n, err)
// 	n, err = bitfield.Get("i8", 9)
// 	fmt.Println(n, err)

// 	bitautofield := NewAutoBitField[uint32](nil, "test-bitautofield", 32, 12, 2)
// 	// defer bitautofield.Del()

// 	a, err := bitautofield.AutoSet(1700000000, 4004, 1)
// 	fmt.Println(a, err)
// 	a, err = bitautofield.AutoGet()
// 	fmt.Println(a, err)
// 	a, err = bitautofield.AutoIncrBy(0, 1, 0)
// 	fmt.Println(a, err)

// }

// func TestHash(t *testing.T) {
// 	hash := NewHash(nil, "test-hash")

// 	hash1 := hash.SubKey("1")
// 	defer hash1.Del()

// 	_, err := hash1.HSet("name", "Alice")
// 	fmt.Println("hset name Alice:", err)
// 	v, has := hash1.HGet("name")
// 	fmt.Println("hget name:", v, has)

// 	hash2 := hash.SubKey("2")
// 	defer hash2.Del()

// 	var u1, u2 User
// 	u1.Id = 2
// 	u1.Name = "Bob"
// 	err = hash2.HMSet(&u1, time.Minute)
// 	fmt.Println("hmset id:", 2, "name:", "Bob")

// 	exists := hash2.HMGet(&u2, "name")
// 	fmt.Println("hmget name:", u2, exists)
// 	exists = hash2.HGetAll(&u2)
// 	fmt.Println("hgetall:", u2, exists)

// 	i, err := hash2.HIncrBy("id", 10)
// 	fmt.Println("hincrby id 10:", i, err)

// }

// func TestHyperLogLog(t *testing.T) {
// 	hyperloglog1 := NewHyperLogLog(nil, "test-hyperloglog1")
// 	defer hyperloglog1.Del()
// 	hyperloglog2 := NewHyperLogLog(nil, "test-hyperloglog2")
// 	defer hyperloglog2.Del()

// 	n := 1000
// 	es := make([]any, 0, n)
// 	for i := 0; i < n; i++ {
// 		es = append(es, i)
// 	}

// 	b, err := hyperloglog1.PFAdd(es...)
// 	fmt.Println("pfadd:", b, err)

// 	size := hyperloglog1.PFCount()
// 	fmt.Println("pfcount:", size)

// 	hyperloglog2.PFAdd(1, 2, 3, 10002, 10003)
// 	b, err = hyperloglog1.PFMerge(hyperloglog2)
// 	fmt.Println("pfmerge:", b, err)

// 	size = hyperloglog1.PFCount()
// 	fmt.Println("pfcount:", size)
// }

// func TestList(t *testing.T) {
// 	list := NewList[int](nil, "test-list")
// 	defer list.Del()

// 	_, err := list.LPush(2, 1)
// 	fmt.Println("lpush 2 1:", err)
// 	_, err = list.RPush(3, 3, 4)
// 	fmt.Println("rpush 3 3 4:", err)

// 	n := list.LLen()
// 	fmt.Println("llen:", n)

// 	r := list.LRange(0, -1)
// 	fmt.Println("lrange:", r)

// 	v, err := list.LPop()
// 	fmt.Println("lpop:", v, err)

// 	v, err = list.RPop()
// 	fmt.Println("rpop:", v, err)

// 	n, err = list.LRem("2", 1)
// 	fmt.Println("lrem 4:", n, err)

// 	n, err = list.Rem("3")
// 	fmt.Println("rem 3:", n, err)

// }

// 	n, err := set.SAdd(1, 2)
// 	fmt.Println("添加元素1|2，添加成功数", n, err)
// 	// assert.Equal(t, 2, n)
// 	n, err = set.SAdd(1, 3)
// 	fmt.Println("添加元素1|3，添加成功数", n, err)
// 	// assert.Equal(t, 1, n)

// 	n = set.SCard()
// 	fmt.Println("总元素数量", n)

// 	all := set.SMembers()
// 	fmt.Println("获取所有元素", all)

// 	n, err = set.SRem(2, 3)
// 	fmt.Println("移除元素2|3，移除成功数量", n, err)

// 	has := set.SIsMember(2)
// 	fmt.Println("元素2是否存在", has)
// }

// func TestZSet(t *testing.T) {
// 	zset := NewZSet[string, float32](nil, "test-zset")
// 	zset.Expire(time.Minute)

// 	n, err := zset.ZAdd(
// 		Z2[string, float32]{Score: 0.5, Member: "A"},
// 		Z2[string, float32]{Score: 1, Member: "C"},
// 		Z2[string, float32]{Score: 0.5, Member: "B"},
// 		Z2[string, float32]{Score: 1, Member: "D"},
// 	)
// 	fmt.Println("添加ABCD:", n, err)

// 	members := zset.ZRange(0, -1)
// 	fmt.Println("正序获取全部成员:", members)

// 	r := zset.ZRangeWithScores(0, -1)
// 	fmt.Println("正序获取全部成员和积分:", r)

// 	n = zset.ZCard()
// 	fmt.Println("成员个数:", n)

// 	score, has := zset.ZScore("B")
// 	fmt.Println("B的积分:", score, has)

// 	n, err = zset.ZRem("B")
// 	fmt.Println("移除B:", n, err)

// 	n = zset.ZLexCount("(A", "[C")
// 	fmt.Println("获取(A,C]之间成员个数:", n)

// 	members = zset.ZRevRange(0, -1)
// 	fmt.Println("逆序获取全部成员:", members)

// 	r = zset.ZRevRangeWithScores(0, -1)
// 	fmt.Println("逆序获取全部成员和积分:", r)

// 	r = zset.ZRangeByScore("1", MaxInf, 1)
// 	fmt.Println("获取积分>=1的1个成员和积分:", r)

// 	n, err = zset.ZRemRangeByRank(-2, -2)
// 	fmt.Println("删除排名倒数第2的成员:", n, err)

// 	n = zset.ZRank("D")
// 	fmt.Println("获取D从小到大的排名:", n)

// 	n = zset.ZRevRank("D")
// 	fmt.Println("获取D从大到小的排名:", n)

// 	// 全部删除
// 	n, err = zset.ZRemRangeByScore(MinInf, MaxInf)
// 	fmt.Println("删除积分为-inf到+inf之间的成员:", n, err)

// }

// func TestGeo(t *testing.T) {
// 	geo := NewGeo(nil, "test-geo")
// 	// defer geo.Del()

// 	// 添加坐标，其中 BCDE距离相近，C和E相同
// 	geo.GeoAdd("A", 13.361389, 38.115556)
// 	geo.GeoAdd("B", 15.077268, 37.502669)
// 	geo.GeoAdd("C", 15.087268, 37.502669)
// 	geo.GeoAdd("D", 15.067268, 37.502669)
// 	geo.GeoAdd("E", 15.087268, 37.502669)

// 	dist := geo.GeoDist("A", "B")
// 	fmt.Println("A和B的距离", dist)

// 	dist = geo.GeoDist("B", "C")
// 	fmt.Println("B和C的距离", dist)

// 	dist = geo.GeoDist("D", "F")
// 	fmt.Println("D和F的距离", dist)

// 	pos := geo.GeoPos("C", "F")
// 	fmt.Println("C的经纬度", pos[0], "F的经纬度", pos[1])

// 	locations := geo.GeoSearchByCoord(pos[0].Longitude, pos[0].Latitude, 2000, 10, false, true, false)
// 	fmt.Println("C半径范围2000米内成员(按坐标)(不排序)", locations)

// 	locations = geo.GeoSearchByMember("C", 2000, 10, false, true, true)
// 	fmt.Println("C半径范围2000米内成员(按成员)(正序)", locations)

// }

// func TestStream(t *testing.T) {
// 	// count := 1000
// 	// // 生产者
// 	// stream := NewStream("test-stream")
// 	// defer stream.Del()
// 	// fmt.Println("stream中插入", count, "条消息")
// 	// wg := &sync.WaitGroup{}
// 	// wg.Add(count)
// 	// for i := 0; i < count; i++ {
// 	// 	go func() {
// 	// 		defer wg.Done()
// 	// 		_, err := stream.XAdd("", map[string]any{"id": ""})
// 	// 		if err != nil {
// 	// 			panic(err)
// 	// 		}
// 	// 	}()
// 	// }

// 	// // 创建组，并且绑定生产者
// 	// group1 := NewGroup("group1", stream)
// 	// fmt.Println("创建组group1读取stream")
// 	// group2 := NewGroup("group2", stream)
// 	// fmt.Println("创建组group2读取stream")

// 	// // 消费者
// 	// fmt.Println("创建2个消费者。其中A和B在group1竞争stream消息，但是A在group2独享stream消息")
// 	// var limit int64 = 10
// 	// a1 := NewConsumer("A", group1, limit)
// 	// b1 := NewConsumer("B", group1, limit)
// 	// a2 := NewConsumer("A", group2, limit)

// 	// cs := []consumer{a1, b1, a2}
// 	// counter := []int64{0, 0, 0}

// 	// for i, c := range cs {
// 	// 	go func(i int, c consumer) {
// 	// 		for {
// 	// 			rs, err := c.Read()
// 	// 			if err != nil {
// 	// 				break
// 	// 			}

// 	// 			atomic.AddInt64(&counter[i], int64(len(rs)))

// 	// 			for _, r := range rs {
// 	// 				go c.Ack(r.ID)
// 	// 			}
// 	// 			// fmt.Println("消费者", cs[i].args.Consumer, "从", cs[i].args.Group, "读取到：", rs)
// 	// 		}
// 	// 	}(i, c)
// 	// }

// 	// wg.Wait()
// 	// time.Sleep(time.Second * 2)

// 	// fmt.Println("分别收到消息数:", counter, "结果校验：", counter[0]+counter[1]+counter[2] == int64(count*2))
// }

// func BenchmarkMutex(b *testing.B) {

// 	wg := &sync.WaitGroup{}
// 	wg.Add(b.N)
// 	start := time.Now()
// 	for i := 0; i < b.N; i++ {
// 		go func(i int) {
// 			defer wg.Done()
// 			mutex := NewMutex("mutex")
// 			mutex.Lock()
// 			defer mutex.UnLock()
// 			time.Sleep(time.Millisecond * 2)
// 		}(i)
// 	}

// 	wg.Wait()
// 	fmt.Println("尝试次数", b.N, "耗时", time.Since(start).Milliseconds(), "ms")
// }
