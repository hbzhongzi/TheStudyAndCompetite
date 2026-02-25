<template>
  <div class="application-review">
    <!-- 搜索和筛选 -->
    <el-card class="filter-card">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-input
            v-model="searchQuery"
            placeholder="搜索申请标题或学生姓名"
            clearable
            @input="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-col>
        <el-col :span="4">
          <el-select v-model="typeFilter" placeholder="类型筛选" clearable @change="handleFilter">
            <el-option label="全部" value="" />
            <el-option label="项目申请" value="项目申请" />
            <el-option label="竞赛报名" value="竞赛报名" />
            <el-option label="其他申请" value="其他申请" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-select v-model="statusFilter" placeholder="状态筛选" clearable @change="handleFilter">
            <el-option label="全部" value="" />
            <el-option label="待审核" value="pending" />
            <el-option label="已通过" value="approved" />
            <el-option label="已拒绝" value="rejected" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-button type="primary" @click="loadApplications">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 申请列表 -->
    <el-card class="application-list">
      <template #header>
        <div class="card-header">
          <span>申请列表 ({{ filteredApplications.length }})</span>
        </div>
      </template>
      
      <el-table
        :data="filteredApplications"
        v-loading="loading"
        style="width: 100%"
        @row-click="handleRowClick"
      >
        <el-table-column prop="title" label="申请标题" min-width="200" />
        <el-table-column prop="studentName" label="申请人" width="100" />
        <el-table-column prop="advisorName" label="指导老师" width="120" />
        <el-table-column prop="type" label="类型" width="100" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="submitTime" label="提交时间" width="160" />
        <el-table-column prop="reviewTime" label="审核时间" width="160" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click.stop="viewApplication(row)">查看</el-button>
            <el-button 
              v-if="row.status === 'pending'"
              size="small" 
              type="success" 
              @click.stop="approveApplication(row)"
            >
              通过
            </el-button>
            <el-button 
              v-if="row.status === 'pending'"
              size="small" 
              type="danger" 
              @click.stop="rejectApplication(row)"
            >
              驳回
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 申请详情对话框 -->
    <el-dialog
      v-model="detailVisible"
      title="申请详情"
      width="60%"
      :before-close="handleCloseDetail"
    >
      <div v-if="currentApplication" class="application-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="申请标题">{{ currentApplication.title }}</el-descriptions-item>
          <el-descriptions-item label="申请人">{{ currentApplication.studentName }}</el-descriptions-item>
          <el-descriptions-item label="指导老师">{{ currentApplication.advisorName || '未指定' }}</el-descriptions-item>
          <el-descriptions-item label="申请类型">{{ currentApplication.type }}</el-descriptions-item>
          <el-descriptions-item label="申请状态">
            <el-tag :type="getStatusType(currentApplication.status)">
              {{ getStatusText(currentApplication.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="提交时间">{{ formatDate(currentApplication.submitTime) }}</el-descriptions-item>
          <el-descriptions-item label="审核时间">{{ formatDate(currentApplication.reviewTime) }}</el-descriptions-item>
          <el-descriptions-item label="申请内容" :span="2">{{ currentApplication.content }}</el-descriptions-item>
        </el-descriptions>

        <!-- 申请附件 -->
        <div class="section">
          <h4>申请附件</h4>
          <el-table :data="currentApplication.attachments || []" style="width: 100%">
            <el-table-column prop="fileName" label="文件名" />
            <el-table-column prop="fileSize" label="文件大小" />
            <el-table-column label="操作">
              <template #default="{ row }">
                <el-button size="small" @click="downloadFile(row)">下载</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <!-- 审核记录 -->
        <div class="section" v-if="currentApplication.reviewRecords && currentApplication.reviewRecords.length > 0">
          <h4>审核记录</h4>
          <el-timeline>
            <el-timeline-item
              v-for="record in currentApplication.reviewRecords"
              :key="record.id"
              :timestamp="formatDate(record.reviewTime)"
              :type="record.status === 'approved' ? 'success' : 'danger'"
            >
              <el-card>
                <h4>{{ record.reviewer }} - {{ record.status === 'approved' ? '通过' : '拒绝' }}</h4>
                <p>{{ record.comments }}</p>
              </el-card>
            </el-timeline-item>
          </el-timeline>
        </div>
      </div>
    </el-dialog>

    <!-- 审核对话框 -->
    <el-dialog
      v-model="reviewVisible"
      :title="reviewType === 'approve' ? '通过申请' : '驳回申请'"
      width="40%"
    >
      <el-form :model="reviewForm" label-width="80px">
        <el-form-item label="审核意见">
          <el-input
            v-model="reviewForm.comment"
            type="textarea"
            :rows="4"
            :placeholder="reviewType === 'approve' ? '请输入通过意见（可选）' : '请输入驳回原因'"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="reviewVisible = false">取消</el-button>
          <el-button type="primary" @click="submitReview" :loading="submitting">
            确认
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Refresh } from '@element-plus/icons-vue'
import teacherService from '@/services/teacherService'
import { ensureArray, validateApiResponse } from '../../utils/dataValidator'

// 响应式数据
const loading = ref(false)
const searchQuery = ref('')
const typeFilter = ref('')
const statusFilter = ref('')
const applications = ref([]) // 确保初始化为空数组
const detailVisible = ref(false)
const reviewVisible = ref(false)
const currentApplication = ref(null)
const reviewType = ref('approve')
const submitting = ref(false)

const reviewForm = ref({
  comment: ''
})

// 模拟申请数据
const mockApplications = [
  {
    id: 1,
    title: '智能校园管理系统项目申请',
    studentName: '张三',
    type: '项目申请',
    status: 'pending',
    content: '申请开展智能校园管理系统项目，该项目将基于物联网技术，实现校园设施的智能化管理，包括照明控制、安防监控、环境监测等功能。',
    submitTime: '2024-01-15 10:30:00',
    reviewTime: null,
    attachments: [
      { fileName: '项目申请书.pdf', fileSize: '2.5MB' },
      { fileName: '技术方案.docx', fileSize: '1.8MB' }
    ],
    reviewRecords: []
  },
  {
    id: 2,
    title: '全国大学生程序设计竞赛报名',
    studentName: '李四',
    type: '竞赛报名',
    status: 'approved',
    content: '申请参加全国大学生程序设计竞赛，希望通过竞赛提升算法和编程能力，为学校争光。',
    submitTime: '2024-01-10 14:20:00',
    reviewTime: '2024-01-12 09:15:00',
    attachments: [
      { fileName: '竞赛报名表.pdf', fileSize: '0.8MB' }
    ],
    reviewRecords: [
      {
        id: 1,
        reviewer: '王教授',
        status: 'approved',
        comments: '学生基础扎实，有较强的编程能力，同意参加竞赛。',
        reviewTime: '2024-01-12 09:15:00'
      }
    ]
  },
  {
    id: 3,
    title: '在线学习平台项目申请',
    studentName: '王五',
    type: '项目申请',
    status: 'rejected',
    content: '申请开发在线学习平台，该平台将支持课程管理、在线学习、作业提交等功能。',
    submitTime: '2024-01-08 16:45:00',
    reviewTime: '2024-01-09 11:30:00',
    attachments: [
      { fileName: '项目计划书.pdf', fileSize: '3.2MB' },
      { fileName: '需求分析.docx', fileSize: '2.1MB' }
    ],
    reviewRecords: [
      {
        id: 2,
        reviewer: '李教授',
        status: 'rejected',
        comments: '项目规模过大，技术难度较高，建议先从小项目开始。',
        reviewTime: '2024-01-09 11:30:00'
      }
    ]
  }
]

// 计算属性
const filteredApplications = computed(() => {
  // 确保applications.value是数组
  let result = Array.isArray(applications.value) ? applications.value : []

  // 类型筛选
  if (typeFilter.value) {
    result = result.filter(a => a && a.type === typeFilter.value)
  }

  // 状态筛选
  if (statusFilter.value) {
    result = result.filter(a => a && a.status === statusFilter.value)
  }

  // 搜索筛选
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(a => 
      a && a.title && a.title.toLowerCase().includes(query) ||
      a && a.studentName && a.studentName.toLowerCase().includes(query)
    )
  }

  return result
})

// 方法
const loadApplications = async () => {
  loading.value = true
  try {
    // 获取项目和竞赛申请
    const params = {}
    if (statusFilter.value) params.status = statusFilter.value
    if (typeFilter.value) params.type = typeFilter.value
    if (searchQuery.value) params.search = searchQuery.value
    
    const [projectsResponse, competitionsResponse] = await Promise.all([
      teacherService.getMyProjects(params),
      teacherService.getCompetitions(params)
    ])
    
    let allApplications = []
    
    // 处理项目申请 - 使用我指导的项目
    const projectsValidation = validateApiResponse(projectsResponse)
    if (projectsValidation.isValid) {
      const projectList = ensureArray(projectsValidation.data)
      
      const projectApplications = projectList.map(project => ({
        id: `project_${project.id}`,
        title: `${project.title} - 项目申请`,
        studentName: project.studentName || project.creatorName,
        advisorName: project.advisorName || '未指定',
        type: '项目申请',
        status: project.status,
        content: project.description,
        submitTime: project.createTime,
        reviewTime: project.updateTime,
        attachments: ensureArray(project.attachments),
        reviewRecords: ensureArray(project.reviewRecords),
        originalData: project
      }))
      allApplications = allApplications.concat(projectApplications)
    }
    
    // 处理竞赛申请
    const competitionsValidation = validateApiResponse(competitionsResponse)
    if (competitionsValidation.isValid) {
      const competitionList = ensureArray(competitionsValidation.data)
      
      const competitionApplications = competitionList.map(competition => ({
        id: `competition_${competition.id}`,
        title: `${competition.title} - 竞赛报名`,
        studentName: competition.studentName || competition.participantName,
        type: '竞赛报名',
        status: competition.status,
        content: competition.description,
        submitTime: competition.createTime,
        reviewTime: competition.updateTime,
        attachments: ensureArray(competition.attachments),
        reviewRecords: ensureArray(competition.reviewRecords),
        originalData: competition
      }))
      allApplications = allApplications.concat(competitionApplications)
    }
    
    applications.value = allApplications
    ElMessage.success('申请列表加载成功')
  } catch (error) {
    console.error('加载申请列表失败:', error)
    applications.value = [] // 确保失败时也设置为空数组
    ElMessage.error(error.message || '加载申请列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  // 搜索逻辑已在计算属性中处理
}

const handleFilter = () => {
  // 筛选逻辑已在计算属性中处理
}

const handleRowClick = (row) => {
  viewApplication(row)
}

const viewApplication = (application) => {
  currentApplication.value = application
  detailVisible.value = true
}

const handleCloseDetail = () => {
  detailVisible.value = false
  currentApplication.value = null
}

const approveApplication = (application) => {
  reviewType.value = 'approve'
  currentApplication.value = application
  reviewForm.value.comment = ''
  reviewVisible.value = true
}

const rejectApplication = (application) => {
  reviewType.value = 'reject'
  currentApplication.value = application
  reviewForm.value.comment = ''
  reviewVisible.value = true
}

const submitReview = async () => {
  if (reviewType.value === 'reject' && !reviewForm.value.comment.trim()) {
    ElMessage.warning('请输入驳回原因')
    return
  }

  submitting.value = true
  try {
    const reviewData = {
      status: reviewType.value === 'approve' ? 'approved' : 'rejected',
      comments: reviewForm.value.comment || (reviewType.value === 'approve' ? '审核通过' : '审核不通过')
    }
    
    let response
    const originalData = currentApplication.value.originalData
    
    // 根据申请类型调用不同的审核API
    if (currentApplication.value.type === '项目申请') {
      response = await teacherService.reviewProject(originalData.id, reviewData)
    } else if (currentApplication.value.type === '竞赛报名') {
      response = await teacherService.reviewCompetition(originalData.id, reviewData)
    }
    
    if (response && response.code === 200) {
      // 更新本地数据
      const application = applications.value.find(a => a.id === currentApplication.value.id)
      if (application) {
        application.status = reviewType.value === 'approve' ? 'approved' : 'rejected'
        application.reviewTime = new Date().toLocaleString()
      }

      ElMessage.success(reviewType.value === 'approve' ? '申请审核通过' : '申请已驳回')
      reviewVisible.value = false
      currentApplication.value = null
      
      // 重新加载申请列表
      loadApplications()
    } else {
      ElMessage.error(response?.message || '审核操作失败')
    }
  } catch (error) {
    console.error('审核操作失败:', error)
    ElMessage.error(error.message || '审核操作失败')
  } finally {
    submitting.value = false
  }
}

const getStatusType = (status) => {
  const statusMap = {
    pending: 'warning',
    approved: 'success',
    rejected: 'danger'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status) => {
  const statusMap = {
    pending: '待审核',
    approved: '已通过',
    rejected: '已驳回'
  }
  return statusMap[status] || status
}

const formatDate = (dateString) => {
  if (!dateString) return '暂无'
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

// 组件挂载时加载数据
onMounted(() => {
  loadApplications()
})
</script>

<style scoped>
.application-review {
  padding: 20px;
}

.filter-card {
  margin-bottom: 20px;
}

.application-list {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.application-detail {
  max-height: 60vh;
  overflow-y: auto;
}

.section {
  margin-top: 20px;
}

.section h4 {
  margin-bottom: 10px;
  color: #2c3e50;
  font-size: 16px;
  font-weight: 600;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style> 