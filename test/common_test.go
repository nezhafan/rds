package test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/nezhafan/rds"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type User struct {
	Name    string   `redis:"name" json:"name"`
	Age     age      `redis:"age" `
	Pet     *string  `redis:"pet,omitempty"`
	Likes   []string `redis:"likes,omitempty"`
	Guns    []gun    `redis:"guns,omitempty"`
	Money   float64  `redis:"money"`
	Nothing string   // 不存储
}

type age uint8

type gun struct {
	Name  string `redis:"name"`
	Price int    `redis:"price"`
}

func TestMain(m *testing.M) {
	err := rds.Connect("127.0.0.1:6379", "", 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rds.DB().(*redis.Client).Close()

	// 打印参数和返回值
	rds.SetDebug(true)

	// 打印错误
	rds.SetErrorHook(func(err error) {
		fmt.Println("报错啦:", err)
	})

	os.Exit(m.Run())
}
