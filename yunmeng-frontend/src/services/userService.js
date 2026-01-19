import axios from 'axios'

// 创建axios实例
const api = axios.create({
  baseURL: 'http://localhost:8080/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器 - 添加认证token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器 - 处理错误
api.interceptors.response.use(
  (response) => {
    // 如果是文件下载，直接返回response
    if (response.config.responseType === 'blob') {
      return response.data
    }
    return response.data
  },
  (error) => {
    if (error.response) {
      // 服务器返回错误状态码
      const { status, data } = error.response
      
      switch (status) {
        case 401:
          // 未认证，清除token但不自动跳转
          localStorage.removeItem('token')
          localStorage.removeItem('user')
          console.log('登录已过期，请重新登录')
          // 不再自动跳转到登录页面，让用户手动访问
          break
        case 403:
          console.error('权限不足:', data.message)
          break
        case 404:
          console.error('资源不存在:', data.message)
          break
        case 500:
          console.error('服务器错误:', data.message)
          break
        default:
          console.error('请求失败:', data.message)
      }
      
      return Promise.reject(data)
    } else if (error.request) {
      // 网络错误
      console.error('网络错误:', error.message)
      return Promise.reject({ message: '网络连接失败，请检查网络设置' })
    } else {
      // 其他错误
      console.error('请求配置错误:', error.message)
      return Promise.reject({ message: '请求配置错误' })
    }
  }
)

// 用户管理API服务
export const userService = {
  // 获取用户列表
  getUserList: (params = {}) => {
    return api.get('/users', { params })
  },

  // 获取用户详情
  getUserById: (id) => {
    return api.get(`/users/${id}`)
  },

  // 创建用户
  createUser: (userData) => {
    return api.post('/users', userData)
  },

  // 更新用户
  updateUser: (id, userData) => {
    return api.put(`/users/${id}`, userData)
      .catch(error => {
        console.error('更新用户失败:', error)
        if (error.response?.status === 403) {
          throw new Error('权限不足，需要管理员权限')
        } else if (error.response?.status === 401) {
          throw new Error('登录已过期，请重新登录')
        } else if (error.response?.data?.message) {
          throw new Error(error.response.data.message)
        } else {
          throw new Error('更新用户失败，请稍后重试')
        }
      })
  },

  // 删除用户
  deleteUser: (id) => {
    return api.delete(`/users/${id}`)
  },

  // 切换用户状态
  toggleUserStatus: (id, status) => {
    return api.patch(`/users/${id}/status`, { status })
  },

  // 重置用户密码
  resetUserPassword: (id) => {
    return api.post(`/users/${id}/reset-password`)
  },

  // 批量删除用户
  batchDeleteUsers: (userIds) => {
    return api.post('/users/batch-delete', { userIds })
  },

  // 获取用户统计信息
  getUserStats: () => {
    return api.get('/users/stats')
  },

  // 导出用户数据
  exportUsers: (params = {}) => {
    return api.get('/users/export', { 
      params,
      responseType: 'blob' // 用于文件下载
    })
  },

  // 获取个人信息
  getProfile: () => {
    return api.get('/users/profile')
  },

  // 更新个人信息
  updateProfile: (profileData) => {
    return api.put('/users/profile', profileData)
  },

  // 修改密码
  changePassword: (passwordData) => {
    return api.put('/users/change-password', passwordData)
  },

  // 上传头像
  uploadAvatar: (file) => {
    const formData = new FormData()
    formData.append('avatar', file)
    return api.post('/users/avatar', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  // 获取学术成果
  getAchievements: (params = {}) => {
    return api.get('/users/achievements', { params })
  },

  // 添加学术成果
  addAchievement: (achievementData) => {
    return api.post('/users/achievements', achievementData)
  },

  // 更新学术成果
  updateAchievement: (id, achievementData) => {
    return api.put(`/users/achievements/${id}`, achievementData)
  },

  // 删除学术成果
  deleteAchievement: (id) => {
    return api.delete(`/users/achievements/${id}`)
  },

  // 导出个人档案
  exportProfile: () => {
    return api.get('/users/profile/export', {
      responseType: 'blob'
    })
  }
}

// 认证API服务
export const authService = {
  // 用户登录
  login: (credentials) => {
    return api.post('/login', credentials)
  },

  // 用户登出
  logout: () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    return Promise.resolve()
  },

  // 检查登录状态
  checkAuth: () => {
    const token = localStorage.getItem('token')
    const user = localStorage.getItem('user')
    return token && user ? JSON.parse(user) : null
  }
}

// 工具函数
export const apiUtils = {
  // 处理分页参数
  buildPaginationParams: (page, size, search, filters = {}) => {
    const params = {
      page: page || 1,
      size: size || 20
    }
    
    if (search) {
      params.search = search
    }
    
    // 添加筛选条件
    Object.keys(filters).forEach(key => {
      if (filters[key] !== null && filters[key] !== undefined && filters[key] !== '') {
        params[key] = filters[key]
      }
    })
    
    return params
  },

  // 处理API错误
  handleApiError: (error, defaultMessage = '操作失败') => {
    const message = error?.message || defaultMessage
    console.error('API错误:', error)
    return message
  },

  // 格式化用户数据
  formatUserData: (user) => {
    return {
      id: user.id,
      username: user.username,
      email: user.email,
      realName: user.realName || '',
      phone: user.phone || '',
      department: user.department || '',
      studentId: user.studentId || '',
      status: user.status,
      roles: user.roles || [],
      createTime: user.createTime,
      lastLogin: user.lastLogin
    }
  },

  // 格式化用户列表数据
  formatUserListData: (response) => {
    return {
      list: response.data.list || [],
      total: response.data.total || 0,
      page: response.data.page || 1,
      size: response.data.size || 20,
      pages: response.data.pages || 1
    }
  },

  // 格式化日期时间
  formatDateTime: (dateTime) => {
    if (!dateTime) return ''
    const date = new Date(dateTime)
    return date.toLocaleString('zh-CN')
  },

  // 验证手机号
  validatePhone: (phone) => {
    const phoneRegex = /^1[3-9]\d{9}$/
    return phoneRegex.test(phone)
  },

  // 验证邮箱
  validateEmail: (email) => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    return emailRegex.test(email)
  }
}

export default userService 