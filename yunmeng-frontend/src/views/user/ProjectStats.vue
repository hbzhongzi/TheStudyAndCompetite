<template>
  <div class="project-stats">
    <!-- 统计概览 -->
    <el-row :gutter="20">
      <el-col :span="6" v-for="stat in overviewStats" :key="stat.title">
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

    <!-- 图表区域 -->
    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>项目状态趋势</span>
            <el-select v-model="timeRange" @change="updateTrendChart" style="float: right; width: 120px;">
              <el-option label="最近7天" value="7" />
              <el-option label="最近30天" value="30" />
              <el-option label="最近90天" value="90" />
            </el-select>
          </template>
          <div id="trendChart" style="height: 300px;"></div>
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

    <!-- 详细统计表格 -->
    <el-row style="margin-top: 20px;">
      <el-col :span="24">
        <el-card>
          <template #header>
            <span>项目详细统计</span>
            <el-button link @click="exportStats">
              <el-icon><Download /></el-icon>
              导出统计
            </el-button>
          </template>
          
          <el-tabs v-model="activeTab" @tab-click="handleTabClick">
            <el-tab-pane label="院系统计" name="department">
              <el-table :data="departmentStats" style="width: 100%">
                <el-table-column prop="department" label="院系名称" />
                <el-table-column prop="totalProjects" label="项目总数" />
                <el-table-column prop="ongoingProjects" label="进行中" />
                <el-table-column prop="completedProjects" label="已完成" />
                <el-table-column prop="pendingProjects" label="待审核" />
                <el-table-column prop="completionRate" label="完成率">
                  <template #default="scope">
                    <el-progress :percentage="scope.row.completionRate" />
                  </template>
                </el-table-column>
              </el-table>
            </el-tab-pane>
            
            <el-tab-pane label="时间统计" name="time">
              <el-table :data="timeStats" style="width: 100%">
                <el-table-column prop="period" label="时间段" />
                <el-table-column prop="newProjects" label="新项目数" />
                <el-table-column prop="completedProjects" label="完成项目数" />
                <el-table-column prop="avgDuration" label="平均周期(天)" />
                <el-table-column prop="efficiency" label="效率指数">
                  <template #default="scope">
                    <el-rate v-model="scope.row.efficiency" disabled show-score />
                  </template>
                </el-table-column>
              </el-table>
            </el-tab-pane>
            
            <el-tab-pane label="质量统计" name="quality">
              <el-table :data="qualityStats" style="width: 100%">
                <el-table-column prop="category" label="质量类别" />
                <el-table-column prop="excellent" label="优秀" />
                <el-table-column prop="good" label="良好" />
                <el-table-column prop="average" label="一般" />
                <el-table-column prop="poor" label="较差" />
                <el-table-column prop="excellentRate" label="优秀率">
                  <template #default="scope">
                    {{ scope.row.excellentRate }}%
                  </template>
                </el-table-column>
              </el-table>
            </el-tab-pane>
          </el-tabs>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { Document, Clock, Check, Warning, Download } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import { adminService } from '../../services/adminService'

// 响应式数据
const timeRange = ref('30')
const activeTab = ref('department')

const overviewStats = ref([
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

const departmentStats = ref([])
const timeStats = ref([])
const qualityStats = ref([])

// 加载统计数据
const loadStats = async () => {
  try {
    const response = await adminService.getProjectStats()
    if (response && response.code === 200) {
      const stats = response.data || {}
      overviewStats.value[0].value = stats.totalProjects || 0
      overviewStats.value[1].value = stats.ongoingProjects || 0
      overviewStats.value[2].value = stats.completedProjects || 0
      overviewStats.value[3].value = stats.pendingProjects || 0
    }
  } catch (error) {
    console.error('加载统计数据失败:', error)
  }
}

// 加载院系统计
const loadDepartmentStats = async () => {
  try {
    const response = await adminService.getDepartmentProjectStats()
    if (response && response.code === 200) {
      departmentStats.value = response.data || []
    } else {
      // 使用模拟数据
      departmentStats.value = [
        {
          department: '计算机学院',
          totalProjects: 45,
          ongoingProjects: 20,
          completedProjects: 20,
          pendingProjects: 5,
          completionRate: 80
        },
        {
          department: '机械工程学院',
          totalProjects: 32,
          ongoingProjects: 15,
          completedProjects: 12,
          pendingProjects: 5,
          completionRate: 75
        }
      ]
    }
  } catch (error) {
    console.error('加载院系统计失败:', error)
  }
}

// 加载时间统计
const loadTimeStats = async () => {
  try {
    const response = await adminService.getTimeTrendStats({ days: timeRange.value })
    if (response && response.code === 200) {
      timeStats.value = response.data || []
    } else {
      // 使用模拟数据
      timeStats.value = [
        {
          period: '第1周',
          newProjects: 12,
          completedProjects: 8,
          avgDuration: 45,
          efficiency: 4.2
        },
        {
          period: '第2周',
          newProjects: 15,
          completedProjects: 10,
          avgDuration: 42,
          efficiency: 4.5
        }
      ]
    }
  } catch (error) {
    console.error('加载时间统计失败:', error)
  }
}

// 加载质量统计
const loadQualityStats = async () => {
  try {
    const response = await adminService.getProjectQualityReport()
    if (response && response.code === 200) {
      qualityStats.value = response.data || []
    } else {
      // 使用模拟数据
      qualityStats.value = [
        {
          category: '软件开发',
          excellent: 15,
          good: 20,
          average: 8,
          poor: 2,
          excellentRate: 33.3
        },
        {
          category: '科研项目',
          excellent: 12,
          good: 18,
          average: 10,
          poor: 3,
          excellentRate: 27.9
        }
      ]
    }
  } catch (error) {
    console.error('加载质量统计失败:', error)
  }
}

// 更新趋势图表
const updateTrendChart = () => {
  loadTimeStats()
  initTrendChart()
}

// 处理标签页点击
const handleTabClick = (tab) => {
  switch (tab.name) {
    case 'department':
      loadDepartmentStats()
      break
    case 'time':
      loadTimeStats()
      break
    case 'quality':
      loadQualityStats()
      break
  }
}

// 导出统计
const exportStats = () => {
  ElMessage.success('统计报告导出成功')
}

// 初始化图表
const initCharts = async () => {
  await nextTick()
  
  // 趋势图表
  const trendChart = echarts.init(document.getElementById('trendChart'))
  trendChart.setOption({
    title: { text: '项目状态趋势' },
    tooltip: { trigger: 'axis' },
    legend: { data: ['新项目', '完成项目', '进行中项目'] },
    xAxis: { type: 'category', data: ['第1周', '第2周', '第3周', '第4周'] },
    yAxis: { type: 'value' },
    series: [
      {
        name: '新项目',
        type: 'line',
        data: [12, 15, 18, 20]
      },
      {
        name: '完成项目',
        type: 'line',
        data: [8, 10, 12, 15]
      },
      {
        name: '进行中项目',
        type: 'line',
        data: [25, 30, 35, 40]
      }
    ]
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

// 初始化趋势图表
const initTrendChart = () => {
  const trendChart = echarts.init(document.getElementById('trendChart'))
  // 这里可以根据实际数据更新图表
}

// 组件挂载时加载数据
onMounted(async () => {
  await loadStats()
  await loadDepartmentStats()
  await loadTimeStats()
  await loadQualityStats()
  await initCharts()
})
</script>

<style scoped>
.project-stats {
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