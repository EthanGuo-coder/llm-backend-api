package storage

import "fmt"

// RedisKey 用于存储 Redis 的键模板
const (
	RedisKeyConversation = "conversation:%s" // 会话的键
	RedisKeyJWT          = "jwt:%s"          // JWT 的键
)

// GenerateRedisKeyConversation 生成会话的 Redis 键
func GenerateRedisKeyConversation(conversationID string) string {
	return fmt.Sprintf(RedisKeyConversation, conversationID)
}

// GenerateRedisKeyJWT 生成 JWT 的 Redis 键
func GenerateRedisKeyJWT(token string) string {
	return fmt.Sprintf(RedisKeyJWT, token)
}
