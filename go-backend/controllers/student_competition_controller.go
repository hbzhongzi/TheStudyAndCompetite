package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"yunmeng-backend/models"
	"yunmeng-backend/utils"
)

// StudentCompetitionController 学生竞赛控制器
type StudentCompetitionController struct {
	DB *gorm.DB
}

// NewStudentCompetitionController 创建学生竞赛控制器
func NewStudentCompetitionController(db *gorm.DB) *StudentCompetitionController {
	return &StudentCompetitionController{DB: db}
}

// GetAvailableCompetitions 获取学生可报名的竞赛
// @Summary 获取学生可报名的竞赛
// @Description 获取学生可报名的竞赛列表，支持按院系、状态筛选
// @Tags 学生竞赛
// @Accept json
// @Produce json
// @Param department query string false "院系筛选"
// @Param status query string false "竞赛状态筛选"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} utils.Response{data=[]models.CompetitionResponse}
// @Router /api/competitions/available [get]
func (cc *StudentCompetitionController) GetAvailableCompetitions(c *gin.Context) {
	var competitions []models.Competition
	query := cc.DB.Model(&models.Competition{}).
		Where("is_open = ? AND start_time > ?", true, time.Now()).
		Preload("CreatedByUser")

	// 院系筛选
	if department := c.Query("department"); department != "" {
		query = query.Where("department_limit IS NULL OR department_limit = '' OR FIND_IN_SET(?, department_limit) > 0", department)
	}

	// 状态筛选
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	offset := (page - 1) * size

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(size).Find(&competitions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取竞赛列表失败: " + err.Error(),
		})
		return
	}

	// 转换为响应格式
	var responses []models.CompetitionResponse
	for _, comp := range competitions {
		response := models.CompetitionResponse{
			ID:                  comp.ID,
			Title:               comp.Title,
			Description:         comp.Description,
			CurrentParticipants: comp.CurrentParticipants,
			Status:              comp.Status,
			CreatedBy:           comp.CreatedBy,
			CreatedAt:           comp.CreatedAt,
			UpdatedAt:           comp.UpdatedAt,
		}

		responses = append(responses, response)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取竞赛列表成功",
		"data": gin.H{
			"competitions": responses,
			"total":        total,
			"page":         page,
			"size":         size,
		},
	})
}

// GetMyRegistrations 获取我的报名记录
// @Summary 获取我的报名记录
// @Description 查看学生的竞赛报名记录及审核状态
// @Tags 学生竞赛
// @Accept json
// @Produce json
// @Param status query string false "报名状态筛选"
// @Param type query string false "竞赛类型筛选"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} utils.Response{data=[]models.CompetitionRegistrationResponse}
// @Router /api/competitions/my-registrations [get]
func (cc *StudentCompetitionController) GetMyRegistrations(c *gin.Context) {
	userID := utils.GetCurrentUserID(c)
	if userID == 0 {
		utils.ResponseError(c, http.StatusUnauthorized, "用户未登录", nil)
		return
	}

	var registrations []models.CompetitionRegistration
	query := cc.DB.Model(&models.CompetitionRegistration{}).
		Where("student_id = ?", userID).
		Preload("Competition.CreatedByUser").
		Preload("Student").
		Preload("Teacher")

	// 状态筛选
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 竞赛类型筛选
	if compType := c.Query("type"); compType != "" {
		query = query.Joins("JOIN competitions ON competition_registrations.competition_id = competitions.id").
			Where("competitions.type = ?", compType)
	}

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	offset := (page - 1) * size

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(size).Order("register_time DESC").Find(&registrations).Error; err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "获取报名记录失败", err)
		return
	}

	// 转换为响应格式
	var responses []models.CompetitionRegistrationResponse
	for _, reg := range registrations {
		response := models.CompetitionRegistrationResponse{
			ID:            reg.ID,
			CompetitionID: reg.CompetitionID,
			RegisterTime:  reg.RegistrationTime,
			Status:        reg.Status,
			TeamName:      reg.TeamName,
			TeamLeader:    reg.TeamLeader,
		}

		// 设置关联数据
		if reg.Competition != nil {
			response.Competition = &models.CompetitionResponseInfo{
				ID:          reg.Competition.ID,
				Title:       reg.Competition.Title,
				Description: reg.Competition.Description,
				Status:      reg.Competition.Status,
			}
		}
		responses = append(responses, response)
	}

	utils.ResponseSuccess(c, gin.H{
		"registrations": responses,
		"total":         total,
		"page":          page,
		"size":          size,
	})
}

// SubmitCompetitionResult 学生提交竞赛成果
func (c *CompetitionController) SubmitCompetitionResult(ctx *gin.Context) {

	//获取学生竞赛提交相关数据s
	var req models.SubmitCompetitionResultRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.Printf("参数绑定失败: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 2. 获取当前学生ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code": 401, "message": "用户未登录",
		})
		return
	}
	studentID := userID.(uint)

	// 3. 校验报名状态（必须 approved）
	var registration models.CompetitionRegistration
	err := c.db.
		Where("competition_id = ? AND student_id = ?", req.CompetitionID, studentID).
		First(&registration).Error

	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "您未报名该竞赛或报名不存在",
		})
		return
	}

	if registration.Status != "approved" {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "报名尚未通过审核，无法提交成果",
		})
		return
	}

	// 4. 接收文件
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请上传成果文件",
		})
		return
	}

	// 5. 保存文件（示例：本地）
	saveDir := "uploads/competition_results"
	_ = os.MkdirAll(saveDir, os.ModePerm)

	filename := fmt.Sprintf(
		"%d_%d_%d_%s",
		req.CompetitionID,
		studentID,
		time.Now().Unix(),
		file.Filename,
	)

	filePath := filepath.Join(saveDir, filename)
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "文件保存失败",
		})
		return
	}

	// 6. 检查是否已有提交（用于版本控制）
	var existing models.CompetitionSubmission

	err = c.db.
		Where("competition_id = ? AND student_id = ?", req.CompetitionID, studentID).
		Order("submitted_at DESC").
		First(&existing).Error
	version := "1.0"
	if err == nil {
		if existing.Locked {
			ctx.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "当前成果已锁定，无法再次提交",
			})
			return
		}
		version = increaseVersion(existing.Version)
	}

	// 7. 保存提交记录
	submission := models.CompetitionSubmission{
		CompetitionID: req.CompetitionID,
		StudentID:     studentID,
		Title:         req.Title,
		Description:   req.Description,
		FileURL:       filePath,
		FileSize:      file.Size,
		Version:       version,
		Locked:        false,
		SubmittedAt:   time.Now(),
	}

	if err := c.db.Save(&submission).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "提交成果失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "成果提交成功",
		"data":    submission,
	})

}

// 版本号增加函数，简单实现为在原版本号基础上加0.1
func increaseVersion(old string) string {
	v, err := strconv.ParseFloat(old, 64)
	if err != nil {
		return "1.0"
	}
	return fmt.Sprintf("%.1f", v+0.1)
}

// GetCompetitionResults 获取竞赛成绩
// @Summary 获取竞赛成绩
// @Description 查看学生的竞赛成绩和获奖情况
// @Tags 学生竞赛
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} utils.Response{data=[]models.CompetitionResultResponse}
// @Router /api/competitions/results [get]
func (cc *StudentCompetitionController) GetCompetitionResults(c *gin.Context) {
	userID := utils.GetCurrentUserID(c)
	if userID == 0 {
		utils.ResponseError(c, http.StatusUnauthorized, "用户未登录", nil)
		return
	}

	var results []models.CompetitionResult
	query := cc.DB.Model(&models.CompetitionResult{}).
		Where("student_id = ?", userID).
		Preload("Competition").
		Preload("Student").
		Preload("Submission").
		Preload("CreatedByUser")

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	offset := (page - 1) * size

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(size).Order("publish_time DESC").Find(&results).Error; err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "获取成绩失败", err)
		return
	}

	// 转换为响应格式
	var responses []models.CompetitionResultResponse
	for _, result := range results {
		response := models.CompetitionResultResponse{
			ID:             result.ID,
			CompetitionID:  result.CompetitionID,
			StudentID:      result.StudentID,
			SubmissionID:   result.SubmissionID,
			AwardLevel:     result.AwardLevel,
			FinalScore:     result.FinalScore,
			CertificateURL: result.CertificateURL,
			PublishTime:    result.PublishTime,
			CreatedBy:      result.CreatedBy,
			FinalizedBy:    result.FinalizedBy,
			FinalizedAt:    result.FinalizedAt,
		}

		// 设置关联数据
		if result.Competition != nil {
			response.Competition = &models.CompetitionResponse{
				ID:     result.Competition.ID,
				Title:  result.Competition.Title,
				Status: result.Competition.Status,
			}
		}

		if result.Student != nil {
			response.Student = &models.Users{
				ID:         result.Student.ID,
				Username:   result.Student.Username,
				Email:      result.Student.Email,
				Status:     result.Student.Status,
				CreateTime: result.Student.CreateTime,
			}
		}

		if result.CreatedByUser != nil {
			response.CreatedByUser = &models.Users{
				ID:         result.CreatedByUser.ID,
				Username:   result.CreatedByUser.Username,
				Email:      result.CreatedByUser.Email,
				Status:     result.CreatedByUser.Status,
				CreateTime: result.CreatedByUser.CreateTime,
			}
		}

		responses = append(responses, response)
	}

	utils.ResponseSuccess(c, gin.H{
		"results": responses,
		"total":   total,
		"page":    page,
		"size":    size,
	})
}
