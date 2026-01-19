package controllers

import (
	"log"
	"net/http"
	"yunmeng-backend/models"
	"yunmeng-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	db *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{db: db}
}

// RefreshToken 刷新Token接口
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	// 从请求头获取当前Token
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未提供认证令牌",
		})
		return
	}

	// 解析当前Token
	claims, err := utils.ParseToken(authHeader[7:]) // 去掉"Bearer "前缀
	if err != nil {
		log.Printf("Token解析失败: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "Token无效或已过期",
		})
		return
	}

	// 生成新的Token
	newToken, err := utils.GenerateToken(claims.UserID, claims.Username, claims.Role)
	if err != nil {
		log.Printf("新Token生成失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Token刷新失败",
		})
		return
	}

	log.Printf("Token刷新成功 - 用户: %s, 角色: %s", claims.Username, claims.Role)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Token刷新成功",
		"token":   newToken,
	})
}

// ValidateToken 验证Token接口
func (c *AuthController) ValidateToken(ctx *gin.Context) {
	// 从请求头获取Token
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未提供认证令牌",
		})
		return
	}

	// 解析Token
	claims, err := utils.ParseToken(authHeader[7:]) // 去掉"Bearer "前缀
	if err != nil {
		log.Printf("Token验证失败: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "Token无效或已过期",
		})
		return
	}

	// 检查用户是否存在
	var userCount int64
	c.db.Model(&models.User{}).Where("id = ? AND status = ?", claims.UserID, "active").Count(&userCount)
	if userCount == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户不存在或已被禁用",
		})
		return
	}

	log.Printf("Token验证成功 - 用户: %s, 角色: %s", claims.Username, claims.Role)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Token有效",
		"data": gin.H{
			"user_id":  claims.UserID,
			"username": claims.Username,
			"role":     claims.Role,
		},
	})
}

// GetUserInfo 获取用户信息接口
func (c *AuthController) GetUserInfo(ctx *gin.Context) {
	// 从上下文获取用户信息（由中间件设置）
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户未登录",
		})
		return
	}

	// 查询用户详细信息
	var user models.User
	if err := c.db.Preload("Profile").First(&user, userID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取用户信息失败",
		})
		return
	}

	// 获取用户角色
	var roles []models.Role
	if err := c.db.Model(&user).Association("Roles").Find(&roles); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取用户角色失败",
		})
		return
	}

	var roleNames []string
	for _, role := range roles {
		roleNames = append(roleNames, role.RoleKey)
	}

	userInfo := gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"status":     user.Status,
		"roles":      roleNames,
		"real_name":  user.Profile.RealName,
		"department": user.Profile.Department,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取用户信息成功",
		"data":    userInfo,
	})
}
