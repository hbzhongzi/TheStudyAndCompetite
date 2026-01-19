<template>
  <div class="competition-guidance">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h2>竞赛指导</h2>
        <p class="header-desc">管理指导学生的竞赛报名、审核、作品评审和成绩查看</p>
      </div>
      <div class="header-right">
        <el-button type="primary" @click="exportGuidanceData">
          <el-icon><Download /></el-icon>
          导出指导数据
        </el-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon pending">
              <el-icon><Clock /></el-icon>
            </div>
            <div class="stat-info">
              <h4>待审核报名</h4>
              <p class="stat-number">{{ stats.pendingCount }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon reviewing">
              <el-icon><Document /></el-icon>
            </div>
            <div class="stat-info">
              <h4>待评审作品</h4>
              <p class="stat-number">{{ stats.reviewingCount }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon completed">
              <el-icon><Check /></el-icon>
            </div>
            <div class="stat-info">
              <h4>已完成指导</h4>
              <p class="stat-number">{{ stats.completedCount }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon total">
              <el-icon><DataAnalysis /></el-icon>
            </div>
            <div class="stat-info">
              <h4>总指导项目</h4>
              <p class="stat-number">{{ stats.totalCount }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 筛选和搜索 -->
    <el-card class="filter-card">
      <template #header>
        <div class="filter-header">
          <span>筛选条件</span>
          <el-button link @click="resetFilters">重置筛选</el-button>
        </div>
      </template>
      
      <el-row :gutter="20">
        <el-col :span="6">
          <el-input
            v-model="filters.search"
            placeholder="搜索竞赛名称、学生姓名"
            clearable
            @input="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-col>
        <el-col :span="4">
          <el-select v-model="filters.status" placeholder="报名状态" clearable @change="handleSearch">
            <el-option label="全部" value="" />
            <el-option label="待审核" value="pending" />
            <el-option label="已通过" value="approved" />
            <el-option label="已驳回" value="rejected" />
            <el-option label="已提交作品" value="submitted" />
            <el-option label="评审中" value="reviewing" />
            <el-option label="已完成" value="completed" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-select v-model="filters.competitionLevel" placeholder="竞赛级别" clearable @change="handleSearch">
            <el-option label="全部" value="" />
            <el-option label="校级" value="校级" />
            <el-option label="省级" value="省级" />
            <el-option label="国家级" value="国家级" />
            <el-option label="国际级" value="国际级" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-select v-model="filters.student" placeholder="选择学生" clearable @change="handleSearch">
            <el-option label="全部学生" value="" />
            <el-option
              v-for="student in studentList"
              :key="student.id"
              :label="student.name"
              :value="student.id"
            />
          </el-select>
        </el-col>
        <el-col :span="6">
          <el-button type="primary" @click="loadCompetitionRegistrations">
            <el-icon><Refresh /></el-icon>
            刷新数据
          </el-button>
          <el-button @click="showAdvancedFilter = !showAdvancedFilter">
            <el-icon><Filter /></el-icon>
            高级筛选
          </el-button>
        </el-col>
      </el-row>

      <!-- 高级筛选 -->
      <el-collapse-transition>
        <div v-show="showAdvancedFilter" class="advanced-filter">
          <el-row :gutter="20">
            <el-col :span="6">
              <el-date-picker
                v-model="filters.dateRange"
                type="daterange"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                format="YYYY-MM-DD"
                value-format="YYYY-MM-DD"
                @change="handleSearch"
              />
            </el-col>
            <el-col :span="4">
              <el-select v-model="filters.hasSubmission" placeholder="作品提交" clearable @change="handleSearch">
                <el-option label="全部" value="" />
                <el-option label="已提交" value="true" />
                <el-option label="未提交" value="false" />
              </el-select>
            </el-col>
            <el-col :span="4">
              <el-select v-model="filters.hasReview" placeholder="评审状态" clearable @change="handleSearch">
                <el-option label="全部" value="" />
                <el-option label="已评审" value="true" />
                <el-option label="未评审" value="false" />
              </el-select>
            </el-col>
          </el-row>
        </div>
      </el-collapse-transition>
    </el-card>

    <!-- 竞赛报名列表 -->
    <el-card class="registration-list">
      <template #header>
        <div class="card-header">
          <span>指导学生的竞赛报名 ({{ filteredRegistrations.length }})</span>
          <div class="header-actions">
            <el-button size="small" @click="batchApprove" :disabled="!hasSelectedPending">
              批量通过
            </el-button>
            <el-button size="small" type="danger" @click="batchReject" :disabled="!hasSelectedPending">
              批量驳回
            </el-button>
          </div>
        </div>
      </template>
      
      <el-table
        :data="filteredRegistrations"
        v-loading="loading"
        style="width: 100%"
        @selection-change="handleSelectionChange"
        @row-click="handleRowClick"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="competitionName" label="竞赛名称" min-width="200" />
        <el-table-column prop="studentName" label="学生姓名" width="100" />
        <el-table-column prop="studentNumber" label="学号" width="120" />
        <el-table-column prop="competitionLevel" label="级别" width="80" />
        <el-table-column prop="registrationStatus" label="报名状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.registrationStatus)">
              {{ getStatusText(row.registrationStatus) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="submissionStatus" label="作品状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getSubmissionStatusType(row.submissionStatus)">
              {{ getSubmissionStatusText(row.submissionStatus) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="reviewStatus" label="评审状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getReviewStatusType(row.reviewStatus)">
              {{ getReviewStatusText(row.reviewStatus) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="submitTime" label="报名时间" width="160" />
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click.stop="viewRegistrationDetail(row)">查看详情</el-button>
            <el-button 
              v-if="row.registrationStatus === 'pending'"
              size="small" 
              type="success" 
              @click.stop="reviewRegistration(row, 'approve')"
            >
              通过
            </el-button>
            <el-button 
              v-if="row.registrationStatus === 'pending'"
              size="small" 
              type="danger" 
              @click.stop="reviewRegistration(row, 'reject')"
            >
              驳回
            </el-button>
            <el-button 
              v-if="row.submissionStatus === 'submitted'"
              size="small" 
              type="primary" 
              @click.stop="reviewSubmission(row)"
            >
              评审作品
            </el-button>
            <el-button 
              v-if="row.reviewStatus === 'completed'"
              size="small" 
              type="info" 
              @click.stop="viewResult(row)"
            >
              查看成绩
            </el-button>
            <el-dropdown @command="(command) => handleCommand(command, row)">
              <el-button size="small">
                更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="view-files">查看附件</el-dropdown-item>
                  <el-dropdown-item command="view-review-records">审核记录</el-dropdown-item>
                  <el-dropdown-item command="download-certificate">下载证书</el-dropdown-item>
                  <el-dropdown-item command="guidance-record" divided>指导记录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="totalRegistrations"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 报名详情对话框 -->
    <el-dialog
      v-model="detailVisible"
      title="报名详情"
      width="80%"
      :close-on-click-modal="false"
    >
      <div v-if="currentRegistration" class="registration-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="竞赛名称">{{ currentRegistration.competitionName }}</el-descriptions-item>
          <el-descriptions-item label="竞赛级别">{{ currentRegistration.competitionLevel }}</el-descriptions-item>
          <el-descriptions-item label="学生姓名">{{ currentRegistration.studentName }}</el-descriptions-item>
          <el-descriptions-item label="学号">{{ currentRegistration.studentNumber }}</el-descriptions-item>
          <el-descriptions-item label="报名状态">
            <el-tag :type="getStatusType(currentRegistration.registrationStatus)">
              {{ getStatusText(currentRegistration.registrationStatus) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="报名时间">{{ currentRegistration.submitTime }}</el-descriptions-item>
          <el-descriptions-item label="项目描述" :span="2">{{ currentRegistration.description }}</el-descriptions-item>
        </el-descriptions>

        <!-- 团队成员 -->
        <div class="section">
          <h4>团队成员</h4>
          <el-table :data="currentRegistration.members || []" style="width: 100%">
            <el-table-column prop="name" label="姓名" />
            <el-table-column prop="studentNumber" label="学号" />
            <el-table-column prop="role" label="角色" />
            <el-table-column prop="major" label="专业" />
          </el-table>
        </div>

        <!-- 附件列表 -->
        <div class="section">
          <h4>报名附件</h4>
          <el-table :data="currentRegistration.attachments || []" style="width: 100%">
            <el-table-column prop="fileName" label="文件名" />
            <el-table-column prop="fileSize" label="文件大小" width="120" />
            <el-table-column prop="uploadTime" label="上传时间" width="160" />
            <el-table-column label="操作" width="120">
              <template #default="{ row }">
                <el-button size="small" @click="downloadFile(row)">下载</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <!-- 审核记录 -->
        <div class="section">
          <h4>审核记录</h4>
          <el-timeline>
            <el-timeline-item
              v-for="record in currentRegistration.reviewRecords || []"
              :key="record.id"
              :timestamp="record.reviewTime"
              :type="record.status === 'approved' ? 'success' : 'danger'"
            >
              <h5>{{ record.reviewerName }} - {{ record.status === 'approved' ? '通过' : '驳回' }}</h5>
              <p>{{ record.comments }}</p>
            </el-timeline-item>
          </el-timeline>
        </div>
      </div>
    </el-dialog>

    <!-- 审核对话框 -->
    <el-dialog
      v-model="reviewVisible"
      :title="reviewType === 'approve' ? '审核通过' : '审核驳回'"
      width="50%"
    >
      <el-form :model="reviewForm" label-width="100px">
        <el-form-item label="审核意见" required>
          <el-input
            v-model="reviewForm.comments"
            type="textarea"
            :rows="4"
            :placeholder="reviewType === 'approve' ? '请输入通过意见（可选）' : '请输入驳回原因'"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="reviewVisible = false">取消</el-button>
          <el-button 
            type="primary" 
            @click="submitReview" 
            :loading="submitting"
            :disabled="reviewType === 'reject' && !reviewForm.comments.trim()"
          >
            确认
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 作品评审对话框 -->
    <el-dialog
      v-model="submissionReviewVisible"
      title="作品评审"
      width="70%"
    >
      <div v-if="currentSubmission" class="submission-review">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="竞赛名称">{{ currentSubmission.competitionName }}</el-descriptions-item>
          <el-descriptions-item label="学生姓名">{{ currentSubmission.studentName }}</el-descriptions-item>
          <el-descriptions-item label="提交时间">{{ currentSubmission.submitTime }}</el-descriptions-item>
          <el-descriptions-item label="作品标题">{{ currentSubmission.title }}</el-descriptions-item>
          <el-descriptions-item label="作品描述" :span="2">{{ currentSubmission.description }}</el-descriptions-item>
        </el-descriptions>

        <!-- 作品文件 -->
        <div class="section">
          <h4>作品文件</h4>
          <el-table :data="currentSubmission.files || []" style="width: 100%">
            <el-table-column prop="fileName" label="文件名" />
            <el-table-column prop="fileSize" label="文件大小" width="120" />
            <el-table-column prop="uploadTime" label="上传时间" width="160" />
            <el-table-column label="操作" width="120">
              <template #default="{ row }">
                <el-button size="small" @click="downloadFile(row)">下载</el-button>
                <el-button size="small" @click="previewFile(row)">预览</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <!-- 评审表单 -->
        <div class="section">
          <h4>评审意见</h4>
          <el-form :model="submissionReviewForm" label-width="100px">
            <el-form-item label="评分" required>
              <el-rate
                v-model="submissionReviewForm.score"
                :max="10"
                :texts="['很差', '较差', '一般', '良好', '优秀', '很好', '非常好', '极好', '完美', '卓越']"
                show-text
              />
            </el-form-item>
            <el-form-item label="评审意见" required>
              <el-input
                v-model="submissionReviewForm.comments"
                type="textarea"
                :rows="4"
                placeholder="请输入详细的评审意见"
              />
            </el-form-item>
            <el-form-item label="建议改进">
              <el-input
                v-model="submissionReviewForm.suggestions"
                type="textarea"
                :rows="3"
                placeholder="请输入改进建议（可选）"
              />
            </el-form-item>
          </el-form>
        </div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="submissionReviewVisible = false">取消</el-button>
          <el-button 
            type="primary" 
            @click="submitSubmissionReview" 
            :loading="submitting"
            :disabled="!submissionReviewForm.score || !submissionReviewForm.comments.trim()"
          >
            提交评审
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 成绩查看对话框 -->
    <el-dialog
      v-model="resultVisible"
      title="竞赛成绩"
      width="60%"
    >
      <div v-if="currentResult" class="competition-result">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="竞赛名称">{{ currentResult.competitionName }}</el-descriptions-item>
          <el-descriptions-item label="学生姓名">{{ currentResult.studentName }}</el-descriptions-item>
          <el-descriptions-item label="最终成绩">{{ currentResult.finalScore }}</el-descriptions-item>
          <el-descriptions-item label="获奖等级">
            <el-tag :type="getAwardType(currentResult.awardLevel)">
              {{ currentResult.awardLevel }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="排名">{{ currentResult.ranking }}</el-descriptions-item>
          <el-descriptions-item label="公布时间">{{ currentResult.publishTime }}</el-descriptions-item>
        </el-descriptions>

        <!-- 评审详情 -->
        <div class="section">
          <h4>评审详情</h4>
          <el-table :data="currentResult.reviewDetails || []" style="width: 100%">
            <el-table-column prop="reviewerName" label="评审教师" />
            <el-table-column prop="score" label="评分" />
            <el-table-column prop="comments" label="评审意见" />
            <el-table-column prop="reviewTime" label="评审时间" width="160" />
          </el-table>
        </div>

        <!-- 证书下载 -->
        <div class="section">
          <h4>证书下载</h4>
          <el-button type="primary" @click="downloadCertificate(currentResult.registrationId)">
            <el-icon><Download /></el-icon>
            下载获奖证书
          </el-button>
          <el-button @click="downloadTranscript">
            <el-icon><Document /></el-icon>
            下载成绩单
          </el-button>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Search, Refresh, Filter, Download, Clock, Document, Check, DataAnalysis,
  ArrowDown, Plus
} from '@element-plus/icons-vue'
import { ensureArray, validateApiResponse } from '../../utils/dataValidator'
import { teacherService, fileService } from '../../services/teacherService'

// 响应式数据
const loading = ref(false)
const showAdvancedFilter = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)
const totalRegistrations = ref(0)
const selectedRegistrations = ref([])

// 筛选条件
const filters = ref({
  search: '',
  status: '',
  competitionLevel: '',
  student: '',
  dateRange: [],
  hasSubmission: '',
  hasReview: ''
})

// 统计数据
const stats = ref({
  pendingCount: 0,
  reviewingCount: 0,
  completedCount: 0,
  totalCount: 0
})

// 学生列表
const studentList = ref([])

// 竞赛报名数据
const registrations = ref([])

// 对话框状态
const detailVisible = ref(false)
const reviewVisible = ref(false)
const submissionReviewVisible = ref(false)
const resultVisible = ref(false)

// 当前操作的数据
const currentRegistration = ref(null)
const currentSubmission = ref(null)
const currentResult = ref(null)
const reviewType = ref('approve')
const submitting = ref(false)

// 表单数据
const reviewForm = ref({
  comments: ''
})

const submissionReviewForm = ref({
  score: 0,
  comments: '',
  suggestions: ''
})

// 计算属性
const filteredRegistrations = computed(() => {
  let result = registrations.value

  // 状态筛选
  if (filters.value.status) {
    result = result.filter(r => r.registrationStatus === filters.value.status)
  }

  // 级别筛选
  if (filters.value.competitionLevel) {
    result = result.filter(r => r.competitionLevel === filters.value.competitionLevel)
  }

  // 学生筛选
  if (filters.value.student) {
    result = result.filter(r => r.studentNumber === filters.value.student)
  }

  // 作品提交筛选
  if (filters.value.hasSubmission !== '') {
    const hasSubmission = filters.value.hasSubmission === 'true'
    result = result.filter(r => 
      hasSubmission ? r.submissionStatus === 'submitted' : r.submissionStatus !== 'submitted'
    )
  }

  // 评审状态筛选
  if (filters.value.hasReview !== '') {
    const hasReview = filters.value.hasReview === 'true'
    result = result.filter(r => 
      hasReview ? r.reviewStatus === 'completed' : r.reviewStatus !== 'completed'
    )
  }

  // 搜索筛选
  if (filters.value.search) {
    const query = filters.value.search.toLowerCase()
    result = result.filter(r => 
      r.competitionName.toLowerCase().includes(query) ||
      r.studentName.toLowerCase().includes(query) ||
      r.studentNumber.includes(query)
    )
  }

  return result
})

const hasSelectedPending = computed(() => {
  return selectedRegistrations.value.some(r => r.registrationStatus === 'pending')
})

// 方法
const loadCompetitionRegistrations = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      size: pageSize.value,
      ...filters.value
    }
    
    const response = await teacherService.getCompetitionRegistrations(params)
    const validation = validateApiResponse(response)
    
    if (validation.isValid) {
      registrations.value = ensureArray(validation.data.registrations, [])
      totalRegistrations.value = validation.data.total || 0
      
      // 更新统计数据
      await loadStats()
      
      ElMessage.success('竞赛报名数据加载成功')
    } else {
      ElMessage.error(validation.error || '加载竞赛报名数据失败')
    }
  } catch (error) {
    console.error('加载竞赛报名数据失败:', error)
    ElMessage.error(error.message || '加载竞赛报名数据失败')
  } finally {
    loading.value = false
  }
}

const loadStats = async () => {
  try {
    const response = await teacherService.getCompetitionGuidanceStats()
    const validation = validateApiResponse(response)
    
    if (validation.isValid) {
      stats.value = validation.data
    }
  } catch (error) {
    console.error('加载统计数据失败:', error)
    // 使用本地计算统计数据作为备选
    const data = registrations.value
    stats.value = {
      pendingCount: data.filter(r => r.registrationStatus === 'pending').length,
      reviewingCount: data.filter(r => r.reviewStatus === 'reviewing').length,
      completedCount: data.filter(r => r.reviewStatus === 'completed').length,
      totalCount: data.length
    }
  }
}

const loadStudentList = async () => {
  try {
    const response = await teacherService.getMyStudents()
    const validation = validateApiResponse(response)
    
    if (validation.isValid) {
      studentList.value = ensureArray(validation.data, [])
    }
  } catch (error) {
    console.error('加载学生列表失败:', error)
  }
}

const handleSearch = () => {
  currentPage.value = 1
  loadCompetitionRegistrations()
}

const resetFilters = () => {
  filters.value = {
    search: '',
    status: '',
    competitionLevel: '',
    student: '',
    dateRange: [],
    hasSubmission: '',
    hasReview: ''
  }
  currentPage.value = 1
  loadCompetitionRegistrations()
}

const handleSelectionChange = (selection) => {
  selectedRegistrations.value = selection
}

const handleRowClick = (row) => {
  viewRegistrationDetail(row)
}

const viewRegistrationDetail = async (registration) => {
  try {
    const response = await teacherService.getCompetitionRegistrationDetail(registration.id)
    const validation = validateApiResponse(response)
    
    if (validation.isValid) {
      currentRegistration.value = validation.data
      detailVisible.value = true
    } else {
      ElMessage.error(validation.error || '获取报名详情失败')
    }
  } catch (error) {
    console.error('获取报名详情失败:', error)
    ElMessage.error(error.message || '获取报名详情失败')
  }
}

const reviewRegistration = (registration, type) => {
  currentRegistration.value = registration
  reviewType.value = type
  reviewForm.value.comments = ''
  reviewVisible.value = true
}

const submitReview = async () => {
  if (reviewType.value === 'reject' && !reviewForm.value.comments.trim()) {
    ElMessage.warning('请输入驳回原因')
    return
  }

  submitting.value = true
  try {
    const reviewData = {
      status: reviewType.value === 'approve' ? 'approved' : 'rejected',
      comments: reviewForm.value.comments
    }
    
    const response = await teacherService.reviewCompetitionRegistration(
      currentRegistration.value.id, 
      reviewData
    )
    const validation = validateApiResponse(response)
    
    if (validation.isValid) {
      ElMessage.success(`审核${reviewType.value === 'approve' ? '通过' : '驳回'}成功`)
      reviewVisible.value = false
      
      // 重新加载数据
      await loadCompetitionRegistrations()
    } else {
      ElMessage.error(validation.error || '审核操作失败')
    }
  } catch (error) {
    console.error('审核操作失败:', error)
    ElMessage.error(error.message || '审核操作失败')
  } finally {
    submitting.value = false
  }
}

const reviewSubmission = async (registration) => {
  try {
    const response = await teacherService.getCompetitionSubmission(registration.submissionId)
    const validation = validateApiResponse(response)
    
    if (validation.isValid) {
      currentSubmission.value = validation.data
      submissionReviewForm.value = {
        score: 0,
        comments: '',
        suggestions: ''
      }
      submissionReviewVisible.value = true
    } else {
      ElMessage.error(validation.error || '获取作品详情失败')
    }
  } catch (error) {
    console.error('获取作品详情失败:', error)
    ElMessage.error(error.message || '获取作品详情失败')
  }
}

const submitSubmissionReview = async () => {
  if (!submissionReviewForm.value.score || !submissionReviewForm.value.comments.trim()) {
    ElMessage.warning('请填写完整的评审信息')
    return
  }

  submitting.value = true
  try {
    const reviewData = {
      score: submissionReviewForm.value.score,
      comments: submissionReviewForm.value.comments,
      suggestions: submissionReviewForm.value.suggestions
    }
    
    const response = await teacherService.submitSubmissionReview(
      currentSubmission.value.id,
      reviewData
    )
    const validation = validateApiResponse(response)
    
    if (validation.isValid) {
      ElMessage.success('评审意见提交成功')
      submissionReviewVisible.value = false
      
      // 重新加载数据
      await loadCompetitionRegistrations()
    } else {
      ElMessage.error(validation.error || '评审意见提交失败')
    }
  } catch (error) {
    console.error('评审意见提交失败:', error)
    ElMessage.error(error.message || '评审意见提交失败')
  } finally {
    submitting.value = false
  }
}

const viewResult = async (registration) => {
  try {
    const response = await teacherService.getCompetitionResult(registration.id)
    const validation = validateApiResponse(response)
    
    if (validation.isValid) {
      currentResult.value = validation.data
      resultVisible.value = true
    } else {
      ElMessage.error(validation.error || '获取竞赛成绩失败')
    }
  } catch (error) {
    console.error('获取竞赛成绩失败:', error)
    ElMessage.error(error.message || '获取竞赛成绩失败')
  }
}

const handleCommand = (command, row) => {
  switch (command) {
    case 'view-files':
      viewRegistrationDetail(row)
      break
    case 'view-review-records':
      viewRegistrationDetail(row)
      break
    case 'download-certificate':
      downloadCertificate(row.id)
      break
    case 'guidance-record':
      ElMessage.info('查看指导记录功能')
      break
  }
}

const batchApprove = async () => {
  if (selectedRegistrations.value.length === 0) {
    ElMessage.warning('请选择要审核的报名')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要批量通过 ${selectedRegistrations.value.length} 个报名吗？`,
      '批量审核',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const registrationIds = selectedRegistrations.value
      .filter(r => r.registrationStatus === 'pending')
      .map(r => r.id)
    
    if (registrationIds.length === 0) {
      ElMessage.warning('没有可审核的报名')
      return
    }

    const reviewData = {
      status: 'approved',
      comments: '批量审核通过'
    }
    
    const response = await teacherService.batchReviewCompetitionRegistrations(registrationIds, reviewData)
    const validation = validateApiResponse(response)
    
    if (validation.isValid) {
      ElMessage.success('批量审核成功')
      await loadCompetitionRegistrations()
    } else {
      ElMessage.error(validation.error || '批量审核失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('批量审核失败:', error)
      ElMessage.error(error.message || '批量审核失败')
    }
  }
}

const batchReject = async () => {
  if (selectedRegistrations.value.length === 0) {
    ElMessage.warning('请选择要审核的报名')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要批量驳回 ${selectedRegistrations.value.length} 个报名吗？`,
      '批量审核',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'danger'
      }
    )

    const registrationIds = selectedRegistrations.value
      .filter(r => r.registrationStatus === 'pending')
      .map(r => r.id)
    
    if (registrationIds.length === 0) {
      ElMessage.warning('没有可审核的报名')
      return
    }

    const reviewData = {
      status: 'rejected',
      comments: '批量审核驳回'
    }
    
    const response = await teacherService.batchReviewCompetitionRegistrations(registrationIds, reviewData)
    const validation = validateApiResponse(response)
    
    if (validation.isValid) {
      ElMessage.success('批量审核成功')
      await loadCompetitionRegistrations()
    } else {
      ElMessage.error(validation.error || '批量审核失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('批量审核失败:', error)
      ElMessage.error(error.message || '批量审核失败')
    }
  }
}

const downloadFile = async (file) => {
  try {
    await fileService.downloadFile(file.fileUrl, file.fileName)
    ElMessage.success('文件下载成功')
  } catch (error) {
    console.error('文件下载失败:', error)
    ElMessage.error(error.message || '文件下载失败')
  }
}

const previewFile = async (file) => {
  try {
    await fileService.previewFile(file.fileUrl)
  } catch (error) {
    console.error('文件预览失败:', error)
    ElMessage.error(error.message || '文件预览失败')
  }
}

const downloadCertificate = async (registrationId) => {
  try {
    const response = await teacherService.downloadCompetitionCertificate(registrationId)
    
    // 创建下载链接
    const blob = new Blob([response])
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `竞赛证书_${registrationId}.pdf`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    
    ElMessage.success('证书下载成功')
  } catch (error) {
    console.error('证书下载失败:', error)
    ElMessage.error(error.message || '证书下载失败')
  }
}

const downloadTranscript = async () => {
  if (!currentResult.value) return
  
  try {
    const response = await teacherService.downloadCompetitionTranscript(currentResult.value.registrationId)
    
    // 创建下载链接
    const blob = new Blob([response])
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `竞赛成绩单_${currentResult.value.studentName}.pdf`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    
    ElMessage.success('成绩单下载成功')
  } catch (error) {
    console.error('成绩单下载失败:', error)
    ElMessage.error(error.message || '成绩单下载失败')
  }
}

const exportGuidanceData = async () => {
  try {
    const response = await teacherService.exportCompetitionGuidanceData(filters.value)
    
    // 创建下载链接
    const blob = new Blob([response])
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `竞赛指导数据_${new Date().toISOString().split('T')[0]}.xlsx`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    
    ElMessage.success('数据导出成功')
  } catch (error) {
    console.error('数据导出失败:', error)
    ElMessage.error(error.message || '数据导出失败')
  }
}

const handleSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
  loadCompetitionRegistrations()
}

const handleCurrentChange = (page) => {
  currentPage.value = page
  loadCompetitionRegistrations()
}

// 状态相关方法
const getStatusType = (status) => {
  const statusMap = {
    'pending': 'warning',
    'approved': 'success',
    'rejected': 'danger'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status) => {
  const statusMap = {
    'pending': '待审核',
    'approved': '已通过',
    'rejected': '已驳回'
  }
  return statusMap[status] || status
}

const getSubmissionStatusType = (status) => {
  const statusMap = {
    'not_submitted': 'info',
    'submitted': 'success',
    'overdue': 'danger'
  }
  return statusMap[status] || 'info'
}

const getSubmissionStatusText = (status) => {
  const statusMap = {
    'not_submitted': '未提交',
    'submitted': '已提交',
    'overdue': '已逾期'
  }
  return statusMap[status] || status
}

const getReviewStatusType = (status) => {
  const statusMap = {
    'not_started': 'info',
    'reviewing': 'warning',
    'completed': 'success'
  }
  return statusMap[status] || 'info'
}

const getReviewStatusText = (status) => {
  const statusMap = {
    'not_started': '未开始',
    'reviewing': '评审中',
    'completed': '已完成'
  }
  return statusMap[status] || status
}

const getAwardType = (award) => {
  const awardMap = {
    '一等奖': 'success',
    '二等奖': 'warning',
    '三等奖': 'info',
    '优秀奖': 'info'
  }
  return awardMap[award] || 'info'
}

// 组件挂载时加载数据
onMounted(() => {
  loadStudentList()
  loadCompetitionRegistrations()
})
</script>

<style scoped>
.competition-guidance {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header-left h2 {
  margin: 0 0 8px 0;
  color: #2c3e50;
  font-size: 24px;
  font-weight: 600;
}

.header-desc {
  margin: 0;
  color: #7f8c8d;
  font-size: 14px;
}

.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  height: 100px;
}

.stat-content {
  display: flex;
  align-items: center;
  height: 100%;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 15px;
  font-size: 24px;
  color: white;
}

.stat-icon.pending {
  background: linear-gradient(135deg, #ff9a56, #ff6b6b);
}

.stat-icon.reviewing {
  background: linear-gradient(135deg, #4ecdc4, #44a08d);
}

.stat-icon.completed {
  background: linear-gradient(135deg, #667eea, #764ba2);
}

.stat-icon.total {
  background: linear-gradient(135deg, #f093fb, #f5576c);
}

.stat-info h4 {
  margin: 0 0 5px 0;
  color: #2c3e50;
  font-size: 14px;
  font-weight: 500;
}

.stat-number {
  margin: 0;
  color: #2c3e50;
  font-size: 24px;
  font-weight: 600;
}

.filter-card {
  margin-bottom: 20px;
}

.filter-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.advanced-filter {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #ebeef5;
}

.registration-list {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.pagination-wrapper {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

.registration-detail,
.submission-review,
.competition-result {
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