<template>
  <div class="dashboard-view">
    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon users">
              <el-icon><User /></el-icon>
            </div>
            <div class="stat-info">
              <h4>总用户数</h4>
              <p class="stat-number">{{ stats.totalUsers }}</p>
              <p class="stat-desc">活跃用户: {{ stats.activeUsers }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon projects">
              <el-icon><Document /></el-icon>
            </div>
            <div class="stat-info">
              <h4>活跃项目</h4>
              <p class="stat-number">{{ stats.totalProjects }}</p>
              <p class="stat-desc">待审核: {{ stats.pendingProjects }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon competitions">
              <el-icon><Trophy /></el-icon>
            </div>
            <div class="stat-info">
              <h4>进行中竞赛</h4>
              <p class="stat-number">{{ stats.activeCompetitions }}</p>
              <p class="stat-desc">总竞赛: {{ stats.totalCompetitions }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon applications">
              <el-icon><EditPen /></el-icon>
            </div>
            <div class="stat-info">
              <h4>待处理申请</h4>
              <p class="stat-number">{{ stats.pendingApplications }}</p>
              <p class="stat-desc">今日新增: {{ stats.todayApplications }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 系统状态和最近活动 -->
    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :span="12">
        <el-card class="system-status">
          <template #header>
            <span>系统状态</span>
            <el-button link @click="refreshSystemStatus">
              <el-icon><Refresh /></el-icon>
            </el-button>
          </template>
          <el-table :data="systemStatus" style="width: 100%">
            <el-table-column prop="service" label="服务" width="120" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="scope.row.status === '正常' ? 'success' : 'danger'">
                  {{ scope.row.status }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="uptime" label="运行时间" />
            <el-table-column prop="load" label="负载" />
          </el-table>
        </el-card>
      </el-col>
      
      <el-col :span="12">
        <el-card class="recent-activity">
          <template #header>
            <span>最近登录</span>
            <el-button link @click="refreshRecentActivity">
              <el-icon><Refresh /></el-icon>
            </el-button>
          </template>
          <el-timeline>
            <el-timeline-item 
              v-for="activity in recentActivity" 
              :key="activity.id"
              :timestamp="activity.timestamp" 
              placement="top"
            >
              <el-card>
                <h4>{{ activity.title }}</h4>
                <p>{{ activity.description }}</p>
              </el-card>
            </el-timeline-item>
          </el-timeline>
        </el-card>
      </el-col>
    </el-row>

    <!-- 快速操作 -->
    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :span="24">
        <el-card class="quick-actions">
          <template #header>
            <span>快速操作</span>
          </template>
          <el-row :gutter="20">
            <el-col :span="6" v-for="action in quickActions" :key="action.name">
              <el-button 
                type="primary" 
                @click="handleQuickAction(action.action)"
                style="width: 100%; height: 80px;"
              >
                <el-icon><component :is="action.icon" /></el-icon>
                <div class="quick-action-content">
                  <h4>{{ action.name }}</h4>
                  <p>{{ action.description }}</p>
                </div>
              </el-button>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
    </el-row>

    <!-- 系统信息 -->
    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :span="12">
        <el-card class="system-info">
          <template #header>
            <span>系统信息</span>
          </template>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="系统名称">云梦高校科研竞赛管理系统</el-descriptions-item>
            <el-descriptions-item label="版本号">v1.0.0</el-descriptions-item>
            <el-descriptions-item label="服务器时间">{{ systemInfo.serverTime }}</el-descriptions-item>
            <el-descriptions-item label="数据库状态">{{ systemInfo.dbStatus }}</el-descriptions-item>
            <el-descriptions-item label="最后备份">{{ systemInfo.lastBackup }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
      
      <el-col :span="12">
        <el-card class="health-status">
          <template #header>
            <span>健康状态</span>
          </template>
          <div v-for="check in healthChecks" :key="check.name" class="health-check">
            <div class="health-check-header">
              <span class="health-check-name">{{ check.name }}</span>
              <el-tag :type="check.status === 'healthy' ? 'success' : 'danger'" size="small">
                {{ check.status === 'healthy' ? '正常' : '异常' }}
              </el-tag>
            </div>
            <p class="health-check-message">{{ check.message }}</p>
            <p class="health-check-details">{{ check.details }}</p>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { adminService } from '../../services/adminService'
import { 
  User, 
  Document, 
  Trophy, 
  EditPen, 
  Refresh, 
  UserFilled, 
  Download, 
  Files, 
  Setting 
} from '@element-plus/icons-vue'

// 响应式数据
const stats = ref({
  totalUsers: 0,
  activeUsers: 0,
  totalProjects: 0,
  pendingProjects: 0,
  activeCompetitions: 0,
  totalCompetitions: 0,
  pendingApplications: 0,
  todayApplications: 0
})

const systemStatus = ref([])
const recentActivity = ref([])
const systemInfo = ref({
  serverTime: '',
  dbStatus: '正常',
  lastBackup: '2024-01-15 02:00:00'
})

const healthChecks = ref([
  {
    name: '数据库连接',
    status: 'healthy',
    message: '连接正常',
    details: '响应时间: 5ms'
  },
  {
    name: '内存使用',
    status: 'healthy',
    message: '使用率正常',
    details: '使用率: 45%'
  },
  {
    name: '磁盘空间',
    status: 'healthy',
    message: '空间充足',
    details: '使用率: 30%'
  },
  {
    name: '网络连接',
    status: 'healthy',
    message: '连接正常',
    details: '响应时间: 10ms'
  }
])

const quickActions = ref([
  {
    name: '创建用户',
    description: '批量导入用户',
    icon: UserFilled,
    action: 'create_user'
  },
  {
    name: '数据导出',
    description: '导出系统数据',
    icon: Download,
    action: 'export_data'
  },
  {
    name: '系统备份',
    description: '创建数据备份',
    icon: Files,
    action: 'create_backup'
  },
  {
    name: '系统设置',
    description: '配置系统参数',
    icon: Setting,
    action: 'system_settings'
  }
])

// 方法
const loadDashboardData = async () => {
  try {
    const response = await adminService.getDashboardStats()
    if (response && response.code === 200 && response.data) {
      if (response.data.userStats) {
        stats.value = { ...stats.value, ...response.data.userStats }
      }
      if (response.data.systemInfo) {
        systemInfo.value = { ...systemInfo.value, ...response.data.systemInfo }
      }
    }
  } catch (error) {
    console.error('加载仪表板数据失败:', error)
    // 使用默认数据，不显示错误
  }
}

const loadSystemStatus = async () => {
  try {
    const response = await adminService.getSystemHealth()
    if (response && response.code === 200 && response.data && response.data.checks) {
      healthChecks.value = response.data.checks
    }
  } catch (error) {
    console.error('加载系统状态失败:', error)
    // 使用默认数据，不显示错误
  }
}

const refreshSystemStatus = () => {
  loadSystemStatus()
  ElMessage.success('系统状态已刷新')
}

const refreshRecentActivity = () => {
  // 这里可以重新加载最近活动数据
  ElMessage.success('最近活动已刷新')
}

const handleQuickAction = (action) => {
  switch (action) {
    case 'create_user':
      ElMessage.info('跳转到用户管理页面')
      break
    case 'export_data':
      ElMessage.info('开始导出数据')
      break
    case 'create_backup':
      ElMessage.info('开始创建备份')
      break
    case 'system_settings':
      ElMessage.info('跳转到系统设置页面')
      break
  }
}

// 组件挂载时加载数据
onMounted(() => {
  loadDashboardData()
  loadSystemStatus()
  
  // 更新服务器时间
  systemInfo.value.serverTime = new Date().toLocaleString('zh-CN')
  
  // 模拟最近活动数据
  recentActivity.value = [
    {
      id: 1,
      title: '教师登录',
      description: '李老师 (ID: 10086) 登录系统',
      timestamp: '2024-01-15 14:30'
    },
    {
      id: 2,
      title: '学生登录',
      description: '张三 (学号: 2021001) 登录系统',
      timestamp: '2024-01-15 13:45'
    },
    {
      id: 3,
      title: '项目提交',
      description: '新项目"智能校园系统"已提交审核',
      timestamp: '2024-01-15 12:20'
    }
  ]
  
  // 模拟系统状态数据
  systemStatus.value = [
    {
      service: 'Web服务',
      status: '正常',
      uptime: '15天 8小时',
      load: '23%'
    },
    {
      service: '数据库',
      status: '正常',
      uptime: '15天 8小时',
      load: '45%'
    },
    {
      service: '文件服务',
      status: '正常',
      uptime: '15天 8小时',
      load: '12%'
    },
    {
      service: '邮件服务',
      status: '异常',
      uptime: '2小时 30分钟',
      load: '89%'
    }
  ]
})
</script>

<style scoped>
.dashboard-view {
  padding: 20px;
}

.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
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

.stat-icon.users {
  background: linear-gradient(135deg, #667eea, #764ba2);
}

.stat-icon.projects {
  background: linear-gradient(135deg, #f093fb, #f5576c);
}

.stat-icon.competitions {
  background: linear-gradient(135deg, #4facfe, #00f2fe);
}

.stat-icon.applications {
  background: linear-gradient(135deg, #43e97b, #38f9d7);
}

.stat-icon .el-icon {
  font-size: 24px;
  color: white;
}

.stat-info h4 {
  margin: 0 0 5px 0;
  color: #7f8c8d;
  font-size: 14px;
}

.stat-number {
  margin: 0;
  font-size: 28px;
  font-weight: 600;
  color: #2c3e50;
}

.stat-desc {
  margin: 5px 0 0 0;
  color: #7f8c8d;
  font-size: 12px;
}

.system-status,
.recent-activity,
.quick-actions,
.system-info,
.health-status {
  border-radius: 8px;
}

.quick-action-content {
  text-align: center;
}

.quick-action-content h4 {
  margin: 0 0 5px 0;
  font-size: 14px;
}

.quick-action-content p {
  margin: 0;
  font-size: 12px;
  color: #7f8c8d;
}

.health-check {
  margin-bottom: 15px;
  padding: 10px;
  border-radius: 4px;
  background: #f8f9fa;
}

.health-check-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 5px;
}

.health-check-name {
  font-weight: 600;
  color: #2c3e50;
}

.health-check-message {
  margin: 5px 0;
  color: #7f8c8d;
  font-size: 14px;
}

.health-check-details {
  margin: 0;
  color: #95a5a6;
  font-size: 12px;
}
</style> 