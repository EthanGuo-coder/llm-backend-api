package services

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EthanGuo-coder/llm-backend-api/models"
	"github.com/EthanGuo-coder/llm-backend-api/storage"
	"github.com/EthanGuo-coder/llm-backend-api/utils"
)

// StreamSendMessage 处理流式消息发送
func StreamSendMessage(c *gin.Context, conversationID int64, message string) error {
	// 获取会话
	conversation, err := getConversationWithMessage(conversationID, message)
	if err != nil {
		return err
	}
	// 构造请求体
	requestData, err := buildRequestBody(conversation)
	if err != nil {
		return err
	}
	// 使用存储的 api_key
	apiKey := conversation.ApiKey
	// 创建 HTTP 请求
	resp, err := sendAPIRequest(apiKey, requestData, conversation.Model)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// 检查响应状态码
	if err := validateResponse(resp); err != nil {
		return err
	}
	// 设置 SSE 响应头
	setSSEHeaders(c)
	// 处理流式响应
	fullResponse, err := handleSSEStream(c, resp.Body)
	if err != nil {
		return err
	}
	// 保存完整的会话到 Redis
	if err := saveConversationWithAIResponse(conversation, fullResponse); err != nil {
		return err
	}
	// 发送完成消息
	sendStreamEndMessage(c, fullResponse)

	return nil
}

// getConversationWithMessage 获取会话并添加用户消息
func getConversationWithMessage(conversationID int64, message string) (*models.Conversation, error) {
	// 从 Redis 获取会话
	conversation, err := storage.GetConversationFromRedis(conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to load conversation: %v", err)
	}
	if conversation == nil {
		return nil, fmt.Errorf("conversation not found")
	}
	// 追加用户消息
	userMessage := models.Message{
		Role:      "user",
		Content:   message,
		MessageID: int32(len(conversation.Messages)),
	}
	conversation.Messages = append(conversation.Messages, userMessage)

	// 将用户消息追加到 Redis
	err = storage.SaveConversationToRedis(conversation)
	if err != nil {
		return nil, fmt.Errorf("failed to append user message: %v", err)
	}

	return conversation, nil
}

// buildRequestBody 构造 API 请求体
func buildRequestBody(conversation *models.Conversation) ([]byte, error) {
	requestBody := map[string]interface{}{
		"model":    conversation.Model,
		"messages": conversation.Messages,
		"stream":   true,
	}
	return json.Marshal(requestBody)
}

// sendAPIRequest 发送 API 请求
func sendAPIRequest(apiKey string, requestData []byte, model string) (*http.Response, error) {
	client := &http.Client{}
	baseURL, err := utils.GetBaseURL(model)
	if err != nil {
		return nil, err
	}
	apiReq, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(requestData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	apiReq.Header.Set("Content-Type", "application/json")
	apiReq.Header.Set("Authorization", "Bearer "+apiKey)

	return client.Do(apiReq)
}

// validateResponse 验证 API 响应状态
func validateResponse(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected response status: %s", string(body))
	}
	return nil
}

// setSSEHeaders 设置 SSE 响应头
func setSSEHeaders(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
}

// handleSSEStream 处理流式 SSE 数据
func handleSSEStream(c *gin.Context, body io.Reader) (string, error) {
	reader := bufio.NewReader(body)
	var fullResponse string

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", fmt.Errorf("error reading stream: %w", err)
		}

		if len(line) == 0 || line[0] == ':' {
			continue
		}

		if bytes.HasPrefix(line, []byte("data: ")) {
			data := bytes.TrimPrefix(line, []byte("data: "))
			data = bytes.TrimSpace(data)

			if string(data) == "[DONE]" {
				break
			}

			fullResponse, err = processSSEData(c, data, fullResponse)
			if err != nil {
				return "", err
			}
		}
	}

	// 发送流式完成消息
	sendSSEEvent(c, "done", "Stream finished")
	return fullResponse, nil
}

// processSSEData 处理单条 SSE 数据
func processSSEData(c *gin.Context, data []byte, fullResponse string) (string, error) {
	var sseResponse *models.SSEResponse
	if err := json.Unmarshal(data, &sseResponse); err != nil {
		sendSSEEvent(c, "error", fmt.Sprintf("Failed to unmarshal SSE data: %v", err))
		return fullResponse, nil
	}

	for _, choice := range sseResponse.Choices {
		content := choice.Delta.Content
		fullResponse += content
		sendSSEEvent(c, "message", content)
	}
	return fullResponse, nil
}

// sendSSEEvent 发送 SSE 消息到客户端
func sendSSEEvent(c *gin.Context, event, data string) {
	message, _ := json.Marshal(map[string]string{
		"event": event,
		"data":  data,
	})
	fmt.Fprintf(c.Writer, "%s\n\n", message)
	c.Writer.Flush()
}

func saveConversationWithAIResponse(conversation *models.Conversation, fullResponse string) error {
	// 构造 AI 回复消息
	aiMessage := models.Message{Role: "assistant", Content: fullResponse, MessageID: int32(len(conversation.Messages))}
	// 追加到会话记录
	conversation.Messages = append(conversation.Messages, aiMessage)
	// 保存对话记录到 Redis
	return storage.SaveConversationToRedis(conversation)
}

// sendStreamEndMessage 发送流结束消息
func sendStreamEndMessage(c *gin.Context, fullResponse string) {
	endMessage, _ := json.Marshal(map[string]string{
		"event": "done",
		"data":  "Stream finished",
	})
	fmt.Fprintf(c.Writer, "%s\n\n", endMessage)
	c.Writer.Flush()

	fmt.Fprintf(c.Writer, "{\"event\":\"full_response\",\"data\":\"%s\"}\n\n", fullResponse)
	c.Writer.Flush()
}
