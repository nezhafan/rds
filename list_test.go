package rds

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	Connect("127.0.0.1:6379", "123456", 0)

	m.Run()
}

func TestList(t *testing.T) {
	l := NewList[string](nil, "list")
	l.LPush("a", "b", "c")
	fmt.Println(l.LRange(0, 2))
	assert.Equal(t, []string{"c", "b", "a"}, l.LRange(0, 2))

}
