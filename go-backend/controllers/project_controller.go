package controllers

import (
	"log"
	"net/http"
	"strconv"

	"yunmeng-backend/models"
	"yunmeng-backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectController struct {
	DB             *gorm.DB
	projectService *services.ProjectService
}

func NewProjectController(projectService *services.ProjectService) *ProjectController {
	return &ProjectController{
		projectService: projectService,
	}
}

// GetProjectList 获取项目列表（教师/管理员）
func (c *ProjectController) GetProjectList(ctx *gin.Context) {
	var params models.ProjectQueryParams
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

	projects, total, err := c.projectService.GetProjectsForTeacher(params)
	if err != nil {
		log.Printf("获取项目列表失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取项目列表失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "查询成功",
		"data": gin.H{
			"list":  projects,
			"total": total,
		},
	})
}

// GetProjectByID 根据ID获取项目详情
func (c *ProjectController) GetProjectByID(ctx *gin.Context) {
	// 从查询参数获取 id
	idStr := ctx.Query("id")

	// 参数为空校验
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "项目ID不能为空",
		})
		return
	}

	// 转换为 uint32
	parsedID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "项目ID格式错误",
		})
		return
	}

	id := uint32(parsedID)

	project, err := c.projectService.GetProjectByID(uint(id))
	if err != nil {
		log.Printf("获取项目详情失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取项目详情失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取项目详情成功",
		"data":    project,
	})
}

// CreateProject 创建项目
func (c *ProjectController) CreateProject(ctx *gin.Context) {
	var req models.ProjectCreateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.Printf("参数绑定失败: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 从JWT中获取学生ID
	studentID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权访问",
		})
		return
	}

	response, err := c.projectService.CreateProject(studentID.(uint), req)
	if err != nil {
		log.Printf("创建项目失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建项目失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"message": "项目提交成功",
		"data":    response,
	})
}

// UpdateProject 更新项目
func (c *ProjectController) UpdateProject(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "项目ID格式错误",
		})
		return
	}

	var req models.ProjectUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("参数绑定失败: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 从JWT中获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权访问",
		})
		return
	}

	err = c.projectService.UpdateProject(uint(id), userID.(uint), req)
	if err != nil {
		log.Printf("更新项目失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新项目失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "项目更新成功",
	})
}

// CreateExtensionApplication 创建项目延期申请
func (c *ProjectController) CreateExtensionApplication(ctx *gin.Context) {
	var req models.ExtensionApplicationRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权访问",
		})
		return
	}

	studentID, ok := userID.(uint)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "用户ID格式错误",
		})
		return
	}

	resp, err := c.projectService.CreateExtensionApplication(studentID, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"message": "项目延期申请提交成功",
		"data":    resp,
	})
}

// DeleteProject 删除项目
func (c *ProjectController) DeleteProject(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "项目ID格式错误",
		})
		return
	}

	err = c.projectService.DeleteProject(uint(id))
	if err != nil {
		log.Printf("删除项目失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除项目失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "项目删除成功",
	})
}

// GetMyExtensionApplications 学生查看我的延期申请
func (c *ProjectController) GetMyExtensionApplications(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code": 401, "message": "未授权访问",
		})
		return
	}

	studentID := userID.(uint)

	var query models.ExtensionListQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400, "message": "参数错误",
		})
		return
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Size <= 0 || query.Size > 50 {
		query.Size = 10
	}

	list, total, err := c.projectService.GetStudentExtensionList(
		studentID, query,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500, "message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "查询成功",
		"data": gin.H{
			"list":  list,
			"total": total,
			"page":  query.Page,
			"size":  query.Size,
		},
	})
}

// ReviewProject 审核项目
func (c *ProjectController) ReviewProject(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "项目ID格式错误",
		})
		return
	}

	var req models.ProjectReviewRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("参数绑定失败: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 从JWT中获取审核者ID
	reviewerID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权访问",
		})
		return
	}

	response, err := c.projectService.ReviewProjectWithResponse(uint(id), reviewerID.(uint), req)
	if err != nil {
		log.Printf("审核项目失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "审核项目失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "审核完成",
		"data":    response,
	})
}

// ForceUpdateProjectStatus 管理员强制更新项目状态
func (c *ProjectController) ForceUpdateProjectStatus(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "项目ID格式错误",
		})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=draft submitted approved rejected archived"`
		Reason string `json:"reason" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("参数绑定失败: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 从JWT中获取操作者ID
	operatorID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权访问",
		})
		return
	}

	// 检查操作者是否为管理员
	userRole, exists := ctx.Get("userRole")
	if !exists || userRole != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "只有管理员可以强制更新项目状态",
		})
		return
	}

	err = c.projectService.ForceUpdateProjectStatus(uint(id), req.Status, req.Reason, operatorID.(uint))
	if err != nil {
		log.Printf("强制更新项目状态失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "强制更新项目状态失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "项目状态强制更新成功",
	})
}

// GetProjectStats 获取项目统计信息
func (c *ProjectController) GetProjectStats(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权访问",
		})
		return
	}

	// userID 类型安全转换
	var uid uint
	switch v := userID.(type) {
	case uint:
		uid = v
	case int:
		uid = uint(v)
	case float64:
		uid = uint(v)
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "用户ID格式错误",
		})
		return
	}

	// 调用 Service
	stats, err := c.projectService.GetProjectStats(uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取项目统计失败",
		})
		return
	}

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "查询成功",
		"data":    stats,
	})
}

// ExportProjects 导出项目数据
func (c *ProjectController) ExportProjects(ctx *gin.Context) {
	var req models.ProjectExportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("参数绑定失败: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// TODO: 实现项目数据导出功能
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "项目数据导出功能开发中",
		"data": gin.H{
			"format":  req.Format,
			"filters": req.Filters,
		},
	})
}

// SubmitProject 提交项目审核
func (c *ProjectController) SubmitProject(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "项目ID格式错误",
		})
		return
	}

	// 从JWT中获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权访问",
		})
		return
	}

	err = c.projectService.SubmitProject(uint(id), userID.(uint))
	if err != nil {
		log.Printf("提交项目失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "提交项目失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "项目提交成功，等待审核",
	})
}

// GetProjectReviews 获取项目审核记录
func (c *ProjectController) GetProjectReviews(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "项目ID格式错误",
		})
		return
	}

	reviews, err := c.projectService.GetProjectReviews(uint(id))
	if err != nil {
		log.Printf("获取审核记录失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取审核记录失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    reviews,
	})
}

func (c *ProjectController) GetMyProjects(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权访问",
		})
		return
	}

	// 类型转换
	var uid uint
	switch v := userID.(type) {
	case uint:
		uid = v
	case int:
		uid = uint(v)
	case float64:
		uid = uint(v)
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "用户ID格式错误",
		})
		return
	}

	// 获取查询参数
	pageStr := ctx.DefaultQuery("page", "1")
	sizeStr := ctx.DefaultQuery("size", "10")
	status := ctx.Query("status")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	size, _ := strconv.Atoi(sizeStr)
	if size < 1 || size > 100 {
		size = 10
	}

	// 调用service，传入分页参数
	projects, total, err := c.projectService.GetMyProjects(uid, status, page, size)
	if err != nil {
		log.Printf("获取我的项目失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取我的项目失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "查询成功",
		"data": gin.H{
			"list":  projects,
			"total": total,
			"page":  page,
			"size":  size,
		},
	})
}

// BindStudentTeacher 绑定学生和教师
func (c *ProjectController) BindStudentTeacher(ctx *gin.Context) {
	var req models.StudentTeacherBindRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("参数绑定失败: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	result, err := c.projectService.BindStudentTeacher(req)
	if err != nil {
		log.Printf("绑定学生教师失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "绑定失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "绑定成功",
		"data":    result,
	})
}

// GetTeacherList 获取教师列表
func (c *ProjectController) GetTeacherList(ctx *gin.Context) {
	teachers, err := c.projectService.GetTeacherList()
	if err != nil {
		log.Printf("获取教师列表失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取教师列表失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取教师列表成功",
		"data":    teachers,
	})
}

// ApproveExtensionApplication 教师审批延期申请
func (c *ProjectController) ApproveExtensionApplication(ctx *gin.Context) {
	var req models.ExtensionApprovalRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权访问",
		})
		return
	}

	teacherID, ok := userID.(uint)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "用户ID格式错误",
		})
		return
	}

	if req.Action == "rejected" && req.Reason == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "驳回延期申请时必须填写原因",
		})
		return
	}

	if err := c.projectService.ApproveExtensionApplication(teacherID, req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "延期申请处理成功",
	})
}

// GetTeacherExtensionApplications 教师查看延期申请列表
func (c *ProjectController) GetTeacherExtensionApplications(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code": 401, "message": "未授权访问",
		})
		return
	}

	teacherID := userID.(uint)

	var query models.ExtensionListQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400, "message": "参数错误",
		})
		return
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Size <= 0 || query.Size > 50 {
		query.Size = 10
	}
	if query.Status == "" {
		query.Status = "pending"
	}

	list, total, err := c.projectService.GetTeacherExtensionList(
		teacherID, query,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500, "message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "查询成功",
		"data": gin.H{
			"list":  list,
			"total": total,
			"page":  query.Page,
			"size":  query.Size,
		},
	})
}

// GetStudentTeachers 获取学生的指导教师列表
func (c *ProjectController) GetStudentTeachers(ctx *gin.Context) {
	studentIDStr := ctx.Param("studentId")
	studentID, err := strconv.ParseUint(studentIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "学生ID格式错误",
		})
		return
	}

	teachers, err := c.projectService.GetStudentTeachers(uint(studentID))
	if err != nil {
		log.Printf("获取学生指导教师列表失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取指导教师列表失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取指导教师列表成功",
		"data":    teachers,
	})
}

// UnbindStudentTeacher 解绑学生和教师
func (c *ProjectController) UnbindStudentTeacher(ctx *gin.Context) {
	studentIDStr := ctx.Param("studentId")
	teacherIDStr := ctx.Param("teacherId")

	studentID, err := strconv.ParseUint(studentIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "学生ID格式错误",
		})
		return
	}

	teacherID, err := strconv.ParseUint(teacherIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "教师ID格式错误",
		})
		return
	}

	err = c.projectService.UnbindStudentTeacher(uint(studentID), uint(teacherID))
	if err != nil {
		log.Printf("解绑学生教师失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "解绑失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "解绑成功",
	})
}

// GetTeacherProjects 获取当前登录教师的所有指导项目
func (c *ProjectController) GetTeacherProjects(ctx *gin.Context) {
	// 从JWT中获取当前用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未获取到用户信息",
		})
		return
	}

	teacherID := userID.(uint)

	var params models.TeacherProjectQueryParams
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

	projects, total, err := c.projectService.GetTeacherProjects(teacherID, params)
	if err != nil {
		log.Printf("获取教师指导项目失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取项目列表失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "查询成功",
		"data": gin.H{
			"list":  projects,
			"total": total,
			"page":  params.Page,
			"size":  params.Size,
		},
	})
}

// GetTeacherListWithFilter 获取教师列表（支持院系筛选）
func (c *ProjectController) GetTeacherListWithFilter(ctx *gin.Context) {
	var params models.TeacherQueryParams
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

	teachers, total, err := c.projectService.GetTeacherListWithFilter(params)
	if err != nil {
		log.Printf("获取教师列表失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取教师列表失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "查询成功",
		"data": gin.H{
			"list":  teachers,
			"total": total,
			"page":  params.Page,
			"size":  params.Size,
		},
	})
}

// BindStudentToTeacher 学生绑定教师（学生端接口）
func (c *ProjectController) BindStudentToTeacher(ctx *gin.Context) {
	// 从JWT中获取当前用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未获取到用户信息",
		})
		return
	}

	studentID := userID.(uint)

	var req models.StudentBindTeacherRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("参数绑定失败: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	if err := c.projectService.BindStudentToTeacher(studentID, req); err != nil {
		log.Printf("学生绑定教师失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "绑定失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "绑定成功",
	})
}

// UpdateProjectWithValidation 更新项目（带验证）
func (c *ProjectController) UpdateProjectWithValidation(ctx *gin.Context) {
	// 从JWT中获取当前用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未获取到用户信息",
		})
		return
	}

	studentID := userID.(uint)

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "项目ID格式错误",
		})
		return
	}

	// 验证项目更新权限
	err = c.projectService.ValidateProjectUpdate(uint(id), map[string]interface{}{"student_id": studentID})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	var req models.ProjectUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("参数绑定失败: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	err = c.projectService.UpdateProject(uint(id), studentID, req)
	if err != nil {
		log.Printf("更新项目失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新项目失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新项目成功",
	})
}

// GetMyStudents 获取当前登录教师指导的学生列表
func (c *ProjectController) GetMyStudents(ctx *gin.Context) {
	// 从JWT中获取当前用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未获取到用户信息",
		})
		return
	}

	teacherID := userID.(uint)

	var params models.StudentQueryParams
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

	students, total, err := c.projectService.GetMyStudents(teacherID, params)
	if err != nil {
		log.Printf("获取教师指导的学生列表失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取学生列表失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "查询成功",
		"data": gin.H{
			"list":  students,
			"total": total,
			"page":  params.Page,
			"size":  params.Size,
		},
	})
}

// =============================================
// 1. 项目状态管理增强 API
// =============================================

// UpdateProjectStatus 更新项目状态（学生/教师/管理员）
func (c *ProjectController) UpdateProjectStatus(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "项目ID格式错1误",
		})
		return
	}

	var req models.ProjectStatusUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 获取当前用户ID
	userID := ctx.GetUint("user_id")

	err = c.projectService.UpdateProjectStatus(uint(id), userID, req)
	if err != nil {
		log.Printf("更新项目状态失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新项目状态失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "项目状态更新成功",
	})
}

// GetProjectStatusHistory 获取项目状态变更历史
func (c *ProjectController) GetProjectStatusHistory(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "项目ID格式错误",
		})
		return
	}

	history, err := c.projectService.GetProjectStatusHistory(uint(id))
	if err != nil {
		log.Printf("获取项目状态历史失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取项目状态历史失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取项目状态历史成功",
		"data":    history,
	})
}

// =============================================
// 2. 项目生命周期管理增强 API
// =============================================

// CreateProjectMilestone 创建项目里程碑
func (c *ProjectController) CreateProjectMilestone(ctx *gin.Context) {
	projectIDStr := ctx.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "项目ID格式错误",
		})
		return
	}

	var req models.ProjectMilestoneCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	userID := ctx.GetUint("user_id")
	milestone, err := c.projectService.CreateProjectMilestone(uint(projectID), userID, req)
	if err != nil {
		log.Printf("创建项目里程碑失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建项目里程碑失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建项目里程碑成功",
		"data":    milestone,
	})
}

// UpdateProjectMilestone 更新项目里程碑
func (c *ProjectController) UpdateProjectMilestone(ctx *gin.Context) {
	milestoneIDStr := ctx.Param("milestoneId")
	milestoneID, err := strconv.ParseUint(milestoneIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "里程碑ID格式错误",
		})
		return
	}

	var req models.ProjectMilestoneUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	userID := ctx.GetUint("user_id")
	err = c.projectService.UpdateProjectMilestone(uint(milestoneID), userID, req)
	if err != nil {
		log.Printf("更新项目里程碑失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新项目里程碑失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新项目里程碑成功",
	})
}

// GetProjectMilestones 获取项目里程碑列表
func (c *ProjectController) GetProjectMilestones(ctx *gin.Context) {
	projectIDStr := ctx.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "项目ID格式错误",
		})
		return
	}

	milestones, err := c.projectService.GetProjectMilestones(uint(projectID))
	if err != nil {
		log.Printf("获取项目里程碑失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取项目里程碑失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取项目里程碑成功",
		"data":    milestones,
	})
}

// UpdateProjectProgress 更新项目进度
func (c *ProjectController) UpdateProjectProgress(ctx *gin.Context) {
	projectIDStr := ctx.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "项目ID格式错误",
		})
		return
	}

	var req models.ProjectProgressUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	userID := ctx.GetUint("user_id")
	err = c.projectService.UpdateProjectProgress(uint(projectID), userID, req)
	if err != nil {
		log.Printf("更新项目进度失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新项目进度失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新项目进度成功",
	})
}

// =============================================
// 3. 成果文件管理增强 API
// =============================================

// UploadProjectFile 上传项目文件（增强版）
func (c *ProjectController) UploadProjectFile(ctx *gin.Context) {
	projectIDStr := ctx.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "项目ID格式错误",
		})
		return
	}

	var req models.ProjectFileUploadRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	userID := ctx.GetUint("user_id")
	file, err := c.projectService.UploadProjectFile(uint(projectID), userID, req)
	if err != nil {
		log.Printf("上传项目文件失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "上传项目文件失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "上传项目文件成功",
		"data":    file,
	})
}

// ReviewProjectFile 审核项目文件（教师/管理员）
func (c *ProjectController) ReviewProjectFile(ctx *gin.Context) {
	fileIDStr := ctx.Param("fileId")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文件ID格式错误",
		})
		return
	}

	var req models.ProjectFileReviewRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	userID := ctx.GetUint("user_id")
	err = c.projectService.ReviewProjectFile(uint(fileID), userID, req)
	if err != nil {
		log.Printf("审核项目文件失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "审核项目文件失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "审核项目文件成功",
	})
}

// GetProjectFilesByType 按类型获取项目文件
func (c *ProjectController) GetProjectFilesByType(ctx *gin.Context) {
	projectIDStr := ctx.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "项目ID格式错误",
		})
		return
	}

	fileType := ctx.Query("type")
	files, err := c.projectService.GetProjectFilesByType(uint(projectID), fileType)
	if err != nil {
		log.Printf("获取项目文件失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取项目文件失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取项目文件成功",
		"data":    files,
	})
}

// GetFileTypeConfigs 获取文件类型配置
func (c *ProjectController) GetFileTypeConfigs(ctx *gin.Context) {
	configs, err := c.projectService.GetFileTypeConfigs()
	if err != nil {
		log.Printf("获取文件类型配置失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取文件类型配置失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取文件类型配置成功",
		"data":    configs,
	})
}

// =============================================
// 4. 项目分类管理增强 API
// =============================================

// CreateProjectType 创建项目分类（管理员）
func (c *ProjectController) CreateProjectType(ctx *gin.Context) {
	var req models.ProjectTypeCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	projectType, err := c.projectService.CreateProjectType(req)
	if err != nil {
		log.Printf("创建项目分类失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建项目分类失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建项目分类成功",
		"data":    projectType,
	})
}

// UpdateProjectType 更新项目分类（管理员）
func (c *ProjectController) UpdateProjectType(ctx *gin.Context) {
	typeIDStr := ctx.Param("id")
	typeID, err := strconv.ParseUint(typeIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "分类ID格式错误",
		})
		return
	}

	var req models.ProjectTypeUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	err = c.projectService.UpdateProjectType(uint(typeID), req)
	if err != nil {
		log.Printf("更新项目分类失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新项目分类失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新项目分类成功",
	})
}

// GetProjectTypeTree 获取项目分类树
func (c *ProjectController) GetProjectTypeTree(ctx *gin.Context) {
	tree, err := c.projectService.GetProjectTypeTree()
	if err != nil {
		log.Printf("获取项目分类树失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取项目分类树失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取项目分类树成功",
		"data":    tree,
	})
}

// GetProjectTypeStats 获取项目分类统计
func (c *ProjectController) GetProjectTypeStats(ctx *gin.Context) {
	stats, err := c.projectService.GetProjectTypeStats()
	if err != nil {
		log.Printf("获取项目分类统计失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取项目分类统计失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取项目分类统计成功",
		"data":    stats,
	})
}

// =============================================
// 5. 审核流程增强 API
// =============================================

// CreateReviewFlow 创建审核流程配置（管理员）
func (c *ProjectController) CreateReviewFlow(ctx *gin.Context) {
	var req models.ReviewFlowCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	flow, err := c.projectService.CreateReviewFlow(req)
	if err != nil {
		log.Printf("创建审核流程失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建审核流程失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建审核流程成功",
		"data":    flow,
	})
}

// DelegateReview 委托审核
func (c *ProjectController) DelegateReview(ctx *gin.Context) {
	reviewIDStr := ctx.Param("reviewId")
	reviewID, err := strconv.ParseUint(reviewIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "审核ID格式错误",
		})
		return
	}

	var req models.ReviewDelegationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	userID := ctx.GetUint("user_id")
	delegation, err := c.projectService.DelegateReview(uint(reviewID), userID, req)
	if err != nil {
		log.Printf("委托审核失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "委托审核失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "委托审核成功",
		"data":    delegation,
	})
}

// GetMyReviewTasks 获取我的审核任务
func (c *ProjectController) GetMyReviewTasks(ctx *gin.Context) {
	var params models.ReviewTaskQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	userID := ctx.GetUint("user_id")
	tasks, total, err := c.projectService.GetMyReviewTasks(userID, params)
	if err != nil {
		log.Printf("获取审核任务失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取审核任务失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取审核任务成功",
		"data": gin.H{
			"list":  tasks,
			"total": total,
		},
	})
}

// GetReviewFlowConfig 获取审核流程配置
func (c *ProjectController) GetReviewFlowConfig(ctx *gin.Context) {
	projectTypeIDStr := ctx.Query("projectTypeId")
	var projectTypeID *uint
	if projectTypeIDStr != "" {
		if id, err := strconv.ParseUint(projectTypeIDStr, 10, 32); err == nil {
			tempID := uint(id)
			projectTypeID = &tempID
		}
	}

	config, err := c.projectService.GetReviewFlowConfig(projectTypeID)
	if err != nil {
		log.Printf("获取审核流程配置失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取审核流程配置失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取审核流程配置成功",
		"data":    config,
	})
}

// GetStudentProjects 获取学生项目列表
func (c *ProjectController) GetStudentProjects(ctx *gin.Context) {
	var params models.ProjectQueryParams
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

	// 获取当前登录用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户未认证",
		})
		return
	}

	// 调用服务层获取学生项目
	projects, total, err := c.projectService.GetStudentProjects(uint(userID.(float64)), params)
	if err != nil {
		log.Printf("获取学生项目列表失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取学生项目列表失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "查询成功",
		"data": gin.H{
			"list":  projects,
			"total": total,
		},
	})
}

// GetStudentProjectStats 获取学生项目统计信息
func (c *ProjectController) GetStudentProjectStats(ctx *gin.Context) {
	// 获取当前登录用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户未认证",
		})
		return
	}

	// 调用服务层获取学生项目统计
	stats, err := c.projectService.GetStudentProjectStats(uint(userID.(float64)))
	if err != nil {
		log.Printf("获取学生项目统计失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取学生项目统计失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取统计信息成功",
		"data":    stats,
	})
}

// GetStudentProjectByID 获取学生项目详情
func (c *ProjectController) GetStudentProjectByID(ctx *gin.Context) {
	projectID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "项目ID格式错误",
		})
		return
	}

	// 调用服务层获取项目详情
	project, err := c.projectService.GetProjectByID(uint(projectID))
	if err != nil {
		log.Printf("获取项目详情失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取项目详情失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取项目详情成功",
		"data":    project,
	})
}

// CreateStudentProject 创建学生项目
func (c *ProjectController) CreateStudentProject(ctx *gin.Context) {
	var req models.ProjectCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("参数绑定失败: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 获取当前登录用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户未认证",
		})
		return
	}

	// 调用服务层创建项目
	response, err := c.projectService.CreateProject(uint(userID.(float64)), req)
	if err != nil {
		log.Printf("创建项目失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建项目失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建项目成功",
		"data":    response,
	})
}

// UpdateStudentProject 更新学生项目
func (c *ProjectController) UpdateStudentProject(ctx *gin.Context) {
	projectID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "项目ID格式错误",
		})
		return
	}

	var req models.ProjectUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("参数绑定失败: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 获取当前登录用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户未认证",
		})
		return
	}

	// 调用服务层更新项目
	err = c.projectService.UpdateProject(uint(projectID), uint(userID.(float64)), req)
	if err != nil {
		log.Printf("更新项目失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新项目失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新项目成功",
	})
}

// DeleteStudentProject 删除学生项目
func (c *ProjectController) DeleteStudentProject(ctx *gin.Context) {
	projectID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "项目ID格式错误",
		})
		return
	}

	// 调用服务层删除项目
	err = c.projectService.DeleteProject(uint(projectID))
	if err != nil {
		log.Printf("删除项目失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除项目失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除项目成功",
	})
}

// SubmitStudentProject 提交学生项目
func (c *ProjectController) SubmitStudentProject(ctx *gin.Context) {
	projectID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "项目ID格式错误",
		})
		return
	}

	// 获取当前登录用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户未认证",
		})
		return
	}

	// 调用服务层提交项目
	err = c.projectService.SubmitProject(uint(projectID), uint(userID.(float64)))
	if err != nil {
		log.Printf("提交项目失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "提交项目失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "提交项目成功",
	})
}

// UpdateStudentProjectProgress 更新学生项目进度
func (c *ProjectController) UpdateStudentProjectProgress(ctx *gin.Context) {
	projectID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "项目ID格式错误",
		})
		return
	}

	var req models.ProjectProgressUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("参数绑定失败: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 获取当前登录用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户未认证",
		})
		return
	}

	// 调用服务层更新项目进度
	err = c.projectService.UpdateProjectProgress(uint(projectID), uint(userID.(float64)), req)
	if err != nil {
		log.Printf("更新项目进度失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新项目进度失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新项目进度成功",
	})
}

// GetStudentProjectFiles 获取学生项目文件
func (c *ProjectController) GetStudentProjectFiles(ctx *gin.Context) {
	projectID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "项目ID格式错误",
		})
		return
	}

	fileType := ctx.Query("type")

	// 调用服务层获取项目文件
	files, err := c.projectService.GetProjectFilesByType(uint(projectID), fileType)
	if err != nil {
		log.Printf("获取项目文件失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取项目文件失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取项目文件成功",
		"data":    files,
	})
}

// UploadStudentProjectFile 上传学生项目文件
func (c *ProjectController) UploadStudentProjectFile(ctx *gin.Context) {
	projectID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "项目ID格式错误",
		})
		return
	}

	var req models.ProjectFileUploadRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("参数绑定失败: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 获取当前登录用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户未认证",
		})
		return
	}

	// 调用服务层上传项目文件
	response, err := c.projectService.UploadProjectFile(uint(projectID), uint(userID.(float64)), req)
	if err != nil {
		log.Printf("上传项目文件失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "上传项目文件失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "上传项目文件成功",
		"data":    response,
	})
}

// DeleteStudentProjectFile 删除学生项目文件
func (c *ProjectController) DeleteStudentProjectFile(ctx *gin.Context) {
	fileID, err := strconv.ParseUint(ctx.Param("fileId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文件ID格式错误",
		})
		return
	}

	// 获取当前登录用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    400,
			"message": "用户未认证",
		})
		return
	}

	// 调用服务层删除项目文件 - 使用拒绝状态来表示删除
	err = c.projectService.ReviewProjectFile(uint(fileID), uint(userID.(float64)), models.ProjectFileReviewRequest{
		ReviewStatus:   "rejected",
		ReviewComments: "文件已删除",
	})
	if err != nil {
		log.Printf("删除项目文件失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除项目文件失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除项目文件成功",
	})
}

// GetStudentProjectReviews 获取学生项目审核记录
func (c *ProjectController) GetStudentProjectReviews(ctx *gin.Context) {
	projectID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "项目ID格式错误",
		})
		return
	}

	// 调用服务层获取项目审核记录
	reviews, err := c.projectService.GetProjectReviews(uint(projectID))
	if err != nil {
		log.Printf("获取项目审核记录失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取项目审核记录失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取项目审核记录成功",
		"data":    reviews,
	})
}

// ExportStudentProjects 导出学生项目数据
func (c *ProjectController) ExportStudentProjects(ctx *gin.Context) {
	// 获取当前登录用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户未认证",
		})
		return
	}

	// 这里需要实现导出学生项目数据的逻辑
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "导出项目数据成功",
		"data": gin.H{
			"userId": userID,
		},
	})
}

// GetStudentProjectReport 获取学生项目报告
func (c *ProjectController) GetStudentProjectReport(ctx *gin.Context) {
	// 获取当前登录用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户未认证",
		})
		return
	}

	// 这里需要实现获取学生项目报告的逻辑
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取项目报告成功",
		"data": gin.H{
			"userId": userID,
		},
	})
}

// GetStudentProjectSuggestions 获取学生项目建议
func (c *ProjectController) GetStudentProjectSuggestions(ctx *gin.Context) {
	// 获取当前登录用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户未认证",
		})
		return
	}

	// 这里需要实现获取学生项目建议的逻辑
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取项目建议成功",
		"data": gin.H{
			"userId": userID,
		},
	})
}
