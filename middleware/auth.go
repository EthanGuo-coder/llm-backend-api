package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/EthanGuo-coder/llm-backend-api/storage"
	"github.com/EthanGuo-coder/llm-backend-api/utils"
)

// AuthMiddleware 用户认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// 提取 Bearer Token
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		// 优先从缓存获取 Token
		cachedData, err := storage.GetCachedJWT(tokenStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch token from cache"})
			c.Abort()
			return
		}

		var userID int64
		if cachedData != nil {
			// 从缓存中解析 user_id
			userID = int64(cachedData["user_id"].(float64))
		} else {
			// 缓存未命中，解析 Token
			claims, err := utils.ParseToken(tokenStr)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
				c.Abort()
				return
			}

			userID = int64(claims["user_id"].(float64))

			// 将解析结果缓存到 Redis，过期时间与 JWT 的剩余有效期一致
			expTime := time.Unix(int64(claims["exp"].(float64)), 0)
			ttl := time.Until(expTime)
			if ttl > 0 {
				if err := storage.CacheJWT(tokenStr, userID, ttl); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cache token"})
					c.Abort()
					return
				}
			}
		}

		// 将 user_id 存入上下文
		c.Set("user_id", userID)
		c.Next()
	}
}
