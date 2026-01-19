<template>
  <div class="competition-audit-logs">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2>竞赛审计日志</h2>
      <p>查看和管理所有竞赛相关操作的审计记录</p>
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
              <h4>总日志数</h4>
              <p class="stat-number">{{ stats.totalLogs }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon today">
              <el-icon><Calendar /></el-icon>
            </div>
            <div class="stat-info">
              <h4>今日日志</h4>
              <p class="stat-number">{{ stats.todayLogs }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon critical">
              <el-icon><Warning /></el-icon>
            </div>
            <div class="stat-info">
              <h4>重要操作</h4>
              <p class="stat-number">{{ stats.criticalLogs }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon users">
              <el-icon><User /></el-icon>
            </div>
            <div class="stat-info">
              <h4>活跃用户</h4>
              <p class="stat-number">{{ stats.activeUsers }}</p>
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
            placeholder="搜索操作内容、用户姓名"
            clearable
            @input="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-col>
        <el-col :span="4">
          <el-select v-model="filters.action" placeholder="操作类型" clearable @change="handleSearch">
            <el-option label="全部操作" value="" />
            <el-option label="创建竞赛" value="competition_created" />
            <el-option label="修改竞赛" value="competition_updated" />
            <el-option label="删除竞赛" value="competition_deleted" />
            <el-option label="报名审核" value="registration_reviewed" />
            <el-option label="作品提交" value="submission_created" />
            <el-option label="评审打分" value="judging_scored" />
            <el-option label="结果发布" value="results_published" />
            <el-option label="用户登录" value="user_login" />
            <el-option label="权限变更" value="permission_changed" />
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
          <el-select v-model="filters.userRole" placeholder="用户角色" clearable @change="handleSearch">
            <el-option label="全部角色" value="" />
            <el-option label="管理员" value="admin" />
            <el-option label="教师" value="teacher" />
            <el-option label="学生" value="student" />
          </el-select>
        </el-col>
        <el-col :span="6">
          <el-date-picker
            v-model="filters.dateRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
            @change="handleSearch"
          />
        </el-col>
      </el-row>
      
      <el-row :gutter="20" style="margin-top: 15px;">
        <el-col :span="6">
          <el-select v-model="filters.severity" placeholder="日志级别" clearable @change="handleSearch">
            <el-option label="全部级别" value="" />
            <el-option label="信息" value="info" />
            <el-option label="警告" value="warning" />
            <el-option label="错误" value="error" />
            <el-option label="严重" value="critical" />
          </el-select>
        </el-col>
        <el-col :span="6">
          <el-input
            v-model="filters.ipAddress"
            placeholder="IP地址"
            clearable
            @input="handleSearch"
          />
        </el-col>
        <el-col :span="12">
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="exportLogs">
            <el-icon><Download /></el-icon>
            导出日志
          </el-button>
          <el-button @click="clearOldLogs" type="warning">
            <el-icon><Delete /></el-icon>
            清理旧日志
          </el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 日志列表 -->
    <el-card class="logs-list-card">
      <el-table 
        :data="auditLogs" 
        style="width: 100%" 
        v-loading="loading"
        @selection-change="handleSelectionChange"
        :row-class-name="getRowClassName"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="日志ID" width="80" />
        <el-table-column label="时间" width="160">
          <template #default="scope">
            <div class="time-info">
              <div class="time">{{ formatDateTime(scope.row.createdAt) }}</div>
              <div class="time-ago">{{ getTimeAgo(scope.row.createdAt) }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="用户信息" width="180">
          <template #default="scope">
            <div class="user-info">
              <h4>{{ scope.row.user?.realName || scope.row.user?.username }}</h4>
              <p class="user-meta">
                <el-tag size="small" :type="getRoleType(scope.row.user?.role)">
                  {{ getRoleText(scope.row.user?.role) }}
                </el-tag>
                <span class="ip-address">{{ scope.row.ipAddress }}</span>
              </p>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="操作信息" width="200">
          <template #default="scope">
            <div class="action-info">
              <h4>{{ getActionText(scope.row.action) }}</h4>
              <p class="action-details">{{ scope.row.details }}</p>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="竞赛信息" width="200">
          <template #default="scope">
            <div v-if="scope.row.competition" class="competition-info">
              <h4>{{ scope.row.competition.title }}</h4>
              <p class="competition-type">{{ scope.row.competition.type }}</p>
            </div>
            <span v-else class="no-competition">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="severity" label="级别" width="80">
          <template #default="scope">
            <el-tag :type="getSeverityType(scope.row.severity)" size="small">
              {{ getSeverityText(scope.row.severity) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="scope">
            <div class="action-buttons">
              <el-button 
                size="small" 
                type="primary" 
                @click="viewLogDetails(scope.row)"
              >
                详情
              </el-button>
              <el-button 
                v-if="scope.row.severity === 'critical' || scope.row.severity === 'error'"
                size="small" 
                type="warning" 
                @click="markAsResolved(scope.row)"
              >
                标记已处理
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
          :page-sizes="[20, 50, 100, 200]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 批量操作 -->
    <el-card v-if="selectedLogs.length > 0" class="batch-actions-card">
      <template #header>
        <span>批量操作 (已选择 {{ selectedLogs.length }} 项)</span>
      </template>
      <div class="batch-actions">
        <el-button 
          type="warning" 
          @click="batchMarkResolved"
          :loading="batchProcessing"
        >
          批量标记已处理
        </el-button>
        <el-button 
          type="danger" 
          @click="batchDelete"
          :loading="batchProcessing"
        >
          批量删除
        </el-button>
        <el-button @click="clearSelection">清除选择</el-button>
      </div>
    </el-card>

    <!-- 日志详情对话框 -->
    <el-dialog
      v-model="showDetailDialog"
      title="审计日志详情"
      width="70%"
    >
      <div v-if="selectedLog" class="log-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="日志ID">
            {{ selectedLog.id }}
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">
            {{ formatDateTime(selectedLog.createdAt) }}
          </el-descriptions-item>
          <el-descriptions-item label="操作用户">
            {{ selectedLog.user?.realName || selectedLog.user?.username }}
          </el-descriptions-item>
          <el-descriptions-item label="用户角色">
            <el-tag :type="getRoleType(selectedLog.user?.role)">
              {{ getRoleText(selectedLog.user?.role) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="操作类型">
            {{ getActionText(selectedLog.action) }}
          </el-descriptions-item>
          <el-descriptions-item label="日志级别">
            <el-tag :type="getSeverityType(selectedLog.severity)">
              {{ getSeverityText(selectedLog.severity) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="IP地址">
            {{ selectedLog.ipAddress }}
          </el-descriptions-item>
          <el-descriptions-item label="用户代理">
            {{ selectedLog.userAgent || '未记录' }}
          </el-descriptions-item>
        </el-descriptions>
        
        <div class="details-section">
          <h4>操作详情</h4>
          <p>{{ selectedLog.details }}</p>
        </div>
        
        <div v-if="selectedLog.competition" class="competition-section">
          <h4>相关竞赛</h4>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="竞赛名称">
              {{ selectedLog.competition.title }}
            </el-descriptions-item>
            <el-descriptions-item label="竞赛类型">
              {{ selectedLog.competition.type }}
            </el-descriptions-item>
            <el-descriptions-item label="竞赛级别">
              {{ selectedLog.competition.level }}
            </el-descriptions-item>
            <el-descriptions-item label="主办方">
              {{ selectedLog.competition.organizer }}
            </el-descriptions-item>
          </el-descriptions>
        </div>
        
        <div v-if="selectedLog.beforeData || selectedLog.afterData" class="data-changes-section">
          <h4>数据变化</h4>
          <el-row :gutter="20">
            <el-col :span="12">
              <h5>变更前</h5>
              <pre v-if="selectedLog.beforeData">{{ JSON.stringify(selectedLog.beforeData, null, 2) }}</pre>
              <p v-else class="no-data">无数据</p>
            </el-col>
            <el-col :span="12">
              <h5>变更后</h5>
              <pre v-if="selectedLog.afterData">{{ JSON.stringify(selectedLog.afterData, null, 2) }}</pre>
              <p v-else class="no-data">无数据</p>
            </el-col>
          </el-row>
        </div>
      </div>
    </el-dialog>

    <!-- 清理旧日志对话框 -->
    <el-dialog
      v-model="showClearDialog"
      title="清理旧日志"
      width="40%"
    >
      <el-form :model="clearForm" label-width="120px">
        <el-form-item label="保留天数" required>
          <el-input-number
            v-model="clearForm.daysToKeep"
            :min="1"
            :max="365"
            placeholder="输入天数"
          />
        </el-form-item>
        <el-form-item label="日志级别">
          <el-select v-model="clearForm.severity" placeholder="选择级别" multiple>
            <el-option label="信息" value="info" />
            <el-option label="警告" value="warning" />
            <el-option label="错误" value="error" />
            <el-option label="严重" value="critical" />
          </el-select>
        </el-form-item>
        <el-form-item label="确认清理">
          <el-checkbox v-model="clearForm.confirm">
            我确认要清理选中的旧日志，此操作不可恢复
          </el-checkbox>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showClearDialog = false">取消</el-button>
          <el-button 
            type="warning" 
            @click="confirmClearLogs" 
            :loading="clearing"
            :disabled="!clearForm.confirm"
          >
            确认清理
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
  DataAnalysis, Calendar, Warning, User, Search, Download, Delete 
} from '@element-plus/icons-vue'
import { formatDate } from '@/utils/dateUtils'

// 响应式数据
const loading = ref(false)
const auditLogs = ref([])
const competitions = ref([])
const selectedLog = ref(null)
const showDetailDialog = ref(false)
const showClearDialog = ref(false)
const clearing = ref(false)
const batchProcessing = ref(false)
const selectedLogs = ref([])

// 分页相关
const currentPage = ref(1)
const pageSize = ref(50)
const total = ref(0)

// 筛选条件
const filters = ref({
  search: '',
  action: '',
  competitionId: '',
  userRole: '',
  dateRange: [],
  severity: '',
  ipAddress: ''
})

// 清理表单
const clearForm = ref({
  daysToKeep: 30,
  severity: [],
  confirm: false
})

// 统计数据
const stats = ref({
  totalLogs: 0,
  todayLogs: 0,
  criticalLogs: 0,
  activeUsers: 0
})

// 计算属性
const filteredLogs = computed(() => {
  let filtered = auditLogs.value

  if (filters.value.search) {
    const search = filters.value.search.toLowerCase()
    filtered = filtered.filter(log => 
      log.details?.toLowerCase().includes(search) ||
      log.user?.realName?.toLowerCase().includes(search) ||
      log.user?.username?.toLowerCase().includes(search)
    )
  }

  if (filters.value.action) {
    filtered = filtered.filter(log => log.action === filters.value.action)
  }

  if (filters.value.competitionId) {
    filtered = filtered.filter(log => log.competitionId === filters.value.competitionId)
  }

  if (filters.value.userRole) {
    filtered = filtered.filter(log => log.user?.role === filters.value.userRole)
  }

  if (filters.value.severity) {
    filtered = filtered.filter(log => log.severity === filters.value.severity)
  }

  if (filters.value.ipAddress) {
    filtered = filtered.filter(log => log.ipAddress?.includes(filters.value.ipAddress))
  }

  if (filters.value.dateRange && filters.value.dateRange.length === 2) {
    filtered = filtered.filter(log => {
      const logDate = new Date(log.createdAt)
      const startDate = new Date(filters.value.dateRange[0])
      const endDate = new Date(filters.value.dateRange[1])
      return logDate >= startDate && logDate <= endDate
    })
  }

  return filtered
})

// 分页后的数据
const paginatedLogs = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredLogs.value.slice(start, end)
})

// 方法
const loadAuditLogs = async () => {
  loading.value = true
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // 模拟数据
    auditLogs.value = [
      {
        id: 1,
        userId: 1,
        user: {
          id: 1,
          username: 'admin',
          realName: '系统管理员',
          role: 'admin'
        },
        action: 'competition_created',
        details: '创建竞赛：全国大学生程序设计竞赛',
        competitionId: 1,
        competition: {
          id: 1,
          title: '全国大学生程序设计竞赛',
          type: '程序设计',
          level: 'national',
          organizer: '计算机学院'
        },
        severity: 'info',
        ipAddress: '192.168.1.100',
        userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36',
        createdAt: new Date('2024-02-15 10:30:00'),
        beforeData: null,
        afterData: {
          title: '全国大学生程序设计竞赛',
          type: '程序设计',
          level: 'national'
        }
      },
      {
        id: 2,
        userId: 2,
        user: {
          id: 2,
          username: 'teacher001',
          realName: '李教授',
          role: 'teacher'
        },
        action: 'registration_reviewed',
        details: '审核学生张三的报名申请：通过',
        competitionId: 1,
        competition: {
          id: 1,
          title: '全国大学生程序设计竞赛',
          type: '程序设计',
          level: 'national',
          organizer: '计算机学院'
        },
        severity: 'info',
        ipAddress: '192.168.1.101',
        userAgent: 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36',
        createdAt: new Date('2024-02-16 14:20:00'),
        beforeData: { status: 'pending' },
        afterData: { status: 'approved' }
      },
      {
        id: 3,
        userId: 3,
        user: {
          id: 3,
          username: 'student001',
          realName: '张三',
          role: 'student'
        },
        action: 'submission_created',
        details: '提交作品：基于深度学习的图像识别系统',
        competitionId: 1,
        competition: {
          id: 1,
          title: '全国大学生程序设计竞赛',
          type: '程序设计',
          level: 'national',
          organizer: '计算机学院'
        },
        severity: 'info',
        ipAddress: '192.168.1.102',
        userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36',
        createdAt: new Date('2024-02-17 09:15:00'),
        beforeData: null,
        afterData: {
          title: '基于深度学习的图像识别系统',
          fileSize: '15MB'
        }
      }
    ]
    
    loadStats()
  } catch (error) {
    console.error('加载审计日志失败:', error)
    ElMessage.error('加载审计日志失败')
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
  const total = auditLogs.value.length
  const today = new Date()
  const todayLogs = auditLogs.value.filter(log => {
    const logDate = new Date(log.createdAt)
    return logDate.toDateString() === today.toDateString()
  }).length
  const critical = auditLogs.value.filter(log => 
    log.severity === 'critical' || log.severity === 'error'
  ).length
  const uniqueUsers = new Set(auditLogs.value.map(log => log.userId)).size
  
  stats.value = {
    totalLogs: total,
    todayLogs: todayLogs,
    criticalLogs: critical,
    activeUsers: uniqueUsers
  }
}

const handleSearch = () => {
  currentPage.value = 1
  // 实际项目中这里会调用API
}

const resetFilters = () => {
  filters.value = {
    search: '',
    action: '',
    competitionId: '',
    userRole: '',
    dateRange: [],
    severity: '',
    ipAddress: ''
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
  selectedLogs.value = selection
}

const clearSelection = () => {
  selectedLogs.value = []
}

const viewLogDetails = (log) => {
  selectedLog.value = log
  showDetailDialog.value = true
}

const markAsResolved = async (log) => {
  try {
    await ElMessageBox.confirm(
      '确定要标记此日志为已处理吗？',
      '确认操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500))
    
    log.resolved = true
    ElMessage.success('日志已标记为已处理')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败：' + (error.message || '未知错误'))
    }
  }
}

const batchMarkResolved = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要批量标记 ${selectedLogs.value.length} 个日志为已处理吗？`,
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
    
    selectedLogs.value.forEach(log => {
      log.resolved = true
    })
    
    selectedLogs.value = []
    ElMessage.success('批量标记成功')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量操作失败：' + (error.message || '未知错误'))
    }
  } finally {
    batchProcessing.value = false
  }
}

const batchDelete = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要批量删除 ${selectedLogs.value.length} 个日志吗？此操作不可恢复！`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'error'
      }
    )
    
    batchProcessing.value = true
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    const idsToDelete = selectedLogs.value.map(log => log.id)
    auditLogs.value = auditLogs.value.filter(log => !idsToDelete.includes(log.id))
    
    selectedLogs.value = []
    loadStats()
    ElMessage.success('批量删除成功')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败：' + (error.message || '未知错误'))
    }
  } finally {
    batchProcessing.value = false
  }
}

const clearOldLogs = () => {
  clearForm.value = {
    daysToKeep: 30,
    severity: [],
    confirm: false
  }
  showClearDialog.value = true
}

const confirmClearLogs = async () => {
  if (!clearForm.value.confirm) {
    ElMessage.warning('请确认清理操作')
    return
  }
  
  clearing.value = true
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    const cutoffDate = new Date()
    cutoffDate.setDate(cutoffDate.getDate() - clearForm.value.daysToKeep)
    
    let logsToRemove = auditLogs.value.filter(log => 
      new Date(log.createdAt) < cutoffDate
    )
    
    if (clearForm.value.severity.length > 0) {
      logsToRemove = logsToRemove.filter(log => 
        clearForm.value.severity.includes(log.severity)
      )
    }
    
    auditLogs.value = auditLogs.value.filter(log => !logsToRemove.includes(log))
    
    loadStats()
    showClearDialog.value = false
    ElMessage.success(`成功清理 ${logsToRemove.length} 条旧日志`)
  } catch (error) {
    ElMessage.error('清理失败：' + (error.message || '未知错误'))
  } finally {
    clearing.value = false
  }
}

const exportLogs = () => {
  ElMessage.success('开始导出审计日志...')
  // 实际项目中这里会调用导出API
}

// 工具方法
const getRowClassName = ({ row }) => {
  if (row.resolved) return 'resolved-log'
  if (row.severity === 'critical') return 'critical-log'
  if (row.severity === 'error') return 'error-log'
  return ''
}

const getActionText = (action) => {
  const actionMap = {
    competition_created: '创建竞赛',
    competition_updated: '修改竞赛',
    competition_deleted: '删除竞赛',
    registration_reviewed: '报名审核',
    submission_created: '作品提交',
    judging_scored: '评审打分',
    results_published: '结果发布',
    user_login: '用户登录',
    permission_changed: '权限变更'
  }
  return actionMap[action] || action
}

const getRoleType = (role) => {
  const roleMap = {
    admin: 'danger',
    teacher: 'warning',
    student: 'primary'
  }
  return roleMap[role] || 'info'
}

const getRoleText = (role) => {
  const roleMap = {
    admin: '管理员',
    teacher: '教师',
    student: '学生'
  }
  return roleMap[role] || role
}

const getSeverityType = (severity) => {
  const severityMap = {
    info: 'info',
    warning: 'warning',
    error: 'danger',
    critical: 'danger'
  }
  return severityMap[severity] || 'info'
}

const getSeverityText = (severity) => {
  const severityMap = {
    info: '信息',
    warning: '警告',
    error: '错误',
    critical: '严重'
  }
  return severityMap[severity] || severity
}

const formatDateTime = (date) => {
  return formatDate(date, 'YYYY-MM-DD HH:mm:ss')
}

const getTimeAgo = (date) => {
  const now = new Date()
  const diff = now - new Date(date)
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)
  
  if (days > 0) return `${days}天前`
  if (hours > 0) return `${hours}小时前`
  if (minutes > 0) return `${minutes}分钟前`
  return '刚刚'
}

// 组件挂载时加载数据
onMounted(() => {
  loadAuditLogs()
  loadCompetitions()
})
</script>

<style scoped>
.competition-audit-logs {
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

.stat-icon.today {
  background: #409eff;
}

.stat-icon.critical {
  background: #f56c6c;
}

.stat-icon.users {
  background: #67c23a;
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

.logs-list-card {
  margin-bottom: 30px;
}

.time-info .time {
  font-weight: 600;
  color: #2c3e50;
}

.time-info .time-ago {
  font-size: 12px;
  color: #7f8c8d;
  margin-top: 2px;
}

.user-info h4 {
  margin: 0 0 5px 0;
  color: #2c3e50;
  font-size: 14px;
}

.user-meta {
  margin: 0;
  display: flex;
  align-items: center;
  gap: 10px;
}

.ip-address {
  font-size: 12px;
  color: #7f8c8d;
}

.action-info h4 {
  margin: 0 0 5px 0;
  color: #2c3e50;
  font-size: 14px;
}

.action-details {
  margin: 0;
  color: #7f8c8d;
  font-size: 12px;
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

.no-competition {
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

/* 日志详情样式 */
.log-detail .details-section,
.log-detail .competition-section,
.log-detail .data-changes-section {
  margin-top: 20px;
}

.log-detail h4 {
  margin: 0 0 15px 0;
  color: #2c3e50;
  font-size: 16px;
  border-bottom: 2px solid #409eff;
  padding-bottom: 10px;
}

.log-detail h5 {
  margin: 0 0 10px 0;
  color: #2c3e50;
  font-size: 14px;
}

.log-detail p {
  margin: 0;
  color: #606266;
  line-height: 1.6;
  background: #f8f9fa;
  padding: 15px;
  border-radius: 6px;
}

.log-detail pre {
  margin: 0;
  color: #606266;
  font-size: 12px;
  background: #f8f9fa;
  padding: 15px;
  border-radius: 6px;
  overflow-x: auto;
}

.no-data {
  color: #c0c4cc;
  font-style: italic;
  text-align: center;
  padding: 20px;
}

/* 表格行样式 */
:deep(.resolved-log) {
  background-color: #f0f9ff;
}

:deep(.critical-log) {
  background-color: #fef0f0;
}

:deep(.error-log) {
  background-color: #fdf6ec;
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
  
  .user-meta {
    flex-direction: column;
    align-items: flex-start;
    gap: 5px;
  }
}
</style> 