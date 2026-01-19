package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// User 用户基础信息表
type User struct {
	ID         uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Username   string    `gorm:"unique;not null;size:50" json:"username"`
	Password   string    `gorm:"not null;size:100" json:"-"` // 密码不返回前端
	Email      string    `gorm:"unique;not null;size:100" json:"email"`
	Status     string    `gorm:"type:enum('active','inactive');default:'active'" json:"status"`
	Department string    `gorm:"size:100;column:department" json:"department"`
	Title      string    `gorm:"size:50;column:title" json:"title"`
	Grade      string    `gorm:"size:20;column:grade" json:"grade"`
	Major      string    `gorm:"size:100;column:major" json:"major"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`

	// 关联关系
	Profile   *UserProfile `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"profile,omitempty"`
	Roles     []Role       `gorm:"many2many:user_roles;" json:"roles,omitempty"`
	LoginLogs []LoginLog   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"loginLogs,omitempty"`
}

func (u *User) TableName() string {
	return "users"
}

// UserProfile 用户扩展信息表
type UserProfile struct {
	UserID     uint       `gorm:"primaryKey;column:user_id" json:"userId"`
	RealName   string     `gorm:"size:50" json:"realName"`
	Phone      string     `gorm:"size:20" json:"phone"`
	Department string     `gorm:"size:100" json:"department"`
	StudentID  string     `gorm:"size:50" json:"studentId"`
	Avatar     string     `gorm:"size:255" json:"avatar"`
	Bio        string     `gorm:"type:text" json:"bio"`
	Interests  JSONArray  `gorm:"type:json" json:"interests"`
	LastLogin  *time.Time `gorm:"column:last_login" json:"lastLogin"`
}

func (up *UserProfile) TableName() string {
	return "user_profiles"
}

// Role 角色定义表
type Role struct {
	ID          uint   `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	RoleKey     string `gorm:"unique;not null;size:30;column:role_key" json:"roleKey"`
	RoleName    string `gorm:"not null;size:50;column:role_name" json:"roleName"`
	Description string `gorm:"type:text" json:"description"`

	// 关联关系
	Users []User `gorm:"many2many:user_roles;" json:"users,omitempty"`
}

func (r *Role) TableName() string {
	return "roles"
}

// UserRole 用户角色关联表（中间表）
type UserRole struct {
	UserID uint `gorm:"primaryKey;column:user_id" json:"userId"`
	RoleID uint `gorm:"primaryKey;column:role_id" json:"roleId"`
}

func (ur *UserRole) TableName() string {
	return "user_roles"
}

// LoginLog 登录记录表
type LoginLog struct {
	ID        uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	UserID    uint      `gorm:"not null;column:user_id" json:"userId"`
	LoginTime time.Time `gorm:"column:login_time;autoCreateTime" json:"loginTime"`
	IPAddress string    `gorm:"size:50;column:ip_address" json:"ipAddress"`
	UserAgent string    `gorm:"type:text;column:user_agent" json:"userAgent"`
}

func (ll *LoginLog) TableName() string {
	return "login_logs"
}

// JSONArray 用于处理JSON数组类型
type JSONArray []string

// Value 实现driver.Valuer接口
func (j JSONArray) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan 实现sql.Scanner接口
func (j *JSONArray) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, j)
	case string:
		return json.Unmarshal([]byte(v), j)
	default:
		return errors.New("cannot scan JSONArray")
	}
}

// UserWithRoles 用户信息（包含角色）
type UserWithRoles struct {
	User
	RoleNames []string `json:"roleNames"`
}

// UserCreateRequest 创建用户请求
type UserCreateRequest struct {
	Username   string   `json:"username" binding:"required,min=3,max=20"`
	Password   string   `json:"password" binding:"required,min=6,max=20"`
	Email      string   `json:"email" binding:"required,email"`
	RealName   string   `json:"realName" binding:"required"`
	Phone      string   `json:"phone"`
	Department string   `json:"department"`
	Title      string   `json:"title"`
	Grade      string   `json:"grade"`
	Major      string   `json:"major"`
	StudentID  string   `json:"studentId"`
	RoleKeys   []string `json:"roleKeys" binding:"required"`
}

// UserUpdateRequest 更新用户请求
type UserUpdateRequest struct {
	RealName   string   `json:"realName"`
	Email      string   `json:"email" binding:"omitempty,email"`
	Phone      string   `json:"phone"`
	Department string   `json:"department"`
	Title      string   `json:"title"`
	Grade      string   `json:"grade"`
	Major      string   `json:"major"`
	StudentID  string   `json:"studentId"`
	RoleKeys   []string `json:"roleKeys"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	ID         uint       `json:"id"`
	Username   string     `json:"username"`
	Email      string     `json:"email"`
	Status     string     `json:"status"`
	RealName   string     `json:"realName"`
	Phone      string     `json:"phone"`
	Department string     `json:"department"`
	Title      string     `json:"title"`
	Grade      string     `json:"grade"`
	Major      string     `json:"major"`
	StudentID  string     `json:"studentId"`
	RoleNames  []string   `json:"roleNames"`
	CreateTime time.Time  `json:"createTime"`
	CreatedAt  time.Time  `json:"createdAt"`
	LastLogin  *time.Time `json:"lastLogin"`
}

// UserDetailResponse 用户详情响应
type UserDetailResponse struct {
	ID         uint      `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Status     string    `json:"status"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
	Profile    struct {
		RealName   string     `json:"realName"`
		Phone      string     `json:"phone"`
		Department string     `json:"department"`
		StudentID  string     `json:"studentId"`
		Avatar     string     `json:"avatar"`
		Bio        string     `json:"bio"`
		Interests  JSONArray  `json:"interests"`
		LastLogin  *time.Time `json:"lastLogin"`
	} `json:"profile"`
	Roles []struct {
		ID          uint   `json:"id"`
		RoleKey     string `json:"roleKey"`
		RoleName    string `json:"roleName"`
		Description string `json:"description"`
	} `json:"roles"`
}

// UserStats 用户统计信息
type UserStats struct {
	TotalUsers      int64            `json:"totalUsers"`
	ActiveUsers     int64            `json:"activeUsers"`
	InactiveUsers   int64            `json:"inactiveUsers"`
	RoleStats       map[string]int64 `json:"roleStats"`
	DepartmentStats map[string]int64 `json:"departmentStats"`
}

// UserQueryParams 用户查询参数
type UserQueryParams struct {
	Page       int    `form:"page" binding:"omitempty,min=1"`
	Size       int    `form:"size" binding:"omitempty,min=1,max=100"`
	Search     string `form:"search"`
	Role       string `form:"role"`
	Status     string `form:"status"`
	Department string `form:"department"`
	SortBy     string `form:"sortBy"`
	SortOrder  string `form:"sortOrder"`
}
