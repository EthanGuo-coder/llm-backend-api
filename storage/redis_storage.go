package storage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/EthanGuo-coder/llm-backend-api/models"
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

// AppendMessageToRedis 追加消息到 Redis 中
func AppendMessageToRedis(conversationID string, message models.Message) error {
	key := fmt.Sprintf("conversation:%s:messages", conversationID)
	messageData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}
	return redisClient.RPush(ctx, key, messageData).Err()
}

// GetMessagesFromRedis 从 Redis 获取会话消息列表
func GetMessagesFromRedis(conversationID string) ([]models.Message, error) {
	key := fmt.Sprintf("conversation:%s:messages", conversationID)
	data, err := redisClient.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch messages: %v", err)
	}

	messages := make([]models.Message, len(data))
	for i, item := range data {
		if err := json.Unmarshal([]byte(item), &messages[i]); err != nil {
			return nil, fmt.Errorf("failed to unmarshal message: %v", err)
		}
	}
	return messages, nil
}
