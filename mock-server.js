const express = require('express')
const cors = require('cors')
const bodyParser = require('body-parser')

const app = express()
const PORT = 8080

// 中间件
app.use(cors())
app.use(bodyParser.json())

// 模拟用户数据
let users = [
  {
    id: 1,
    username: 'admin',
    realName: '系统管理员',
    email: 'admin@yunmeng.edu.cn',
    role: 'admin',
    status: 'active',
    createTime: '2024-01-01T00:00:00.000Z',
    lastLogin: '2024-01-15T14:30:00.000Z',
    phone: '13800138000',
    department: '信息技术部'
  },
  {
    id: 2,
    username: 'teacher001',
    realName: '李老师',
    email: 'li.teacher@yunmeng.edu.cn',
    role: 'teacher',
    status: 'active',
    createTime: '2024-01-02T10:00:00.000Z',
    lastLogin: '2024-01-15T13:45:00.000Z',
    phone: '13800138001',
    department: '计算机学院'
  },
  {
    id: 3,
    username: 'student001',
    realName: '张三',
    email: 'zhangsan@yunmeng.edu.cn',
    role: 'student',
    status: 'active',
    createTime: '2024-01-03T14:20:00.000Z',
    lastLogin: '2024-01-15T12:20:00.000Z',
    phone: '13800138002',
    department: '计算机学院',
    studentId: '2021001'
  }
]

let nextId = 4

// 通用响应格式
const createResponse = (code, message, data = null) => {
  return {
    code,
    message,
    data,
    timestamp: new Date().toISOString()
  }
}

// 分页响应格式
const createPageResponse = (list, total, page, size) => {
  return {
    code: 200,
    message: '操作成功',
    data: {
      list,
      total,
      page,
      size,
      pages: Math.ceil(total / size)
    },
    timestamp: new Date().toISOString()
  }
}

// 获取用户列表
app.get('/api/users', (req, res) => {
  try {
    const { page = 1, size = 20, search = '', role = '', status = '' } = req.query
    
    let filteredUsers = [...users]
    
    // 搜索过滤
    if (search) {
      const searchLower = search.toLowerCase()
      filteredUsers = filteredUsers.filter(user => 
        user.username.toLowerCase().includes(searchLower) ||
        user.realName.toLowerCase().includes(searchLower) ||
        user.email.toLowerCase().includes(searchLower)
      )
    }
    
    // 角色过滤
    if (role) {
      filteredUsers = filteredUsers.filter(user => user.role === role)
    }
    
    // 状态过滤
    if (status) {
      filteredUsers = filteredUsers.filter(user => user.status === status)
    }
    
    const total = filteredUsers.length
    const startIndex = (page - 1) * size
    const endIndex = startIndex + parseInt(size)
    const paginatedUsers = filteredUsers.slice(startIndex, endIndex)
    
    res.json(createPageResponse(paginatedUsers, total, parseInt(page), parseInt(size)))
  } catch (error) {
    res.status(500).json(createResponse(500, '服务器内部错误'))
  }
})

// 获取用户详情
app.get('/api/users/:id', (req, res) => {
  try {
    const id = parseInt(req.params.id)
    const user = users.find(u => u.id === id)
    
    if (!user) {
      return res.status(404).json(createResponse(404, '用户不存在'))
    }
    
    res.json(createResponse(200, '获取用户详情成功', user))
  } catch (error) {
    res.status(500).json(createResponse(500, '服务器内部错误'))
  }
})

// 创建用户
app.post('/api/users', (req, res) => {
  try {
    const { username, realName, email, password, role, phone, department, studentId } = req.body
    
    // 验证必填字段
    if (!username || !realName || !email || !password || !role) {
      return res.status(400).json(createResponse(400, '缺少必填字段'))
    }
    
    // 检查用户名是否已存在
    if (users.find(u => u.username === username)) {
      return res.status(409).json(createResponse(409, '用户名已存在'))
    }
    
    // 检查邮箱是否已存在
    if (users.find(u => u.email === email)) {
      return res.status(409).json(createResponse(409, '邮箱已存在'))
    }
    
    const newUser = {
      id: nextId++,
      username,
      realName,
      email,
      role,
      status: 'active',
      createTime: new Date().toISOString(),
      lastLogin: null,
      phone: phone || '',
      department: department || '',
      studentId: studentId || ''
    }
    
    users.push(newUser)
    
    res.status(201).json(createResponse(201, '用户创建成功', newUser))
  } catch (error) {
    res.status(500).json(createResponse(500, '服务器内部错误'))
  }
})

// 更新用户信息
app.put('/api/users/:id', (req, res) => {
  try {
    const id = parseInt(req.params.id)
    const { realName, email, role, phone, department, studentId } = req.body
    
    const userIndex = users.findIndex(u => u.id === id)
    if (userIndex === -1) {
      return res.status(404).json(createResponse(404, '用户不存在'))
    }
    
    // 检查邮箱是否已被其他用户使用
    if (email && users.find(u => u.email === email && u.id !== id)) {
      return res.status(409).json(createResponse(409, '邮箱已被其他用户使用'))
    }
    
    // 更新用户信息
    users[userIndex] = {
      ...users[userIndex],
      ...(realName && { realName }),
      ...(email && { email }),
      ...(role && { role }),
      ...(phone !== undefined && { phone }),
      ...(department !== undefined && { department }),
      ...(studentId !== undefined && { studentId }),
      updateTime: new Date().toISOString()
    }
    
    res.json(createResponse(200, '用户信息更新成功', users[userIndex]))
  } catch (error) {
    res.status(500).json(createResponse(500, '服务器内部错误'))
  }
})

// 删除用户
app.delete('/api/users/:id', (req, res) => {
  try {
    const id = parseInt(req.params.id)
    const userIndex = users.findIndex(u => u.id === id)
    
    if (userIndex === -1) {
      return res.status(404).json(createResponse(404, '用户不存在'))
    }
    
    users.splice(userIndex, 1)
    
    res.json(createResponse(200, '用户删除成功'))
  } catch (error) {
    res.status(500).json(createResponse(500, '服务器内部错误'))
  }
})

// 启用/禁用用户
app.patch('/api/users/:id/status', (req, res) => {
  try {
    const id = parseInt(req.params.id)
    const { status } = req.body
    
    if (!status || !['active', 'inactive'].includes(status)) {
      return res.status(400).json(createResponse(400, '状态值无效'))
    }
    
    const userIndex = users.findIndex(u => u.id === id)
    if (userIndex === -1) {
      return res.status(404).json(createResponse(404, '用户不存在'))
    }
    
    users[userIndex].status = status
    users[userIndex].updateTime = new Date().toISOString()
    
    res.json(createResponse(200, '用户状态更新成功', {
      id: users[userIndex].id,
      status: users[userIndex].status,
      updateTime: users[userIndex].updateTime
    }))
  } catch (error) {
    res.status(500).json(createResponse(500, '服务器内部错误'))
  }
})

// 重置用户密码
app.post('/api/users/:id/reset-password', (req, res) => {
  try {
    const id = parseInt(req.params.id)
    const user = users.find(u => u.id === id)
    
    if (!user) {
      return res.status(404).json(createResponse(404, '用户不存在'))
    }
    
    // 模拟密码重置逻辑
    const resetTime = new Date().toISOString()
    
    res.json(createResponse(200, '密码重置成功，新密码已发送到用户邮箱', {
      id: user.id,
      resetTime
    }))
  } catch (error) {
    res.status(500).json(createResponse(500, '服务器内部错误'))
  }
})

// 批量删除用户
app.post('/api/users/batch-delete', (req, res) => {
  try {
    const { userIds } = req.body
    
    if (!userIds || !Array.isArray(userIds)) {
      return res.status(400).json(createResponse(400, '用户ID列表无效'))
    }
    
    const deletedIds = []
    const originalLength = users.length
    
    users = users.filter(user => {
      if (userIds.includes(user.id)) {
        deletedIds.push(user.id)
        return false
      }
      return true
    })
    
    const deletedCount = originalLength - users.length
    
    res.json(createResponse(200, '批量删除成功', {
      deletedCount,
      deletedIds
    }))
  } catch (error) {
    res.status(500).json(createResponse(500, '服务器内部错误'))
  }
})

// 获取用户统计信息
app.get('/api/users/stats', (req, res) => {
  try {
    const totalUsers = users.length
    const activeUsers = users.filter(u => u.status === 'active').length
    const inactiveUsers = totalUsers - activeUsers
    
    const roleStats = {
      admin: users.filter(u => u.role === 'admin').length,
      teacher: users.filter(u => u.role === 'teacher').length,
      student: users.filter(u => u.role === 'student').length
    }
    
    const departmentStats = {}
    users.forEach(user => {
      if (user.department) {
        departmentStats[user.department] = (departmentStats[user.department] || 0) + 1
      }
    })
    
    res.json(createResponse(200, '获取统计信息成功', {
      totalUsers,
      activeUsers,
      inactiveUsers,
      roleStats,
      departmentStats
    }))
  } catch (error) {
    res.status(500).json(createResponse(500, '服务器内部错误'))
  }
})

// 启动服务器
app.listen(PORT, () => {
  console.log(`模拟服务器运行在 http://localhost:${PORT}`)
  console.log('用户管理API已就绪')
  console.log('可用接口:')
  console.log('  GET    /api/users - 获取用户列表')
  console.log('  GET    /api/users/:id - 获取用户详情')
  console.log('  POST   /api/users - 创建用户')
  console.log('  PUT    /api/users/:id - 更新用户')
  console.log('  DELETE /api/users/:id - 删除用户')
  console.log('  PATCH  /api/users/:id/status - 更新用户状态')
  console.log('  POST   /api/users/:id/reset-password - 重置密码')
  console.log('  POST   /api/users/batch-delete - 批量删除')
  console.log('  GET    /api/users/stats - 获取统计信息')
}) 