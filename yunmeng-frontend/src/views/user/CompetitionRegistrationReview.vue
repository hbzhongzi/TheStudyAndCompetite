<template>
  <div class="competition-registration-review">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2>竞赛报名审核</h2>
      <p>审核和管理学生的竞赛报名申请</p>
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
              <h4>待审核</h4>
              <p class="stat-number">{{ stats.pendingCount }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon approved">
              <el-icon><Check /></el-icon>
            </div>
            <div class="stat-info">
              <h4>已通过</h4>
              <p class="stat-number">{{ stats.approvedCount }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon rejected">
              <el-icon><Close /></el-icon>
            </div>
            <div class="stat-info">
              <h4>已拒绝</h4>
              <p class="stat-number">{{ stats.rejectedCount }}</p>
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
              <h4>总报名</h4>
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
            <el-option label="已拒绝" value="rejected" />
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
          <el-button @click="exportRegistrations">
            <el-icon><Download /></el-icon>
            导出数据
          </el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 报名列表 -->
    <el-card class="registration-list-card">
      <el-table 
        :data="registrations" 
        style="width: 100%" 
        v-loading="loading"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="报名ID" width="80" />
        <el-table-column label="竞赛信息" width="250">
          <template #default="scope">
            <div class="competition-info">
              <h4>{{ scope.row.competition?.title }}</h4>
              <p class="competition-type">{{ scope.row.competition?.type }}</p>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="学生信息" width="200">
          <template #default="scope">
            <div class="student-info">
              <h4>{{ scope.row.student?.realName }}</h4>
              <p>{{ scope.row.student?.department }} - {{ scope.row.student?.studentId }}</p>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="团队信息" width="150">
          <template #default="scope">
            <div class="team-info">
              <p v-if="scope.row.teamName">{{ scope.row.teamName }}</p>
              <p v-else class="no-team">个人参赛</p>
              <p v-if="scope.row.teamLeader" class="team-leader">队长</p>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="指导教师" width="120">
          <template #default="scope">
            <span v-if="scope.row.teacher">{{ scope.row.teacher.realName }}</span>
            <span v-else class="no-teacher">未绑定</span>
          </template>
        </el-table-column>
        <el-table-column prop="registerTime" label="报名时间" width="150">
          <template #default="scope">
            {{ formatDate(scope.row.registerTime) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="scope">
            <div class="action-buttons">
              <el-button 
                v-if="scope.row.status === 'pending'"
                size="small" 
                type="success" 
                @click="approveRegistration(scope.row)"
              >
                通过
              </el-button>
              <el-button 
                v-if="scope.row.status === 'pending'"
                size="small" 
                type="danger" 
                @click="rejectRegistration(scope.row)"
              >
                拒绝
              </el-button>
              <el-button 
                size="small" 
                type="info" 
                @click="viewDetails(scope.row)"
              >
                详情
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

    <!-- 批量操作 -->
    <el-card v-if="selectedRegistrations.length > 0" class="batch-actions-card">
      <template #header>
        <span>批量操作 (已选择 {{ selectedRegistrations.length }} 项)</span>
      </template>
      <div class="batch-actions">
        <el-button 
          type="success" 
          @click="batchApprove"
          :loading="batchProcessing"
        >
          批量通过
        </el-button>
        <el-button 
          type="danger" 
          @click="batchReject"
          :loading="batchProcessing"
        >
          批量拒绝
        </el-button>
        <el-button @click="clearSelection">清除选择</el-button>
      </div>
    </el-card>

    <!-- 报名详情对话框 -->
    <el-dialog
      v-model="showDetailDialog"
      title="报名详情"
      width="60%"
    >
      <div v-if="selectedRegistration" class="registration-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="竞赛名称">
            {{ selectedRegistration.competition?.title }}
          </el-descriptions-item>
          <el-descriptions-item label="竞赛类型">
            {{ selectedRegistration.competition?.type }}
          </el-descriptions-item>
          <el-descriptions-item label="学生姓名">
            {{ selectedRegistration.student?.realName }}
          </el-descriptions-item>
          <el-descriptions-item label="学号">
            {{ selectedRegistration.student?.studentId }}
          </el-descriptions-item>
          <el-descriptions-item label="学院">
            {{ selectedRegistration.student?.department }}
          </el-descriptions-item>
          <el-descriptions-item label="联系电话">
            {{ selectedRegistration.contactPhone || '未提供' }}
          </el-descriptions-item>
          <el-descriptions-item label="联系邮箱">
            {{ selectedRegistration.contactEmail || '未提供' }}
          </el-descriptions-item>
          <el-descriptions-item label="团队名称">
            {{ selectedRegistration.teamName || '个人参赛' }}
          </el-descriptions-item>
          <el-descriptions-item label="指导教师">
            {{ selectedRegistration.teacher?.realName || '未绑定' }}
          </el-descriptions-item>
          <el-descriptions-item label="报名时间">
            {{ formatDate(selectedRegistration.registerTime) }}
          </el-descriptions-item>
          <el-descriptions-item label="当前状态">
            <el-tag :type="getStatusType(selectedRegistration.status)">
              {{ getStatusText(selectedRegistration.status) }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>

        <div v-if="selectedRegistration.additionalInfo" class="additional-info">
          <h4>额外信息</h4>
          <pre>{{ JSON.stringify(selectedRegistration.additionalInfo, null, 2) }}</pre>
        </div>
      </div>
    </el-dialog>

    <!-- 拒绝原因对话框 -->
    <el-dialog
      v-model="showRejectDialog"
      title="拒绝报名"
      width="40%"
    >
      <el-form :model="rejectForm" label-width="100px">
        <el-form-item label="拒绝原因" required>
          <el-input
            v-model="rejectForm.reason"
            type="textarea"
            :rows="4"
            placeholder="请输入拒绝原因"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showRejectDialog = false">取消</el-button>
          <el-button type="danger" @click="confirmReject" :loading="rejecting">
            确认拒绝
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Clock, Check, Close, DataAnalysis, Search, Download 
} from '@element-plus/icons-vue'
import { formatDate } from '@/utils/dateUtils'

// 响应式数据
const loading = ref(false)
const registrations = ref([])
const competitions = ref([])
const selectedRegistration = ref(null)
const selectedRegistrations = ref([])
const showDetailDialog = ref(false)
const showRejectDialog = ref(false)
const rejecting = ref(false)
const batchProcessing = ref(false)

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

// 拒绝表单
const rejectForm = ref({
  reason: ''
})

// 统计数据
const stats = ref({
  pendingCount: 0,
  approvedCount: 0,
  rejectedCount: 0,
  totalCount: 0
})

// 计算属性
const filteredRegistrations = computed(() => {
  let filtered = registrations.value

  if (filters.value.search) {
    const search = filters.value.search.toLowerCase()
    filtered = filtered.filter(r => 
      r.competition?.title?.toLowerCase().includes(search) ||
      r.student?.realName?.toLowerCase().includes(search) ||
      r.student?.studentId?.includes(search)
    )
  }

  if (filters.value.status) {
    filtered = filtered.filter(r => r.status === filters.value.status)
  }

  if (filters.value.competitionId) {
    filtered = filtered.filter(r => r.competitionId === filters.value.competitionId)
  }

  if (filters.value.department) {
    filtered = filtered.filter(r => r.student?.department === filters.value.department)
  }

  return filtered
})

// 分页后的数据
const paginatedRegistrations = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredRegistrations.value.slice(start, end)
})

// 方法
const loadRegistrations = async () => {
  loading.value = true
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // 模拟数据
    registrations.value = [
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
        teacher: {
          id: 1,
          realName: '李教授'
        },
        teamName: '代码王者队',
        teamLeader: true,
        contactPhone: '13800138000',
        contactEmail: 'zhangsan@example.com',
        registerTime: new Date('2024-01-15'),
        status: 'pending',
        additionalInfo: {
          programmingLanguages: ['C++', 'Python'],
          experience: '参加过校级比赛'
        }
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
        teacher: null,
        teamName: '',
        teamLeader: false,
        contactPhone: '13800138001',
        contactEmail: 'lisi@example.com',
        registerTime: new Date('2024-01-16'),
        status: 'approved',
        additionalInfo: null
      }
    ]
    
    loadStats()
  } catch (error) {
    console.error('加载报名数据失败:', error)
    ElMessage.error('加载报名数据失败')
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
  const pending = registrations.value.filter(r => r.status === 'pending').length
  const approved = registrations.value.filter(r => r.status === 'approved').length
  const rejected = registrations.value.filter(r => r.status === 'rejected').length
  
  stats.value = {
    pendingCount: pending,
    approvedCount: approved,
    rejectedCount: rejected,
    totalCount: registrations.value.length
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
  selectedRegistrations.value = selection
}

const clearSelection = () => {
  selectedRegistrations.value = []
}

const viewDetails = (registration) => {
  selectedRegistration.value = registration
  showDetailDialog.value = true
}

const approveRegistration = async (registration) => {
  try {
    await ElMessageBox.confirm(
      `确定要通过学生 "${registration.student.realName}" 的报名申请吗？`,
      '确认操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500))
    
    registration.status = 'approved'
    loadStats()
    ElMessage.success('报名审核通过成功')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败：' + (error.message || '未知错误'))
    }
  }
}

const rejectRegistration = (registration) => {
  selectedRegistration.value = registration
  rejectForm.value.reason = ''
  showRejectDialog.value = true
}

const confirmReject = async () => {
  if (!rejectForm.value.reason.trim()) {
    ElMessage.warning('请输入拒绝原因')
    return
  }
  
  rejecting.value = true
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500))
    
    selectedRegistration.value.status = 'rejected'
    selectedRegistration.value.rejectionReason = rejectForm.value.reason
    loadStats()
    
    showRejectDialog.value = false
    ElMessage.success('报名已拒绝')
  } catch (error) {
    ElMessage.error('操作失败：' + (error.message || '未知错误'))
  } finally {
    rejecting.value = false
  }
}

const batchApprove = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要批量通过 ${selectedRegistrations.value.length} 个报名申请吗？`,
      '确认批量操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    batchProcessing.value = true
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    selectedRegistrations.value.forEach(r => {
      r.status = 'approved'
    })
    
    loadStats()
    selectedRegistrations.value = []
    ElMessage.success('批量审核通过成功')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量操作失败：' + (error.message || '未知错误'))
    }
  } finally {
    batchProcessing.value = false
  }
}

const batchReject = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要批量拒绝 ${selectedRegistrations.value.length} 个报名申请吗？`,
      '确认批量操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    batchProcessing.value = true
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    selectedRegistrations.value.forEach(r => {
      r.status = 'rejected'
      r.rejectionReason = '批量拒绝'
    })
    
    loadStats()
    selectedRegistrations.value = []
    ElMessage.success('批量拒绝成功')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量操作失败：' + (error.message || '未知错误'))
    }
  } finally {
    batchProcessing.value = false
  }
}

const exportRegistrations = () => {
  ElMessage.success('开始导出报名数据...')
  // 实际项目中这里会调用导出API
}

// 状态相关方法
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
    rejected: '已拒绝'
  }
  return statusMap[status] || status
}

// 组件挂载时加载数据
onMounted(() => {
  loadRegistrations()
  loadCompetitions()
})
</script>

<style scoped>
.competition-registration-review {
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

.stat-icon.approved {
  background: #67c23a;
}

.stat-icon.rejected {
  background: #f56c6c;
}

.stat-icon.total {
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

.registration-list-card {
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

.team-info p {
  margin: 5px 0;
  color: #606266;
  font-size: 12px;
}

.team-leader {
  color: #409eff !important;
  font-weight: 600;
}

.no-team, .no-teacher {
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

.batch-actions-card {
  margin-top: 20px;
}

.batch-actions {
  display: flex;
  gap: 15px;
}

.additional-info {
  margin-top: 20px;
  padding: 15px;
  background: #f8f9fa;
  border-radius: 6px;
}

.additional-info h4 {
  margin: 0 0 10px 0;
  color: #2c3e50;
}

.additional-info pre {
  margin: 0;
  color: #606266;
  font-size: 12px;
  white-space: pre-wrap;
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
  
  .batch-actions {
    flex-direction: column;
  }
}
</style> 