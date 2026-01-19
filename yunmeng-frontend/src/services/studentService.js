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
      const response = await api.get('/student/projects', { params })
      return response
    } catch (error) {
      console.error('获取我的项目列表失败:', error)
      // 返回模拟数据作为后备
      return {
        code: 200,
        data: [
          {
            id: 1,
            name: '智能校园系统',
            type: '软件开发',
            status: '进行中',
            progress: 75,
            createTime: '2024-01-15',
            deadline: '2024-07-15',
            description: '基于物联网技术的智能校园管理系统',
            plan: '预计6个月完成，分为需求分析、设计、开发、测试四个阶段'
          },
          {
            id: 2,
            name: '数据分析平台',
            type: '科研项目',
            status: '待审核',
            progress: 90,
            createTime: '2024-01-14',
            deadline: '2024-06-14',
            description: '大数据分析平台，支持多种数据源和算法',
            plan: '预计8个月完成，包括数据采集、预处理、分析、可视化等模块'
          },
          {
            id: 3,
            name: '在线教育平台',
            type: '创新项目',
            status: '已完成',
            progress: 100,
            createTime: '2024-01-10',
            deadline: '2024-05-10',
            description: '基于Web的在线教育学习平台',
            plan: '预计4个月完成，包括用户管理、课程管理、学习跟踪等模块'
          }
        ]
      }
    }
  },

  /**
   * 获取项目统计数据
   */
  async getProjectStats() {
    try {
      const response = await api.get('/student/projects/stats')
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
   * 申请项目延期
   */
  async requestExtension(projectId, extensionData) {
    try {
      const response = await api.post(`/student/projects/${projectId}/extension`, extensionData)
      return response
    } catch (error) {
      console.error('申请延期失败:', error)
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
   * 获取项目成员列表
   */
  async getProjectMembers(projectId) {
    try {
      const response = await api.get(`/student/projects/${projectId}/members`)
      return response
    } catch (error) {
      console.error('获取项目成员失败:', error)
      throw error
    }
  },

  /**
   * 邀请项目成员
   */
  async inviteProjectMember(projectId, memberData) {
    try {
      const response = await api.post(`/student/projects/${projectId}/members`, memberData)
      return response
    } catch (error) {
      console.error('邀请成员失败:', error)
      throw error
    }
  },

  /**
   * 移除项目成员
   */
  async removeProjectMember(projectId, memberId) {
    try {
      const response = await api.delete(`/student/projects/${projectId}/members/${memberId}`)
      return response
    } catch (error) {
      console.error('移除成员失败:', error)
      throw error
    }
  },

  /**
   * 获取项目审核记录
   */
  async getProjectReviews(projectId) {
    try {
      const response = await api.get(`/student/projects/${projectId}/reviews`)
      return response
    } catch (error) {
      console.error('获取审核记录失败:', error)
      throw error
    }
  },

  /**
   * 获取项目时间线
   */
  async getProjectTimeline(projectId) {
    try {
      const response = await api.get(`/student/projects/${projectId}/timeline`)
      return response
    } catch (error) {
      console.error('获取项目时间线失败:', error)
      throw error
    }
  },

  /**
   * 获取项目模板列表
   */
  async getProjectTemplates() {
    try {
      const response = await api.get('/student/project-templates')
      return response
    } catch (error) {
      console.error('获取项目模板失败:', error)
      throw error
    }
  },

  /**
   * 获取项目模板详情
   */
  async getProjectTemplate(templateId) {
    try {
      const response = await api.get(`/student/project-templates/${templateId}`)
      return response
    } catch (error) {
      console.error('获取项目模板详情失败:', error)
      throw error
    }
  },

  /**
   * 从模板创建项目
   */
  async createProjectFromTemplate(templateId, projectData) {
    try {
      const response = await api.post(`/student/project-templates/${templateId}/create`, projectData)
      return response
    } catch (error) {
      console.error('从模板创建项目失败:', error)
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

  /**
   * 获取项目建议
   */
  async getProjectSuggestions() {
    try {
      const response = await api.get('/student/projects/suggestions')
      return response
    } catch (error) {
      console.error('获取项目建议失败:', error)
      throw error
    }
  },

  /**
   * 获取项目通知
   */
  async getProjectNotifications(params = {}) {
    try {
      const response = await api.get('/student/projects/notifications', { params })
      return response
    } catch (error) {
      console.error('获取项目通知失败:', error)
      throw error
    }
  },

  /**
   * 标记通知为已读
   */
  async markNotificationRead(notificationId) {
    try {
      const response = await api.put(`/student/notifications/${notificationId}/read`)
      return response
    } catch (error) {
      console.error('标记通知已读失败:', error)
      throw error
    }
  },

  /**
   * 获取项目协作记录
   */
  async getProjectCollaboration(projectId) {
    try {
      const response = await api.get(`/student/projects/${projectId}/collaboration`)
      return response
    } catch (error) {
      console.error('获取协作记录失败:', error)
      throw error
    }
  },

  /**
   * 添加协作记录
   */
  async addCollaborationRecord(projectId, recordData) {
    try {
      const response = await api.post(`/student/projects/${projectId}/collaboration`, recordData)
      return response
    } catch (error) {
      console.error('添加协作记录失败:', error)
      throw error
    }
  }
}

export default studentService 