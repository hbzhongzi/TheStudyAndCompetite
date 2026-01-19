<template>
  <div class="teacher-project-review">
    <el-card>
      <template #header>
        <div class="header-content">
          <span>项目审核管理</span>
          <el-button type="primary" @click="refreshReviews">刷新</el-button>
        </div>
      </template>

      <!-- 审核统计 -->
      <el-row :gutter="20" class="review-stats">
        <el-col :span="6" v-for="stat in reviewStats" :key="stat.label">
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

      <!-- 筛选和搜索 -->
      <el-card style="margin: 20px 0;">
        <el-row :gutter="20">
          <el-col :span="6">
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
          </el-col>
          <el-col :span="4">
            <el-select v-model="statusFilter" placeholder="审核状态" clearable @change="handleFilter">
              <el-option label="全部状态" :value="''" />
              <el-option label="待审核" value="pending" />
              <el-option label="已通过" value="approved" />
              <el-option label="已拒绝" value="rejected" />
              <el-option label="需修改" value="need_revision" />
            </el-select>
          </el-col>
          <el-col :span="4">
            <el-select v-model="priorityFilter" placeholder="优先级" clearable @change="handleFilter">
              <el-option label="全部优先级" :value="''" />
              <el-option label="低" value="low" />
              <el-option label="中" value="medium" />
              <el-option label="高" value="high" />
              <el-option label="紧急" value="urgent" />
            </el-select>
          </el-col>
          <el-col :span="4">
            <el-select v-model="sortBy" placeholder="排序方式" @change="handleSort">
              <el-option label="提交时间" value="submitTime" />
              <el-option label="优先级" value="priority" />
              <el-option label="项目名称" value="name" />
              <el-option label="学生姓名" value="studentName" />
            </el-select>
          </el-col>
          <el-col :span="6">
            <el-button-group>
              <el-button @click="batchApprove" :disabled="selectedReviews.length === 0" type="success">
                批量通过
              </el-button>
              <el-button @click="batchReject" :disabled="selectedReviews.length === 0" type="danger">
                批量拒绝
              </el-button>
            </el-button-group>
          </el-col>
        </el-row>
      </el-card>

      <!-- 审核列表 -->
      <el-card>
        <template #header>
          <span>审核列表 ({{ filteredReviews.length }})</span>
        </template>
        
        <el-table :data="filteredReviews" @selection-change="handleSelectionChange" v-loading="loading">
          <el-table-column type="selection" width="55" />
          <el-table-column prop="name" label="项目名称" min-width="150">
            <template #default="scope">
              <el-link type="primary" @click="viewProject(scope.row)">{{ scope.row.name }}</el-link>
            </template>
          </el-table-column>
          <el-table-column prop="studentName" label="学生姓名" width="100" />
          <el-table-column prop="type" label="项目类型" width="120" />
          <el-table-column prop="priority" label="优先级" width="100">
            <template #default="scope">
              <el-tag :type="getPriorityType(scope.row.priority)" size="small">
                {{ getPriorityLabel(scope.row.priority) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="审核状态" width="100">
            <template #default="scope">
              <el-tag :type="getStatusType(scope.row.status)" size="small">
                {{ getStatusLabel(scope.row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="submitTime" label="提交时间" width="150" />
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
              <el-button 
                v-if="scope.row.status === 'pending'"
                size="small" 
                type="success" 
                @click="approveProject(scope.row)"
              >
                通过
              </el-button>
              <el-button 
                v-if="scope.row.status === 'pending'"
                size="small" 
                type="danger" 
                @click="rejectProject(scope.row)"
              >
                拒绝
              </el-button>
              <el-button 
                v-if="scope.row.status === 'pending'"
                size="small" 
                type="warning" 
                @click="requestRevision(scope.row)"
              >
                需要修改
              </el-button>
              <el-button size="small" type="info" @click="viewReviewHistory(scope.row)">
                审核历史
              </el-button>
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页 -->
        <div class="pagination-container" v-if="filteredReviews.length > pageSize">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[10, 20, 50]"
            :total="filteredReviews.length"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </el-card>
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
          <el-descriptions-item label="审核状态">{{ getStatusLabel(currentProject.status) }}</el-descriptions-item>
          <el-descriptions-item label="优先级">{{ getPriorityLabel(currentProject.priority) }}</el-descriptions-item>
          <el-descriptions-item label="提交时间">{{ currentProject.submitTime }}</el-descriptions-item>
          <el-descriptions-item label="截止时间">{{ currentProject.deadline }}</el-descriptions-item>
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
        
        <h4>提交材料</h4>
        <div v-if="currentProject.materials && currentProject.materials.length > 0">
          <el-table :data="currentProject.materials" style="width: 100%">
            <el-table-column prop="name" label="材料名称" />
            <el-table-column prop="type" label="类型" />
            <el-table-column prop="size" label="大小" />
            <el-table-column prop="uploadTime" label="上传时间" />
            <el-table-column label="操作" width="150">
              <template #default="scope">
                <el-button size="small" @click="downloadMaterial(scope.row)">下载</el-button>
                <el-button size="small" type="primary" @click="previewMaterial(scope.row)">预览</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
        <el-empty v-else description="暂无提交材料" />
      </div>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="projectDetailVisible = false">关闭</el-button>
          <el-button 
            v-if="currentProject && currentProject.status === 'pending'"
            type="success" 
            @click="approveProject(currentProject)"
          >
            通过审核
          </el-button>
          <el-button 
            v-if="currentProject && currentProject.status === 'pending'"
            type="danger" 
            @click="rejectProject(currentProject)"
          >
            拒绝审核
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 审核操作对话框 -->
    <el-dialog
      v-model="reviewDialogVisible"
      :title="reviewDialogTitle"
      width="60%"
    >
      <el-form :model="reviewForm" :rules="reviewRules" ref="reviewFormRef" label-width="100px">
        <el-form-item label="审核结果" prop="result">
          <el-radio-group v-model="reviewForm.result">
            <el-radio label="approved">通过</el-radio>
            <el-radio label="rejected">拒绝</el-radio>
            <el-radio label="revision">需要修改</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item label="审核意见" prop="comments">
          <el-input
            v-model="reviewForm.comments"
            type="textarea"
            :rows="4"
            placeholder="请详细说明审核意见"
          />
        </el-form-item>
        
        <el-form-item label="修改建议" v-if="reviewForm.result === 'revision'">
          <el-input
            v-model="reviewForm.suggestions"
            type="textarea"
            :rows="3"
            placeholder="请提供具体的修改建议"
          />
        </el-form-item>
        
        <el-form-item label="修改期限" v-if="reviewForm.result === 'revision'">
          <el-date-picker
            v-model="reviewForm.revisionDeadline"
            type="date"
            placeholder="选择修改截止日期"
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
          <el-button type="primary" @click="submitReview">提交审核</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 审核历史对话框 -->
    <el-dialog
      v-model="historyDialogVisible"
      title="审核历史"
      width="70%"
    >
      <div v-if="currentProject && currentProject.reviewHistory">
        <el-timeline>
          <el-timeline-item
            v-for="review in currentProject.reviewHistory"
            :key="review.id"
            :timestamp="review.reviewTime"
            :type="getReviewType(review.result)"
          >
            <el-card class="review-history-card">
              <div class="review-header">
                <h4>{{ getReviewResultLabel(review.result) }}</h4>
                <el-tag :type="getReviewType(review.result)" size="small">
                  {{ getReviewResultLabel(review.result) }}
                </el-tag>
              </div>
              <p><strong>审核人:</strong> {{ review.reviewer }}</p>
              <p><strong>审核时间:</strong> {{ review.reviewTime }}</p>
              <p><strong>审核意见:</strong></p>
              <p class="review-comments">{{ review.comments }}</p>
              <div v-if="review.suggestions">
                <p><strong>修改建议:</strong></p>
                <p class="review-suggestions">{{ review.suggestions }}</p>
              </div>
              <div v-if="review.score">
                <p><strong>评分:</strong> {{ review.score }}/10</p>
              </div>
            </el-card>
          </el-timeline-item>
        </el-timeline>
      </div>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="historyDialogVisible = false">关闭</el-button>
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
const loading = ref(false)
const reviewList = ref([])
const projectDetailVisible = ref(false)
const reviewDialogVisible = ref(false)
const historyDialogVisible = ref(false)
const currentProject = ref(null)
const selectedReviews = ref([])
const searchQuery = ref('')
const statusFilter = ref('')
const priorityFilter = ref('')
const sortBy = ref('submitTime')
const currentPage = ref(1)
const pageSize = ref(10)

const reviewForm = ref({
  result: '',
  comments: '',
  suggestions: '',
  revisionDeadline: '',
  score: 5
})

const reviewRules = {
  result: [{ required: true, message: '请选择审核结果', trigger: 'change' }],
  comments: [{ required: true, message: '请输入审核意见', trigger: 'blur' }]
}

const reviewFormRef = ref()

// 审核统计
const reviewStats = ref([
  {
    title: '待审核',
    value: 0,
    description: '等待审核',
    icon: Clock,
    type: 'pending'
  },
  {
    title: '已通过',
    value: 0,
    description: '审核通过',
    icon: Check,
    type: 'approved'
  },
  {
    title: '已拒绝',
    value: 0,
    description: '审核拒绝',
    icon: Warning,
    type: 'rejected'
  },
  {
    title: '需要修改',
    value: 0,
    description: '要求修改',
    icon: Document,
    type: 'revision'
  }
])

// 计算属性
const filteredReviews = computed(() => {
  let filtered = reviewList.value

  // 搜索过滤
  if (searchQuery.value) {
    filtered = filtered.filter(review => 
      review.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      review.studentName.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      review.description.toLowerCase().includes(searchQuery.value.toLowerCase())
    )
  }

  // 状态过滤
  if (statusFilter.value) {
    filtered = filtered.filter(review => review.status === statusFilter.value)
  }

  // 优先级过滤
  if (priorityFilter.value) {
    filtered = filtered.filter(review => review.priority === priorityFilter.value)
  }

  // 排序
  filtered.sort((a, b) => {
    switch (sortBy.value) {
      case 'name':
        return a.name.localeCompare(b.name)
      case 'priority':
        return getPriorityWeight(b.priority) - getPriorityWeight(a.priority)
      case 'studentName':
        return a.studentName.localeCompare(b.studentName)
      case 'submitTime':
      default:
        return new Date(b.submitTime) - new Date(a.submitTime)
    }
  })

  return filtered
})

// 获取优先级权重
const getPriorityWeight = (priority) => {
  const weightMap = { high: 3, medium: 2, low: 1 }
  return weightMap[priority] || 0
}

// 获取优先级类型
const getPriorityType = (priority) => {
  const typeMap = {
    high: 'danger',
    medium: 'warning',
    low: 'info'
  }
  return typeMap[priority] || 'info'
}

// 获取优先级标签
const getPriorityLabel = (priority) => {
  const labelMap = {
    high: '高',
    medium: '中',
    low: '低'
  }
  return labelMap[priority] || '未知'
}

// 获取状态类型
const getStatusType = (status) => {
  const typeMap = {
    pending: 'warning',
    approved: 'success',
    rejected: 'danger',
    revision: 'info'
  }
  return typeMap[status] || 'info'
}

// 获取状态标签
const getStatusLabel = (status) => {
  const labelMap = {
    pending: '待审核',
    approved: '已通过',
    rejected: '已拒绝',
    revision: '需要修改'
  }
  return labelMap[status] || '未知'
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

// 获取审核类型
const getReviewType = (result) => {
  const typeMap = {
    approved: 'success',
    rejected: 'danger',
    revision: 'warning'
  }
  return typeMap[result] || 'info'
}

// 获取审核结果标签
const getReviewResultLabel = (result) => {
  const labelMap = {
    approved: '审核通过',
    rejected: '审核拒绝',
    revision: '需要修改'
  }
  return labelMap[result] || '未知'
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

// 选择变化处理
const handleSelectionChange = (selection) => {
  selectedReviews.value = selection
}

// 加载审核列表
const loadReviews = async () => {
  loading.value = true
  try {
    const response = await teacherService.getProjectReviews()
    if (response && response.code === 200) {
      reviewList.value = response.data || []
      updateStats()
    } else {
      // 使用模拟数据
      loadMockData()
    }
  } catch (error) {
    console.error('加载审核列表失败:', error)
    // API调用失败时，使用模拟数据
    loadMockData()
  } finally {
    loading.value = false
  }
}

// 加载模拟数据
const loadMockData = () => {
  reviewList.value = [
    {
      id: 1,
      name: '智能校园系统',
      studentName: '张三',
      type: '软件开发',
      priority: 'high',
      status: 'pending',
      submitTime: '2024-06-01 10:30:00',
      deadline: '2024-07-15',
      description: '基于物联网技术的智能校园管理系统',
      plan: '预计6个月完成，分为需求分析、设计、开发、测试四个阶段',
      teacherName: '李教授',
      materials: [
        { name: '项目申请书.pdf', type: 'PDF', size: '2.5MB', uploadTime: '2024-06-01 10:30:00' },
        { name: '技术方案.docx', type: 'Word', size: '1.8MB', uploadTime: '2024-06-01 10:30:00' },
        { name: '项目计划.xlsx', type: 'Excel', size: '0.5MB', uploadTime: '2024-06-01 10:30:00' }
      ],
      milestones: [
        { id: 1, title: '需求分析', description: '完成用户需求调研和分析', date: '2024-01-20', completed: true },
        { id: 2, title: '系统设计', description: '完成系统架构和数据库设计', date: '2024-02-15', completed: true },
        { id: 3, title: '功能开发', description: '核心功能模块开发', date: '2024-05-15', completed: false },
        { id: 4, title: '系统测试', description: '功能测试和性能测试', date: '2024-06-15', completed: false }
      ],
      reviewHistory: []
    },
    {
      id: 2,
      name: '数据分析平台',
      studentName: '李四',
      type: '科研项目',
      priority: 'medium',
      status: 'approved',
      submitTime: '2024-05-15 14:20:00',
      deadline: '2024-06-14',
      description: '大数据分析平台，支持多种数据源和算法',
      plan: '预计8个月完成，包括数据采集、预处理、分析、可视化等模块',
      teacherName: '王教授',
      materials: [
        { name: '研究报告.pdf', type: 'PDF', size: '3.2MB', uploadTime: '2024-05-15 14:20:00' },
        { name: '实验数据.xlsx', type: 'Excel', size: '1.5MB', uploadTime: '2024-05-15 14:20:00' }
      ],
      milestones: [
        { id: 1, title: '数据采集', description: '建立数据采集管道', date: '2024-01-25', completed: true },
        { id: 2, title: '算法实现', description: '核心分析算法开发', date: '2024-03-15', completed: true },
        { id: 3, title: '平台集成', description: '各模块集成测试', date: '2024-05-15', completed: false }
      ],
      reviewHistory: [
        {
          id: 1,
          result: 'approved',
          comments: '项目设计合理，技术方案可行，同意通过',
          suggestions: '',
          score: 8,
          reviewer: '王教授',
          reviewTime: '2024-05-20 09:15:00'
        }
      ]
    },
    {
      id: 3,
      name: '在线教育平台',
      studentName: '王五',
      type: '创新项目',
      priority: 'low',
      status: 'revision',
      submitTime: '2024-05-10 16:45:00',
      deadline: '2024-06-10',
      description: '基于Web的在线教育学习平台',
      plan: '预计4个月完成，包括用户管理、课程管理、学习跟踪等模块',
      teacherName: '赵教授',
      materials: [
        { name: '项目提案.docx', type: 'Word', size: '2.1MB', uploadTime: '2024-05-10 16:45:00' },
        { name: 'UI设计稿.png', type: '图片', size: '0.8MB', uploadTime: '2024-05-10 16:45:00' }
      ],
      milestones: [
        { id: 1, title: '用户系统', description: '用户注册登录功能', date: '2024-02-10', completed: true },
        { id: 2, title: '课程管理', description: '课程创建和发布功能', date: '2024-03-10', completed: true },
        { id: 3, title: '学习跟踪', description: '学习进度和成绩统计', date: '2024-04-10', completed: true }
      ],
      reviewHistory: [
        {
          id: 1,
          result: 'revision',
          comments: '项目整体设计不错，但需要完善技术实现细节',
          suggestions: '建议补充技术架构图和数据库设计文档',
          score: 6,
          reviewer: '赵教授',
          reviewTime: '2024-05-15 11:30:00'
        }
      ]
    }
  ]
  updateStats()
}

// 更新统计数据
const updateStats = () => {
  const pending = reviewList.value.filter(r => r.status === 'pending').length
  const approved = reviewList.value.filter(r => r.status === 'approved').length
  const rejected = reviewList.value.filter(r => r.status === 'rejected').length
  const revision = reviewList.value.filter(r => r.status === 'revision').length
  
  reviewStats.value[0].value = pending
  reviewStats.value[1].value = approved
  reviewStats.value[2].value = rejected
  reviewStats.value[3].value = revision
}

// 查看项目详情
const viewProject = (project) => {
  currentProject.value = project
  projectDetailVisible.value = true
}

// 通过项目
const approveProject = (project) => {
  currentProject.value = project
  reviewForm.value = {
    result: 'approved',
    comments: '',
    suggestions: '',
    revisionDeadline: '',
    score: 8
  }
  reviewDialogTitle.value = '审核通过'
  reviewDialogVisible.value = true
}

// 拒绝项目
const rejectProject = (project) => {
  currentProject.value = project
  reviewForm.value = {
    result: 'rejected',
    comments: '',
    suggestions: '',
    revisionDeadline: '',
    score: 0
  }
  reviewDialogTitle.value = '审核拒绝'
  reviewDialogVisible.value = true
}

// 要求修改
const requestRevision = (project) => {
  currentProject.value = project
  reviewForm.value = {
    result: 'revision',
    comments: '',
    suggestions: '',
    revisionDeadline: '',
    score: 0
  }
  reviewDialogTitle.value = '要求修改'
  reviewDialogVisible.value = true
}

// 查看审核历史
const viewReviewHistory = (project) => {
  currentProject.value = project
  historyDialogVisible.value = true
}

// 下载材料
const downloadMaterial = (material) => {
  ElMessage.success(`开始下载材料: ${material.name}`)
}

// 预览材料
const previewMaterial = (material) => {
  ElMessage.info(`预览材料: ${material.name}`)
}

// 提交审核
const submitReview = async () => {
  try {
    await reviewFormRef.value.validate()
    
    // 这里应该调用实际的API
    // await teacherService.submitProjectReview(currentProject.value.id, reviewForm.value)
    
    // 添加审核记录
    const newReview = {
      id: Date.now(),
      result: reviewForm.value.result,
      comments: reviewForm.value.comments,
      suggestions: reviewForm.value.suggestions,
      score: reviewForm.value.score,
      reviewer: '当前用户',
      reviewTime: new Date().toLocaleString('zh-CN')
    }
    
    if (!currentProject.value.reviewHistory) {
      currentProject.value.reviewHistory = []
    }
    currentProject.value.reviewHistory.push(newReview)
    
    // 更新项目状态
    currentProject.value.status = reviewForm.value.result
    
    ElMessage.success('审核提交成功')
    reviewDialogVisible.value = false
    updateStats()
  } catch (error) {
    console.error('提交审核失败:', error)
    ElMessage.error('提交失败')
  }
}

// 批量通过
const batchApprove = async () => {
  try {
    await ElMessageBox.confirm(`确定要批量通过选中的 ${selectedReviews.value.length} 个项目吗？`, '确认操作', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'info'
    })
    
    // 这里应该调用实际的API
    selectedReviews.value.forEach(review => {
      review.status = 'approved'
    })
    
    ElMessage.success('批量审核通过成功')
    updateStats()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('批量操作失败:', error)
      ElMessage.error('操作失败')
    }
  }
}

// 批量拒绝
const batchReject = async () => {
  try {
    await ElMessageBox.confirm(`确定要批量拒绝选中的 ${selectedReviews.value.length} 个项目吗？`, '确认操作', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    // 这里应该调用实际的API
    selectedReviews.value.forEach(review => {
      review.status = 'rejected'
    })
    
    ElMessage.success('批量审核拒绝成功')
    updateStats()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('批量操作失败:', error)
      ElMessage.error('操作失败')
    }
  }
}

// 刷新审核列表
const refreshReviews = () => {
  loadReviews()
  ElMessage.success('审核列表已刷新')
}

// 组件挂载时加载数据
onMounted(() => {
  loadReviews()
})
</script>

<style scoped>
.teacher-project-review {
  padding: 20px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.review-stats {
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

.stat-card.pending {
  border-left: 4px solid #f093fb;
}

.stat-card.approved {
  border-left: 4px solid #4facfe;
}

.stat-card.rejected {
  border-left: 4px solid #43e97b;
}

.stat-card.revision {
  border-left: 4px solid #667eea;
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

.stat-icon.pending {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.stat-icon.approved {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-icon.rejected {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.stat-icon.revision {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
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

.review-history-card {
  margin-bottom: 10px;
}

.review-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.review-header h4 {
  margin: 0;
  color: #2c3e50;
}

.review-comments, .review-suggestions {
  background-color: #f5f7fa;
  padding: 10px;
  border-radius: 4px;
  margin: 10px 0;
  line-height: 1.6;
}
</style> 