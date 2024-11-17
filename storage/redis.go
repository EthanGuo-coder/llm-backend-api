package storage

import (
	"encoding/json"
	"fmt"

	"github.com/EthanGuo-coder/llm-backend-api/models"
)

// SaveConversation 保存会话元信息
func SaveConversation(conversation *models.Conversation) error {
	conversationKey := fmt.Sprintf("conversation:%s:meta", conversation.ID)
	data, err := json.Marshal(conversation)
	if err != nil {
		return fmt.Errorf("failed to marshal conversation: %v", err)
	}
	return redisClient.Set(ctx, conversationKey, data, 0).Err()
}

// GetConversation 获取会话元信息
func GetConversation(conversationID string) (*models.Conversation, error) {
	conversationKey := fmt.Sprintf("conversation:%s:meta", conversationID)
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

// SaveMessages 批量保存消息
func SaveMessages(conversationID string, messages []models.Message) error {
	messageKey := fmt.Sprintf("conversation:%s:messages", conversationID)
	pipe := redisClient.Pipeline()
	for _, message := range messages {
		data, err := json.Marshal(message)
		if err != nil {
			return fmt.Errorf("failed to marshal message: %v", err)
		}
		pipe.RPush(ctx, messageKey, data)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to save messages: %v", err)
	}
	return nil
}

// GetMessages 获取所有消息
func GetMessages(conversationID string) ([]models.Message, error) {
	messageKey := fmt.Sprintf("conversation:%s:messages", conversationID)
	data, err := redisClient.LRange(ctx, messageKey, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %v", err)
	}

	var messages []models.Message
	for _, item := range data {
		var msg models.Message
		if err := json.Unmarshal([]byte(item), &msg); err != nil {
			return nil, fmt.Errorf("failed to unmarshal message: %v", err)
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
