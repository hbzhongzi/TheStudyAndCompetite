package services

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"yunmeng-backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectService struct {
	db *gorm.DB
}

func NewProjectService(db *gorm.DB) *ProjectService {
	return &ProjectService{db: db}
}

// GetProjectList 获取项目列表
func (s *ProjectService) GetProjectList(params models.ProjectQueryParams) ([]models.ProjectListResponse, int64, error) {
	var projects []models.Project
	var total int64

	// 构建查询
	query := s.db.Model(&models.Project{}).
		Preload("Student.Profile").
		Preload("Teacher.Profile").
		Preload("Members").
		Preload("Files").
		Preload("Reviews")

	// 搜索条件
	if params.Search != "" {
		search := "%" + params.Search + "%"
		query = query.Where("projects.title LIKE ? OR projects.description LIKE ?", search, search)
	}

	// 类型筛�?
	if params.Type != "" {
		query = query.Where("projects.type = ?", params.Type)
	}

	// 状态筛�?
	if params.Status != "" {
		query = query.Where("projects.status = ?", params.Status)
	}

	// 学生ID筛�?
	if params.StudentID > 0 {
		query = query.Where("projects.student_id = ?", params.StudentID)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		log.Printf("获取项目总数失败: %v", err)
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
		query = query.Order("projects.created_at DESC")
	}

	// 分页
	if params.Page > 0 && params.Size > 0 {
		offset := (params.Page - 1) * params.Size
		query = query.Offset(offset).Limit(params.Size)
	}

	// 执行查询
	if err := query.Find(&projects).Error; err != nil {
		log.Printf("查询项目列表失败: %v", err)
		return nil, 0, err
	}

	// 转换为响应格�?
	var responses []models.ProjectListResponse
	for _, project := range projects {
		response := models.ProjectListResponse{
			ID:          project.ID,
			Title:       project.Title,
			Description: project.Description,
			Type:        project.Type,
			Status:      project.Status,
			TeacherID:   project.TeacherID,
			SubmittedAt: project.SubmittedAt,
			CreatedAt:   project.CreatedAt,
			UpdatedAt:   project.UpdatedAt,
		}

		// 添加学生信息
		if project.Student != nil && project.Student.Profile != nil {
			response.StudentName = project.Student.Profile.RealName
			response.StudentID = project.Student.Profile.StudentID
		}

		// 添加教师信息
		if project.Teacher != nil && project.Teacher.Profile != nil {
			response.TeacherName = project.Teacher.Profile.RealName
		}

		responses = append(responses, response)
	}

	log.Printf("项目列表查询完成 - 总数: %d, 返回数量: %d", total, len(responses))
	return responses, total, nil
}

// GetProjectByID 根据ID获取项目详情
func (s *ProjectService) GetProjectByID(id uint) (*models.ProjectDetailResponse, error) {
	var project models.Project
	err := s.db.Preload("Student.Profile").
		Preload("Teacher.Profile").
		First(&project, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("项目不存在")
		}
		log.Printf("查询项目详情失败 - 项目ID: %d, 错误: %v", id, err)
		return nil, err
	}

	// 构建响应数据
	response := &models.ProjectDetailResponse{
		ID:          project.ID,
		Title:       project.Title,
		Description: project.Description,
		IsApproved:  project.IsApproved,
		Plan:        project.Plan,
		Type:        project.Type,
		Status:      project.Status,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	}

	// 添加学生信息
	if project.Student != nil && project.Student.Profile != nil {
		response.Student.ID = project.Student.ID
		response.Student.Username = project.Student.Username
		response.Student.RealName = project.Student.Profile.RealName
		response.Student.Email = project.Student.Email
		response.Student.Phone = project.Student.Profile.Phone
		response.Student.Department = project.Student.Profile.Department
		response.Student.StudentID = project.Student.Profile.StudentID
	}

	// 添加教师信息
	if project.Teacher != nil && project.Teacher.Profile != nil {
		response.Teacher.ID = project.Teacher.ID
		response.Teacher.Username = project.Teacher.Username
		response.Teacher.RealName = project.Teacher.Profile.RealName
		response.Teacher.Email = project.Teacher.Email
		response.Teacher.Phone = project.Teacher.Profile.Phone
		response.Teacher.Department = project.Teacher.Profile.Department
	}

	log.Printf("项目详情查询完成 - 项目ID: %d, 标题: %s", id, project.Title)
	return response, nil
}

// CreateProject 创建项目
func (s *ProjectService) CreateProject(studentID uint, req models.ProjectCreateRequest) (*models.ProjectCreateResponse, error) {
	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建项目
	project := &models.Project{
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		StudentID:   studentID,
		TeacherID:   req.TeacherID,
		Plan:        req.Plan,
		FinishTime:  req.FinishedAt,
		Status:      "draft",
		Progress:    0,
		Deleted:     false,
		IsApproved:  false,
	}

	if err := tx.Create(&project).Error; err != nil {
		tx.Rollback()
		log.Printf("创建项目失败: %v", err)
		return nil, errors.New("创建项目失败")
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Printf("提交事务失败: %v", err)
		return nil, errors.New("创建项目失败")
	}

	log.Printf("项目创建成功 - 项目ID: %d, 标题: %s", project.ID, project.Title)
	return &models.ProjectCreateResponse{ProjectID: project.ID}, nil
}

// CreateExtensionApplication 创建延期申请
func (s *ProjectService) CreateExtensionApplication(
	studentID uint,
	req models.ExtensionApplicationRequest,
) (*models.ExtensionApplicationRequest, error) {

	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, errors.New("事务启动失败")
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1️⃣ 校验项目是否存在 & 是否属于该学生 & 是否已通过
	var project models.Project
	if err := tx.Where(
		"id = ? AND student_id = ? AND status = ? AND deleted = 0",
		req.ProjectID,
		studentID,
		"approved",
	).First(&project).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("项目不存在或未通过审批，无法申请延期")
	}

	// 2️⃣ 去重：是否已有待审批延期申请
	var count int64
	tx.Model(&models.ProjectExtensionApplication{}).
		Where("project_id = ? AND status = 'pending'", req.ProjectID).
		Count(&count)

	if count > 0 {
		tx.Rollback()
		return nil, errors.New("该项目已有待审批的延期申请")
	}

	// 3️⃣ 创建延期申请
	application := &models.ProjectExtensionApplication{
		ProjectID:           project.ID,
		StudentID:           studentID,
		TeacherID:           project.TeacherID,
		OriginalFinishTime:  project.FinishTime,
		RequestedFinishTime: req.RequestedFinishTime,
		ApplyReason:         req.ApplyReason,
		Status:              "pending",
	}

	if err := tx.Create(application).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("创建延期申请失败")
	}

	// 4️⃣ 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("提交事务失败")
	}

	return &models.ExtensionApplicationRequest{
		ProjectID:           application.ID,
		RequestedFinishTime: application.RequestedFinishTime,
		ApplyReason:         application.ApplyReason,
	}, nil
}

// GetStudentExtensionList 获取学生的延期申请列表
func (s *ProjectService) GetStudentExtensionList(
	studentID uint,
	query models.ExtensionListQuery,
) ([]models.ExtensionApplicationListResponse, int64, error) {

	db := s.db.Table("project_extension_applications pea").
		Joins("LEFT JOIN projects p ON pea.project_id = p.id").
		Joins("LEFT JOIN users u ON pea.student_id = u.id").
		Where("pea.student_id = ?", studentID)

	if query.Status != "" {
		db = db.Where("pea.status = ?", query.Status)
	}
	// 1️⃣ 统计总数（只 JOIN 一次）
	var total int64
	countDB := s.db.
		Table("project_extension_applications AS pea").
		Joins("LEFT JOIN projects p ON pea.project_id = p.id").
		Where("pea.student_id = ?", studentID)

	if err := countDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 2️⃣ 查询列表（全新 db 链，重新 JOIN）
	var list []models.ExtensionApplicationListResponse
	listDB := s.db.
		Table("project_extension_applications AS pea").
		Joins("LEFT JOIN projects p ON pea.project_id = p.id").
		Joins("LEFT JOIN users t ON pea.teacher_id = t.id").
		Where("pea.student_id = ?", studentID)

	err := listDB.
		Select(`
		pea.id,
		pea.project_id,
		p.title AS project_title,

		pea.student_id,
		NULL AS student_name,

		pea.teacher_id,
		t.username AS teacher_name,

		pea.original_finish_time,
		pea.requested_finish_time,

		pea.apply_reason,
		pea.status,
		pea.review_reason,
		pea.reviewed_at,
		pea.created_at
	`).
		Order("pea.created_at DESC").
		Offset((query.Page - 1) * query.Size).
		Limit(query.Size).
		Scan(&list).Error

	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

/* =========================
   项目状态校验
========================= */

func (s *ProjectService) CheckProjectEditable(projectID uint) error {
	var status string
	log.Printf("projectID: %d", projectID)
	err := s.db.
		Table("projects").
		Where("id = ?", projectID).
		Pluck("status", &status).Error

	if err != nil {
		return err
	}
	if status != "approved" {
		return errors.New("项目未通过审核，禁止上传文件")
	}
	return nil
}

/* =========================
   文件列表
========================= */

func (s *ProjectService) GetProjectFiles(
	projectID uint,
	userID uint,
) ([]models.File, error) {

	var files []models.File
	err := s.db.
		Where("project_id = ?  AND deleted_at IS NULL ", projectID).
		Order("created_at DESC").
		Find(&files).Error

	return files, err
}

/* =========================
   上传文件
========================= */

func (s *ProjectService) SaveProjectFiles(
	projectID uint,
	userID uint,
	files []*multipart.FileHeader,
) error {

	basePath := fmt.Sprintf("uploads/projects/%d", projectID)
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		return err
	}

	for _, f := range files {
		ext := path.Ext(f.Filename)
		newName := uuid.New().String() + ext
		fullPath := path.Join(basePath, newName)

		if err := saveUploadedFile(f, fullPath); err != nil {
			return err
		}

		record := models.File{
			FileName:     newName,
			OriginalName: f.Filename,
			FilePath:     fullPath,
			FileSize:     f.Size,
			Status:       "draft", // 默认状态为 draft
			FileExt:      ext,
			UploadedBy:   userID,
			ProjectID:    projectID,
			CreatedAt:    time.Now(),
		}

		if err := s.db.Create(&record).Error; err != nil {
			return err
		}
	}

	return nil
}

/* =========================
   删除文件（软删除）
========================= */

func (s *ProjectService) DeleteProjectFile(
	fileID uint,
	userID uint,
	projectID uint,
) error {

	// 1. 先检查文件是否存在且属于该项目
	var file models.File
	err := s.db.Where("id = ? AND project_id = ?", fileID, projectID).
		First(&file).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文件不存在或不属于该项目")
		}
		return err
	}

	// 2. 权限检查：文件上传者可以删除
	canDelete := false

	// 方式1：直接检查文件的上传者
	if file.UploadedBy == userID {
		canDelete = true
	}

	// 3. 执行软删除（设置 deleted_at）
	if !canDelete {
		return errors.New("无删除权限")
	}
	// 3. 执行软删除（设置 deleted_at）
	now := time.Now()
	result := s.db.Model(&models.File{}).
		Where("id = ?", fileID).
		Updates(map[string]interface{}{
			"deleted_at": &now, // 设置当前时间
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("删除失败")
	}

	return nil
}

/* =========================
   工具函数
========================= */

func saveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func detectCategory(ext string) string {
	switch strings.ToLower(ext) {
	case ".pdf", ".doc", ".docx":
		return "document"
	case ".jpg", ".png":
		return "image"
	case ".zip", ".rar":
		return "archive"
	case ".mp4":
		return "video"
	default:
		return "other"
	}
}

// UpdateProject 更新项目
func (s *ProjectService) UpdateProject(id uint, studentID uint, req models.ProjectUpdateRequest) error {
	// 检查项目是否存�?
	var project models.Project
	if err := s.db.First(&project, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("项目不存在")
		}
		return err
	}

	// 检查权限：只有项目创建者可以修�?
	if project.StudentID != studentID {
		return errors.New("无权限修改此项目")
	}

	// 检查状态：只有草稿或已驳回状态的项目可以修改
	if project.Status != "draft" && project.Status != "rejected" {
		return errors.New("只有草稿或已驳回状态的项目可以修改")
	}

	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新项目基本信息
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.TeacherID > 0 {
		updates["teacher_id"] = req.TeacherID
	}

	if len(updates) > 0 {
		if err := tx.Model(&project).Updates(updates).Error; err != nil {
			tx.Rollback()
			log.Printf("更新项目失败: %v", err)
			return errors.New("更新项目失败")
		}
	}

	// 更新项目文件
	if req.Files != nil {
		// 删除现有文件
		if err := tx.Where("project_id = ?", id).Delete(&models.ProjectFile{}).Error; err != nil {
			tx.Rollback()
			log.Printf("删除项目文件失败: %v", err)
			return errors.New("更新项目文件失败")
		}

		// 创建新文件
		for _, fileReq := range req.Files {
			file := models.ProjectFile{
				ProjectID: id,
				FileName:  fileReq.FileName,
				FileURL:   fileReq.FileURL,
			}

			if err := tx.Create(&file).Error; err != nil {
				tx.Rollback()
				log.Printf("创建项目文件失败: %v", err)
				return errors.New("更新项目文件失败")
			}
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Printf("提交事务失败: %v", err)
		return errors.New("更新项目失败")
	}

	log.Printf("项目更新成功 - 项目ID: %d", id)
	return nil
}

// DeleteProject 删除项目
func (s *ProjectService) DeleteProject(id uint) error {
	// 检查项目是否存�?
	var project models.Project
	if err := s.db.First(&project, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("项目不存在")
		}
		return err
	}

	// 检查状态：只有草稿或已驳回状态的项目可以删除
	//if project.Status != "draft" && project.Status != "rejected" {
	//	return errors.New("只有草稿或已驳回状态的项目可以删除")
	//}

	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除项目（关联数据会通过CASCADE自动删除）
	if err := tx.Delete(&project).Error; err != nil {
		tx.Rollback()
		log.Printf("删除项目失败: %v", err)
		return errors.New("删除项目失败")
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Printf("提交事务失败: %v", err)
		return errors.New("删除项目失败")
	}

	log.Printf("项目删除成功 - 项目ID: %d", id)
	return nil
}

// ReviewProject 审核项目
func (s *ProjectService) ReviewProject(projectID uint, reviewerID uint, req models.ProjectReviewRequest) error {
	// 检查项目是否存�?
	var project models.Project
	if err := s.db.First(&project, projectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("项目不存在")
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

	// 创建审核记录
	review := models.ProjectReview{
		ProjectID:  projectID,
		ReviewerID: reviewerID,
		Status:     req.Status,
		Comments:   req.Comments,
	}

	if err := tx.Create(&review).Error; err != nil {
		tx.Rollback()
		log.Printf("创建审核记录失败: %v", err)
		return errors.New("创建审核记录失败")
	}

	// 更新项目状�?
	if err := tx.Model(&project).Update("status", req.Status).Error; err != nil {
		tx.Rollback()
		log.Printf("更新项目状态失败: %v", err)
		return errors.New("更新项目状态失败")
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Printf("提交事务失败: %v", err)
		return errors.New("审核项目失败")
	}

	log.Printf("项目审核完成 - 项目ID: %d, 审核结果: %s", projectID, req.Status)
	return nil
}

// GetMyProjects 获取我的项目列表
func (s *ProjectService) GetMyProjects(studentID uint, status string, page, size int) ([]models.ProjectMyListResponse, int64, error) {
	var projects []models.Project
	query := s.db.Where("student_id = ?", studentID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	var total int64
	query.Model(&models.Project{}).Count(&total)

	// 分页查询
	offset := (page - 1) * size
	err := query.Offset(offset).Limit(size).Order("created_at DESC").Find(&projects).Error
	if err != nil {
		log.Printf("获取我的项目列表失败: %v", err)
		return nil, 0, err
	}

	var responses []models.ProjectMyListResponse
	for _, project := range projects {
		responses = append(responses, models.ProjectMyListResponse{
			ID:          project.ID,
			Title:       project.Title,
			Type:        project.Type,
			Status:      project.Status,
			CreatedAt:   project.CreatedAt,
			Deadline:    project.UpdatedAt,
			Description: project.Description,
			Progress:    project.Progress,
			Plan:        project.Plan,
		})
	}

	return responses, total, nil
}

// ProjectStats 项目统计返回结构
type ProjectStats struct {
	TotalProjects      int64 `json:"totalProjects"`
	OngoingProjects    int64 `json:"ongoingProjects"`
	CompletedProjects  int64 `json:"completedProjects"`
	PendingProjects    int64 `json:"pendingProjects"`
	TotalCompetitions  int64 `json:"totalCompetitions"`
	ActiveCompetitions int64 `json:"activeCompetitions"`
}

// GetProjectStats 获取学生项目与竞赛统计信息
func (s *ProjectService) GetProjectStats(userID uint) (*ProjectStats, error) {
	stats := &ProjectStats{}

	// 项目总数
	if err := s.db.
		Table("projects").
		Where("student_id = ?", userID).
		Count(&stats.TotalProjects).Error; err != nil {
		return nil, err
	}

	// 进行中项目
	s.db.Table("projects").
		Where("student_id = ? AND status = ?", userID, "ongoing").
		Count(&stats.OngoingProjects)

	// 已完成项目
	s.db.Table("projects").
		Where("student_id = ? AND status = ?", userID, "completed").
		Count(&stats.CompletedProjects)

	// 待审核项目
	s.db.Table("projects").
		Where("student_id = ? AND status = ?", userID, "pending").
		Count(&stats.PendingProjects)

	return stats, nil
}

// GetProjectsForTeacher 获取教师查看的项目列�?
func (s *ProjectService) GetProjectsForTeacher(params models.ProjectQueryParams) ([]models.ProjectListForTeacherResponse, int64, error) {
	var projects []models.Project
	var total int64

	query := s.db.Model(&models.Project{}).Preload("Student.Profile")

	// 状态筛�?
	if params.Status != "" {
		query = query.Where("projects.status = ?", params.Status)
	}

	// 类型筛�?
	if params.Type != "" {
		query = query.Where("projects.type = ?", params.Type)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	if params.Page > 0 && params.Size > 0 {
		offset := (params.Page - 1) * params.Size
		query = query.Offset(offset).Limit(params.Size)
	}

	// 执行查询
	if err := query.Order("projects.created_at DESC").Find(&projects).Error; err != nil {
		return nil, 0, err
	}

	var responses []models.ProjectListForTeacherResponse
	for _, project := range projects {
		response := models.ProjectListForTeacherResponse{
			ID:     project.ID,
			Title:  project.Title,
			Type:   project.Type,
			Status: project.Status,
		}

		if project.Student != nil && project.Student.Profile != nil {
			response.Student.Name = project.Student.Profile.RealName
			response.Student.StudentID = project.Student.Profile.StudentID
		}

		responses = append(responses, response)
	}

	return responses, total, nil
}

// ReviewProjectWithResponse 审核项目并返回响�?
func (s *ProjectService) ReviewProjectWithResponse(projectID uint, reviewerID uint, req models.ProjectReviewRequest) (*models.ProjectReviewResponse, error) {
	// 检查项目是否存�?
	var project models.Project
	if err := s.db.First(&project, projectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("项目不存在")
		}
		return nil, err
	}

	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建审核记录
	review := models.ProjectReview{
		ProjectID:  projectID,
		ReviewerID: reviewerID,
		Status:     req.Status,
		Comments:   req.Comments,
	}

	if err := tx.Create(&review).Error; err != nil {
		tx.Rollback()
		log.Printf("创建审核记录失败: %v", err)
		return nil, errors.New("创建审核记录失败")
	}

	// 更新项目状�?
	if err := tx.Model(&project).Update("status", req.Status).Error; err != nil {
		tx.Rollback()
		log.Printf("更新项目状态失败: %v", err)
		return nil, errors.New("更新项目状态失败")
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Printf("提交事务失败: %v", err)
		return nil, errors.New("审核项目失败")
	}

	// 获取审核者信�?
	var reviewer models.User
	if err := s.db.Preload("Profile").First(&reviewer, reviewerID).Error; err != nil {
		log.Printf("获取审核者信息失�? %v", err)
	}

	reviewerName := reviewer.Username
	if reviewer.Profile != nil {
		reviewerName = reviewer.Profile.RealName
	}

	return &models.ProjectReviewResponse{
		Reviewer:   reviewerName,
		ReviewTime: review.ReviewTime,
	}, nil
}

// GetProjectReviews 获取项目审核记录
func (s *ProjectService) GetProjectReviews(projectID uint) ([]models.ProjectReviewRecordResponse, error) {
	var reviews []models.ProjectReview
	err := s.db.Preload("Reviewer.Profile").Where("project_id = ?", projectID).Order("review_time DESC").Find(&reviews).Error
	if err != nil {
		return nil, err
	}

	var responses []models.ProjectReviewRecordResponse
	for _, review := range reviews {
		reviewerName := review.Reviewer.Username
		if review.Reviewer.Profile != nil {
			reviewerName = review.Reviewer.Profile.RealName
		}

		responses = append(responses, models.ProjectReviewRecordResponse{
			Reviewer:   reviewerName,
			Status:     review.Status,
			Comments:   review.Comments,
			ReviewTime: review.ReviewTime,
		})
	}

	return responses, nil
}

// GetDB 获取数据库实�?
// BindStudentTeacher 绑定学生和教�?
func (s *ProjectService) BindStudentTeacher(req models.StudentTeacherBindRequest) (*models.StudentTeacherBindResponse, error) {
	// 检查学生是否存�?
	var student models.User
	if err := s.db.Preload("Profile").First(&student, req.StudentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("学生不存在")
		}
		return nil, err
	}

	// 检查教师是否存在
	var teacher models.User
	if err := s.db.Preload("Profile").First(&teacher, req.TeacherID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("教师不存在")
		}
		return nil, err
	}

	// 检查是否已经绑定
	var existingBind models.StudentTeacher
	if err := s.db.Where("student_id = ? AND teacher_id = ?", req.StudentID, req.TeacherID).First(&existingBind).Error; err == nil {
		return nil, errors.New("该学生和教师已经绑定")
	}

	// 创建绑定关系
	bind := models.StudentTeacher{
		StudentID: req.StudentID,
		TeacherID: req.TeacherID,
	}

	if err := s.db.Create(&bind).Error; err != nil {
		log.Printf("创建学生教师绑定失败: %v", err)
		return nil, errors.New("创建绑定关系失败")
	}

	// 构建响应
	response := &models.StudentTeacherBindResponse{
		ID:        bind.ID,
		StudentID: bind.StudentID,
		TeacherID: bind.TeacherID,
		BindTime:  bind.BindTime,
	}

	// 添加学生信息
	if student.Profile != nil {
		response.Student.ID = student.ID
		response.Student.Username = student.Username
		response.Student.RealName = student.Profile.RealName
	}

	// 添加教师信息
	if teacher.Profile != nil {
		response.Teacher.ID = teacher.ID
		response.Teacher.Username = teacher.Username
		response.Teacher.RealName = teacher.Profile.RealName
	}

	log.Printf("学生教师绑定成功 - 学生ID: %d, 教师ID: %d", req.StudentID, req.TeacherID)
	return response, nil
}

// GetTeacherList 获取教师列表
func (s *ProjectService) GetTeacherList() ([]models.TeacherListResponse, error) {
	var teachers []models.User
	err := s.db.Preload("Profile").
		Joins("JOIN user_roles ur ON users.id = ur.user_id").
		Joins("JOIN roles r ON ur.role_id = r.id").
		Where("r.role_key = ?", "teacher").
		Where("users.status = ?", "active").
		Find(&teachers).Error

	if err != nil {
		log.Printf("获取教师列表失败: %v", err)
		return nil, err
	}

	var responses []models.TeacherListResponse
	for _, teacher := range teachers {
		response := models.TeacherListResponse{
			ID:       teacher.ID,
			Username: teacher.Username,
			Email:    teacher.Email,
		}

		if teacher.Profile != nil {
			response.RealName = teacher.Profile.RealName
			response.Phone = teacher.Profile.Phone
			response.Department = teacher.Profile.Department
			response.Bio = teacher.Profile.Bio
		}

		responses = append(responses, response)
	}

	return responses, nil
}

// ApproveExtensionApplication 教师审批延期申请
func (s *ProjectService) ApproveExtensionApplication(
	teacherID uint,
	req models.ExtensionApprovalRequest,
) error {

	tx := s.db.Begin()
	if tx.Error != nil {
		return errors.New("事务启动失败")
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1️⃣ 查询延期申请
	var application models.ProjectExtensionApplication
	if err := tx.Where(
		"id = ? AND status = 'pending'",
		req.ApplicationID,
	).First(&application).Error; err != nil {
		tx.Rollback()
		return errors.New("延期申请不存在或已处理")
	}

	// 2️⃣ 校验教师权限
	if application.TeacherID != teacherID {
		tx.Rollback()
		return errors.New("无权审批该延期申请")
	}

	now := time.Now()

	// 3️⃣ 根据审批动作处理
	switch req.Action {

	case "approved":
		// 3.1 更新延期申请
		if err := tx.Model(&application).Updates(map[string]interface{}{
			"status":      "approved",
			"reviewed_at": now,
		}).Error; err != nil {
			tx.Rollback()
			return errors.New("审批延期申请失败")
		}

		// 3.2 同步更新项目完成时间
		if err := tx.Model(&models.Project{}).
			Where("id = ?", application.ProjectID).
			Update("finish_time", application.RequestedFinishTime).
			Error; err != nil {
			tx.Rollback()
			return errors.New("更新项目完成时间失败")
		}

	case "rejected":
		// 3.3 驳回延期申请
		if err := tx.Model(&application).Updates(map[string]interface{}{
			"status":         "rejected",
			"reviewed_at":    now,
			"reviewed_by":    teacherID,
			"review_comment": req.Reason,
		}).Error; err != nil {
			tx.Rollback()
			return errors.New("驳回延期申请失败")
		}

	default:
		tx.Rollback()
		return errors.New("非法审批操作")
	}

	// 4️⃣ 提交事务
	if err := tx.Commit().Error; err != nil {
		return errors.New("事务提交失败")
	}

	return nil
}

// GetTeacherExtensionList 获取教师的延期申请列表
func (s *ProjectService) GetTeacherExtensionList(
	teacherID uint,
	query models.ExtensionListQuery,
) ([]models.ExtensionApplicationListResponse, int64, error) {

	db := s.db.Table("project_extension_applications pea").
		Joins("LEFT JOIN projects p ON pea.project_id = p.id").
		Joins("LEFT JOIN users u ON pea.student_id = u.id").
		Where("p.teacher_id = ?", teacherID)

	if query.Status != "" {
		db = db.Where("pea.status = ?", query.Status)
	}
	// count
	var total int64
	countDB := s.db.
		Table("project_extension_applications AS pea").
		Joins("LEFT JOIN projects p ON pea.project_id = p.id").
		Where("pea.teacher_id = ?", teacherID)

	if err := countDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// list
	var list []models.ExtensionApplicationListResponse
	listDB := s.db.
		Table("project_extension_applications AS pea").
		Joins("LEFT JOIN projects p ON pea.project_id = p.id").
		Joins("LEFT JOIN users s ON pea.student_id = s.id").
		Where("pea.teacher_id = ?", teacherID)

	err := listDB.
		Select(`
		pea.id,
		pea.project_id,
		p.title AS project_title,

		pea.student_id,
		s.username AS student_name,

		pea.teacher_id,
		NULL AS teacher_name,

		pea.original_finish_time,
		pea.requested_finish_time,

		pea.apply_reason,
		pea.status,
		pea.review_reason,
		pea.reviewed_at,
		pea.created_at
	`).
		Order("pea.created_at DESC").
		Offset((query.Page - 1) * query.Size).
		Limit(query.Size).
		Scan(&list).Error

	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// GetStudentTeachers 获取学生的指导教师列�?
func (s *ProjectService) GetStudentTeachers(studentID uint) ([]models.TeacherListResponse, error) {
	var teachers []models.User
	err := s.db.Preload("Profile").
		Joins("JOIN student_teacher st ON users.id = st.teacher_id").
		Where("st.student_id = ?", studentID).
		Where("users.status = ?", "active").
		Find(&teachers).Error

	if err != nil {
		log.Printf("获取学生指导教师列表失败: %v", err)
		return nil, err
	}

	var responses []models.TeacherListResponse
	for _, teacher := range teachers {
		response := models.TeacherListResponse{
			ID:       teacher.ID,
			Username: teacher.Username,
			Email:    teacher.Email,
		}

		if teacher.Profile != nil {
			response.RealName = teacher.Profile.RealName
			response.Phone = teacher.Profile.Phone
			response.Department = teacher.Profile.Department
			response.Bio = teacher.Profile.Bio
		}

		responses = append(responses, response)
	}

	return responses, nil
}

// SubmitProject 提交项目审核
func (s *ProjectService) SubmitProject(projectID, studentID uint) error {
	// 检查项目是否存?
	var project models.Project
	if err := s.db.First(&project, projectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("项目不存在")
		}
		return err
	}

	// 检查权限：只有项目创建者可以提�?
	if project.StudentID != studentID {
		return errors.New("无权限提交此项目")
	}

	// 检查状态：只有草稿状态的项目可以提交
	if project.Status != "draft" {
		return errors.New("只有草稿状态的项目可以提交")
	}

	// 更新项目状态和提交时间
	now := time.Now()
	updates := map[string]interface{}{
		"status":       "submitted",
		"submitted_at": &now,
	}

	if err := s.db.Model(&project).Updates(updates).Error; err != nil {
		log.Printf("提交项目失败: %v", err)
		return errors.New("提交项目失败")
	}

	log.Printf("项目提交成功 - 项目ID: %d", projectID)
	return nil
}

// UnbindStudentTeacher 解绑学生和教�?
func (s *ProjectService) UnbindStudentTeacher(studentID, teacherID uint) error {
	result := s.db.Where("student_id = ? AND teacher_id = ?", studentID, teacherID).Delete(&models.StudentTeacher{})
	if result.Error != nil {
		log.Printf("解绑学生教师失败: %v", result.Error)
		return errors.New("解绑失败")
	}

	if result.RowsAffected == 0 {
		return errors.New("绑定关系不存在")
	}

	log.Printf("学生教师解绑成功 - 学生ID: %d, 教师ID: %d", studentID, teacherID)
	return nil
}

func (s *ProjectService) GetDB() *gorm.DB {
	return s.db
}

// =============================================
// 1. 项目状态管理增强 - 新增方法
// =============================================

// UpdateProjectStatus 更新项目状态
func (s *ProjectService) UpdateProjectStatus(projectID uint, userID uint, req models.ProjectStatusUpdateRequest) error {
	// 检查项目是否存在
	var project models.Project
	if err := s.db.First(&project, projectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("项目不存在")
		}
		return err
	}

	// 检查权限：只有项目相关人员可以更新状态
	if project.StudentID != userID && project.TeacherID != userID {
		// 检查是否为管理员
		var user models.User
		if err := s.db.Preload("Roles").First(&user, userID).Error; err != nil {
			return errors.New("用户不存在")
		}

		isAdmin := false
		for _, role := range user.Roles {
			if role.RoleKey == "admin" {
				isAdmin = true
				break
			}
		}

		if !isAdmin {
			return errors.New("无权限更新此项目状态")
		}
	}

	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 记录状态变更历史
	statusHistory := models.ProjectStatusHistory{
		ProjectID:    projectID,
		OldStatus:    project.Status,
		NewStatus:    req.Status,
		ChangeReason: req.StatusChangeReason,
		ChangedBy:    userID,
		ChangedAt:    time.Now(),
	}

	if err := tx.Create(&statusHistory).Error; err != nil {
		tx.Rollback()
		log.Printf("创建状态变更历史失败: %v", err)
		return errors.New("创建状态变更历史失败")
	}

	// 更新项目状态
	updates := map[string]interface{}{
		"status":               req.Status,
		"status_change_reason": req.StatusChangeReason,
		"status_changed_by":    userID,
		"status_changed_at":    time.Now(),
		"updated_at":           time.Now(),
	}

	if err := tx.Model(&project).Updates(updates).Error; err != nil {
		tx.Rollback()
		log.Printf("更新项目状态失败: %v", err)
		return errors.New("更新项目状态失败")
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Printf("提交事务失败: %v", err)
		return errors.New("更新项目状态失败")
	}

	log.Printf("项目状态更新成功 - 项目ID: %d, 新状态: %s", projectID, req.Status)
	return nil
}

// GetProjectStatusHistory 获取项目状态变更历史
func (s *ProjectService) GetProjectStatusHistory(projectID uint) ([]models.ProjectStatusHistoryResponse, error) {
	var histories []models.ProjectStatusHistory
	err := s.db.Where("project_id = ?", projectID).
		Preload("ChangedByUser.Profile").
		Order("changed_at DESC").
		Find(&histories).Error

	if err != nil {
		return nil, err
	}

	var responses []models.ProjectStatusHistoryResponse
	for _, history := range histories {
		response := models.ProjectStatusHistoryResponse{
			ID:           history.ID,
			ProjectID:    history.ProjectID,
			OldStatus:    history.OldStatus,
			NewStatus:    history.NewStatus,
			ChangeReason: history.ChangeReason,
			ChangedBy:    history.ChangedBy,
			ChangedAt:    history.ChangedAt,
		}

		// 获取操作人姓名
		if history.ChangedByUser != nil && history.ChangedByUser.Profile != nil {
			response.OperatorName = history.ChangedByUser.Profile.RealName
		} else {
			response.OperatorName = history.ChangedByUser.Username
		}

		responses = append(responses, response)
	}

	return responses, nil
}

// =============================================
// 2. 项目生命周期管理增强 - 新增方法
// =============================================

// CreateProjectMilestone 创建项目里程碑
func (s *ProjectService) CreateProjectMilestone(projectID uint, userID uint, req models.ProjectMilestoneCreateRequest) (*models.ProjectMilestoneResponse, error) {
	// 检查项目是否存在
	var project models.Project
	if err := s.db.First(&project, projectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("项目不存在")
		}
		return nil, err
	}

	// 检查权限：只有项目相关人员可以创建里程碑
	if project.StudentID != userID && project.TeacherID != userID {
		return nil, errors.New("无权限为此项目创建里程碑")
	}

	// 创建里程碑
	milestone := models.ProjectMilestone{
		ProjectID:   projectID,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Status:      "pending",
		Progress:    0,
	}

	if err := s.db.Create(&milestone).Error; err != nil {
		log.Printf("创建项目里程碑失败: %v", err)
		return nil, errors.New("创建项目里程碑失败")
	}

	// 转换为响应格式
	response := &models.ProjectMilestoneResponse{
		ID:            milestone.ID,
		ProjectID:     milestone.ProjectID,
		Title:         milestone.Title,
		Description:   milestone.Description,
		DueDate:       milestone.DueDate,
		CompletedDate: milestone.CompletedDate,
		Status:        milestone.Status,
		Progress:      milestone.Progress,
		CreatedAt:     milestone.CreatedAt,
		UpdatedAt:     milestone.UpdatedAt,
	}

	log.Printf("项目里程碑创建成功 - 项目ID: %d, 里程碑ID: %d", projectID, milestone.ID)
	return response, nil
}

// UpdateProjectMilestone 更新项目里程碑
func (s *ProjectService) UpdateProjectMilestone(milestoneID uint, userID uint, req models.ProjectMilestoneUpdateRequest) error {
	// 检查里程碑是否存在
	var milestone models.ProjectMilestone
	if err := s.db.First(&milestone, milestoneID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("里程碑不存在")
		}
		return nil
	}

	// 检查权限：只有项目相关人员可以更新里程碑
	var project models.Project
	if err := s.db.First(&project, milestone.ProjectID).Error; err != nil {
		return errors.New("项目不存在")
	}

	if project.StudentID != userID && project.TeacherID != userID {
		return errors.New("无权限更新此里程碑")
	}

	// 更新里程碑
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.DueDate != nil {
		updates["due_date"] = *req.DueDate
	}
	if req.Progress != nil {
		updates["progress"] = *req.Progress
		// 如果进度为100%，自动标记为完成
		if *req.Progress >= 100 {
			updates["status"] = "completed"
			updates["completed_date"] = time.Now()
		}
	}

	if len(updates) > 0 {
		updates["updated_at"] = time.Now()
		if err := s.db.Model(&milestone).Updates(updates).Error; err != nil {
			log.Printf("更新项目里程碑失败: %v", err)
			return errors.New("更新项目里程碑失败")
		}
	}

	log.Printf("项目里程碑更新成功 - 里程碑ID: %d", milestoneID)
	return nil
}

// GetProjectMilestones 获取项目里程碑列表
func (s *ProjectService) GetProjectMilestones(projectID uint) ([]models.ProjectMilestoneResponse, error) {
	var milestones []models.ProjectMilestone
	err := s.db.Where("project_id = ?", projectID).
		Order("due_date ASC").
		Find(&milestones).Error

	if err != nil {
		return nil, err
	}

	var responses []models.ProjectMilestoneResponse
	for _, milestone := range milestones {
		response := models.ProjectMilestoneResponse{
			ID:            milestone.ID,
			ProjectID:     milestone.ProjectID,
			Title:         milestone.Title,
			Description:   milestone.Description,
			DueDate:       milestone.DueDate,
			CompletedDate: milestone.CompletedDate,
			Status:        milestone.Status,
			Progress:      milestone.Progress,
			CreatedAt:     milestone.CreatedAt,
			UpdatedAt:     milestone.UpdatedAt,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

// UpdateProjectProgress 更新项目进度
func (s *ProjectService) UpdateProjectProgress(projectID uint, userID uint, req models.ProjectProgressUpdateRequest) error {
	// 检查项目是否存在
	var project models.Project
	if err := s.db.First(&project, projectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("项目不存在")
		}
		return nil
	}

	// 检查权限：只有项目相关人员可以更新进度
	if project.StudentID != userID && project.TeacherID != userID {
		return errors.New("无权限更新此项目进度")
	}

	// 检查项目状态：只有进行中的项目可以更新进度
	if project.Status != "in_progress" && project.Status != "approved" {
		return errors.New("只有进行中的项目可以更新进度")
	}

	// 更新项目进度
	updates := map[string]interface{}{
		"progress":   req.Progress,
		"updated_at": time.Now(),
	}

	// 如果进度为100%，自动标记为完成
	if req.Progress >= 100 {
		updates["status"] = "completed"
		updates["actual_end_date"] = time.Now()
	}

	if err := s.db.Model(&project).Updates(updates).Error; err != nil {
		log.Printf("更新项目进度失败: %v", err)
		return errors.New("更新项目进度失败")
	}

	log.Printf("项目进度更新成功 - 项目ID: %d, 新进度: %d%%", projectID, req.Progress)
	return nil
}

// =============================================
// 3. 成果文件管理增强 - 新增方法
// =============================================

// UploadProjectFile 上传项目文件（增强版）
func (s *ProjectService) UploadProjectFile(projectID uint, userID uint, req models.ProjectFileUploadRequest) (*models.ProjectFileEnhancedResponse, error) {
	// 检查项目是否存在
	var project models.Project
	if err := s.db.First(&project, projectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("项目不存在")
		}
		return nil, err
	}

	// 检查权限：只有项目相关人员可以上传文件
	if project.StudentID != userID && project.TeacherID != userID {
		return nil, errors.New("无权限为此项目上传文件")
	}

	// 创建文件记录
	file := models.ProjectFile{
		ProjectID:    projectID,
		FileName:     req.FileName,
		FileURL:      req.FileURL,
		FileType:     req.FileType,
		FileVersion:  req.FileVersion,
		ReviewStatus: "pending",
		IsPublic:     req.IsPublic,
		UploadTime:   time.Now(),
	}

	if err := s.db.Create(&file).Error; err != nil {
		log.Printf("创建项目文件记录失败: %v", err)
		return nil, errors.New("创建项目文件记录失败")
	}

	// 转换为响应格式
	response := &models.ProjectFileEnhancedResponse{
		ID:             file.ID,
		ProjectID:      file.ProjectID,
		FileName:       file.FileName,
		FileURL:        file.FileURL,
		FileType:       file.FileType,
		FileVersion:    file.FileVersion,
		ReviewStatus:   file.ReviewStatus,
		ReviewComments: file.ReviewComments,
		ReviewedBy:     file.ReviewedBy,
		ReviewedAt:     file.ReviewedAt,
		FileSize:       file.FileSize,
		DownloadCount:  file.DownloadCount,
		IsPublic:       file.IsPublic,
		UploadTime:     file.UploadTime,
	}

	log.Printf("项目文件上传成功 - 项目ID: %d, 文件ID: %d", projectID, file.ID)
	return response, nil
}

// ReviewProjectFile 审核项目文件
func (s *ProjectService) ReviewProjectFile(fileID uint, userID uint, req models.ProjectFileReviewRequest) error {
	// 检查文件是否存在
	var file models.ProjectFile
	if err := s.db.First(&file, fileID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文件不存在")
		}
		return nil
	}

	// 检查权限：只有指导教师或管理员可以审核文件
	var project models.Project
	if err := s.db.First(&project, file.ProjectID).Error; err != nil {
		return errors.New("项目不存在")
	}

	if project.TeacherID != userID {
		// 检查是否为管理员
		var user models.User
		if err := s.db.Preload("Roles").First(&user, userID).Error; err != nil {
			return errors.New("用户不存在")
		}

		isAdmin := false
		for _, role := range user.Roles {
			if role.RoleKey == "admin" {
				isAdmin = true
				break
			}
		}

		if !isAdmin {
			return errors.New("无权限审核此文件")
		}
	}

	// 更新文件审核状态
	updates := map[string]interface{}{
		"review_status":   req.ReviewStatus,
		"review_comments": req.ReviewComments,
		"reviewed_by":     userID,
		"reviewed_at":     time.Now(),
	}

	if err := s.db.Model(&file).Updates(updates).Error; err != nil {
		log.Printf("更新文件审核状态失败: %v", err)
		return errors.New("更新文件审核状态失败")
	}

	log.Printf("项目文件审核完成 - 文件ID: %d, 审核结果: %s", fileID, req.ReviewStatus)
	return nil
}

// GetProjectFilesByType 按类型获取项目文件
func (s *ProjectService) GetProjectFilesByType(projectID uint, fileType string) ([]models.ProjectFileEnhancedResponse, error) {
	var files []models.ProjectFile
	query := s.db.Where("project_id = ?", projectID)

	if fileType != "" {
		query = query.Where("file_type = ?", fileType)
	}

	err := query.Order("upload_time DESC").Find(&files).Error
	if err != nil {
		return nil, err
	}

	var responses []models.ProjectFileEnhancedResponse
	for _, file := range files {
		response := models.ProjectFileEnhancedResponse{
			ID:             file.ID,
			ProjectID:      file.ProjectID,
			FileName:       file.FileName,
			FileURL:        file.FileURL,
			FileType:       file.FileType,
			FileVersion:    file.FileVersion,
			ReviewStatus:   file.ReviewStatus,
			ReviewComments: file.ReviewComments,
			ReviewedBy:     file.ReviewedBy,
			ReviewedAt:     file.ReviewedAt,
			FileSize:       file.FileSize,
			DownloadCount:  file.DownloadCount,
			IsPublic:       file.IsPublic,
			UploadTime:     file.UploadTime,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

// GetFileTypeConfigs 获取文件类型配置
func (s *ProjectService) GetFileTypeConfigs() ([]models.FileTypeConfigResponse, error) {
	var configs []models.FileTypeConfig
	err := s.db.Where("is_active = ?", true).Order("sort_order").Find(&configs).Error
	if err != nil {
		return nil, err
	}

	var responses []models.FileTypeConfigResponse
	for _, config := range configs {
		response := models.FileTypeConfigResponse{
			ID:                config.ID,
			FileType:          config.FileType,
			DisplayName:       config.DisplayName,
			Description:       config.Description,
			IsRequired:        config.IsRequired,
			MaxFileSize:       config.MaxFileSize,
			AllowedExtensions: config.AllowedExtensions,
			SortOrder:         config.SortOrder,
			IsActive:          config.IsActive,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

// =============================================
// 4. 项目分类管理增强 - 新增方法
// =============================================

// CreateProjectType 创建项目分类
func (s *ProjectService) CreateProjectType(req models.ProjectTypeCreateRequest) (*models.ProjectTypeEnhancedResponse, error) {
	// 创建项目分类
	projectType := models.ProjectType{
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
		Level:       req.Level,
		SortOrder:   req.SortOrder,
		IsActive:    req.IsActive,
		Icon:        req.Icon,
		Color:       req.Color,
	}

	if err := s.db.Create(&projectType).Error; err != nil {
		log.Printf("创建项目分类失败: %v", err)
		return nil, errors.New("创建项目分类失败")
	}

	// 转换为响应格式
	response := &models.ProjectTypeEnhancedResponse{
		ID:           projectType.ID,
		Name:         projectType.Name,
		Description:  projectType.Description,
		ParentID:     projectType.ParentID,
		Level:        projectType.Level,
		SortOrder:    projectType.SortOrder,
		IsActive:     projectType.IsActive,
		Icon:         projectType.Icon,
		Color:        projectType.Color,
		ProjectCount: 0,
		CreatedAt:    projectType.CreatedAt,
		UpdatedAt:    projectType.UpdatedAt,
	}

	log.Printf("项目分类创建成功 - 分类ID: %d, 名称: %s", projectType.ID, projectType.Name)
	return response, nil
}

// UpdateProjectType 更新项目分类
func (s *ProjectService) UpdateProjectType(typeID uint, req models.ProjectTypeUpdateRequest) error {
	// 检查分类是否存在
	var projectType models.ProjectType
	if err := s.db.First(&projectType, typeID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("项目分类不存在")
		}
		return nil
	}

	// 更新分类信息
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.ParentID != nil {
		updates["parent_id"] = *req.ParentID
	}
	if req.Level > 0 {
		updates["level"] = req.Level
	}
	if req.SortOrder > 0 {
		updates["sort_order"] = req.SortOrder
	}
	if req.Icon != "" {
		updates["icon"] = req.Icon
	}
	if req.Color != "" {
		updates["color"] = req.Color
	}

	if len(updates) > 0 {
		updates["updated_at"] = time.Now()
		if err := s.db.Model(&projectType).Updates(updates).Error; err != nil {
			log.Printf("更新项目分类失败: %v", err)
			return errors.New("更新项目分类失败")
		}
	}

	log.Printf("项目分类更新成功 - 分类ID: %d", typeID)
	return nil
}

// GetProjectTypeTree 获取项目分类树
func (s *ProjectService) GetProjectTypeTree() ([]models.ProjectTypeTreeResponse, error) {
	var projectTypes []models.ProjectType
	err := s.db.Where("is_active = ?", true).Order("sort_order").Find(&projectTypes).Error
	if err != nil {
		return nil, err
	}

	// 构建分类树
	return s.buildProjectTypeTree(projectTypes, nil), nil
}

// buildProjectTypeTree 构建项目分类树
func (s *ProjectService) buildProjectTypeTree(types []models.ProjectType, parentID *uint) []models.ProjectTypeTreeResponse {
	var tree []models.ProjectTypeTreeResponse

	for _, t := range types {
		if (parentID == nil && t.ParentID == nil) || (parentID != nil && t.ParentID != nil && *t.ParentID == *parentID) {
			node := models.ProjectTypeTreeResponse{
				ID:           t.ID,
				Name:         t.Name,
				Description:  t.Description,
				Level:        t.Level,
				SortOrder:    t.SortOrder,
				Icon:         t.Icon,
				Color:        t.Color,
				ProjectCount: int(t.ProjectCount),
				Children:     s.buildProjectTypeTree(types, &t.ID),
			}
			tree = append(tree, node)
		}
	}

	return tree
}

// GetProjectTypeStats 获取项目分类统计
func (s *ProjectService) GetProjectTypeStats() ([]models.ProjectTypeStatsResponse, error) {
	var stats []models.ProjectTypeStats
	err := s.db.Find(&stats).Error
	if err != nil {
		return nil, err
	}

	var responses []models.ProjectTypeStatsResponse
	for _, stat := range stats {
		response := models.ProjectTypeStatsResponse{
			TypeID:             stat.TypeID,
			TypeName:           stat.TypeName,
			TotalProjects:      int(stat.TotalProjects),
			DraftProjects:      int(stat.DraftProjects),
			SubmittedProjects:  int(stat.SubmittedProjects),
			ApprovedProjects:   int(stat.ApprovedProjects),
			InProgressProjects: int(stat.InProgressProjects),
			CompletedProjects:  int(stat.CompletedProjects),
			RejectedProjects:   int(stat.RejectedProjects),
			LastUpdated:        stat.LastUpdated,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

// =============================================
// 5. 审核流程增强 - 新增方法
// =============================================

// CreateReviewFlow 创建审核流程配置
func (s *ProjectService) CreateReviewFlow(req models.ReviewFlowCreateRequest) (*models.ReviewFlowResponse, error) {
	// 检查项目类型ID
	if req.ProjectTypeID == nil {
		return nil, errors.New("项目类型ID不能为空")
	}

	// 创建审核流程配置
	flow := models.ProjectReviewFlow{
		ProjectTypeID:      *req.ProjectTypeID,
		ReviewLevel:        req.ReviewLevel,
		ReviewerRole:       req.ReviewerRole,
		ReviewerDepartment: req.ReviewerDepartment,
		ReviewOrder:        req.ReviewOrder,
		IsRequired:         req.IsRequired,
		DeadlineHours:      req.DeadlineHours,
		AutoApprove:        req.AutoApprove,
		CanDelegate:        req.CanDelegate,
	}

	if err := s.db.Create(&flow).Error; err != nil {
		log.Printf("创建审核流程配置失败: %v", err)
		return nil, errors.New("创建审核流程配置失败")
	}

	// 转换为响应格式
	response := &models.ReviewFlowResponse{
		ID:                 flow.ID,
		ProjectTypeID:      &flow.ProjectTypeID,
		ReviewLevel:        flow.ReviewLevel,
		ReviewerRole:       flow.ReviewerRole,
		ReviewerDepartment: flow.ReviewerDepartment,
		ReviewOrder:        flow.ReviewOrder,
		IsRequired:         flow.IsRequired,
		DeadlineHours:      flow.DeadlineHours,
		AutoApprove:        flow.AutoApprove,
		CanDelegate:        flow.CanDelegate,
		CreatedAt:          flow.CreatedAt,
		UpdatedAt:          flow.UpdatedAt,
	}

	log.Printf("审核流程配置创建成功 - 流程ID: %d", flow.ID)
	return response, nil
}

// DelegateReview 委托审核
func (s *ProjectService) DelegateReview(reviewID uint, userID uint, req models.ReviewDelegationRequest) (*models.ReviewDelegationResponse, error) {
	// 检查审核记录是否存在
	var review models.ProjectReview
	if err := s.db.First(&review, reviewID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("审核记录不存在")
		}
		return nil, err
	}

	// 检查权限：只有原审核人可以委托
	if review.ReviewerID != userID {
		return nil, errors.New("无权限委托此审核任务")
	}

	// 创建委托记录
	delegation := models.ReviewDelegation{
		OriginalReviewerID:  userID,
		DelegatedReviewerID: req.DelegatedReviewerID,
		ProjectID:           review.ProjectID,
		Reason:              req.Reason,
		StartDate:           time.Now(),
		EndDate:             req.EndDate,
		Status:              "active",
	}

	if err := s.db.Create(&delegation).Error; err != nil {
		log.Printf("创建审核委托失败: %v", err)
		return nil, errors.New("创建审核委托失败")
	}

	// 转换为响应格式
	response := &models.ReviewDelegationResponse{
		ID:                  delegation.ID,
		OriginalReviewerID:  delegation.OriginalReviewerID,
		DelegatedReviewerID: delegation.DelegatedReviewerID,
		ProjectID:           delegation.ProjectID,
		Reason:              delegation.Reason,
		StartDate:           delegation.StartDate,
		EndDate:             delegation.EndDate,
		Status:              delegation.Status,
		CreatedAt:           delegation.CreatedAt,
	}

	log.Printf("审核委托创建成功 - 委托ID: %d", delegation.ID)
	return response, nil
}

// GetMyReviewTasks 获取我的审核任务
func (s *ProjectService) GetMyReviewTasks(userID uint, params models.ReviewTaskQueryParams) ([]models.ReviewTaskResponse, int64, error) {
	var reviews []models.ProjectReview
	var total int64

	query := s.db.Model(&models.ProjectReview{}).Where("reviewer_id = ?", userID)

	// 应用查询参数
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
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
	if err := query.Preload("Project").Find(&reviews).Error; err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	var responses []models.ReviewTaskResponse
	for _, review := range reviews {
		response := models.ReviewTaskResponse{
			ID:           review.ID,
			ProjectID:    review.ProjectID,
			ProjectTitle: review.Project.Title,
			ReviewLevel:  review.ReviewLevel,
			ReviewOrder:  review.ReviewOrder,
			Deadline:     review.Deadline,
			IsUrgent:     review.IsUrgent,
			Status:       review.Status,
			CreatedAt:    review.CreatedAt,
		}

		// 获取学生姓名
		if review.Project.Student != nil && review.Project.Student.Profile != nil {
			response.StudentName = review.Project.Student.Profile.RealName
		}

		// 获取项目类型
		response.ProjectType = review.Project.Type

		responses = append(responses, response)
	}

	return responses, total, nil
}

// GetReviewFlowConfig 获取审核流程配置
func (s *ProjectService) GetReviewFlowConfig(projectTypeID *uint) ([]models.ReviewFlowResponse, error) {
	var flows []models.ProjectReviewFlow
	query := s.db.Model(&models.ProjectReviewFlow{})

	if projectTypeID != nil {
		query = query.Where("project_type_id = ?", *projectTypeID)
	}

	err := query.Order("review_level, review_order").Find(&flows).Error
	if err != nil {
		return nil, err
	}

	var responses []models.ReviewFlowResponse
	for _, flow := range flows {
		response := models.ReviewFlowResponse{
			ID:                 flow.ID,
			ProjectTypeID:      &flow.ProjectTypeID,
			ReviewLevel:        flow.ReviewLevel,
			ReviewerRole:       flow.ReviewerRole,
			ReviewerDepartment: flow.ReviewerDepartment,
			ReviewOrder:        flow.ReviewOrder,
			IsRequired:         flow.IsRequired,
			DeadlineHours:      flow.DeadlineHours,
			AutoApprove:        flow.AutoApprove,
			CanDelegate:        flow.CanDelegate,
			CreatedAt:          flow.CreatedAt,
			UpdatedAt:          flow.UpdatedAt,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

// ForceUpdateProjectStatus 强制更新项目状态
func (s *ProjectService) ForceUpdateProjectStatus(projectID uint, status string, reason string, operatorID uint) error {
	// 检查项目是否存在
	var project models.Project
	if err := s.db.First(&project, projectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("项目不存在")
		}
		return err
	}

	// 更新项目状态
	updates := map[string]interface{}{
		"status":              status,
		"force_status_reason": reason,
		"updated_at":          time.Now(),
	}

	if err := s.db.Model(&project).Updates(updates).Error; err != nil {
		log.Printf("强制更新项目状态失败: %v", err)
		return errors.New("强制更新项目状态失败")
	}

	log.Printf("项目状态强制更新成功 - 项目ID: %d, 新状态: %s", projectID, status)
	return nil
}

// SoftDeleteProject 软删除项目
func (s *ProjectService) SoftDeleteProject(projectID uint, operatorID uint) error {
	// 检查项目是否存在
	var project models.Project
	if err := s.db.First(&project, projectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("项目不存在")
		}
		return err
	}

	// 软删除项目
	if err := s.db.Model(&project).Update("deleted", true).Error; err != nil {
		log.Printf("软删除项目失败: %v", err)
		return errors.New("软删除项目失败")
	}

	log.Printf("项目软删除成功 - 项目ID: %d", projectID)
	return nil
}

// RestoreProject 恢复软删除的项目
func (s *ProjectService) RestoreProject(projectID uint, operatorID uint) error {
	// 检查项目是否存在
	var project models.Project
	if err := s.db.First(&project, projectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("项目不存在")
		}
		return err
	}

	// 恢复项目
	if err := s.db.Model(&project).Update("deleted", false).Error; err != nil {
		log.Printf("恢复项目失败: %v", err)
		return errors.New("恢复项目失败")
	}

	log.Printf("项目恢复成功 - 项目ID: %d", projectID)
	return nil
}

// GetTeacherProjects 获取教师项目列表
func (s *ProjectService) GetTeacherProjects(teacherID uint, params models.TeacherProjectQueryParams) ([]models.ProjectListResponse, int64, error) {
	var projects []models.Project
	var total int64

	query := s.db.Model(&models.Project{}).Where("deleted = ?", false)

	// 应用查询参数
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.Type != "" {
		query = query.Where("type = ?", params.Type)
	}
	if params.Level != "" {
		query = query.Where("level = ?", params.Level)
	}
	if params.CategoryID != nil {
		query = query.Where("category_id = ?", *params.CategoryID)
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
	if err := query.Preload("Student").Preload("Teacher").Find(&projects).Error; err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	var responses []models.ProjectListResponse
	for _, project := range projects {
		response := models.ProjectListResponse{
			ID:          project.ID,
			Title:       project.Title,
			Description: project.Description,
			Type:        project.Type,
			Status:      project.Status,
			StudentName: project.Student.Username,
			StudentID:   strconv.FormatUint(uint64(project.Student.ID), 10),
			TeacherName: project.Teacher.Username,
			TeacherID:   project.Teacher.ID,
			SubmittedAt: project.SubmittedAt,
			CreatedAt:   project.CreatedAt,
			UpdatedAt:   project.UpdatedAt,
		}

		// 获取成员数量
		var memberCount int64
		s.db.Model(&models.ProjectMember{}).Where("project_id = ?", project.ID).Count(&memberCount)
		response.MemberCount = int(memberCount)

		responses = append(responses, response)
	}

	return responses, total, nil
}

func (s *ProjectService) GetStudentProjectsFiles(
	teacherID uint,
	projectID uint,
	studentID uint,
) ([]models.File, error) {

	var files []models.File

	// 基础查询：files 表
	query := s.db.Model(&models.File{}).
		Joins("JOIN projects ON projects.id = files.project_id").
		Where("projects.teacher_id = ? and files.deleted_at is NULL", teacherID)

	// 按项目过滤
	if projectID != 0 {
		query = query.Where("files.project_id = ?", projectID)
	}

	// 按学生过滤
	if studentID != 0 {
		query = query.Where("files.uploaded_by = ?", studentID)
	}

	// 排序（最新上传的在前）
	err := query.
		Order("files.created_at DESC").
		Find(&files).Error

	if err != nil {
		return nil, err
	}

	return files, nil
}

// GetTeacherListWithFilter 获取教师列表（带过滤）
func (s *ProjectService) GetTeacherListWithFilter(params models.TeacherQueryParams) ([]models.TeacherListResponse, int64, error) {
	var teachers []models.User
	var total int64

	query := s.db.Model(&models.User{}).Where("role = 'teacher' AND deleted = ?", false)

	// 应用查询参数
	if params.Department != "" {
		query = query.Where("department = ?", params.Department)
	}
	if params.Title != "" {
		query = query.Where("title = ?", params.Title)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
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
	if err := query.Preload("Profile").Find(&teachers).Error; err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	var responses []models.TeacherListResponse
	for _, teacher := range teachers {
		response := models.TeacherListResponse{
			ID:         teacher.ID,
			Username:   teacher.Username,
			Email:      teacher.Email,
			Department: teacher.Department,
			Title:      teacher.Title,
			Status:     teacher.Status,
			CreatedAt:  teacher.CreatedAt,
		}

		if teacher.Profile != nil {
			response.RealName = teacher.Profile.RealName
			response.Phone = teacher.Profile.Phone
		}

		responses = append(responses, response)
	}

	return responses, total, nil
}

// BindStudentToTeacher 绑定学生到教师
func (s *ProjectService) BindStudentToTeacher(studentID uint, req models.StudentBindTeacherRequest) error {
	// 检查学生是否存在
	var student models.User
	if err := s.db.First(&student, studentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("学生不存在")
		}
		return err
	}

	// 检查教师是否存在
	var teacher models.User
	if err := s.db.First(&teacher, req.TeacherID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("教师不存在")
		}
		return err
	}

	// 检查是否已经绑定
	var existingBinding models.StudentTeacher
	if err := s.db.Where("student_id = ? AND teacher_id = ?", studentID, req.TeacherID).First(&existingBinding).Error; err == nil {
		return errors.New("学生和教师已经绑定")
	}

	// 创建绑定关系
	binding := models.StudentTeacher{
		StudentID: studentID,
		TeacherID: req.TeacherID,
		BindTime:  time.Now(),
	}

	if err := s.db.Create(&binding).Error; err != nil {
		log.Printf("绑定学生到教师失败: %v", err)
		return errors.New("绑定学生到教师失败")
	}

	log.Printf("学生绑定到教师成功 - 学生ID: %d, 教师ID: %d", studentID, req.TeacherID)
	return nil
}

// ValidateProjectUpdate 验证项目更新
func (s *ProjectService) ValidateProjectUpdate(projectID uint, updates map[string]interface{}) error {
	// 检查项目是否存在
	var project models.Project
	if err := s.db.First(&project, projectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("项目不存在")
		}
		return err
	}

	// 检查项目状态是否允许更新
	if project.Status == "submitted" || project.Status == "approved" {
		return errors.New("项目已提交或已审核，不允许修改")
	}

	// 验证更新字段
	for field, value := range updates {
		switch field {
		case "title":
			if title, ok := value.(string); ok && len(title) > 100 {
				return errors.New("项目标题不能超过100个字符")
			}
		case "description":
			if desc, ok := value.(string); ok && len(desc) > 1000 {
				return errors.New("项目描述不能超过1000个字符")
			}
		case "type":
			if projectType, ok := value.(string); ok {
				if projectType != "科研" && projectType != "竞赛" {
					return errors.New("项目类型只能是科研或竞赛")
				}
			}
		case "level":
			if level, ok := value.(string); ok {
				if level != "校级" && level != "省级" && level != "国家级" {
					return errors.New("项目级别只能是校级、省级或国家级")
				}
			}
		}
	}

	return nil
}

// GetMyStudents 获取我的学生列表
func (s *ProjectService) GetMyStudents(teacherID uint, params models.StudentQueryParams) ([]models.StudentListResponse, int64, error) {
	var students []models.User
	var total int64

	query := s.db.Model(&models.User{}).Where("role_name = ?", "student")

	// 应用查询参数
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.Grade != "" {
		query = query.Where("grade = ?", params.Grade)
	}
	if params.Major != "" {
		query = query.Where("major = ?", params.Major)
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
	if err := query.Preload("Profile").Find(&students).Error; err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	var responses []models.StudentListResponse
	for _, student := range students {
		response := models.StudentListResponse{
			ID:        student.ID,
			Name:      student.Username,
			Email:     student.Email,
			Major:     student.Major,
			Grade:     student.Grade,
			Status:    student.Status,
			CreatedAt: student.CreatedAt,
		}

		if student.Profile != nil {
			response.RealName = student.Profile.RealName
			response.Phone = student.Profile.Phone
		}

		responses = append(responses, response)
	}

	return responses, total, nil
}

// GetStudentProjects 获取学生项目列表（支持分页和查询）
func (s *ProjectService) GetStudentProjects(studentID uint, params models.ProjectQueryParams) ([]models.ProjectMyListResponse, int64, error) {
	var projects []models.Project
	var total int64

	query := s.db.Model(&models.Project{}).Where("student_id = ?", studentID)

	// 搜索条件
	if params.Search != "" {
		search := "%" + params.Search + "%"
		query = query.Where("title LIKE ? OR description LIKE ?", search, search)
	}

	// 类型筛选
	if params.Type != "" {
		query = query.Where("type = ?", params.Type)
	}

	// 状态筛选
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		log.Printf("获取学生项目总数失败: %v", err)
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
	if err := query.Find(&projects).Error; err != nil {
		log.Printf("获取学生项目列表失败: %v", err)
		return nil, 0, err
	}

	// 转换为响应格式
	var responses []models.ProjectMyListResponse
	for _, project := range projects {
		responses = append(responses, models.ProjectMyListResponse{
			ID:        project.ID,
			Title:     project.Title,
			Type:      project.Type,
			Status:    project.Status,
			CreatedAt: project.CreatedAt,
		})
	}

	return responses, total, nil
}

// GetStudentProjectStats 获取学生项目统计信息
func (s *ProjectService) GetStudentProjectStats(studentID uint) (*models.ProjectStats, error) {
	var stats models.ProjectStats

	// 获取项目总数
	if err := s.db.Model(&models.Project{}).Where("student_id = ?", studentID).Count(&stats.TotalProjects).Error; err != nil {
		log.Printf("获取学生项目总数失败: %v", err)
		return nil, err
	}

	// 获取各状态项目数量
	statuses := []string{"draft", "pending", "approved", "rejected"}
	for _, status := range statuses {
		var count int64
		if err := s.db.Model(&models.Project{}).Where("student_id = ? AND status = ?", studentID, status).Count(&count).Error; err != nil {
			log.Printf("获取状态 %s 项目数量失败: %v", status, err)
			continue
		}

		switch status {
		case "draft":
			stats.DraftProjects = count
		case "pending":
			stats.PendingProjects = count
		case "approved":
			stats.ApprovedProjects = count
		case "rejected":
			stats.RejectedProjects = count
		}
	}

	// 获取类型统计
	if err := s.db.Model(&models.Project{}).Where("student_id = ? AND type = ?", studentID, "科研").Count(&stats.ResearchProjects).Error; err != nil {
		log.Printf("获取科研项目数量失败: %v", err)
	}
	if err := s.db.Model(&models.Project{}).Where("student_id = ? AND type = ?", studentID, "竞赛").Count(&stats.CompetitionProjects).Error; err != nil {
		log.Printf("获取竞赛项目数量失败: %v", err)
	}

	return &stats, nil
}
