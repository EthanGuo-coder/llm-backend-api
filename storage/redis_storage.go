package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/EthanGuo-coder/llm-backend-api/config"
	"github.com/EthanGuo-coder/llm-backend-api/models"
)

var redisClient *redis.Client
var ctx = context.Background()

// InitializeRedis 初始化 Redis 客户端
func InitializeRedis() error {
	cfg := config.AppConfig.Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	// 测试连接
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}
	fmt.Println("Connected to Redis successfully!")
	return nil
}

// SaveConversationToRedis 保存完整会话到 Redis
func SaveConversationToRedis(conversation *models.Conversation) error {
	conversationKey := fmt.Sprintf("conversation:%s", conversation.ID)
	data, err := json.Marshal(conversation)
	if err != nil {
		return fmt.Errorf("failed to marshal conversation: %v", err)
	}
	return redisClient.Set(ctx, conversationKey, data, 0).Err()
}

// GetConversationFromRedis 从 Redis 获取完整会话
func GetConversationFromRedis(conversationID string) (*models.Conversation, error) {
	conversationKey := fmt.Sprintf("conversation:%s", conversationID)
	data, err := redisClient.Get(ctx, conversationKey).Result()
	if err == redis.Nil {
		return nil, nil // 会话不存在
	} else if err != nil {
		return nil, fmt.Errorf("failed to get conversation from redis: %v", err)
	}

	var conversation models.Conversation
	if err := json.Unmarshal([]byte(data), &conversation); err != nil {
		return nil, fmt.Errorf("failed to unmarshal conversation: %v", err)
	}
	return &conversation, nil
}

// DeleteConversationFromRedis 从 Redis 中删除完整会话
func DeleteConversationFromRedis(conversationID string) error {
	conversationKey := fmt.Sprintf("conversation:%s", conversationID)
	return redisClient.Del(ctx, conversationKey).Err()
}

// CacheJWT 将 Token 存入 Redis
func CacheJWT(tokenStr string, userID int64, ttl time.Duration) error {
	key := fmt.Sprintf("jwt:%s", tokenStr)
	value := map[string]interface{}{
		"user_id": userID,
	}
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal token data: %w", err)
	}

	return redisClient.Set(ctx, key, data, ttl).Err()
}

// GetCachedJWT 从 Redis 中获取缓存的 Token
func GetCachedJWT(tokenStr string) (map[string]interface{}, error) {
	key := fmt.Sprintf("jwt:%s", tokenStr)
	data, err := redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // 未命中缓存
	} else if err != nil {
		return nil, fmt.Errorf("failed to get token from redis: %w", err)
	}

	var value map[string]interface{}
	if err := json.Unmarshal([]byte(data), &value); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token data: %w", err)
	}

	return value, nil
}
