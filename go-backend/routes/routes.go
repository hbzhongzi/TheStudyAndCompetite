package routes

import (
	"yunmeng-backend/controllers"
	"yunmeng-backend/middlewares"
	"yunmeng-backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/api")
	{
		// 健康检查路由（无需认证）
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"code":    200,
				"message": "服务运行正常",
				"status":  "healthy",
			})
		})

		// 登录路由（无需认证）
		api.POST("/login", controllers.LoginHandler(db))

		// 认证相关路由
		authController := controllers.NewAuthController(db)
		api.POST("/refresh-token", authController.RefreshToken)  // Token刷新
		api.GET("/validate-token", authController.ValidateToken) // Token验证

		// 需要认证的路由组
		auth := api.Group("")
		auth.Use(middlewares.AuthMiddleware())
		{
			// 用户管理路由（仅管理员）
			userService := services.NewUserService(db)
			userController := controllers.NewUserController(userService)

			users := auth.Group("/users")
			users.Use(middlewares.AdminOnly())
			{
				users.GET("", userController.GetUserList)                           // 获取用户列表
				users.GET("/:id", userController.GetUserByID)                       // 获取用户详情
				users.POST("", userController.CreateUser)                           // 创建用户
				users.PUT("/:id", userController.UpdateUser)                        // 更新用户
				users.DELETE("/:id", userController.DeleteUser)                     // 删除用户
				users.PATCH("/:id/status", userController.ToggleUserStatus)         // 切换用户状态
				users.POST("/:id/reset-password", userController.ResetUserPassword) // 重置密码
				users.POST("/batch-delete", userController.BatchDeleteUsers)        // 批量删除
				users.GET("/stats", userController.GetUserStats)                    // 获取统计信息
				users.GET("/export", userController.ExportUsers)                    // 导出用户数据
			}

			// 项目模块路由
			projectService := services.NewProjectService(db)
			projectController := controllers.NewProjectController(projectService)
			fileController := controllers.NewFileController()

			// 项目分类管理路由（仅管理员）
			projectTypeController := controllers.NewProjectTypeController(db)
			projectTypes := auth.Group("/project-types")
			projectTypes.Use(middlewares.AdminOnly())
			{
				projectTypes.GET("", projectTypeController.GetProjectTypeList)        // 获取项目分类列表
				projectTypes.GET("/stats", projectTypeController.GetProjectTypeStats) // 获取项目分类统计
				projectTypes.GET("/:id", projectTypeController.GetProjectTypeByID)    // 获取项目分类详情
				projectTypes.POST("", projectTypeController.CreateProjectType)        // 创建项目分类
				projectTypes.PUT("/:id", projectTypeController.UpdateProjectType)     // 更新项目分类
				projectTypes.DELETE("/:id", projectTypeController.DeleteProjectType)  // 删除项目分类
			}

			// 教师管理路由
			teachers := auth.Group("/teachers")
			teachers.Use(middlewares.RoleMiddleware("teacher", "admin"))
			{
				teachers.GET("", projectController.GetTeacherList)                                               // 获取教师列表
				teachers.POST("/ApproveExtensionApplication", projectController.ApproveExtensionApplication)     // 审批延期申请
				teachers.GET("/TeacherExtensionApplications", projectController.GetTeacherExtensionApplications) // 教师查看延期申请列表
				teachers.GET("/Stufiles", projectController.GetStudentProjectsFiles)                             // 获取指定学生的项目对应的文件列表

				teachers.GET("/filter", projectController.GetTeacherListWithFilter)                                 // 获取教师列表（支持院系筛选）
				teachers.POST("/bind", projectController.BindStudentTeacher)                                        // 绑定学生和教师
				teachers.GET("/students", projectController.GetMyStudents)                                          // 获取当前登录教师指导的学生
				teachers.GET("/students/:studentId", projectController.GetStudentTeachers)                          // 获取学生的指导教师
				teachers.DELETE("/students/:studentId/teachers/:teacherId", projectController.UnbindStudentTeacher) // 解绑学生和教师

				teachers.GET("/projects", projectController.GetTeacherProjects) // 获取当前登录教师的所有指导项目

			}

			// 教师/管理员项目路由
			teacherProjects := auth.Group("/teacher-projects")
			teacherProjects.Use(middlewares.RoleMiddleware("teacher", "admin"))
			{
				teacherProjects.GET("", projectController.GetProjectList)                  // 获取所有项目列表
				teacherProjects.POST("/review", projectController.ReviewProject)           // 审核项目
				teacherProjects.GET("/complete/list", projectController.GetMyCompleteList) // 查看分配的竞赛审核项目列表

				teacherProjects.GET("/:id/reviews", projectController.GetProjectReviews) // 获取审核记录

				// =============================================
				// 3. 成果文件管理增强路由（教师/管理员）
				// =============================================
				teacherProjects.PUT("/files/:fileId/review", projectController.ReviewProjectFile) // 审核项目文件

			}

			// 管理员项目管理路由
			adminProjects := auth.Group("/admin/projects")
			//adminProjects.Use(middlewares.AdminOnly())
			{
				adminProjects.PUT("/:id/force-status", projectController.ForceUpdateProjectStatus) // 强制更新项目状态

				adminProjects.GET("/stats", projectController.GetProjectStats)  // 获取项目统计
				adminProjects.POST("/export", projectController.ExportProjects) // 导出项目数据

				// =============================================
				// 4. 项目分类管理增强路由（管理员）
				// =============================================
				adminProjects.POST("/types", projectController.CreateProjectType)        // 创建项目分类
				adminProjects.PUT("/types/:id", projectController.UpdateProjectType)     // 更新项目分类
				adminProjects.GET("/types/tree", projectController.GetProjectTypeTree)   // 获取项目分类树
				adminProjects.GET("/types/stats", projectController.GetProjectTypeStats) // 获取项目分类统计

				// =============================================
				// 6. 新增的管理员项目统计路由（兼容前端调用）
				// =============================================
				adminProjects.GET("/quality-report", projectController.GetProjectStats)   // 获取项目质量报告（使用项目统计）
				adminProjects.GET("/type-stats", projectController.GetProjectTypeStats)   // 获取项目类型统计
				adminProjects.GET("/department-stats", projectController.GetProjectStats) // 获取院系项目统计（使用项目统计）
				adminProjects.GET("/time-trend", projectController.GetProjectStats)       // 获取时间趋势统计（使用项目统计）
			}

			// 学生专用路由
			students := auth.Group("/students")
			students.Use(middlewares.RoleMiddleware("student"))
			{
				students.POST("/bind-teacher", projectController.BindStudentToTeacher) // 学生绑定教师
			}

			// 学生获取教师列表路由（学生需要选择指导老师）
			studentTeachers := auth.Group("/student-teachers")
			//studentTeachers.Use(middlewares.RoleMiddleware("student"))
			{
				studentTeachers.GET("", projectController.GetTeacherList) // 学生获取教师列表
			}

			// 通用项目路由（所有认证用户）
			projects := auth.Group("/projects")
			{
				projects.GET("/my", projectController.GetMyProjects)                                   // 学生获取我的项目
				projects.GET("/status", projectController.GetProjectStats)                             // 获取项目统计信息
				projects.GET("/detail", projectController.GetProjectByID)                              // 查看项目详情
				projects.POST("", projectController.CreateProject)                                     // 学生创建项目
				projects.DELETE("/delete", projectController.DeleteProject)                            // 学生删除项目
				projects.POST("/submit", projectController.SubmitProject)                              // 学生提交项目审核项目
				projects.POST("/extensionapplication", projectController.CreateExtensionApplication)   // 学生申请项目延期
				projects.GET("/MyExtensionApplications", projectController.GetMyExtensionApplications) // 学生获取我的延期申请

				projects.GET("/files/getfiles", projectController.GetFiles)        // 获取项目成果文件列表
				projects.POST("/files/uploadfiles", projectController.UploadFiles) // 上传项目成果文件
				projects.DELETE("/files/delete", projectController.DeleteFile)     // 删除项目成果文件

				projects.POST("/progress/update", projectController.UpdateProjectProgress) // 学生更新项目进度
				projects.GET("/progress", projectController.GetProjectProgress)            // 学生和教师获取指定项目的进度

			}

			// 文件上传路由（所有认证用户）
			files := auth.Group("/files")
			{
				files.POST("/upload", fileController.UploadFile) // 上传文件
			}

			// 竞赛管理路由
			competitionController := controllers.NewCompetitionController(db)
			competitions := auth.Group("/competitions")
			{
				competitions.GET("", competitionController.GetCompetitionList) // 获取竞赛列表

				competitions.GET("/stats", competitionController.GetCompetitionStats) // 获取竞赛统计
			}

			// 学生竞赛路由
			studentCompetitions := auth.Group("/student-competitions")
			//studentCompetitions.Use(middlewares.RoleMiddleware("student"))
			{
				studentCompetitions.POST("/register", competitionController.RegisterCompetition)        // 报名竞赛
				studentCompetitions.GET("/my", competitionController.GetMyRegistrations)                // 查看自己的报名记录
				studentCompetitions.POST("/submissions", competitionController.SubmitCompetitionResult) // 上传参赛成果
			}

			// 教师竞赛路由
			teacherCompetitions := auth.Group("/teacher-competitions")
			teacherCompetitions.Use(middlewares.RoleMiddleware("teacher", "admin"))
			{
				teacherCompetitions.GET("/:id/submissions", competitionController.GetCompetitionSubmissions) // 查看竞赛提交作品
				teacherCompetitions.POST("/score", competitionController.ScoreCompetition)                   // 提交评语

				// 竞赛评审路由
				teacherCompetitions.POST("/submissions/:submissionId/scores", competitionController.SubmitScore)        // 提交评分
				teacherCompetitions.GET("/submissions/:submissionId/scores", competitionController.GetSubmissionScores) // 获取提交的评分列表
			}

			// 管理员竞赛管理路由
			adminCompetitions := auth.Group("/admin/competitions")
			//adminCompetitions.Use(middlewares.AdminOnly())
			{
				adminCompetitions.POST("", competitionController.CreateCompetition)                // 创建竞赛
				adminCompetitions.GET("/:id/detail", competitionController.GetCompetitionDetail)   // 查看竞赛详情
				adminCompetitions.POST("/:id/isopen", competitionController.ToggleCompetitionOpen) // 切换竞赛开放状态

				adminCompetitions.DELETE("/:id", competitionController.DeleteCompetition)                      // 删除竞赛
				adminCompetitions.GET("/:id/registrations", competitionController.GetCompetitionRegistrations) // 查看某竞赛所有报名
				adminCompetitions.POST("/:id/verify", competitionController.VerifyRegistration)                // 管理员审核学生报名

				adminCompetitions.POST("/:id/result", competitionController.SubmitResult)         // 登记成绩/获奖信息
				adminCompetitions.GET("/:id/export", competitionController.ExportCompetitionData) // 导出竞赛数据

				// 竞赛评审管理
				adminCompetitions.POST("/judges/distribute", competitionController.AssignJudge) // 分配评审教师

				adminCompetitions.GET("/:id/judges", competitionController.GetCompetitionJudges) // 获取评审教师列表

				adminCompetitions.GET("/:id/judging-progress", competitionController.GetJudgingProgress) // 获取评审进度
				adminCompetitions.POST("/:id/finalize", competitionController.FinalizeResults)           // 最终确认成绩
			}

			// 通知系统路由
			notificationService := services.NewNotificationService(db)
			notificationController := controllers.NewNotificationController(notificationService)
			notifications := auth.Group("/notifications")
			{
				notifications.GET("", notificationController.GetMyNotifications)                  // 获取我的通知列表
				notifications.GET("/unread-count", notificationController.GetUnreadCount)         // 获取未读通知数量
				notifications.PUT("/:id/read", notificationController.MarkNotificationAsRead)     // 标记通知为已读
				notifications.PUT("/read-all", notificationController.MarkAllNotificationsAsRead) // 标记所有通知为已读
				notifications.DELETE("/:id", notificationController.DeleteNotification)           // 删除通知
			}

			// 管理员通知管理路由
			adminNotifications := auth.Group("/admin/notifications")
			adminNotifications.Use(middlewares.AdminOnly())
			{
				adminNotifications.GET("/templates", notificationController.GetNotificationTemplates)       // 获取通知模板列表
				adminNotifications.PUT("/templates/:id", notificationController.UpdateNotificationTemplate) // 更新通知模板
				adminNotifications.POST("/send", notificationController.SendNotification)                   // 发送通知
			}

			// 管理员专用路由
			adminController := controllers.NewAdminController(userService)
			systemController := controllers.NewSystemController(db)

			admin := auth.Group("/admin")
			admin.Use(middlewares.AdminOnly())
			{
				// 仪表板
				admin.GET("/dashboard", adminController.GetDashboardStats)       // 获取仪表板数据
				admin.GET("/dashboard/stats", adminController.GetDashboardStats) // 获取仪表板统计数据 (兼容前端调用)
				admin.GET("/overview", adminController.GetUserOverview)          // 获取用户概览

				// 系统管理
				admin.GET("/logs", systemController.GetSystemLogs)                         // 获取系统日志
				admin.GET("/logs/summary", systemController.GetSystemLogsSummary)          // 获取系统日志统计
				admin.GET("/logs/health", systemController.GetSystemHealthLogs)            // 获取系统健康日志
				admin.GET("/logs/health/summary", systemController.GetSystemHealthSummary) // 获取系统健康统计
				admin.POST("/logs/health", systemController.RecordSystemHealth)            // 记录系统健康状态
				admin.GET("/settings", systemController.GetSystemSettings)                 // 获取系统设置
				admin.GET("/settings/:key", systemController.GetSystemSettingByKey)        // 根据键获取系统设置
				admin.PUT("/settings/:key", systemController.UpdateSystemSetting)          // 更新系统设置
				admin.GET("/config", systemController.GetSystemConfig)                     // 获取系统配置结构
				admin.GET("/health", systemController.GetSystemHealth)                     // 获取系统健康状态
				admin.GET("/system/health", systemController.GetSystemHealth)              // 获取系统健康状态 (兼容前端调用)
				admin.GET("/stats", systemController.GetSystemStats)                       // 获取系统统计
				admin.PUT("/maintenance", systemController.UpdateMaintenanceMode)          // 更新维护模式
				admin.POST("/logs/cleanup", systemController.CleanupOldLogs)               // 清理过期日志

				// 系统性能监控
				admin.GET("/performance", systemController.GetSystemPerformance)            // 获取系统性能数据
				admin.POST("/performance/record", systemController.RecordSystemPerformance) // 记录系统性能数据

				// 系统告警管理
				admin.GET("/alerts", systemController.GetSystemAlerts)                   // 获取系统告警列表
				admin.POST("/alerts/:id/acknowledge", systemController.AcknowledgeAlert) // 确认告警
				admin.POST("/alerts/:id/resolve", systemController.ResolveAlert)         // 解决告警

				// 系统诊断
				admin.GET("/diagnostics", systemController.GetSystemDiagnostics)      // 获取系统诊断记录
				admin.POST("/diagnostics/run", systemController.RunSystemDiagnostics) // 运行系统诊断

				// 备份管理
				admin.GET("/backups", systemController.GetBackupRecords)               // 获取备份记录
				admin.POST("/backups", systemController.CreateBackup)                  // 创建备份
				admin.GET("/backups/statistics", systemController.GetBackupStatistics) // 获取备份统计
			}
		}
	}
}
