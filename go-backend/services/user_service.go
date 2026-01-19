package services

import (
	"errors"
	"log"
	"time"
	"yunmeng-backend/models"
	"yunmeng-backend/utils"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// GetUserList 获取用户列表
func (s *UserService) GetUserList(params models.UserQueryParams) ([]models.UserListResponse, int64, error) {
	var users []models.User
	var total int64

	// 构建查询
	query := s.db.Model(&models.User{}).
		Preload("Profile").
		Preload("Roles")

	// 搜索条件
	if params.Search != "" {
		search := "%" + params.Search + "%"
		query = query.Where("users.username LIKE ? OR users.email LIKE ?", search, search)
	}

	// 角色筛选
	if params.Role != "" {
		query = query.Joins("JOIN user_roles ON users.id = user_roles.user_id").
			Joins("JOIN roles ON user_roles.role_id = roles.id").
			Where("roles.role_key = ?", params.Role)
	}

	// 状态筛选
	if params.Status != "" {
		query = query.Where("users.status = ?", params.Status)
	}

	// 部门筛选
	if params.Department != "" {
		query = query.Joins("LEFT JOIN user_profiles ON users.id = user_profiles.user_id").
			Where("user_profiles.department = ?", params.Department)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		log.Printf("获取用户总数失败: %v", err)
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
		query = query.Order("users.create_time DESC")
	}

	// 分页
	if params.Page > 0 && params.Size > 0 {
		offset := (params.Page - 1) * params.Size
		query = query.Offset(offset).Limit(params.Size)
	}

	// 执行查询
	if err := query.Find(&users).Error; err != nil {
		log.Printf("查询用户列表失败: %v", err)
		return nil, 0, err
	}

	// 转换为响应格式
	var responses []models.UserListResponse
	for _, user := range users {
		response := models.UserListResponse{
			ID:         user.ID,
			Username:   user.Username,
			Email:      user.Email,
			Status:     user.Status,
			CreateTime: user.CreateTime,
		}

		// 添加用户详细信息
		if user.Profile != nil && user.Profile.UserID != 0 {
			response.RealName = user.Profile.RealName
			response.Phone = user.Profile.Phone
			response.Department = user.Profile.Department
			response.StudentID = user.Profile.StudentID
			response.LastLogin = user.Profile.LastLogin
		}

		// 添加角色信息
		var roleNames []string
		for _, role := range user.Roles {
			roleNames = append(roleNames, role.RoleName)
		}
		response.RoleNames = roleNames

		responses = append(responses, response)
	}

	log.Printf("用户列表查询完成 - 总数: %d, 返回数量: %d", total, len(responses))
	return responses, total, nil
}

// GetUserByID 根据ID获取用户详情
func (s *UserService) GetUserByID(id uint) (*models.UserDetailResponse, error) {
	var user models.User
	err := s.db.Preload("Profile").Preload("Roles").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		log.Printf("查询用户详情失败 - 用户ID: %d, 错误: %v", id, err)
		return nil, err
	}

	// 构建响应数据
	response := &models.UserDetailResponse{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		Status:     user.Status,
		CreateTime: user.CreateTime,
		UpdateTime: user.UpdateTime,
	}

	// 添加用户详细信息
	if user.Profile != nil {
		response.Profile.RealName = user.Profile.RealName
		response.Profile.Phone = user.Profile.Phone
		response.Profile.Department = user.Profile.Department
		response.Profile.StudentID = user.Profile.StudentID
		response.Profile.Avatar = user.Profile.Avatar
		response.Profile.Bio = user.Profile.Bio
		response.Profile.Interests = user.Profile.Interests
		response.Profile.LastLogin = user.Profile.LastLogin
	}

	// 添加角色信息
	for _, role := range user.Roles {
		response.Roles = append(response.Roles, struct {
			ID          uint   `json:"id"`
			RoleKey     string `json:"roleKey"`
			RoleName    string `json:"roleName"`
			Description string `json:"description"`
		}{
			ID:          role.ID,
			RoleKey:     role.RoleKey,
			RoleName:    role.RoleName,
			Description: role.Description,
		})
	}

	log.Printf("用户详情查询完成 - 用户ID: %d, 用户名: %s", id, user.Username)
	return response, nil
}

// CreateUser 创建用户
func (s *UserService) CreateUser(req models.UserCreateRequest) (*models.User, error) {
	// 检查用户名是否已存在
	var existingUser models.User
	if err := s.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("邮箱已存在")
	}

	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建用户基础信息
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("密码加密失败")
	}

	user := models.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		Status:   "active",
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		log.Printf("创建用户失败: %v", err)
		return nil, errors.New("创建用户失败")
	}

	// 创建用户详细信息
	profile := models.UserProfile{
		UserID:     user.ID,
		RealName:   req.RealName,
		Phone:      req.Phone,
		Department: req.Department,
		StudentID:  req.StudentID,
	}

	if err := tx.Create(&profile).Error; err != nil {
		tx.Rollback()
		log.Printf("创建用户详细信息失败: %v", err)
		return nil, errors.New("创建用户详细信息失败")
	}

	// 分配角色
	for _, roleKey := range req.RoleKeys {
		var role models.Role
		if err := tx.Where("role_key = ?", roleKey).First(&role).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("角色不存在: " + roleKey)
		}

		userRole := models.UserRole{
			UserID: user.ID,
			RoleID: role.ID,
		}

		if err := tx.Create(&userRole).Error; err != nil {
			tx.Rollback()
			log.Printf("分配角色失败: %v", err)
			return nil, errors.New("分配角色失败")
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Printf("提交事务失败: %v", err)
		return nil, errors.New("创建用户失败")
	}

	log.Printf("用户创建成功 - 用户ID: %d, 用户名: %s", user.ID, user.Username)
	return &user, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(id uint, req models.UserUpdateRequest) error {
	// 检查用户是否存在
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新用户详细信息
	updates := make(map[string]interface{})
	if req.RealName != "" {
		updates["real_name"] = req.RealName
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Department != "" {
		updates["department"] = req.Department
	}
	if req.StudentID != "" {
		updates["student_id"] = req.StudentID
	}

	if len(updates) > 0 {
		if err := tx.Model(&models.UserProfile{}).Where("user_id = ?", id).Updates(updates).Error; err != nil {
			tx.Rollback()
			log.Printf("更新用户详细信息失败: %v", err)
			return errors.New("更新用户详细信息失败")
		}
	}

	// 更新角色（如果提供）
	if len(req.RoleKeys) > 0 {
		// 删除现有角色
		if err := tx.Where("user_id = ?", id).Delete(&models.UserRole{}).Error; err != nil {
			tx.Rollback()
			log.Printf("删除现有角色失败: %v", err)
			return errors.New("更新角色失败")
		}

		// 分配新角色
		for _, roleKey := range req.RoleKeys {
			var role models.Role
			if err := tx.Where("role_key = ?", roleKey).First(&role).Error; err != nil {
				tx.Rollback()
				return errors.New("角色不存在: " + roleKey)
			}

			userRole := models.UserRole{
				UserID: id,
				RoleID: role.ID,
			}

			if err := tx.Create(&userRole).Error; err != nil {
				tx.Rollback()
				log.Printf("分配新角色失败: %v", err)
				return errors.New("更新角色失败")
			}
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Printf("提交事务失败: %v", err)
		return errors.New("更新用户失败")
	}

	log.Printf("用户更新成功 - 用户ID: %d", id)
	return nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id uint) error {
	// 检查用户是否存在
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除用户（级联删除相关数据）
	if err := tx.Delete(&user).Error; err != nil {
		tx.Rollback()
		log.Printf("删除用户失败: %v", err)
		return errors.New("删除用户失败")
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Printf("提交事务失败: %v", err)
		return errors.New("删除用户失败")
	}

	log.Printf("用户删除成功 - 用户ID: %d", id)
	return nil
}

// ToggleUserStatus 切换用户状态
func (s *UserService) ToggleUserStatus(id uint, status string) error {
	// 检查用户是否存在
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	// 更新状态
	if err := s.db.Model(&user).Update("status", status).Error; err != nil {
		log.Printf("更新用户状态失败: %v", err)
		return errors.New("更新用户状态失败")
	}

	log.Printf("用户状态更新成功 - 用户ID: %d, 新状态: %s", id, status)
	return nil
}

// ResetUserPassword 重置用户密码
func (s *UserService) ResetUserPassword(id uint) (string, error) {
	// 检查用户是否存在
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("用户不存在")
		}
		return "", err
	}

	// 生成新密码
	newPassword := utils.GenerateRandomPassword(8)
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		log.Printf("密码加密失败: %v", err)
		return "", errors.New("密码加密失败")
	}

	// 更新密码
	if err := s.db.Model(&user).Update("password", hashedPassword).Error; err != nil {
		log.Printf("更新密码失败: %v", err)
		return "", errors.New("更新密码失败")
	}

	log.Printf("用户密码重置成功 - 用户ID: %d", id)
	return newPassword, nil
}

// BatchDeleteUsers 批量删除用户
func (s *UserService) BatchDeleteUsers(userIDs []uint) (int64, error) {
	if len(userIDs) == 0 {
		return 0, errors.New("用户ID列表不能为空")
	}

	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 批量删除用户
	result := tx.Where("id IN ?", userIDs).Delete(&models.User{})
	if result.Error != nil {
		tx.Rollback()
		log.Printf("批量删除用户失败: %v", result.Error)
		return 0, errors.New("批量删除用户失败")
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Printf("提交事务失败: %v", err)
		return 0, errors.New("批量删除用户失败")
	}

	log.Printf("批量删除用户成功 - 删除数量: %d", result.RowsAffected)
	return result.RowsAffected, nil
}

// GetUserStats 获取用户统计信息
func (s *UserService) GetUserStats() (*models.UserStats, error) {
	stats := &models.UserStats{
		RoleStats:       make(map[string]int64),
		DepartmentStats: make(map[string]int64),
	}

	// 获取总用户数
	if err := s.db.Model(&models.User{}).Count(&stats.TotalUsers).Error; err != nil {
		log.Printf("获取总用户数失败: %v", err)
		return nil, err
	}

	// 获取活跃用户数
	if err := s.db.Model(&models.User{}).Where("status = ?", "active").Count(&stats.ActiveUsers).Error; err != nil {
		log.Printf("获取活跃用户数失败: %v", err)
		return nil, err
	}

	// 获取禁用用户数
	if err := s.db.Model(&models.User{}).Where("status = ?", "inactive").Count(&stats.InactiveUsers).Error; err != nil {
		log.Printf("获取禁用用户数失败: %v", err)
		return nil, err
	}

	// 获取角色统计
	var roleStats []struct {
		RoleKey string `json:"roleKey"`
		Count   int64  `json:"count"`
	}
	if err := s.db.Table("roles").
		Select("roles.role_key, COUNT(user_roles.user_id) as count").
		Joins("LEFT JOIN user_roles ON roles.id = user_roles.role_id").
		Group("roles.id, roles.role_key").
		Scan(&roleStats).Error; err != nil {
		log.Printf("获取角色统计失败: %v", err)
		return nil, err
	}

	for _, rs := range roleStats {
		stats.RoleStats[rs.RoleKey] = rs.Count
	}

	// 获取部门统计
	var deptStats []struct {
		Department string `json:"department"`
		Count      int64  `json:"count"`
	}
	if err := s.db.Table("user_profiles").
		Select("department, COUNT(*) as count").
		Where("department IS NOT NULL AND department != ''").
		Group("department").
		Scan(&deptStats).Error; err != nil {
		log.Printf("获取部门统计失败: %v", err)
		return nil, err
	}

	for _, ds := range deptStats {
		stats.DepartmentStats[ds.Department] = ds.Count
	}

	log.Printf("用户统计信息获取成功 - 总用户数: %d", stats.TotalUsers)
	return stats, nil
}

// RecordLoginLog 记录登录日志
func (s *UserService) RecordLoginLog(userID uint, ipAddress, userAgent string) error {
	loginLog := models.LoginLog{
		UserID:    userID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}

	if err := s.db.Create(&loginLog).Error; err != nil {
		log.Printf("记录登录日志失败: %v", err)
		return err
	}

	return nil
}

// UpdateLastLogin 更新最后登录时间
func (s *UserService) UpdateLastLogin(userID uint) error {
	now := time.Now()
	if err := s.db.Model(&models.UserProfile{}).Where("user_id = ?", userID).Update("last_login", now).Error; err != nil {
		log.Printf("更新最后登录时间失败: %v", err)
		return err
	}

	return nil
}

// GetDB 获取数据库连接
func (s *UserService) GetDB() *gorm.DB {
	return s.db
}
