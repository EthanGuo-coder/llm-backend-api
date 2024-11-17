package services

import (
	"fmt"
	"sync"

	"github.com/EthanGuo-coder/llm-backend-api/models"
	"github.com/EthanGuo-coder/llm-backend-api/storage"
)

func SendMessage(conversationID, model, apiKey, message string) (map[string]string, error) {
	conversation, err := storage.GetConversation(conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch conversation: %w", err)
	}

	// 用户消息
	userMessage := models.Message{Role: "user", Content: message}
	conversation.Messages = append(conversation.Messages, userMessage)

	// 模拟 AI 回复
	aiResponseContent := fmt.Sprintf("你好！我看到你说: %s", message)
	aiMessage := models.Message{Role: "assistant", Content: aiResponseContent}
	conversation.Messages = append(conversation.Messages, aiMessage)

	// 异步更新 Redis
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := storage.SaveConversation(conversation); err != nil {
			fmt.Printf("Failed to save conversation: %v\n", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := storage.SaveMessages(conversationID, []models.Message{userMessage, aiMessage}); err != nil {
			fmt.Printf("Failed to save messages: %v\n", err)
		}
	}()

	wg.Wait()

	return map[string]string{
		"user_message": message,
		"ai_response":  aiResponseContent,
	}, nil
}
