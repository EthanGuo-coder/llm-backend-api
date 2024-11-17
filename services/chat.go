package services

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/EthanGuo-coder/llm-backend-api/constant"
	"github.com/EthanGuo-coder/llm-backend-api/models"
	"github.com/EthanGuo-coder/llm-backend-api/storage"
	"github.com/gin-gonic/gin"
)

// StreamSendMessage 处理流式消息发送
func StreamSendMessage(c *gin.Context, conversationID, model, apiKey, message string) error {
	// 获取会话上下文
	conversation, err := storage.GetConversation(conversationID)
	if err != nil || conversation == nil {
		return fmt.Errorf("conversation not found")
	}

	// 用户消息添加到上下文
	userMessage := models.Message{Role: "user", Content: message}
	conversation.Messages = append(conversation.Messages, userMessage)

	// 构造 AI 请求体
	requestBody := map[string]interface{}{
		"model":    model,
		"messages": conversation.Messages, // 包含上下文的消息
		"stream":   true,
	}
	requestData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	// 创建 HTTP 请求
	client := &http.Client{}
	apiReq, err := http.NewRequest("POST", constant.BaseURL, bytes.NewBuffer(requestData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// 请求头
	apiReq.Header.Set("Content-Type", "application/json")
	apiReq.Header.Set("Authorization", "Bearer "+apiKey)

	// 发送 HTTP 请求
	resp, err := client.Do(apiReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected response status: %s", string(body))
	}

	// 设置 SSE 响应头
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	// 解析 SSE 数据
	reader := bufio.NewReader(resp.Body)
	var fullResponse string
	var wg sync.WaitGroup

	// 存储用户消息到 Redis
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := storage.SaveMessages(conversationID, []models.Message{userMessage}); err != nil {
			fmt.Printf("Failed to save user message: %v\n", err)
		}
	}()

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error reading stream: %w", err)
		}

		// 忽略注释行或空行
		if len(line) == 0 || line[0] == ':' {
			continue
		}

		// 解析SSE数据行
		if bytes.HasPrefix(line, []byte("data: ")) {
			data := bytes.TrimPrefix(line, []byte("data: "))
			data = bytes.TrimSpace(data)

			// 检测结束标志
			if string(data) == "[DONE]" {
				// 发送结束事件
				endMessage, _ := json.Marshal(map[string]string{
					"event": "done",
					"data":  "Stream finished",
				})
				fmt.Fprintf(c.Writer, "%s\n\n", endMessage)
				c.Writer.Flush()

				// 打印完整的返回信息
				fmt.Fprintf(c.Writer, "{\"event\":\"full_response\",\"data\":\"%s\"}\n\n", fullResponse)
				c.Writer.Flush()
				break
			}

			// 使用结构体解析JSON格式数据
			var sseResponse *models.SSEResponse
			if err := json.Unmarshal(data, &sseResponse); err != nil {
				errorMessage, _ := json.Marshal(map[string]string{
					"event": "error",
					"data":  fmt.Sprintf("Failed to unmarshal SSE data: %v", err),
				})
				fmt.Fprintf(c.Writer, "%s\n\n", errorMessage)
				c.Writer.Flush()
				continue
			}

			// 将每个增量内容封装为JSON
			for _, choice := range sseResponse.Choices {
				content := choice.Delta.Content
				fullResponse += content // 累加内容到完整响应
				message, _ := json.Marshal(map[string]string{
					"event": "message",
					"data":  content,
				})
				fmt.Fprintf(c.Writer, "%s\n\n", message)
				c.Writer.Flush() // 确保数据实时发送
			}
		}
	}

	// 保存完整 AI 回复
	wg.Add(1)
	go func() {
		defer wg.Done()
		aiMessage := models.Message{Role: "assistant", Content: fullResponse}
		if err := storage.SaveMessages(conversationID, []models.Message{aiMessage}); err != nil {
			fmt.Printf("Failed to save AI message: %v\n", err)
		}
	}()

	wg.Wait()
	return nil
}
