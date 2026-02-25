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
    console.log('API响应原始数据:', response.data) // 调试日志
    return response.data
  },
  async error => {
    console.error('API请求错误:', error)
    
    // 处理不同的错误状态码
    if (error.response?.status === 401) {
      // 尝试刷新Token
      try {
        const tokenManager = await import('./tokenManager.js')
        const newToken = await tokenManager.default.refreshToken()
        console.log('Token自动刷新成功')
        
        // 重新发送原始请求
        const originalRequest = error.config
        originalRequest.headers['Authorization'] = `Bearer ${newToken}`
        return api(originalRequest)
      } catch (refreshError) {
        console.error('Token刷新失败:', refreshError)
        // 清除过期的认证信息
        localStorage.removeItem('token')
        localStorage.removeItem('userInfo')
        localStorage.removeItem('userRole')
        console.log('登录已过期，请重新登录')
      }
    } else if (error.response?.status === 403) {
      console.error('权限不足，需要管理员权限')
      // 可以在这里添加全局的错误提示
    }
    
    return Promise.reject(error)
  }
)

// 竞赛相关API服务
export const competitionService = {
  // ==================== 公共API ====================
  
  // 获取公开竞赛列表
  async getPublicCompetitions(params = {}) {
    try {
      const response = await api.get('/competitions/public', { params })
      return response
    } catch (error) {
      console.error('获取公开竞赛列表失败:', error)
      throw error
    }
  },


  // 获取竞赛统计信息
  async getCompetitionStats() {
    try {
      const response = await api.get('/competitions/stats')
      return response
    } catch (error) {
      console.error('获取竞赛统计信息失败:', error)
      throw error
    }
  },

  // ==================== 管理员API ====================
  
  // 获取所有竞赛列表（管理员）
  async getCompetitions(params = {}) {
    try {
      const response = await api.get('/competitions', { params })
      return response
    } catch (error) {
      console.error('获取竞赛列表失败:', error)
      throw error
    }
  },



 //管理员获取评审教师列表
async getJudgesByCompetition(competitionId) {
  try{
    const response = await api.get(`/admin/competitions/${competitionId}/judges`)
    return response
  }catch(error){
    console.error('获取评审教师列表失败:', error)
  }
},


  // 创建竞赛（管理员）
  async createCompetition(competitionData) {
    try {
      const response = await api.post('/admin/competitions', competitionData)
      return response
    } catch (error) {
      console.error('创建竞赛失败:', error)
      throw error
    }
  },

// 获取竞赛详情
async getCompetitionDetail(id) {
  try {
    const response = await api.get(`/admin/competitions/${id}/detail`)
    return response
  } catch (error) {
    console.error('获取竞赛详情失败:', error)
    throw error
  }
},


//分配新的评审教师
async distributeJudge(distributeData) {
  try {
    const response = await api.post('/admin/competitions/judges/distribute', distributeData)
    return response
  } catch (error) {
    console.error('分配评审教师失败:', error)
    throw error
  }
},


// 切换开放状态
async toggleCompetitionOpen(id) {
  try {
    const response = await api.post(`/admin/competitions/${id}/isopen`)
    return response
  } catch (error) {
    console.error('切换报名状态失败:', error)
    throw error
  }
},

  // 更新竞赛（管理员）
  async updateCompetition(id, competitionData) {
    try {
      const response = await api.put(`/admin/competitions/${id}`, competitionData)
      return response
    } catch (error) {
      console.error('更新竞赛失败:', error)
      throw error
    }
  },

  // 删除竞赛（管理员）
  async deleteCompetition(id) {
    try {
      const response = await api.delete(`/admin/competitions/${id}`)
      return response
    } catch (error) {
      console.error('删除竞赛失败:', error)
      throw error
    }
  },

  // 获取竞赛报名记录（管理员）
  async getCompetitionRegistrations(competitionId, params = {}) {
    try {
      const response = await api.get(`/admin/competitions/${competitionId}/registrations`, { params })
      return response
    } catch (error) {
      console.error('获取竞赛报名记录失败:', error)
      throw error
    }
  },

  // 获取竞赛提交作品（管理员）
  async getCompetitionSubmissions(competitionId, params = {}) {
    try {
      const response = await api.get(`/competitions/${competitionId}/submissions`, { params })
      return response
    } catch (error) {
      console.error('获取竞赛提交作品失败:', error)
      throw error
    }
  },

  // 登记获奖结果（管理员）
  async submitResult(competitionId, resultData) {
    try {
      const response = await api.post(`/admin/competitions/${competitionId}/result`, resultData)
      return response
    } catch (error) {
      console.error('登记获奖结果失败:', error)
      throw error
    }
  },

  // 导出竞赛数据（管理员）
  async exportCompetitionData(competitionId, type = 'registrations') {
    try {
      const response = await api.get(`/admin/competitions/${competitionId}/export`, {
        params: { type },
        responseType: 'blob'
      })
      return response
    } catch (error) {
      console.error('导出竞赛数据失败:', error)
      throw error
    }
  },

  // 分配评审教师（管理员）
  async assignJudge(competitionId, judgeData) {
    try {
      const response = await api.post(`/admin/competitions/${competitionId}/judges`, judgeData)
      return response
    } catch (error) {
      console.error('分配评审教师失败:', error)
      throw error
    }
  },

  // 获取评审教师列表（管理员）
  async getCompetitionJudges(competitionId) {
    try {
      const response = await api.get(`/admin/competitions/${competitionId}/judges`)
      return response
    } catch (error) {
      console.error('获取评审教师列表失败:', error)
      throw error
    }
  },

  // 获取评审进度（管理员）
  async getJudgingProgress(competitionId) {
    try {
      const response = await api.get(`/admin/competitions/${competitionId}/judging-progress`)
      return response
    } catch (error) {
      console.error('获取评审进度失败:', error)
      throw error
    }
  },

  // 最终确认成绩（管理员）
  async finalizeResults(competitionId) {
    try {
      const response = await api.post(`/admin/competitions/${competitionId}/finalize`)
      return response
    } catch (error) {
      console.error('最终确认成绩失败:', error)
      throw error
    }
  },

  // ==================== 学生API ====================
  
  // 获取可报名竞赛（学生）
  async getAvailableCompetitions(params = {}) {
    try {
      const response = await api.get('/competitions/available', { params })
      return response
    } catch (error) {
      console.error('获取可报名竞赛失败:', error)
      throw error
    }
  },

  // 报名竞赛（学生）
  async registerCompetition( registrationData) {
    try {
      const response = await api.post(`/student-competitions/register`, registrationData)
      return response
    } catch (error) {
      console.error('报名竞赛失败:', error)
      throw error
    }
  },

  // 获取我的报名记录（学生）
  async getMyRegistrations() {
    try {
      const response = await api.get('/student-competitions/my')
      return response
    } catch (error) {
      console.error('获取我的报名记录失败:', error)
      throw error
    }
  },

  // 更新报名信息（学生）
  async updateRegistration(registrationId, registrationData) {
    try {
      const response = await api.put(`/competition-registrations/${registrationId}`, registrationData)
      return response
    } catch (error) {
      console.error('更新报名信息失败:', error)
      throw error
    }
  },

  // 提交参赛作品（学生）
  async submitWork(competitionId, submissionData) {
    try {
      const formData = new FormData()
      formData.append('file_url', submissionData.file_url)
      formData.append('file_name', submissionData.file_name)
      formData.append('file_size', submissionData.file_size)
      formData.append('description', submissionData.description)
      
      // 添加文件
      if (submissionData.file) {
        formData.append('file', submissionData.file)
      }

      const response = await api.post(`/competitions/${competitionId}/upload`, formData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      })
      return response
    } catch (error) {
      console.error('提交参赛作品失败:', error)
      throw error
    }
  },

  // 获取竞赛结果（学生）
  async getCompetitionResults(params = {}) {
    try {
      const response = await api.get('/competitions/results', { params })
      return response
    } catch (error) {
      console.error('获取竞赛结果失败:', error)
      throw error
    }
  },

  // ==================== 教师API ====================
  
  // 获取指导学生报名记录（教师）
  async getStudentRegistrations(params = {}) {
    try {
      const response = await api.get('/teacher/competition-registrations', { params })
      return response
    } catch (error) {
      console.error('获取指导学生报名记录失败:', error)
      throw error
    }
  },

  // 审核报名（教师）
  async reviewRegistration(registrationId, reviewData) {
    try {
      const response = await api.put(`/teacher/competition-registrations/${registrationId}/review`, reviewData)
      return response
    } catch (error) {
      console.error('审核报名失败:', error)
      throw error
    }
  },

  // 查看学生作品（教师）
  async viewStudentSubmission(competitionId, params = {}) {
    try {
      const response = await api.get(`/teacher/competitions/${competitionId}/submissions`, { params })
      return response
    } catch (error) {
      console.error('查看学生作品失败:', error)
      throw error
    }
  },

  // 提交评语（教师）
  async submitFeedback(competitionId, feedbackData) {
    try {
      const response = await api.post(`/competitions/${competitionId}/feedback`, feedbackData)
      return response
    } catch (error) {
      console.error('提交评语失败:', error)
      throw error
    }
  },

  // 获取评语历史（教师）
  async getFeedbackHistory(competitionId, params = {}) {
    try {
      const response = await api.get(`/teacher/competitions/${competitionId}/feedback-history`, { params })
      return response
    } catch (error) {
      console.error('获取评语历史失败:', error)
      throw error
    }
  },

  // 提交评分（教师）
  async submitScore(submissionId, scoreData) {
    try {
      const response = await api.post(`/competitions/submissions/${submissionId}/scores`, scoreData)
      return response
    } catch (error) {
      console.error('提交评分失败:', error)
      throw error
    }
  },

  // 获取提交评分（教师）
  async getSubmissionScores(submissionId) {
    try {
      const response = await api.get(`/competitions/submissions/${submissionId}/scores`)
      return response
    } catch (error) {
      console.error('获取提交评分失败:', error)
      throw error
    }
  },

  // 获取学生结果（教师）
  async getStudentResults(competitionId, params = {}) {
    try {
      const response = await api.get(`/teacher/competitions/${competitionId}/results`, { params })
      return response
    } catch (error) {
      console.error('获取学生结果失败:', error)
      throw error
    }
  }
}

export default competitionService 