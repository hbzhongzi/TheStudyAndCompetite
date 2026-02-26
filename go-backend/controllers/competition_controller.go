package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"
	"yunmeng-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CompetitionController struct {
	db *gorm.DB
}

func NewCompetitionController(db *gorm.DB) *CompetitionController {
	return &CompetitionController{db: db}
}

// GetCompetitionList 获取竞赛列表（单表查询）
func (c *CompetitionController) GetCompetitionList(ctx *gin.Context) {
	var list []models.Competition
	var total int64

	// 分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	// 查询参数
	search := ctx.Query("search")
	level := ctx.Query("level")
	category := ctx.Query("category")
	status := ctx.Query("status")

	query := c.db.Model(&models.Competition{})

	// 模糊搜索
	if search != "" {
		query = query.Where(
			"title LIKE ? OR description LIKE ?",
			"%"+search+"%",
			"%"+search+"%",
		)
	}

	if level != "" {
		query = query.Where("level = ?", level)
	}

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 总数
	if err := query.Count(&total).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取竞赛总数失败",
		})
		return
	}

	// 分页数据
	offset := (page - 1) * size
	if err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(size).
		Find(&list).Error; err != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取竞赛列表失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取竞赛列表成功",
		"data": gin.H{
			"list":  list,
			"total": total,
			"page":  page,
			"size":  size,
		},
	})
}

// CreateCompetition 创建竞赛
func (c *CompetitionController) CreateCompetition(ctx *gin.Context) {
	var req models.CompetitionCreateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 获取当前用户
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户未登录",
		})
		return
	}

	now := time.Now()

	// ===== time.Time → *time.Time（基于 IsZero）=====
	var (
		regStart *time.Time
		regEnd   *time.Time
		subStart *time.Time
		subEnd   *time.Time
	)

	if !req.RegistrationStart.IsZero() {
		regStart = &req.RegistrationStart
	}

	if !req.RegistrationEnd.IsZero() {
		regEnd = &req.RegistrationEnd
	}

	if !req.StartTime.IsZero() {
		subStart = &req.StartTime
	}

	if !req.EndTime.IsZero() {
		subEnd = &req.EndTime
	}

	// ===== 时间逻辑校验 =====

	if regStart != nil && regEnd != nil && regStart.After(*regEnd) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "报名开始时间不能晚于报名截止时间",
		})
		return
	}

	if regEnd != nil && subStart != nil && regEnd.After(*subStart) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "报名截止时间不能晚于提交开始时间",
		})
		return
	}

	if subStart != nil && subEnd != nil && subStart.After(*subEnd) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "提交开始时间不能晚于提交截止时间",
		})
		return
	}

	if subStart != nil && subStart.Before(now) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "提交开始时间不能早于当前时间",
		})
		return
	}

	if subEnd != nil && subEnd.Before(now) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "提交截止时间不能早于当前时间",
		})
		return
	}

	// ===== 构建模型 =====
	competition := models.Competition{
		Title:             req.Title,
		Description:       req.Description,
		Category:          req.Type,
		RegistrationStart: regStart,
		Level:             req.Level,
		RegistrationEnd:   regEnd,
		SubmissionStart:   subStart,
		SubmissionEnd:     subEnd,
		MaxParticipants:   req.MaxParticipants,
		Status:            "draft",
		CreatedBy:         userID.(uint),
	}

	if err := c.db.Create(&competition).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建竞赛成功",
		"data":    competition,
	})
}

// GetCompetitionDetail 获取竞赛详情
func (c *CompetitionController) GetCompetitionDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "竞赛ID格式错误",
		})
		return
	}

	var competition models.Competition
	if err := c.db.First(&competition, id).Error; err != nil {
		log.Printf("获取竞赛详情失败: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "竞赛不存在",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取竞赛详情成功",
		"data":    competition,
	})
}

// ToggleCompetitionOpen 切换竞赛开放状态
func (c *CompetitionController) ToggleCompetitionOpen(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "竞赛ID格式错误",
		})
		return
	}

	var competition models.Competition
	if err := c.db.First(&competition, id).Error; err != nil {
		log.Printf("获取竞赛详情失败: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "竞赛不存在",
		})
		return
	}

	isOpen := !competition.IsOpen

	if err := c.db.Model(&competition).Update("is_open", isOpen).Error; err != nil {
		log.Printf("切换竞赛开放状态失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "切换竞赛开放状态失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "切换竞赛开放状态成功",
	})
}

// DeleteCompetition 删除竞赛
func (c *CompetitionController) DeleteCompetition(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "竞赛ID格式错误",
		})
		return
	}

	// 检查竞赛是否存在
	var competition models.Competition
	if err := c.db.First(&competition, id).Error; err != nil {
		log.Printf("竞赛不存在: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "竞赛不存在",
		})
		return
	}

	// 检查是否有报名记录
	var registrationCount int64
	c.db.Model(&models.CompetitionRegistration{}).Where("competition_id = ?", id).Count(&registrationCount)
	if registrationCount > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "该竞赛已有报名记录，无法删除",
		})
		return
	}

	if err := c.db.Delete(&competition).Error; err != nil {
		log.Printf("删除竞赛失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除竞赛失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除竞赛成功",
	})
}

// GetCompetitionStats 获取竞赛统计信息
func (c *CompetitionController) GetCompetitionStats(ctx *gin.Context) {
	var stats []gin.H

	// 获取各类型竞赛数量
	rows, err := c.db.Raw(`
		SELECT 
			type,
			COUNT(*) as count,
			SUM(CASE WHEN is_open = true THEN 1 ELSE 0 END) as open_count
		FROM competitions
		GROUP BY type
		ORDER BY count DESC
	`).Rows()
	if err != nil {
		log.Printf("获取竞赛统计失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取竞赛统计失败",
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var competitionType string
		var count, openCount int
		rows.Scan(&competitionType, &count, &openCount)
		stats = append(stats, gin.H{
			"type":       competitionType,
			"count":      count,
			"open_count": openCount,
		})
	}

	// 获取总体统计
	var totalCompetitions, totalRegistrations, totalSubmissions, totalResults int64
	c.db.Model(&models.Competition{}).Count(&totalCompetitions)
	c.db.Model(&models.CompetitionRegistration{}).Count(&totalRegistrations)
	c.db.Model(&models.CompetitionSubmission{}).Count(&totalSubmissions)
	c.db.Model(&models.CompetitionResult{}).Count(&totalResults)

	overallStats := gin.H{
		"total_competitions":  totalCompetitions,
		"total_registrations": totalRegistrations,
		"total_submissions":   totalSubmissions,
		"total_results":       totalResults,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取竞赛统计成功",
		"data": gin.H{
			"type_stats": stats,
			"overall":    overallStats,
		},
	})
}

// RegisterCompetition 报名竞赛（学生）
func (c *CompetitionController) RegisterCompetition(ctx *gin.Context) {
	idStr := ctx.PostForm("id")
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少竞赛ID参数",
		})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "竞赛ID格式错误",
		})
		return
	}

	// 获取当前用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户未登录",
		})
		return
	}

	// 检查竞赛是否存在且开放报名
	var competition models.Competition
	if err := c.db.First(&competition, id).Error; err != nil {
		log.Printf("竞赛不存在: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "竞赛不存在",
		})
		return
	}

	// 检查是否已经报名
	var existingRegistration models.CompetitionRegistration
	if err := c.db.Where("competition_id = ? AND student_id = ?", id, userID).First(&existingRegistration).Error; err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "您已经报名该竞赛",
		})
		return
	}

	var req models.CompetitionRegistrationRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.Printf("参数绑定失败: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 创建报名记录
	registration := models.CompetitionRegistration{
		CompetitionID:    uint(id),
		StudentID:        userID.(uint),
		RegistrationTime: time.Now(),
		TeamName:         req.TeamName,
		TeamLeader:       req.TeamLeader,
		Status:           "pending",
	}

	if err := c.db.Create(&registration).Error; err != nil {
		log.Printf("报名竞赛失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "报名竞赛失败",
		})
		return
	}

	// 更新竞赛当前参与人数
	c.db.Model(&competition).Update("current_participants", competition.CurrentParticipants+1)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "报名竞赛成功",
		"data":    registration,
	})
}

// GetMyRegistrations 查看自己的报名记录（学生）
func (c *CompetitionController) GetMyRegistrations(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户未登录",
		})
		return
	}

	var registrations []models.CompetitionRegistration
	var total int64

	// 获取查询参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

	query := c.db.Model(&models.CompetitionRegistration{}).Where("student_id = ?", userID).Preload("Competition")

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		log.Printf("获取报名记录总数失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取报名记录失败",
		})
		return
	}

	// 分页查询
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("registration_time DESC").Find(&registrations).Error; err != nil {
		log.Printf("获取报名记录失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取报名记录失败",
		})
		return
	}

	// 构建响应数据
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

		if reg.Competition != nil {
			response.Competition = &models.CompetitionResponseInfo{
				ID:          reg.Competition.ID,
				Title:       reg.Competition.Title,
				Description: reg.Competition.Description,
			}
		}

		responses = append(responses, response)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取报名记录成功",
		"data": gin.H{
			"list":  responses,
			"total": total,
			"page":  page,
			"size":  size,
		},
	})
}

// GetCompetitionRegistrations 查看某竞赛所有报名（管理员）
func (c *CompetitionController) GetCompetitionRegistrations(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "竞赛ID格式错误",
		})
		return
	}

	// 检查竞赛是否存在
	var competition models.Competition
	if err := c.db.First(&competition, id).Error; err != nil {
		log.Printf("竞赛不存在: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "竞赛不存在",
		})
		return
	}

	var registrations []models.CompetitionRegistration
	var total int64

	// 获取查询参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	status := ctx.Query("status")

	query := c.db.Model(&models.CompetitionRegistration{}).Where("competition_id = ?", id).Preload("Competition")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		log.Printf("获取报名记录总数失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取报名记录失败",
		})
		return
	}

	// 分页查询
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("registration_time DESC").Find(&registrations).Error; err != nil {
		log.Printf("获取报名记录失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取报名记录失败",
		})
		return
	}

	// 构建响应数据
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

		responses = append(responses, response)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取报名记录成功",
		"data": gin.H{
			"list":  responses,
			"total": total,
			"page":  page,
			"size":  size,
		},
	})
}

// GetCompetitionSubmissions 查看竞赛提交作品（教师/管理员）
func (c *CompetitionController) GetCompetitionSubmissions(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "竞赛ID格式错误",
		})
		return
	}

	// 检查竞赛是否存在
	var competition models.Competition
	if err := c.db.First(&competition, id).Error; err != nil {
		log.Printf("竞赛不存在: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "竞赛不存在",
		})
		return
	}

	var submissions []models.CompetitionSubmission
	var total int64

	// 获取查询参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	status := ctx.Query("status")

	query := c.db.Model(&models.CompetitionSubmission{}).Where("competition_id = ?", id).Preload("Student.Profile")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		log.Printf("获取提交作品总数失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取提交作品失败",
		})
		return
	}

	// 分页查询
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("submitted_at DESC").Find(&submissions).Error; err != nil {
		log.Printf("获取提交作品失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取提交作品失败",
		})
		return
	}

	// 构建响应数据
	var responses []models.CompetitionSubmissionResponse
	for _, sub := range submissions {
		response := models.CompetitionSubmissionResponse{
			ID:            sub.ID,
			CompetitionID: sub.CompetitionID,
			StudentID:     sub.StudentID,
			TeacherViewed: sub.TeacherViewed,
			SubmitTime:    sub.SubmittedAt,
			FileURL:       sub.FileURL,
			Version:       sub.Version,
			FileSize:      sub.FileSize,
			Description:   sub.Description,
		}

		if sub.Student != nil {
			response.Student = &models.Users{
				ID:         sub.Student.ID,
				Username:   sub.Student.Username,
				Email:      sub.Student.Email,
				Status:     sub.Student.Status,
				CreateTime: sub.Student.CreateTime,
			}
			// 如果有用户档案信息，添加详细信息
			if sub.Student.Profile != nil {
				response.Student.RealName = sub.Student.Profile.RealName
				response.Student.Phone = sub.Student.Profile.Phone
				response.Student.Department = sub.Student.Profile.Department
			}
		}

		responses = append(responses, response)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取提交作品成功",
		"data": gin.H{
			"list":  responses,
			"total": total,
			"page":  page,
			"size":  size,
		},
	})
}

// ScoreCompetition 评审提交作品（教师）
func (c *CompetitionController) ScoreCompetition(ctx *gin.Context) {

	var req models.CompetitionScoreRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "参数错误",
		})
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, models.Response{
			Code:    401,
			Message: "未登录",
		})
		return
	}

	judgeID := userID.(uint)

	var submission models.CompetitionSubmission
	if err := c.db.First(&submission, req.SubmissionID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "提交记录不存在",
		})
		return
	}

	score := models.CompetitionScore{
		JudgeID:       judgeID,
		SubmissionID:  submission.ID,
		CompetitionID: submission.CompetitionID,
		StudentID:     submission.StudentID,
		Score:         req.Score,
		Comment:       req.Comment,
	}

	if err := c.db.Create(&score).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "评分失败",
		})
		return
	}

	// 构造返回数据
	resp := models.CompetitionScoreResponse{
		ID:            score.ID,
		JudgeID:       score.JudgeID,
		SubmissionID:  score.SubmissionID,
		CompetitionID: score.CompetitionID,
		StudentID:     score.StudentID,
		Score:         score.Score,
		Comment:       score.Comment,
		ScoredAt:      score.ScoredAt,
	}

	ctx.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "评分成功",
		Data:    resp,
	})
}

// SubmitResult 登记成绩/获奖信息（管理员）
func (c *CompetitionController) SubmitResult(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "竞赛ID格式错误",
		})
		return
	}

	// 获取当前用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户未登录",
		})
		return
	}

	var req struct {
		StudentID  uint   `json:"student_id" binding:"required"`
		AwardLevel string `json:"award_level" binding:"required"`
		FinalScore *int   `json:"final_score"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("参数绑定失败: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 检查竞赛是否存在
	var competition models.Competition
	if err := c.db.First(&competition, id).Error; err != nil {
		log.Printf("竞赛不存在: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "竞赛不存在",
		})
		return
	}

	// 检查学生是否提交了作品
	var submission models.CompetitionSubmission
	if err := c.db.Where("competition_id = ? AND student_id = ?", id, req.StudentID).First(&submission).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "该学生尚未提交作品",
		})
		return
	}

	// 检查是否已有结果记录
	var existingResult models.CompetitionResult
	if err := c.db.Where("competition_id = ? AND student_id = ?", id, req.StudentID).First(&existingResult).Error; err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "该学生已有成绩记录",
		})
		return
	}

	// 创建结果记录
	result := models.CompetitionResult{
		CompetitionID: uint(id),
		StudentID:     req.StudentID,
		SubmissionID:  submission.ID,
		AwardLevel:    req.AwardLevel,
		FinalScore:    req.FinalScore,
		CreatedBy:     userID.(uint),
	}

	if err := c.db.Create(&result).Error; err != nil {
		log.Printf("登记成绩失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "登记成绩失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登记成绩成功",
		"data":    result,
	})
}

// ExportCompetitionData 导出竞赛报名或结果数据（管理员）
func (c *CompetitionController) ExportCompetitionData(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "竞赛ID格式错误",
		})
		return
	}

	exportType := ctx.Query("type")
	if exportType != "registrations" && exportType != "results" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "导出类型必须是 registrations 或 results",
		})
		return
	}

	// 检查竞赛是否存在
	var competition models.Competition
	if err := c.db.First(&competition, id).Error; err != nil {
		log.Printf("竞赛不存在: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "竞赛不存在",
		})
		return
	}

	var data []gin.H

	if exportType == "registrations" {
		// 导出报名数据
		var registrations []models.CompetitionRegistration
		if err := c.db.Where("competition_id = ?", id).Preload("Student").Find(&registrations).Error; err != nil {
			log.Printf("获取报名数据失败: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "导出报名数据失败",
			})
			return
		}

		for _, reg := range registrations {
			item := gin.H{
				"id":            reg.ID,
				"register_time": reg.RegistrationTime,
				"status":        reg.Status,
				"team_name":     reg.TeamName,
				"team_leader":   reg.TeamLeader,
			}
			data = append(data, item)
		}
	} else {
		// 导出结果数据
		var results []models.CompetitionResult
		if err := c.db.Where("competition_id = ?", id).Preload("Student").Find(&results).Error; err != nil {
			log.Printf("获取结果数据失败: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "导出结果数据失败",
			})
			return
		}

		for _, result := range results {
			item := gin.H{
				"id":           result.ID,
				"student_id":   result.StudentID,
				"award_level":  result.AwardLevel,
				"final_score":  result.FinalScore,
				"publish_time": result.PublishTime,
			}

			if result.Student != nil {
				item["student_name"] = result.Student.Username
				item["student_email"] = result.Student.Email
			}

			data = append(data, item)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "导出数据成功",
		"data": gin.H{
			"competition_id":    id,
			"competition_title": competition.Title,
			"export_type":       exportType,
			"total_count":       len(data),
			"data":              data,
		},
	})
}

// ==================== 竞赛评审相关方法 ====================

// AssignJudge 分配评审教师
func (c *CompetitionController) AssignJudge(ctx *gin.Context) {

	var request models.CompetitionJudgeRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	// 检查竞赛是否存在
	var competition models.Competition
	if err := c.db.First(&competition, request.CompetitionID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "竞赛不存在",
		})
		return
	}

	// 检查教师是否存在且具有教师角色
	var teacher models.User
	if err := c.db.Preload("Roles").First(&teacher, request.TeacherID).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "教师不存在",
		})
		return
	}

	// 验证教师角色
	isTeacher := false
	for _, role := range teacher.Roles {
		if role.RoleKey == "teacher" {
			isTeacher = true
			break
		}
	}
	if !isTeacher {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "该用户不是教师",
		})
		return
	}

	// 检查是否已经分配过
	var existingJudge models.CompetitionJudge
	if err := c.db.Where("competition_id = ? AND teacher_id = ?", request.CompetitionID, request.TeacherID).First(&existingJudge).Error; err == nil {
		// 更新状态
		existingJudge.Status = request.Status
		if err := c.db.Save(&existingJudge).Error; err != nil {
			log.Printf("更新评审教师状态失败: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "更新评审教师状态失败",
			})
			return
		}
	} else {
		// 创建新的分配记录
		judge := models.CompetitionJudge{
			CompetitionID: uint(request.CompetitionID),
			TeacherID:     request.TeacherID,
			Status:        request.Status,
		}
		if err := c.db.Create(&judge).Error; err != nil {
			log.Printf("分配评审教师失败: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "分配评审教师失败",
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "分配评审教师成功",
	})
}

// GetCompetitionJudges 获取竞赛评审教师列表
func (c *CompetitionController) GetCompetitionJudges(ctx *gin.Context) {
	competitionID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的竞赛ID",
		})
		return
	}

	var judges []models.CompetitionJudge
	if err := c.db.Where("competition_id = ?", competitionID).
		Preload("Teacher.Profile").
		Find(&judges).Error; err != nil {
		log.Printf("获取评审教师列表失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取评审教师列表失败",
		})
		return
	}

	var responses []models.CompetitionJudgeResponse
	for _, judge := range judges {
		response := models.CompetitionJudgeResponse{
			ID:            judge.ID,
			CompetitionID: judge.CompetitionID,
			AssignedAt:    judge.AssignedAt,
			Status:        judge.Status,
		}

		if judge.Teacher != nil {
			response.Teacher = &models.Users{
				ID:         judge.Teacher.ID,
				Username:   judge.Teacher.Username,
				Email:      judge.Teacher.Email,
				Status:     judge.Teacher.Status,
				CreateTime: judge.Teacher.CreateTime,
			}
			if judge.Teacher.Profile != nil {
				response.Teacher.RealName = judge.Teacher.Profile.RealName
				response.Teacher.Phone = judge.Teacher.Profile.Phone
				response.Teacher.Department = judge.Teacher.Profile.Department
			}
		}

		responses = append(responses, response)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取评审教师列表成功",
		"data":    responses,
	})
}

// SubmitScore 提交评分
func (c *CompetitionController) SubmitScore(ctx *gin.Context) {
	submissionID, err := strconv.ParseUint(ctx.Param("submissionId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的提交ID",
		})
		return
	}

	var request models.CompetitionScoreRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	// 获取当前用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	// 检查提交是否存在
	var submission models.CompetitionSubmission
	if err := c.db.First(&submission, submissionID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "提交记录不存在",
		})
		return
	}

	// 检查用户是否为该竞赛的评审教师
	var judge models.CompetitionJudge
	if err := c.db.Where("competition_id = ? AND teacher_id = ? AND status = ?",
		submission.CompetitionID, userID, "active").First(&judge).Error; err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "您不是该竞赛的评审教师",
		})
		return
	}

	// 检查是否已经评分过
	var existingScore models.CompetitionScore
	if err := c.db.Where("submission_id = ? AND judge_id = ?", submissionID, userID).First(&existingScore).Error; err == nil {
		// 更新评分
		existingScore.Score = request.Score
		existingScore.Comment = request.Comment
		if err := c.db.Save(&existingScore).Error; err != nil {
			log.Printf("更新评分失败: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "更新评分失败",
			})
			return
		}
	} else {
		// 创建新的评分记录
		score := models.CompetitionScore{
			SubmissionID:  uint(submissionID),
			JudgeID:       userID.(uint),
			StudentID:     submission.StudentID,
			CompetitionID: submission.CompetitionID,
			Score:         request.Score,
			Comment:       request.Comment,
		}
		if err := c.db.Create(&score).Error; err != nil {
			log.Printf("提交评分失败: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "提交评分失败",
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "提交评分成功",
	})
}

// GetSubmissionScores 获取提交的评分列表
func (c *CompetitionController) GetSubmissionScores(ctx *gin.Context) {
	submissionID, err := strconv.ParseUint(ctx.Param("submissionId"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的提交ID",
		})
		return
	}

	var scores []models.CompetitionScore
	if err := c.db.Where("submission_id = ?", submissionID).
		Preload("Judge.Profile").
		Find(&scores).Error; err != nil {
		log.Printf("获取评分列表失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取评分列表失败",
		})
		return
	}

	var responses []models.CompetitionScoreResponse
	for _, score := range scores {
		response := models.CompetitionScoreResponse{
			ID:           score.ID,
			SubmissionID: score.SubmissionID,
			JudgeID:      score.JudgeID,
			Score:        score.Score,
			Comment:      score.Comment,
			ScoredAt:     score.ScoredAt,
		}

		if score.Judge != nil {
			response.Judge = &models.Users{
				ID:         score.Judge.ID,
				Username:   score.Judge.Username,
				Email:      score.Judge.Email,
				Status:     score.Judge.Status,
				CreateTime: score.Judge.CreateTime,
			}
			if score.Judge.Profile != nil {
				response.Judge.RealName = score.Judge.Profile.RealName
				response.Judge.Phone = score.Judge.Profile.Phone
				response.Judge.Department = score.Judge.Profile.Department
			}
		}

		responses = append(responses, response)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取评分列表成功",
		"data":    responses,
	})
}

// GetJudgingProgress 获取评审进度
func (c *CompetitionController) GetJudgingProgress(ctx *gin.Context) {
	competitionID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的竞赛ID",
		})
		return
	}

	// 获取评审教师数量
	var totalJudges, activeJudges int64
	c.db.Model(&models.CompetitionJudge{}).Where("competition_id = ?", competitionID).Count(&totalJudges)
	c.db.Model(&models.CompetitionJudge{}).Where("competition_id = ? AND status = ?", competitionID, "active").Count(&activeJudges)

	// 获取提交数量
	var totalSubmissions int64
	c.db.Model(&models.CompetitionSubmission{}).Where("competition_id = ?", competitionID).Count(&totalSubmissions)

	// 获取已评分的提交数量
	var scoredSubmissions int64
	c.db.Model(&models.CompetitionScore{}).
		Joins("JOIN competition_submissions ON competition_scores.submission_id = competition_submissions.id").
		Where("competition_submissions.competition_id = ?", competitionID).
		Distinct("competition_scores.submission_id").
		Count(&scoredSubmissions)

	// 计算进度百分比
	var progressPercentage float64
	if totalSubmissions > 0 {
		progressPercentage = float64(scoredSubmissions) / float64(totalSubmissions) * 100
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取评审进度成功",
		"data": gin.H{
			"total_judges":        totalJudges,
			"active_judges":       activeJudges,
			"total_submissions":   totalSubmissions,
			"scored_submissions":  scoredSubmissions,
			"progress_percentage": progressPercentage,
		},
	})
}

// FinalizeResults 最终确认成绩
func (c *CompetitionController) FinalizeResults(ctx *gin.Context) {
	competitionID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的竞赛ID",
		})
		return
	}

	// 获取当前用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	// 检查竞赛是否存在
	var competition models.Competition
	if err := c.db.First(&competition, competitionID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "竞赛不存在",
		})
		return
	}

	// 更新所有结果记录的最终确认信息
	now := time.Now()
	if err := c.db.Model(&models.CompetitionResult{}).
		Where("competition_id = ? AND finalized_by IS NULL", competitionID).
		Updates(map[string]interface{}{
			"finalized_by": userID,
			"finalized_at": now,
		}).Error; err != nil {
		log.Printf("最终确认成绩失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "最终确认成绩失败",
		})
		return
	}

	// 更新竞赛状态为已完成
	if err := c.db.Model(&competition).Update("status", "completed").Error; err != nil {
		log.Printf("更新竞赛状态失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新竞赛状态失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "最终确认成绩成功",
	})
}
