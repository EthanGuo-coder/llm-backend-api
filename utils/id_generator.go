package utils

import (
	"math/rand"
	"time"
)

func GenerateID() int64 {
	rand.Seed(time.Now().UnixNano())
	return int64(rand.Intn(1000000)) // 返回随机生成的 int64 类型
}
