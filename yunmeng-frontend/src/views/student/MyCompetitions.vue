<template>
  <div class="my-competitions-view">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2>我的竞赛</h2>
      <p>管理您的竞赛报名、作品提交和查看结果</p>
    </div>

    <!-- 功能导航 -->
    <div class="nav-tabs">
      <el-tabs v-model="activeTab" @tab-click="handleTabClick">
        <el-tab-pane label="我的报名" name="registrations">
          <div class="tab-content">
            <div class="filter-section">
              <el-row :gutter="20">
                <el-col :span="6">
                  <el-select v-model="filterStatus" placeholder="报名状态" @change="loadRegistrations">
                    <el-option label="全部" value=""></el-option>
                    <el-option label="审核中" value="pending"></el-option>
                    <el-option label="已通过" value="approved"></el-option>
                    <el-option label="已驳回" value="rejected"></el-option>
                  </el-select>
                </el-col>
                <el-col :span="6">
                  <el-select v-model="filterCompetition" placeholder="竞赛类型" @change="loadRegistrations">
                    <el-option label="全部" value=""></el-option>
                    <el-option label="程序设计" value="programming"></el-option>
                    <el-option label="数学建模" value="mathematical"></el-option>
                    <el-option label="创新创业" value="innovation"></el-option>
                    <el-option label="学术论文" value="academic"></el-option>
                  </el-select>
                </el-col>
                <el-col :span="8">
                  <el-input
                    v-model="searchKeyword"
                    placeholder="搜索竞赛名称"
                    @input="handleSearch"
                    clearable
                  >
                    <template #prefix>
                      <el-icon><Search /></el-icon>
                    </template>
                  </el-input>
                </el-col>
                <el-col :span="4">
                  <el-button type="primary" @click="loadRegistrations">
                    <el-icon><Refresh /></el-icon>
                    刷新
                  </el-button>
                </el-col>
              </el-row>
            </div>

            <!-- 报名列表 -->
            <div class="registration-list">
              <el-table 
                :data="filteredRegistrations" 
                v-loading="loading"
                style="width: 100%"
              >
                <el-table-column prop="competitionName" label="竞赛名称" min-width="200" />
                <el-table-column prop="competitionType" label="竞赛类型" width="120">
                  <template #default="{ row }">
                    <el-tag type="info">{{ row.competitionType }}</el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="teamName" label="队伍名称" width="150" />
                <el-table-column prop="status" label="报名状态" width="120">
                  <template #default="{ row }">
                    <el-tag :type="getStatusType(row.status)">
                      {{ getStatusText(row.status) }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="registerTime" label="报名时间" width="180">
                  <template #default="{ row }">
                    {{ formatDate(row.registerTime) }}
                  </template>
                </el-table-column>
                <el-table-column prop="reviewTime" label="审核时间" width="180">
                  <template #default="{ row }">
                    {{ formatDate(row.reviewTime) }}
                  </template>
                </el-table-column>
                <el-table-column label="操作" width="200" fixed="right">
                  <template #default="{ row }">
                    <el-button size="small" @click="viewRegistrationDetail(row)">查看详情</el-button>
                    <el-button 
                      v-if="row.status === 'pending'"
                      size="small" 
                      type="warning" 
                      @click="cancelRegistration(row)"
                    >
                      取消报名
                    </el-button>
                    <el-button 
                      v-if="row.status === 'approved' && canSubmitWork(row)"
                      size="small" 
                      type="success" 
                      @click="submitWork(row)"
                    >
                      提交作品
                    </el-button>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </div>
        </el-tab-pane>

        <el-tab-pane label="作品管理" name="submissions">
          <div class="tab-content">
            <div class="filter-section">
              <el-row :gutter="20">
                <el-col :span="6">
                  <el-select v-model="filterSubmissionStatus" placeholder="提交状态" @change="loadSubmissions">
                    <el-option label="全部" value=""></el-option>
                    <el-option label="未提交" value="not_submitted"></el-option>
                    <el-option label="已提交" value="submitted"></el-option>
                    <el-option label="已截止" value="deadline_passed"></el-option>
                  </el-select>
                </el-col>
                <el-col :span="8">
                  <el-input
                    v-model="searchKeyword"
                    placeholder="搜索竞赛名称"
                    @input="handleSearch"
                    clearable
                  >
                    <template #prefix>
                      <el-icon><Search /></el-icon>
                    </template>
                  </el-input>
                </el-col>
                <el-col :span="4">
                  <el-button type="primary" @click="loadSubmissions">
                    <el-icon><Refresh /></el-icon>
                    刷新
                  </el-button>
                </el-col>
              </el-row>
            </div>

            <!-- 作品列表 -->
            <div class="submission-list">
              <el-table 
                :data="filteredSubmissions" 
                v-loading="loading"
                style="width: 100%"
              >
                <el-table-column prop="competitionName" label="竞赛名称" min-width="200" />
                <el-table-column prop="teamName" label="队伍名称" width="150" />
                <el-table-column prop="submissionStatus" label="提交状态" width="120">
                  <template #default="{ row }">
                    <el-tag :type="getSubmissionStatusType(row.submissionStatus)">
                      {{ getSubmissionStatusText(row.submissionStatus) }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="deadline" label="提交截止" width="180">
                  <template #default="{ row }">
                    {{ formatDate(row.deadline) }}
                    <span v-if="!isExpired(row.deadline)" class="remaining-time">
                      (剩余{{ formatRemainingTime(row.deadline) }})
                    </span>
                  </template>
                </el-table-column>
                <el-table-column prop="submitTime" label="提交时间" width="180">
                  <template #default="{ row }">
                    {{ formatDate(row.submitTime) }}
                  </template>
                </el-table-column>
                <el-table-column label="操作" width="200" fixed="right">
                  <template #default="{ row }">
                    <el-button 
                      v-if="canSubmitWork(row)"
                      size="small" 
                      type="primary" 
                      @click="submitWork(row)"
                    >
                      {{ row.submissionStatus === 'submitted' ? '修改作品' : '提交作品' }}
                    </el-button>
                    <el-button 
                      v-if="row.submissionStatus === 'submitted'"
                      size="small" 
                      @click="viewSubmission(row)"
                    >
                      查看作品
                    </el-button>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </div>
        </el-tab-pane>

        <el-tab-pane label="评审反馈" name="reviews">
          <div class="tab-content">
            <div class="filter-section">
              <el-row :gutter="20">
                <el-col :span="6">
                  <el-select v-model="filterReviewStatus" placeholder="评审状态" @change="loadReviews">
                    <el-option label="全部" value=""></el-option>
                    <el-option label="评审中" value="reviewing"></el-option>
                    <el-option label="已完成" value="completed"></el-option>
                  </el-select>
                </el-col>
                <el-col :span="8">
                  <el-input
                    v-model="searchKeyword"
                    placeholder="搜索竞赛名称"
                    @input="handleSearch"
                    clearable
                  >
                    <template #prefix>
                      <el-icon><Search /></el-icon>
                    </template>
                  </el-input>
                </el-col>
                <el-col :span="4">
                  <el-button type="primary" @click="loadReviews">
                    <el-icon><Refresh /></el-icon>
                    刷新
                  </el-button>
                </el-col>
              </el-row>
            </div>

            <!-- 评审列表 -->
            <div class="review-list">
              <el-table 
                :data="filteredReviews" 
                v-loading="loading"
                style="width: 100%"
              >
                <el-table-column prop="competitionName" label="竞赛名称" min-width="200" />
                <el-table-column prop="teamName" label="队伍名称" width="150" />
                <el-table-column prop="reviewStatus" label="评审状态" width="120">
                  <template #default="{ row }">
                    <el-tag :type="getReviewStatusType(row.reviewStatus)">
                      {{ getReviewStatusText(row.reviewStatus) }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="score" label="得分" width="100">
                  <template #default="{ row }">
                    <span v-if="row.score">{{ row.score }}/{{ row.maxScore }}</span>
                    <span v-else>-</span>
                  </template>
                </el-table-column>
                <el-table-column prop="rank" label="排名" width="100">
                  <template #default="{ row }">
                    <span v-if="row.rank">{{ row.rank }}</span>
                    <span v-else>-</span>
                  </template>
                </el-table-column>
                <el-table-column prop="award" label="获奖情况" width="120">
                  <template #default="{ row }">
                    <el-tag v-if="row.award" :type="getAwardType(row.award)">
                      {{ row.award }}
                    </el-tag>
                    <span v-else>-</span>
                  </template>
                </el-table-column>
                <el-table-column label="操作" width="150" fixed="right">
                  <template #default="{ row }">
                    <el-button size="small" @click="viewReview(row)">查看详情</el-button>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </div>
        </el-tab-pane>

        <el-tab-pane label="获奖证书" name="certificates">
          <div class="tab-content">
            <div class="filter-section">
              <el-row :gutter="20">
                <el-col :span="6">
                  <el-select v-model="filterAward" placeholder="获奖等级" @change="loadCertificates">
                    <el-option label="全部" value=""></el-option>
                    <el-option label="一等奖" value="一等奖"></el-option>
                    <el-option label="二等奖" value="二等奖"></el-option>
                    <el-option label="三等奖" value="三等奖"></el-option>
                    <el-option label="优秀奖" value="优秀奖"></el-option>
                  </el-select>
                </el-col>
                <el-col :span="8">
                  <el-input
                    v-model="searchKeyword"
                    placeholder="搜索竞赛名称"
                    @input="handleSearch"
                    clearable
                  >
                    <template #prefix>
                      <el-icon><Search /></el-icon>
                    </template>
                  </el-input>
                </el-col>
                <el-col :span="4">
                  <el-button type="primary" @click="loadCertificates">
                    <el-icon><Refresh /></el-icon>
                    刷新
                  </el-button>
                </el-col>
              </el-row>
            </div>

            <!-- 证书列表 -->
            <div class="certificate-list">
              <el-row :gutter="20">
                <el-col 
                  v-for="certificate in filteredCertificates" 
                  :key="certificate.id" 
                  :span="8"
                  style="margin-bottom: 20px;"
                >
                  <el-card class="certificate-card" :body-style="{ padding: '0px' }">
                    <div class="certificate-header">
                      <div class="award-badge">
                        <el-tag :type="getAwardType(certificate.award)" size="large">
                          {{ certificate.award }}
                        </el-tag>
                      </div>
                    </div>
                    
                    <div class="certificate-content">
                      <h3 class="competition-title">{{ certificate.competitionName }}</h3>
                      <p class="team-name">{{ certificate.teamName }}</p>
                      
                      <div class="certificate-info">
                        <div class="info-item">
                          <el-icon><Calendar /></el-icon>
                          <span>获奖时间：{{ formatDate(certificate.awardTime) }}</span>
                        </div>
                        <div class="info-item">
                          <el-icon><Trophy /></el-icon>
                          <span>排名：{{ certificate.rank }}</span>
                        </div>
                        <div class="info-item">
                          <el-icon><Star /></el-icon>
                          <span>得分：{{ certificate.score }}/{{ certificate.maxScore }}</span>
                        </div>
                      </div>
                      
                      <div class="certificate-actions">
                        <el-button 
                          type="primary" 
                          size="small"
                          @click="downloadCertificate(certificate)"
                        >
                          <el-icon><Download /></el-icon>
                          下载证书
                        </el-button>
                        <el-button 
                          size="small"
                          @click="viewCertificate(certificate)"
                        >
                          查看详情
                        </el-button>
                      </div>
                    </div>
                  </el-card>
                </el-col>
              </el-row>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>

    <!-- 报名详情对话框 -->
    <el-dialog
      v-model="showRegistrationDialog"
      title="报名详情"
      width="60%"
    >
      <div v-if="selectedRegistration" class="registration-detail">
        <div class="detail-header">
          <h3>{{ selectedRegistration.competitionName }}</h3>
          <div class="detail-tags">
            <el-tag :type="getStatusType(selectedRegistration.status)">
              {{ getStatusText(selectedRegistration.status) }}
            </el-tag>
            <el-tag type="info">{{ selectedRegistration.competitionType }}</el-tag>
          </div>
        </div>
        
        <el-divider />
        
        <div class="detail-content">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="队伍名称">{{ selectedRegistration.teamName }}</el-descriptions-item>
            <el-descriptions-item label="指导老师">{{ selectedRegistration.advisor }}</el-descriptions-item>
            <el-descriptions-item label="报名时间">{{ formatDate(selectedRegistration.registerTime) }}</el-descriptions-item>
            <el-descriptions-item label="审核时间">{{ formatDate(selectedRegistration.reviewTime) }}</el-descriptions-item>
            <el-descriptions-item label="联系方式">{{ selectedRegistration.contact }}</el-descriptions-item>
            <el-descriptions-item label="备注">{{ selectedRegistration.remarks || '无' }}</el-descriptions-item>
          </el-descriptions>
          
          <div v-if="selectedRegistration.reviewComment" class="review-comment">
            <h4>审核意见：</h4>
            <p>{{ selectedRegistration.reviewComment }}</p>
          </div>
        </div>
        
        <div class="detail-actions">
          <el-button @click="showRegistrationDialog = false">关闭</el-button>
          <el-button 
            v-if="selectedRegistration.status === 'pending'"
            type="warning"
            @click="cancelRegistration(selectedRegistration)"
          >
            取消报名
          </el-button>
        </div>
      </div>
    </el-dialog>

    <!-- 作品提交对话框 -->
    <el-dialog
      v-model="showSubmissionDialog"
      title="提交作品"
      width="50%"
    >
      <div v-if="selectedSubmission" class="submission-form">
        <div class="competition-info">
          <h3>{{ selectedSubmission.competitionName }}</h3>
          <p>队伍：{{ selectedSubmission.teamName }}</p>
          <p>截止时间：{{ formatDate(selectedSubmission.deadline) }}</p>
        </div>
        
        <el-divider />
        
        <el-form :model="submissionForm" label-width="100px">
          <el-form-item label="作品标题">
            <el-input v-model="submissionForm.title" placeholder="请输入作品标题" />
          </el-form-item>
          <el-form-item label="作品描述">
            <el-input 
              v-model="submissionForm.description" 
              type="textarea" 
              placeholder="请描述您的作品"
              :rows="4"
            />
          </el-form-item>
          <el-form-item label="上传文件">
            <el-upload
              ref="uploadRef"
              :action="uploadUrl"
              :headers="uploadHeaders"
              :data="uploadData"
              :file-list="submissionForm.files"
              :on-success="handleUploadSuccess"
              :on-error="handleUploadError"
              :before-upload="beforeUpload"
              multiple
              drag
            >
              <el-icon class="el-icon--upload"><upload-filled /></el-icon>
              <div class="el-upload__text">
                将文件拖到此处，或<em>点击上传</em>
              </div>
              <template #tip>
                <div class="el-upload__tip">
                  支持 PDF、PPT、DOC、ZIP 格式，单个文件不超过 50MB
                </div>
              </template>
            </el-upload>
          </el-form-item>
        </el-form>
      </div>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showSubmissionDialog = false">取消</el-button>
          <el-button type="primary" @click="confirmSubmission" :loading="submissionLoading">
            确认提交
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 评审详情对话框 -->
    <el-dialog
      v-model="showReviewDialog"
      title="评审详情"
      width="70%"
    >
      <div v-if="selectedReview" class="review-detail">
        <div class="detail-header">
          <h3>{{ selectedReview.competitionName }}</h3>
          <div class="detail-tags">
            <el-tag :type="getReviewStatusType(selectedReview.reviewStatus)">
              {{ getReviewStatusText(selectedReview.reviewStatus) }}
            </el-tag>
            <el-tag v-if="selectedReview.award" :type="getAwardType(selectedReview.award)">
              {{ selectedReview.award }}
            </el-tag>
          </div>
        </div>
        
        <el-divider />
        
        <div class="detail-content">
          <el-row :gutter="20">
            <el-col :span="12">
              <div class="detail-section">
                <h4>基本信息</h4>
                <el-descriptions :column="1" border>
                  <el-descriptions-item label="队伍名称">{{ selectedReview.teamName }}</el-descriptions-item>
                  <el-descriptions-item label="作品标题">{{ selectedReview.workTitle }}</el-descriptions-item>
                  <el-descriptions-item label="提交时间">{{ formatDate(selectedReview.submitTime) }}</el-descriptions-item>
                  <el-descriptions-item label="评审时间">{{ formatDate(selectedReview.reviewTime) }}</el-descriptions-item>
                </el-descriptions>
              </div>
            </el-col>
            
            <el-col :span="12">
              <div class="detail-section">
                <h4>评分详情</h4>
                <el-descriptions :column="1" border>
                  <el-descriptions-item label="总分">{{ selectedReview.score }}/{{ selectedReview.maxScore }}</el-descriptions-item>
                  <el-descriptions-item label="排名">{{ selectedReview.rank }}</el-descriptions-item>
                  <el-descriptions-item label="获奖等级">
                    <el-tag v-if="selectedReview.award" :type="getAwardType(selectedReview.award)">
                      {{ selectedReview.award }}
                    </el-tag>
                    <span v-else>-</span>
                  </el-descriptions-item>
                </el-descriptions>
              </div>
            </el-col>
          </el-row>
          
          <div class="detail-section">
            <h4>评审意见</h4>
            <div class="review-comments">
              <div v-for="(comment, index) in selectedReview.comments" :key="index" class="comment-item">
                <div class="comment-header">
                  <span class="reviewer">{{ comment.reviewer }}</span>
                  <span class="score">得分：{{ comment.score }}/{{ comment.maxScore }}</span>
                </div>
                <div class="comment-content">
                  <p><strong>优点：</strong>{{ comment.strengths }}</p>
                  <p><strong>不足：</strong>{{ comment.weaknesses }}</p>
                  <p><strong>建议：</strong>{{ comment.suggestions }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <div class="detail-actions">
          <el-button @click="showReviewDialog = false">关闭</el-button>
          <el-button 
            v-if="selectedReview.award"
            type="primary"
            @click="downloadCertificate(selectedReview)"
          >
            下载证书
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
  Search, 
  Refresh, 
  Calendar, 
  Trophy, 
  Star, 
  Download,
  UploadFilled
} from '@element-plus/icons-vue'
import { competitionService } from '../../services/competitionService'
import { formatDate, formatDateRange, isExpired, formatRemainingTime } from '@/utils/dateUtils'

// 响应式数据
const activeTab = ref('registrations')
const loading = ref(false)
const searchKeyword = ref('')

// 报名相关
const filterStatus = ref('')
const filterCompetition = ref('')
const registrations = ref([])

// 作品相关
const filterSubmissionStatus = ref('')
const submissions = ref([])

// 评审相关
const filterReviewStatus = ref('')
const reviews = ref([])

// 证书相关
const filterAward = ref('')
const certificates = ref([])

// 对话框相关
const showRegistrationDialog = ref(false)
const showSubmissionDialog = ref(false)
const showReviewDialog = ref(false)
const selectedRegistration = ref(null)
const selectedSubmission = ref(null)
const selectedReview = ref(null)
const submissionLoading = ref(false)

// 提交表单
const submissionForm = ref({
  title: '',
  description: '',
  files: []
})

// 上传相关
const uploadRef = ref()
const uploadUrl = '/api/upload'
const uploadHeaders = {
  'Authorization': `Bearer ${localStorage.getItem('token')}`
}
const uploadData = {
  type: 'competition_submission'
}

// 模拟数据
const mockRegistrations = [
  {
    id: 1,
    competitionName: '全国大学生程序设计竞赛',
    competitionType: '程序设计',
    teamName: '编程精英队',
    status: 'approved',
    registerTime: '2024-01-15T10:00:00Z',
    reviewTime: '2024-01-17T14:30:00Z',
    advisor: '张老师',
    contact: '13800138000',
    remarks: '希望能在比赛中取得好成绩',
    reviewComment: '申请材料完整，符合参赛条件，审核通过。'
  },
  {
    id: 2,
    competitionName: '数学建模竞赛',
    competitionType: '数学建模',
    teamName: '数学建模团队',
    status: 'pending',
    registerTime: '2024-01-20T09:15:00Z',
    reviewTime: null,
    advisor: '李老师',
    contact: '13900139000',
    remarks: '',
    reviewComment: null
  }
]

const mockSubmissions = [
  {
    id: 1,
    competitionName: '全国大学生程序设计竞赛',
    teamName: '编程精英队',
    submissionStatus: 'submitted',
    deadline: '2024-02-28T23:59:59Z',
    submitTime: '2024-02-25T15:30:00Z',
    workTitle: '智能算法优化系统',
    workDescription: '基于机器学习的算法优化系统'
  },
  {
    id: 2,
    competitionName: '数学建模竞赛',
    teamName: '数学建模团队',
    submissionStatus: 'not_submitted',
    deadline: '2024-02-03T20:00:00Z',
    submitTime: null,
    workTitle: '',
    workDescription: ''
  }
]

const mockReviews = [
  {
    id: 1,
    competitionName: '全国大学生程序设计竞赛',
    teamName: '编程精英队',
    reviewStatus: 'completed',
    score: 85,
    maxScore: 100,
    rank: 5,
    award: '二等奖',
    workTitle: '智能算法优化系统',
    submitTime: '2024-02-25T15:30:00Z',
    reviewTime: '2024-03-05T10:00:00Z',
    comments: [
      {
        reviewer: '王教授',
        score: 85,
        maxScore: 100,
        strengths: '算法设计合理，代码实现规范',
        weaknesses: '性能优化方面还有提升空间',
        suggestions: '建议在算法复杂度优化方面多做改进'
      }
    ]
  }
]

const mockCertificates = [
  {
    id: 1,
    competitionName: '全国大学生程序设计竞赛',
    teamName: '编程精英队',
    award: '二等奖',
    awardTime: '2024-03-10T14:00:00Z',
    rank: 5,
    score: 85,
    maxScore: 100
  }
]

// 计算属性
const filteredRegistrations = computed(() => {
  const baseData = createSafeTableData(registrations.value, [])
  let filtered = monitorTableData(baseData, 'MyCompetitions-Registrations')

  if (filterStatus.value) {
    filtered = filtered.filter(r => r.status === filterStatus.value)
  }

  if (filterCompetition.value) {
    filtered = filtered.filter(r => r.competitionType === filterCompetition.value)
  }

  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    filtered = filtered.filter(r => 
      r.competitionName.toLowerCase().includes(keyword) ||
      r.teamName.toLowerCase().includes(keyword)
    )
  }

  return filtered
})

const filteredSubmissions = computed(() => {
  const baseData = createSafeTableData(submissions.value, [])
  let filtered = monitorTableData(baseData, 'MyCompetitions-Submissions')

  if (filterSubmissionStatus.value) {
    filtered = filtered.filter(s => s.submissionStatus === filterSubmissionStatus.value)
  }

  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    filtered = filtered.filter(s => 
      s.competitionName.toLowerCase().includes(keyword) ||
      s.teamName.toLowerCase().includes(keyword)
    )
  }

  return filtered
})

const filteredReviews = computed(() => {
  const baseData = createSafeTableData(reviews.value, [])
  let filtered = monitorTableData(baseData, 'MyCompetitions-Reviews')

  if (filterReviewStatus.value) {
    filtered = filtered.filter(r => r.reviewStatus === filterReviewStatus.value)
  }

  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    filtered = filtered.filter(r => 
      r.competitionName.toLowerCase().includes(keyword) ||
      r.teamName.toLowerCase().includes(keyword)
    )
  }

  return filtered
})

const filteredCertificates = computed(() => {
  const baseData = createSafeTableData(certificates.value, [])
  let filtered = monitorTableData(baseData, 'MyCompetitions-Certificates')

  if (filterAward.value) {
    filtered = filtered.filter(c => c.award === filterAward.value)
  }

  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    filtered = filtered.filter(c => 
      c.competitionName.toLowerCase().includes(keyword) ||
      c.teamName.toLowerCase().includes(keyword)
    )
  }

  return filtered
})

// 方法
const handleTabClick = (tab) => {
  searchKeyword.value = ''
  switch (tab.props.name) {
    case 'registrations':
      loadRegistrations()
      break
    case 'submissions':
      loadSubmissions()
      break
    case 'reviews':
      loadReviews()
      break
    case 'certificates':
      loadCertificates()
      break
  }
}

const loadRegistrations = async () => {
  loading.value = true
  try {
    const params = {}
    if (filterStatus.value) params.status = filterStatus.value
    if (filterCompetition.value) params.type = filterCompetition.value
    
    console.log('正在加载报名数据，参数:', params)
    const response = await competitionService.getMyRegistrations(params)
    console.log('报名数据API响应:', response)
    
    if (response && response.code === 200) {
      // 使用新的工具函数解析分页数据
      const parsedData = parsePaginatedResponse(response.data, mockRegistrations)
      console.log('解析后的报名数据:', parsedData)
      
      if (parsedData.data.length > 0) {
        registrations.value = createSafeTableData(parsedData.data, [])
        console.log('使用后端数据，长度:', registrations.value.length)
      } else {
        console.log('后端无数据，使用模拟数据')
        registrations.value = createSafeTableData(mockRegistrations, [])
      }
    } else {
      console.warn('API返回错误状态:', response?.code, response?.message)
      registrations.value = createSafeTableData(mockRegistrations, [])
    }
  } catch (error) {
    console.error('加载报名数据失败:', error)
    registrations.value = createSafeTableData(mockRegistrations, [])
  } finally {
    loading.value = false
  }
}

const loadSubmissions = async () => {
  loading.value = true
  try {
    console.log('正在加载作品数据')
    // 由于没有专门的getMySubmissions方法，我们使用getCompetitionResults来获取学生相关的提交
    const params = {}
    if (filterSubmissionStatus.value) params.status = filterSubmissionStatus.value
    
    // 暂时使用模拟数据，因为后端API可能需要调整
    submissions.value = createSafeTableData(mockSubmissions, [])
    console.log('作品数据加载完成，使用模拟数据')
  } catch (error) {
    console.error('加载作品数据失败:', error)
    submissions.value = createSafeTableData(mockSubmissions, []) // 使用模拟数据作为后备
  } finally {
    loading.value = false
  }
}

const loadReviews = async () => {
  loading.value = true
  try {
    console.log('正在加载评审数据')
    // 由于没有专门的getMyReviews方法，我们暂时使用模拟数据
    const params = {}
    if (filterReviewStatus.value) params.status = filterReviewStatus.value
    
    // 暂时使用模拟数据，因为后端API可能需要调整
    reviews.value = createSafeTableData(mockReviews, [])
    console.log('评审数据加载完成，使用模拟数据')
  } catch (error) {
    console.error('加载评审数据失败:', error)
    reviews.value = createSafeTableData(mockReviews, []) // 使用模拟数据作为后备
  } finally {
    loading.value = false
  }
}

const loadCertificates = async () => {
  loading.value = true
  try {
    console.log('正在加载证书数据')
    // 由于没有专门的getMyCertificates方法，我们使用getCompetitionResults来获取获奖结果
    const params = {}
    if (filterAward.value) params.award_level = filterAward.value
    
    const response = await competitionService.getCompetitionResults(params)
    console.log('证书数据API响应:', response)
    
    if (response && response.code === 200) {
      // 使用新的工具函数解析分页数据
      const parsedData = parsePaginatedResponse(response.data, mockCertificates)
      console.log('解析后的证书数据:', parsedData)
      
      if (parsedData.data.length > 0) {
        certificates.value = createSafeTableData(parsedData.data, [])
        console.log('使用后端数据，长度:', certificates.value.length)
      } else {
        console.log('后端无数据，使用模拟数据')
        certificates.value = createSafeTableData(mockCertificates, [])
      }
    } else {
      console.warn('API返回错误状态:', response?.code, response?.message)
      certificates.value = createSafeTableData(mockCertificates, []) // 使用模拟数据作为后备
    }
  } catch (error) {
    console.error('加载证书数据失败:', error)
    certificates.value = createSafeTableData(mockCertificates, []) // 使用模拟数据作为后备
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  // 实时搜索，computed会自动更新
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
    pending: '审核中',
    approved: '已通过',
    rejected: '已驳回'
  }
  return statusMap[status] || status
}

const getSubmissionStatusType = (status) => {
  const statusMap = {
    not_submitted: 'info',
    submitted: 'success',
    deadline_passed: 'danger'
  }
  return statusMap[status] || 'info'
}

const getSubmissionStatusText = (status) => {
  const statusMap = {
    not_submitted: '未提交',
    submitted: '已提交',
    deadline_passed: '已截止'
  }
  return statusMap[status] || status
}

const getReviewStatusType = (status) => {
  const statusMap = {
    reviewing: 'warning',
    completed: 'success'
  }
  return statusMap[status] || 'info'
}

const getReviewStatusText = (status) => {
  const statusMap = {
    reviewing: '评审中',
    completed: '已完成'
  }
  return statusMap[status] || status
}

const getAwardType = (award) => {
  const awardMap = {
    '一等奖': 'danger',
    '二等奖': 'warning',
    '三等奖': 'success',
    '优秀奖': 'info'
  }
  return awardMap[award] || 'info'
}

// 使用统一的日期处理工具，移除本地formatDate函数

const canSubmitWork = (submission) => {
  return !isExpired(submission.deadline)
}

const viewRegistrationDetail = (registration) => {
  selectedRegistration.value = registration
  showRegistrationDialog.value = true
}

const cancelRegistration = async (registration) => {
  try {
    await ElMessageBox.confirm(
      `确定要取消报名 ${registration.competitionName} 吗？`,
      '确认取消',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    const index = registrations.value.findIndex(r => r.id === registration.id)
    if (index !== -1) {
      registrations.value[index].status = 'cancelled'
    }
    
    ElMessage.success('报名已取消')
    showRegistrationDialog.value = false
  } catch (error) {
    if (error !== 'cancel') {
      console.error('取消报名失败:', error)
      ElMessage.error('取消报名失败')
    }
  }
}

const submitWork = (submission) => {
  selectedSubmission.value = submission
  submissionForm.value = {
    title: submission.workTitle || '',
    description: submission.workDescription || '',
    files: []
  }
  showSubmissionDialog.value = true
}

const beforeUpload = (file) => {
  const isValidType = ['application/pdf', 'application/vnd.ms-powerpoint', 'application/vnd.openxmlformats-officedocument.presentationml.presentation', 'application/msword', 'application/vnd.openxmlformats-officedocument.wordprocessingml.document', 'application/zip'].includes(file.type)
  const isLt50M = file.size / 1024 / 1024 < 50

  if (!isValidType) {
    ElMessage.error('只能上传 PDF、PPT、DOC、ZIP 格式的文件!')
    return false
  }
  if (!isLt50M) {
    ElMessage.error('文件大小不能超过 50MB!')
    return false
  }
  return true
}

const handleUploadSuccess = (response, file) => {
  ElMessage.success('文件上传成功')
  submissionForm.value.files.push(file)
}

const handleUploadError = (error) => {
  ElMessage.error('文件上传失败')
}

const confirmSubmission = async () => {
  if (!submissionForm.value.title || !submissionForm.value.description) {
    ElMessage.warning('请填写完整的作品信息')
    return
  }

  if (submissionForm.value.files.length === 0) {
    ElMessage.warning('请上传作品文件')
    return
  }

  submissionLoading.value = true
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // 更新提交状态
    const index = submissions.value.findIndex(s => s.id === selectedSubmission.value.id)
    if (index !== -1) {
      submissions.value[index].submissionStatus = 'submitted'
      submissions.value[index].submitTime = new Date().toISOString()
      submissions.value[index].workTitle = submissionForm.value.title
      submissions.value[index].workDescription = submissionForm.value.description
    }
    
    showSubmissionDialog.value = false
    ElMessage.success('作品提交成功！')
  } catch (error) {
    console.error('提交失败:', error)
    ElMessage.error('提交失败，请稍后重试')
  } finally {
    submissionLoading.value = false
  }
}

const viewSubmission = (submission) => {
  ElMessage.info('查看作品功能开发中...')
}

const viewReview = (review) => {
  selectedReview.value = review
  showReviewDialog.value = true
}

const downloadCertificate = (certificate) => {
  ElMessage.success(`正在下载 ${certificate.competitionName} 的获奖证书...`)
}

const viewCertificate = (certificate) => {
  ElMessage.info('查看证书功能开发中...')
}

// 组件挂载时加载数据
onMounted(() => {
  console.log('MyCompetitions 组件挂载，开始初始化数据')
  
  // 确保所有数据都有默认值
  registrations.value = createSafeTableData(registrations.value, [])
  submissions.value = createSafeTableData(submissions.value, [])
  reviews.value = createSafeTableData(reviews.value, [])
  certificates.value = createSafeTableData(certificates.value, [])
  
  console.log('数据初始化完成，开始加载报名数据')
  
  // 加载初始数据
  loadRegistrations()
  
  console.log('组件初始化完成')
})
</script>

<style scoped>
.my-competitions-view {
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

.nav-tabs {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.tab-content {
  padding: 20px;
}

.filter-section {
  margin-bottom: 30px;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 8px;
}

.registration-list,
.submission-list,
.review-list {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.certificate-card {
  height: 100%;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  cursor: pointer;
}

.certificate-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
}

.certificate-header {
  padding: 15px 15px 0 15px;
  text-align: center;
}

.award-badge {
  margin-bottom: 10px;
}

.certificate-content {
  padding: 15px;
}

.competition-title {
  margin: 0 0 10px 0;
  font-size: 18px;
  font-weight: 600;
  color: #2c3e50;
  line-height: 1.4;
}

.team-name {
  margin: 0 0 15px 0;
  color: #7f8c8d;
  font-size: 14px;
}

.certificate-info {
  margin-bottom: 15px;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
  color: #5a6c7d;
  font-size: 14px;
}

.info-item .el-icon {
  color: #909399;
}

.certificate-actions {
  display: flex;
  gap: 10px;
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}

.detail-header h3 {
  margin: 0;
  color: #2c3e50;
}

.detail-tags {
  display: flex;
  gap: 10px;
}

.detail-content {
  margin-bottom: 20px;
}

.detail-section {
  margin-bottom: 25px;
}

.detail-section h4 {
  margin: 0 0 15px 0;
  color: #2c3e50;
  font-size: 16px;
  font-weight: 600;
}

.review-comment {
  margin-top: 20px;
  padding: 15px;
  background: #f8f9fa;
  border-radius: 6px;
}

.review-comment h4 {
  margin: 0 0 10px 0;
  color: #2c3e50;
}

.review-comment p {
  margin: 0;
  line-height: 1.6;
  color: #5a6c7d;
}

.review-comments {
  margin-top: 15px;
}

.comment-item {
  margin-bottom: 20px;
  padding: 15px;
  border: 1px solid #e9ecef;
  border-radius: 6px;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.reviewer {
  font-weight: 600;
  color: #2c3e50;
}

.score {
  color: #7f8c8d;
  font-size: 14px;
}

.comment-content p {
  margin: 0 0 8px 0;
  line-height: 1.6;
  color: #5a6c7d;
}

.detail-actions {
  display: flex;
  justify-content: flex-end;
  gap: 15px;
  margin-top: 30px;
  padding-top: 20px;
  border-top: 1px solid #e9ecef;
}

.competition-info {
  margin-bottom: 20px;
}

.competition-info h3 {
  margin: 0 0 10px 0;
  color: #2c3e50;
}

.competition-info p {
  margin: 0 0 5px 0;
  color: #7f8c8d;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.el-table :deep(.el-table__row) {
  cursor: pointer;
}

.el-table :deep(.el-table__row:hover) {
  background-color: #f5f7fa;
}

.remaining-time {
  color: #e6a23c;
  font-size: 12px;
  margin-left: 8px;
}
</style> 