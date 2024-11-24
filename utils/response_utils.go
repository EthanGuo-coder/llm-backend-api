package utils

import "github.com/gin-gonic/gin"

func JSONError(c *gin.Context, code int, err error) {
	c.JSON(code, gin.H{"error": err.Error()})
}
