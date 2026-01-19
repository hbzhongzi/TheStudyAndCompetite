package utils

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ResponseError 返回错误响应
func ResponseError(c *gin.Context, statusCode int, message string, err error) {
	errorMsg := message
	if err != nil {
		errorMsg = message + ": " + err.Error()
	}

	c.JSON(statusCode, gin.H{
		"code":    statusCode,
		"message": errorMsg,
	})
}

// ResponseSuccess 返回成功响应
func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "操作成功",
		"data":    data,
	})
}

// GetCurrentUserID 从上下文中获取当前用户ID
func GetCurrentUserID(c *gin.Context) uint {
	userID, exists := c.Get("userID")
	if !exists {
		return 0
	}

	if id, ok := userID.(uint); ok {
		return id
	}

	return 0
}

// GenerateFileName 生成唯一的文件名
func GenerateFileName(originalName string) string {
	timestamp := time.Now().UnixNano()
	ext := ""
	if idx := strings.LastIndex(originalName, "."); idx != -1 {
		ext = originalName[idx:]
	}
	return fmt.Sprintf("%d%s", timestamp, ext)
}
