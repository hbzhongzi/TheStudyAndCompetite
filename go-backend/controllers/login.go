package controllers

import (
	"log"
	"net/http"
	"yunmeng-backend/models"
	"yunmeng-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type LoginResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Token   string      `json:"token,omitempty"`
}

func LoginHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		// 解析请求参数
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("参数解析失败: %v", err)
			c.JSON(http.StatusOK, LoginResponse{Code: 400, Message: "参数错误/缺失字段"})
			return
		}

		// 打印接收到的登录请求
		log.Printf("收到登录请求 - 用户名: %s, 角色: %s", req.Username, req.Role)

		var user models.User

		// 查询用户 - 通过用户名查找
		err := db.Where("username = ?", req.Username).First(&user).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				log.Printf("用户不存在 - 用户名: %s, 角色: %s", req.Username, req.Role)
				c.JSON(http.StatusOK, LoginResponse{Code: 401, Message: "账号或密码错误"})
			} else {
				log.Printf("数据库查询错误: %v", err)
				c.JSON(http.StatusOK, LoginResponse{Code: 500, Message: "服务端异常"})
			}
			return
		}

		// 打印查询到的用户信息（不包含密码）
		log.Printf("查询到用户 - ID: %d, 用户名: %s, 状态: %s",
			user.ID, user.Username, user.Status)

		// 检查用户状态
		if user.Status == "inactive" {
			log.Printf("用户被禁用 - 用户名: %s", req.Username)
			c.JSON(http.StatusOK, LoginResponse{Code: 403, Message: "账户已被禁用"})
			return
		}

		// 密码校验 - 使用bcrypt比较
		if !utils.CheckPassword(req.Password, user.Password) {
			log.Printf("密码校验失败 - 用户名: %s", req.Username)
			c.JSON(http.StatusOK, LoginResponse{Code: 401, Message: "账号或密码错误"})
			return
		}

		log.Printf("密码校验成功 - 用户名: %s", req.Username)

		// 获取用户角色
		var roles []models.Role
		if err := db.Model(&user).Association("Roles").Find(&roles); err != nil {
			log.Printf("获取用户角色失败: %v", err)
			c.JSON(http.StatusOK, LoginResponse{Code: 500, Message: "获取用户角色失败"})
			return
		}

		// 检查请求的角色是否匹配
		roleMatched := false
		var primaryRole string
		for _, role := range roles {
			if role.RoleKey == req.Role {
				roleMatched = true
				primaryRole = role.RoleKey
				break
			}
		}

		if !roleMatched {
			log.Printf("用户角色不匹配 - 用户名: %s, 请求角色: %s", req.Username, req.Role)
			c.JSON(http.StatusOK, LoginResponse{Code: 403, Message: "用户角色不匹配"})
			return
		}

		// 生成JWT Token
		token, err := utils.GenerateToken(user.ID, user.Username, primaryRole)
		if err != nil {
			log.Printf("Token生成失败: %v", err)
			c.JSON(http.StatusOK, LoginResponse{Code: 500, Message: "Token生成失败"})
			return
		}

		// 获取用户详细信息
		var profile models.UserProfile
		db.Where("user_id = ?", user.ID).First(&profile)

		// 返回用户信息（不含密码）
		userInfo := map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     primaryRole,
			"realName": profile.RealName,
		}

		log.Printf("登录成功 - 用户名: %s, 角色: %s", req.Username, req.Role)

		c.JSON(http.StatusOK, LoginResponse{
			Code:    200,
			Message: "登录成功",
			Data:    userInfo,
			Token:   token,
		})
	}
}
