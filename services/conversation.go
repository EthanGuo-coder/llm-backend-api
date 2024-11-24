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
func CreateConversation(userID int64, title, model, apiKey string) (*models.Conversation, error) {
	// 生成唯一会话 ID
	conversationID := utils.GenerateID()

	// 构造会话对象
	conversation := &models.Conversation{
		ID:     conversationID,
		Title:  title,
		Model:  model,
		ApiKey: apiKey, // 存储 api_key
		Messages: []models.Message{
			{Role: "system", Content: constant.SystemPrompt},
		},
		CreatedTime: time.Now().Unix(),
	}

	// 保存会话元信息到数据库
	err := storage.SaveConversationToDB(userID, conversation)
	if err != nil {
		return nil, errors.New("failed to save conversation to database: " + err.Error())
	}

	// 初始化会话记录到 Redis
	err = storage.SaveConversationToRedis(conversation)
	if err != nil {
		return nil, errors.New("failed to save conversation to redis: " + err.Error())
	}

	return conversation, nil
}

// GetConversationHistory 获取完整的会话历史
func GetConversationHistory(conversationID string) (*models.ConversationHistory, error) {
	// 从 Redis 获取完整会话的记录
	conversation, err := storage.GetConversationFromRedis(conversationID)
	if err != nil {
		return nil, errors.New("failed to fetch conversation from redis: " + err.Error())
	}
	if conversation == nil {
		return nil, errors.New("conversation not found")
	}

	// 过滤掉 role 为 "system" 的消息
	filteredMessages := make([]models.Message, 0)
	for _, message := range conversation.Messages {
		if message.Role != "system" {
			filteredMessages = append(filteredMessages, message)
		}
	}

	// 转换为 ConversationHistory
	history := &models.ConversationHistory{
		ID:       conversation.ID,
		Title:    conversation.Title,
		Model:    conversation.Model,
		Messages: filteredMessages, // 使用过滤后的消息
	}

	return history, nil
}

// DeleteUserConversation 删除指定的用户对话
func DeleteUserConversation(userID int64, conversationID string) error {
	// 从数据库删除会话元信息
	err := storage.DeleteConversationFromDB(userID, conversationID)
	if err != nil {
		return errors.New("failed to delete conversation from database: " + err.Error())
	}

	// 从 Redis 删除会话记录
	err = storage.DeleteConversationFromRedis(conversationID)
	if err != nil {
		return errors.New("failed to delete conversation from redis: " + err.Error())
	}

	return nil
}

// GetUserConversations 查询用户的所有对话
func GetUserConversations(userID int64) ([]models.ConversationSummary, error) {
	// 获取用户的所有会话
	conversations, err := storage.FetchConversationsByUserID(userID)
	if err != nil {
		return nil, errors.New("failed to fetch conversations: " + err.Error())
	}

	// 转换为 ConversationSummary
	var summaries []models.ConversationSummary
	for _, conversation := range conversations {
		summaries = append(summaries, models.ConversationSummary{
			ID:          conversation.ID,
			Title:       conversation.Title,
			CreatedTime: conversation.CreatedTime,
		})
	}

	return summaries, nil
}
