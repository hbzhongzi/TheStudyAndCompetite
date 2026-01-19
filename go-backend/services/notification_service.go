package services

import (
	"errors"
	"time"

	"yunmeng-backend/models"

	"gorm.io/gorm"
)

type NotificationService struct {
	db *gorm.DB
}

func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{db: db}
}

// GetMyNotifications 获取我的通知列表
func (s *NotificationService) GetMyNotifications(userID uint, params models.NotificationQueryParams) ([]models.ProjectNotificationResponse, int64, error) {
	var notifications []models.ProjectNotification
	var total int64

	query := s.db.Model(&models.ProjectNotification{}).Where("user_id = ?", userID)

	// 应用查询参数
	if params.Type != "" {
		query = query.Where("type = ?", params.Type)
	}
	if params.IsRead != nil {
		query = query.Where("is_read = ?", *params.IsRead)
	}
	if params.Priority != "" {
		query = query.Where("priority = ?", params.Priority)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	if params.SortBy != "" {
		order := params.SortBy
		if params.SortOrder == "desc" {
			order += " DESC"
		}
		query = query.Order(order)
	} else {
		query = query.Order("created_at DESC")
	}

	// 分页
	if params.Page > 0 && params.Size > 0 {
		offset := (params.Page - 1) * params.Size
		query = query.Offset(offset).Limit(params.Size)
	}

	// 执行查询
	if err := query.Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	var responses []models.ProjectNotificationResponse
	for _, notification := range notifications {
		responses = append(responses, models.ProjectNotificationResponse{
			ID:        notification.ID,
			ProjectID: notification.ProjectID,
			Type:      notification.Type,
			Title:     notification.Title,
			Content:   notification.Content,
			IsRead:    notification.IsRead,
			Priority:  notification.Priority,
			CreatedAt: notification.CreatedAt,
		})
	}

	return responses, total, nil
}

// MarkNotificationAsRead 标记通知为已读
func (s *NotificationService) MarkNotificationAsRead(notificationID, userID uint) error {
	result := s.db.Model(&models.ProjectNotification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("通知不存在或无权限")
	}

	return nil
}

// MarkAllNotificationsAsRead 标记所有通知为已读
func (s *NotificationService) MarkAllNotificationsAsRead(userID uint) error {
	result := s.db.Model(&models.ProjectNotification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// GetUnreadCount 获取未读通知数量
func (s *NotificationService) GetUnreadCount(userID uint) (int64, error) {
	var count int64
	err := s.db.Model(&models.ProjectNotification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error

	return count, err
}

// DeleteNotification 删除通知
func (s *NotificationService) DeleteNotification(notificationID, userID uint) error {
	result := s.db.Where("id = ? AND user_id = ?", notificationID, userID).
		Delete(&models.ProjectNotification{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("通知不存在或无权限")
	}

	return nil
}

// GetNotificationTemplates 获取通知模板列表
func (s *NotificationService) GetNotificationTemplates() ([]models.NotificationTemplateResponse, error) {
	var templates []models.NotificationTemplate
	err := s.db.Where("is_active = ?", true).Order("id").Find(&templates).Error
	if err != nil {
		return nil, err
	}

	var responses []models.NotificationTemplateResponse
	for _, template := range templates {
		responses = append(responses, models.NotificationTemplateResponse{
			ID:              template.ID,
			TemplateKey:     template.TemplateKey,
			TitleTemplate:   template.TitleTemplate,
			ContentTemplate: template.ContentTemplate,
			Variables:       template.Variables,
			IsActive:        template.IsActive,
			CreatedAt:       template.CreatedAt,
			UpdatedAt:       template.UpdatedAt,
		})
	}

	return responses, nil
}

// UpdateNotificationTemplate 更新通知模板
func (s *NotificationService) UpdateNotificationTemplate(templateID uint, req models.NotificationTemplateUpdateRequest) error {
	result := s.db.Model(&models.NotificationTemplate{}).
		Where("id = ?", templateID).
		Updates(map[string]interface{}{
			"title_template":   req.TitleTemplate,
			"content_template": req.ContentTemplate,
			"variables":        req.Variables,
			"is_active":        req.IsActive,
			"updated_at":       time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("通知模板不存在")
	}

	return nil
}

// SendNotification 发送通知
func (s *NotificationService) SendNotification(senderID uint, req models.NotificationSendRequest) (*models.ProjectNotificationResponse, error) {
	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 为每个用户创建通知
	for _, userID := range req.UserIDs {
		notification := models.ProjectNotification{
			ProjectID: req.ProjectID,
			UserID:    userID,
			Type:      req.Type,
			Title:     req.Title,
			Content:   req.Content,
			Priority:  req.Priority,
			IsRead:    false,
		}

		if err := tx.Create(&notification).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// 返回第一个通知作为响应
	return &models.ProjectNotificationResponse{
		ProjectID: req.ProjectID,
		Type:      req.Type,
		Title:     req.Title,
		Content:   req.Content,
		Priority:  req.Priority,
		IsRead:    false,
		CreatedAt: time.Now(),
	}, nil
}

// GetNotificationStats 获取通知统计信息
func (s *NotificationService) GetNotificationStats() (map[string]interface{}, error) {
	var stats map[string]interface{}

	// 总通知数
	var totalCount int64
	if err := s.db.Model(&models.ProjectNotification{}).Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// 未读通知数
	var unreadCount int64
	if err := s.db.Model(&models.ProjectNotification{}).Where("is_read = ?", false).Count(&unreadCount).Error; err != nil {
		return nil, err
	}

	// 按类型统计
	var typeStats []struct {
		Type  string `json:"type"`
		Count int64  `json:"count"`
	}
	if err := s.db.Model(&models.ProjectNotification{}).
		Select("type, count(*) as count").
		Group("type").
		Find(&typeStats).Error; err != nil {
		return nil, err
	}

	// 按优先级统计
	var priorityStats []struct {
		Priority string `json:"priority"`
		Count    int64  `json:"count"`
	}
	if err := s.db.Model(&models.ProjectNotification{}).
		Select("priority, count(*) as count").
		Group("priority").
		Find(&priorityStats).Error; err != nil {
		return nil, err
	}

	stats = map[string]interface{}{
		"totalCount":    totalCount,
		"unreadCount":   unreadCount,
		"typeStats":     typeStats,
		"priorityStats": priorityStats,
		"readRate":      float64(totalCount-unreadCount) / float64(totalCount) * 100,
	}

	return stats, nil
}

// BatchSendNotifications 批量发送通知
func (s *NotificationService) BatchSendNotifications(notifications []models.NotificationSendRequest) error {
	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, req := range notifications {
		for _, userID := range req.UserIDs {
			notification := models.ProjectNotification{
				ProjectID: req.ProjectID,
				UserID:    userID,
				Type:      req.Type,
				Title:     req.Title,
				Content:   req.Content,
				Priority:  req.Priority,
				IsRead:    false,
			}

			if err := tx.Create(&notification).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// 提交事务
	return tx.Commit().Error
}

// GetNotificationHistory 获取通知发送历史
func (s *NotificationService) GetNotificationHistory(params models.NotificationQueryParams) ([]models.ProjectNotificationResponse, int64, error) {
	var notifications []models.ProjectNotification
	var total int64

	query := s.db.Model(&models.ProjectNotification{})

	// 应用查询参数
	if params.Type != "" {
		query = query.Where("type = ?", params.Type)
	}
	if params.Priority != "" {
		query = query.Where("priority = ?", params.Priority)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	if params.SortBy != "" {
		order := params.SortBy
		if params.SortOrder == "desc" {
			order += " DESC"
		}
		query = query.Order(order)
	} else {
		query = query.Order("created_at DESC")
	}

	// 分页
	if params.Page > 0 && params.Size > 0 {
		offset := (params.Page - 1) * params.Size
		query = query.Offset(offset).Limit(params.Size)
	}

	// 执行查询
	if err := query.Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	var responses []models.ProjectNotificationResponse
	for _, notification := range notifications {
		responses = append(responses, models.ProjectNotificationResponse{
			ID:        notification.ID,
			ProjectID: notification.ProjectID,
			Type:      notification.Type,
			Title:     notification.Title,
			Content:   notification.Content,
			IsRead:    notification.IsRead,
			Priority:  notification.Priority,
			CreatedAt: notification.CreatedAt,
		})
	}

	return responses, total, nil
}

// SetNotificationPreferences 设置通知偏好
func (s *NotificationService) SetNotificationPreferences(userID uint, preferences map[string]bool) error {
	// 这里可以实现用户通知偏好设置
	// 暂时返回成功
	return nil
}

// GetNotificationPreferences 获取通知偏好
func (s *NotificationService) GetNotificationPreferences(userID uint) (map[string]bool, error) {
	// 这里可以实现获取用户通知偏好
	// 暂时返回默认值
	return map[string]bool{
		"email":  true,
		"sms":    false,
		"push":   true,
		"in_app": true,
	}, nil
}

// SubscribeNotificationChannel 订阅通知频道
func (s *NotificationService) SubscribeNotificationChannel(userID, channelID uint) error {
	// 这里可以实现通知频道订阅
	// 暂时返回成功
	return nil
}

// UnsubscribeNotificationChannel 取消订阅通知频道
func (s *NotificationService) UnsubscribeNotificationChannel(userID, channelID uint) error {
	// 这里可以实现取消通知频道订阅
	// 暂时返回成功
	return nil
}

// GetAvailableNotificationChannels 获取可用的通知频道
func (s *NotificationService) GetAvailableNotificationChannels() ([]map[string]interface{}, error) {
	// 这里可以实现获取可用通知频道
	// 暂时返回默认值
	return []map[string]interface{}{
		{"id": 1, "name": "系统通知", "description": "系统级别的通知"},
		{"id": 2, "name": "项目通知", "description": "项目相关的通知"},
		{"id": 3, "name": "审核通知", "description": "审核流程相关的通知"},
	}, nil
}

// TestNotification 测试通知发送
func (s *NotificationService) TestNotification(templateID uint, testData map[string]string) error {
	// 这里可以实现测试通知发送
	// 暂时返回成功
	return nil
}
