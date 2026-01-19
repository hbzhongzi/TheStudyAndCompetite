<template>
  <div class="system-logs">
    <div class="page-header">
      <h2>系统日志</h2>
      <div class="header-actions">
        <el-button type="primary" @click="exportLogs">
          <i class="el-icon-download"></i>
          导出日志
        </el-button>
        <el-button type="warning" @click="clearLogs">
          <i class="el-icon-delete"></i>
          清空日志
        </el-button>
      </div>
    </div>

    <!-- 搜索和筛选 -->
    <el-card class="search-card">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-input
            v-model="searchQuery"
            placeholder="搜索日志内容"
            clearable
            @input="handleSearch"
          >
            <template #prefix>
              <i class="el-icon-search"></i>
            </template>
          </el-input>
        </el-col>
        <el-col :span="4">
          <el-select v-model="levelFilter" placeholder="日志级别" clearable @change="handleSearch">
            <el-option label="全部级别" value="" />
            <el-option label="INFO" value="INFO" />
            <el-option label="WARNING" value="WARNING" />
            <el-option label="ERROR" value="ERROR" />
            <el-option label="DEBUG" value="DEBUG" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-select v-model="moduleFilter" placeholder="模块筛选" clearable @change="handleSearch">
            <el-option label="全部模块" value="" />
            <el-option label="用户管理" value="user" />
            <el-option label="项目管理" value="project" />
            <el-option label="竞赛管理" value="competition" />
            <el-option label="系统管理" value="system" />
            <el-option label="认证授权" value="auth" />
          </el-select>
        </el-col>
        <el-col :span="6">
          <el-date-picker
            v-model="dateRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
            @change="handleSearch"
            style="width: 100%"
          />
        </el-col>
        <el-col :span="4">
          <el-button type="primary" @click="handleSearch">搜索</el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 日志列表 -->
    <el-card class="log-list-card">
      <div class="log-controls">
        <el-switch
          v-model="autoRefresh"
          active-text="自动刷新"
          inactive-text="手动刷新"
          @change="toggleAutoRefresh"
        />
        <span class="refresh-info" v-if="autoRefresh">
          每 {{ refreshInterval }} 秒自动刷新
        </span>
      </div>
      
      <el-table :data="filteredLogs" style="width: 100%" v-loading="loading" max-height="600">
        <el-table-column prop="timestamp" label="时间" width="180" sortable>
          <template #default="scope">
            {{ formatTime(scope.row.timestamp) }}
          </template>
        </el-table-column>
        <el-table-column prop="level" label="级别" width="100">
          <template #default="scope">
            <el-tag :type="getLevelTagType(scope.row.level)" size="small">
              {{ scope.row.level }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="module" label="模块" width="120">
          <template #default="scope">
            <el-tag type="info" size="small">{{ getModuleLabel(scope.row.module) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="user" label="用户" width="120" />
        <el-table-column prop="ip" label="IP地址" width="140" />
        <el-table-column prop="action" label="操作" width="150" />
        <el-table-column prop="message" label="日志内容" min-width="300">
          <template #default="scope">
            <div class="log-message">
              <span class="message-text">{{ scope.row.message }}</span>
              <el-button 
                v-if="scope.row.details" 
                link 
                size="small" 
                @click="showLogDetails(scope.row)"
              >
                详情
              </el-button>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="scope">
            <el-button size="small" @click="viewLogDetails(scope.row)">查看</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[20, 50, 100, 200]"
          :total="totalLogs"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 日志详情对话框 -->
    <el-dialog
      v-model="showDetailsDialog"
      title="日志详情"
      width="800px"
      :close-on-click-modal="false"
    >
      <div v-if="selectedLog" class="log-details">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="时间">{{ formatTime(selectedLog.timestamp) }}</el-descriptions-item>
          <el-descriptions-item label="级别">
            <el-tag :type="getLevelTagType(selectedLog.level)">{{ selectedLog.level }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="模块">{{ getModuleLabel(selectedLog.module) }}</el-descriptions-item>
          <el-descriptions-item label="用户">{{ selectedLog.user }}</el-descriptions-item>
          <el-descriptions-item label="IP地址">{{ selectedLog.ip }}</el-descriptions-item>
          <el-descriptions-item label="操作">{{ selectedLog.action }}</el-descriptions-item>
        </el-descriptions>
        
        <div class="log-message-detail">
          <h4>日志内容：</h4>
          <div class="message-content">{{ selectedLog.message }}</div>
        </div>
        
        <div v-if="selectedLog.details" class="log-details-extra">
          <h4>详细信息：</h4>
          <pre class="details-content">{{ selectedLog.details }}</pre>
        </div>
        
        <div v-if="selectedLog.stackTrace" class="log-stack-trace">
          <h4>堆栈信息：</h4>
          <pre class="stack-content">{{ selectedLog.stackTrace }}</pre>
        </div>
      </div>
    </el-dialog>

    <!-- 清空日志确认对话框 -->
    <el-dialog
      v-model="showClearDialog"
      title="确认清空日志"
      width="400px"
      :close-on-click-modal="false"
    >
      <div class="clear-confirmation">
        <p>确定要清空所有日志吗？此操作不可恢复。</p>
        <el-form :model="clearForm" label-width="100px">
          <el-form-item label="保留天数">
            <el-input-number
              v-model="clearForm.keepDays"
              :min="1"
              :max="365"
              style="width: 200px"
            />
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showClearDialog = false">取消</el-button>
          <el-button type="danger" @click="confirmClearLogs" :loading="clearing">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

// 响应式数据
const loading = ref(false)
const clearing = ref(false)
const searchQuery = ref('')
const levelFilter = ref('')
const moduleFilter = ref('')
const dateRange = ref([])
const currentPage = ref(1)
const pageSize = ref(50)
const totalLogs = ref(0)
const autoRefresh = ref(false)
const refreshInterval = ref(30)
const showDetailsDialog = ref(false)
const showClearDialog = ref(false)
const selectedLog = ref(null)

// 清空表单
const clearForm = reactive({
  keepDays: 30
})

// 模拟日志数据
const logs = ref([
  {
    id: 1,
    timestamp: '2024-01-15 14:30:25',
    level: 'INFO',
    module: 'user',
    user: 'admin',
    ip: '192.168.1.100',
    action: '用户登录',
    message: '用户 admin 成功登录系统',
    details: '登录时间: 2024-01-15 14:30:25\n浏览器: Chrome 120.0.0.0\n操作系统: Windows 10'
  },
  {
    id: 2,
    timestamp: '2024-01-15 14:25:18',
    level: 'WARNING',
    module: 'auth',
    user: 'teacher001',
    ip: '192.168.1.101',
    action: '登录失败',
    message: '用户 teacher001 登录失败，密码错误',
    details: '尝试次数: 3\n锁定状态: 已锁定\n锁定时间: 30分钟'
  },
  {
    id: 3,
    timestamp: '2024-01-15 14:20:45',
    level: 'ERROR',
    module: 'system',
    user: 'system',
    ip: '127.0.0.1',
    action: '系统异常',
    message: '数据库连接超时',
    details: '连接池状态: 已满\n超时时间: 30秒\n重试次数: 3',
    stackTrace: 'java.sql.SQLException: Connection timeout\n    at com.mysql.jdbc.ConnectionImpl.connect(ConnectionImpl.java:1234)\n    at com.mysql.jdbc.ConnectionImpl.<init>(ConnectionImpl.java:567)\n    at com.mysql.jdbc.Driver.connect(Driver.java:890)'
  },
  {
    id: 4,
    timestamp: '2024-01-15 14:15:32',
    level: 'INFO',
    module: 'project',
    user: 'teacher001',
    ip: '192.168.1.101',
    action: '创建项目',
    message: '用户 teacher001 创建了新项目：人工智能研究项目',
    details: '项目ID: PRJ001\n项目类型: 科研项目\n预算: 50000元'
  },
  {
    id: 5,
    timestamp: '2024-01-15 14:10:15',
    level: 'DEBUG',
    module: 'competition',
    user: 'student001',
    ip: '192.168.1.102',
    action: '查看竞赛',
    message: '用户 student001 查看了竞赛：全国大学生数学建模竞赛',
    details: '竞赛ID: COMP001\n查看时长: 2分钟\n页面: 竞赛详情页'
  }
])

// 计算属性
const filteredLogs = computed(() => {
  let result = logs.value

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(log => 
      log.message.toLowerCase().includes(query) ||
      log.action.toLowerCase().includes(query) ||
      log.user.toLowerCase().includes(query)
    )
  }

  if (levelFilter.value) {
    result = result.filter(log => log.level === levelFilter.value)
  }

  if (moduleFilter.value) {
    result = result.filter(log => log.module === moduleFilter.value)
  }

  if (dateRange.value && dateRange.value.length === 2) {
    const startTime = new Date(dateRange.value[0]).getTime()
    const endTime = new Date(dateRange.value[1]).getTime()
    result = result.filter(log => {
      const logTime = new Date(log.timestamp).getTime()
      return logTime >= startTime && logTime <= endTime
    })
  }

  return result
})

// 定时器
let refreshTimer = null

// 方法
const getLevelTagType = (level) => {
  const typeMap = {
    INFO: 'success',
    WARNING: 'warning',
    ERROR: 'danger',
    DEBUG: 'info'
  }
  return typeMap[level] || 'info'
}

const getModuleLabel = (module) => {
  const moduleMap = {
    user: '用户管理',
    project: '项目管理',
    competition: '竞赛管理',
    system: '系统管理',
    auth: '认证授权'
  }
  return moduleMap[module] || module
}

const formatTime = (timestamp) => {
  return new Date(timestamp).toLocaleString('zh-CN')
}

const handleSearch = () => {
  currentPage.value = 1
  // 这里可以调用API进行搜索
}

const handleSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
  // 重新加载数据
}

const handleCurrentChange = (page) => {
  currentPage.value = page
  // 重新加载数据
}

const toggleAutoRefresh = (value) => {
  if (value) {
    startAutoRefresh()
  } else {
    stopAutoRefresh()
  }
}

const startAutoRefresh = () => {
  refreshTimer = setInterval(() => {
    // 模拟刷新数据
    console.log('自动刷新日志数据')
  }, refreshInterval.value * 1000)
}

const stopAutoRefresh = () => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
}

const viewLogDetails = (log) => {
  selectedLog.value = log
  showDetailsDialog.value = true
}

const showLogDetails = (log) => {
  selectedLog.value = log
  showDetailsDialog.value = true
}

const exportLogs = () => {
  ElMessage.success('日志导出成功')
  // 这里可以实现实际的导出功能
}

const clearLogs = () => {
  showClearDialog.value = true
}

const confirmClearLogs = async () => {
  try {
    clearing.value = true
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 2000))
    
    ElMessage.success(`日志清空成功，保留最近 ${clearForm.keepDays} 天的日志`)
    showClearDialog.value = false
    
    // 重新加载数据
    // loadLogs()
  } catch (error) {
    ElMessage.error('清空日志失败')
  } finally {
    clearing.value = false
  }
}

// 生命周期
onMounted(() => {
  totalLogs.value = logs.value.length
})

onUnmounted(() => {
  stopAutoRefresh()
})
</script>

<style scoped>
.system-logs {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
  color: #2c3e50;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.search-card {
  margin-bottom: 20px;
}

.log-list-card {
  margin-bottom: 20px;
}

.log-controls {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 15px;
  padding: 10px;
  background: #f8f9fa;
  border-radius: 5px;
}

.refresh-info {
  color: #7f8c8d;
  font-size: 14px;
}

.log-message {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.message-text {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.log-details {
  padding: 20px 0;
}

.log-message-detail {
  margin-top: 20px;
}

.log-message-detail h4 {
  margin: 0 0 10px 0;
  color: #2c3e50;
}

.message-content {
  padding: 15px;
  background: #f8f9fa;
  border-radius: 5px;
  border-left: 4px solid #409EFF;
  font-family: 'Courier New', monospace;
  white-space: pre-wrap;
  word-break: break-all;
}

.log-details-extra {
  margin-top: 20px;
}

.log-details-extra h4 {
  margin: 0 0 10px 0;
  color: #2c3e50;
}

.details-content {
  padding: 15px;
  background: #f8f9fa;
  border-radius: 5px;
  border-left: 4px solid #67C23A;
  font-family: 'Courier New', monospace;
  white-space: pre-wrap;
  word-break: break-all;
  margin: 0;
}

.log-stack-trace {
  margin-top: 20px;
}

.log-stack-trace h4 {
  margin: 0 0 10px 0;
  color: #2c3e50;
}

.stack-content {
  padding: 15px;
  background: #f8f9fa;
  border-radius: 5px;
  border-left: 4px solid #F56C6C;
  font-family: 'Courier New', monospace;
  white-space: pre-wrap;
  word-break: break-all;
  margin: 0;
  max-height: 300px;
  overflow-y: auto;
}

.clear-confirmation {
  padding: 20px 0;
}

.clear-confirmation p {
  margin: 0 0 20px 0;
  color: #F56C6C;
  font-weight: 500;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style> 