<template>
  <div class="competition-judging">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2>竞赛作品评审</h2>
      <p>评审学生提交的竞赛作品并给出评分</p>
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
              <h4>待评审</h4>
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
              <h4>评审中</h4>
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
              <h4>已完成</h4>
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
              <h4>总作品</h4>
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
            placeholder="搜索作品标题、学生姓名"
            clearable
            @input="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-col>
        <el-col :span="4">
          <el-select v-model="filters.status" placeholder="评审状态" clearable @change="handleSearch">
            <el-option label="全部" value="" />
            <el-option label="待评审" value="pending" />
            <el-option label="评审中" value="reviewing" />
            <el-option label="已完成" value="completed" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-select v-model="filters.competitionId" placeholder="选择竞赛" clearable @change="handleSearch">
            <el-option label="全部竞赛" value="" />
            <el-option 
              v-for="comp in competitions" 
              :key="comp.id" 
              :label="comp.title" 
              :value="comp.id"
            />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-select v-model="filters.department" placeholder="选择学院" clearable @change="handleSearch">
            <el-option label="全部学院" value="" />
            <el-option label="计算机学院" value="computer" />
            <el-option label="数学学院" value="mathematics" />
            <el-option label="物理学院" value="physics" />
            <el-option label="化学学院" value="chemistry" />
            <el-option label="工程学院" value="engineering" />
            <el-option label="商学院" value="business" />
          </el-select>
        </el-col>
        <el-col :span="6">
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="exportJudgingData">
            <el-icon><Download /></el-icon>
            导出评审数据
          </el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 作品列表 -->
    <el-card class="submission-list-card">
      <el-table 
        :data="submissions" 
        style="width: 100%" 
        v-loading="loading"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="作品ID" width="80" />
        <el-table-column label="竞赛信息" width="200">
          <template #default="scope">
            <div class="competition-info">
              <h4>{{ scope.row.competition?.title }}</h4>
              <p class="competition-type">{{ scope.row.competition?.type }}</p>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="学生信息" width="180">
          <template #default="scope">
            <div class="student-info">
              <h4>{{ scope.row.student?.realName }}</h4>
              <p>{{ scope.row.student?.department }} - {{ scope.row.student?.studentId }}</p>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="title" label="作品标题" width="200" />
        <el-table-column prop="submitTime" label="提交时间" width="150">
          <template #default="scope">
            {{ formatDate(scope.row.submitTime) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="评审状态" width="100">
          <template #default="scope">
            <el-tag :type="getJudgingStatusType(scope.row.status)">
              {{ getJudgingStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="score" label="当前评分" width="100">
          <template #default="scope">
            <span v-if="scope.row.score" class="score">{{ scope.row.score }}/100</span>
            <span v-else class="no-score">未评分</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="scope">
            <div class="action-buttons">
              <el-button 
                size="small" 
                type="primary" 
                @click="viewSubmission(scope.row)"
              >
                查看
              </el-button>
              <el-button 
                v-if="scope.row.status !== 'completed'"
                size="small" 
                type="success" 
                @click="startJudging(scope.row)"
              >
                开始评审
              </el-button>
              <el-button 
                v-if="scope.row.status === 'reviewing'"
                size="small" 
                type="warning" 
                @click="continueJudging(scope.row)"
              >
                继续评审
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 评审表单对话框 -->
    <el-dialog
      v-model="showJudgingDialog"
      title="作品评审"
      width="80%"
      class="judging-dialog"
    >
      <div v-if="currentJudgingSubmission" class="judging-content">
        <!-- 作品基本信息 -->
        <el-card class="submission-info-card">
          <template #header>
            <h3>作品基本信息</h3>
          </template>
          
          <el-descriptions :column="3" border>
            <el-descriptions-item label="作品标题">
              {{ currentJudgingSubmission.title }}
            </el-descriptions-item>
            <el-descriptions-item label="学生姓名">
              {{ currentJudgingSubmission.student?.realName }}
            </el-descriptions-item>
            <el-descriptions-item label="提交时间">
              {{ formatDate(currentJudgingSubmission.submitTime) }}
            </el-descriptions-item>
            <el-descriptions-item label="竞赛名称">
              {{ currentJudgingSubmission.competition?.title }}
            </el-descriptions-item>
            <el-descriptions-item label="作品版本">
              {{ currentJudgingSubmission.version }}
            </el-descriptions-item>
            <el-descriptions-item label="文件大小">
              {{ formatFileSize(currentJudgingSubmission.fileSize) }}
            </el-descriptions-item>
          </el-descriptions>
          
          <div class="description-section">
            <h4>作品描述</h4>
            <p>{{ currentJudgingSubmission.description }}</p>
          </div>
          
          <div class="file-section">
            <h4>作品文件</h4>
            <el-button 
              type="primary" 
              @click="downloadFile(currentJudgingSubmission)"
              :icon="Download"
            >
              下载作品文件
            </el-button>
          </div>
        </el-card>

        <!-- 评审表单 -->
        <el-card class="judging-form-card">
          <template #header>
            <h3>评审评分</h3>
          </template>
          
          <el-form
            ref="judgingFormRef"
            :model="judgingForm"
            :rules="judgingRules"
            label-width="120px"
          >
            <!-- 多维度评分 -->
            <div class="scoring-section">
              <h4>多维度评分 (总分100分)</h4>
              
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="创新性" prop="innovation">
                    <el-input-number
                      v-model="judgingForm.innovation"
                      :min="0"
                      :max="25"
                      :step="1"
                      placeholder="0-25分"
                    />
                    <span class="score-hint">(25分)</span>
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="技术性" prop="technical">
                    <el-input-number
                      v-model="judgingForm.technical"
                      :min="0"
                      :max="25"
                      :step="1"
                      placeholder="0-25分"
                    />
                    <span class="score-hint">(25分)</span>
                  </el-form-item>
                </el-col>
              </el-row>
              
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="完整性" prop="completeness">
                    <el-input-number
                      v-model="judgingForm.completeness"
                      :min="0"
                      :max="20"
                      :step="1"
                      placeholder="0-20分"
                    />
                    <span class="score-hint">(20分)</span>
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="实用性" prop="practicality">
                    <el-input-number
                      v-model="judgingForm.practicality"
                      :min="0"
                      :max="20"
                      :step="1"
                      placeholder="0-20分"
                    />
                    <span class="score-hint">(20分)</span>
                  </el-form-item>
                </el-col>
              </el-row>
              
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="文档质量" prop="documentation">
                    <el-input-number
                      v-model="judgingForm.documentation"
                      :min="0"
                      :max="10"
                      :step="1"
                      placeholder="0-10分"
                    />
                    <span class="score-hint">(10分)</span>
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="总分" prop="totalScore">
                    <el-input
                      v-model="judgingForm.totalScore"
                      disabled
                      class="total-score-input"
                    />
                    <span class="score-hint">(自动计算)</span>
                  </el-form-item>
                </el-col>
              </el-row>
            </div>
            
            <!-- 评语 -->
            <el-form-item label="评审评语" prop="comment">
              <el-input
                v-model="judgingForm.comment"
                type="textarea"
                :rows="6"
                placeholder="请详细说明评分理由，指出作品的优点和不足，提出改进建议..."
                maxlength="1000"
                show-word-limit
              />
            </el-form-item>
            
            <!-- 评审建议 -->
            <el-form-item label="评审建议" prop="suggestion">
              <el-select v-model="judgingForm.suggestion" placeholder="请选择评审建议">
                <el-option label="推荐获奖" value="recommend_award" />
                <el-option label="建议修改后重新提交" value="suggest_revision" />
                <el-option label="基本合格" value="basic_qualified" />
                <el-option label="需要重大改进" value="need_major_improvement" />
              </el-select>
            </el-form-item>
            
            <!-- 操作按钮 -->
            <el-form-item>
              <el-button 
                type="primary" 
                @click="submitJudging"
                :loading="submitting"
              >
                提交评审
              </el-button>
              <el-button @click="saveAsDraft">保存草稿</el-button>
              <el-button @click="showJudgingDialog = false">取消</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </div>
    </el-dialog>

    <!-- 作品详情对话框 -->
    <el-dialog
      v-model="showDetailDialog"
      title="作品详情"
      width="70%"
    >
      <div v-if="selectedSubmission" class="submission-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="作品标题">
            {{ selectedSubmission.title }}
          </el-descriptions-item>
          <el-descriptions-item label="学生姓名">
            {{ selectedSubmission.student?.realName }}
          </el-descriptions-item>
          <el-descriptions-item label="竞赛名称">
            {{ selectedSubmission.competition?.title }}
          </el-descriptions-item>
          <el-descriptions-item label="提交时间">
            {{ formatDate(selectedSubmission.submitTime) }}
          </el-descriptions-item>
          <el-descriptions-item label="文件大小">
            {{ formatFileSize(selectedSubmission.fileSize) }}
          </el-descriptions-item>
          <el-descriptions-item label="评审状态">
            <el-tag :type="getJudgingStatusType(selectedSubmission.status)">
              {{ getJudgingStatusText(selectedSubmission.status) }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>
        
        <div class="description-section">
          <h4>作品描述</h4>
          <p>{{ selectedSubmission.description }}</p>
        </div>
        
        <div class="file-section">
          <h4>作品文件</h4>
          <el-button 
            type="primary" 
            @click="downloadFile(selectedSubmission)"
            :icon="Download"
          >
            下载作品文件
          </el-button>
        </div>
        
        <div v-if="selectedSubmission.score" class="score-section">
          <h4>评审结果</h4>
          <el-descriptions :column="3" border>
            <el-descriptions-item label="总分">
              <span class="total-score">{{ selectedSubmission.score }}/100</span>
            </el-descriptions-item>
            <el-descriptions-item label="评审教师">
              {{ selectedSubmission.judge?.realName }}
            </el-descriptions-item>
            <el-descriptions-item label="评审时间">
              {{ formatDate(selectedSubmission.judgeTime) }}
            </el-descriptions-item>
          </el-descriptions>
          
          <div v-if="selectedSubmission.comment" class="comment-section">
            <h5>评审评语</h5>
            <p>{{ selectedSubmission.comment }}</p>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Clock, Document, Check, DataAnalysis, Search, Download 
} from '@element-plus/icons-vue'
import { formatDate } from '@/utils/dateUtils'

// 响应式数据
const loading = ref(false)
const submissions = ref([])
const competitions = ref([])
const selectedSubmission = ref(null)
const currentJudgingSubmission = ref(null)
const showDetailDialog = ref(false)
const showJudgingDialog = ref(false)
const submitting = ref(false)
const selectedSubmissions = ref([])

// 分页相关
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

// 筛选条件
const filters = ref({
  search: '',
  status: '',
  competitionId: '',
  department: ''
})

// 评审表单
const judgingForm = ref({
  innovation: 0,
  technical: 0,
  completeness: 0,
  practicality: 0,
  documentation: 0,
  totalScore: 0,
  comment: '',
  suggestion: ''
})

// 表单引用
const judgingFormRef = ref()

// 表单验证规则
const judgingRules = {
  innovation: [
    { required: true, message: '请评分创新性', trigger: 'blur' }
  ],
  technical: [
    { required: true, message: '请评分技术性', trigger: 'blur' }
  ],
  completeness: [
    { required: true, message: '请评分完整性', trigger: 'blur' }
  ],
  practicality: [
    { required: true, message: '请评分实用性', trigger: 'blur' }
  ],
  documentation: [
    { required: true, message: '请评分文档质量', trigger: 'blur' }
  ],
  comment: [
    { required: true, message: '请填写评审评语', trigger: 'blur' },
    { min: 10, max: 1000, message: '评语长度在 10 到 1000 个字符', trigger: 'blur' }
  ],
  suggestion: [
    { required: true, message: '请选择评审建议', trigger: 'blur' }
  ]
}

// 统计数据
const stats = ref({
  pendingCount: 0,
  reviewingCount: 0,
  completedCount: 0,
  totalCount: 0
})

// 计算属性
const filteredSubmissions = computed(() => {
  let filtered = submissions.value

  if (filters.value.search) {
    const search = filters.value.search.toLowerCase()
    filtered = filtered.filter(s => 
      s.title?.toLowerCase().includes(search) ||
      s.student?.realName?.toLowerCase().includes(search) ||
      s.competition?.title?.toLowerCase().includes(search)
    )
  }

  if (filters.value.status) {
    filtered = filtered.filter(s => s.status === filters.value.status)
  }

  if (filters.value.competitionId) {
    filtered = filtered.filter(s => s.competitionId === filters.value.competitionId)
  }

  if (filters.value.department) {
    filtered = filtered.filter(s => s.student?.department === filters.value.department)
  }

  return filtered
})

// 分页后的数据
const paginatedSubmissions = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredSubmissions.value.slice(start, end)
})

// 监听评分变化，自动计算总分
watch(
  () => [
    judgingForm.value.innovation,
    judgingForm.value.technical,
    judgingForm.value.completeness,
    judgingForm.value.practicality,
    judgingForm.value.documentation
  ],
  (newScores) => {
    const total = newScores.reduce((sum, score) => sum + (score || 0), 0)
    judgingForm.value.totalScore = total
  },
  { immediate: true }
)

// 方法
const loadSubmissions = async () => {
  loading.value = true
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // 模拟数据
    submissions.value = [
      {
        id: 1,
        competitionId: 1,
        competition: {
          id: 1,
          title: '全国大学生程序设计竞赛',
          type: '程序设计'
        },
        student: {
          id: 1,
          realName: '张三',
          studentId: '2021001',
          department: '计算机学院'
        },
        title: '基于深度学习的图像识别系统',
        description: '使用Python和TensorFlow实现的图像识别系统，能够识别多种物体类别',
        submitTime: new Date('2024-02-15'),
        status: 'pending',
        score: null,
        fileSize: 1024 * 1024 * 15, // 15MB
        version: '1.0'
      },
      {
        id: 2,
        competitionId: 1,
        competition: {
          id: 1,
          title: '全国大学生程序设计竞赛',
          type: '程序设计'
        },
        student: {
          id: 2,
          realName: '李四',
          studentId: '2021002',
          department: '计算机学院'
        },
        title: '智能聊天机器人',
        description: '基于自然语言处理技术的智能聊天机器人',
        submitTime: new Date('2024-02-16'),
        status: 'reviewing',
        score: 85,
        fileSize: 1024 * 1024 * 8, // 8MB
        version: '1.0',
        judge: { realName: '王教授' },
        judgeTime: new Date('2024-02-18'),
        comment: '作品技术实现较好，但创新性有待提高'
      }
    ]
    
    loadStats()
  } catch (error) {
    console.error('加载作品数据失败:', error)
    ElMessage.error('加载作品数据失败')
  } finally {
    loading.value = false
  }
}

const loadCompetitions = async () => {
  try {
    // 模拟加载竞赛列表
    competitions.value = [
      { id: 1, title: '全国大学生程序设计竞赛' },
      { id: 2, title: '数学建模竞赛' },
      { id: 3, title: '创新创业大赛' }
    ]
  } catch (error) {
    console.error('加载竞赛列表失败:', error)
  }
}

const loadStats = () => {
  const pending = submissions.value.filter(s => s.status === 'pending').length
  const reviewing = submissions.value.filter(s => s.status === 'reviewing').length
  const completed = submissions.value.filter(s => s.status === 'completed').length
  
  stats.value = {
    pendingCount: pending,
    reviewingCount: reviewing,
    completedCount: completed,
    totalCount: submissions.value.length
  }
}

const handleSearch = () => {
  currentPage.value = 1
  // 实际项目中这里会调用API
}

const resetFilters = () => {
  filters.value = {
    search: '',
    status: '',
    competitionId: '',
    department: ''
  }
  handleSearch()
}

const handleSizeChange = (val) => {
  pageSize.value = val
  currentPage.value = 1
}

const handleCurrentChange = (val) => {
  currentPage.value = val
}

const handleSelectionChange = (selection) => {
  selectedSubmissions.value = selection
}

const viewSubmission = (submission) => {
  selectedSubmission.value = submission
  showDetailDialog.value = true
}

const startJudging = (submission) => {
  currentJudgingSubmission.value = submission
  resetJudgingForm()
  showJudgingDialog.value = true
}

const continueJudging = (submission) => {
  currentJudgingSubmission.value = submission
  // 加载已有的评审数据
  loadExistingJudging(submission)
  showJudgingDialog.value = true
}

const resetJudgingForm = () => {
  judgingForm.value = {
    innovation: 0,
    technical: 0,
    completeness: 0,
    practicality: 0,
    documentation: 0,
    totalScore: 0,
    comment: '',
    suggestion: ''
  }
}

const loadExistingJudging = (submission) => {
  // 模拟加载已有的评审数据
  if (submission.score) {
    judgingForm.value = {
      innovation: 20,
      technical: 22,
      completeness: 18,
      practicality: 16,
      documentation: 9,
      totalScore: submission.score,
      comment: submission.comment || '',
      suggestion: 'recommend_award'
    }
  }
}

const submitJudging = async () => {
  try {
    await judgingFormRef.value.validate()
    
    submitting.value = true
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // 更新作品状态和评分
    const index = submissions.value.findIndex(s => s.id === currentJudgingSubmission.value.id)
    if (index !== -1) {
      submissions.value[index].status = 'completed'
      submissions.value[index].score = judgingForm.value.totalScore
      submissions.value[index].comment = judgingForm.value.comment
      submissions.value[index].judge = { realName: '当前教师' }
      submissions.value[index].judgeTime = new Date()
    }
    
    loadStats()
    showJudgingDialog.value = false
    ElMessage.success('评审提交成功！')
  } catch (error) {
    ElMessage.error('提交失败：' + (error.message || '未知错误'))
  } finally {
    submitting.value = false
  }
}

const saveAsDraft = () => {
  ElMessage.info('草稿保存功能开发中...')
}

const downloadFile = (submission) => {
  ElMessage.success('开始下载文件...')
  // 实际项目中这里会调用下载API
}

const exportJudgingData = () => {
  ElMessage.success('开始导出评审数据...')
  // 实际项目中这里会调用导出API
}

// 状态相关方法
const getJudgingStatusType = (status) => {
  const statusMap = {
    pending: 'warning',
    reviewing: 'primary',
    completed: 'success'
  }
  return statusMap[status] || 'info'
}

const getJudgingStatusText = (status) => {
  const statusMap = {
    pending: '待评审',
    reviewing: '评审中',
    completed: '已完成'
  }
  return statusMap[status] || status
}

const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 组件挂载时加载数据
onMounted(() => {
  loadSubmissions()
  loadCompetitions()
})
</script>

<style scoped>
.competition-judging {
  padding: 20px;
}

.page-header {
  margin-bottom: 30px;
  text-align: center;
}

.page-header h2 {
  margin: 0 0 10px 0;
  color: #2c3e50;
  font-size: 28px;
  font-weight: 600;
}

.page-header p {
  margin: 0;
  color: #7f8c8d;
  font-size: 16px;
}

.stats-row {
  margin-bottom: 30px;
}

.stat-card {
  height: 120px;
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
  margin-right: 20px;
  font-size: 24px;
  color: white;
}

.stat-icon.pending {
  background: #e6a23c;
}

.stat-icon.reviewing {
  background: #409eff;
}

.stat-icon.completed {
  background: #67c23a;
}

.stat-icon.total {
  background: #909399;
}

.stat-info h4 {
  margin: 0 0 10px 0;
  color: #606266;
  font-size: 14px;
}

.stat-number {
  margin: 0;
  color: #2c3e50;
  font-size: 24px;
  font-weight: 600;
}

.filter-card {
  margin-bottom: 30px;
}

.filter-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.submission-list-card {
  margin-bottom: 30px;
}

.competition-info h4 {
  margin: 0 0 5px 0;
  color: #2c3e50;
  font-size: 14px;
}

.competition-type {
  margin: 0;
  color: #7f8c8d;
  font-size: 12px;
}

.student-info h4 {
  margin: 0 0 5px 0;
  color: #2c3e50;
  font-size: 14px;
}

.student-info p {
  margin: 0;
  color: #7f8c8d;
  font-size: 12px;
}

.score {
  color: #67c23a;
  font-weight: 600;
}

.no-score {
  color: #c0c4cc;
  font-style: italic;
}

.action-buttons {
  display: flex;
  gap: 5px;
}

.pagination-container {
  margin-top: 20px;
  text-align: center;
}

/* 评审对话框样式 */
.judging-dialog .el-dialog__body {
  max-height: 70vh;
  overflow-y: auto;
}

.submission-info-card,
.judging-form-card {
  margin-bottom: 20px;
}

.submission-info-card h3,
.judging-form-card h3 {
  margin: 0;
  color: #2c3e50;
}

.description-section,
.file-section,
.score-section,
.comment-section {
  margin-top: 20px;
}

.description-section h4,
.file-section h4,
.score-section h4,
.comment-section h5 {
  margin: 0 0 10px 0;
  color: #2c3e50;
  font-size: 16px;
}

.description-section p,
.comment-section p {
  margin: 0;
  color: #606266;
  line-height: 1.6;
  background: #f8f9fa;
  padding: 15px;
  border-radius: 6px;
}

.scoring-section {
  margin-bottom: 30px;
}

.scoring-section h4 {
  margin: 0 0 20px 0;
  color: #2c3e50;
  font-size: 18px;
  border-bottom: 2px solid #409eff;
  padding-bottom: 10px;
}

.score-hint {
  margin-left: 10px;
  color: #909399;
  font-size: 12px;
}

.total-score-input {
  background: #f0f9ff;
  border-color: #409eff;
}

.total-score {
  color: #67c23a;
  font-weight: 600;
  font-size: 16px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .stats-row .el-col {
    margin-bottom: 15px;
  }
  
  .filter-card .el-col {
    margin-bottom: 10px;
  }
  
  .action-buttons {
    flex-direction: column;
  }
  
  .scoring-section .el-col {
    margin-bottom: 15px;
  }
}
</style> 