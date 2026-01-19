<template>
  <div class="teacher-review-tasks">
    <el-card>
      <template #header>
        <div class="header-content">
          <span>审核任务管理</span>
          <el-button type="primary" @click="refreshTasks">刷新</el-button>
        </div>
      </template>

      <!-- 审核任务统计 -->
      <el-row :gutter="20" class="task-stats">
        <el-col :span="6" v-for="stat in taskStats" :key="stat.label">
          <el-card class="stat-card" :class="stat.type">
            <div class="stat-content">
              <div class="stat-icon">
                <el-icon><component :is="stat.icon" /></el-icon>
              </div>
              <div class="stat-info">
                <h4>{{ stat.label }}</h4>
                <p class="stat-number">{{ stat.value }}</p>
                <p class="stat-desc">{{ stat.description }}</p>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 搜索和筛选 -->
      <el-card class="search-card">
        <el-form :inline="true" class="search-form">
          <el-form-item label="搜索">
            <el-input
              v-model="searchQuery"
              placeholder="搜索项目名称或学生姓名"
              clearable
              @input="handleSearch"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item label="任务类型">
            <el-select v-model="selectedType" placeholder="所有类型" clearable @change="handleSearch">
              <el-option label="所有类型" value="" />
              <el-option label="项目审核" value="project" />
              <el-option label="里程碑审核" value="milestone" />
              <el-option label="文件审核" value="file" />
              <el-option label="延期审核" value="extension" />
            </el-select>
          </el-form-item>
          <el-form-item label="优先级">
            <el-select v-model="selectedPriority" placeholder="所有优先级" clearable @change="handleSearch">
              <el-option label="所有优先级" value="" />
              <el-option label="低" value="low" />
              <el-option label="中" value="medium" />
              <el-option label="高" value="high" />
              <el-option label="紧急" value="urgent" />
            </el-select>
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="selectedStatus" placeholder="所有状态" clearable @change="handleSearch">
              <el-option label="所有状态" value="" />
              <el-option label="待审核" value="pending" />
              <el-option label="审核中" value="reviewing" />
              <el-option label="已完成" value="completed" />
            </el-select>
          </el-form-item>
          <el-form-item label="排序">
            <el-select v-model="sortBy" @change="handleSearch">
              <el-option label="创建时间" value="createTime" />
              <el-option label="截止时间" value="deadline" />
              <el-option label="优先级" value="priority" />
            </el-select>
          </el-form-item>
        </el-form>
      </el-card>

      <!-- 审核任务列表 -->
      <el-table :data="filteredTasks" style="width: 100%" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column prop="taskName" label="任务名称" min-width="200">
          <template #default="scope">
            <div class="task-name">
              <el-icon class="task-icon" :class="getTaskIconClass(scope.row.type)">
                <component :is="getTaskIcon(scope.row.type)" />
              </el-icon>
              <span>{{ scope.row.taskName }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="projectName" label="项目名称" width="150" />
        <el-table-column prop="studentName" label="申请人" width="100" />
        <el-table-column prop="type" label="任务类型" width="100">
          <template #default="scope">
            <el-tag :type="getTaskTypeTag(scope.row.type)" size="small">
              {{ getTaskTypeLabel(scope.row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="priority" label="优先级" width="100">
          <template #default="scope">
            <el-tag :type="getPriorityType(scope.row.priority)" size="small">
              {{ getPriorityLabel(scope.row.priority) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)" size="small">
              {{ getStatusLabel(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createTime" label="创建时间" width="150" />
        <el-table-column prop="deadline" label="截止时间" width="150">
          <template #default="scope">
            <span :class="getDeadlineClass(scope.row.deadline)">
              {{ scope.row.deadline }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="scope">
            <el-button size="small" @click="viewTaskDetail(scope.row)">查看详情</el-button>
            <el-button 
              v-if="scope.row.status === 'pending'"
              size="small" 
              type="primary" 
              @click="startReview(scope.row)"
            >
              开始审核
            </el-button>
            <el-button 
              v-if="scope.row.status === 'reviewing'"
              size="small" 
              type="success" 
              @click="completeReview(scope.row)"
            >
              完成审核
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 批量操作 -->
      <div v-if="selectedTasks.length > 0" class="batch-actions">
        <el-button type="primary" @click="batchStartReview">批量开始审核</el-button>
        <el-button type="success" @click="batchCompleteReview">批量完成审核</el-button>
        <el-button @click="clearSelection">清除选择</el-button>
      </div>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="filteredTasks.length"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 任务详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="审核任务详情"
      width="70%"
    >
      <div v-if="currentTask" class="task-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="任务名称">{{ currentTask.taskName }}</el-descriptions-item>
          <el-descriptions-item label="项目名称">{{ currentTask.projectName }}</el-descriptions-item>
          <el-descriptions-item label="申请人">{{ currentTask.studentName }}</el-descriptions-item>
          <el-descriptions-item label="任务类型">
            <el-tag :type="getTaskTypeTag(currentTask.type)" size="small">
              {{ getTaskTypeLabel(currentTask.type) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="优先级">
            <el-tag :type="getPriorityType(currentTask.priority)" size="small">
              {{ getPriorityLabel(currentTask.priority) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="当前状态">
            <el-tag :type="getStatusType(currentTask.status)" size="small">
              {{ getStatusLabel(currentTask.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ currentTask.createTime }}</el-descriptions-item>
          <el-descriptions-item label="截止时间">
            <span :class="getDeadlineClass(currentTask.deadline)">
              {{ currentTask.deadline }}
            </span>
          </el-descriptions-item>
        </el-descriptions>
        
        <el-divider />
        
        <h4>任务描述</h4>
        <p>{{ currentTask.description }}</p>
        
        <el-divider />
        
        <h4>审核要求</h4>
        <div v-if="currentTask.requirements && currentTask.requirements.length > 0">
          <ul>
            <li v-for="req in currentTask.requirements" :key="req.id">
              {{ req.description }}
            </li>
          </ul>
        </div>
        <p v-else>暂无具体要求</p>
        
        <el-divider />
        
        <h4>相关材料</h4>
        <div v-if="currentTask.materials && currentTask.materials.length > 0">
          <el-table :data="currentTask.materials" style="width: 100%">
            <el-table-column prop="name" label="材料名称" />
            <el-table-column prop="type" label="类型" />
            <el-table-column prop="size" label="大小" />
            <el-table-column label="操作" width="150">
              <template #default="scope">
                <el-button size="small" @click="downloadMaterial(scope.row)">下载</el-button>
                <el-button size="small" type="primary" @click="previewMaterial(scope.row)">预览</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
        <el-empty v-else description="暂无相关材料" />
        
        <el-divider v-if="currentTask.reviewNotes" />
        
        <h4 v-if="currentTask.reviewNotes">审核记录</h4>
        <div v-if="currentTask.reviewNotes && currentTask.reviewNotes.length > 0" class="review-notes">
          <el-timeline>
            <el-timeline-item
              v-for="note in currentTask.reviewNotes"
              :key="note.id"
              :timestamp="note.time"
              type="success"
            >
              <h6>{{ note.title }}</h6>
              <p>{{ note.content }}</p>
              <p><strong>审核人:</strong> {{ note.reviewer }}</p>
            </el-timeline-item>
          </el-timeline>
        </div>
      </div>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="detailDialogVisible = false">关闭</el-button>
          <el-button 
            v-if="currentTask && currentTask.status === 'pending'"
            type="primary" 
            @click="startReview(currentTask)"
          >
            开始审核
          </el-button>
          <el-button 
            v-if="currentTask && currentTask.status === 'reviewing'"
            type="success" 
            @click="completeReview(currentTask)"
          >
            完成审核
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 审核对话框 -->
    <el-dialog
      v-model="reviewDialogVisible"
      :title="reviewForm.action === 'start' ? '开始审核' : '完成审核'"
      width="60%"
    >
      <el-form :model="reviewForm" :rules="reviewRules" ref="reviewFormRef" label-width="100px">
        <el-form-item label="审核标题" prop="title">
          <el-input v-model="reviewForm.title" placeholder="请输入审核标题" />
        </el-form-item>
        
        <el-form-item label="审核内容" prop="content">
          <el-input
            v-model="reviewForm.content"
            type="textarea"
            :rows="5"
            placeholder="请详细描述审核内容"
          />
        </el-form-item>
        
        <el-form-item label="审核结果" v-if="reviewForm.action === 'complete'">
          <el-radio-group v-model="reviewForm.result">
            <el-radio label="approved">通过</el-radio>
            <el-radio label="rejected">拒绝</el-radio>
            <el-radio label="revision">需要修改</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item label="修改要求" v-if="reviewForm.result === 'revision'">
          <el-input
            v-model="reviewForm.revisionRequirements"
            type="textarea"
            :rows="3"
            placeholder="请说明需要修改的具体内容"
          />
        </el-form-item>
        
        <el-form-item label="下次审核时间" v-if="reviewForm.result === 'revision'">
          <el-date-picker
            v-model="reviewForm.nextReviewTime"
            type="date"
            placeholder="选择下次审核时间"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        
        <el-form-item label="评分" v-if="reviewForm.result === 'approved'">
          <el-rate
            v-model="reviewForm.score"
            :max="10"
            show-score
            :texts="['很差', '差', '一般', '好', '很好']"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="reviewDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitReview">提交</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Document, Clock, Check, Warning, Search } from '@element-plus/icons-vue'
import { teacherService } from '../../services/teacherService'

// 响应式数据
const searchQuery = ref('')
const selectedType = ref('')
const selectedPriority = ref('')
const selectedStatus = ref('')
const sortBy = ref('createTime')
const currentPage = ref(1)
const pageSize = ref(20)
const selectedTasks = ref([])
const detailDialogVisible = ref(false)
const reviewDialogVisible = ref(false)
const currentTask = ref(null)

const reviewForm = ref({
  action: 'start',
  title: '',
  content: '',
  result: 'approved',
  revisionRequirements: '',
  nextReviewTime: '',
  score: 8
})

const reviewRules = {
  title: [{ required: true, message: '请输入审核标题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入审核内容', trigger: 'blur' }]
}

const reviewFormRef = ref()

// 审核任务数据
const tasks = ref([
  {
    id: 1,
    taskName: '项目中期审核',
    projectName: '智能校园系统',
    studentName: '张三',
    type: 'project',
    priority: 'high',
    status: 'pending',
    createTime: '2024-05-20 09:00:00',
    deadline: '2024-05-25 18:00:00',
    description: '对智能校园系统项目进行中期进度审核，评估项目完成情况和技术实现质量',
    requirements: [
      { id: 1, description: '检查项目进度是否符合计划' },
      { id: 2, description: '评估技术实现质量' },
      { id: 3, description: '检查文档完整性' }
    ],
    materials: [
      { name: '项目进度报告.pdf', type: 'PDF', size: '2.1MB' },
      { name: '技术文档.docx', type: 'Word', size: '1.8MB' },
      { name: '演示视频.mp4', type: 'Video', size: '15.2MB' }
    ]
  },
  {
    id: 2,
    taskName: '里程碑审核',
    projectName: '数据分析平台',
    studentName: '李四',
    type: 'milestone',
    priority: 'medium',
    status: 'reviewing',
    createTime: '2024-05-18 14:30:00',
    deadline: '2024-05-22 18:00:00',
    description: '审核数据分析平台项目的需求分析里程碑完成情况',
    requirements: [
      { id: 1, description: '检查需求分析报告质量' },
      { id: 2, description: '验证需求完整性' }
    ],
    materials: [
      { name: '需求分析报告.pdf', type: 'PDF', size: '3.2MB' }
    ],
    reviewNotes: [
      { id: 1, title: '开始审核', content: '开始对需求分析里程碑进行审核', time: '2024-05-19 10:00:00', reviewer: '当前用户' }
    ]
  },
  {
    id: 3,
    taskName: '文件审核',
    projectName: '在线教育平台',
    studentName: '王五',
    type: 'file',
    priority: 'urgent',
    status: 'pending',
    createTime: '2024-05-21 16:00:00',
    deadline: '2024-05-23 18:00:00',
    description: '审核在线教育平台项目的系统设计文档',
    requirements: [
      { id: 1, description: '检查设计文档规范性' },
      { id: 2, description: '验证技术方案可行性' }
    ],
    materials: [
      { name: '系统设计文档.pdf', type: 'PDF', size: '4.5MB' },
      { name: '架构图.vsdx', type: 'Visio', size: '2.8MB' }
    ]
  },
  {
    id: 4,
    taskName: '延期申请审核',
    projectName: '智能校园系统',
    studentName: '张三',
    type: 'extension',
    priority: 'high',
    status: 'completed',
    createTime: '2024-05-15 11:00:00',
    deadline: '2024-05-17 18:00:00',
    description: '审核智能校园系统项目的延期申请',
    requirements: [
      { id: 1, description: '评估延期理由合理性' },
      { id: 2, description: '分析延期影响' }
    ],
    materials: [
      { name: '延期申请说明.docx', type: 'Word', size: '512KB' }
    ],
    reviewNotes: [
      { id: 1, title: '开始审核', content: '开始审核延期申请', time: '2024-05-16 09:00:00', reviewer: '当前用户' },
      { id: 2, title: '完成审核', content: '延期申请审核通过，延期理由充分', time: '2024-05-16 16:00:00', reviewer: '当前用户' }
    ]
  }
])

// 审核任务统计
const taskStats = computed(() => [
  {
    label: '总任务数',
    value: tasks.value.length,
    description: '所有审核任务',
    icon: Document,
    type: 'total'
  },
  {
    label: '待审核',
    value: tasks.value.filter(t => t.status === 'pending').length,
    description: '需要开始审核',
    icon: Clock,
    type: 'pending'
  },
  {
    label: '审核中',
    value: tasks.value.filter(t => t.status === 'reviewing').length,
    description: '正在审核',
    icon: Warning,
    type: 'reviewing'
  },
  {
    label: '已完成',
    value: tasks.value.filter(t => t.status === 'completed').length,
    description: '审核完成',
    icon: Check,
    type: 'completed'
  }
])

// 过滤后的审核任务
const filteredTasks = computed(() => {
  let result = [...tasks.value]
  
  // 搜索过滤
  if (searchQuery.value) {
    result = result.filter(task => 
      task.taskName.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      task.projectName.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      task.studentName.toLowerCase().includes(searchQuery.value.toLowerCase())
    )
  }
  
  // 类型过滤
  if (selectedType.value) {
    result = result.filter(task => task.type === selectedType.value)
  }
  
  // 优先级过滤
  if (selectedPriority.value) {
    result = result.filter(task => task.priority === selectedPriority.value)
  }
  
  // 状态过滤
  if (selectedStatus.value) {
    result = result.filter(task => task.status === selectedStatus.value)
  }
  
  // 排序
  result.sort((a, b) => {
    switch (sortBy.value) {
      case 'createTime':
        return new Date(b.createTime) - new Date(a.createTime)
      case 'deadline':
        return new Date(a.deadline) - new Date(b.deadline)
      case 'priority':
        return getPriorityValue(b.priority) - getPriorityValue(a.priority)
      default:
        return 0
    }
  })
  
  return result
})

// 获取优先级值
const getPriorityValue = (priority) => {
  const valueMap = {
    low: 1,
    medium: 2,
    high: 3,
    urgent: 4
  }
  return valueMap[priority] || 1
}

// 获取任务图标
const getTaskIcon = (type) => {
  const iconMap = {
    project: Document,
    milestone: Clock,
    file: Document,
    extension: Warning
  }
  return iconMap[type] || Document
}

// 获取任务图标样式类
const getTaskIconClass = (type) => {
  const classMap = {
    project: 'task-icon-project',
    milestone: 'task-icon-milestone',
    file: 'task-icon-file',
    extension: 'task-icon-extension'
  }
  return classMap[type] || 'task-icon-default'
}

// 获取任务类型标签
const getTaskTypeLabel = (type) => {
  const labelMap = {
    project: '项目审核',
    milestone: '里程碑审核',
    file: '文件审核',
    extension: '延期审核'
  }
  return labelMap[type] || '其他'
}

// 获取任务类型标签样式
const getTaskTypeTag = (type) => {
  const tagMap = {
    project: 'primary',
    milestone: 'success',
    file: 'warning',
    extension: 'danger'
  }
  return tagMap[type] || 'info'
}

// 获取优先级类型
const getPriorityType = (priority) => {
  const typeMap = {
    low: 'info',
    medium: 'warning',
    high: 'danger',
    urgent: 'danger'
  }
  return typeMap[priority] || 'info'
}

// 获取优先级标签
const getPriorityLabel = (priority) => {
  const labelMap = {
    low: '低',
    medium: '中',
    high: '高',
    urgent: '紧急'
  }
  return labelMap[priority] || '中'
}

// 获取状态类型
const getStatusType = (status) => {
  const typeMap = {
    pending: 'warning',
    reviewing: 'info',
    completed: 'success'
  }
  return typeMap[status] || 'info'
}

// 获取状态标签
const getStatusLabel = (status) => {
  const labelMap = {
    pending: '待审核',
    reviewing: '审核中',
    completed: '已完成'
  }
  return labelMap[status] || '未知'
}

// 获取截止时间样式类
const getDeadlineClass = (deadline) => {
  const now = new Date()
  const deadlineDate = new Date(deadline)
  const diffDays = Math.ceil((deadlineDate - now) / (1000 * 60 * 60 * 24))
  
  if (diffDays < 0) return 'deadline-overdue'
  if (diffDays <= 1) return 'deadline-urgent'
  if (diffDays <= 3) return 'deadline-warning'
  return 'deadline-normal'
}

// 搜索处理
const handleSearch = () => {
  currentPage.value = 1
}

// 选择处理
const handleSelectionChange = (selection) => {
  selectedTasks.value = selection
}

// 查看任务详情
const viewTaskDetail = (task) => {
  currentTask.value = task
  detailDialogVisible.value = true
}

// 开始审核
const startReview = (task) => {
  currentTask.value = task
  reviewForm.value = {
    action: 'start',
    title: '开始审核',
    content: '',
    result: 'approved',
    revisionRequirements: '',
    nextReviewTime: '',
    score: 8
  }
  reviewDialogVisible.value = true
}

// 完成审核
const completeReview = (task) => {
  currentTask.value = task
  reviewForm.value = {
    action: 'complete',
    title: '完成审核',
    content: '',
    result: 'approved',
    revisionRequirements: '',
    nextReviewTime: '',
    score: 8
  }
  reviewDialogVisible.value = true
}

// 提交审核
const submitReview = async () => {
  try {
    await reviewFormRef.value.validate()
    
    // 这里应该调用实际的API
    // await teacherService.submitTaskReview(currentTask.value.id, reviewForm.value)
    
    // 添加审核记录
    const newNote = {
      id: Date.now(),
      title: reviewForm.value.title,
      content: reviewForm.value.content,
      time: new Date().toLocaleString('zh-CN'),
      reviewer: '当前用户'
    }
    
    if (!currentTask.value.reviewNotes) {
      currentTask.value.reviewNotes = []
    }
    currentTask.value.reviewNotes.push(newNote)
    
    // 更新任务状态
    if (reviewForm.value.action === 'start') {
      currentTask.value.status = 'reviewing'
    } else if (reviewForm.value.action === 'complete') {
      currentTask.value.status = 'completed'
    }
    
    ElMessage.success('审核提交成功')
    reviewDialogVisible.value = false
    detailDialogVisible.value = false
  } catch (error) {
    console.error('提交审核失败:', error)
    ElMessage.error('提交失败')
  }
}

// 批量开始审核
const batchStartReview = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要开始审核选中的 ${selectedTasks.value.length} 个任务吗？`,
      '批量操作确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 这里应该调用实际的API
    // await teacherService.batchStartReview(selectedTasks.value.map(t => t.id))
    
    // 更新本地状态
    selectedTasks.value.forEach(task => {
      task.status = 'reviewing'
      if (!task.reviewNotes) {
        task.reviewNotes = []
      }
      task.reviewNotes.push({
        id: Date.now(),
        title: '批量开始审核',
        content: '通过批量操作开始审核',
        time: new Date().toLocaleString('zh-CN'),
        reviewer: '当前用户'
      })
    })
    
    ElMessage.success('批量开始审核成功')
    selectedTasks.value = []
  } catch (error) {
    if (error !== 'cancel') {
      console.error('批量开始审核失败:', error)
      ElMessage.error('批量操作失败')
    }
  }
}

// 批量完成审核
const batchCompleteReview = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要完成审核选中的 ${selectedTasks.value.length} 个任务吗？`,
      '批量操作确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 这里应该调用实际的API
    // await teacherService.batchCompleteReview(selectedTasks.value.map(t => t.id))
    
    // 更新本地状态
    selectedTasks.value.forEach(task => {
      task.status = 'completed'
      if (!task.reviewNotes) {
        task.reviewNotes = []
      }
      task.reviewNotes.push({
        id: Date.now(),
        title: '批量完成审核',
        content: '通过批量操作完成审核',
        time: new Date().toLocaleString('zh-CN'),
        reviewer: '当前用户'
      })
    })
    
    ElMessage.success('批量完成审核成功')
    selectedTasks.value = []
  } catch (error) {
    if (error !== 'cancel') {
      console.error('批量完成审核失败:', error)
      ElMessage.error('批量操作失败')
    }
  }
}

// 清除选择
const clearSelection = () => {
  selectedTasks.value = []
}

// 下载材料
const downloadMaterial = (material) => {
  ElMessage.success(`开始下载材料: ${material.name}`)
}

// 预览材料
const previewMaterial = (material) => {
  ElMessage.info(`预览材料: ${material.name}`)
}

// 分页处理
const handleSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
}

const handleCurrentChange = (page) => {
  currentPage.value = page
}

// 刷新审核任务
const refreshTasks = () => {
  ElMessage.success('审核任务列表已刷新')
}

// 组件挂载时加载数据
onMounted(() => {
  // 这里可以加载实际的审核任务数据
})
</script>

<style scoped>
.teacher-review-tasks {
  padding: 20px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.task-stats {
  margin-bottom: 20px;
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

.stat-card.pending {
  border-left: 4px solid #f093fb;
}

.stat-card.reviewing {
  border-left: 4px solid #4facfe;
}

.stat-card.completed {
  border-left: 4px solid #43e97b;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 15px;
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

.stat-icon.pending {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.stat-icon.reviewing {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-icon.completed {
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

.search-card {
  margin-bottom: 20px;
}

.search-form {
  margin: 0;
}

.batch-actions {
  margin: 20px 0;
  padding: 15px;
  background-color: #f5f7fa;
  border-radius: 8px;
  text-align: center;
}

.batch-actions .el-button {
  margin: 0 10px;
}

.pagination-container {
  text-align: center;
  margin-top: 20px;
}

.task-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.task-icon {
  font-size: 18px;
}

.task-icon-project {
  color: #409EFF;
}

.task-icon-milestone {
  color: #67C23A;
}

.task-icon-file {
  color: #E6A23C;
}

.task-icon-extension {
  color: #F56C6C;
}

.task-icon-default {
  color: #909399;
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
  color: #606266;
}

.task-detail {
  padding: 20px;
}

.task-detail h4 {
  margin: 20px 0 10px 0;
  color: #2c3e50;
  font-size: 16px;
}

.task-detail p {
  margin: 10px 0;
  color: #606266;
  line-height: 1.6;
}

.task-detail ul {
  margin: 10px 0;
  padding-left: 20px;
}

.task-detail li {
  margin: 5px 0;
  color: #606266;
}

.review-notes h6 {
  margin: 0 0 5px 0;
  color: #2c3e50;
  font-size: 14px;
}

.review-notes p {
  margin: 5px 0;
  color: #606266;
  line-height: 1.6;
}

.dialog-footer {
  text-align: right;
}

:deep(.el-table .el-table__row:hover) {
  background-color: #f5f7fa;
}

:deep(.el-card__header) {
  padding: 15px 20px;
}

:deep(.el-card__body) {
  padding: 20px;
}

:deep(.el-form-item) {
  margin-bottom: 15px;
}

:deep(.el-pagination) {
  justify-content: center;
}
</style> 