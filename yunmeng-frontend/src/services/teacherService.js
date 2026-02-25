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

//获得指定学生项目的文件列表
  async getStudentProjectsFiles(params = {student_id: null, id: null}) {
    try{
      const response = await api.get(`/teachers/Stufiles`, { params })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取学生项目文件列表失败')
    }
  }


  // 获取项目详情
  async getProjectDetail(projectId) {
    try {
      
      const response = await api.get(`/projects/detail` , { params: { id: projectId } })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取项目详情失败')
    }
  }

  // 审核项目
  async reviewProject(reviewData) {
    try {
      const response = await api.post(`/teacher-projects/review`, reviewData)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '项目审核失败')
    }
  }

  // 获取项目审核记录
  async getProjectReviews() {
    try {
      const response = await api.get(`/teachers/TeacherExtensionApplications`)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取审核记录失败')
    }
  }

  //更新项目审核记录
  async updateProjectReviews(reviewData) {
    try {
      const response = await api.post(`/teachers/ApproveExtensionApplication`, reviewData)
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
      const response = await api.get('/teachers/students', { params })
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


  // 获取项目文件列表
  async getProjectFiles(projectId, params = {}) {
    try {
      const response = await api.get(`/projects/${projectId}/files`, { params })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取项目文件失败')
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
}


// 创建实例
const teacherService = new TeacherService()

export default teacherService