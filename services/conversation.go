package services

import (
	"errors"
	"time"

	"github.com/EthanGuo-coder/llm-backend-api/models"
	"github.com/EthanGuo-coder/llm-backend-api/storage"
	"github.com/EthanGuo-coder/llm-backend-api/utils"
)

// CreateConversation 创建新的会话
func CreateConversation(title string) (*models.Conversation, error) {
	// 生成唯一的会话 ID
	conversationID := utils.GenerateID()
	createdTime := time.Now().Format(time.RFC3339) // 使用 ISO 8601 格式记录时间

	// 初始化会话结构
	conversation := &models.Conversation{
		ID:          conversationID,
		Title:       title,
		Messages:    []models.Message{},
		CreatedTime: createdTime,
	}

	// 保存到存储层
	err := storage.SaveConversation(conversation)
	if err != nil {
		return nil, errors.New("failed to save conversation: " + err.Error())
	}

	return conversation, nil
}

// GetConversationHistory 获取会话的完整历史
func GetConversationHistory(conversationID string) (*models.Conversation, error) {
	// 获取会话元信息
	conversation, err := storage.GetConversation(conversationID)
	if err != nil {
		return nil, errors.New("failed to fetch conversation: " + err.Error())
	}
	if conversation == nil {
		return nil, errors.New("conversation not found")
	}

	// 获取该会话的消息记录
	messages, err := storage.GetMessages(conversationID)
	if err != nil {
		return nil, errors.New("failed to fetch messages: " + err.Error())
	}
	conversation.Messages = messages

	return conversation, nil
}
