package middlewares

import (
	"log"
	"net/http"
	"strings"
	"yunmeng-backend/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("AuthMiddleware失败 - 未提供认证令牌")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未提供认证令牌",
			})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			log.Printf("AuthMiddleware失败 - 认证令牌格式错误: %s", authHeader)
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "认证令牌格式错误",
			})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// 解析JWT Token
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			log.Printf("AuthMiddleware失败 - Token解析错误: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "认证令牌无效或已过期",
			})
			c.Abort()
			return
		}

		log.Printf("AuthMiddleware成功 - 用户ID: %d, 用户名: %s, 角色: %s",
			claims.UserID, claims.Username, claims.Role)

		// 将用户信息存储到上下文中
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// RoleMiddleware 角色权限中间件
func RoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未找到用户角色信息",
			})
			c.Abort()
			return
		}

		role := userRole.(string)
		log.Printf("权限检查 - 用户角色: %s, 需要角色: %v", role, requiredRoles)

		hasPermission := false

		// admin角色拥有所有权限
		if role == "admin" {
			hasPermission = true
			log.Printf("管理员权限验证通过")
		} else {
			// 检查用户角色是否在所需角色列表中
			for _, requiredRole := range requiredRoles {
				if role == requiredRole {
					hasPermission = true
					log.Printf("角色权限验证通过: %s", role)
					break
				}
			}
		}

		if !hasPermission {
			log.Printf("权限验证失败 - 用户角色: %s, 需要角色: %v", role, requiredRoles)
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "权限不足",
			})
			c.Abort()
			return
		}

		log.Printf("权限验证成功 - 用户角色: %s", role)
		c.Next()
	}
}

// AdminOnly 仅管理员权限
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			log.Printf("AdminOnly权限检查失败 - 未找到用户角色信息")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未找到用户角色信息",
			})
			c.Abort()
			return
		}

		role := userRole.(string)
		log.Printf("AdminOnly权限检查 - 用户角色: %s", role)

		if role != "admin" {
			log.Printf("AdminOnly权限检查失败 - 用户角色: %s, 需要管理员权限", role)
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "需要管理员权限",
			})
			c.Abort()
			return
		}

		log.Printf("AdminOnly权限检查通过 - 用户角色: %s", role)
		c.Next()
	}
}

// TeacherOrAdmin 教师或管理员权限
func TeacherOrAdmin() gin.HandlerFunc {
	return RoleMiddleware("admin", "teacher")
}
