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
      localStorage.removeItem('token')
      localStorage.removeItem('userInfo')
      localStorage.removeItem('userRole')
      console.log('登录已过期，请重新登录')
    }
    return Promise.reject(error)
  }
)

// 通知服务类
class NotificationService {
  // 获取我的通知列表
  async getMyNotifications(params = {}) {
    try {
      const response = await api.get('/notifications', { params })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取通知列表失败')
    }
  }

  // 标记通知为已读
  async markNotificationAsRead(notificationId) {
    try {
      const response = await api.put(`/notifications/${notificationId}/read`)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '标记通知为已读失败')
    }
  }

  // 标记所有通知为已读
  async markAllNotificationsAsRead() {
    try {
      const response = await api.put('/notifications/read-all')
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '标记所有通知为已读失败')
    }
  }

  // 获取未读通知数量
  async getUnreadCount() {
    try {
      const response = await api.get('/notifications/unread-count')
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取未读通知数量失败')
    }
  }

  // 删除通知
  async deleteNotification(notificationId) {
    try {
      const response = await api.delete(`/notifications/${notificationId}`)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '删除通知失败')
    }
  }

  // 获取通知模板列表（管理员）
  async getNotificationTemplates() {
    try {
      const response = await api.get('/admin/notifications/templates')
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取通知模板失败')
    }
  }

  // 更新通知模板（管理员）
  async updateNotificationTemplate(templateId, templateData) {
    try {
      const response = await api.put(`/admin/notifications/templates/${templateId}`, templateData)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '更新通知模板失败')
    }
  }

  // 发送通知（管理员/系统）
  async sendNotification(notificationData) {
    try {
      const response = await api.post('/admin/notifications/send', notificationData)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '发送通知失败')
    }
  }

  // 获取通知统计信息
  async getNotificationStats() {
    try {
      const response = await api.get('/admin/notifications/stats')
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取通知统计失败')
    }
  }

  // 批量发送通知
  async batchSendNotifications(notifications) {
    try {
      const response = await api.post('/admin/notifications/batch-send', { notifications })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '批量发送通知失败')
    }
  }

  // 获取通知发送历史
  async getNotificationHistory(params = {}) {
    try {
      const response = await api.get('/admin/notifications/history', { params })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取通知历史失败')
    }
  }

  // 设置通知偏好
  async setNotificationPreferences(preferences) {
    try {
      const response = await api.put('/notifications/preferences', preferences)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '设置通知偏好失败')
    }
  }

  // 获取通知偏好
  async getNotificationPreferences() {
    try {
      const response = await api.get('/notifications/preferences')
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取通知偏好失败')
    }
  }

  // 订阅通知频道
  async subscribeNotificationChannel(channelId) {
    try {
      const response = await api.post(`/notifications/channels/${channelId}/subscribe`)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '订阅通知频道失败')
    }
  }

  // 取消订阅通知频道
  async unsubscribeNotificationChannel(channelId) {
    try {
      const response = await api.delete(`/notifications/channels/${channelId}/subscribe`)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '取消订阅通知频道失败')
    }
  }

  // 获取可用的通知频道
  async getAvailableNotificationChannels() {
    try {
      const response = await api.get('/notifications/channels')
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取通知频道失败')
    }
  }

  // 测试通知发送
  async testNotification(templateId, testData) {
    try {
      const response = await api.post(`/admin/notifications/templates/${templateId}/test`, testData)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '测试通知发送失败')
    }
  }
}

export default new NotificationService() 