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

// 项目服务类
class ProjectService {
  // 创建项目
  async createProject(projectData) {
    try {
      const response = await api.post('/projects', projectData)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '创建项目失败')
    }
  }

  // 获取我的项目列表
  async getMyProjects(status = '') {
    try {
      const params = status ? { status } : {}
      const response = await api.get('/projects/my', { params })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取项目列表失败')
    }
  }

  

 // 项目详情
  async getProjectDetail(id) {
    return api.get('/projects/detail', {
      params: { id }
    })
  }
  
  // 项目成果文件列表
  async getProjectFiles(id) {
    return api.get('/projects/files/getfiles', {
      params: { id }
    })
  }

  // 上传项目成果文件
  async uploadProjectFiles(id, files) {
    const formData = new FormData()
    formData.append('id', id)
    files.forEach(f => formData.append('files', f))

    return api.post('/projects/files/uploadfiles', formData)
  }

  // 删除项目成果文件
  async deleteProjectFile(projectId, fileId) {
    return api.delete('/projects/files/delete', {
      params: {
        id: projectId,
        fileId
      }
    })
  }

// 获取项目文件列表
async getProjectFiles(projectId) {
  if (!projectId) throw new Error('项目ID不能为空')
  return api.get(`/projects/files/getfiles?id=${projectId}`)
}



  // 获取项目审核记录
  async getProjectReviews(projectId) {
    try {
      if (!projectId) {
        throw new Error('项目ID不能为空')
      }
      const response = await api.get(`/projects/${projectId}/reviews`)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取项目审核记录失败')
    }
  }


  // 获取文件类型配置
  async getFileTypeConfigs() {
    try {
      const response = await api.get('/projects/file-type-configs')
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取文件类型配置失败')
    }
  }



  // 获取教师项目列表
  async getTeacherProjects(params = {}) {
    try {
      const response = await api.get('/teachers/projects', { params })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取教师项目列表失败')
    }
  }

  // 获取教师列表
  async getTeacherList(params = {}) {
    try {
      const response = await api.get('/teachers', { params })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取教师列表失败')
    }
  }



  // 获取我的学生列表
  async getMyStudents(params = {}) {
    try {
      const response = await api.get('/teachers/students', { params })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取我的学生列表失败')
    }
  }


  // 获取学生的指导教师列表
  async getStudentTeachers(studentId) {
    try {
      const response = await api.get(`/teachers/students/${studentId}`)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取学生指导教师列表失败')
    }
  }

  // 提交项目审核
  async submitProject(id) {
    try {
      if (!id) {
        throw new Error('项目ID不能为空')
      }
        // 使用 FormData
        const formData = new FormData()
        formData.append('id', id)

      const response = await api.post(`/projects/submit`, formData)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '提交项目审核失败')
    }
  }


  // 获取项目统计信息
  async getProjectStats() {
    try {
      const response = await api.get('/admin/projects/stats')
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取项目统计信息失败')
    }
  }

  // 获取项目列表（管理员）
  async getProjectList(params = {}) {
    try {
      const response = await api.get('/teacher-projects', { params })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取项目列表失败')
    }
  }

  // 强制更新项目状态
  async forceUpdateProjectStatus(projectId, statusData) {
    try {
      if (!projectId) {
        throw new Error('项目ID不能为空')
      }
      const response = await api.put(`/admin/projects/${projectId}/force-status`, statusData)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '强制更新项目状态失败')
    }
  }



  // 恢复软删除的项目
  async restoreProject(projectId) {
    try {
      if (!projectId) {
        throw new Error('项目ID不能为空')
      }
      const response = await api.put(`/admin/projects/${projectId}/restore`)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '恢复项目失败')
    }
  }

  // 导出项目数据
  async exportProjects(exportData) {
    try {
      const response = await api.post('/admin/projects/export', exportData)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '导出项目数据失败')
    }
  }
}



export const projectService = new ProjectService()