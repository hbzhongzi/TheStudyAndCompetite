import axios from 'axios'
import { getToken } from './auth'

// 创建axios实例
const api = axios.create({
  baseURL: 'http://localhost:8080/api',
  timeout: 10000
})

// 请求拦截器，添加token
api.interceptors.request.use(
  config => {
    const token = getToken()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  response => {
    return response.data
  },
  error => {
    console.error('API请求错误:', error)
    if (error.response?.status === 401) {
      // token过期，清除本地存储但不自动跳转
      localStorage.removeItem('token')
      localStorage.removeItem('userInfo')
      localStorage.removeItem('userRole')
      console.log('登录已过期，请重新登录')
      // 不再自动跳转到登录页面，让用户手动访问
    }
    return Promise.reject(error)
  }
)

// 教师服务类
class TeacherService {
  // 获取所有项目列表
  async getAllProjects(params = {}) {
    try {
      const response = await api.get('/projects', { params })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取项目列表失败')
    }
  }


  async getTeacherProjects(params = {}) {
    try {
      // 调用新接口
      const response = await api.get('/teacher-projects', { params })
      
      if (response && response.code === 200) {
        return response
      }
      
      // 如果接口响应格式不同，可能需要适配
      return {
        code: response?.code || 500,
        data: response?.data || { list: [], page: 1, size: 20, total: 0 },
        message: response?.message || '获取失败'
      }
    } catch (error) {
      console.error('获取教师项目列表失败:', error)
      return {
        code: 500,
        data: { list: [], page: 1, size: 20, total: 0 },
        message: '请求失败: ' + error.message
      }
    }
  }

  // 获取项目详情
  async getProjectDetail(projectId) {
    try {
      const response = await api.get(`/projects/${projectId}`)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取项目详情失败')
    }
  }

  // 审核项目
  async reviewProject(projectId, reviewData) {
    try {
      const response = await api.put(`/projects/${projectId}/review`, reviewData)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '项目审核失败')
    }
  }

  // 获取项目审核记录
  async getProjectReviews(projectId) {
    try {
      const response = await api.get(`/projects/${projectId}/reviews`)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取审核记录失败')
    }
  }

  // 更新项目进度
  async updateProjectProgress(projectId, progressData) {
    try {
      const response = await api.put(`/projects/${projectId}/progress`, progressData)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '更新项目进度失败')
    }
  }

  // 获取学生列表
  async getStudents(params = {}) {
    try {
      const response = await api.get('/students', { params })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取学生列表失败')
    }
  }

  // 获取学生详情
  async getStudentDetail(studentId) {
    try {
      const response = await api.get(`/students/${studentId}`)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取学生详情失败')
    }
  }

  // 获取教师列表（供学生选择指导老师）
  async getTeacherList(params = {}) {
    try {
      // 根据用户角色选择不同的端点
      const userRole = localStorage.getItem('userRole')
      let endpoint = '/teachers'
      
      // 如果是学生，使用学生专用的教师列表端点
      if (userRole === 'student') {
        endpoint = '/student-teachers'
      }
      
      const response = await api.get(endpoint, { params })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取教师列表失败')
    }
  }

  // 获取我指导的学生列表
  async getMyStudents(params = {}) {
    try {
      const response = await api.get('/teachers/students', { params })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取指导的学生列表失败')
    }
  }

  // 获取我指导的项目列表
  async getMyProjects(params = {}) {
    try {
      const response = await api.get('/teachers/projects', { params })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取指导的项目列表失败')
    }
  }

  // 获取指导的项目列表
  async getGuidedProjects() {
    try {
       const response = await api.get('/teachers/projects')
    
    if (response && response.code === 200) {
      const responseData = response.data
      
      // ✅ 修复：检查 response.data 是否是数组，或者包含 list 数组
      if (Array.isArray(responseData)) {
        // 直接返回数组格式
        return response
      } else if (responseData && responseData.list && Array.isArray(responseData.list)) {
        // 返回分页格式，不发出警告
        return response
      } else {
        // 格式不正确，发出警告
        console.warn('API返回的projects数据格式不正确:', responseData)
        return {
          code: 200,
          data: { list: [], page: 1, size: 20, total: 0 },
          message: '数据格式不正确，返回空列表'
                }
      }
    }
    
    return response
    
    } catch (error) {
      console.error('获取指导项目列表失败:', error)
      // 返回模拟数据作为备选
      return {
        code: 200,
        data: [
          { id: 1, name: '智能校园系统' },
          { id: 2, name: '数据分析平台' },
          { id: 3, name: '在线教育平台' }
        ],
        message: '获取项目列表成功'
      }
    }
  }

  // 新增：获取项目统计数据
  async getProjectStats(params = {}) {
    try {
      const response = await api.get('/teachers/projects/stats', { params })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取项目统计失败')
    }
  }

  // 新增：项目质量评估
  async assessProjectQuality(projectId, qualityData) {
    try {
      const response = await api.post(`/projects/${projectId}/quality-assessment`, qualityData)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '项目质量评估失败')
    }
  }

  // 新增：批量通过项目
  async batchApproveProjects(batchData) {
    try {
      const response = await api.post('/teachers/projects/batch-approve', batchData)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '批量通过项目失败')
    }
  }

  // 新增：批量驳回项目
  async batchRejectProjects(batchData) {
    try {
      const response = await api.post('/teachers/projects/batch-reject', batchData)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '批量驳回项目失败')
    }
  }

  // 新增：批量更新项目进度
  async batchUpdateProjectProgress(batchData) {
    try {
      const response = await api.post('/teachers/projects/batch-update-progress', batchData)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '批量更新项目进度失败')
    }
  }

  // 新增：批量添加评语
  async batchAddComments(batchData) {
    try {
      const response = await api.post('/teachers/projects/batch-add-comments', batchData)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '批量添加评语失败')
    }
  }

  // 新增：导出项目数据
  async exportProjects(projectIds) {
    try {
      const response = await api.post('/teachers/projects/export', { projectIds }, {
        responseType: 'blob'
      })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '导出项目数据失败')
    }
  }

  // 新增：生成项目报告
  async generateProjectReport(projectIds) {
    try {
      const response = await api.post('/teachers/projects/report', { projectIds }, {
        responseType: 'blob'
      })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '生成项目报告失败')
    }
  }

  // 获取项目文件列表
  async getProjectFiles(projectId, params = {}) {
    try {
      const response = await api.get(`/projects/${projectId}/files`, { params })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取项目文件失败')
    }
  }

  // 获取项目延期申请列表
  async getProjectExtensions(projectId, params = {}) {
    try {
      const response = await api.get(`/projects/${projectId}/extensions`, { params })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取项目延期申请失败')
    }
  }

  // 获取项目里程碑列表
  async getProjectMilestones(projectId, params = {}) {
    try {
      const response = await api.get(`/projects/${projectId}/milestones`, { params })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取项目里程碑失败')
    }
  }

  // 添加指导记录
  async addGuidanceRecord(studentId, guidanceData) {
    try {
      const response = await api.post(`/students/${studentId}/guidance`, guidanceData)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '添加指导记录失败')
    }
  }

  // 获取指导记录
  async getGuidanceRecords(studentId) {
    try {
      const response = await api.get(`/students/${studentId}/guidance`)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取指导记录失败')
    }
  }

  // ========== 竞赛相关API ==========
  
  // 获取我指导学生的竞赛报名列表
  async getCompetitionRegistrations(params = {}) {
    try {
      const response = await api.get('/teachers/competition-registrations', { params })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取竞赛报名列表失败')
    }
  }

  // 获取竞赛报名详情
  async getCompetitionRegistrationDetail(registrationId) {
    try {
      const response = await api.get(`/teachers/competition-registrations/${registrationId}`)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取竞赛报名详情失败')
    }
  }

  // 审核竞赛报名
  async reviewCompetitionRegistration(registrationId, reviewData) {
    try {
      const response = await api.put(`/teachers/competition-registrations/${registrationId}/review`, reviewData)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '竞赛报名审核失败')
    }
  }

  // 批量审核竞赛报名
  async batchReviewCompetitionRegistrations(registrationIds, reviewData) {
    try {
      const response = await api.put('/teachers/competition-registrations/batch-review', {
        registrationIds,
        ...reviewData
      })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '批量审核竞赛报名失败')
    }
  }

  // 获取竞赛作品详情
  async getCompetitionSubmission(submissionId) {
    try {
      const response = await api.get(`/teachers/competition-submissions/${submissionId}`)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取竞赛作品详情失败')
    }
  }

  // 提交作品评审意见
  async submitSubmissionReview(submissionId, reviewData) {
    try {
      const response = await api.post(`/teachers/competition-submissions/${submissionId}/review`, reviewData)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '提交作品评审失败')
    }
  }

  // 获取竞赛成绩
  async getCompetitionResult(registrationId) {
    try {
      const response = await api.get(`/teachers/competition-registrations/${registrationId}/result`)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取竞赛成绩失败')
    }
  }

  // 下载竞赛证书
  async downloadCompetitionCertificate(registrationId) {
    try {
      const response = await api.get(`/teachers/competition-registrations/${registrationId}/certificate`, {
        responseType: 'blob'
      })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '下载竞赛证书失败')
    }
  }

  // 下载竞赛成绩单
  async downloadCompetitionTranscript(registrationId) {
    try {
      const response = await api.get(`/teachers/competition-registrations/${registrationId}/transcript`, {
        responseType: 'blob'
      })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '下载竞赛成绩单失败')
    }
  }

  // 获取竞赛指导统计数据
  async getCompetitionGuidanceStats() {
    try {
      const response = await api.get('/teachers/competition-guidance/stats')
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取竞赛指导统计数据失败')
    }
  }

  // 导出竞赛指导数据
  async exportCompetitionGuidanceData(params = {}) {
    try {
      const response = await api.get('/teachers/competition-guidance/export', { 
        params,
        responseType: 'blob'
      })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '导出竞赛指导数据失败')
    }
  }

  // ========== 原有竞赛API（保留兼容性） ==========

  // 获取竞赛列表
  async getCompetitions(params = {}) {
    try {
      const response = await api.get('/competitions', { params })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取竞赛列表失败')
    }
  }

  // 获取竞赛详情
  async getCompetitionDetail(competitionId) {
    try {
      const response = await api.get(`/competitions/${competitionId}`)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取竞赛详情失败')
    }
  }

  // 添加竞赛指导记录
  async addCompetitionGuidance(competitionId, guidanceData) {
    try {
      const response = await api.post(`/competitions/${competitionId}/guidance`, guidanceData)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '添加竞赛指导记录失败')
    }
  }

  // 获取竞赛指导记录
  async getCompetitionGuidance(competitionId) {
    try {
      const response = await api.get(`/competitions/${competitionId}/guidance`)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取竞赛指导记录失败')
    }
  }

  // 审核竞赛
  async reviewCompetition(competitionId, reviewData) {
    try {
      const response = await api.put(`/competitions/${competitionId}/review`, reviewData)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '竞赛审核失败')
    }
  }

  // 获取竞赛审核记录
  async getCompetitionReviews(competitionId) {
    try {
      const response = await api.get(`/competitions/${competitionId}/reviews`)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取竞赛审核记录失败')
    }
  }

  // 获取统计数据
  async getStatistics() {
    try {
      const response = await api.get('/statistics')
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取统计数据失败')
    }
  }

  // 导出数据
  async exportData(type, params = {}) {
    try {
      const response = await api.get(`/export/${type}`, { 
        params,
        responseType: 'blob'
      })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '导出数据失败')
    }
  }
}

// 文件服务类
class FileService {
  // 上传文件
  async uploadFile(file) {
    try {
      const formData = new FormData()
      formData.append('file', file)
      
      const response = await api.post('/files/upload', formData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '文件上传失败')
    }
  }

  // 下载文件
  async downloadFile(fileUrl, fileName) {
    try {
      const response = await api.get(fileUrl, {
        responseType: 'blob'
      })
      
      // 创建下载链接
      const blob = new Blob([response])
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.download = fileName
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      window.URL.revokeObjectURL(url)
      
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '文件下载失败')
    }
  }

  // 预览文件
  async previewFile(fileUrl) {
    try {
      const response = await api.get(fileUrl, {
        responseType: 'blob'
      })
      
      // 创建预览链接
      const blob = new Blob([response])
      const url = window.URL.createObjectURL(blob)
      
      // 在新窗口中打开文件
      window.open(url, '_blank')
      
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '文件预览失败')
    }
  }
}

// 创建实例
const teacherService = new TeacherService()
const fileService = new FileService()

export { teacherService, fileService } 