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

// 项目分类服务类
class ProjectTypeService {
  // 获取项目分类列表
  async getProjectTypeList() {
    try {
      const response = await api.get('/project-types')
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取项目分类列表失败')
    }
  }

  // 获取项目分类统计
  async getProjectTypeStats() {
    try {
      const response = await api.get('/project-types/stats')
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取项目分类统计失败')
    }
  }

  // 获取项目分类详情
  async getProjectTypeDetail(id) {
    try {
      const response = await api.get(`/project-types/${id}`)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取项目分类详情失败')
    }
  }

  // 创建项目分类
  async createProjectType(data) {
    try {
      const response = await api.post('/project-types', data)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '创建项目分类失败')
    }
  }

  // 更新项目分类
  async updateProjectType(id, data) {
    try {
      const response = await api.put(`/project-types/${id}`, data)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '更新项目分类失败')
    }
  }

  // 删除项目分类
  async deleteProjectType(id) {
    try {
      const response = await api.delete(`/project-types/${id}`)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '删除项目分类失败')
    }
  }
}

export const projectTypeService = new ProjectTypeService() 