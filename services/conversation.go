package services

import (
	"errors"
	"time"

	"github.com/EthanGuo-coder/llm-backend-api/constant"
	"github.com/EthanGuo-coder/llm-backend-api/models"
	"github.com/EthanGuo-coder/llm-backend-api/storage"
	"github.com/EthanGuo-coder/llm-backend-api/utils"
)

// CreateConversation 创建新的会话
func CreateConversation(title, model string) (*models.Conversation, error) {
	conversationID := utils.GenerateID()
	conversation := &models.Conversation{
		ID:    conversationID,
		Title: title,
		Model: model,
		Messages: []models.Message{
			{Role: "system", Content: constant.SystemPrompt},
		},
		CreatedTime: time.Now().Unix(),
	}
	if err := storage.SaveConversation(conversation); err != nil {
		return nil, errors.New("failed to save conversation: " + err.Error())
	}
	return conversation, nil
}

// GetConversationHistory 获取完整的会话历史
func GetConversationHistory(conversationID string) (*models.Conversation, error) {
	conversation, err := storage.GetConversation(conversationID)
	if err != nil {
		return nil, errors.New("failed to fetch conversation: " + err.Error())
	}
	if conversation == nil {
		return nil, errors.New("conversation not found")
	}
	return conversation, nil
}
