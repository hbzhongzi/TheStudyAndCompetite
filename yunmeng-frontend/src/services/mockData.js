// 模拟数据服务
export const mockProjects = [
  {
    id: 1,
    title: '智能校园管理系统',
    description: '基于物联网技术的智能校园管理系统，实现校园设施的智能化管理',
    type: '科研',
    status: 'pending',
    createdAt: '2024-01-15T10:00:00Z',
    updatedAt: '2024-01-15T10:00:00Z'
  },
  {
    id: 2,
    title: '在线学习平台开发',
    description: '开发一个支持多种学习模式的在线教育平台',
    type: '竞赛',
    status: 'draft',
    createdAt: '2024-01-10T14:30:00Z',
    updatedAt: '2024-01-10T14:30:00Z'
  },
  {
    id: 3,
    title: '人工智能图像识别系统',
    description: '基于深度学习的图像识别系统，用于医疗影像分析',
    type: '科研',
    status: 'approved',
    createdAt: '2024-01-05T09:15:00Z',
    updatedAt: '2024-01-05T09:15:00Z'
  }
]

// 模拟API响应
export const createMockResponse = (data, message = '获取成功') => {
  return {
    code: 200,
    message,
    data
  }
}

// 模拟项目服务
export const mockProjectService = {
  async getMyProjects(status = '') {
    // 模拟网络延迟
    await new Promise(resolve => setTimeout(resolve, 500))
    
    let filteredProjects = mockProjects
    
    if (status) {
      filteredProjects = mockProjects.filter(project => project.status === status)
    }
    
    return createMockResponse(filteredProjects)
  },
  
  async getProjectDetail(projectId) {
    await new Promise(resolve => setTimeout(resolve, 300))
    
    const project = mockProjects.find(p => p.id === parseInt(projectId))
    if (project) {
      return createMockResponse(project)
    } else {
      throw new Error('项目不存在')
    }
  },
  
  async createProject(projectData) {
    await new Promise(resolve => setTimeout(resolve, 800))
    
    const newProject = {
      id: mockProjects.length + 1,
      ...projectData,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    }
    
    mockProjects.push(newProject)
    
    return createMockResponse({ projectId: newProject.id }, '项目创建成功')
  },
  
  async updateProject(projectId, projectData) {
    await new Promise(resolve => setTimeout(resolve, 600))
    
    const index = mockProjects.findIndex(p => p.id === parseInt(projectId))
    if (index !== -1) {
      mockProjects[index] = {
        ...mockProjects[index],
        ...projectData,
        updatedAt: new Date().toISOString()
      }
      return createMockResponse(null, '项目更新成功')
    } else {
      throw new Error('项目不存在')
    }
  },
  
  async deleteProject(projectId) {
    await new Promise(resolve => setTimeout(resolve, 400))
    
    const index = mockProjects.findIndex(p => p.id === parseInt(projectId))
    if (index !== -1) {
      mockProjects.splice(index, 1)
      return createMockResponse(null, '项目删除成功')
    } else {
      throw new Error('项目不存在')
    }
  }
} 