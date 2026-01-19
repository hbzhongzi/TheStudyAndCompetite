<template>
  <div class="competition-view">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2>竞赛信息</h2>
      <p>查看和报名参加各类学术竞赛</p>
    </div>

    <!-- 筛选和搜索 -->
    <div class="filter-section">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-select v-model="filterStatus" placeholder="竞赛状态" @change="handleFilterChange">
            <el-option label="全部" value=""></el-option>
            <el-option label="即将开始" value="upcoming"></el-option>
            <el-option label="进行中" value="ongoing"></el-option>
            <el-option label="已结束" value="ended"></el-option>
          </el-select>
        </el-col>
        <el-col :span="6">
          <el-select v-model="filterType" placeholder="竞赛类型" @change="handleFilterChange">
            <el-option label="全部" value=""></el-option>
            <el-option label="程序设计" value="程序设计"></el-option>
            <el-option label="数学建模" value="数学建模"></el-option>
            <el-option label="创新创业" value="创新创业"></el-option>
            <el-option label="学术论文" value="学术论文"></el-option>
            <el-option label="工程设计" value="工程设计"></el-option>
            <el-option label="实验技能" value="实验技能"></el-option>
            <el-option label="语言技能" value="语言技能"></el-option>
          </el-select>
        </el-col>
        <el-col :span="6">
          <el-select v-model="filterCollege" placeholder="适用学院" @change="handleFilterChange">
            <el-option label="全部" value=""></el-option>
            <el-option label="计算机学院" value="computer"></el-option>
            <el-option label="数学学院" value="mathematics"></el-option>
            <el-option label="物理学院" value="physics"></el-option>
            <el-option label="化学学院" value="chemistry"></el-option>
            <el-option label="工程学院" value="engineering"></el-option>
            <el-option label="商学院" value="business"></el-option>
            <el-option label="学术研究" value="academic"></el-option>
            <el-option label="语言学院" value="language"></el-option>
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
          <el-button type="primary" @click="loadCompetitions">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </el-col>
      </el-row>
    </div>

    <!-- 竞赛列表 -->
    <div class="competition-list">
      <el-row :gutter="20">
        <el-col 
          v-for="competition in paginatedCompetitions" 
          :key="competition.id" 
          :span="8"
          style="margin-bottom: 20px;"
        >
          <el-card class="competition-card" :body-style="{ padding: '0px' }">
            <div class="competition-header">
              <div class="competition-status">
                <el-tag :type="getStatusType(getCompetitionStatus(competition))">
                  {{ getStatusText(getCompetitionStatus(competition)) }}
                </el-tag>
              </div>
              <div class="competition-type">
                <el-tag size="small" type="info">{{ getCompetitionType(competition) }}</el-tag>
              </div>
            </div>
            
            <div class="competition-content">
              <h3 class="competition-title">{{ competition.title }}</h3>
              <p class="competition-description">{{ competition.description }}</p>
              <div class="competition-meta" v-if="competition.organizer || competition.location">
                <span v-if="competition.organizer" class="meta-item">
                  <el-icon><OfficeBuilding /></el-icon>
                  {{ competition.organizer }}
                </span>
                <span v-if="competition.location" class="meta-item">
                  <el-icon><Location /></el-icon>
                  {{ competition.location }}
                </span>
              </div>
              
              <div class="competition-info">
                <div class="info-item">
                  <el-icon><Calendar /></el-icon>
                  <span>报名截止：{{ formatRegistrationDeadline(competition) }}</span>
                  <span v-if="!isRegistrationExpired(competition)" class="remaining-time">
                    (剩余{{ formatRemainingTime(getRegistrationDeadline(competition)) }})
                  </span>
                </div>
                <div class="info-item">
                  <el-icon><Clock /></el-icon>
                  <span>比赛时间：{{ formatCompetitionTime(competition) }}</span>
                </div>
                <div class="info-item">
                  <el-icon><User /></el-icon>
                  <span>已报名：{{ getRegisteredCount(competition) }}/{{ getMaxParticipants(competition) }}</span>
                </div>
                <div class="info-item" v-if="competition.organizer">
                  <el-icon><OfficeBuilding /></el-icon>
                  <span>主办方：{{ competition.organizer }}</span>
                </div>
                <div class="info-item" v-if="competition.location">
                  <el-icon><Location /></el-icon>
                  <span>地点：{{ competition.location }}</span>
                </div>
              </div>
              
              <div class="competition-actions">
                <el-button 
                  type="primary" 
                  size="small"
                  @click="viewCompetitionDetail(competition)"
                >
                  查看详情
                </el-button>
                <el-button 
                  v-if="canRegister(competition)"
                  type="success" 
                  size="small"
                  @click="registerCompetition(competition)"
                  class="register-btn"
                >
                  <el-icon><User /></el-icon>
                  立即报名
                </el-button>
                <el-button 
                  v-else-if="competition.isRegistered"
                  type="warning" 
                  size="small"
                  disabled
                  class="registered-btn"
                >
                  <el-icon><Check /></el-icon>
                  已报名
                </el-button>
                <el-button 
                  v-else-if="!competition.is_open"
                  type="info" 
                  size="small"
                  disabled
                  class="closed-btn"
                >
                  <el-icon><Lock /></el-icon>
                  报名关闭
                </el-button>
                <el-button 
                  v-else-if="getCompetitionStatus(competition) === 'ended'"
                  type="info" 
                  size="small"
                  disabled
                  class="ended-btn"
                >
                  <el-icon><Clock /></el-icon>
                  已结束
                </el-button>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 空状态 -->
    <el-empty 
      v-if="!loading && filteredCompetitions.length === 0" 
      description="暂无竞赛信息"
      style="margin-top: 40px;"
    />

    <!-- 分页 -->
    <div v-if="filteredCompetitions.length > 0" class="pagination-container">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[6, 12, 18, 24]"
        :total="filteredCompetitions.length"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    <!-- 竞赛详情对话框 -->
    <el-dialog
      v-model="showDetailDialog"
      title="竞赛详情"
      width="70%"
      :before-close="handleCloseDetail"
      class="competition-detail-dialog"
    >
      <div v-if="selectedCompetition" class="competition-detail">
        <div class="detail-header">
          <h2>{{ selectedCompetition.title }}</h2>
          <div class="detail-tags">
            <el-tag :type="getStatusType(getCompetitionStatus(selectedCompetition))">
              {{ getStatusText(getCompetitionStatus(selectedCompetition)) }}
            </el-tag>
            <el-tag type="info">{{ getCompetitionType(selectedCompetition) }}</el-tag>
          </div>
        </div>
        
        <el-divider />
        
        <div class="detail-content">
          <div class="detail-section">
            <h3>竞赛简介</h3>
            <p>{{ selectedCompetition.description }}</p>
          </div>
          
          <div class="detail-section">
            <h3>基本信息</h3>
            <el-descriptions :column="2" border>
              <el-descriptions-item label="主办方">{{ selectedCompetition.organizer || '云梦高校' }}</el-descriptions-item>
              <el-descriptions-item label="适用学院">{{ getCollegeName(selectedCompetition.college) }}</el-descriptions-item>
              <el-descriptions-item label="竞赛地点">{{ selectedCompetition.location || '待定' }}</el-descriptions-item>
              <el-descriptions-item label="联系方式">{{ selectedCompetition.contact || '竞赛办公室' }}</el-descriptions-item>
              <el-descriptions-item label="报名人数">{{ getRegisteredCount(selectedCompetition) }}/{{ getMaxParticipants(selectedCompetition) }}</el-descriptions-item>
              <el-descriptions-item label="报名状态">
                <el-tag :type="selectedCompetition.is_open ? 'success' : 'danger'">
                  {{ selectedCompetition.is_open ? '开放报名' : '报名关闭' }}
                </el-tag>
              </el-descriptions-item>
            </el-descriptions>
          </div>
          
          <div class="detail-section">
            <h3>竞赛规则</h3>
            <div v-if="selectedCompetition.rules" v-html="selectedCompetition.rules"></div>
            <div v-else class="no-content">
              <p>暂无竞赛规则信息</p>
            </div>
          </div>
          
          <div class="detail-section">
            <h3>时间安排</h3>
            <el-descriptions :column="1" border>
              <el-descriptions-item label="报名开始时间">
                {{ formatRegistrationStartTime(selectedCompetition) }}
              </el-descriptions-item>
              <el-descriptions-item label="报名截止时间">
                {{ formatRegistrationDeadline(selectedCompetition) }}
                <span v-if="!isRegistrationExpired(selectedCompetition)" class="remaining-time">
                  (剩余{{ formatRemainingTime(getRegistrationDeadline(selectedCompetition)) }})
                </span>
              </el-descriptions-item>
              <el-descriptions-item label="比赛开始时间">
                {{ formatCompetitionStartTime(selectedCompetition) }}
              </el-descriptions-item>
              <el-descriptions-item label="比赛结束时间">
                {{ formatCompetitionEndTime(selectedCompetition) }}
              </el-descriptions-item>
              <el-descriptions-item label="作品提交截止时间" v-if="getSubmissionDeadline(selectedCompetition)">
                {{ formatSubmissionDeadline(selectedCompetition) }}
              </el-descriptions-item>
            </el-descriptions>
          </div>
          
          <div class="detail-section">
            <h3>奖项设置</h3>
            <div v-if="selectedCompetition.awards && selectedCompetition.awards.length > 0">
              <el-table :data="selectedCompetition.awards" border>
                <el-table-column prop="rank" label="奖项" width="120" />
                <el-table-column prop="description" label="说明" />
                <el-table-column prop="count" label="数量" width="80" />
              </el-table>
            </div>
            <div v-else class="no-content">
              <p>暂无奖项设置信息</p>
            </div>
          </div>
          
          <div class="detail-section">
            <h3>报名要求</h3>
            <div v-if="selectedCompetition.requirements" v-html="selectedCompetition.requirements"></div>
            <div v-else class="no-content">
              <p>暂无报名要求信息</p>
            </div>
          </div>
          
          <div class="detail-section">
            <h3>文件要求</h3>
            <el-descriptions :column="1" border>
              <el-descriptions-item label="支持格式">{{ selectedCompetition.fileFormats || 'PDF、PPT、DOC、ZIP' }}</el-descriptions-item>
              <el-descriptions-item label="文件大小限制">{{ selectedCompetition.fileSizeLimit || '单个文件不超过50MB' }}</el-descriptions-item>
              <el-descriptions-item label="提交截止时间" v-if="getSubmissionDeadline(selectedCompetition)">
                {{ formatSubmissionDeadline(selectedCompetition) }}
              </el-descriptions-item>
              <el-descriptions-item label="提交截止时间" v-else>
                待定
              </el-descriptions-item>
            </el-descriptions>
          </div>
          
          <div class="detail-section">
            <h3>评审方式</h3>
            <div v-if="selectedCompetition.judgingMethod" v-html="selectedCompetition.judgingMethod"></div>
            <div v-else class="no-content">
              <p>1. 专家评审</p>
              <p>2. 现场答辩</p>
              <p>3. 综合评分</p>
            </div>
          </div>
          
          <div class="detail-section">
            <h3>相关资源</h3>
            <el-descriptions :column="1" border>
              <el-descriptions-item label="竞赛附件">
                <el-button v-if="selectedCompetition.attachment" type="primary" size="small" @click="downloadAttachment(selectedCompetition)">
                  <el-icon><Download /></el-icon>
                  下载附件
                </el-button>
                <span v-else class="no-content-text">暂无附件</span>
              </el-descriptions-item>
              <el-descriptions-item label="相关链接">
                <el-link v-if="selectedCompetition.website" :href="selectedCompetition.website" target="_blank" type="primary">
                  访问官网
                </el-link>
                <span v-else class="no-content-text">暂无相关链接</span>
              </el-descriptions-item>
              <el-descriptions-item label="QQ群">
                <span v-if="selectedCompetition.qqGroup">{{ selectedCompetition.qqGroup }}</span>
                <span v-else class="no-content-text">暂无QQ群</span>
              </el-descriptions-item>
            </el-descriptions>
          </div>
          
          <div class="detail-section">
            <h3>注意事项</h3>
            <div class="notice-content">
              <el-alert
                v-if="selectedCompetition.importantNotes"
                :title="selectedCompetition.importantNotes"
                type="warning"
                :closable="false"
                show-icon
              />
              <div v-else>
                <p>1. 请仔细阅读竞赛规则和要求</p>
                <p>2. 确保在截止时间前完成报名</p>
                <p>3. 如有疑问请联系相关负责人</p>
                <p>4. 请遵守竞赛纪律和学术诚信</p>
              </div>
            </div>
          </div>
        </div>
        
        <div class="detail-actions">
          <el-button @click="showDetailDialog = false">关闭</el-button>
        </div>
      </div>
    </el-dialog>

    <!-- 报名确认对话框 -->
    <el-dialog
      v-model="showRegisterDialog"
      title="确认报名"
      width="40%"
    >
      <div class="register-confirm">
        <p>您确定要报名参加 <strong>{{ selectedCompetition?.title }}</strong> 吗？</p>
        <p class="warning-text">报名后不可取消，请确认您的参赛资格。</p>
        
        <el-form :model="registerForm" label-width="100px" :rules="registerRules" ref="registerFormRef">
          <el-form-item label="参赛队伍" prop="teamName" required>
            <el-input v-model="registerForm.teamName" placeholder="请输入队伍名称" />
          </el-form-item>
          <el-form-item label="指导老师" prop="advisorId" required>
            <el-select v-model="registerForm.advisorId" placeholder="请选择指导老师" style="width: 100%">
              <el-option 
                v-for="advisor in availableAdvisors" 
                :key="advisor.id" 
                :label="advisor.name" 
                :value="advisor.id"
              >
                <span>{{ advisor.name }}</span>
                <span style="float: right; color: #8492a6; font-size: 13px">{{ advisor.department }}</span>
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="联系电话" prop="contact" required>
            <el-input v-model="registerForm.contact" placeholder="请输入联系电话" />
          </el-form-item>
          <el-form-item label="联系邮箱" prop="email" required>
            <el-input v-model="registerForm.email" placeholder="请输入联系邮箱" />
          </el-form-item>
          <el-form-item label="队伍成员">
            <el-input 
              v-model="registerForm.members" 
              type="textarea" 
              placeholder="请输入队伍成员名单（可选）"
              :rows="3"
            />
          </el-form-item>
          <el-form-item label="备注">
            <el-input 
              v-model="registerForm.remarks" 
              type="textarea" 
              placeholder="请输入备注信息（可选）"
              :rows="3"
            />
          </el-form-item>
        </el-form>
      </div>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showRegisterDialog = false">取消</el-button>
          <el-button type="primary" @click="confirmRegister" :loading="registerLoading">
            确认报名
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Calendar, Clock, User, Download, OfficeBuilding, Location, Check, Lock } from '@element-plus/icons-vue'
import { formatDate, formatDateRange, isExpired, formatRemainingTime } from '@/utils/dateUtils'
import { parsePaginatedResponse } from '@/utils/dataValidator'
import competitionService from '@/services/competitionService'

// 响应式数据
const loading = ref(false)
const filterStatus = ref('')
const filterType = ref('')
const filterCollege = ref('')
const searchKeyword = ref('')
const competitions = ref([])
const showDetailDialog = ref(false)
const showRegisterDialog = ref(false)
const selectedCompetition = ref(null)
const registerLoading = ref(false)

// 分页相关
const currentPage = ref(1)
const pageSize = ref(6)

// 报名表单
const registerForm = ref({
  teamName: '',
  advisorId: '',
  contact: '',
  email: '', // 新增邮箱字段
  members: '',
  remarks: ''
})

// 表单引用
const registerFormRef = ref()

// 表单验证规则
const registerRules = {
  teamName: [{ required: true, message: '请输入队伍名称', trigger: 'blur' }],
  advisorId: [{ required: true, message: '请选择指导老师', trigger: 'change' }],
  contact: [{ required: true, message: '请输入联系电话', trigger: 'blur' }],
  email: [
    { required: true, message: '请输入联系邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
}

// 可用指导老师列表
const availableAdvisors = ref([
  { id: 1, name: '张教授', department: '计算机学院' },
  { id: 2, name: '李副教授', department: '数学学院' },
  { id: 3, name: '王老师', department: '物理学院' },
  { id: 4, name: '陈教授', department: '化学学院' },
  { id: 5, name: '刘老师', department: '工程学院' }
])

// 模拟竞赛数据
const mockCompetitions = [
  {
    id: 1,
    title: '全国大学生程序设计竞赛',
    description: '面向全国大学生的程序设计竞赛，旨在提高学生的编程能力和算法思维',
    type: '程序设计',
    college: 'computer',
    status: 'upcoming',
    organizer: '计算机学院',
    location: '计算机学院实验楼A301',
    contact: '张老师 13800138000',
    registrationStart: '2024-01-01T00:00:00Z',
    registrationDeadline: '2024-02-15T23:59:59Z',
    startDate: '2024-03-01T09:00:00Z',
    endDate: '2024-03-01T18:00:00Z',
    registeredCount: 45,
    maxParticipants: 100,
    isRegistered: false,
    is_open: true,
    rules: '<p>1. 参赛者需为在校大学生</p><p>2. 比赛时长3小时</p><p>3. 禁止使用网络搜索</p>',
    requirements: '<p>1. 具备基本的编程能力</p><p>2. 熟悉C++、Java或Python</p>',
    fileFormats: 'PDF、PPT、DOC、ZIP',
    fileSizeLimit: '单个文件不超过50MB',
    submissionDeadline: '2024-02-28T23:59:59Z',
    judgingMethod: '<p>1. 代码质量评审（40%）</p><p>2. 算法效率测试（30%）</p><p>3. 现场答辩（30%）</p>',
    awards: [
      { rank: '一等奖', description: '奖金5000元', count: 3 },
      { rank: '二等奖', description: '奖金3000元', count: 5 },
      { rank: '三等奖', description: '奖金1000元', count: 10 }
    ],
    attachment: '程序设计竞赛规则.pdf',
    website: 'https://acm.ecnu.edu.cn',
    qqGroup: '123456789',
    importantNotes: '请提前熟悉编程环境，比赛期间禁止使用手机和网络搜索'
  },
  {
    id: 2,
    title: '数学建模竞赛',
    description: '全国大学生数学建模竞赛，培养数学建模和解决实际问题的能力',
    type: '数学建模',
    college: 'mathematics',
    status: 'ongoing',
    organizer: '数学学院',
    location: '数学学院报告厅',
    contact: '李老师 13900139000',
    registrationStart: '2024-01-01T00:00:00Z',
    registrationDeadline: '2024-01-31T23:59:59Z',
    startDate: '2024-02-01T08:00:00Z',
    endDate: '2024-02-03T20:00:00Z',
    registeredCount: 78,
    maxParticipants: 80,
    isRegistered: true,
    is_open: false,
    rules: '<p>1. 三人组队参赛</p><p>2. 比赛时长72小时</p><p> 3. 可使用计算机和网络</p>',
    requirements: '<p>1. 具备数学建模基础</p><p>2. 熟悉MATLAB或Python</p>',
    fileFormats: 'PDF、DOC、ZIP',
    fileSizeLimit: '单个文件不超过100MB',
    submissionDeadline: '2024-02-03T20:00:00Z',
    judgingMethod: '<p>1. 模型建立（30%）</p><p>2. 求解方法（30%）</p><p>3. 结果分析（25%）</p><p>4. 论文质量（15%）</p>',
    awards: [
      { rank: '特等奖', description: '奖金8000元', count: 2 },
      { rank: '一等奖', description: '奖金5000元', count: 5 },
      { rank: '二等奖', description: '奖金3000元', count: 15 }
    ],
    attachment: '数学建模竞赛说明.pdf',
    website: 'https://www.comap.com',
    qqGroup: '987654321',
    importantNotes: '请提前准备MATLAB软件，比赛期间可查阅资料'
  },
  {
    id: 3,
    title: '创新创业大赛',
    description: '大学生创新创业大赛，鼓励学生创新思维和创业实践',
    type: '创新创业',
    college: 'business',
    status: 'ended',
    organizer: '商学院',
    location: '商学院多功能厅',
    contact: '王老师 13700137000',
    registrationStart: '2023-10-01T00:00:00Z',
    registrationDeadline: '2023-11-30T23:59:59Z',
    startDate: '2023-12-01T09:00:00Z',
    endDate: '2023-12-15T18:00:00Z',
    registeredCount: 120,
    maxParticipants: 150,
    isRegistered: false,
    is_open: false,
    rules: '<p>1. 团队参赛，3-5人</p><p>2. 提交商业计划书</p><p>3. 现场答辩展示</p>',
    requirements: '<p>1. 有创新想法或项目</p><p>2. 具备商业分析能力</p>',
    fileFormats: 'PDF、PPT、DOC、ZIP',
    fileSizeLimit: '单个文件不超过50MB',
    submissionDeadline: '2023-12-10T18:00:00Z',
    judgingMethod: '<p>1. 创新性（30%）</p><p>2. 可行性（25%）</p><p>3. 商业价值（25%）</p><p>4. 现场展示（20%）</p>',
    awards: [
      { rank: '金奖', description: '奖金10000元', count: 3 },
      { rank: '银奖', description: '奖金6000元', count: 8 },
      { rank: '铜奖', description: '奖金3000元', count: 15 }
    ],
    attachment: '创新创业大赛说明.pdf',
    website: 'https://www.cyds.com',
    qqGroup: '456789123',
    importantNotes: '请提前准备商业计划书，比赛期间需要现场展示'
  },
  {
    id: 4,
    title: '学术论文竞赛',
    description: '全国大学生学术论文竞赛，提升学术研究能力和论文写作水平',
    type: '学术论文',
    college: 'academic',
    status: 'upcoming',
    organizer: '学术研究处',
    location: '图书馆学术报告厅',
    contact: '陈老师 13600136000',
    registrationStart: '2024-02-01T00:00:00Z',
    registrationDeadline: '2024-03-31T23:59:59Z',
    startDate: '2024-04-01T09:00:00Z',
    endDate: '2024-05-31T18:00:00Z',
    registeredCount: 25,
    maxParticipants: 100,
    isRegistered: false,
    is_open: true,
    rules: '<p>1. 个人或团队参赛</p><p>2. 提交学术论文</p><p>3. 现场答辩</p>',
    requirements: '<p>1. 具备学术研究基础</p><p>2. 熟悉论文写作规范</p>',
    fileFormats: 'PDF、DOC',
    fileSizeLimit: '单个文件不超过20MB',
    submissionDeadline: '2024-05-15T18:00:00Z',
    judgingMethod: '<p>1. 论文质量（40%）</p><p>2. 创新性（30%）</p><p>3. 研究方法（20%）</p><p>4. 答辩表现（10%）</p>',
    awards: [
      { rank: '特等奖', description: '奖金8000元', count: 2 },
      { rank: '一等奖', description: '奖金5000元', count: 5 },
      { rank: '二等奖', description: '奖金3000元', count: 10 },
      { rank: '三等奖', description: '奖金1000元', count: 20 }
    ]
  },
  {
    id: 5,
    title: '机器人设计竞赛',
    description: '大学生机器人设计竞赛，培养工程设计和创新能力',
    type: '工程设计',
    college: 'engineering',
    status: 'ongoing',
    organizer: '工程学院',
    location: '工程学院实验中心',
    contact: '刘老师 13500135000',
    registrationStart: '2024-01-15T00:00:00Z',
    registrationDeadline: '2024-02-28T23:59:59Z',
    startDate: '2024-03-01T08:00:00Z',
    endDate: '2024-03-15T18:00:00Z',
    registeredCount: 35,
    maxParticipants: 50,
    isRegistered: false,
    is_open: true,
    rules: '<p>1. 团队参赛，4-6人</p><p>2. 设计制作机器人</p><p>3. 现场竞技</p>',
    requirements: '<p>1. 具备机械设计基础</p><p>2. 熟悉电子电路</p><p>3. 了解编程控制</p>',
    fileFormats: 'PDF、ZIP、视频文件',
    fileSizeLimit: '单个文件不超过200MB',
    submissionDeadline: '2024-03-10T18:00:00Z',
    judgingMethod: '<p>1. 设计创新性（25%）</p><p>2. 制作工艺（25%）</p><p>3. 功能实现（30%）</p><p>4. 竞技表现（20%）</p>',
    awards: [
      { rank: '冠军', description: '奖金12000元', count: 1 },
      { rank: '亚军', description: '奖金8000元', count: 1 },
      { rank: '季军', description: '奖金5000元', count: 1 },
      { rank: '优秀奖', description: '奖金2000元', count: 5 }
    ]
  },
  {
    id: 6,
    title: '化学实验技能竞赛',
    description: '大学生化学实验技能竞赛，提升实验操作和科学研究能力',
    type: '实验技能',
    college: 'chemistry',
    status: 'upcoming',
    organizer: '化学学院',
    location: '化学学院实验室',
    contact: '赵老师 13400134000',
    registrationStart: '2024-02-15T00:00:00Z',
    registrationDeadline: '2024-03-15T23:59:59Z',
    startDate: '2024-03-20T09:00:00Z',
    endDate: '2024-03-22T18:00:00Z',
    registeredCount: 18,
    maxParticipants: 40,
    isRegistered: false,
    is_open: true,
    rules: '<p>1. 个人参赛</p><p>2. 现场实验操作</p><p>3. 实验报告撰写</p>',
    requirements: '<p>1. 具备化学实验基础</p><p>2. 熟悉实验安全规范</p><p>3. 掌握实验报告写作</p>',
    fileFormats: 'PDF、DOC',
    fileSizeLimit: '单个文件不超过10MB',
    submissionDeadline: '2024-03-21T18:00:00Z',
    judgingMethod: '<p>1. 实验操作（40%）</p><p>2. 实验设计（25%）</p><p>3. 数据分析（20%）</p><p>4. 报告质量（15%）</p>',
    awards: [
      { rank: '一等奖', description: '奖金6000元', count: 3 },
      { rank: '二等奖', description: '奖金4000元', count: 6 },
      { rank: '三等奖', description: '奖金2000元', count: 10 }
    ],
    attachment: '化学实验技能竞赛说明.pdf',
    website: 'https://www.chemistry.edu.cn',
    qqGroup: '789123456',
    importantNotes: '请提前熟悉实验操作流程，注意实验安全'
  },
  {
    id: 7,
    title: '物理创新实验竞赛',
    description: '大学生物理创新实验竞赛，培养物理思维和实验创新能力',
    type: '实验技能',
    college: 'physics',
    status: 'upcoming',
    organizer: '物理学院',
    location: '物理学院实验楼',
    contact: '孙老师 13300133000',
    registrationStart: '2024-03-01T00:00:00Z',
    registrationDeadline: '2024-04-15T23:59:59Z',
    startDate: '2024-04-20T09:00:00Z',
    endDate: '2024-04-25T18:00:00Z',
    registeredCount: 22,
    maxParticipants: 60,
    isRegistered: false,
    is_open: true,
    rules: '<p>1. 个人或团队参赛</p><p>2. 设计创新实验</p><p>3. 现场展示答辩</p>',
    requirements: '<p>1. 具备物理实验基础</p><p>2. 有创新实验想法</p><p>3. 掌握实验设计方法</p>',
    fileFormats: 'PDF、PPT、视频文件',
    fileSizeLimit: '单个文件不超过100MB',
    submissionDeadline: '2024-04-22T18:00:00Z',
    judgingMethod: '<p>1. 实验创新性（35%）</p><p>2. 实验设计（25%）</p><p>3. 实验操作（20%）</p><p>4. 答辩表现（20%）</p>',
    awards: [
      { rank: '特等奖', description: '奖金10000元', count: 2 },
      { rank: '一等奖', description: '奖金6000元', count: 4 },
      { rank: '二等奖', description: '奖金4000元', count: 8 },
      { rank: '三等奖', description: '奖金2000元', count: 15 }
    ],
    attachment: '物理创新实验竞赛说明.pdf',
    website: 'https://www.physics.edu.cn',
    qqGroup: '321654987',
    importantNotes: '请提前准备实验设计方案，注意实验创新性'
  },
  {
    id: 8,
    title: '英语演讲竞赛',
    description: '大学生英语演讲竞赛，提升英语口语表达和跨文化交流能力',
    type: '语言技能',
    college: 'language',
    status: 'ongoing',
    organizer: '语言学院',
    location: '语言学院报告厅',
    contact: '周老师 13200132000',
    registrationStart: '2024-01-20T00:00:00Z',
    registrationDeadline: '2024-02-20T23:59:59Z',
    startDate: '2024-02-25T09:00:00Z',
    endDate: '2024-02-28T18:00:00Z',
    registeredCount: 45,
    maxParticipants: 80,
    isRegistered: false,
    is_open: true,
    rules: '<p>1. 个人参赛</p><p>2. 现场英语演讲</p><p>3. 即兴问答</p>',
    requirements: '<p>1. 具备良好的英语口语</p><p>2. 有演讲经验</p><p>3. 具备应变能力</p>',
    fileFormats: 'PDF、PPT',
    fileSizeLimit: '单个文件不超过20MB',
    submissionDeadline: '2024-02-26T18:00:00Z',
    judgingMethod: '<p>1. 演讲内容（30%）</p><p>2. 语言表达（30%）</p><p>3. 演讲技巧（25%）</p><p>4. 问答表现（15%）</p>',
    awards: [
      { rank: '一等奖', description: '奖金5000元', count: 3 },
      { rank: '二等奖', description: '奖金3000元', count: 6 },
      { rank: '三等奖', description: '奖金1500元', count: 12 }
    ],
    attachment: '英语演讲竞赛说明.pdf',
    website: 'https://www.english.edu.cn',
    qqGroup: '147258369',
    importantNotes: '请提前准备演讲稿，注意英语发音和表达技巧'
  }
]

// 过滤后的竞赛列表
const filteredCompetitions = computed(() => {
  let filtered = competitions.value

  // 状态筛选
  if (filterStatus.value) {
    filtered = filtered.filter(c => c.status === filterStatus.value)
  }

  // 类型筛选
  if (filterType.value) {
    filtered = filtered.filter(c => c.type === filterType.value)
  }

  // 学院筛选
  if (filterCollege.value) {
    filtered = filtered.filter(c => c.college === filterCollege.value)
  }

  // 关键词搜索
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    filtered = filtered.filter(c => 
      c.title.toLowerCase().includes(keyword) || 
      c.description.toLowerCase().includes(keyword)
    )
  }

  return filtered
})

// 分页后的竞赛列表
const paginatedCompetitions = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredCompetitions.value.slice(start, end)
})

// 分页处理函数
const handleSizeChange = (val) => {
  pageSize.value = val
  currentPage.value = 1 // 重置到第一页
}

const handleCurrentChange = (val) => {
  currentPage.value = val
}

// 筛选处理函数
const handleFilterChange = () => {
  // 重置分页到第一页
  currentPage.value = 1
  // 重新加载数据
  loadCompetitions()
}

// 加载竞赛数据
const loadCompetitions = async () => {
  loading.value = true
  try {
    // 构建查询参数
    const params = {
      page: 1,
      size: 50, // 获取更多数据
      search: searchKeyword.value,
      status: filterStatus.value,
      type: filterType.value,
      is_open: true // 只显示开放的竞赛
    }
    
    console.log('加载竞赛数据，参数:', params)
    
    const response = await competitionService.getCompetitions(params)
    console.log('竞赛数据响应:', response)
    
    // 使用parsePaginatedResponse函数解析数据
    if (response && response.code === 200) {
      const parsedData = parsePaginatedResponse(response.data, mockCompetitions)
      console.log('解析后的竞赛数据:', parsedData)
      
      if (parsedData.data.length > 0) {
        competitions.value = parsedData.data
        console.log('使用后端数据，长度:', competitions.value.length)
    } else {
        console.log('后端无数据，使用模拟数据')
        competitions.value = mockCompetitions
      }
    } else {
      console.warn('API返回错误状态:', response?.code, response?.message)
      competitions.value = mockCompetitions
    }
    
    // 如果没有数据，使用模拟数据作为备选
    if (competitions.value.length === 0) {
      console.log('没有真实数据，使用模拟数据')
      competitions.value = mockCompetitions
    }
    
  } catch (error) {
    console.error('加载竞赛数据失败:', error)
    ElMessage.error('加载竞赛数据失败，使用模拟数据')
    // 使用模拟数据作为备选
    competitions.value = mockCompetitions
  } finally {
    loading.value = false
  }
}

// 状态类型映射
const getStatusType = (status) => {
  const statusMap = {
    upcoming: 'warning',
    ongoing: 'success',
    ended: 'info'
  }
  return statusMap[status] || 'info'
}

// 状态文本映射
const getStatusText = (status) => {
  const statusMap = {
    upcoming: '即将开始',
    ongoing: '进行中',
    ended: '已结束'
  }
  return statusMap[status] || status
}

// 学院名称映射
const getCollegeName = (college) => {
  const collegeMap = {
    computer: '计算机学院',
    mathematics: '数学学院',
    physics: '物理学院',
    chemistry: '化学学院',
    engineering: '工程学院',
    business: '商学院',
    academic: '学术研究',
    language: '语言学院'
  }
  return collegeMap[college] || college
}

// 获取竞赛状态
const getCompetitionStatus = (competition) => {
  const now = new Date()
  const start = new Date(competition.start_time)
  const end = new Date(competition.end_time)

  if (now < start) {
    return 'upcoming'
  } else if (now >= start && now <= end) {
    return 'ongoing'
  } else {
    return 'ended'
  }
}

// 获取竞赛类型
const getCompetitionType = (competition) => {
  return competition.type
}

// 使用统一的日期处理工具，移除本地formatDate函数

// 获取报名截止时间（如果没有明确的报名截止时间，使用比赛开始时间的前一天）
const getRegistrationDeadline = (competition) => {
  if (competition.registration_deadline) {
    return competition.registration_deadline
  }
  if (competition.start_time) {
    const startTime = new Date(competition.start_time)
    const deadline = new Date(startTime.getTime() - 24 * 60 * 60 * 1000) // 比赛开始前一天
    return deadline.toISOString()
  }
  return null
}

// 格式化报名截止时间显示
const formatRegistrationDeadline = (competition) => {
  const deadline = getRegistrationDeadline(competition)
  if (deadline) {
    return formatDate(deadline)
  }
  return '未设置'
}

// 检查报名是否已截止
const isRegistrationExpired = (competition) => {
  const deadline = getRegistrationDeadline(competition)
  if (deadline) {
    return isExpired(deadline)
  }
  return false
}

// 格式化比赛时间显示
const formatCompetitionTime = (competition) => {
  if (competition.start_time && competition.end_time) {
    return formatDateRange(competition.start_time, competition.end_time, 'date')
  } else if (competition.start_time) {
    return formatDate(competition.start_time)
  } else if (competition.end_time) {
    return formatDate(competition.end_time)
  }
  return '未设置'
}

// 获取已报名人数
const getRegisteredCount = (competition) => {
  return competition.registration_count || competition.registeredCount || 0
}

// 获取最大参与人数
const getMaxParticipants = (competition) => {
  if (competition.max_participants || competition.maxParticipants) {
    return competition.max_participants || competition.maxParticipants
  }
  return '不限'
}

// 格式化报名开始时间
const formatRegistrationStartTime = (competition) => {
  if (competition.registration_start) {
    return formatDate(competition.registration_start)
  } else if (competition.registrationStart) {
    return formatDate(competition.registrationStart)
  }
  return '未设置'
}

// 格式化比赛开始时间
const formatCompetitionStartTime = (competition) => {
  if (competition.start_time) {
    return formatDate(competition.start_time)
  } else if (competition.startDate) {
    return formatDate(competition.startDate)
  }
  return '未设置'
}

// 格式化比赛结束时间
const formatCompetitionEndTime = (competition) => {
  if (competition.end_time) {
    return formatDate(competition.end_time)
  } else if (competition.endDate) {
    return formatDate(competition.endDate)
  }
  return '未设置'
}

// 获取作品提交截止时间
const getSubmissionDeadline = (competition) => {
  if (competition.submission_deadline) {
    return competition.submission_deadline
  } else if (competition.submissionDeadline) {
    return competition.submissionDeadline
  }
  return null
}

// 格式化作品提交截止时间
const formatSubmissionDeadline = (competition) => {
  const deadline = getSubmissionDeadline(competition)
  if (deadline) {
    return formatDate(deadline)
  }
  return '未设置'
}

// 搜索处理
const handleSearch = () => {
  // 重置分页到第一页
  currentPage.value = 1
  // 实时搜索，computed会自动更新
}

// 查看竞赛详情
const viewCompetitionDetail = (competition) => {
  selectedCompetition.value = competition
  showDetailDialog.value = true
}

// 关闭详情对话框
const handleCloseDetail = () => {
  showDetailDialog.value = false
  selectedCompetition.value = null
}

// 检查是否可以报名
const canRegister = (competition) => {
  const status = getCompetitionStatus(competition)
  const isRegistered = competition.isRegistered || false
  const registeredCount = getRegisteredCount(competition)
  const maxParticipants = competition.max_participants || competition.maxParticipants
  
  return status === 'upcoming' && 
         !isRegistered && 
         competition.is_open &&
         (!maxParticipants || registeredCount < maxParticipants)
}

// 报名竞赛
const registerCompetition = (competition) => {
  selectedCompetition.value = competition
  registerForm.value = {
    teamName: '',
    advisorId: '',
    contact: '',
    email: '', // 新增邮箱字段
    members: '',
    remarks: ''
  }
  showRegisterDialog.value = true
}

// 下载附件
const downloadAttachment = (competition) => {
  if (competition.attachment) {
    // 模拟下载
    ElMessage.success('开始下载附件...')
    // 这里可以添加实际的下载逻辑
  } else {
    ElMessage.warning('暂无附件可下载')
  }
}

// 确认报名
const confirmRegister = async () => {
  // 表单验证
  if (!registerFormRef.value) {
    ElMessage.error('表单引用未找到')
    return
  }
  
  try {
    await registerFormRef.value.validate()
  } catch (error) {
    console.log('表单验证失败:', error)
    return
  }

  registerLoading.value = true
  try {
    console.log('开始报名竞赛:', selectedCompetition.value.title)
    console.log('报名表单数据:', registerForm.value)
    
    // 构建报名数据
    const registrationData = {
      teamName: registerForm.value.teamName,
      teacherId: registerForm.value.advisorId,
      contactPhone: registerForm.value.contact,
      contactEmail: registerForm.value.email,
      teamLeader: true, // 默认为队长
      additionalInfo: {
        members: registerForm.value.members,
        remarks: registerForm.value.remarks
      }
    }
    
    console.log('发送报名数据:', registrationData)
    
    // 调用后端报名API
    const response = await competitionService.registerCompetition(
      selectedCompetition.value.id, 
      registrationData
    )
    
    console.log('报名API响应:', response)
    
    if (response && response.code === 200) {
      // 报名成功，更新本地数据
    const index = competitions.value.findIndex(c => c.id === selectedCompetition.value.id)
    if (index !== -1) {
        // 更新报名状态和人数
      competitions.value[index].isRegistered = true
        competitions.value[index].registration_count = (competitions.value[index].registration_count || 0) + 1
        competitions.value[index].current_participants = (competitions.value[index].current_participants || 0) + 1
        
        // 如果后端返回了更新后的数据，使用后端数据
        if (response.data) {
          console.log('使用后端返回的更新数据')
          // 可以在这里更新更多字段
        }
    }
    
    showRegisterDialog.value = false
    ElMessage.success('报名成功！')
      
      // 重新加载竞赛数据以确保数据同步
      await loadCompetitions()
      
      // 重置表单
      registerForm.value = {
        teamName: '',
        advisorId: '',
        contact: '',
        email: '',
        members: '',
        remarks: ''
      }
      
      // 重置表单验证状态
      if (registerFormRef.value) {
        registerFormRef.value.resetFields()
      }
    } else {
      // API返回错误
      const errorMsg = response?.message || '报名失败，请稍后重试'
      ElMessage.error(errorMsg)
      console.error('报名API返回错误:', response)
    }
    
  } catch (error) {
    console.error('报名失败:', error)
    
    // 根据错误类型显示不同的错误信息
    let errorMsg = '报名失败，请稍后重试'
    
    if (error.response) {
      const status = error.response.status
      const data = error.response.data
      
      switch (status) {
        case 400:
          errorMsg = data?.message || '报名信息有误，请检查后重试'
          break
        case 401:
          errorMsg = '登录已过期，请重新登录'
          break
        case 403:
          errorMsg = '权限不足，无法报名此竞赛'
          break
        case 404:
          errorMsg = '竞赛不存在或已删除'
          break
        case 409:
          errorMsg = '您已经报名过此竞赛'
          break
        default:
          errorMsg = data?.message || `报名失败 (${status})`
      }
    } else if (error.request) {
      errorMsg = '网络连接失败，请检查网络后重试'
    }
    
    ElMessage.error(errorMsg)
  } finally {
    registerLoading.value = false
  }
}

// 组件挂载时加载数据
onMounted(() => {
  loadCompetitions()
})
</script>

<style scoped>
.competition-view {
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

.filter-section {
  margin-bottom: 30px;
  padding: 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.competition-card {
  height: auto;
  min-height: 320px;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  cursor: pointer;
  display: flex;
  flex-direction: column;
}

.competition-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
}

.competition-actions {
  padding: 15px;
  display: flex;
  gap: 10px;
  justify-content: center;
  border-top: 1px solid #f0f0f0;
}

.register-btn {
  background: linear-gradient(135deg, #67c23a 0%, #85ce61 100%);
  border: none;
  color: white;
  font-weight: 600;
  transition: all 0.3s ease;
}

.register-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(103, 194, 58, 0.4);
}

.registered-btn {
  background: linear-gradient(135deg, #e6a23c 0%, #f0c78a 100%);
  border: none;
  color: white;
  font-weight: 600;
}

.closed-btn {
  background: linear-gradient(135deg, #909399 0%, #c0c4cc 100%);
  border: none;
  color: white;
  font-weight: 600;
}

.ended-btn {
  background: linear-gradient(135deg, #f56c6c 0%, #f78989 100%);
  border: none;
  color: white;
  font-weight: 600;
}

.competition-info {
  margin: 15px 0;
}

.info-item {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
  font-size: 14px;
  color: #606266;
}

.info-item .el-icon {
  margin-right: 8px;
  color: #909399;
  width: 16px;
}

.remaining-time {
  color: #e6a23c;
  font-size: 12px;
  margin-left: 8px;
  font-weight: 600;
}

.competition-header {
  padding: 15px 15px 0 15px;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.competition-content {
  padding: 15px;
  flex: 1;
  display: flex;
  flex-direction: column;
}

.competition-title {
  margin: 0 0 10px 0;
  font-size: 18px;
  font-weight: 600;
  color: #2c3e50;
  line-height: 1.4;
}

.competition-meta {
  display: flex;
  gap: 15px;
  margin-bottom: 15px;
  flex-wrap: wrap;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 5px;
  color: #606266;
  font-size: 13px;
}

.meta-item .el-icon {
  color: #909399;
  font-size: 14px;
}

.competition-description {
  margin: 0 0 15px 0;
  color: #7f8c8d;
  line-height: 1.6;
  /* 移除行数限制，允许完整显示 */
  max-height: none;
  overflow: visible;
}

.competition-actions {
  display: flex;
  gap: 10px;
  margin-top: auto;
  padding-top: 15px;
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}

.detail-header h2 {
  margin: 0;
  color: #2c3e50;
}

.detail-tags {
  display: flex;
  gap: 10px;
}

.detail-section {
  margin-bottom: 25px;
}

.detail-section h3 {
  margin: 0 0 15px 0;
  color: #2c3e50;
  font-size: 18px;
  font-weight: 600;
}

.detail-section p {
  margin: 0;
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

.register-confirm {
  padding: 20px 0;
}

.warning-text {
  color: #e6a23c;
  font-size: 14px;
  margin: 10px 0 20px 0;
}

.notice-content {
  background: #f8f9fa;
  padding: 15px;
  border-radius: 6px;
  border-left: 4px solid #e6a23c;
}

.notice-content p {
  margin: 8px 0;
  color: #5a6c7d;
  line-height: 1.6;
}

.no-content {
  background: #f8f9fa;
  padding: 15px;
  border-radius: 6px;
  border-left: 4px solid #e4e7ed;
  color: #909399;
  font-style: italic;
}

.no-content p {
  margin: 5px 0;
  color: #909399;
}

.no-content-text {
  color: #909399;
  font-style: italic;
}

.detail-section .el-descriptions {
  margin-top: 10px;
}

.detail-section .el-button {
  margin-right: 10px;
}

.detail-section .el-link {
  font-size: 14px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

/* 竞赛详情对话框样式优化 */
.competition-detail-dialog .el-dialog__body {
  max-height: 70vh;
  overflow-y: auto;
}

.competition-detail-dialog .el-dialog__header {
  border-bottom: 1px solid #e9ecef;
  padding-bottom: 15px;
}

.competition-detail-dialog .el-dialog__title {
  font-size: 20px;
  font-weight: 600;
  color: #2c3e50;
}

/* 分页样式 */
.pagination-container {
  margin-top: 30px;
  text-align: center;
  padding: 20px 0;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

/* 竞赛卡片响应式布局 */
@media (max-width: 1200px) {
  .competition-list .el-col {
    width: 50% !important;
  }
  
  .competition-card {
    min-height: 350px;
  }
}

@media (max-width: 768px) {
  .competition-list .el-col {
    width: 100% !important;
  }
  
  .filter-section .el-col {
    margin-bottom: 10px;
  }
  
  .competition-actions {
    flex-direction: column;
  }
  
  .competition-actions .el-button {
    width: 100%;
  }
  
  .competition-card {
    min-height: auto;
  }
  
  .competition-meta {
    flex-direction: column;
    gap: 8px;
  }
  
  .competition-detail-dialog {
    width: 95% !important;
  }
}
</style> 