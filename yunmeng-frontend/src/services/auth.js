import axios from 'axios'

// 创建axios实例
const api = axios.create({
  baseURL: 'http://localhost:8080/api', // 后端API地址
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器 - 添加token
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
    return response.data
  },
  (error) => {
    if (error.response?.status === 401) {
      // token过期，清除本地存储但不自动跳转到登录页
      localStorage.removeItem('token')
      localStorage.removeItem('userInfo')
      localStorage.removeItem('userRole')
      console.log('登录已过期，请重新登录')
      // 不再自动跳转到登录页面，让用户手动访问
    }
    return Promise.reject(error)
  }
)

// 登录接口
export const login = async (loginData) => {
  try {
    const response = await api.post('/login', loginData)
    // 兼容后端返回格式
    const { code, message, data, token } = response
    if (code === 200 && token) {
      return {
        success: true,
        data,
        token,
        message: message || '登录成功'
      }
    } else {
      return {
        success: false,
        code,
        message: message || '登录失败'
      }
    }
  } catch (error) {
    // 处理异常响应
    if (error.response) {
      const { code, message } = error.response.data || {}
      return {
        success: false,
        code: code || error.response.status,
        message: message || '登录失败，请检查用户名和密码'
      }
    } else {
      return {
        success: false,
        code: 500,
        message: '网络异常或服务器无响应'
      }
    }
  }
}

// 退出登录
export const logout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('userInfo')
  localStorage.removeItem('userRole')
}

// 获取用户信息
export const getUserInfo = async () => {
  try {
    const response = await api.get('/auth/user-info')
    return {
      success: true,
      data: response.data
    }
  } catch (error) {
    return {
      success: false,
      message: error.response?.data?.message || '获取用户信息失败'
    }
  }
}

// 获取token
export const getToken = () => {
  return localStorage.getItem('token')
}

// 检查登录状态
export const checkAuth = () => {
  const token = localStorage.getItem('token')
  const userInfo = localStorage.getItem('userInfo')
  const userRole = localStorage.getItem('userRole')
  
  return {
    isLoggedIn: !!token,
    userInfo: userInfo ? JSON.parse(userInfo) : null,
    userRole: userRole
  }
}

export default api 