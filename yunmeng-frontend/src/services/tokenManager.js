import axios from 'axios'

// Token管理器
class TokenManager {
  constructor() {
    this.baseURL = 'http://localhost:8080/api'
    this.isRefreshing = false
    this.failedQueue = []
  }

  // 检查Token是否即将过期（提前5分钟刷新）
  isTokenExpiringSoon() {
    const token = localStorage.getItem('token')
    if (!token) return true

    try {
      const parts = token.split('.')
      if (parts.length !== 3) return true

      const payload = JSON.parse(atob(parts[1]))
      if (!payload.exp) return true

      const expTime = payload.exp * 1000
      const now = Date.now()
      const fiveMinutes = 5 * 60 * 1000

      return (expTime - now) < fiveMinutes
    } catch (error) {
      console.error('解析Token失败:', error)
      return true
    }
  }

  // 刷新Token
  async refreshToken() {
    if (this.isRefreshing) {
      return new Promise((resolve, reject) => {
        this.failedQueue.push({ resolve, reject })
      })
    }

    this.isRefreshing = true

    try {
      const currentToken = localStorage.getItem('token')
      if (!currentToken) {
        throw new Error('没有可用的Token')
      }

      // 调用后端的Token刷新接口
      const response = await axios.post(`${this.baseURL}/refresh-token`, {}, {
        headers: {
          'Authorization': `Bearer ${currentToken}`
        }
      })

      if (response.data.code === 200 && response.data.token) {
        localStorage.setItem('token', response.data.token)
        this.processQueue(null, response.data.token)
        console.log('Token刷新成功')
        return response.data.token
      } else {
        throw new Error(response.data.message || '刷新Token失败')
      }
    } catch (error) {
      console.error('刷新Token失败:', error)
      this.processQueue(error, null)
      
      // 清除过期的认证信息
      localStorage.removeItem('token')
      localStorage.removeItem('userInfo')
      localStorage.removeItem('userRole')
      
      throw error
    } finally {
      this.isRefreshing = false
    }
  }

  // 处理队列中的请求
  processQueue(error, token = null) {
    this.failedQueue.forEach(({ resolve, reject }) => {
      if (error) {
        reject(error)
      } else {
        resolve(token)
      }
    })
    this.failedQueue = []
  }

  // 获取有效的Token
  async getValidToken() {
    if (this.isTokenExpiringSoon()) {
      try {
        return await this.refreshToken()
      } catch (error) {
        // 如果刷新失败，清除认证信息并抛出错误
        localStorage.clear()
        throw new Error('Token已过期，请重新登录')
      }
    }
    return localStorage.getItem('token')
  }

  // 清除认证信息
  clearAuth() {
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
    localStorage.removeItem('userRole')
  }

  // 检查认证状态
  checkAuthStatus() {
    const token = localStorage.getItem('token')
    const userInfo = localStorage.getItem('userInfo')
    const userRole = localStorage.getItem('userRole')

    return {
      hasToken: !!token,
      hasUserInfo: !!userInfo,
      hasUserRole: !!userRole,
      isAdmin: userRole === 'admin',
      isExpiringSoon: this.isTokenExpiringSoon()
    }
  }
}

// 创建单例实例
const tokenManager = new TokenManager()

export default tokenManager 