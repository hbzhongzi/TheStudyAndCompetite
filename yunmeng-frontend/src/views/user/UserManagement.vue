<template>
  <div class="user-management">
    <div class="page-header">
      <h2>用户管理</h2>
      <div class="header-actions">
        <el-button type="success" @click="exportUsers" :loading="exporting">
          <i class="el-icon-download"></i>
          导出用户
        </el-button>
        <el-button type="primary" @click="showAddUserDialog = true">
          <i class="el-icon-plus"></i>
          添加用户
        </el-button>
      </div>
    </div>

    <!-- 统计信息卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <el-card class="stats-card">
          <div class="stats-content">
            <div class="stats-number">{{ stats.totalUsers }}</div>
            <div class="stats-label">总用户数</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stats-card">
          <div class="stats-content">
            <div class="stats-number">{{ stats.activeUsers }}</div>
            <div class="stats-label">活跃用户</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stats-card">
          <div class="stats-content">
            <div class="stats-number">{{ stats.newUsersToday }}</div>
            <div class="stats-label">今日新增</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stats-card">
          <div class="stats-content">
            <div class="stats-number">{{ stats.onlineUsers }}</div>
            <div class="stats-label">在线用户</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 搜索和筛选 -->
    <el-card class="search-card">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-input
            v-model="searchQuery"
            placeholder="搜索用户名、姓名或邮箱"
            clearable
            @input="handleSearch"
          >
            <template #prefix>
              <i class="el-icon-search"></i>
            </template>
          </el-input>
        </el-col>
        <el-col :span="4">
          <el-select v-model="roleFilter" placeholder="角色筛选" clearable @change="handleSearch">
            <el-option label="全部角色" value="" />
            <el-option label="学生" value="student" />
            <el-option label="教师" value="teacher" />
            <el-option label="管理员" value="admin" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-select v-model="statusFilter" placeholder="状态筛选" clearable @change="handleSearch">
            <el-option label="全部状态" value="" />
            <el-option label="活跃" value="active" />
            <el-option label="禁用" value="inactive" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 用户列表 -->
    <el-card class="user-list-card">
      <div class="table-header">
        <div class="table-title">用户列表</div>
        <div class="table-actions">
          <el-button 
            type="danger" 
            :disabled="selectedUsers.length === 0"
            @click="batchDeleteUsers"
          >
            批量删除 ({{ selectedUsers.length }})
          </el-button>
        </div>
      </div>

      <el-table 
        :data="safeUsers" 
        style="width: 100%" 
        v-loading="loading"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="realName" label="姓名" width="120" />
        <el-table-column prop="email" label="邮箱" width="200" />
        <el-table-column prop="phone" label="电话" width="130" />
        <el-table-column prop="department" label="部门/学院" width="150" />
        <el-table-column prop="role" label="角色" width="100">
          <template #default="scope">
            <div v-if="scope.row.roleNames && scope.row.roleNames.length > 0">
              <el-tag 
                v-for="roleName in scope.row.roleNames" 
                :key="roleName"
                :type="getRoleTagType(getRoleKeyFromName(roleName))"
                style="margin-right: 5px;"
              >
                {{ roleName }}
              </el-tag>
            </div>
            <span v-else>未分配角色</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'active' ? 'success' : 'danger'">
              {{ scope.row.status === 'active' ? '活跃' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createTime" label="创建时间" width="180">
          <template #default="scope">
            {{ formatDateTime(scope.row.createTime) }}
          </template>
        </el-table-column>
        <el-table-column prop="lastLogin" label="最后登录" width="180">
          <template #default="scope">
            {{ scope.row.lastLogin ? formatDateTime(scope.row.lastLogin) : '从未登录' }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="scope">
            <el-button size="small" @click="viewUser(scope.row)">查看</el-button>
            <el-button size="small" type="primary" @click="editUser(scope.row)">编辑</el-button>
            <el-button 
              size="small" 
              :type="scope.row.status === 'active' ? 'danger' : 'success'"
              @click="toggleUserStatus(scope.row)"
            >
              {{ scope.row.status === 'active' ? '禁用' : '启用' }}
            </el-button>
            <el-button size="small" type="warning" @click="resetPassword(scope.row)">重置密码</el-button>
            <el-button size="small" type="danger" @click="deleteUser(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="totalUsers"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 添加用户对话框 -->
    <el-dialog
      v-model="showAddUserDialog"
      title="添加用户"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form :model="userForm" :rules="userRules" ref="userFormRef" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="用户名" prop="username">
              <el-input v-model="userForm.username" placeholder="请输入用户名" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="姓名" prop="realName">
              <el-input v-model="userForm.realName" placeholder="请输入真实姓名" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="userForm.email" placeholder="请输入邮箱地址" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="电话" prop="phone">
              <el-input v-model="userForm.phone" placeholder="请输入电话号码" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="部门/学院" prop="department">
              <el-input v-model="userForm.department" placeholder="请输入部门或学院" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="学号/工号" prop="studentId">
              <el-input v-model="userForm.studentId" placeholder="请输入学号或工号" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="角色" prop="role">
              <el-select v-model="userForm.role" placeholder="请选择角色" style="width: 100%">
                <el-option label="学生" value="student" />
                <el-option label="教师" value="teacher" />
                <el-option label="管理员" value="admin" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-select v-model="userForm.status" placeholder="请选择状态" style="width: 100%">
                <el-option label="活跃" value="active" />
                <el-option label="禁用" value="inactive" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="初始密码" prop="password">
              <el-input v-model="userForm.password" type="password" placeholder="请输入初始密码" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="确认密码" prop="confirmPassword">
              <el-input v-model="userForm.confirmPassword" type="password" placeholder="请确认密码" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showAddUserDialog = false">取消</el-button>
          <el-button type="primary" @click="addUser" :loading="submitting">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 编辑用户对话框 -->
    <el-dialog
      v-model="showEditUserDialog"
      title="编辑用户"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form :model="editForm" :rules="editRules" ref="editFormRef" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="用户名" prop="username">
              <el-input v-model="editForm.username" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="姓名" prop="realName">
              <el-input v-model="editForm.realName" placeholder="请输入真实姓名" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="editForm.email" placeholder="请输入邮箱地址" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="电话" prop="phone">
              <el-input v-model="editForm.phone" placeholder="请输入电话号码" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="部门/学院" prop="department">
              <el-input v-model="editForm.department" placeholder="请输入部门或学院" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="学号/工号" prop="studentId">
              <el-input v-model="editForm.studentId" placeholder="请输入学号或工号" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="角色" prop="role">
              <el-select v-model="editForm.role" placeholder="请选择角色" style="width: 100%">
                <el-option label="学生" value="student" />
                <el-option label="教师" value="teacher" />
                <el-option label="管理员" value="admin" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-select v-model="editForm.status" placeholder="请选择状态" style="width: 100%">
                <el-option label="活跃" value="active" />
                <el-option label="禁用" value="inactive" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showEditUserDialog = false">取消</el-button>
          <el-button type="primary" @click="updateUser" :loading="submitting">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 用户详情对话框 -->
    <el-dialog
      v-model="showUserDetailDialog"
      title="用户详情"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-descriptions :column="2" border>
        <el-descriptions-item label="用户ID">{{ currentUser.id }}</el-descriptions-item>
        <el-descriptions-item label="用户名">{{ currentUser.username }}</el-descriptions-item>
        <el-descriptions-item label="真实姓名">{{ currentUser.realName }}</el-descriptions-item>
        <el-descriptions-item label="邮箱">{{ currentUser.email }}</el-descriptions-item>
        <el-descriptions-item label="电话">{{ currentUser.phone || '未设置' }}</el-descriptions-item>
        <el-descriptions-item label="部门/学院">{{ currentUser.department || '未设置' }}</el-descriptions-item>
        <el-descriptions-item label="学号/工号">{{ currentUser.studentId || '未设置' }}</el-descriptions-item>
        <el-descriptions-item label="角色">
          <div v-if="currentUser.roles && currentUser.roles.length > 0">
            <el-tag 
              v-for="role in currentUser.roles" 
              :key="role.id"
              :type="getRoleTagType(role.roleKey)"
              style="margin-right: 5px;"
            >
              {{ role.roleName }}
            </el-tag>
          </div>
          <span v-else>未分配角色</span>
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="currentUser.status === 'active' ? 'success' : 'danger'">
            {{ currentUser.status === 'active' ? '活跃' : '禁用' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDateTime(currentUser.createTime) }}</el-descriptions-item>
        <el-descriptions-item label="最后登录">
          {{ currentUser.lastLogin ? formatDateTime(currentUser.lastLogin) : '从未登录' }}
        </el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showUserDetailDialog = false">关闭</el-button>
          <el-button type="primary" @click="editUser(currentUser)">编辑用户</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import userService from '@/services/userService'
import { ensureArray, safeSelectionChange } from '@/utils/dataValidator'

// 响应式数据
const loading = ref(false)
const submitting = ref(false)
const exporting = ref(false)
const searchQuery = ref('')
const roleFilter = ref('')
const statusFilter = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const totalUsers = ref(0)
const showAddUserDialog = ref(false)
const showEditUserDialog = ref(false)
const showUserDetailDialog = ref(false)
const selectedUsers = ref([])
const currentUser = ref({})

// 统计信息
const stats = ref({
  totalUsers: 0,
  activeUsers: 0,
  newUsersToday: 0,
  onlineUsers: 0
})

// 表单引用
const userFormRef = ref()
const editFormRef = ref()

// 用户表单
const userForm = reactive({
  username: '',
  realName: '',
  email: '',
  phone: '',
  department: '',
  studentId: '',
  role: '',
  status: 'active',
  password: '',
  confirmPassword: ''
})

// 编辑表单
const editForm = reactive({
  id: '',
  username: '',
  realName: '',
  email: '',
  phone: '',
  department: '',
  studentId: '',
  role: '',
  status: ''
})

// 表单验证规则
const userRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  realName: [
    { required: true, message: '请输入真实姓名', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' }
  ],
  status: [
    { required: true, message: '请选择状态', trigger: 'change' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于 6 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== userForm.password) {
          callback(new Error('两次输入密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

const editRules = {
  realName: [
    { required: true, message: '请输入真实姓名', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' }
  ],
  status: [
    { required: true, message: '请选择状态', trigger: 'change' }
  ]
}

// 用户数据
const users = ref([])

// 确保 users 始终是一个数组
const safeUsers = computed(() => {
  return Array.isArray(users.value) ? users.value : []
})

// 安全的选择变化处理函数
const handleSelectionChange = (selection) => {
  safeSelectionChange(selection, (value) => {
    selectedUsers.value = value
  }, [])
}

// 方法
const getRoleLabel = (role) => {
  const roleMap = {
    admin: '管理员',
    teacher: '教师',
    student: '学生'
  }
  return roleMap[role] || role
}

// 根据角色名称获取角色键值
const getRoleKeyFromName = (roleName) => {
  const roleKeyMap = {
    '管理员': 'admin',
    '教师': 'teacher',
    '学生': 'student'
  }
  return roleKeyMap[roleName] || 'info'
}

const getRoleTagType = (role) => {
  const typeMap = {
    admin: 'danger',
    teacher: 'warning',
    student: 'success'
  }
  return typeMap[role] || 'info'
}

// 获取用户角色的辅助函数
const getUserRole = (user) => {
  // 处理用户详情对话框中的角色数据（roles数组，每个元素包含roleKey和roleName）
  if (Array.isArray(user.roles) && user.roles.length > 0) {
    // 如果是角色对象数组，返回第一个角色的roleKey
    if (user.roles[0] && typeof user.roles[0] === 'object' && user.roles[0].roleKey) {
      return user.roles[0].roleKey
    }
    // 如果是字符串数组，返回第一个角色
    return user.roles[0]
  }
  // 处理用户列表中的角色数据（roleNames数组或role字符串）
  if (Array.isArray(user.roleNames) && user.roleNames.length > 0) {
    return user.roleNames[0]
  } else if (Array.isArray(user.role)) {
    return user.role[0]
  }
  return user.role || ''
}

const formatDateTime = (dateTime) => {
  if (!dateTime) return ''
  const date = new Date(dateTime)
  return date.toLocaleString('zh-CN')
}

// 加载用户数据
const loadUsers = async () => {
  try {
    loading.value = true
    const params = {
      page: currentPage.value,
      size: pageSize.value,
      search: searchQuery.value,
      role: roleFilter.value,
      status: statusFilter.value
    }
    
    const response = await userService.getUserList(params)
    users.value = response.data.list || []
    totalUsers.value = response.data.total || 0
  } catch (error) {
    ElMessage.error('加载用户数据失败')
    console.error('加载用户数据失败:', error)
  } finally {
    loading.value = false
  }
}

// 加载统计信息
const loadStats = async () => {
  try {
    const response = await userService.getUserStats()
    stats.value = response.data || {
      totalUsers: 0,
      activeUsers: 0,
      newUsersToday: 0,
      onlineUsers: 0
    }
  } catch (error) {
    console.error('加载统计信息失败:', error)
  }
}

const handleSearch = () => {
  currentPage.value = 1
  loadUsers()
}

// 处理角色筛选的辅助函数
const getRoleFilterValue = () => {
  return roleFilter.value || ''
}

const resetFilters = () => {
  searchQuery.value = ''
  roleFilter.value = ''
  statusFilter.value = ''
  currentPage.value = 1
  loadUsers()
}

const handleSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
  loadUsers()
}

const handleCurrentChange = (page) => {
  currentPage.value = page
  loadUsers()
}

const viewUser = async (user) => {
  try {
    const response = await userService.getUserById(user.id)
    currentUser.value = response.data
    showUserDetailDialog.value = true
  } catch (error) {
    ElMessage.error('获取用户详情失败')
  }
}

const editUser = (user) => {
  // 处理角色数据，优先使用roleNames数组中的第一个角色
  let userRole = ''
  if (Array.isArray(user.roleNames) && user.roleNames.length > 0) {
    // 将中文角色名称转换为角色键值
    const roleName = user.roleNames[0]
    if (roleName === '管理员') userRole = 'admin'
    else if (roleName === '教师') userRole = 'teacher'
    else if (roleName === '学生') userRole = 'student'
  } else if (Array.isArray(user.roles) && user.roles.length > 0) {
    // 如果是角色对象数组，使用第一个角色的roleKey
    if (user.roles[0] && typeof user.roles[0] === 'object' && user.roles[0].roleKey) {
      userRole = user.roles[0].roleKey
    }
  } else if (user.role) {
    userRole = user.role
  }
  
  Object.assign(editForm, {
    id: user.id,
    username: user.username,
    realName: user.realName,
    email: user.email,
    phone: user.phone || '',
    department: user.department || '',
    studentId: user.studentId || '',
    role: userRole,
    status: user.status
  })
  showEditUserDialog.value = true
}

const updateUser = async () => {
  try {
    await editFormRef.value.validate()
    submitting.value = true
    
    await userService.updateUser(editForm.id, {
      realName: editForm.realName,
      email: editForm.email,
      phone: editForm.phone,
      department: editForm.department,
      studentId: editForm.studentId,
      RoleKeys: [editForm.role], // 修改为RoleKeys数组格式
      status: editForm.status
    })
    
    ElMessage.success('用户信息更新成功')
    showEditUserDialog.value = false
    loadUsers()
    loadStats()
  } catch (error) {
    ElMessage.error('更新失败：' + (error.message || '未知错误'))
  } finally {
    submitting.value = false
  }
}

const addUser = async () => {
  try {
    await userFormRef.value.validate()
    submitting.value = true
    
    await userService.createUser({
      username: userForm.username,
      realName: userForm.realName,
      email: userForm.email,
      phone: userForm.phone,
      department: userForm.department,
      studentId: userForm.studentId,
      RoleKeys: [userForm.role], // 修改为RoleKeys数组格式
      status: userForm.status,
      password: userForm.password
    })
    
    ElMessage.success('用户添加成功')
    showAddUserDialog.value = false
    
    // 重置表单
    Object.assign(userForm, {
      username: '',
      realName: '',
      email: '',
      phone: '',
      department: '',
      studentId: '',
      role: '',
      status: 'active',
      password: '',
      confirmPassword: ''
    })
    
    loadUsers()
    loadStats()
  } catch (error) {
    ElMessage.error('添加失败：' + (error.message || '未知错误'))
  } finally {
    submitting.value = false
  }
}

const toggleUserStatus = async (user) => {
  try {
    const action = user.status === 'active' ? '禁用' : '启用'
    await ElMessageBox.confirm(
      `确定要${action}用户 "${user.realName}" 吗？`,
      '确认操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    const newStatus = user.status === 'active' ? 'inactive' : 'active'
    await userService.toggleUserStatus(user.id, newStatus)
    
    ElMessage.success(`用户${action}成功`)
    loadUsers()
    loadStats()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败：' + (error.message || '未知错误'))
    }
  }
}

const resetPassword = async (user) => {
  try {
    await ElMessageBox.confirm(
      `确定要重置用户 "${user.realName}" 的密码吗？`,
      '确认操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await userService.resetUserPassword(user.id)
    
    ElMessage.success('密码重置成功，新密码已发送到用户邮箱')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('密码重置失败：' + (error.message || '未知错误'))
    }
  }
}

const deleteUser = async (user) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除用户 "${user.realName}" 吗？此操作不可恢复！`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'error'
      }
    )
    
    await userService.deleteUser(user.id)
    
    ElMessage.success('用户删除成功')
    loadUsers()
    loadStats()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败：' + (error.message || '未知错误'))
    }
  }
}

const batchDeleteUsers = async () => {
  if (selectedUsers.value.length === 0) {
    ElMessage.warning('请选择要删除的用户')
    return
  }

  try {
    const userNames = selectedUsers.value.map(user => user.realName).join('、')
    await ElMessageBox.confirm(
      `确定要删除以下用户吗？此操作不可恢复！\n${userNames}`,
      '批量删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'error'
      }
    )
    
    const userIds = selectedUsers.value.map(user => user.id)
    await userService.batchDeleteUsers(userIds)
    
    ElMessage.success('批量删除成功')
    selectedUsers.value = []
    loadUsers()
    loadStats()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败：' + (error.message || '未知错误'))
    }
  }
}

const exportUsers = async () => {
  try {
    exporting.value = true
    const params = {
      search: searchQuery.value,
      role: roleFilter.value,
      status: statusFilter.value
    }
    
    const response = await userService.exportUsers(params)
    
    // 创建下载链接
    const blob = new Blob([response], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `用户列表_${new Date().toISOString().split('T')[0]}.xlsx`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    
    ElMessage.success('用户数据导出成功')
  } catch (error) {
    ElMessage.error('导出失败：' + (error.message || '未知错误'))
  } finally {
    exporting.value = false
  }
}

// 生命周期
onMounted(() => {
  loadUsers()
  loadStats()
})
</script>

<style scoped>
.user-management {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
  color: #2c3e50;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.stats-row {
  margin-bottom: 20px;
}

.stats-card {
  text-align: center;
}

.stats-content {
  padding: 10px;
}

.stats-number {
  font-size: 24px;
  font-weight: bold;
  color: #409eff;
  margin-bottom: 5px;
}

.stats-label {
  font-size: 14px;
  color: #666;
}

.search-card {
  margin-bottom: 20px;
}

.user-list-card {
  margin-bottom: 20px;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.table-title {
  font-size: 16px;
  font-weight: bold;
  color: #2c3e50;
}

.table-actions {
  display: flex;
  gap: 10px;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style> 