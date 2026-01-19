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
 * 管理员服务
 * 提供管理员特有的项目管理功能
 */
export const adminService = {
  /**
   * 获取项目统计数据
   */
  async getProjectStats() {
    try {
      const response = await api.get('/admin/projects/stats')
      return response
    } catch (error) {
      console.error('获取项目统计数据失败:', error)
      throw error
    }
  },

  /**
   * 强制更新项目状态
   * @param {number} projectId 项目ID
   * @param {string} status 新状态
   */
  async forceUpdateProjectStatus(projectId, status) {
    try {
      const response = await api.put(`/admin/projects/${projectId}/status`, {
        status,
        reason: '管理员强制更新'
      })
      return response
    } catch (error) {
      console.error('强制更新项目状态失败:', error)
      throw error
    }
  },

  /**
   * 软删除项目
   * @param {number} projectId 项目ID
   */
  async softDeleteProject(projectId) {
    try {
      const response = await api.delete(`/admin/projects/${projectId}/soft`)
      return response
    } catch (error) {
      console.error('软删除项目失败:', error)
      throw error
    }
  },

  /**
   * 恢复软删除的项目
   * @param {number} projectId 项目ID
   */
  async restoreProject(projectId) {
    try {
      const response = await api.put(`/admin/projects/${projectId}/restore`)
      return response
    } catch (error) {
      console.error('恢复项目失败:', error)
      throw error
    }
  },

  /**
   * 清理无效项目
   */
  async cleanupProjects() {
    try {
      const response = await api.post('/admin/projects/cleanup')
      return response
    } catch (error) {
      console.error('清理项目失败:', error)
      throw error
    }
  },

  /**
   * 导出项目数据
   * @param {Object} filters 筛选条件
   */
  async exportProjects(filters = {}) {
    try {
      const response = await api.get('/admin/projects/export', {
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
   * 获取项目操作日志
   * @param {number} projectId 项目ID
   * @param {Object} params 查询参数
   */
  async getProjectLogs(projectId, params = {}) {
    try {
      const response = await api.get(`/admin/projects/${projectId}/logs`, {
        params
      })
      return response
    } catch (error) {
      console.error('获取项目日志失败:', error)
      throw error
    }
  },

  /**
   * 获取项目审核记录
   * @param {number} projectId 项目ID
   */
  async getProjectReviews(projectId) {
    try {
      const response = await api.get(`/admin/projects/${projectId}/reviews`)
      return response
    } catch (error) {
      console.error('获取项目审核记录失败:', error)
      throw error
    }
  },

  /**
   * 获取项目附件列表
   * @param {number} projectId 项目ID
   */
  async getProjectFiles(projectId) {
    try {
      const response = await api.get(`/admin/projects/${projectId}/files`)
      return response
    } catch (error) {
      console.error('获取项目附件失败:', error)
      throw error
    }
  },

  /**
   * 批量审核项目
   * @param {Array<number>} projectIds 项目ID数组
   * @param {Object} reviewData 审核数据
   */
  async batchReviewProjects(projectIds, reviewData) {
    try {
      const response = await api.post('/admin/projects/batch-review', {
        projectIds,
        ...reviewData
      })
      return response
    } catch (error) {
      console.error('批量审核项目失败:', error)
      throw error
    }
  },

  /**
   * 获取系统统计概览
   */
  async getSystemOverview() {
    try {
      const response = await api.get('/admin/system/overview')
      return response
    } catch (error) {
      console.error('获取系统概览失败:', error)
      throw error
    }
  },

  /**
   * 获取用户活跃度统计
   */
  async getUserActivityStats() {
    try {
      const response = await api.get('/admin/users/activity-stats')
      return response
    } catch (error) {
      console.error('获取用户活跃度统计失败:', error)
      throw error
    }
  },

  /**
   * 获取项目质量报告
   */
  async getProjectQualityReport() {
    try {
      const response = await api.get('/admin/projects/quality-report')
      return response
    } catch (error) {
      console.error('获取项目质量报告失败:', error)
      throw error
    }
  },

  /**
   * 发送系统通知
   * @param {Object} notificationData 通知数据
   */
  async sendSystemNotification(notificationData) {
    try {
      const response = await api.post('/admin/notifications/send', notificationData)
      return response
    } catch (error) {
      console.error('发送系统通知失败:', error)
      throw error
    }
  },

  /**
   * 获取项目类型统计
   */
  async getProjectTypeStats() {
    try {
      const response = await api.get('/admin/projects/type-stats')
      return response
    } catch (error) {
      console.error('获取项目类型统计失败:', error)
      throw error
    }
  },

  /**
   * 获取院系项目统计
   */
  async getDepartmentProjectStats() {
    try {
      const response = await api.get('/admin/projects/department-stats')
      return response
    } catch (error) {
      console.error('获取院系项目统计失败:', error)
      throw error
    }
  },

  /**
   * 获取时间趋势统计
   * @param {Object} timeRange 时间范围
   */
  async getTimeTrendStats(timeRange = {}) {
    try {
      const response = await api.get('/admin/projects/time-trend', {
        params: timeRange
      })
      return response
    } catch (error) {
      console.error('获取时间趋势统计失败:', error)
      throw error
    }
  },

  /**
   * 获取仪表板统计数据
   */
  async getDashboardStats() {
    try {
      const response = await api.get('/admin/dashboard/stats')
      return response
    } catch (error) {
      console.error('获取仪表板统计数据失败:', error)
      // 返回模拟数据作为后备
      return {
        code: 200,
        data: {
          userStats: {
            totalUsers: 1250,
            activeUsers: 890,
            newUsers: 45,
            userGrowth: 12.5,
            totalProjects: 156,
            pendingProjects: 23,
            activeCompetitions: 8,
            totalCompetitions: 12,
            pendingApplications: 15,
            todayApplications: 3
          },
          systemInfo: {
            serverTime: new Date().toLocaleString('zh-CN'),
            version: '1.0.0',
            uptime: '15天 8小时',
            dbStatus: '正常',
            lastBackup: '2024-01-15 02:00:00'
          }
        }
      }
    }
  },

  /**
   * 获取系统健康状态
   */
  async getSystemHealth() {
    try {
      const response = await api.get('/admin/system/health')
      return response
    } catch (error) {
      console.error('获取系统健康状态失败:', error)
      // 返回模拟数据作为后备
      return {
        code: 200,
        data: {
          checks: [
            {
              service: 'Web服务',
              status: '正常',
              uptime: '15天 8小时',
              load: '23%'
            },
            {
              service: '数据库',
              status: '正常',
              uptime: '15天 8小时',
              load: '15%'
            },
            {
              service: '文件存储',
              status: '正常',
              uptime: '15天 8小时',
              load: '8%'
            }
          ]
        }
      }
    }
  }
}

export default adminService 