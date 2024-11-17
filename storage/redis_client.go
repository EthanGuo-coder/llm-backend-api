package storage

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var ctx = context.Background()

// InitializeRedis 初始化 Redis 客户端
func InitializeRedis(addr, password string, db int) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// 测试连接
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}
	fmt.Println("Connected to Redis successfully!")
}

// GetRedisClient 获取 Redis 客户端
func GetRedisClient() *redis.Client {
	return redisClient
}
