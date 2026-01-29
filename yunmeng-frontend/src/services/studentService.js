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


/**
 * 学生服务
 * 提供学生特有的项目管理功能
 */
export const studentService = {
  /**
   * 获取我的项目列表
   */
  async getMyProjects(params = {}) {
    try {

      // 使用 api 而不是 request
      const response = await api.get('/projects/my', { params })
      
      console.log('✅ 请求成功，响应:', response)
      return response
    } catch (error)  {
      console.error('获取我的项目列表失败:', error)
      throw error
    }
  },

  /**
   * 获取项目统计数据
   */
  async getProjectStats() {
    try {
      const response = await api.get('/projects/status')
      return response
    } catch (error) {
      console.error('获取项目统计数据失败:', error)
      // 返回模拟数据作为后备
      return {
        code: 200,
        data: {
          totalProjects: 3,
          ongoingProjects: 1,
          completedProjects: 1,
          pendingProjects: 1,
          totalCompetitions: 2,
          activeCompetitions: 1
        }
      }
    }
  },

  /**
   * 创建新项目
   */
  async createProject(projectData) {
    try {
      const response = await api.post('/student/projects', projectData)
      return response
    } catch (error) {
      console.error('创建项目失败:', error)
      throw error
    }
  },

  /**
   * 更新项目
   */
  async updateProject(projectId, updateData) {
    try {
      const response = await api.put(`/student/projects/${projectId}`, updateData)
      return response
    } catch (error) {
      console.error('更新项目失败:', error)
      throw error
    }
  },

  /**
   * 删除项目
   */
  async deleteProject(projectId) {
    try {
      const response = await api.delete(`/student/projects/${projectId}`)
      return response
    } catch (error) {
      console.error('删除项目失败:', error)
      throw error
    }
  },

    /**
     * 获取我的项目延期申请
     */
    async getMyExtensionApplications() {
    try {
      const response = await api.get('/projects/MyExtensionApplications')
      return response
    } catch (error) {
      console.error('获取延期申请失败:', error)
      throw error
    }
  },

  

    /**
   * 创建项目延期申请
   */
    async createExtensionApplication(extensionData) {
    try {
      const response = await api.post('/projects/extensionapplication', extensionData)
      return response
    } catch (error) {
      console.error('创建延期申请失败:', error)
      throw error
    }
  },




  /**
   * 获取项目详情
   */
  async getProjectDetail(projectId) {
    try {
      const response = await api.get(`/student/projects/${projectId}`)
      return response
    } catch (error) {
      console.error('获取项目详情失败:', error)
      throw error
    }
  },

  /**
   * 提交项目审核
   */
  async submitProject(projectId) {
    try {
      const response = await api.post(`/student/projects/${projectId}/submit`)
      return response
    } catch (error) {
      console.error('提交项目失败:', error)
      throw error
    }
  },

  /**
   * 更新项目进度
   */
  async updateProjectProgress(projectId, progressData) {
    try {
      const response = await api.put(`/student/projects/${projectId}/progress`, progressData)
      return response
    } catch (error) {
      console.error('更新项目进度失败:', error)
      throw error
    }
  },


  /**
   * 获取项目附件列表
   */
  async getProjectFiles(projectId) {
    try {
      const response = await api.get(`/student/projects/${projectId}/files`)
      return response
    } catch (error) {
      console.error('获取项目附件失败:', error)
      throw error
    }
  },

  /**
   * 上传项目附件
   */
  async uploadProjectFile(projectId, fileData) {
    try {
      const response = await api.post(`/student/projects/${projectId}/files`, fileData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      })
      return response
    } catch (error) {
      console.error('上传附件失败:', error)
      throw error
    }
  },

  /**
   * 删除项目附件
   */
  async deleteProjectFile(projectId, fileId) {
    try {
      const response = await api.delete(`/student/projects/${projectId}/files/${fileId}`)
      return response
    } catch (error) {
      console.error('删除附件失败:', error)
      throw error
    }
  },


  /**
   * 导出我的项目
   */
  async exportMyProjects(filters = {}) {
    try {
      const response = await api.get('/student/projects/export', {
        params: filters,
        responseType: 'blob'
      })
      return response
    } catch (error) {
      console.error('导出项目失败:', error)
      throw error
    }
  },

  /**
   * 获取项目统计报告
   */
  async getProjectReport(params = {}) {
    try {
      const response = await api.get('/student/projects/report', { params })
      return response
    } catch (error) {
      console.error('获取项目报告失败:', error)
      throw error
    }
  },


}

export default studentService 