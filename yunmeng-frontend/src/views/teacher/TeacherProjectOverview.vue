<template>
  <div class="teacher-project-overview">
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
        <span>指导的项目</span>
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
              <el-link type="primary" @click="viewProject(scope.row)">{{ scope.row.title }}</el-link>
            </template>
          </el-table-column>
          <el-table-column prop="studentName" label="学生姓名" width="100" />
          <el-table-column prop="type" label="项目类型" width="120" />
          <el-table-column prop="status" label="项目状态" width="100">
            <template #default="scope">
              <el-tag :type="getStatusType(scope.row.status)">
                {{ getStatusText(scope.row.status) }}
              </el-tag>
            </template>
            </el-table-column>
          <el-table-column prop="createdAt" label="创建时间" width="140">
              <template #default="scope">
                    <span>{{ formatDateTime(scope.row.createdAt) }}</span>
              </template>
          </el-table-column>

          <el-table-column prop="updatedAt" label="更新时间" width="140">
              <template #default="scope">
                  <span>{{ formatDateTime(scope.row.updatedAt) }}</span>
              </template>
          </el-table-column>
          <el-table-column label="操作" width="280" fixed="right">
            <template #default="scope">
              <el-button size="small" type="success" @click="updateProgress(scope.row)">项目审批</el-button>
              <el-button size="small" type="warning" @click="provideGuidance(scope.row)">指导建议</el-button>
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
                  <span class="project-name">{{ project.title }}</span>
                  <el-tag :type="getStatusType(project.status)" size="small">
                    {{ getStatusText(project.status) }}
                  </el-tag>
                </div>
              </template>
              
              <div class="card-content">
                <p><strong>学生:</strong> {{ project.studentName }}</p>
                <p><strong>类型:</strong> {{ project.type }}</p>
                <p><strong>描述:</strong> {{ project.description }}</p>
                <p><strong>创建时间:</strong> {{formatRelativeTime(project.createdAt) }}</p>
                
                <div class="card-actions">
                  <el-button size="small" @click="viewProject(project)">查看</el-button>
                  <el-button size="small" type="warning" @click="provideGuidance(project)">指导</el-button>
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
      width="80%"
    >
      <div v-if="currentProject">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="项目名称">{{ currentProject.name }}</el-descriptions-item>
          <el-descriptions-item label="学生姓名">{{ currentProject.studentName }}</el-descriptions-item>
          <el-descriptions-item label="项目类型">{{ currentProject.type }}</el-descriptions-item>
          <el-descriptions-item label="项目状态">{{ currentProject.status }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ currentProject.createTime }}</el-descriptions-item>
          <el-descriptions-item label="截止时间">{{ currentProject.deadline }}</el-descriptions-item>
          <el-descriptions-item label="项目进度">{{ currentProject.progress }}%</el-descriptions-item>
          <el-descriptions-item label="指导教师">{{ currentProject.teacherName || '当前用户' }}</el-descriptions-item>
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

        <el-divider />
        
        <h4>指导记录</h4>
        <div v-if="currentProject.guidanceRecords && currentProject.guidanceRecords.length > 0">
          <el-timeline>
            <el-timeline-item
              v-for="record in currentProject.guidanceRecords"
              :key="record.id"
              :timestamp="record.date"
              type="success"
            >
              <h5>{{ record.title }}</h5>
              <p>{{ record.content }}</p>
              <p><strong>指导类型:</strong> {{ record.type }}</p>
              <p><strong>指导人:</strong> {{ record.teacherName }}</p>
            </el-timeline-item>
          </el-timeline>
        </div>
        <el-empty v-else description="暂无指导记录" />
      </div>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="projectDetailVisible = false">关闭</el-button>
          <el-button type="primary" @click="editProject(currentProject)">编辑项目</el-button>
          <el-button type="warning" @click="provideGuidance(currentProject)">添加指导</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 更新进度对话框 -->
    <el-dialog
      v-model="progressDialogVisible"
      title="更新项目进度"
      width="50%"
    >
      <el-form :model="progressForm" label-width="100px">
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
        <el-form-item label="评估意见">
          <el-input
            v-model="progressForm.evaluation"
            type="textarea"
            :rows="3"
            placeholder="请提供进度评估意见"
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

    <!-- 指导建议对话框 -->
    <el-dialog
      v-model="guidanceDialogVisible"
      title="提供指导建议"
      width="60%"
    >
      <el-form :model="guidanceForm" :rules="guidanceRules" ref="guidanceFormRef" label-width="100px">
        <el-form-item label="指导类型" prop="type">
          <el-select v-model="guidanceForm.type" placeholder="选择指导类型">
            <el-option label="技术指导" value="technical" />
            <el-option label="进度指导" value="progress" />
            <el-option label="方法指导" value="methodology" />
            <el-option label="问题解决" value="problem_solving" />
            <el-option label="其他" value="other" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="指导标题" prop="title">
          <el-input v-model="guidanceForm.title" placeholder="请输入指导标题" />
        </el-form-item>
        
        <el-form-item label="指导内容" prop="content">
          <el-input
            v-model="guidanceForm.content"
            type="textarea"
            :rows="5"
            placeholder="请详细描述指导内容"
          />
        </el-form-item>
        
        <el-form-item label="建议措施">
          <el-input
            v-model="guidanceForm.suggestions"
            type="textarea"
            :rows="3"
            placeholder="请提供具体的建议措施"
          />
        </el-form-item>
        
        <el-form-item label="预期效果">
          <el-input
            v-model="guidanceForm.expectedOutcome"
            type="textarea"
            :rows="2"
            placeholder="请描述预期达到的效果"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="guidanceDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitGuidance">提交指导</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Document, Clock, Check, Warning, Search, Grid, List } from '@element-plus/icons-vue'
import teacherService from '@/services/teacherService'

// 响应式数据
const loading = ref(false)
const projectList = ref([])
const projectDetailVisible = ref(false)
const progressDialogVisible = ref(false)
const guidanceDialogVisible = ref(false)
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
  description: '',
  evaluation: ''
})

const guidanceForm = ref({
  type: '',
  title: '',
  content: '',
  suggestions: '',
  expectedOutcome: ''
})

const guidanceRules = {
  type: [{ required: true, message: '请选择指导类型', trigger: 'change' }],
  title: [{ required: true, message: '请输入指导标题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入指导内容', trigger: 'blur' }]
}

const guidanceFormRef = ref()

const projectStats = ref([
  {
    title: '总项目数',
    value: 0,
    description: '指导的所有项目',
    icon: Document,
    type: 'total'
  },
  {
    title: '已驳回',
    value: 0,
    description: '已经驳回的项目',
    icon: Clock,
    type: 'rejected'  // 对应 isApproved = -1 或 status = 'rejected'
  },
  {
    title: '已通过',
    value: 0,
    description: '已经通过的项目',
    icon: Check,
    type: 'approved'  // 对应 isApproved = 1
  },
  {
    title: '待审核',
    value: 0,
    description: '等待审核的项目',
    icon: Warning,
    type: 'pending'   // 对应 isApproved = 0
  }
])



const formatDateTime = (dateString) => {
  if (!dateString) return '--'
  
  try {
    const date = new Date(dateString)
    
    // 格式：YYYY-MM-DD HH:mm:ss
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    const hours = String(date.getHours()).padStart(2, '0')
    const minutes = String(date.getMinutes()).padStart(2, '0')
    const seconds = String(date.getSeconds()).padStart(2, '0')
    
    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
  } catch (error) {
    console.error('日期格式化错误:', error)
    return dateString
  }
}

// 格式化相对时间
const formatRelativeTime = (dateString) => {
  if (!dateString) return '--'
  const date = new Date(dateString)
  const now = new Date()
  const diffMs = now - date
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMs / 3600000)
  const diffDays = Math.floor(diffMs / 86400000)
  
  if (diffMins < 60) return `${diffMins}分钟前`
  if (diffHours < 24) return `${diffHours}小时前`
  if (diffDays < 30) return `${diffDays}天前`
  return formatDate(dateString)
}



// 计算属性
const filteredProjects = computed(() => {
  let filtered = projectList.value

  // 搜索过滤
  if (searchQuery.value) {
    filtered = filtered.filter(project => 
      project.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      project.studentName.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
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
      case 'studentName':
        return a.studentName.localeCompare(b.studentName)
      case 'createTime':
      default:
        return new Date(b.createTime) - new Date(a.createTime)
    }
  })

  return filtered
})

// 获取状态类型
const getStatusType = (status) => {
  const typeMap = {
    'submitted': 'info',
    'reviewing': 'warning', 
    'approved': 'success',
    'rejected': 'danger',
    'in_progress': '',
    'completed': 'success'
  }
  return typeMap[status] || 'info'
}

const getStatusText = (status) => {
  const textMap = {
    'submitted': '已提交',
    'reviewing': '审核中',
    'approved': '已通过',
    'rejected': '已拒绝',
    'in_progress': '进行中',
    'completed': '已完成'
  }
  return textMap[status] || status
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
    const response = await teacherService.getGuidedProjects()
    if (response && response.code === 200) {
      // 数据在 response.data.list 中
      projectList.value = response.data.list || []  // 改为 response.data.list
      updateStats()
    } else {
      // 使用模拟数据
      loadMockData()
    }
  } catch (error) {
    console.error('加载项目列表失败:', error)
    // API调用失败时，使用模拟数据
    loadMockData()
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
      studentName: '张三',
      type: '软件开发',
      status: '进行中',
      progress: 75,
      createTime: '2024-01-15',
      deadline: '2024-07-15',
      description: '基于物联网技术的智能校园管理系统',
      plan: '预计6个月完成，分为需求分析、设计、开发、测试四个阶段',
      teacherName: '李教授',
      milestones: [
        { id: 1, title: '需求分析', description: '完成用户需求调研和分析', date: '2024-01-20', completed: true },
        { id: 2, title: '系统设计', description: '完成系统架构和数据库设计', date: '2024-02-15', completed: true },
        { id: 3, title: '功能开发', description: '核心功能模块开发', date: '2024-05-15', completed: false },
        { id: 4, title: '系统测试', description: '功能测试和性能测试', date: '2024-06-15', completed: false }
      ],
      guidanceRecords: [
        { id: 1, title: '架构设计指导', content: '建议采用微服务架构，提高系统可扩展性', type: '技术指导', date: '2024-02-10', teacherName: '李教授' },
        { id: 2, title: '进度管理建议', content: '建议采用敏捷开发方法，每周进行进度回顾', type: '方法指导', date: '2024-03-15', teacherName: '李教授' }
      ]
    },
    {
      id: 2,
      name: '数据分析平台',
      studentName: '李四',
      type: '科研项目',
      status: '待审核',
      progress: 90,
      createTime: '2024-01-14',
      deadline: '2024-06-14',
      description: '大数据分析平台，支持多种数据源和算法',
      plan: '预计8个月完成，包括数据采集、预处理、分析、可视化等模块',
      teacherName: '王教授',
      milestones: [
        { id: 1, title: '数据采集', description: '建立数据采集管道', date: '2024-01-25', completed: true },
        { id: 2, title: '算法实现', description: '核心分析算法开发', date: '2024-03-15', completed: true },
        { id: 3, title: '平台集成', description: '各模块集成测试', date: '2024-05-15', completed: false }
      ],
      guidanceRecords: [
        { id: 1, title: '算法优化建议', content: '建议使用并行计算提高算法效率', type: '技术指导', date: '2024-02-20', teacherName: '王教授' }
      ]
    },
    {
      id: 3,
      name: '在线教育平台',
      studentName: '王五',
      type: '创新项目',
      status: '已完成',
      progress: 100,
      createTime: '2024-01-10',
      deadline: '2024-05-10',
      description: '基于Web的在线教育学习平台',
      plan: '预计4个月完成，包括用户管理、课程管理、学习跟踪等模块',
      teacherName: '赵教授',
      milestones: [
        { id: 1, title: '用户系统', description: '用户注册登录功能', date: '2024-02-10', completed: true },
        { id: 2, title: '课程管理', description: '课程创建和发布功能', date: '2024-03-10', completed: true },
        { id: 3, title: '学习跟踪', description: '学习进度和成绩统计', date: '2024-04-10', completed: true }
      ],
      guidanceRecords: [
        { id: 1, title: '用户体验优化', content: '建议增加学习路径推荐功能', type: '方法指导', date: '2024-03-05', teacherName: '赵教授' },
        { id: 2, title: '性能优化建议', content: '建议使用CDN加速静态资源加载', type: '技术指导', date: '2024-04-01', teacherName: '赵教授' }
      ]
    }
  ]
  updateStats()
}

const updateStats = () => {
  const projects = projectList.value
  
  // 根据实际字段计算
  const total = projects.length
  
  // 已通过：isApproved === 1
  const approved = projects.filter(p => p.isApproved === 1).length
  
  // 已驳回：isApproved === -1 (如果有这个状态)
  const rejected = projects.filter(p => p.isApproved === -1).length
  
  // 待审核：isApproved === 0 或 status === 'submitted'/'reviewing'
  const pending = projects.filter(p => 
    p.isApproved === 0 || 
    ['submitted', 'reviewing'].includes(p.status)
  ).length
  
  // 更新统计卡片
  projectStats.value[0].value = total     // 总项目数
  projectStats.value[1].value = rejected  // 已驳回
  projectStats.value[2].value = approved  // 已通过
  projectStats.value[3].value = pending   // 待审核
}

// 查看项目详情
const viewProject = (project) => {
  currentProject.value = project
  projectDetailVisible.value = true
}

// 编辑项目
const editProject = (project) => {
  ElMessage.info(`编辑项目：${project.name}`)
}

// 更新进度
const updateProgress = (project) => {
  progressForm.value.progress = project.progress
  progressForm.value.description = ''
  progressForm.value.evaluation = ''
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

// 提供指导建议
const provideGuidance = (project) => {
  currentProject.value = project
  guidanceForm.value = {
    type: '',
    title: '',
    content: '',
    suggestions: '',
    expectedOutcome: ''
  }
  guidanceDialogVisible.value = true
}

// 提交指导建议
const submitGuidance = async () => {
  try {
    await guidanceFormRef.value.validate()
    
    // 这里应该调用实际的API
    // await teacherService.submitGuidance(currentProject.value.id, guidanceForm.value)
    
    // 添加指导记录到项目
    const newGuidance = {
      id: Date.now(),
      title: guidanceForm.value.title,
      content: guidanceForm.value.content,
      type: getGuidanceTypeLabel(guidanceForm.value.type),
      date: new Date().toLocaleString('zh-CN'),
      teacherName: '当前用户'
    }
    
    if (!currentProject.value.guidanceRecords) {
      currentProject.value.guidanceRecords = []
    }
    currentProject.value.guidanceRecords.push(newGuidance)
    
    ElMessage.success('指导建议提交成功')
    guidanceDialogVisible.value = false
  } catch (error) {
    console.error('提交指导建议失败:', error)
    ElMessage.error('提交失败')
  }
}

// 查看指导记录
const viewGuidanceRecords = (project) => {
  currentProject.value = project
  projectDetailVisible.value = true
}

// 获取指导类型标签
const getGuidanceTypeLabel = (type) => {
  const labelMap = {
    'technical': '技术指导',
    'progress': '进度指导',
    'methodology': '方法指导',
    'problem_solving': '问题解决',
    'other': '其他'
  }
  return labelMap[type] || '其他'
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
.teacher-project-overview {
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