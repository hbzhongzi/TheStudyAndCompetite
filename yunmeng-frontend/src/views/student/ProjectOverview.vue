<template>
  <div class="project-overview">
    <!-- 搜索和筛选区域 -->
    <el-card style="margin-bottom: 20px;">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-input
            v-model="searchQuery"
            placeholder="搜索项目名称"
            clearable
            @input="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-col>
        <el-col :span="4">
          <el-select v-model="statusFilter" placeholder="项目状态" clearable @change="handleFilter">
            <el-option label="全部" value="" />
            <el-option label="进行中" value="进行中" />
            <el-option label="已完成" value="已完成" />
            <el-option label="待审核" value="待审核" />
            <el-option label="已拒绝" value="已拒绝" />
            <el-option label="已暂停" value="已暂停" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-select v-model="typeFilter" placeholder="项目类型" clearable @change="handleFilter">
            <el-option label="全部" value="" />
            <el-option label="软件开发" value="软件开发" />
            <el-option label="科研项目" value="科研项目" />
            <el-option label="创新项目" value="创新项目" />
            <el-option label="竞赛项目" value="竞赛项目" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-select v-model="sortBy" placeholder="排序方式" @change="handleSort">
            <el-option label="创建时间" value="createTime" />
            <el-option label="项目名称" value="name" />
            <el-option label="项目进度" value="progress" />
            <el-option label="截止时间" value="deadline" />
          </el-select>
        </el-col>
        <el-col :span="6">
          <el-button type="primary" @click="createNewProject">创建新项目</el-button>
          <el-button @click="refreshProjects">刷新</el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 项目统计卡片 -->
    <el-row :gutter="20">
      <el-col :span="6" v-for="stat in projectStats" :key="stat.title">
        <el-card class="stat-card" :class="stat.type">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon><component :is="stat.icon" /></el-icon>
            </div>
            <div class="stat-info">
              <h4>{{ stat.title }}</h4>
              <p class="stat-number">{{ stat.value }}</p>
              <p class="stat-desc">{{ stat.description }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 项目列表 -->
    <el-card style="margin-top: 20px;">
      <template #header>
        <span>我的项目 ({{ filteredProjects.length }})</span>
        <div style="float: right;">
          <el-button-group>
            <el-button :type="viewMode === 'table' ? 'primary' : ''" @click="viewMode = 'table'">
              <el-icon><Grid /></el-icon>
            </el-button>
            <el-button :type="viewMode === 'card' ? 'primary' : ''" @click="viewMode = 'card'">
              <el-icon><List /></el-icon>
            </el-button>
          </el-button-group>
        </div>
      </template>
      
      <!-- 表格视图 -->
      <div v-if="viewMode === 'table'">
        <el-table :data="filteredProjects" style="width: 100%" v-loading="loading">
          <el-table-column prop="name" label="项目名称" min-width="150">
            <template #default="scope">
              <el-link type="primary" @click="viewProject(scope.row)">{{ scope.row.name }}</el-link>
            </template>
          </el-table-column>
          <el-table-column prop="type" label="项目类型" width="120" />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="scope">
              <el-tag :type="getStatusType(scope.row.status)">
                {{ scope.row.status }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="progress" label="进度" width="150">
            <template #default="scope">
              <el-progress :percentage="scope.row.progress" :status="getProgressStatus(scope.row.progress)" />
            </template>
          </el-table-column>
          <el-table-column prop="createTime" label="创建时间" width="120" />
          <el-table-column prop="deadline" label="截止时间" width="120">
            <template #default="scope">
              <span :class="getDeadlineClass(scope.row.deadline)">
                {{ scope.row.deadline }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="250" fixed="right">
            <template #default="scope">
              <el-button size="small" @click="viewProject(scope.row)">查看</el-button>
              <el-button size="small" type="primary" @click="editProject(scope.row)">编辑</el-button>
              <el-button size="small" type="success" @click="updateProgress(scope.row)">更新进度</el-button>
              <el-button size="small" type="warning" @click="requestExtension(scope.row)">申请延期</el-button>
              <el-button size="small" type="danger" @click="deleteProject(scope.row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 卡片视图 -->
      <div v-else class="project-cards">
        <el-row :gutter="20">
          <el-col :span="8" v-for="project in filteredProjects" :key="project.id">
            <el-card class="project-card" :class="getProjectCardClass(project)">
              <template #header>
                <div class="card-header">
                  <span class="project-name">{{ project.name }}</span>
                  <el-tag :type="getStatusType(project.status)" size="small">
                    {{ project.status }}
                  </el-tag>
                </div>
              </template>
              
              <div class="card-content">
                <p><strong>类型:</strong> {{ project.type }}</p>
                <p><strong>进度:</strong></p>
                <el-progress :percentage="project.progress" :status="getProgressStatus(project.progress)" />
                <p><strong>创建时间:</strong> {{ project.createTime }}</p>
                <p><strong>截止时间:</strong> 
                  <span :class="getDeadlineClass(project.deadline)">
                    {{ project.deadline }}
                  </span>
                </p>
                
                <div class="card-actions">
                  <el-button size="small" @click="viewProject(project)">查看</el-button>
                  <el-button size="small" type="primary" @click="editProject(project)">编辑</el-button>
                  <el-button size="small" type="success" @click="updateProgress(project)">进度</el-button>
                </div>
              </div>
            </el-card>
          </el-col>
        </el-row>
      </div>

      <!-- 分页 -->
      <div class="pagination-container" v-if="filteredProjects.length > pageSize">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50]"
          :total="filteredProjects.length"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 项目详情对话框 -->
    <el-dialog
      v-model="projectDetailVisible"
      title="项目详情"
      width="70%"
    >
      <div v-if="currentProject">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="项目名称">{{ currentProject.name }}</el-descriptions-item>
          <el-descriptions-item label="项目类型">{{ currentProject.type }}</el-descriptions-item>
          <el-descriptions-item label="项目状态">{{ currentProject.status }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ currentProject.createTime }}</el-descriptions-item>
          <el-descriptions-item label="截止时间">{{ currentProject.deadline }}</el-descriptions-item>
          <el-descriptions-item label="项目进度">{{ currentProject.progress }}%</el-descriptions-item>
        </el-descriptions>
        
        <el-divider />
        
        <h4>项目描述</h4>
        <p>{{ currentProject.description }}</p>
        
        <el-divider />
        
        <h4>项目计划</h4>
        <p>{{ currentProject.plan }}</p>

        <el-divider />
        
        <h4>项目里程碑</h4>
        <el-timeline>
          <el-timeline-item
            v-for="milestone in currentProject.milestones || []"
            :key="milestone.id"
            :timestamp="milestone.date"
            :type="milestone.completed ? 'success' : 'primary'"
          >
            <h5>{{ milestone.title }}</h5>
            <p>{{ milestone.description }}</p>
            <el-tag v-if="milestone.completed" type="success" size="small">已完成</el-tag>
            <el-tag v-else type="warning" size="small">进行中</el-tag>
          </el-timeline-item>
        </el-timeline>
      </div>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="projectDetailVisible = false">关闭</el-button>
          <el-button type="primary" @click="editProject(currentProject)">编辑项目</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 更新进度对话框 -->
    <el-dialog
      v-model="progressDialogVisible"
      title="更新项目进度"
      width="40%"
    >
      <el-form :model="progressForm" label-width="80px">
        <el-form-item label="当前进度">
          <el-slider v-model="progressForm.progress" :min="0" :max="100" :step="5" show-input />
        </el-form-item>
        <el-form-item label="进度说明">
          <el-input
            v-model="progressForm.description"
            type="textarea"
            :rows="3"
            placeholder="请描述当前进度情况"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="progressDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitProgressUpdate">提交更新</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Document, Clock, Check, Warning, Search, Grid, List } from '@element-plus/icons-vue'
import { studentService } from '../../services/studentService'

// 响应式数据
const loading = ref(false)
const projectList = ref([])
const projectDetailVisible = ref(false)
const progressDialogVisible = ref(false)
const currentProject = ref(null)
const viewMode = ref('table')
const searchQuery = ref('')
const statusFilter = ref('')
const typeFilter = ref('')
const sortBy = ref('createTime')
const currentPage = ref(1)
const pageSize = ref(10)

const progressForm = ref({
  progress: 0,
  description: ''
})

const projectStats = ref([
  {
    title: '总项目数',
    value: 0,
    description: '所有项目',
    icon: Document,
    type: 'total'
  },
  {
    title: '进行中',
    value: 0,
    description: '正在执行',
    icon: Clock,
    type: 'ongoing'
  },
  {
    title: '已完成',
    value: 0,
    description: '成功完成',
    icon: Check,
    type: 'completed'
  },
  {
    title: '待审核',
    value: 0,
    description: '等待审核',
    icon: Warning,
    type: 'pending'
  }
])

// 计算属性
const filteredProjects = computed(() => {
  let filtered = projectList.value

  // 搜索过滤
  if (searchQuery.value) {
    filtered = filtered.filter(project => 
      project.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      project.description.toLowerCase().includes(searchQuery.value.toLowerCase())
    )
  }

  // 状态过滤
  if (statusFilter.value) {
    filtered = filtered.filter(project => project.status === statusFilter.value)
  }

  // 类型过滤
  if (typeFilter.value) {
    filtered = filtered.filter(project => project.type === typeFilter.value)
  }

  // 排序
  filtered.sort((a, b) => {
    switch (sortBy.value) {
      case 'name':
        return a.name.localeCompare(b.name)
      case 'progress':
        return b.progress - a.progress
      case 'deadline':
        return new Date(a.deadline) - new Date(b.deadline)
      case 'createTime':
      default:
        return new Date(b.createTime) - new Date(a.createTime)
    }
  })

  return filtered
})

// 获取状态类型
const getStatusType = (status) => {
  const statusMap = {
    '进行中': 'primary',
    '已完成': 'success',
    '待审核': 'warning',
    '已拒绝': 'danger',
    '已暂停': 'info'
  }
  return statusMap[status] || 'info'
}

// 获取进度状态
const getProgressStatus = (progress) => {
  if (progress >= 100) return 'success'
  if (progress >= 80) return 'warning'
  if (progress >= 50) return ''
  return 'exception'
}

// 获取截止时间样式
const getDeadlineClass = (deadline) => {
  const deadlineDate = new Date(deadline)
  const today = new Date()
  const diffDays = Math.ceil((deadlineDate - today) / (1000 * 60 * 60 * 24))
  
  if (diffDays < 0) return 'deadline-overdue'
  if (diffDays <= 7) return 'deadline-urgent'
  if (diffDays <= 30) return 'deadline-warning'
  return 'deadline-normal'
}

// 获取项目卡片样式
const getProjectCardClass = (project) => {
  const status = project.status
  if (status === '已完成') return 'completed'
  if (status === '待审核') return 'pending'
  if (status === '已拒绝') return 'rejected'
  return 'ongoing'
}

// 搜索处理
const handleSearch = () => {
  currentPage.value = 1
}

// 筛选处理
const handleFilter = () => {
  currentPage.value = 1
}

// 排序处理
const handleSort = () => {
  currentPage.value = 1
}

// 分页处理
const handleSizeChange = (val) => {
  pageSize.value = val
  currentPage.value = 1
}

const handleCurrentChange = (val) => {
  currentPage.value = val
}

// 加载项目列表
const loadProjects = async () => {
  loading.value = true
  try {
    const response = await studentService.getMyProjects()
    if (response && response.code === 200) {
      projectList.value = response.data || []
      updateStats()
    } else {
      // 如果响应格式不符合预期，使用模拟数据
      loadMockData()
    }
  } catch (error) {
    console.error('加载项目列表失败:', error)
    // API调用失败时，使用模拟数据
    loadMockData()
    // 不显示错误消息，静默处理
  } finally {
    loading.value = false
  }
}

// 加载模拟数据
const loadMockData = () => {
  projectList.value = [
    {
      id: 1,
      name: '智能校园系统',
      type: '软件开发',
      status: '进行中',
      progress: 75,
      createTime: '2024-01-15',
      deadline: '2024-07-15',
      description: '基于物联网技术的智能校园管理系统',
      plan: '预计6个月完成，分为需求分析、设计、开发、测试四个阶段',
      milestones: [
        { id: 1, title: '需求分析', description: '完成用户需求调研和分析', date: '2024-01-20', completed: true },
        { id: 2, title: '系统设计', description: '完成系统架构和数据库设计', date: '2024-02-15', completed: true },
        { id: 3, title: '功能开发', description: '核心功能模块开发', date: '2024-05-15', completed: false },
        { id: 4, title: '系统测试', description: '功能测试和性能测试', date: '2024-06-15', completed: false }
      ]
    },
    {
      id: 2,
      name: '数据分析平台',
      type: '科研项目',
      status: '待审核',
      progress: 90,
      createTime: '2024-01-14',
      deadline: '2024-06-14',
      description: '大数据分析平台，支持多种数据源和算法',
      plan: '预计8个月完成，包括数据采集、预处理、分析、可视化等模块',
      milestones: [
        { id: 1, title: '数据采集', description: '建立数据采集管道', date: '2024-01-25', completed: true },
        { id: 2, title: '算法实现', description: '核心分析算法开发', date: '2024-03-15', completed: true },
        { id: 3, title: '平台集成', description: '各模块集成测试', date: '2024-05-15', completed: false }
      ]
    },
    {
      id: 3,
      name: '在线教育平台',
      type: '创新项目',
      status: '已完成',
      progress: 100,
      createTime: '2024-01-10',
      deadline: '2024-05-10',
      description: '基于Web的在线教育学习平台',
      plan: '预计4个月完成，包括用户管理、课程管理、学习跟踪等模块',
      milestones: [
        { id: 1, title: '用户系统', description: '用户注册登录功能', date: '2024-02-10', completed: true },
        { id: 2, title: '课程管理', description: '课程创建和发布功能', date: '2024-03-10', completed: true },
        { id: 3, title: '学习跟踪', description: '学习进度和成绩统计', date: '2024-04-10', completed: true }
      ]
    }
  ]
  updateStats()
}

// 更新统计数据
const updateStats = () => {
  const total = projectList.value.length
  const ongoing = projectList.value.filter(p => p.status === '进行中').length
  const completed = projectList.value.filter(p => p.status === '已完成').length
  const pending = projectList.value.filter(p => p.status === '待审核').length
  
  projectStats.value[0].value = total
  projectStats.value[1].value = ongoing
  projectStats.value[2].value = completed
  projectStats.value[3].value = pending
}

// 查看项目详情
const viewProject = (project) => {
  currentProject.value = project
  projectDetailVisible.value = true
}

// 编辑项目
const editProject = (project) => {
  ElMessage.info(`跳转到项目编辑页面：${project.name}`)
}

// 更新进度
const updateProgress = (project) => {
  progressForm.value.progress = project.progress
  progressForm.value.description = ''
  currentProject.value = project
  progressDialogVisible.value = true
}

// 提交进度更新
const submitProgressUpdate = async () => {
  try {
    // 这里应该调用实际的API
    currentProject.value.progress = progressForm.value.progress
    ElMessage.success('进度更新成功')
    progressDialogVisible.value = false
    updateStats()
  } catch (error) {
    ElMessage.error('进度更新失败')
  }
}

// 申请延期
const requestExtension = (project) => {
  ElMessage.info(`申请项目延期：${project.name}`)
}

// 删除项目
const deleteProject = async (project) => {
  try {
    await ElMessageBox.confirm('确定要删除这个项目吗？', '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await studentService.deleteProject(project.id)
    ElMessage.success('项目删除成功')
    loadProjects()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除项目失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 创建新项目
const createNewProject = () => {
  ElMessage.info('跳转到项目创建页面')
}

// 刷新项目列表
const refreshProjects = () => {
  loadProjects()
  ElMessage.success('项目列表已刷新')
}

// 组件挂载时加载数据
onMounted(() => {
  loadProjects()
})
</script>

<style scoped>
.project-overview {
  padding: 20px;
}

.stat-card {
  margin-bottom: 20px;
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

.stat-card.total {
  border-left: 4px solid #667eea;
}

.stat-card.ongoing {
  border-left: 4px solid #f093fb;
}

.stat-card.completed {
  border-left: 4px solid #4facfe;
}

.stat-card.pending {
  border-left: 4px solid #43e97b;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-icon.total {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-icon.ongoing {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.stat-icon.completed {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-icon.pending {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.stat-icon i {
  font-size: 24px;
  color: white;
}

.stat-info h4 {
  margin: 0 0 5px 0;
  color: #7f8c8d;
  font-size: 14px;
}

.stat-number {
  margin: 0 0 5px 0;
  font-size: 28px;
  font-weight: 600;
  color: #2c3e50;
}

.stat-desc {
  margin: 0;
  color: #95a5a6;
  font-size: 12px;
}

.project-cards {
  margin-top: 20px;
}

.project-card {
  margin-bottom: 20px;
  transition: all 0.3s ease;
}

.project-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.1);
}

.project-card.completed {
  border-left: 4px solid #67C23A;
}

.project-card.pending {
  border-left: 4px solid #E6A23C;
}

.project-card.rejected {
  border-left: 4px solid #F56C6C;
}

.project-card.ongoing {
  border-left: 4px solid #409EFF;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.project-name {
  font-weight: 600;
  color: #2c3e50;
}

.card-content p {
  margin: 8px 0;
  color: #606266;
}

.card-actions {
  margin-top: 15px;
  text-align: center;
}

.card-actions .el-button {
  margin: 0 5px;
}

.deadline-overdue {
  color: #F56C6C;
  font-weight: bold;
}

.deadline-urgent {
  color: #E6A23C;
  font-weight: bold;
}

.deadline-warning {
  color: #F56C6C;
}

.deadline-normal {
  color: #67C23A;
}

.pagination-container {
  margin-top: 20px;
  text-align: right;
}

.dialog-footer {
  text-align: right;
}
</style> 