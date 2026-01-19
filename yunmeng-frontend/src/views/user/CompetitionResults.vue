<template>
  <div class="competition-results">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2>竞赛结果管理</h2>
      <p>管理竞赛获奖结果、发布成绩和生成证书</p>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon total">
              <el-icon><DataAnalysis /></el-icon>
            </div>
            <div class="stat-info">
              <h4>总竞赛数</h4>
              <p class="stat-number">{{ stats.totalCompetitions }}</p>
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
              <p class="stat-number">{{ stats.completedCompetitions }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon pending">
              <el-icon><Clock /></el-icon>
            </div>
            <div class="stat-info">
              <h4>待发布</h4>
              <p class="stat-number">{{ stats.pendingResults }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon published">
              <el-icon><Trophy /></el-icon>
            </div>
            <div class="stat-info">
              <h4>已发布</h4>
              <p class="stat-number">{{ stats.publishedResults }}</p>
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
            placeholder="搜索竞赛名称"
            clearable
            @input="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-col>
        <el-col :span="4">
          <el-select v-model="filters.status" placeholder="竞赛状态" clearable @change="handleSearch">
            <el-option label="全部" value="" />
            <el-option label="评审中" value="review" />
            <el-option label="已完成" value="completed" />
            <el-option label="结果已发布" value="published" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-select v-model="filters.level" placeholder="竞赛级别" clearable @change="handleSearch">
            <el-option label="全部级别" value="" />
            <el-option label="校级" value="school" />
            <el-option label="省级" value="provincial" />
            <el-option label="国家级" value="national" />
            <el-option label="国际级" value="international" />
          </el-select>
        </el-col>
        <el-col :span="4">
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
        <el-col :span="6">
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="exportResults">
            <el-icon><Download /></el-icon>
            导出结果
          </el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 竞赛列表 -->
    <el-card class="competition-list-card">
      <el-table 
        :data="competitions" 
        style="width: 100%" 
        v-loading="loading"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="竞赛ID" width="80" />
        <el-table-column label="竞赛信息" width="300">
          <template #default="scope">
            <div class="competition-info">
              <h4>{{ scope.row.title }}</h4>
              <p class="competition-meta">
                <span class="type">{{ scope.row.type }}</span>
                <span class="level">{{ scope.row.level }}</span>
                <span class="organizer">{{ scope.row.organizer }}</span>
              </p>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="participantCount" label="参赛人数" width="100" />
        <el-table-column prop="status" label="竞赛状态" width="100">
          <template #default="scope">
            <el-tag :type="getCompetitionStatusType(scope.row.status)">
              {{ getCompetitionStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="resultStatus" label="结果状态" width="120">
          <template #default="scope">
            <el-tag :type="getResultStatusType(scope.row.resultStatus)">
              {{ getResultStatusText(scope.row.resultStatus) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="completedAt" label="完成时间" width="150">
          <template #default="scope">
            {{ scope.row.completedAt ? formatDate(scope.row.completedAt) : '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="publishedAt" label="发布时间" width="150">
          <template #default="scope">
            {{ scope.row.publishedAt ? formatDate(scope.row.publishedAt) : '-' }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="scope">
            <div class="action-buttons">
              <el-button 
                size="small" 
                type="primary" 
                @click="viewCompetition(scope.row)"
              >
                查看
              </el-button>
              <el-button 
                v-if="scope.row.status === 'completed' && !scope.row.resultStatus"
                size="small" 
                type="success" 
                @click="manageResults(scope.row)"
              >
                管理结果
              </el-button>
              <el-button 
                v-if="scope.row.resultStatus === 'pending'"
                size="small" 
                type="warning" 
                @click="publishResults(scope.row)"
              >
                发布结果
              </el-button>
              <el-button 
                v-if="scope.row.resultStatus === 'published'"
                size="small" 
                type="info" 
                @click="viewResults(scope.row)"
              >
                查看结果
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

    <!-- 结果管理对话框 -->
    <el-dialog
      v-model="showResultsDialog"
      title="竞赛结果管理"
      width="80%"
      class="results-dialog"
    >
      <div v-if="currentCompetition" class="results-content">
        <!-- 竞赛基本信息 -->
        <el-card class="competition-info-card">
          <template #header>
            <h3>{{ currentCompetition.title }} - 结果管理</h3>
          </template>
          
          <el-descriptions :column="3" border>
            <el-descriptions-item label="竞赛类型">
              {{ currentCompetition.type }}
            </el-descriptions-item>
            <el-descriptions-item label="竞赛级别">
              {{ currentCompetition.level }}
            </el-descriptions-item>
            <el-descriptions-item label="主办方">
              {{ currentCompetition.organizer }}
            </el-descriptions-item>
            <el-descriptions-item label="参赛人数">
              {{ currentCompetition.participantCount }}
            </el-descriptions-item>
            <el-descriptions-item label="完成时间">
              {{ formatDate(currentCompetition.completedAt) }}
            </el-descriptions-item>
            <el-descriptions-item label="当前状态">
              <el-tag :type="getCompetitionStatusType(currentCompetition.status)">
                {{ getCompetitionStatusText(currentCompetition.status) }}
              </el-tag>
            </el-descriptions-item>
          </el-descriptions>
        </el-card>

        <!-- 获奖结果录入 -->
        <el-card class="results-entry-card">
          <template #header>
            <h3>获奖结果录入</h3>
          </template>
          
          <el-table :data="resultsList" style="width: 100%">
            <el-table-column prop="rank" label="排名" width="80" />
            <el-table-column label="学生信息" width="200">
              <template #default="scope">
                <div class="student-info">
                  <h4>{{ scope.row.student?.realName }}</h4>
                  <p>{{ scope.row.student?.department }} - {{ scope.row.student?.studentId }}</p>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="awardLevel" label="获奖等级" width="120">
              <template #default="scope">
                <el-select v-model="scope.row.awardLevel" placeholder="选择等级">
                  <el-option label="特等奖" value="special" />
                  <el-option label="一等奖" value="first" />
                  <el-option label="二等奖" value="second" />
                  <el-option label="三等奖" value="third" />
                  <el-option label="优秀奖" value="honorable" />
                </el-select>
              </template>
            </el-table-column>
            <el-table-column prop="finalScore" label="最终得分" width="120">
              <template #default="scope">
                <el-input-number
                  v-model="scope.row.finalScore"
                  :min="0"
                  :max="100"
                  :step="0.1"
                  :precision="1"
                  placeholder="0-100"
                />
              </template>
            </el-table-column>
            <el-table-column prop="certificateUrl" label="证书" width="150">
              <template #default="scope">
                <el-upload
                  :action="uploadAction"
                  :show-file-list="false"
                  :on-success="(res) => handleCertificateUpload(res, scope.row)"
                  :on-error="handleUploadError"
                  accept=".pdf,.jpg,.png"
                >
                  <el-button size="small" type="primary">
                    {{ scope.row.certificateUrl ? '重新上传' : '上传证书' }}
                  </el-button>
                </el-upload>
                <div v-if="scope.row.certificateUrl" class="certificate-info">
                  <el-link :href="scope.row.certificateUrl" target="_blank" type="primary">
                    查看证书
                  </el-link>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="remarks" label="备注" width="200">
              <template #default="scope">
                <el-input
                  v-model="scope.row.remarks"
                  placeholder="备注信息"
                  maxlength="100"
                  show-word-limit
                />
              </template>
            </el-table-column>
          </el-table>
          
          <div class="results-actions">
            <el-button type="primary" @click="saveResults" :loading="saving">
              保存结果
            </el-button>
            <el-button @click="generateCertificates">
              生成证书
            </el-button>
            <el-button @click="previewResults">
              预览结果
            </el-button>
          </div>
        </el-card>
      </div>
    </el-dialog>

    <!-- 结果预览对话框 -->
    <el-dialog
      v-model="showPreviewDialog"
      title="结果预览"
      width="70%"
    >
      <div v-if="previewData" class="preview-content">
        <h3>{{ previewData.competitionTitle }} - 获奖名单</h3>
        
        <el-table :data="previewData.results" style="width: 100%">
          <el-table-column prop="rank" label="排名" width="80" />
          <el-table-column label="学生信息" width="200">
            <template #default="scope">
              <div class="student-info">
                <h4>{{ scope.row.student.realName }}</h4>
                <p>{{ scope.row.student.department }} - {{ scope.row.student.studentId }}</p>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="awardLevel" label="获奖等级" width="120">
            <template #default="scope">
              <el-tag :type="getAwardLevelType(scope.row.awardLevel)">
                {{ getAwardLevelText(scope.row.awardLevel) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="finalScore" label="最终得分" width="100" />
          <el-table-column prop="remarks" label="备注" width="200" />
        </el-table>
        
        <div class="preview-summary">
          <h4>统计摘要</h4>
          <el-descriptions :column="3" border>
            <el-descriptions-item label="特等奖">
              {{ previewData.summary.special }}
            </el-descriptions-item>
            <el-descriptions-item label="一等奖">
              {{ previewData.summary.first }}
            </el-descriptions-item>
            <el-descriptions-item label="二等奖">
              {{ previewData.summary.second }}
            </el-descriptions-item>
            <el-descriptions-item label="三等奖">
              {{ previewData.summary.third }}
            </el-descriptions-item>
            <el-descriptions-item label="优秀奖">
              {{ previewData.summary.honorable }}
            </el-descriptions-item>
            <el-descriptions-item label="总获奖人数">
              {{ previewData.summary.total }}
            </el-descriptions-item>
          </el-descriptions>
        </div>
      </div>
    </el-dialog>

    <!-- 竞赛详情对话框 -->
    <el-dialog
      v-model="showDetailDialog"
      title="竞赛详情"
      width="70%"
    >
      <div v-if="selectedCompetition" class="competition-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="竞赛名称">
            {{ selectedCompetition.title }}
          </el-descriptions-item>
          <el-descriptions-item label="竞赛类型">
            {{ selectedCompetition.type }}
          </el-descriptions-item>
          <el-descriptions-item label="竞赛级别">
            {{ selectedCompetition.level }}
          </el-descriptions-item>
          <el-descriptions-item label="主办方">
            {{ selectedCompetition.organizer }}
          </el-descriptions-item>
          <el-descriptions-item label="参赛人数">
            {{ selectedCompetition.participantCount }}
          </el-descriptions-item>
          <el-descriptions-item label="完成时间">
            {{ formatDate(selectedCompetition.completedAt) }}
          </el-descriptions-item>
        </el-descriptions>
        
        <div v-if="selectedCompetition.results && selectedCompetition.results.length > 0" class="results-section">
          <h4>获奖结果</h4>
          <el-table :data="selectedCompetition.results" style="width: 100%">
            <el-table-column prop="rank" label="排名" width="80" />
            <el-table-column label="学生信息" width="200">
              <template #default="scope">
                <div class="student-info">
                  <h4>{{ scope.row.student.realName }}</h4>
                  <p>{{ scope.row.student.department }} - {{ scope.row.student.studentId }}</p>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="awardLevel" label="获奖等级" width="120">
              <template #default="scope">
                <el-tag :type="getAwardLevelType(scope.row.awardLevel)">
                  {{ getAwardLevelText(scope.row.awardLevel) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="finalScore" label="最终得分" width="100" />
          </el-table>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  DataAnalysis, Check, Clock, Trophy, Search, Download 
} from '@element-plus/icons-vue'
import { formatDate } from '@/utils/dateUtils'

// 响应式数据
const loading = ref(false)
const competitions = ref([])
const selectedCompetition = ref(null)
const currentCompetition = ref(null)
const showDetailDialog = ref(false)
const showResultsDialog = ref(false)
const showPreviewDialog = ref(false)
const saving = ref(false)
const selectedCompetitions = ref([])
const resultsList = ref([])
const previewData = ref(null)

// 分页相关
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

// 筛选条件
const filters = ref({
  search: '',
  status: '',
  level: '',
  dateRange: []
})

// 上传配置
const uploadAction = '/api/files/upload'

// 统计数据
const stats = ref({
  totalCompetitions: 0,
  completedCompetitions: 0,
  pendingResults: 0,
  publishedResults: 0
})

// 计算属性
const filteredCompetitions = computed(() => {
  let filtered = competitions.value

  if (filters.value.search) {
    const search = filters.value.search.toLowerCase()
    filtered = filtered.filter(c => c.title.toLowerCase().includes(search))
  }

  if (filters.value.status) {
    filtered = filtered.filter(c => c.status === filters.value.status)
  }

  if (filters.value.level) {
    filtered = filtered.filter(c => c.level === filters.value.level)
  }

  if (filters.value.dateRange && filters.value.dateRange.length === 2) {
    filtered = filtered.filter(c => {
      const completedDate = new Date(c.completedAt)
      const startDate = new Date(filters.value.dateRange[0])
      const endDate = new Date(filters.value.dateRange[1])
      return completedDate >= startDate && completedDate <= endDate
    })
  }

  return filtered
})

// 分页后的数据
const paginatedCompetitions = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredCompetitions.value.slice(start, end)
})

// 方法
const loadCompetitions = async () => {
  loading.value = true
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // 模拟数据
    competitions.value = [
      {
        id: 1,
        title: '全国大学生程序设计竞赛',
        type: '程序设计',
        level: 'national',
        organizer: '计算机学院',
        participantCount: 45,
        status: 'completed',
        resultStatus: 'pending',
        completedAt: new Date('2024-02-15'),
        publishedAt: null
      },
      {
        id: 2,
        title: '数学建模竞赛',
        type: '数学建模',
        level: 'provincial',
        organizer: '数学学院',
        participantCount: 78,
        status: 'completed',
        resultStatus: 'published',
        completedAt: new Date('2024-02-10'),
        publishedAt: new Date('2024-02-20')
      },
      {
        id: 3,
        title: '创新创业大赛',
        type: '创新创业',
        level: 'school',
        organizer: '商学院',
        participantCount: 120,
        status: 'review',
        resultStatus: null,
        completedAt: null,
        publishedAt: null
      }
    ]
    
    loadStats()
  } catch (error) {
    console.error('加载竞赛数据失败:', error)
    ElMessage.error('加载竞赛数据失败')
  } finally {
    loading.value = false
  }
}

const loadStats = () => {
  const total = competitions.value.length
  const completed = competitions.value.filter(c => c.status === 'completed').length
  const pending = competitions.value.filter(c => c.resultStatus === 'pending').length
  const published = competitions.value.filter(c => c.resultStatus === 'published').length
  
  stats.value = {
    totalCompetitions: total,
    completedCompetitions: completed,
    pendingResults: pending,
    publishedResults: published
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
    level: '',
    dateRange: []
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
  selectedCompetitions.value = selection
}

const viewCompetition = (competition) => {
  selectedCompetition.value = competition
  showDetailDialog.value = true
}

const manageResults = (competition) => {
  currentCompetition.value = competition
  loadResultsList(competition)
  showResultsDialog.value = true
}

const loadResultsList = (competition) => {
  // 模拟加载参赛学生列表
  resultsList.value = [
    {
      rank: 1,
      student: {
        realName: '张三',
        department: '计算机学院',
        studentId: '2021001'
      },
      awardLevel: 'first',
      finalScore: 95.5,
      certificateUrl: '',
      remarks: ''
    },
    {
      rank: 2,
      student: {
        realName: '李四',
        department: '计算机学院',
        studentId: '2021002'
      },
      awardLevel: 'second',
      finalScore: 92.0,
      certificateUrl: '',
      remarks: ''
    },
    {
      rank: 3,
      student: {
        realName: '王五',
        department: '计算机学院',
        studentId: '2021003'
      },
      awardLevel: 'second',
      finalScore: 89.5,
      certificateUrl: '',
      remarks: ''
    }
  ]
}

const saveResults = async () => {
  saving.value = true
  try {
    // 验证必填字段
    const invalidResults = resultsList.value.filter(r => !r.awardLevel || !r.finalScore)
    if (invalidResults.length > 0) {
      ElMessage.warning('请完善所有获奖结果信息')
      return
    }
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // 更新竞赛状态
    const index = competitions.value.findIndex(c => c.id === currentCompetition.value.id)
    if (index !== -1) {
      competitions.value[index].resultStatus = 'pending'
    }
    
    loadStats()
    ElMessage.success('结果保存成功！')
  } catch (error) {
    ElMessage.error('保存失败：' + (error.message || '未知错误'))
  } finally {
    saving.value = false
  }
}

const generateCertificates = () => {
  ElMessage.info('证书生成功能开发中...')
}

const previewResults = () => {
  // 生成预览数据
  const summary = {
    special: resultsList.value.filter(r => r.awardLevel === 'special').length,
    first: resultsList.value.filter(r => r.awardLevel === 'first').length,
    second: resultsList.value.filter(r => r.awardLevel === 'second').length,
    third: resultsList.value.filter(r => r.awardLevel === 'third').length,
    honorable: resultsList.value.filter(r => r.awardLevel === 'honorable').length
  }
  summary.total = Object.values(summary).reduce((sum, count) => sum + count, 0)
  
  previewData.value = {
    competitionTitle: currentCompetition.value.title,
    results: resultsList.value,
    summary
  }
  
  showPreviewDialog.value = true
}

const publishResults = async (competition) => {
  try {
    await ElMessageBox.confirm(
      `确定要发布竞赛 "${competition.title}" 的结果吗？发布后将无法修改。`,
      '确认发布',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // 更新状态
    const index = competitions.value.findIndex(c => c.id === competition.id)
    if (index !== -1) {
      competitions.value[index].resultStatus = 'published'
      competitions.value[index].publishedAt = new Date()
    }
    
    loadStats()
    ElMessage.success('结果发布成功！')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('发布失败：' + (error.message || '未知错误'))
    }
  }
}

const viewResults = (competition) => {
  selectedCompetition.value = competition
  showDetailDialog.value = true
}

const handleCertificateUpload = (response, row) => {
  row.certificateUrl = response.url
  ElMessage.success('证书上传成功')
}

const handleUploadError = () => {
  ElMessage.error('证书上传失败')
}

const exportResults = () => {
  ElMessage.success('开始导出结果数据...')
  // 实际项目中这里会调用导出API
}

// 状态相关方法
const getCompetitionStatusType = (status) => {
  const statusMap = {
    draft: 'info',
    registration: 'warning',
    submission: 'primary',
    review: 'warning',
    completed: 'success'
  }
  return statusMap[status] || 'info'
}

const getCompetitionStatusText = (status) => {
  const statusMap = {
    draft: '草稿',
    registration: '报名中',
    submission: '提交中',
    review: '评审中',
    completed: '已完成'
  }
  return statusMap[status] || status
}

const getResultStatusType = (status) => {
  const statusMap = {
    pending: 'warning',
    published: 'success'
  }
  return statusMap[status] || 'info'
}

const getResultStatusText = (status) => {
  const statusMap = {
    pending: '待发布',
    published: '已发布'
  }
  return statusMap[status] || '未管理'
}

const getAwardLevelType = (level) => {
  const levelMap = {
    special: 'danger',
    first: 'success',
    second: 'warning',
    third: 'info',
    honorable: 'primary'
  }
  return levelMap[level] || 'info'
}

const getAwardLevelText = (level) => {
  const levelMap = {
    special: '特等奖',
    first: '一等奖',
    second: '二等奖',
    third: '三等奖',
    honorable: '优秀奖'
  }
  return levelMap[level] || level
}

// 组件挂载时加载数据
onMounted(() => {
  loadCompetitions()
})
</script>

<style scoped>
.competition-results {
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

.stat-icon.total {
  background: #909399;
}

.stat-icon.completed {
  background: #67c23a;
}

.stat-icon.pending {
  background: #e6a23c;
}

.stat-icon.published {
  background: #409eff;
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

.competition-list-card {
  margin-bottom: 30px;
}

.competition-info h4 {
  margin: 0 0 5px 0;
  color: #2c3e50;
  font-size: 14px;
}

.competition-meta {
  margin: 0;
  color: #7f8c8d;
  font-size: 12px;
}

.competition-meta .type,
.competition-meta .level,
.competition-meta .organizer {
  margin-right: 10px;
}

.competition-meta .type {
  color: #409eff;
}

.competition-meta .level {
  color: #67c23a;
}

.action-buttons {
  display: flex;
  gap: 5px;
}

.pagination-container {
  margin-top: 20px;
  text-align: center;
}

/* 对话框样式 */
.results-dialog .el-dialog__body {
  max-height: 70vh;
  overflow-y: auto;
}

.competition-info-card,
.results-entry-card {
  margin-bottom: 20px;
}

.competition-info-card h3,
.results-entry-card h3 {
  margin: 0;
  color: #2c3e50;
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

.results-actions {
  margin-top: 20px;
  text-align: center;
}

.certificate-info {
  margin-top: 5px;
  font-size: 12px;
}

.preview-content h3 {
  margin: 0 0 20px 0;
  color: #2c3e50;
  text-align: center;
}

.preview-summary {
  margin-top: 30px;
}

.preview-summary h4 {
  margin: 0 0 15px 0;
  color: #2c3e50;
  font-size: 16px;
}

.results-section {
  margin-top: 20px;
}

.results-section h4 {
  margin: 0 0 15px 0;
  color: #2c3e50;
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
  
  .results-actions {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }
}
</style> 