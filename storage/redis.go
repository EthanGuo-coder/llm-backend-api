package storage

import (
	"encoding/json"
	"fmt"

	"github.com/EthanGuo-coder/llm-backend-api/models"
)

// SaveConversation 保存完整会话到 Redis
func SaveConversation(conversation *models.Conversation) error {
	conversationKey := fmt.Sprintf("conversation:%s", conversation.ID)
	data, err := json.Marshal(conversation)
	if err != nil {
		return fmt.Errorf("failed to marshal conversation: %v", err)
	}
	return redisClient.Set(ctx, conversationKey, data, 0).Err()
}

// GetConversation 从 Redis 获取完整会话
func GetConversation(conversationID string) (*models.Conversation, error) {
	conversationKey := fmt.Sprintf("conversation:%s", conversationID)
	data, err := redisClient.Get(ctx, conversationKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %v", err)
	}
	var conversation models.Conversation
	if err := json.Unmarshal([]byte(data), &conversation); err != nil {
		return nil, fmt.Errorf("failed to unmarshal conversation: %v", err)
	}
	return &conversation, nil
}
