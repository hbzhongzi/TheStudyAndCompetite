package models

import (
	"time"
)

// Project 项目表（严格映射 projects 表）
type Project struct {
	ID          uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Title       string `gorm:"column:title;type:varchar(255);not null" json:"title"`
	Description string `gorm:"column:description;varchar(255)" json:"description"`
	Type        string `gorm:"column:type;type:varchar(255)" json:"type"` // 冗余字段，便于展示
	StudentID   uint   `gorm:"column:student_id;not null" json:"studentId"`
	TeacherID   uint   `gorm:"column:teacher_id" json:"teacherId,omitempty"`

	Status          string     `gorm:"column:status;type:enum('draft','submitted','reviewing','approved','rejected','completed');default:'draft'" json:"status"`
	SubmittedAt     *time.Time `gorm:"column:submitted_at" json:"submittedAt"`
	ApprovedAt      *time.Time `gorm:"column:approved_at" json:"approvedAt"`
	ApprovedBy      *uint      `gorm:"column:approved_by" json:"approvedBy,omitempty"`
	RejectionReason string     `gorm:"column:rejection_reason;type:text" json:"rejectionReason"`

	Plan       string     `gorm:"column:plan;type:varchar(255)" json:"plan"`
	Progress   int        `gorm:"column:progress;default:0" json:"progress"`
	FinishTime *time.Time `gorm:"column:finish_time" json:"finishTime"`

	Deleted    bool `gorm:"column:deleted;default:0" json:"deleted"`
	IsApproved bool `gorm:"column:is_approved;default:0" json:"isApproved"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`

	// 关联关系
	Student *User `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Teacher *User `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
}

func (Project) TableName() string {
	return "projects"
}

// ProjectMember 项目成员表
type ProjectMember struct {
	ID            uint   `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	ProjectID     uint   `gorm:"not null;column:project_id" json:"projectId"`
	Name          string `gorm:"not null;size:50" json:"name"`
	StudentNumber string `gorm:"size:30;column:student_number" json:"studentNumber"`
	Role          string `gorm:"size:30" json:"role"`
}

func (pm *ProjectMember) TableName() string {
	return "project_members"
}

// ProjectFile 项目附件表
type ProjectFile struct {
	ID             uint       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	ProjectID      uint       `gorm:"not null;column:project_id" json:"projectId"`
	FileName       string     `gorm:"size:100;column:file_name" json:"fileName"`
	FileURL        string     `gorm:"size:255;column:file_url" json:"fileUrl"`
	FileType       string     `gorm:"size:50;column:file_type" json:"fileType"`
	FileVersion    string     `gorm:"size:20;column:file_version" json:"fileVersion"`
	ReviewStatus   string     `gorm:"size:20;default:'pending';column:review_status" json:"reviewStatus"`
	ReviewComments string     `gorm:"type:text;column:review_comments" json:"reviewComments"`
	ReviewedBy     *uint      `gorm:"column:reviewed_by" json:"reviewedBy"`
	ReviewedAt     *time.Time `gorm:"column:reviewed_at" json:"reviewedAt"`
	FileSize       int64      `gorm:"column:file_size" json:"fileSize"`
	DownloadCount  int        `gorm:"default:0;column:download_count" json:"downloadCount"`
	IsPublic       bool       `gorm:"default:false;column:is_public" json:"isPublic"`
	UploadTime     time.Time  `gorm:"column:upload_time;autoCreateTime" json:"uploadTime"`
}

func (pf *ProjectFile) TableName() string {
	return "project_files"
}

// ProjectReview 项目审核记录表
type ProjectReview struct {
	ID          uint       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	ProjectID   uint       `gorm:"not null;column:project_id" json:"projectId"`
	ReviewerID  uint       `gorm:"not null;column:reviewer_id" json:"reviewerId"`
	Status      string     `gorm:"type:enum('approved','rejected','archived');not null" json:"status"`
	Comments    string     `gorm:"type:text" json:"comments"`
	ReviewTime  time.Time  `gorm:"column:review_time;autoCreateTime" json:"reviewTime"`
	IsForce     bool       `gorm:"default:false;column:is_force" json:"isForce"`
	ReviewLevel int        `gorm:"column:review_level" json:"reviewLevel"`
	ReviewOrder int        `gorm:"column:review_order" json:"reviewOrder"`
	Deadline    *time.Time `gorm:"column:deadline" json:"deadline"`
	IsUrgent    bool       `gorm:"default:false;column:is_urgent" json:"isUrgent"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`

	// 关联关系
	Reviewer *User    `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
	Project  *Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

func (pr *ProjectReview) TableName() string {
	return "project_reviews"
}

// StudentTeacher 学生教师绑定关系表
type StudentTeacher struct {
	ID        uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	StudentID uint      `gorm:"not null;column:student_id" json:"studentId"`
	TeacherID uint      `gorm:"not null;column:teacher_id" json:"teacherId"`
	BindTime  time.Time `gorm:"column:bind_time;autoCreateTime" json:"bindTime"`

	// 关联关系
	Student *User `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Teacher *User `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
}

func (st *StudentTeacher) TableName() string {
	return "student_teacher"
}

// ProjectCreateRequest 创建项目请求
type ProjectCreateRequest struct {
	Title       string `json:"title" binding:"required,max=100"`
	Description string `json:"description"`
	Type        string `json:"type" `
	Status      string `json:"status" `
	TeacherID   uint   `json:"teacherId" binding:"required"`
	Plan        string `json:"plan"`
	FinishedAt  string `json:"finishedAt"`
}

// FileRequest 文件请求
type FileRequest struct {
	FileName string `json:"fileName" binding:"required,max=100"`
	FileURL  string `json:"fileUrl" binding:"required,max=255"`
}

// ProjectUpdateRequest 更新项目请求
type ProjectUpdateRequest struct {
	Title       string        `json:"title" binding:"max=100"`
	Description string        `json:"description"`
	Type        string        `json:"type" binding:"omitempty,oneof=科研 竞赛"`
	Status      string        `json:"status" binding:"omitempty,oneof=draft submitted approved rejected archived"`
	TeacherID   uint          `json:"teacherId" binding:"omitempty"`
	Level       string        `json:"level" binding:"omitempty,oneof=校级 省级 国家级"`
	CategoryID  *uint         `json:"categoryId"`
	Files       []FileRequest `json:"files"`
}

// ProjectReviewRequest 项目审核请求
type ProjectReviewRequest struct {
	Status   string `json:"status" binding:"required,oneof=approved rejected archived"`
	Comments string `json:"comments"`
	IsForce  bool   `json:"isForce"`
}

// ProjectListResponse 项目列表响应
type ProjectListResponse struct {
	ID          uint       `json:"id"`
	IsApproved  int        `json:"isApproved"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Type        string     `json:"type"`
	Status      string     `json:"status"`
	StudentName string     `json:"studentName"`
	StudentID   string     `json:"studentId"`
	TeacherName string     `json:"teacherName"`
	TeacherID   uint       `json:"teacherId"`
	SubmittedAt *time.Time `json:"submittedAt"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	MemberCount int        `json:"memberCount"`
	FileCount   int        `json:"fileCount"`
	ReviewCount int        `json:"reviewCount"`
}

// ProjectDetailResponse 项目详情响应
type ProjectDetailResponse struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Type        string     `json:"type"`
	Status      string     `json:"status"`
	SubmittedAt *time.Time `json:"submittedAt"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	Student     struct {
		ID         uint   `json:"id"`
		Username   string `json:"username"`
		RealName   string `json:"realName"`
		Email      string `json:"email"`
		Phone      string `json:"phone"`
		Department string `json:"department"`
		StudentID  string `json:"studentId"`
	} `json:"student"`
	Teacher struct {
		ID         uint   `json:"id"`
		Username   string `json:"username"`
		RealName   string `json:"realName"`
		Email      string `json:"email"`
		Phone      string `json:"phone"`
		Department string `json:"department"`
	} `json:"teacher"`
	Members []struct {
		ID            uint   `json:"id"`
		Name          string `json:"name"`
		StudentNumber string `json:"studentNumber"`
		Role          string `json:"role"`
	} `json:"members"`
	Files []struct {
		ID         uint      `json:"id"`
		FileName   string    `json:"fileName"`
		FileURL    string    `json:"fileUrl"`
		UploadTime time.Time `json:"uploadTime"`
	} `json:"files"`
	Reviews []struct {
		ID         uint      `json:"id"`
		Status     string    `json:"status"`
		Comments   string    `json:"comments"`
		ReviewTime time.Time `json:"reviewTime"`
		Reviewer   struct {
			ID       uint   `json:"id"`
			Username string `json:"username"`
			RealName string `json:"realName"`
		} `json:"reviewer"`
	} `json:"reviews"`
}

// ProjectStats 项目统计信息
type ProjectStats struct {
	TotalProjects       int64            `json:"totalProjects"`
	DraftProjects       int64            `json:"draftProjects"`
	PendingProjects     int64            `json:"pendingProjects"`
	ApprovedProjects    int64            `json:"approvedProjects"`
	RejectedProjects    int64            `json:"rejectedProjects"`
	ResearchProjects    int64            `json:"researchProjects"`
	CompetitionProjects int64            `json:"competitionProjects"`
	TypeStats           map[string]int64 `json:"typeStats"`
	StatusStats         map[string]int64 `json:"statusStats"`
	MonthlyStats        map[string]int64 `json:"monthlyStats"`
}

// ProjectQueryParams 项目查询参数
type ProjectQueryParams struct {
	Page       int    `form:"page" binding:"omitempty,min=1"`
	Size       int    `form:"size" binding:"omitempty,min=1,max=100"`
	Search     string `form:"search"`
	Type       string `form:"type"`
	Status     string `form:"status"`
	Level      string `form:"level"`
	CategoryID *uint  `form:"categoryId"`
	StudentID  uint   `form:"studentId"`
	SortBy     string `form:"sortBy"`
	SortOrder  string `form:"sortOrder"`
}

// ProjectExportRequest 项目导出请求
type ProjectExportRequest struct {
	Format  string             `json:"format" binding:"required,oneof=excel csv"`
	Filters ProjectQueryParams `json:"filters"`
}

// ProjectCreateResponse 创建项目响应
type ProjectCreateResponse struct {
	ProjectID uint `json:"projectId"`
}

// ProjectMyListResponse 我的项目列表响应
type ProjectMyListResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Plan        string    `json:"plan"`
	Progress    int       `json:"progress"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	Deadline    time.Time `json:"updated_at"`
}

// ProjectListForTeacherResponse 教师查看的项目列表响应
type ProjectListForTeacherResponse struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Type    string `json:"type"`
	Status  string `json:"status"`
	Student struct {
		Name      string `json:"name"`
		StudentID string `json:"studentId"`
	} `json:"student"`
}

// ProjectReviewResponse 项目审核响应
type ProjectReviewResponse struct {
	Reviewer   string    `json:"reviewer"`
	ReviewTime time.Time `json:"reviewTime"`
}

// ProjectReviewRecordResponse 审核记录响应
type ProjectReviewRecordResponse struct {
	Reviewer   string    `json:"reviewer"`
	Status     string    `json:"status"`
	Comments   string    `json:"comments"`
	ReviewTime time.Time `json:"reviewTime"`
}

// FileUploadResponse 文件上传响应
type FileUploadResponse struct {
	FileName string `json:"fileName"`
	FileURL  string `json:"fileUrl"`
}

// StudentTeacherBindRequest 学生教师绑定请求
type StudentTeacherBindRequest struct {
	StudentID uint `json:"studentId" binding:"required"`
	TeacherID uint `json:"teacherId" binding:"required"`
}

// StudentTeacherBindResponse 学生教师绑定响应
type StudentTeacherBindResponse struct {
	ID        uint      `json:"id"`
	StudentID uint      `json:"studentId"`
	TeacherID uint      `json:"teacherId"`
	BindTime  time.Time `json:"bindTime"`
	Student   struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		RealName string `json:"realName"`
	} `json:"student"`
	Teacher struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		RealName string `json:"realName"`
	} `json:"teacher"`
}

// TeacherQueryParams 教师查询参数
type TeacherQueryParams struct {
	Page       int    `form:"page" binding:"omitempty,min=1"`
	Size       int    `form:"size" binding:"omitempty,min=1,max=100"`
	Department string `form:"department"`
	Title      string `form:"title"`
	Status     string `form:"status"`
	SortBy     string `form:"sortBy"`
	SortOrder  string `form:"sortOrder"`
}

// TeacherListResponse 教师列表响应
type TeacherListResponse struct {
	ID         uint      `json:"id"`
	Username   string    `json:"username"`
	RealName   string    `json:"realName"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Department string    `json:"department"`
	Title      string    `json:"title"`
	Status     string    `json:"status"`
	Bio        string    `json:"bio"`
	CreatedAt  time.Time `json:"createdAt"`
}

// StudentQueryParams 学生查询参数
type StudentQueryParams struct {
	Page      int    `form:"page" binding:"omitempty,min=1"`
	Size      int    `form:"size" binding:"omitempty,min=1,max=100"`
	Status    string `form:"status"`
	Grade     string `form:"grade"`
	Major     string `form:"major"`
	SortBy    string `form:"sortBy"`
	SortOrder string `form:"sortOrder"`
}

// StudentListResponse 学生列表响应
type StudentListResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Grade     string    `json:"grade"`
	Major     string    `json:"major"`
	Status    string    `json:"status"`
	RealName  string    `json:"realName"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"createdAt"`
}

// TeacherProjectQueryParams 教师指导项目查询参数
type TeacherProjectQueryParams struct {
	Page       int    `form:"page" binding:"omitempty,min=1"`
	Size       int    `form:"size" binding:"omitempty,min=1,max=100"`
	Search     string `form:"search"`
	Type       string `form:"type"`
	Status     string `form:"status"`
	Level      string `form:"level"`
	CategoryID *uint  `form:"categoryId"`
	StudentID  uint   `form:"studentId"`
	SortBy     string `form:"sortBy"`
	SortOrder  string `form:"sortOrder"`
}

// TeacherProjectResponse 教师指导项目响应
type TeacherProjectResponse struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Type        string     `json:"type"`
	Status      string     `json:"status"`
	SubmittedAt *time.Time `json:"submittedAt"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	Student     struct {
		ID         uint   `json:"id"`
		Username   string `json:"username"`
		RealName   string `json:"realName"`
		Email      string `json:"email"`
		Phone      string `json:"phone"`
		Department string `json:"department"`
		StudentID  string `json:"studentId"`
	} `json:"student"`
	MemberCount int `json:"memberCount"`
	FileCount   int `json:"fileCount"`
	ReviewCount int `json:"reviewCount"`
}

// StudentBindTeacherRequest 学生绑定教师请求
type StudentBindTeacherRequest struct {
	TeacherID uint `json:"teacherId" binding:"required"`
}

// StudentBindTeacherResponse 学生绑定教师响应
type StudentBindTeacherResponse struct {
	ID        uint      `json:"id"`
	StudentID uint      `json:"studentId"`
	TeacherID uint      `json:"teacherId"`
	BindTime  time.Time `json:"bindTime"`
	Teacher   struct {
		ID         uint   `json:"id"`
		Username   string `json:"username"`
		RealName   string `json:"realName"`
		Email      string `json:"email"`
		Phone      string `json:"phone"`
		Department string `json:"department"`
		Bio        string `json:"bio"`
	} `json:"teacher"`
}

// ProjectType 项目分类表
type ProjectType struct {
	ID           uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Name         string    `gorm:"not null;size:100" json:"name"`
	Description  string    `gorm:"type:text" json:"description"`
	ParentID     *uint     `gorm:"column:parent_id" json:"parentId"`
	Level        int       `gorm:"default:1" json:"level"`
	SortOrder    int       `gorm:"default:0;column:sort_order" json:"sortOrder"`
	IsActive     bool      `gorm:"default:true;column:is_active" json:"isActive"`
	Icon         string    `gorm:"size:100" json:"icon"`
	Color        string    `gorm:"size:20" json:"color"`
	ProjectCount int64     `gorm:"-" json:"projectCount"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`

	// 关联关系
	Projects []Project     `gorm:"foreignKey:CategoryID" json:"projects,omitempty"`
	Parent   *ProjectType  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []ProjectType `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

func (pt *ProjectType) TableName() string {
	return "project_types"
}

// ProjectTypeCreateRequest 创建项目分类请求
type ProjectTypeCreateRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description"`
	ParentID    *uint  `json:"parentId"`
	Level       int    `json:"level" binding:"omitempty,min=1"`
	SortOrder   int    `json:"sortOrder" binding:"omitempty,min=0"`
	IsActive    bool   `json:"isActive"`
	Icon        string `json:"icon" binding:"omitempty,max=100"`
	Color       string `json:"color" binding:"omitempty,max=20"`
}

// ProjectTypeUpdateRequest 更新项目分类请求
type ProjectTypeUpdateRequest struct {
	Name        string `json:"name" binding:"omitempty,max=100"`
	Description string `json:"description"`
	ParentID    *uint  `json:"parentId"`
	Level       int    `json:"level" binding:"omitempty,min=1"`
	SortOrder   int    `json:"sortOrder" binding:"omitempty,min=0"`
	IsActive    *bool  `json:"isActive"`
	Icon        string `json:"icon" binding:"omitempty,max=100"`
	Color       string `json:"color" binding:"omitempty,max=20"`
}

// ProjectTypeResponse 项目分类响应
type ProjectTypeResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	ProjectCount int64     `json:"projectCount"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// =============================================
// 1. 项目状态管理增强 - 新增模型
// =============================================

// ProjectStatusUpdateRequest 项目状态更新请求
type ProjectStatusUpdateRequest struct {
	Status             string `json:"status" binding:"required,oneof=draft submitted approved rejected archived in_progress completed suspended need_revision"`
	StatusChangeReason string `json:"statusChangeReason" binding:"required"`
}

// ProjectStatusHistoryResponse 项目状态变更历史响应
type ProjectStatusHistoryResponse struct {
	ID           uint      `json:"id"`
	ProjectID    uint      `json:"projectId"`
	OldStatus    string    `json:"oldStatus"`
	NewStatus    string    `json:"newStatus"`
	ChangeReason string    `json:"changeReason"`
	ChangedBy    uint      `json:"changedBy"`
	ChangedAt    time.Time `json:"changedAt"`
	OperatorName string    `json:"operatorName"`
}

// =============================================
// 2. 项目生命周期管理增强 - 新增模型
// =============================================

// ProjectMilestoneCreateRequest 创建项目里程碑请求
type ProjectMilestoneCreateRequest struct {
	Title       string    `json:"title" binding:"required,max=200"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"dueDate" binding:"required"`
}

// ProjectMilestoneUpdateRequest 更新项目里程碑请求
type ProjectMilestoneUpdateRequest struct {
	Title       string     `json:"title" binding:"omitempty,max=200"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"dueDate"`
	Progress    *int       `json:"progress" binding:"omitempty,min=0,max=100"`
}

// ProjectMilestoneResponse 项目里程碑响应
type ProjectMilestoneResponse struct {
	ID            uint       `json:"id"`
	ProjectID     uint       `json:"projectId"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	DueDate       time.Time  `json:"dueDate"`
	CompletedDate *time.Time `json:"completedDate"`
	Status        string     `json:"status"`
	Progress      int        `json:"progress"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

// ProjectExtensionRequest 项目延期申请请求
type ProjectExtensionRequest struct {
	Reason           string    `json:"reason" binding:"required"`
	RequestedEndDate time.Time `json:"requestedEndDate" binding:"required"`
}

// ProjectExtensionReviewRequest 项目延期申请审核请求
type ProjectExtensionReviewRequest struct {
	Status         string `json:"status" binding:"required,oneof=approved rejected"`
	ReviewComments string `json:"reviewComments"`
}

// ProjectExtensionResponse 项目延期申请响应
type ProjectExtensionResponse struct {
	ID               uint       `json:"id"`
	ProjectID        uint       `json:"projectId"`
	ApplicantID      uint       `json:"applicantId"`
	Reason           string     `json:"reason"`
	OriginalEndDate  time.Time  `json:"originalEndDate"`
	RequestedEndDate time.Time  `json:"requestedEndDate"`
	Status           string     `json:"status"`
	ReviewerID       *uint      `json:"reviewerId"`
	ReviewComments   string     `json:"reviewComments"`
	ReviewedAt       *time.Time `json:"reviewedAt"`
	CreatedAt        time.Time  `json:"createdAt"`
}

// ProjectProgressUpdateRequest 项目进度更新请求
type ProjectProgressUpdateRequest struct {
	Progress int `json:"progress" binding:"required,min=0,max=100"`
}

// =============================================
// 3. 成果文件管理增强 - 新增模型
// =============================================

// ProjectFileUploadRequest 项目文件上传请求（增强版）
type ProjectFileUploadRequest struct {
	FileName    string `json:"fileName" binding:"required,max=100"`
	FileURL     string `json:"fileUrl" binding:"required,max=255"`
	FileType    string `json:"fileType" binding:"required,oneof=proposal midterm final achievement other"`
	FileVersion string `json:"fileVersion" binding:"omitempty,max=20"`
	IsPublic    bool   `json:"isPublic"`
}

// ProjectFileReviewRequest 项目文件审核请求
type ProjectFileReviewRequest struct {
	ReviewStatus   string `json:"reviewStatus" binding:"required,oneof=approved rejected"`
	ReviewComments string `json:"reviewComments"`
}

// ProjectFileEnhancedResponse 项目文件增强响应
type ProjectFileEnhancedResponse struct {
	ID             uint       `json:"id"`
	ProjectID      uint       `json:"projectId"`
	FileName       string     `json:"fileName"`
	FileURL        string     `json:"fileUrl"`
	FileType       string     `json:"fileType"`
	FileVersion    string     `json:"fileVersion"`
	ReviewStatus   string     `json:"reviewStatus"`
	ReviewComments string     `json:"reviewComments"`
	ReviewedBy     *uint      `json:"reviewedBy"`
	ReviewedAt     *time.Time `json:"reviewedAt"`
	FileSize       int64      `json:"fileSize"`
	DownloadCount  int        `json:"downloadCount"`
	IsPublic       bool       `json:"isPublic"`
	UploadTime     time.Time  `json:"uploadTime"`
}

// FileTypeConfigResponse 文件类型配置响应
type FileTypeConfigResponse struct {
	ID                uint   `json:"id"`
	FileType          string `json:"fileType"`
	DisplayName       string `json:"displayName"`
	Description       string `json:"description"`
	IsRequired        bool   `json:"isRequired"`
	MaxFileSize       int64  `json:"maxFileSize"`
	AllowedExtensions string `json:"allowedExtensions"`
	SortOrder         int    `json:"sortOrder"`
	IsActive          bool   `json:"isActive"`
}

// =============================================
// 4. 项目分类管理增强 - 新增模型
// =============================================

// ProjectTypeEnhancedResponse 项目分类增强响应
type ProjectTypeEnhancedResponse struct {
	ID           uint                          `json:"id"`
	Name         string                        `json:"name"`
	Description  string                        `json:"description"`
	ParentID     *uint                         `json:"parentId"`
	Level        int                           `json:"level"`
	SortOrder    int                           `json:"sortOrder"`
	IsActive     bool                          `json:"isActive"`
	Icon         string                        `json:"icon"`
	Color        string                        `json:"color"`
	ProjectCount int                           `json:"projectCount"`
	Children     []ProjectTypeEnhancedResponse `json:"children,omitempty"`
	CreatedAt    time.Time                     `json:"createdAt"`
	UpdatedAt    time.Time                     `json:"updatedAt"`
}

// ProjectTypeTreeResponse 项目分类树响应
type ProjectTypeTreeResponse struct {
	ID           uint                      `json:"id"`
	Name         string                    `json:"name"`
	Description  string                    `json:"description"`
	Level        int                       `json:"level"`
	SortOrder    int                       `json:"sortOrder"`
	Icon         string                    `json:"icon"`
	Color        string                    `json:"color"`
	ProjectCount int                       `json:"projectCount"`
	Children     []ProjectTypeTreeResponse `json:"children,omitempty"`
}

// ProjectTypeStatsResponse 项目分类统计响应
type ProjectTypeStatsResponse struct {
	TypeID             uint      `json:"typeId"`
	TypeName           string    `json:"typeName"`
	TotalProjects      int       `json:"totalProjects"`
	DraftProjects      int       `json:"draftProjects"`
	SubmittedProjects  int       `json:"submittedProjects"`
	ApprovedProjects   int       `json:"approvedProjects"`
	InProgressProjects int       `json:"inProgressProjects"`
	CompletedProjects  int       `json:"completedProjects"`
	RejectedProjects   int       `json:"rejectedProjects"`
	LastUpdated        time.Time `json:"lastUpdated"`
}

// =============================================
// 5. 审核流程增强 - 新增模型
// =============================================

// ReviewFlowCreateRequest 创建审核流程配置请求
type ReviewFlowCreateRequest struct {
	ProjectTypeID      *uint  `json:"projectTypeId"`
	ReviewLevel        int    `json:"reviewLevel" binding:"required,min=1,max=5"`
	ReviewerRole       string `json:"reviewerRole" binding:"required"`
	ReviewerDepartment string `json:"reviewerDepartment"`
	ReviewOrder        int    `json:"reviewOrder" binding:"required,min=1"`
	IsRequired         bool   `json:"isRequired"`
	DeadlineHours      int    `json:"deadlineHours" binding:"required,min=1"`
	AutoApprove        bool   `json:"autoApprove"`
	CanDelegate        bool   `json:"canDelegate"`
}

// ReviewFlowResponse 审核流程配置响应
type ReviewFlowResponse struct {
	ID                 uint      `json:"id"`
	ProjectTypeID      *uint     `json:"projectTypeId"`
	ReviewLevel        int       `json:"reviewLevel"`
	ReviewerRole       string    `json:"reviewerRole"`
	ReviewerDepartment string    `json:"reviewerDepartment"`
	ReviewOrder        int       `json:"reviewOrder"`
	IsRequired         bool      `json:"isRequired"`
	DeadlineHours      int       `json:"deadlineHours"`
	AutoApprove        bool      `json:"autoApprove"`
	CanDelegate        bool      `json:"canDelegate"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

// ReviewDelegationRequest 审核委托请求
type ReviewDelegationRequest struct {
	DelegatedReviewerID uint      `json:"delegatedReviewerId" binding:"required"`
	Reason              string    `json:"reason" binding:"required"`
	EndDate             time.Time `json:"endDate" binding:"required"`
}

// ReviewDelegationResponse 审核委托响应
type ReviewDelegationResponse struct {
	ID                  uint      `json:"id"`
	OriginalReviewerID  uint      `json:"originalReviewerId"`
	DelegatedReviewerID uint      `json:"delegatedReviewerId"`
	ProjectID           uint      `json:"projectId"`
	Reason              string    `json:"reason"`
	StartDate           time.Time `json:"startDate"`
	EndDate             time.Time `json:"endDate"`
	Status              string    `json:"status"`
	CreatedAt           time.Time `json:"createdAt"`
}

// ReviewTaskQueryParams 审核任务查询参数
type ReviewTaskQueryParams struct {
	Page      int    `form:"page" binding:"omitempty,min=1"`
	Size      int    `form:"size" binding:"omitempty,min=1,max=100"`
	Status    string `form:"status"`
	Priority  string `form:"priority"`
	SortBy    string `form:"sortBy"`
	SortOrder string `form:"sortOrder"`
}

// ReviewTaskResponse 审核任务响应
type ReviewTaskResponse struct {
	ID           uint       `json:"id"`
	ProjectID    uint       `json:"projectId"`
	ProjectTitle string     `json:"projectTitle"`
	ReviewLevel  int        `json:"reviewLevel"`
	ReviewOrder  int        `json:"reviewOrder"`
	Deadline     *time.Time `json:"deadline"`
	IsUrgent     bool       `json:"isUrgent"`
	Status       string     `json:"status"`
	StudentName  string     `json:"studentName"`
	ProjectType  string     `json:"projectType"`
	CreatedAt    time.Time  `json:"createdAt"`
}

// =============================================
// 通知系统相关模型
// =============================================

// ProjectNotificationResponse 项目通知响应
type ProjectNotificationResponse struct {
	ID        uint      `json:"id"`
	ProjectID uint      `json:"projectId"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	IsRead    bool      `json:"isRead"`
	Priority  string    `json:"priority"`
	CreatedAt time.Time `json:"createdAt"`
}

// NotificationTemplateResponse 通知模板响应
type NotificationTemplateResponse struct {
	ID              uint      `json:"id"`
	TemplateKey     string    `json:"templateKey"`
	TitleTemplate   string    `json:"titleTemplate"`
	ContentTemplate string    `json:"contentTemplate"`
	Variables       string    `json:"variables"`
	IsActive        bool      `json:"isActive"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// NotificationQueryParams 通知查询参数
type NotificationQueryParams struct {
	Page      int    `form:"page" binding:"omitempty,min=1"`
	Size      int    `form:"size" binding:"omitempty,min=1,max=100"`
	Type      string `form:"type"`
	IsRead    *bool  `form:"isRead"`
	Priority  string `form:"priority"`
	SortBy    string `form:"sortBy"`
	SortOrder string `form:"sortOrder"`
}

// NotificationTemplateUpdateRequest 通知模板更新请求
type NotificationTemplateUpdateRequest struct {
	TitleTemplate   string `json:"titleTemplate" binding:"required,max=200"`
	ContentTemplate string `json:"contentTemplate" binding:"required"`
	Variables       string `json:"variables"`
	IsActive        bool   `json:"isActive"`
}

// NotificationSendRequest 发送通知请求
type NotificationSendRequest struct {
	ProjectID uint   `json:"projectId" binding:"required"`
	UserIDs   []uint `json:"userIds" binding:"required"`
	Type      string `json:"type" binding:"required"`
	Title     string `json:"title" binding:"required,max=200"`
	Content   string `json:"content" binding:"required"`
	Priority  string `json:"priority" binding:"omitempty,oneof=low normal high urgent"`
}

// =============================================
// 通知相关模型
// =============================================

// ProjectNotification 项目通知表
type ProjectNotification struct {
	ID        uint       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	ProjectID uint       `gorm:"not null;column:project_id" json:"projectId"`
	UserID    uint       `gorm:"not null;column:user_id" json:"userId"`
	Type      string     `gorm:"size:50;not null" json:"type"`
	Title     string     `gorm:"size:200;not null" json:"title"`
	Content   string     `gorm:"type:text" json:"content"`
	Priority  string     `gorm:"size:20;default:'normal'" json:"priority"`
	IsRead    bool       `gorm:"default:false;column:is_read" json:"isRead"`
	ReadAt    *time.Time `gorm:"column:read_at" json:"readAt"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`

	// 关联关系
	Project *Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (pn *ProjectNotification) TableName() string {
	return "project_notifications"
}

// NotificationTemplate 通知模板表
type NotificationTemplate struct {
	ID              uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	TemplateKey     string    `gorm:"size:100;unique;not null;column:template_key" json:"templateKey"`
	TitleTemplate   string    `gorm:"size:200;not null;column:title_template" json:"titleTemplate"`
	ContentTemplate string    `gorm:"type:text;not null;column:content_template" json:"contentTemplate"`
	Variables       string    `gorm:"type:text;column:variables" json:"variables"`
	IsActive        bool      `gorm:"default:true;column:is_active" json:"isActive"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (nt *NotificationTemplate) TableName() string {
	return "notification_templates"
}

// =============================================
// 通知相关响应模型
// =============================================
// 项目状态管理增强模型
// =============================================

// ProjectStatusHistory 项目状态变更历史表
type ProjectStatusHistory struct {
	ID           uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	ProjectID    uint      `gorm:"not null;column:project_id" json:"projectId"`
	OldStatus    string    `gorm:"size:50;not null;column:old_status" json:"oldStatus"`
	NewStatus    string    `gorm:"size:50;not null;column:new_status" json:"newStatus"`
	ChangeReason string    `gorm:"type:text;not null;column:change_reason" json:"changeReason"`
	ChangedBy    uint      `gorm:"not null;column:changed_by" json:"changedBy"`
	ChangedAt    time.Time `gorm:"column:changed_at;autoCreateTime" json:"changedAt"`

	// 关联关系
	Project       *Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	ChangedByUser *User    `gorm:"foreignKey:ChangedBy" json:"changedByUser,omitempty"`
}

func (psh *ProjectStatusHistory) TableName() string {
	return "project_status_history"
}

// =============================================
// 项目生命周期管理增强模型
// =============================================

// ProjectMilestone 项目里程碑表
type ProjectMilestone struct {
	ID            uint       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	ProjectID     uint       `gorm:"not null;column:project_id" json:"projectId"`
	Title         string     `gorm:"size:200;not null" json:"title"`
	Description   string     `json:"description"`
	DueDate       time.Time  `gorm:"not null;column:due_date" json:"dueDate"`
	CompletedDate *time.Time `gorm:"column:completed_date" json:"completedDate"`
	Status        string     `gorm:"size:20;default:'pending'" json:"status"`
	Progress      int        `gorm:"default:0" json:"progress"`
	CreatedAt     time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`

	// 关联关系
	Project *Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

func (pm *ProjectMilestone) TableName() string {
	return "project_milestones"
}

// ProjectExtension 项目延期申请表
type ProjectExtension struct {
	ID               uint       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	ProjectID        uint       `gorm:"not null;column:project_id" json:"projectId"`
	ApplicantID      uint       `gorm:"not null;column:applicant_id" json:"applicantId"`
	Reason           string     `gorm:"type:text;not null" json:"reason"`
	OriginalEndDate  time.Time  `gorm:"not null;column:original_end_date" json:"originalEndDate"`
	RequestedEndDate time.Time  `gorm:"not null;column:requested_end_date" json:"requestedEndDate"`
	Status           string     `gorm:"size:20;default:'pending'" json:"status"`
	ReviewerID       *uint      `gorm:"column:reviewer_id" json:"reviewerId"`
	ReviewComments   string     `gorm:"type:text;column:review_comments" json:"reviewComments"`
	ReviewedAt       *time.Time `gorm:"column:reviewed_at" json:"reviewedAt"`
	CreatedAt        time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt        time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`

	// 关联关系
	Project   *Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Applicant *User    `gorm:"foreignKey:ApplicantID" json:"applicant,omitempty"`
	Reviewer  *User    `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
}

func (pe *ProjectExtension) TableName() string {
	return "project_extensions"
}

// =============================================
// 成果文件管理增强模型
// =============================================

// FileTypeConfig 文件类型配置表
type FileTypeConfig struct {
	ID                uint   `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	FileType          string `gorm:"size:50;unique;not null;column:file_type" json:"fileType"`
	DisplayName       string `gorm:"size:100;not null;column:display_name" json:"displayName"`
	Description       string `gorm:"type:text" json:"description"`
	IsRequired        bool   `gorm:"default:false;column:is_required" json:"isRequired"`
	MaxFileSize       int64  `gorm:"default:52428800;column:max_file_size" json:"maxFileSize"`
	AllowedExtensions string `gorm:"type:text;column:allowed_extensions" json:"allowedExtensions"`
	SortOrder         int    `gorm:"default:0;column:sort_order" json:"sortOrder"`
	IsActive          bool   `gorm:"default:true;column:is_active" json:"isActive"`
}

func (ftc *FileTypeConfig) TableName() string {
	return "file_type_configs"
}

// ProjectTypeStats 项目分类统计表
type ProjectTypeStats struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	TypeID             uint      `gorm:"not null;column:type_id" json:"typeId"`
	TypeName           string    `gorm:"size:100;not null;column:type_name" json:"typeName"`
	TotalProjects      int64     `gorm:"default:0;column:total_projects" json:"totalProjects"`
	DraftProjects      int64     `gorm:"default:0;column:draft_projects" json:"draftProjects"`
	SubmittedProjects  int64     `gorm:"default:0;column:submitted_projects" json:"submittedProjects"`
	ApprovedProjects   int64     `gorm:"default:0;column:approved_projects" json:"approvedProjects"`
	InProgressProjects int64     `gorm:"default:0;column:in_progress_projects" json:"inProgressProjects"`
	CompletedProjects  int64     `gorm:"default:0;column:completed_projects" json:"completedProjects"`
	RejectedProjects   int64     `gorm:"default:0;column:rejected_projects" json:"rejectedProjects"`
	LastUpdated        time.Time `gorm:"column:last_updated;autoUpdateTime" json:"lastUpdated"`
}

func (pts *ProjectTypeStats) TableName() string {
	return "project_type_stats"
}

// =============================================
// 审核流程增强模型
// =============================================

// ProjectReviewFlow 项目审核流程配置表
type ProjectReviewFlow struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	ProjectTypeID      uint      `gorm:"not null;column:project_type_id" json:"projectTypeId"`
	ReviewLevel        int       `gorm:"not null;column:review_level" json:"reviewLevel"`
	ReviewerRole       string    `gorm:"size:50;not null;column:reviewer_role" json:"reviewerRole"`
	ReviewerDepartment string    `gorm:"size:100;column:reviewer_department" json:"reviewerDepartment"`
	ReviewOrder        int       `gorm:"not null;column:review_order" json:"reviewOrder"`
	IsRequired         bool      `gorm:"default:true;column:is_required" json:"isRequired"`
	DeadlineHours      int       `gorm:"default:72;column:deadline_hours" json:"deadlineHours"`
	AutoApprove        bool      `gorm:"default:false;column:auto_approve" json:"autoApprove"`
	CanDelegate        bool      `gorm:"default:false;column:can_delegate" json:"canDelegate"`
	CreatedAt          time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt          time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (prf *ProjectReviewFlow) TableName() string {
	return "project_review_flows"
}

// ReviewDelegation 审核委托表
type ReviewDelegation struct {
	ID                  uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	OriginalReviewerID  uint      `gorm:"not null;column:original_reviewer_id" json:"originalReviewerId"`
	DelegatedReviewerID uint      `gorm:"not null;column:delegated_reviewer_id" json:"delegatedReviewerId"`
	ProjectID           uint      `gorm:"not null;column:project_id" json:"projectId"`
	Reason              string    `gorm:"type:text;not null" json:"reason"`
	StartDate           time.Time `gorm:"not null;column:start_date" json:"startDate"`
	EndDate             time.Time `gorm:"not null;column:end_date" json:"endDate"`
	Status              string    `gorm:"size:20;default:'active'" json:"status"`
	CreatedAt           time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt           time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`

	// 关联关系
	OriginalReviewer  *User    `gorm:"foreignKey:OriginalReviewerID" json:"originalReviewer,omitempty"`
	DelegatedReviewer *User    `gorm:"foreignKey:DelegatedReviewerID" json:"delegatedReviewer,omitempty"`
	Project           *Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

func (rd *ReviewDelegation) TableName() string {
	return "review_delegations"
}
