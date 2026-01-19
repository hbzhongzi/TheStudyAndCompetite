package controllers

import (
	"net/http"
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
			Type:                comp.Type,
			Organizer:           comp.Organizer,
			StartTime:           comp.StartTime,
			EndTime:             comp.EndTime,
			Description:         comp.Description,
			Attachment:          comp.Attachment,
			IsOpen:              comp.IsOpen,
			MaxParticipants:     comp.MaxParticipants,
			CurrentParticipants: comp.CurrentParticipants,
			Status:              comp.Status,
			AwardConfig:         comp.AwardConfig,
			CreatedBy:           comp.CreatedBy,
			CreatedAt:           comp.CreatedAt,
			UpdatedAt:           comp.UpdatedAt,
		}

		if comp.CreatedByUser != nil {
			response.CreatedByUser = &models.CompetitionUserResponse{
				ID:         comp.CreatedByUser.ID,
				Username:   comp.CreatedByUser.Username,
				Email:      comp.CreatedByUser.Email,
				Status:     comp.CreatedByUser.Status,
				CreateTime: comp.CreatedByUser.CreateTime,
			}
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

// RegisterCompetition 报名竞赛
// @Summary 报名竞赛
// @Description 学生报名参加竞赛，需要选择指导老师
// @Tags 学生竞赛
// @Accept json
// @Produce json
// @Param id path int true "竞赛ID"
// @Param request body models.CompetitionRegistrationRequest true "报名信息"
// @Success 200 {object} utils.Response{data=models.CompetitionRegistrationResponse}
// @Router /api/competitions/{id}/register [post]
func (cc *StudentCompetitionController) RegisterCompetition(c *gin.Context) {
	competitionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "无效的竞赛ID", err)
		return
	}

	// 获取当前用户ID
	userID := utils.GetCurrentUserID(c)
	if userID == 0 {
		utils.ResponseError(c, http.StatusUnauthorized, "用户未登录", nil)
		return
	}

	var request models.CompetitionRegistrationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "请求参数错误", err)
		return
	}

	// 检查竞赛是否存在且开放报名
	var competition models.Competition
	if err := cc.DB.First(&competition, competitionID).Error; err != nil {
		utils.ResponseError(c, http.StatusNotFound, "竞赛不存在", nil)
		return
	}

	if !competition.IsOpen {
		utils.ResponseError(c, http.StatusBadRequest, "竞赛未开放报名", nil)
		return
	}

	// 检查报名时间是否在有效范围内
	now := time.Now()

	// 移除报名开始时间不能早于当前时间的限制，允许设置过去的报名时间
	// 如果设置了报名开始时间，检查是否已开始报名
	// if competition.RegistrationStart != nil && now.Before(*competition.RegistrationStart) {
	// 	utils.ResponseError(c, http.StatusBadRequest, "报名尚未开始，请等待报名开放", nil)
	// 	return
	// }

	// 如果设置了报名截止时间，检查是否已截止报名
	if competition.RegistrationEnd != nil && now.After(*competition.RegistrationEnd) {
		utils.ResponseError(c, http.StatusBadRequest, "报名已截止，无法报名", nil)
		return
	}

	// 如果设置了比赛开始时间，检查比赛是否已开始
	if competition.StartTime != nil && now.After(*competition.StartTime) {
		utils.ResponseError(c, http.StatusBadRequest, "比赛已开始，无法报名", nil)
		return
	}

	// 检查是否已经报名
	var existingRegistration models.CompetitionRegistration
	if err := cc.DB.Where("competition_id = ? AND student_id = ?", competitionID, userID).First(&existingRegistration).Error; err == nil {
		utils.ResponseError(c, http.StatusBadRequest, "您已经报名过此竞赛", nil)
		return
	}

	// 检查是否需要指导老师
	if competition.TeacherLimit && request.TeacherID == nil {
		utils.ResponseError(c, http.StatusBadRequest, "此竞赛需要选择指导老师", nil)
		return
	}

	// 创建报名记录
	registration := models.CompetitionRegistration{
		CompetitionID:       uint(competitionID),
		StudentID:           userID,
		TeacherID:           request.TeacherID,
		Status:              "registered",
		TeacherReviewStatus: "pending",
		TeamName:            request.TeamName,
		TeamLeader:          request.TeamLeader,
		ContactPhone:        request.ContactPhone,
		ContactEmail:        request.ContactEmail,
		AdditionalInfo:      request.AdditionalInfo,
		RegisterTime:        time.Now(),
	}

	if err := cc.DB.Create(&registration).Error; err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "报名失败", err)
		return
	}

	// 更新竞赛参与人数
	cc.DB.Model(&competition).Update("current_participants", gorm.Expr("current_participants + ?", 1))

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "报名成功",
		"data":    registration,
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
			ID:                   reg.ID,
			CompetitionID:        reg.CompetitionID,
			StudentID:            reg.StudentID,
			TeacherID:            reg.TeacherID,
			RegisterTime:         reg.RegisterTime,
			Status:               reg.Status,
			TeacherReviewStatus:  reg.TeacherReviewStatus,
			TeacherReviewComment: reg.TeacherReviewComment,
			TeacherReviewTime:    reg.TeacherReviewTime,
			TeamName:             reg.TeamName,
			TeamLeader:           reg.TeamLeader,
			ContactPhone:         reg.ContactPhone,
			ContactEmail:         reg.ContactEmail,
			AdditionalInfo:       reg.AdditionalInfo,
		}

		// 设置关联数据
		if reg.Competition != nil {
			response.Competition = &models.CompetitionResponse{
				ID:          reg.Competition.ID,
				Title:       reg.Competition.Title,
				Type:        reg.Competition.Type,
				Organizer:   reg.Competition.Organizer,
				StartTime:   reg.Competition.StartTime,
				EndTime:     reg.Competition.EndTime,
				Description: reg.Competition.Description,
				Status:      reg.Competition.Status,
			}
		}

		if reg.Student != nil {
			response.Student = &models.CompetitionUserResponse{
				ID:         reg.Student.ID,
				Username:   reg.Student.Username,
				Email:      reg.Student.Email,
				Status:     reg.Student.Status,
				CreateTime: reg.Student.CreateTime,
			}
		}

		if reg.Teacher != nil {
			response.Teacher = &models.CompetitionUserResponse{
				ID:         reg.Teacher.ID,
				Username:   reg.Teacher.Username,
				Email:      reg.Teacher.Email,
				Status:     reg.Teacher.Status,
				CreateTime: reg.Teacher.CreateTime,
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

// UpdateRegistration 更新报名信息
// @Summary 更新报名信息
// @Description 仅限未提交或被驳回状态的报名可以更新
// @Tags 学生竞赛
// @Accept json
// @Produce json
// @Param id path int true "报名记录ID"
// @Param request body models.CompetitionRegistrationRequest true "更新信息"
// @Success 200 {object} utils.Response{data=models.CompetitionRegistrationResponse}
// @Router /api/competition-registrations/{id} [put]
func (cc *StudentCompetitionController) UpdateRegistration(c *gin.Context) {
	registrationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "无效的报名记录ID", err)
		return
	}

	userID := utils.GetCurrentUserID(c)
	if userID == 0 {
		utils.ResponseError(c, http.StatusUnauthorized, "用户未登录", nil)
		return
	}

	var request models.CompetitionRegistrationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "请求参数错误", err)
		return
	}

	// 检查报名记录是否存在且属于当前用户
	var registration models.CompetitionRegistration
	if err := cc.DB.Where("id = ? AND student_id = ?", registrationID, userID).First(&registration).Error; err != nil {
		utils.ResponseError(c, http.StatusNotFound, "报名记录不存在", err)
		return
	}

	// 检查是否可以更新
	if registration.Status != "registered" && registration.Status != "rejected" {
		utils.ResponseError(c, http.StatusBadRequest, "当前状态不允许更新", nil)
		return
	}

	// 更新报名信息
	updates := map[string]interface{}{
		"team_name":       request.TeamName,
		"team_leader":     request.TeamLeader,
		"contact_phone":   request.ContactPhone,
		"contact_email":   request.ContactEmail,
		"additional_info": request.AdditionalInfo,
	}

	if request.TeacherID != nil {
		updates["teacher_id"] = request.TeacherID
		updates["teacher_review_status"] = "pending"
		updates["teacher_review_comment"] = ""
		updates["teacher_review_time"] = nil
	}

	if err := cc.DB.Model(&registration).Updates(updates).Error; err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "更新失败", err)
		return
	}

	utils.ResponseSuccess(c, registration)
}

// SubmitWork 提交竞赛作品
// @Summary 提交竞赛作品
// @Description 学生提交竞赛作品文件
// @Tags 学生竞赛
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "竞赛ID"
// @Param file formData file true "作品文件"
// @Param description formData string false "作品描述"
// @Success 200 {object} utils.Response{data=models.CompetitionSubmissionResponse}
// @Router /api/competition-submissions/{id} [post]
func (cc *StudentCompetitionController) SubmitWork(c *gin.Context) {
	competitionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "无效的竞赛ID", err)
		return
	}

	userID := utils.GetCurrentUserID(c)
	if userID == 0 {
		utils.ResponseError(c, http.StatusUnauthorized, "用户未登录", nil)
		return
	}

	// 检查是否已报名
	var registration models.CompetitionRegistration
	if err := cc.DB.Where("competition_id = ? AND student_id = ? AND status = ?", competitionID, userID, "approved").First(&registration).Error; err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "您尚未报名或报名未通过审核", err)
		return
	}

	// 处理文件上传
	file, err := c.FormFile("file")
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "文件上传失败", err)
		return
	}

	// 保存文件
	filename := utils.GenerateFileName(file.Filename)
	filepath := "uploads/competitions/" + filename
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "文件保存失败", err)
		return
	}

	description := c.PostForm("description")

	// 创建提交记录
	submission := models.CompetitionSubmission{
		CompetitionID: uint(competitionID),
		StudentID:     userID,
		FileURL:       filepath,
		FileName:      file.Filename,
		FileSize:      file.Size,
		Description:   description,
		Version:       1,
		Status:        "submitted",
		SubmitTime:    time.Now(),
	}

	if err := cc.DB.Create(&submission).Error; err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "提交失败", err)
		return
	}

	utils.ResponseSuccess(c, submission)
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
				Type:   result.Competition.Type,
				Status: result.Competition.Status,
			}
		}

		if result.Student != nil {
			response.Student = &models.CompetitionUserResponse{
				ID:         result.Student.ID,
				Username:   result.Student.Username,
				Email:      result.Student.Email,
				Status:     result.Student.Status,
				CreateTime: result.Student.CreateTime,
			}
		}

		if result.CreatedByUser != nil {
			response.CreatedByUser = &models.CompetitionUserResponse{
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
