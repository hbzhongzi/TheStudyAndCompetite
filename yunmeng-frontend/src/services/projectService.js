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

  // 获取项目详情
  async getProjectDetail(projectId) {
    try {
      if (!projectId) {
        throw new Error('项目ID不能为空')
      }
      const response = await api.get(`/projects/${projectId}`)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取项目详情失败')
    }
  }

  // 更新项目
  async updateProject(projectId, projectData) {
    try {
      if (!projectId) {
        throw new Error('项目ID不能为空')
      }
      const response = await api.put(`/projects/${projectId}`, projectData)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '更新项目失败')
    }
  }

  // 删除项目
  async deleteProject(projectId) {
    try {
      if (!projectId) {
        throw new Error('项目ID不能为空')
      }
      const response = await api.delete(`/projects/${projectId}`)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '删除项目失败')
    }
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

  // =============================================
  // 1. 项目状态管理增强 API
  // =============================================

  // 更新项目状态
  async updateProjectStatus(projectId, statusData) {
    try {
      if (!projectId) {
        throw new Error('项目ID不能为空')
      }
      const response = await api.put(`/projects/${projectId}/status`, statusData)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '更新项目状态失败')
    }
  }

  // 获取项目状态变更历史
  async getProjectStatusHistory(projectId) {
    try {
      if (!projectId) {
        throw new Error('项目ID不能为空')
      }
      const response = await api.get(`/projects/${projectId}/status-history`)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取项目状态历史失败')
    }
  }

  // =============================================
  // 2. 项目生命周期管理增强 API
  // =============================================

  // 创建项目里程碑
  async createProjectMilestone(projectId, milestoneData) {
    try {
      if (!projectId) {
        throw new Error('项目ID不能为空')
      }
      const response = await api.post(`/projects/${projectId}/milestones`, milestoneData)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '创建项目里程碑失败')
    }
  }

  // 更新项目里程碑
  async updateProjectMilestone(milestoneId, milestoneData) {
    try {
      const response = await api.put(`/projects/milestones/${milestoneId}`, milestoneData)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '更新项目里程碑失败')
    }
  }

  // 获取项目里程碑列表
  async getProjectMilestones(projectId) {
    try {
      if (!projectId) {
        throw new Error('项目ID不能为空')
      }
      const response = await api.get(`/projects/${projectId}/milestones`)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取项目里程碑失败')
    }
  }

  // 申请项目延期
  async applyProjectExtension(projectId, extensionData) {
    try {
      if (!projectId) {
        throw new Error('项目ID不能为空')
      }
      const response = await api.post(`/projects/${projectId}/extensions`, extensionData)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '申请项目延期失败')
    }
  }

  // 审核项目延期申请
  async reviewProjectExtension(extensionId, reviewData) {
    try {
      const response = await api.put(`/projects/extensions/${extensionId}/review`, reviewData)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '审核项目延期申请失败')
    }
  }

  // 更新项目进度
  async updateProjectProgress(projectId, progressData) {
    try {
      if (!projectId) {
        throw new Error('项目ID不能为空')
      }
      const response = await api.put(`/projects/${projectId}/progress`, progressData)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '更新项目进度失败')
    }
  }

  // =============================================
  // 3. 成果文件管理增强 API
  // =============================================

  // 上传项目文件（增强版）
  async uploadProjectFile(projectId, fileData) {
    try {
      if (!projectId) {
        throw new Error('项目ID不能为空')
      }
      const response = await api.post(`/projects/${projectId}/files`, fileData)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '上传项目文件失败')
    }
  }

  // 审核项目文件
  async reviewProjectFile(fileId, reviewData) {
    try {
      const response = await api.put(`/projects/files/${fileId}/review`, reviewData)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '审核项目文件失败')
    }
  }

  // 按类型获取项目文件
  async getProjectFilesByType(projectId, fileType = '') {
    try {
      if (!projectId) {
        throw new Error('项目ID不能为空')
      }
      const params = fileType ? { type: fileType } : {}
      const response = await api.get(`/projects/${projectId}/files`, { params })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取项目文件失败')
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

  // =============================================
  // 4. 项目分类管理增强 API
  // =============================================

  // 创建项目分类
  async createProjectType(typeData) {
    try {
      const response = await api.post('/admin/projects/types', typeData)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '创建项目分类失败')
    }
  }

  // 更新项目分类
  async updateProjectType(typeId, typeData) {
    try {
      const response = await api.put(`/admin/projects/types/${typeId}`, typeData)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '更新项目分类失败')
    }
  }

  // 获取项目分类树
  async getProjectTypeTree() {
    try {
      const response = await api.get('/admin/projects/types/tree')
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取项目分类树失败')
    }
  }

  // 获取项目分类统计
  async getProjectTypeStats() {
    try {
      const response = await api.get('/admin/projects/types/stats')
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取项目分类统计失败')
    }
  }

  // =============================================
  // 5. 审核流程增强 API
  // =============================================

  // 创建审核流程配置
  async createReviewFlow(flowData) {
    try {
      const response = await api.post('/admin/projects/review-flows', flowData)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '创建审核流程失败')
    }
  }

  // 委托审核
  async delegateReview(reviewId, delegationData) {
    try {
      const response = await api.post(`/projects/reviews/${reviewId}/delegate`, delegationData)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '委托审核失败')
    }
  }

  // 获取我的审核任务
  async getMyReviewTasks(params = {}) {
    try {
      const response = await api.get('/projects/my-review-tasks', { params })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取审核任务失败')
    }
  }

  // 获取审核流程配置
  async getReviewFlowConfig(projectTypeId = null) {
    try {
      const params = projectTypeId ? { projectTypeId } : {}
      const response = await api.get('/projects/review-flow-config', { params })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取审核流程配置失败')
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

  // 绑定学生和教师
  async bindStudentTeacher(bindData) {
    try {
      const response = await api.post('/teachers/bind-student', bindData)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '绑定学生和教师失败')
    }
  }

  // 解绑学生和教师
  async unbindStudentTeacher(studentId, teacherId) {
    try {
      const response = await api.delete(`/teachers/students/${studentId}/teachers/${teacherId}`)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '解绑学生和教师失败')
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

  // 学生绑定教师
  async bindStudentToTeacher(teacherId) {
    try {
      const response = await api.post('/students/bind-teacher', { teacherId })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '绑定教师失败')
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
  async submitProject(projectId) {
    try {
      if (!projectId) {
        throw new Error('项目ID不能为空')
      }
      const response = await api.post(`/projects/submit/${projectId}`)
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

  // 软删除项目
  async softDeleteProject(projectId) {
    try {
      if (!projectId) {
        throw new Error('项目ID不能为空')
      }
      const response = await api.delete(`/admin/projects/${projectId}/soft`)
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '软删除项目失败')
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

// 文件上传服务类
class FileService {
  // 上传文件
  async uploadFile(file) {
    try {
      const formData = new FormData()
      formData.append('file', file)
      
      const response = await api.post('/files/upload', formData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      })
      return response
    } catch (error) {
      throw new Error(error.response?.data?.message || '文件上传失败')
    }
  }
}

export const projectService = new ProjectService()
export const fileService = new FileService() 