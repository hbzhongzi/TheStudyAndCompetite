package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/datatypes"
)

// CompetitionUserResponse 竞赛模块用户响应（简化版）
type Users struct {
	ID         uint      `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Status     string    `json:"status"`
	RealName   string    `json:"realName"`
	Phone      string    `json:"phone"`
	Department string    `json:"department"`
	CreateTime time.Time `json:"createTime"`
}

// Competition 竞赛信息表
type Competition struct {
	ID                  uint           `gorm:"column:id;primaryKey" json:"id"`
	Title               string         `gorm:"column:title" json:"title"`
	Description         string         `gorm:"column:description" json:"description"`
	Level               string         `gorm:"column:level" json:"level"`
	Category            string         `gorm:"column:category" json:"category"`
	RegistrationStart   *time.Time     `gorm:"column:registration_start" json:"registrationStart"`
	RegistrationEnd     *time.Time     `gorm:"column:registration_end" json:"registrationEnd"`
	SubmissionStart     *time.Time     `gorm:"column:submission_start" json:"submissionStart"`
	SubmissionEnd       *time.Time     `gorm:"column:submission_end" json:"submissionEnd"`
	MaxParticipants     int            `gorm:"column:max_participants" json:"maxParticipants"`
	CurrentParticipants int            `gorm:"column:current_participants" json:"currentParticipants"`
	IsOpen              bool           `gorm:"column:is_open" json:"isOpen"`
	Status              string         `gorm:"column:status" json:"status"`
	AwardConfig         datatypes.JSON `gorm:"column:award_config" json:"awardConfig"`
	CreatedBy           uint           `gorm:"column:created_by" json:"createdBy"`
	CreatedAt           *time.Time     `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt           *time.Time     `gorm:"column:updated_at" json:"updatedAt"`
}

func (Competition) TableName() string {
	return "competitions"
}

// CompetitionRegistration 竞赛报名记录表
type CompetitionRegistration struct {
	ID                   uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	CompetitionID        uint       `json:"competition_id" gorm:"not null;comment:竞赛ID"`
	StudentID            uint       `json:"student_id" gorm:"not null;comment:学生ID"`
	TeacherID            *uint      `json:"teacher_id" gorm:"comment:指导教师ID（从 users.id 外键）"`
	RegisterTime         time.Time  `json:"registration_time" gorm:"autoCreateTime;comment:报名时间"`
	Status               string     `json:"status" gorm:"type:enum('registered','withdrawn','approved','rejected');default:registered;comment:报名状态"`
	TeacherReviewStatus  string     `json:"teacher_review_status" gorm:"type:enum('pending','approved','rejected');default:pending;comment:教师审核状态"`
	TeacherReviewComment string     `json:"teacher_review_comment" gorm:"type:text;comment:教师审核意见"`
	TeacherReviewTime    *time.Time `json:"teacher_review_time" gorm:"comment:教师审核时间"`
	TeamName             string     `json:"team_name" gorm:"type:varchar(100);comment:团队名称"`
	TeamLeader           int        `json:"team_leader" gorm:"default:0;comment:团队负责人id"`
	ContactPhone         string     `json:"contact_phone" gorm:"type:varchar(20);comment:联系电话"`
	ContactEmail         string     `json:"contact_email" gorm:"type:varchar(100);comment:联系邮箱"`
	AdditionalInfo       JSONMap    `json:"additional_info" gorm:"type:json;comment:额外信息"`

	// 关联关系
	Competition *Competition `json:"competition" gorm:"foreignKey:CompetitionID"`
	Student     *User        `json:"student" gorm:"foreignKey:StudentID"`
	Teacher     *User        `json:"teacher" gorm:"foreignKey:TeacherID"`
}

// TableName 指定表名
func (CompetitionRegistration) TableName() string {
	return "competition_registrations"
}

// CompetitionSubmission 竞赛成果提交表
type CompetitionSubmission struct {
	ID              uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	CompetitionID   uint       `json:"competition_id" gorm:"not null;comment:竞赛ID"`
	StudentID       uint       `json:"student_id" gorm:"not null;comment:学生ID"`
	FileURL         string     `json:"file_url" gorm:"type:varchar(255);comment:文件URL"`
	FileName        string     `json:"file_name" gorm:"type:varchar(100);comment:文件名"`
	FileSize        int64      `json:"file_size" gorm:"comment:文件大小"`
	Description     string     `json:"description" gorm:"type:text;comment:成果描述"`
	Version         int        `json:"version" gorm:"default:1;comment:版本号"`
	SubmitTime      time.Time  `json:"submit_time" gorm:"autoCreateTime;comment:提交时间"`
	Status          string     `json:"status" gorm:"type:enum('submitted','reviewing','approved','rejected');default:submitted;comment:提交状态"`
	ReviewComments  string     `json:"review_comments" gorm:"type:text;comment:评审意见"`
	Locked          bool       `json:"locked" gorm:"default:false;comment:是否锁定"`
	TeacherViewed   bool       `json:"teacher_viewed" gorm:"default:false;comment:教师是否查看过作品"`
	TeacherFeedback string     `json:"teacher_feedback" gorm:"type:text;comment:教师对作品的非正式反馈（不是评分）"`
	LastViewTime    *time.Time `json:"last_view_time" gorm:"comment:教师最后查看时间"`

	// 关联关系
	Competition *Competition          `json:"competition" gorm:"foreignKey:CompetitionID"`
	Student     *User                 `json:"student" gorm:"foreignKey:StudentID"`
	Feedback    []CompetitionFeedback `json:"feedback" gorm:"foreignKey:SubmissionID"`
	Results     []CompetitionResult   `json:"results" gorm:"foreignKey:SubmissionID"`
	Scores      []CompetitionScore    `json:"scores" gorm:"foreignKey:SubmissionID"`
}

// TableName 指定表名
func (CompetitionSubmission) TableName() string {
	return "competition_submissions"
}

// CompetitionFeedback 竞赛教师评语表
type CompetitionFeedback struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	CompetitionID uint      `json:"competition_id" gorm:"not null;comment:竞赛ID"`
	StudentID     uint      `json:"student_id" gorm:"not null;comment:学生ID"`
	TeacherID     uint      `json:"teacher_id" gorm:"not null;comment:教师ID"`
	ReviewerID    *uint     `json:"reviewer_id" gorm:"comment:评审教师ID（users.id 外键）"`
	SubmissionID  uint      `json:"submission_id" gorm:"not null;comment:提交记录ID"`
	Comment       string    `json:"comment" gorm:"type:text;comment:评语内容"`
	Score         *float64  `json:"score" gorm:"type:decimal(5,2);comment:评审分数"`
	FeedbackTime  time.Time `json:"feedback_time" gorm:"autoCreateTime;comment:评语时间"`
	IsFinal       bool      `json:"is_final" gorm:"default:false;comment:是否为最终评语"`

	// 关联关系
	Competition *Competition           `json:"competition" gorm:"foreignKey:CompetitionID"`
	Student     *User                  `json:"student" gorm:"foreignKey:StudentID"`
	Teacher     *User                  `json:"teacher" gorm:"foreignKey:TeacherID"`
	Reviewer    *User                  `json:"reviewer" gorm:"foreignKey:ReviewerID"`
	Submission  *CompetitionSubmission `json:"submission" gorm:"foreignKey:SubmissionID"`
}

// TableName 指定表名
func (CompetitionFeedback) TableName() string {
	return "competition_feedback"
}

// CompetitionJudge 竞赛评审教师表
type CompetitionJudge struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	CompetitionID uint      `json:"competition_id" gorm:"not null;comment:竞赛ID"`
	TeacherID     uint      `json:"teacher_id" gorm:"not null;comment:教师ID"`
	AssignedAt    time.Time `json:"assigned_at" gorm:"autoCreateTime;comment:分配时间"`
	Status        string    `json:"status" gorm:"type:enum('active','inactive');default:active;comment:是否参与"`

	// 关联关系
	Competition *Competition       `json:"competition" gorm:"foreignKey:CompetitionID"`
	Teacher     *User              `json:"teacher" gorm:"foreignKey:TeacherID"`
	Scores      []CompetitionScore `json:"scores" gorm:"foreignKey:JudgeID"`
}

// TableName 指定表名
func (CompetitionJudge) TableName() string {
	return "competition_judges"
}

// CompetitionScore 竞赛评分记录表
type CompetitionScore struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	SubmissionID uint      `json:"submission_id" gorm:"not null;comment:提交记录ID"`
	JudgeID      uint      `json:"judge_id" gorm:"not null;comment:评审教师ID"`
	Score        float64   `json:"score" gorm:"type:decimal(5,2);not null;comment:评分"`
	Comment      string    `json:"comment" gorm:"type:text;comment:评语"`
	ScoredAt     time.Time `json:"scored_at" gorm:"autoCreateTime;comment:评分时间"`

	// 关联关系
	Submission *CompetitionSubmission `json:"submission" gorm:"foreignKey:SubmissionID"`
	Judge      *User                  `json:"judge" gorm:"foreignKey:JudgeID"`
}

// TableName 指定表名
func (CompetitionScore) TableName() string {
	return "competition_scores"
}

// CompetitionResult 竞赛获奖登记表
type CompetitionResult struct {
	ID             uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	CompetitionID  uint       `json:"competition_id" gorm:"not null;comment:竞赛ID"`
	StudentID      uint       `json:"student_id" gorm:"not null;comment:学生ID"`
	SubmissionID   uint       `json:"submission_id" gorm:"not null;comment:提交记录ID"`
	AwardLevel     string     `json:"award_level" gorm:"type:varchar(50);comment:获奖等级"`
	FinalScore     *int       `json:"final_score" gorm:"comment:最终得分"`
	CertificateURL string     `json:"certificate_url" gorm:"type:varchar(255);comment:证书URL"`
	PublishTime    time.Time  `json:"publish_time" gorm:"autoCreateTime;comment:公布时间"`
	CreatedBy      uint       `json:"created_by" gorm:"not null;comment:创建者ID"`
	FinalizedBy    *uint      `json:"finalized_by" gorm:"comment:最终确认成绩的管理员ID"`
	FinalizedAt    *time.Time `json:"finalized_at" gorm:"comment:确认时间"`

	// 关联关系
	Competition     *Competition           `json:"competition" gorm:"foreignKey:CompetitionID"`
	Student         *User                  `json:"student" gorm:"foreignKey:StudentID"`
	Submission      *CompetitionSubmission `json:"submission" gorm:"foreignKey:SubmissionID"`
	CreatedByUser   *User                  `json:"created_by_user" gorm:"foreignKey:CreatedBy"`
	FinalizedByUser *User                  `json:"finalized_by_user" gorm:"foreignKey:FinalizedBy"`
}

// TableName 指定表名
func (CompetitionResult) TableName() string {
	return "competition_results"
}

// JSONMap 用于处理JSON字段
type JSONMap map[string]interface{}

// Value 实现driver.Valuer接口
func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan 实现sql.Scanner接口
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("cannot scan non-string value into JSONMap")
	}

	return json.Unmarshal(bytes, j)
}

type CompetitionCreateRequest struct {
	Title             string    `json:"title" binding:"required"`
	Type              string    `json:"type"`
	Organizer         string    `json:"organizer"`
	RegistrationStart time.Time `json:"registration_start"`
	RegistrationEnd   time.Time `json:"registration_end"`
	StartTime         time.Time `json:"start_time"`
	EndTime           time.Time `json:"end_time"`
	Level             string    `json:"level"`
	Description       string    `json:"description"`
	Attachment        string    `json:"attachment"`
	IsOpen            bool      `json:"is_open"`
	MaxParticipants   int       `json:"max_participants"`
}

// CompetitionUpdateRequest 更新竞赛请求
type CompetitionUpdateRequest struct {
	Title             string     `json:"title"`
	Type              string     `json:"type"`
	Organizer         string     `json:"organizer"`
	RegistrationStart time.Time  `json:"registration_start"`
	RegistrationEnd   time.Time  `json:"registration_end"`
	StartTime         *time.Time `json:"start_time"`
	EndTime           *time.Time `json:"end_time"`
	Description       string     `json:"description"`
	Attachment        string     `json:"attachment"`
	IsOpen            *bool      `json:"is_open"`
	MaxParticipants   *int       `json:"max_participants"`
	Status            string     `json:"status" binding:"omitempty,oneof=draft registration submission review completed"`
	AwardConfig       JSONMap    `json:"award_config"`
}

// CompetitionRegistrationRequest 竞赛报名请求
type CompetitionRegistrationRequest struct {
	TeamName   string `json:"team_name"`
	TeamLeader int    `json:"team_leader"`
}

// CompetitionSubmissionRequest 竞赛提交请求
type CompetitionSubmissionRequest struct {
	FileURL     string `json:"file_url" binding:"required"`
	FileName    string `json:"file_name" binding:"required"`
	FileSize    int64  `json:"file_size"`
	Description string `json:"description"`
}

// CompetitionFeedbackRequest 竞赛评语请求
type CompetitionFeedbackRequest struct {
	Comment string   `json:"comment" binding:"required"`
	Score   *float64 `json:"score"`
	IsFinal bool     `json:"is_final"`
}

// CompetitionResultRequest 竞赛获奖请求
type CompetitionResultRequest struct {
	AwardLevel     string `json:"award_level" binding:"required"`
	FinalScore     *int   `json:"final_score"`
	CertificateURL string `json:"certificate_url"`
}

// CompetitionJudgeRequest 竞赛评审教师分配请求
type CompetitionJudgeRequest struct {
	CompetitionID uint   `json:"competition_id" binding:"required"`
	TeacherID     uint   `json:"teacher_id" binding:"required"`
	Status        string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// CompetitionScoreRequest 竞赛评分请求
type CompetitionScoreRequest struct {
	Score   float64 `json:"score" binding:"required,min=0,max=100"`
	Comment string  `json:"comment"`
}

// CompetitionJudgeResponse 竞赛评审教师响应
type CompetitionJudgeResponse struct {
	ID            uint      `json:"id"`
	CompetitionID uint      `json:"competition_id"`
	TeacherID     uint      `json:"teacher_id"`
	AssignedAt    time.Time `json:"assigned_at"`
	Status        string    `json:"status"`

	// 关联数据
	Teacher *Users `json:"teacher"`
}

// CompetitionScoreResponse 竞赛评分响应
type CompetitionScoreResponse struct {
	ID           uint      `json:"id"`
	SubmissionID uint      `json:"submission_id"`
	JudgeID      uint      `json:"judge_id"`
	Score        float64   `json:"score"`
	Comment      string    `json:"comment"`
	ScoredAt     time.Time `json:"scored_at"`

	// 关联数据
	Judge      *Users                         `json:"judge"`
	Submission *CompetitionSubmissionResponse `json:"submission"`
}

// CompetitionResponse 竞赛响应
type CompetitionResponse struct {
	ID                  uint       `json:"id"`
	Title               string     `json:"title"`
	Type                string     `json:"type"`
	Organizer           string     `json:"organizer"`
	RegistrationStart   *time.Time `json:"registration_start"`
	RegistrationEnd     *time.Time `json:"registration_end"`
	StartTime           *time.Time `json:"start_time"`
	EndTime             *time.Time `json:"end_time"`
	Description         string     `json:"description"`
	Attachment          string     `json:"attachment"`
	IsOpen              bool       `json:"is_open"`
	MaxParticipants     *int       `json:"max_participants"`
	CurrentParticipants int        `json:"current_participants"`
	Status              string     `json:"status"`
	AwardConfig         JSONMap    `json:"award_config"`
	CreatedBy           uint       `json:"created_by"`
	CreatedAt           *time.Time `json:"created_at"`
	UpdatedAt           *time.Time `json:"updated_at"`

	// 关联数据
	CreatedByUser     *Users `json:"created_by_user"`
	RegistrationCount int    `json:"registration_count"`
	SubmissionCount   int    `json:"submission_count"`
	ResultCount       int    `json:"result_count"`
	JudgeCount        int    `json:"judge_count"`
}

// CompetitionRegistrationResponse 竞赛报名响应
type CompetitionRegistrationResponse struct {
	ID                   uint       `json:"id"`
	CompetitionID        uint       `json:"competition_id"`
	StudentID            uint       `json:"student_id"`
	TeacherID            *uint      `json:"teacher_id"`
	RegisterTime         time.Time  `json:"register_time"`
	Status               string     `json:"status"`
	TeacherReviewStatus  string     `json:"teacher_review_status"`
	TeacherReviewComment string     `json:"teacher_review_comment"`
	TeacherReviewTime    *time.Time `json:"teacher_review_time"`
	TeamName             string     `json:"team_name"`
	TeamLeader           int        `json:"team_leader"`
	ContactPhone         string     `json:"contact_phone"`
	ContactEmail         string     `json:"contact_email"`
	AdditionalInfo       JSONMap    `json:"additional_info"`

	// 关联数据
	Competition *CompetitionResponse `json:"competition"`
	Student     *Users               `json:"student"`
	Teacher     *Users               `json:"teacher"`
}

// CompetitionSubmissionResponse 竞赛提交响应
type CompetitionSubmissionResponse struct {
	ID              uint       `json:"id"`
	CompetitionID   uint       `json:"competition_id"`
	StudentID       uint       `json:"student_id"`
	FileURL         string     `json:"file_url"`
	FileName        string     `json:"file_name"`
	FileSize        int64      `json:"file_size"`
	Description     string     `json:"description"`
	Version         int        `json:"version"`
	SubmitTime      time.Time  `json:"submit_time"`
	Status          string     `json:"status"`
	ReviewComments  string     `json:"review_comments"`
	Locked          bool       `json:"locked"`
	TeacherViewed   bool       `json:"teacher_viewed"`
	TeacherFeedback string     `json:"teacher_feedback"`
	LastViewTime    *time.Time `json:"last_view_time"`

	// 关联数据
	Competition *CompetitionResponse          `json:"competition"`
	Student     *Users                        `json:"student"`
	Feedback    []CompetitionFeedbackResponse `json:"feedback"`
	Scores      []CompetitionScoreResponse    `json:"scores"`
}

// CompetitionFeedbackResponse 竞赛评语响应
type CompetitionFeedbackResponse struct {
	ID            uint      `json:"id"`
	CompetitionID uint      `json:"competition_id"`
	StudentID     uint      `json:"student_id"`
	TeacherID     uint      `json:"teacher_id"`
	ReviewerID    *uint     `json:"reviewer_id"`
	SubmissionID  uint      `json:"submission_id"`
	Comment       string    `json:"comment"`
	Score         *float64  `json:"score"`
	FeedbackTime  time.Time `json:"feedback_time"`
	IsFinal       bool      `json:"is_final"`

	// 关联数据
	Competition *CompetitionResponse           `json:"competition"`
	Student     *Users                         `json:"student"`
	Teacher     *Users                         `json:"teacher"`
	Reviewer    *Users                         `json:"reviewer"`
	Submission  *CompetitionSubmissionResponse `json:"submission"`
}

// CompetitionResultResponse 竞赛获奖响应
type CompetitionResultResponse struct {
	ID             uint       `json:"id"`
	CompetitionID  uint       `json:"competition_id"`
	StudentID      uint       `json:"student_id"`
	SubmissionID   uint       `json:"submission_id"`
	AwardLevel     string     `json:"award_level"`
	FinalScore     *int       `json:"final_score"`
	CertificateURL string     `json:"certificate_url"`
	PublishTime    time.Time  `json:"publish_time"`
	CreatedBy      uint       `json:"created_by"`
	FinalizedBy    *uint      `json:"finalized_by"`
	FinalizedAt    *time.Time `json:"finalized_at"`

	// 关联数据
	Competition     *CompetitionResponse           `json:"competition"`
	Student         *Users                         `json:"student"`
	Submission      *CompetitionSubmissionResponse `json:"submission"`
	CreatedByUser   *Users                         `json:"created_by_user"`
	FinalizedByUser *Users                         `json:"finalized_by_user"`
}
