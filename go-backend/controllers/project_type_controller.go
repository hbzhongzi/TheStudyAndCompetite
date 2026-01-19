package controllers

import (
	"log"
	"net/http"
	"strconv"
	"yunmeng-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectTypeController struct {
	db *gorm.DB
}

func NewProjectTypeController(db *gorm.DB) *ProjectTypeController {
	return &ProjectTypeController{db: db}
}

// GetProjectTypeList 获取项目分类列表
func (c *ProjectTypeController) GetProjectTypeList(ctx *gin.Context) {
	var projectTypes []models.ProjectType
	var total int64

	// 获取总数
	if err := c.db.Model(&models.ProjectType{}).Count(&total).Error; err != nil {
		log.Printf("获取项目分类总数失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取项目分类列表失败",
		})
		return
	}

	// 获取列表
	if err := c.db.Find(&projectTypes).Error; err != nil {
		log.Printf("获取项目分类列表失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取项目分类列表失败",
		})
		return
	}

	// 构建响应数据
	var responses []models.ProjectTypeResponse
	for _, pt := range projectTypes {
		var projectCount int64
		c.db.Model(&models.Project{}).Where("category_id = ?", pt.ID).Count(&projectCount)

		response := models.ProjectTypeResponse{
			ID:           pt.ID,
			Name:         pt.Name,
			Description:  pt.Description,
			ProjectCount: projectCount,
			CreatedAt:    pt.CreatedAt,
			UpdatedAt:    pt.UpdatedAt,
		}
		responses = append(responses, response)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取项目分类列表成功",
		"data": gin.H{
			"list":  responses,
			"total": total,
		},
	})
}

// GetProjectTypeByID 根据ID获取项目分类详情
func (c *ProjectTypeController) GetProjectTypeByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "项目分类ID格式错误",
		})
		return
	}

	var projectType models.ProjectType
	if err := c.db.First(&projectType, id).Error; err != nil {
		log.Printf("获取项目分类详情失败: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "项目分类不存在",
		})
		return
	}

	// 获取项目数量
	var projectCount int64
	c.db.Model(&models.Project{}).Where("category_id = ?", projectType.ID).Count(&projectCount)

	response := models.ProjectTypeResponse{
		ID:           projectType.ID,
		Name:         projectType.Name,
		Description:  projectType.Description,
		ProjectCount: projectCount,
		CreatedAt:    projectType.CreatedAt,
		UpdatedAt:    projectType.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取项目分类详情成功",
		"data":    response,
	})
}

// CreateProjectType 创建项目分类
func (c *ProjectTypeController) CreateProjectType(ctx *gin.Context) {
	var req models.ProjectTypeCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("参数绑定失败: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 检查分类名称是否已存在
	var count int64
	c.db.Model(&models.ProjectType{}).Where("name = ?", req.Name).Count(&count)
	if count > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "项目分类名称已存在",
		})
		return
	}

	projectType := models.ProjectType{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := c.db.Create(&projectType).Error; err != nil {
		log.Printf("创建项目分类失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建项目分类失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建项目分类成功",
		"data":    projectType,
	})
}

// UpdateProjectType 更新项目分类
func (c *ProjectTypeController) UpdateProjectType(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "项目分类ID格式错误",
		})
		return
	}

	var req models.ProjectTypeUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("参数绑定失败: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 检查项目分类是否存在
	var projectType models.ProjectType
	if err := c.db.First(&projectType, id).Error; err != nil {
		log.Printf("项目分类不存在: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "项目分类不存在",
		})
		return
	}

	// 检查分类名称是否已存在（排除当前记录）
	var count int64
	c.db.Model(&models.ProjectType{}).Where("name = ? AND id != ?", req.Name, id).Count(&count)
	if count > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "项目分类名称已存在",
		})
		return
	}

	// 更新项目分类
	updates := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
	}

	if err := c.db.Model(&projectType).Updates(updates).Error; err != nil {
		log.Printf("更新项目分类失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新项目分类失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新项目分类成功",
		"data":    projectType,
	})
}

// DeleteProjectType 删除项目分类
func (c *ProjectTypeController) DeleteProjectType(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "项目分类ID格式错误",
		})
		return
	}

	// 检查项目分类是否存在
	var projectType models.ProjectType
	if err := c.db.First(&projectType, id).Error; err != nil {
		log.Printf("项目分类不存在: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "项目分类不存在",
		})
		return
	}

	// 检查是否有项目使用此分类
	var projectCount int64
	c.db.Model(&models.Project{}).Where("category_id = ?", id).Count(&projectCount)
	if projectCount > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "该分类下还有项目，无法删除",
		})
		return
	}

	// 删除项目分类
	if err := c.db.Delete(&projectType).Error; err != nil {
		log.Printf("删除项目分类失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除项目分类失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除项目分类成功",
	})
}

// GetProjectTypeStats 获取项目分类统计信息
func (c *ProjectTypeController) GetProjectTypeStats(ctx *gin.Context) {
	var stats []gin.H

	// 获取所有项目分类
	var projectTypes []models.ProjectType
	if err := c.db.Find(&projectTypes).Error; err != nil {
		log.Printf("获取项目分类列表失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取项目分类统计失败",
		})
		return
	}

	// 为每个分类统计项目数量
	for _, pt := range projectTypes {
		var projectCount int64

		// 尝试多种查询方式
		query := c.db.Model(&models.Project{}).Where("category_id = ?", pt.ID)

		// 检查deleted字段是否存在
		if err := query.Count(&projectCount).Error; err != nil {
			// 如果失败，尝试不包含deleted字段的查询
			log.Printf("包含deleted字段的查询失败，尝试简化查询: %v", err)
			if err := c.db.Model(&models.Project{}).Where("category_id = ?", pt.ID).Count(&projectCount).Error; err != nil {
				log.Printf("统计分类 %d 项目数量失败: %v", pt.ID, err)
				projectCount = 0
			}
		}

		stats = append(stats, gin.H{
			"id":           pt.ID,
			"name":         pt.Name,
			"projectCount": projectCount,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取项目分类统计成功",
		"data":    stats,
	})
}
