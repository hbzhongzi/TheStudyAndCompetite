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

// TeacherCompetitionController 教师竞赛控制器
type TeacherCompetitionController struct {
	DB *gorm.DB
}

// NewTeacherCompetitionController 创建教师竞赛控制器
func NewTeacherCompetitionController(db *gorm.DB) *TeacherCompetitionController {
	return &TeacherCompetitionController{DB: db}
}

// GetStudentRegistrations 获取指导学生报名记录
// @Summary 获取指导学生报名记录
// @Description 按竞赛名称/学生/状态筛选指导学生的报名记录
// @Tags 教师竞赛
// @Accept json
// @Produce json
// @Param competition_name query string false "竞赛名称筛选"
// @Param student_name query string false "学生姓名筛选"
// @Param status query string false "报名状态筛选"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} utils.Response{data=[]models.CompetitionRegistrationResponse}
// @Router /api/teacher/competition-registrations [get]
func (tc *TeacherCompetitionController) GetStudentRegistrations(c *gin.Context) {
	teacherID := utils.GetCurrentUserID(c)
	if teacherID == 0 {
		utils.ResponseError(c, http.StatusUnauthorized, "用户未登录", nil)
		return
	}

	var registrations []models.CompetitionRegistration
	query := tc.DB.Model(&models.CompetitionRegistration{}).
		Where("teacher_id = ?", teacherID).
		Preload("Competition.CreatedByUser").
		Preload("Student").
		Preload("Teacher")

	// 竞赛名称筛选
	if competitionName := c.Query("competition_name"); competitionName != "" {
		query = query.Joins("JOIN competitions ON competition_registrations.competition_id = competitions.id").
			Where("competitions.title LIKE ?", "%"+competitionName+"%")
	}

	// 学生姓名筛选
	if studentName := c.Query("student_name"); studentName != "" {
		query = query.Joins("JOIN users ON competition_registrations.student_id = users.id").
			Joins("JOIN user_profiles ON users.id = user_profiles.user_id").
			Where("user_profiles.real_name LIKE ?", "%"+studentName+"%")
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

// ReviewRegistration 审核报名
// @Summary 审核报名
// @Description 审核通过/驳回学生的报名申请，附带意见
// @Tags 教师竞赛
// @Accept json
// @Produce json
// @Param id path int true "报名记录ID"
// @Param request body map[string]interface{} true "审核信息"
// @Success 200 {object} utils.Response
// @Router /api/teacher/competition-registrations/{id}/review [put]
func (tc *TeacherCompetitionController) ReviewRegistration(c *gin.Context) {
	registrationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "无效的报名记录ID", err)
		return
	}

	teacherID := utils.GetCurrentUserID(c)
	if teacherID == 0 {
		utils.ResponseError(c, http.StatusUnauthorized, "用户未登录", nil)
		return
	}

	var request struct {
		Status  string `json:"status" binding:"required,oneof=approved rejected"`
		Comment string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "请求参数错误", err)
		return
	}

	// 检查报名记录是否存在且属于当前教师
	var registration models.CompetitionRegistration
	if err := tc.DB.Where("id = ? AND teacher_id = ?", registrationID, teacherID).First(&registration).Error; err != nil {
		utils.ResponseError(c, http.StatusNotFound, "报名记录不存在", err)
		return
	}

	// 检查是否可以审核
	if registration.TeacherReviewStatus != "pending" {
		utils.ResponseError(c, http.StatusBadRequest, "当前状态不允许审核", nil)
		return
	}

	// 更新审核状态
	updates := map[string]interface{}{
		"teacher_review_status":  request.Status,
		"teacher_review_comment": request.Comment,
		"teacher_review_time":    time.Now(),
	}

	// 如果审核通过，更新报名状态
	if request.Status == "approved" {
		updates["status"] = "approved"
	} else if request.Status == "rejected" {
		updates["status"] = "rejected"
	}

	if err := tc.DB.Model(&registration).Updates(updates).Error; err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "审核失败", err)
		return
	}

	utils.ResponseSuccess(c, gin.H{
		"message": "审核成功",
		"status":  request.Status,
	})
}

// ViewStudentSubmission 查看学生作品
// @Summary 查看学生作品
// @Description 查看学生上传的文件及时间
// @Tags 教师竞赛
// @Accept json
// @Produce json
// @Param id path int true "提交记录ID"
// @Success 200 {object} utils.Response{data=models.CompetitionSubmissionResponse}
// @Router /api/teacher/competition-submissions/{id} [get]
func (tc *TeacherCompetitionController) ViewStudentSubmission(c *gin.Context) {
	submissionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "无效的提交记录ID", err)
		return
	}

	teacherID := utils.GetCurrentUserID(c)
	if teacherID == 0 {
		utils.ResponseError(c, http.StatusUnauthorized, "用户未登录", nil)
		return
	}

	// 检查提交记录是否存在且属于当前教师指导的学生
	var submission models.CompetitionSubmission
	if err := tc.DB.Joins("JOIN competition_registrations ON competition_submissions.competition_id = competition_registrations.competition_id AND competition_submissions.student_id = competition_registrations.student_id").
		Where("competition_submissions.id = ? AND competition_registrations.teacher_id = ?", submissionID, teacherID).
		Preload("Competition").
		Preload("Student").
		First(&submission).Error; err != nil {
		utils.ResponseError(c, http.StatusNotFound, "提交记录不存在", err)
		return
	}

	// 记录查看时间
	tc.DB.Model(&submission).Updates(map[string]interface{}{
		"teacher_viewed": true,
		"last_view_time": time.Now(),
	})

	// 转换为响应格式
	response := models.CompetitionSubmissionResponse{
		ID:              submission.ID,
		CompetitionID:   submission.CompetitionID,
		StudentID:       submission.StudentID,
		FileURL:         submission.FileURL,
		FileName:        submission.FileName,
		FileSize:        submission.FileSize,
		Description:     submission.Description,
		Version:         submission.Version,
		SubmitTime:      submission.SubmitTime,
		Status:          submission.Status,
		ReviewComments:  submission.ReviewComments,
		Locked:          submission.Locked,
		TeacherViewed:   submission.TeacherViewed,
		TeacherFeedback: submission.TeacherFeedback,
		LastViewTime:    submission.LastViewTime,
	}

	if submission.Competition != nil {
		response.Competition = &models.CompetitionResponse{
			ID:     submission.Competition.ID,
			Title:  submission.Competition.Title,
			Type:   submission.Competition.Type,
			Status: submission.Competition.Status,
		}
	}

	if submission.Student != nil {
		response.Student = &models.CompetitionUserResponse{
			ID:         submission.Student.ID,
			Username:   submission.Student.Username,
			Email:      submission.Student.Email,
			Status:     submission.Student.Status,
			CreateTime: submission.Student.CreateTime,
		}
	}

	utils.ResponseSuccess(c, response)
}

// SubmitFeedback 提交评审意见
// @Summary 提交评审意见
// @Description 给出评分和评语
// @Tags 教师竞赛
// @Accept json
// @Produce json
// @Param id path int true "提交记录ID"
// @Param request body models.CompetitionFeedbackRequest true "评审意见"
// @Success 200 {object} utils.Response{data=models.CompetitionFeedbackResponse}
// @Router /api/teacher/competition-feedback/{id} [post]
func (tc *TeacherCompetitionController) SubmitFeedback(c *gin.Context) {
	submissionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "无效的提交记录ID", err)
		return
	}

	teacherID := utils.GetCurrentUserID(c)
	if teacherID == 0 {
		utils.ResponseError(c, http.StatusUnauthorized, "用户未登录", nil)
		return
	}

	var request models.CompetitionFeedbackRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "请求参数错误", err)
		return
	}

	// 检查提交记录是否存在且属于当前教师指导的学生
	var submission models.CompetitionSubmission
	if err := tc.DB.Joins("JOIN competition_registrations ON competition_submissions.competition_id = competition_registrations.competition_id AND competition_submissions.student_id = competition_registrations.student_id").
		Where("competition_submissions.id = ? AND competition_registrations.teacher_id = ?", submissionID, teacherID).
		First(&submission).Error; err != nil {
		utils.ResponseError(c, http.StatusNotFound, "提交记录不存在", err)
		return
	}

	// 创建评审记录
	feedback := models.CompetitionFeedback{
		CompetitionID: submission.CompetitionID,
		StudentID:     submission.StudentID,
		TeacherID:     teacherID,
		ReviewerID:    &teacherID,
		SubmissionID:  uint(submissionID),
		Comment:       request.Comment,
		Score:         request.Score,
		IsFinal:       request.IsFinal,
		FeedbackTime:  time.Now(),
	}

	if err := tc.DB.Create(&feedback).Error; err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "提交评审意见失败", err)
		return
	}

	// 如果是最终评审，更新提交状态
	if request.IsFinal {
		tc.DB.Model(&submission).Update("status", "reviewing")
	}

	utils.ResponseSuccess(c, feedback)
}

// GetFeedbackHistory 查看评审历史
// @Summary 查看评审历史
// @Description 查看自己评审过的记录
// @Tags 教师竞赛
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} utils.Response{data=[]models.CompetitionFeedbackResponse}
// @Router /api/teacher/competition-feedback/history [get]
func (tc *TeacherCompetitionController) GetFeedbackHistory(c *gin.Context) {
	teacherID := utils.GetCurrentUserID(c)
	if teacherID == 0 {
		utils.ResponseError(c, http.StatusUnauthorized, "用户未登录", nil)
		return
	}

	var feedbacks []models.CompetitionFeedback
	query := tc.DB.Model(&models.CompetitionFeedback{}).
		Where("reviewer_id = ?", teacherID).
		Preload("Competition").
		Preload("Student").
		Preload("Submission")

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	offset := (page - 1) * size

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(size).Order("feedback_time DESC").Find(&feedbacks).Error; err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "获取评审历史失败", err)
		return
	}

	// 转换为响应格式
	var responses []models.CompetitionFeedbackResponse
	for _, fb := range feedbacks {
		response := models.CompetitionFeedbackResponse{
			ID:            fb.ID,
			CompetitionID: fb.CompetitionID,
			StudentID:     fb.StudentID,
			TeacherID:     fb.TeacherID,
			ReviewerID:    fb.ReviewerID,
			SubmissionID:  fb.SubmissionID,
			Comment:       fb.Comment,
			Score:         fb.Score,
			FeedbackTime:  fb.FeedbackTime,
			IsFinal:       fb.IsFinal,
		}

		// 设置关联数据
		if fb.Competition != nil {
			response.Competition = &models.CompetitionResponse{
				ID:     fb.Competition.ID,
				Title:  fb.Competition.Title,
				Type:   fb.Competition.Type,
				Status: fb.Competition.Status,
			}
		}

		if fb.Student != nil {
			response.Student = &models.CompetitionUserResponse{
				ID:         fb.Student.ID,
				Username:   fb.Student.Username,
				Email:      fb.Student.Email,
				Status:     fb.Student.Status,
				CreateTime: fb.Student.CreateTime,
			}
		}

		if fb.Submission != nil {
			response.Submission = &models.CompetitionSubmissionResponse{
				ID:             fb.Submission.ID,
				CompetitionID:  fb.Submission.CompetitionID,
				StudentID:      fb.Submission.StudentID,
				FileURL:        fb.Submission.FileURL,
				FileName:       fb.Submission.FileName,
				FileSize:       fb.Submission.FileSize,
				Description:    fb.Submission.Description,
				Version:        fb.Submission.Version,
				SubmitTime:     fb.Submission.SubmitTime,
				Status:         fb.Submission.Status,
				ReviewComments: fb.Submission.ReviewComments,
				Locked:         fb.Submission.Locked,
			}
		}

		responses = append(responses, response)
	}

	utils.ResponseSuccess(c, gin.H{
		"feedbacks": responses,
		"total":     total,
		"page":      page,
		"size":      size,
	})
}

// GetStudentResults 查看指导学生成绩
// @Summary 查看指导学生成绩
// @Description 查看指导学生的成绩及获奖情况
// @Tags 教师竞赛
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} utils.Response{data=[]models.CompetitionResultResponse}
// @Router /api/teacher/competition-results [get]
func (tc *TeacherCompetitionController) GetStudentResults(c *gin.Context) {
	teacherID := utils.GetCurrentUserID(c)
	if teacherID == 0 {
		utils.ResponseError(c, http.StatusUnauthorized, "用户未登录", nil)
		return
	}

	var results []models.CompetitionResult
	query := tc.DB.Model(&models.CompetitionResult{}).
		Joins("JOIN competition_registrations ON competition_results.competition_id = competition_registrations.competition_id AND competition_results.student_id = competition_registrations.student_id").
		Where("competition_registrations.teacher_id = ?", teacherID).
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
