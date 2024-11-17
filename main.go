package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/EthanGuo-coder/llm-backend-api/constant"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SSEDelta 定义结构体以匹配 JSON 数据格式
type SSEDelta struct {
	Content string `json:"content"`
}

type SSEChoice struct {
	Delta SSEDelta `json:"delta"`
}

type SSEResponse struct {
	Choices []SSEChoice `json:"choices"`
}

func streamChat(c *gin.Context) {
	// 设置响应头为 SSE
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	// 接收参数
	var req struct {
		Model  string `json:"model" binding:"required"`
		ApiKey string `json:"api_key" binding:"required"`
		Query  string `json:"query" binding:"required"` // 新增 query 参数
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 构造请求的消息体
	messages := []map[string]string{
		{"role": "system", "content": "你是一个乐于回答各种问题的小助手"},
		{"role": "user", "content": req.Query}, // 用户传入的问题作为内容
	}

	requestBody := map[string]interface{}{
		"model":    req.Model,
		"messages": messages,
		"stream":   true, // 开启SSE流式返回
	}

	// 将请求体转换为JSON格式
	requestData, err := json.Marshal(requestBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request body"})
		return
	}

	// 创建HTTP请求
	client := &http.Client{}
	apiReq, err := http.NewRequest("POST", constant.BaseURL, bytes.NewBuffer(requestData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// 设置请求头
	apiReq.Header.Set("Content-Type", "application/json")
	apiReq.Header.Set("Authorization", "Bearer "+req.ApiKey)

	// 发送HTTP请求
	resp, err := client.Do(apiReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
		return
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.JSON(resp.StatusCode, gin.H{"error": "Unexpected response status", "details": string(body)})
		return
	}

	// 解析SSE流数据并实时以JSON格式返回
	reader := bufio.NewReader(resp.Body)
	var fullResponse string // 用于存储完整的返回信息
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading stream", "details": err.Error()})
			return
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
			var sseResponse SSEResponse
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
}

func main() {
	r := gin.Default()

	// 定义 POST 路由
	r.POST("/api/chat", streamChat)

	// 启动服务器
	r.Run(":8080") // 启动在8080端口
}
