package controllers

import (
	"log"
	"net/http"
	"strconv"
	"yunmeng-backend/models"
	"yunmeng-backend/services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// GetUserList 获取用户列表
func (c *UserController) GetUserList(ctx *gin.Context) {
	var params models.UserQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		log.Printf("获取用户列表参数错误: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 设置默认值
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Size <= 0 {
		params.Size = 20
	}

	log.Printf("获取用户列表 - 页码: %d, 每页数量: %d, 搜索: %s, 角色: %s, 状态: %s",
		params.Page, params.Size, params.Search, params.Role, params.Status)

	users, total, err := c.userService.GetUserList(params)
	if err != nil {
		log.Printf("获取用户列表失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取用户列表失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	totalPages := (total + int64(params.Size) - 1) / int64(params.Size)

	log.Printf("用户列表获取成功 - 总数: %d, 当前页: %d, 总页数: %d", total, params.Page, totalPages)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取用户列表成功",
		"data": gin.H{
			"list":  users,
			"total": total,
			"page":  params.Page,
			"size":  params.Size,
			"pages": totalPages,
		},
	})
}

// GetUserByID 根据ID获取用户详情
func (c *UserController) GetUserByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("无效的用户ID: %s", idStr)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的用户ID",
			"data":    nil,
		})
		return
	}

	log.Printf("获取用户详情 - 用户ID: %d", id)

	user, err := c.userService.GetUserByID(uint(id))
	if err != nil {
		log.Printf("获取用户详情失败 - 用户ID: %d, 错误: %v", id, err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "用户不存在: " + err.Error(),
			"data":    nil,
		})
		return
	}

	log.Printf("用户详情获取成功 - 用户ID: %d, 用户名: %s", id, user.Username)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取用户详情成功",
		"data":    user,
	})
}

// CreateUser 创建用户
func (c *UserController) CreateUser(ctx *gin.Context) {
	var req models.UserCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("创建用户参数错误: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	log.Printf("创建用户 - 用户名: %s, 邮箱: %s, 角色: %v", req.Username, req.Email, req.RoleKeys)

	user, err := c.userService.CreateUser(req)
	if err != nil {
		log.Printf("创建用户失败 - 用户名: %s, 错误: %v", req.Username, err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "创建用户失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	log.Printf("用户创建成功 - 用户ID: %d, 用户名: %s", user.ID, user.Username)

	ctx.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"message": "用户创建成功",
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"status":   user.Status,
		},
	})
}

// UpdateUser 更新用户信息
func (c *UserController) UpdateUser(ctx *gin.Context) {
	// 获取当前用户信息
	currentUserID, exists := ctx.Get("userID")
	if !exists {
		log.Printf("UpdateUser失败 - 未找到当前用户ID")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未找到用户认证信息",
			"data":    nil,
		})
		return
	}

	currentUserRole, exists := ctx.Get("role")
	if !exists {
		log.Printf("UpdateUser失败 - 未找到当前用户角色")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未找到用户角色信息",
			"data":    nil,
		})
		return
	}

	log.Printf("UpdateUser权限检查 - 当前用户ID: %v, 角色: %v", currentUserID, currentUserRole)

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("无效的用户ID: %s", idStr)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的用户ID",
			"data":    nil,
		})
		return
	}

	var req models.UserUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("更新用户参数错误 - 用户ID: %d, 错误: %v", id, err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	log.Printf("更新用户 - 目标用户ID: %d, 当前用户ID: %v, 角色: %v, 真实姓名: %s, 邮箱: %s",
		id, currentUserID, currentUserRole, req.RealName, req.Email)

	err = c.userService.UpdateUser(uint(id), req)
	if err != nil {
		log.Printf("更新用户失败 - 用户ID: %d, 错误: %v", id, err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "更新用户失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	log.Printf("用户更新成功 - 用户ID: %d", id)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "用户信息更新成功",
		"data":    nil,
	})
}

// DeleteUser 删除用户
func (c *UserController) DeleteUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("无效的用户ID: %s", idStr)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的用户ID",
			"data":    nil,
		})
		return
	}

	log.Printf("删除用户 - 用户ID: %d", id)

	err = c.userService.DeleteUser(uint(id))
	if err != nil {
		log.Printf("删除用户失败 - 用户ID: %d, 错误: %v", id, err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "删除用户失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	log.Printf("用户删除成功 - 用户ID: %d", id)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "用户删除成功",
		"data":    nil,
	})
}

// ToggleUserStatus 切换用户状态
func (c *UserController) ToggleUserStatus(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("无效的用户ID: %s", idStr)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的用户ID",
			"data":    nil,
		})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=active inactive"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("切换用户状态参数错误 - 用户ID: %d, 错误: %v", id, err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	log.Printf("切换用户状态 - 用户ID: %d, 新状态: %s", id, req.Status)

	err = c.userService.ToggleUserStatus(uint(id), req.Status)
	if err != nil {
		log.Printf("切换用户状态失败 - 用户ID: %d, 错误: %v", id, err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "更新用户状态失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	log.Printf("用户状态更新成功 - 用户ID: %d, 状态: %s", id, req.Status)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "用户状态更新成功",
		"data": gin.H{
			"id":     id,
			"status": req.Status,
		},
	})
}

// ResetUserPassword 重置用户密码
func (c *UserController) ResetUserPassword(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("无效的用户ID: %s", idStr)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的用户ID",
			"data":    nil,
		})
		return
	}

	log.Printf("重置用户密码 - 用户ID: %d", id)

	newPassword, err := c.userService.ResetUserPassword(uint(id))
	if err != nil {
		log.Printf("重置用户密码失败 - 用户ID: %d, 错误: %v", id, err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "重置密码失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	log.Printf("用户密码重置成功 - 用户ID: %d", id)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "密码重置成功",
		"data": gin.H{
			"id":          id,
			"newPassword": newPassword, // 实际项目中不应该返回密码
		},
	})
}

// BatchDeleteUsers 批量删除用户
func (c *UserController) BatchDeleteUsers(ctx *gin.Context) {
	var req struct {
		UserIDs []uint `json:"userIds" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("批量删除用户参数错误: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	log.Printf("批量删除用户 - 用户ID列表: %v", req.UserIDs)

	deletedCount, err := c.userService.BatchDeleteUsers(req.UserIDs)
	if err != nil {
		log.Printf("批量删除用户失败 - 错误: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "批量删除失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	log.Printf("批量删除用户成功 - 删除数量: %d", deletedCount)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "批量删除成功",
		"data": gin.H{
			"deletedCount": deletedCount,
			"deletedIds":   req.UserIDs,
		},
	})
}

// GetUserStats 获取用户统计信息
func (c *UserController) GetUserStats(ctx *gin.Context) {
	log.Printf("获取用户统计信息")

	stats, err := c.userService.GetUserStats()
	if err != nil {
		log.Printf("获取用户统计信息失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取统计信息失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	log.Printf("用户统计信息获取成功 - 总用户数: %d", stats.TotalUsers)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取统计信息成功",
		"data":    stats,
	})
}

// ExportUsers 导出用户数据
func (c *UserController) ExportUsers(ctx *gin.Context) {
	var params models.UserQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		log.Printf("导出用户数据参数错误: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	format := ctx.Query("format")
	if format == "" {
		format = "json"
	}

	log.Printf("导出用户数据 - 格式: %s, 角色: %s, 状态: %s", format, params.Role, params.Status)

	// 获取所有用户数据（不分页）
	params.Page = 0
	params.Size = 0
	users, _, err := c.userService.GetUserList(params)
	if err != nil {
		log.Printf("导出用户数据失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "导出用户数据失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	log.Printf("用户数据导出成功 - 用户数量: %d", len(users))

	// 根据格式返回不同的响应
	switch format {
	case "excel":
		// TODO: 实现Excel导出
		ctx.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "Excel导出功能开发中",
			"data":    users,
		})
	case "csv":
		// TODO: 实现CSV导出
		ctx.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "CSV导出功能开发中",
			"data":    users,
		})
	default:
		// 返回JSON格式
		ctx.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "导出用户数据成功",
			"data":    users,
		})
	}
}
