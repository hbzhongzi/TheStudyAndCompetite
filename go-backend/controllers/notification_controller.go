package controllers

import (
	"log"
	"net/http"
	"strconv"

	"yunmeng-backend/models"
	"yunmeng-backend/services"

	"github.com/gin-gonic/gin"
)

type NotificationController struct {
	notificationService *services.NotificationService
}

func NewNotificationController(notificationService *services.NotificationService) *NotificationController {
	return &NotificationController{
		notificationService: notificationService,
	}
}

// GetMyNotifications 获取我的通知列表
func (c *NotificationController) GetMyNotifications(ctx *gin.Context) {
	var params models.NotificationQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		log.Printf("参数绑定失败: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 设置默认分页参数
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Size <= 0 {
		params.Size = 20
	}

	userID := ctx.GetUint("user_id")
	notifications, total, err := c.notificationService.GetMyNotifications(userID, params)
	if err != nil {
		log.Printf("获取通知列表失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取通知列表失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取通知列表成功",
		"data": gin.H{
			"list":  notifications,
			"total": total,
		},
	})
}

// MarkNotificationAsRead 标记通知为已读
func (c *NotificationController) MarkNotificationAsRead(ctx *gin.Context) {
	notificationIDStr := ctx.Param("id")
	notificationID, err := strconv.ParseUint(notificationIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "通知ID格式错误",
		})
		return
	}

	userID := ctx.GetUint("user_id")
	err = c.notificationService.MarkNotificationAsRead(uint(notificationID), userID)
	if err != nil {
		log.Printf("标记通知为已读失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "标记通知为已读失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "标记通知为已读成功",
	})
}

// MarkAllNotificationsAsRead 标记所有通知为已读
func (c *NotificationController) MarkAllNotificationsAsRead(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	err := c.notificationService.MarkAllNotificationsAsRead(userID)
	if err != nil {
		log.Printf("标记所有通知为已读失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "标记所有通知为已读失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "标记所有通知为已读成功",
	})
}

// GetUnreadCount 获取未读通知数量
func (c *NotificationController) GetUnreadCount(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	count, err := c.notificationService.GetUnreadCount(userID)
	if err != nil {
		log.Printf("获取未读通知数量失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取未读通知数量失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取未读通知数量成功",
		"data": gin.H{
			"unreadCount": count,
		},
	})
}

// DeleteNotification 删除通知
func (c *NotificationController) DeleteNotification(ctx *gin.Context) {
	notificationIDStr := ctx.Param("id")
	notificationID, err := strconv.ParseUint(notificationIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "通知ID格式错误",
		})
		return
	}

	userID := ctx.GetUint("user_id")
	err = c.notificationService.DeleteNotification(uint(notificationID), userID)
	if err != nil {
		log.Printf("删除通知失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除通知失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除通知成功",
	})
}

// GetNotificationTemplates 获取通知模板列表（管理员）
func (c *NotificationController) GetNotificationTemplates(ctx *gin.Context) {
	templates, err := c.notificationService.GetNotificationTemplates()
	if err != nil {
		log.Printf("获取通知模板失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取通知模板失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取通知模板成功",
		"data":    templates,
	})
}

// UpdateNotificationTemplate 更新通知模板（管理员）
func (c *NotificationController) UpdateNotificationTemplate(ctx *gin.Context) {
	templateIDStr := ctx.Param("id")
	templateID, err := strconv.ParseUint(templateIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "模板ID格式错误",
		})
		return
	}

	var req models.NotificationTemplateUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	err = c.notificationService.UpdateNotificationTemplate(uint(templateID), req)
	if err != nil {
		log.Printf("更新通知模板失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新通知模板失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新通知模板成功",
	})
}

// SendNotification 发送通知（管理员/系统）
func (c *NotificationController) SendNotification(ctx *gin.Context) {
	var req models.NotificationSendRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	userID := ctx.GetUint("user_id")
	notification, err := c.notificationService.SendNotification(userID, req)
	if err != nil {
		log.Printf("发送通知失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "发送通知失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "发送通知成功",
		"data":    notification,
	})
}
