<template>
  <div class="project-overview">
    <el-row :gutter="20">
      <!-- 统计卡片 -->
      <el-col :span="6" v-for="stat in projectStats" :key="stat.title">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" :class="stat.type">
              <el-icon><component :is="stat.icon" /></el-icon>
            </div>
            <div class="stat-info">
              <h4>{{ stat.title }}</h4>
              <p class="stat-number">{{ stat.value }}</p>
              <p class="stat-desc">{{ stat.description }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 项目状态分布 -->
    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>项目状态分布</span>
          </template>
          <div id="statusChart" style="height: 300px;"></div>
        </el-card>
      </el-col>
      
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>项目类型分布</span>
          </template>
          <div id="typeChart" style="height: 300px;"></div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 最近项目 -->
    <el-row style="margin-top: 20px;">
      <el-col :span="24">
        <el-card>
          <template #header>
            <span>最近项目</span>
            <el-button link @click="refreshRecentProjects">
              <el-icon><Refresh /></el-icon>
            </el-button>
          </template>
          <el-table :data="recentProjects" style="width: 100%">
            <el-table-column prop="name" label="项目名称" />
            <el-table-column prop="type" label="项目类型" />
            <el-table-column prop="status" label="状态">
              <template #default="scope">
                <el-tag :type="getStatusType(scope.row.status)">
                  {{ scope.row.status }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="createTime" label="创建时间" />
            <el-table-column prop="progress" label="进度">
              <template #default="scope">
                <el-progress :percentage="scope.row.progress" />
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { Document, Clock, Check, Warning, Refresh } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import { adminService } from '../../services/adminService'

// 响应式数据
const projectStats = ref([
  {
    title: '总项目数',
    value: 0,
    description: '所有项目',
    icon: Document,
    type: 'total'
  },
  {
    title: '进行中',
    value: 0,
    description: '正在执行',
    icon: Clock,
    type: 'ongoing'
  },
  {
    title: '已完成',
    value: 0,
    description: '成功完成',
    icon: Check,
    type: 'completed'
  },
  {
    title: '待审核',
    value: 0,
    description: '等待审核',
    icon: Warning,
    type: 'pending'
  }
])

const recentProjects = ref([])

// 获取状态类型
const getStatusType = (status) => {
  const statusMap = {
    '进行中': 'primary',
    '已完成': 'success',
    '待审核': 'warning',
    '已拒绝': 'danger',
    '已暂停': 'info'
  }
  return statusMap[status] || 'info'
}

// 加载项目统计数据
const loadProjectStats = async () => {
  try {
    const response = await adminService.getProjectStats()
    if (response && response.code === 200) {
      const stats = response.data || {}
      projectStats.value[0].value = stats.totalProjects || 0
      projectStats.value[1].value = stats.ongoingProjects || 0
      projectStats.value[2].value = stats.completedProjects || 0
      projectStats.value[3].value = stats.pendingProjects || 0
    }
  } catch (error) {
    console.error('加载项目统计数据失败:', error)
    // 使用默认数据
  }
}

// 加载最近项目
const loadRecentProjects = async () => {
  try {
    // 这里可以调用实际的API，暂时使用模拟数据
    recentProjects.value = [
      {
        name: '智能校园系统',
        type: '软件开发',
        status: '进行中',
        createTime: '2024-01-15',
        progress: 75
      },
      {
        name: '数据分析平台',
        type: '科研项目',
        status: '待审核',
        createTime: '2024-01-14',
        progress: 90
      },
      {
        name: '在线教育平台',
        type: '创新项目',
        status: '已完成',
        createTime: '2024-01-13',
        progress: 100
      }
    ]
  } catch (error) {
    console.error('加载最近项目失败:', error)
  }
}

// 刷新最近项目
const refreshRecentProjects = () => {
  loadRecentProjects()
  ElMessage.success('项目列表已刷新')
}

// 初始化图表
const initCharts = async () => {
  await nextTick()
  
  // 状态分布图表
  const statusChart = echarts.init(document.getElementById('statusChart'))
  statusChart.setOption({
    title: { text: '项目状态分布' },
    tooltip: { trigger: 'item' },
    series: [{
      type: 'pie',
      radius: '50%',
      data: [
        { value: 45, name: '进行中' },
        { value: 30, name: '已完成' },
        { value: 15, name: '待审核' },
        { value: 10, name: '已暂停' }
      ]
    }]
  })

  // 类型分布图表
  const typeChart = echarts.init(document.getElementById('typeChart'))
  typeChart.setOption({
    title: { text: '项目类型分布' },
    tooltip: { trigger: 'item' },
    series: [{
      type: 'pie',
      radius: '50%',
      data: [
        { value: 40, name: '软件开发' },
        { value: 25, name: '科研项目' },
        { value: 20, name: '创新项目' },
        { value: 15, name: '竞赛项目' }
      ]
    }]
  })
}

// 组件挂载时加载数据
onMounted(async () => {
  await loadProjectStats()
  await loadRecentProjects()
  await initCharts()
})
</script>

<style scoped>
.project-overview {
  padding: 20px;
}

.stat-card {
  margin-bottom: 20px;
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

.stat-icon.total {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-icon.ongoing {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.stat-icon.completed {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-icon.pending {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.stat-icon i {
  font-size: 24px;
  color: white;
}

.stat-info h4 {
  margin: 0 0 5px 0;
  color: #7f8c8d;
  font-size: 14px;
}

.stat-number {
  margin: 0 0 5px 0;
  font-size: 28px;
  font-weight: 600;
  color: #2c3e50;
}

.stat-desc {
  margin: 0;
  color: #95a5a6;
  font-size: 12px;
}
</style> 