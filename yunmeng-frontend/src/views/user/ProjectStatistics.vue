<template>
  <div class="project-statistics">
    <!-- 统计概览 -->
    <el-row :gutter="20" class="stats-overview">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon research">
              <i class="el-icon-science"></i>
            </div>
            <div class="stat-info">
              <h4>科研项目</h4>
              <p class="stat-number">{{ stats.researchCount }}</p>
              <p class="stat-percent">{{ stats.researchPercent }}%</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon competition">
              <i class="el-icon-trophy"></i>
            </div>
            <div class="stat-info">
              <h4>竞赛项目</h4>
              <p class="stat-number">{{ stats.competitionCount }}</p>
              <p class="stat-percent">{{ stats.competitionPercent }}%</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon participation">
              <i class="el-icon-user"></i>
            </div>
            <div class="stat-info">
              <h4>参与学生</h4>
              <p class="stat-number">{{ stats.studentCount }}</p>
              <p class="stat-percent">参与率 {{ stats.participationRate }}%</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon teacher">
              <i class="el-icon-s-custom"></i>
            </div>
            <div class="stat-info">
              <h4>指导老师</h4>
              <p class="stat-number">{{ stats.teacherCount }}</p>
              <p class="stat-percent">活跃度 {{ stats.teacherActivity }}%</p>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 图表区域 -->
    <el-row :gutter="20" class="charts-row">
      <el-col :span="12">
        <el-card class="chart-card">
          <template #header>
            <span>项目状态分布</span>
          </template>
          <div class="chart-container">
            <div ref="statusChart" class="chart"></div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card class="chart-card">
          <template #header>
            <span>院系项目分布</span>
          </template>
          <div class="chart-container">
            <div ref="departmentChart" class="chart"></div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="charts-row">
      <el-col :span="12">
        <el-card class="chart-card">
          <template #header>
            <span>月度项目趋势</span>
          </template>
          <div class="chart-container">
            <div ref="trendChart" class="chart"></div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card class="chart-card">
          <template #header>
            <span>项目类型对比</span>
          </template>
          <div class="chart-container">
            <div ref="typeChart" class="chart"></div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 详细统计表格 -->
    <el-row :gutter="20" class="table-row">
      <el-col :span="24">
        <el-card class="table-card">
          <template #header>
            <span>院系项目统计详情</span>
          </template>
          <el-table :data="departmentStats" style="width: 100%">
            <el-table-column prop="department" label="院系" />
            <el-table-column prop="totalProjects" label="总项目数" />
            <el-table-column prop="researchProjects" label="科研项目" />
            <el-table-column prop="competitionProjects" label="竞赛项目" />
            <el-table-column prop="pendingProjects" label="待审核" />
            <el-table-column prop="approvedProjects" label="已通过" />
            <el-table-column prop="rejectedProjects" label="已驳回" />
            <el-table-column prop="completionRate" label="完成率">
              <template #default="{ row }">
                <el-progress :percentage="row.completionRate" :color="getProgressColor(row.completionRate)" />
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <!-- 时间范围选择 -->
    <el-row :gutter="20" class="filter-row">
      <el-col :span="24">
        <el-card class="filter-card">
          <template #header>
            <span>统计时间范围</span>
          </template>
          <el-form :inline="true" :model="timeFilter">
            <el-form-item label="开始时间">
              <el-date-picker
                v-model="timeFilter.startDate"
                type="date"
                placeholder="选择开始时间"
                @change="updateStatistics"
              />
            </el-form-item>
            <el-form-item label="结束时间">
              <el-date-picker
                v-model="timeFilter.endDate"
                type="date"
                placeholder="选择结束时间"
                @change="updateStatistics"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="updateStatistics">更新统计</el-button>
              <el-button @click="resetTimeFilter">重置</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'

// 响应式数据
const stats = ref({
  researchCount: 0,
  competitionCount: 0,
  studentCount: 0,
  teacherCount: 0,
  researchPercent: 0,
  competitionPercent: 0,
  participationRate: 0,
  teacherActivity: 0
})

const departmentStats = ref([])
const timeFilter = reactive({
  startDate: null,
  endDate: null
})

// 图表引用
const statusChart = ref(null)
const departmentChart = ref(null)
const trendChart = ref(null)
const typeChart = ref(null)

// 图表实例
let statusChartInstance = null
let departmentChartInstance = null
let trendChartInstance = null
let typeChartInstance = null

// 加载统计数据
const loadStatistics = async () => {
  try {
    // 模拟数据，实际应该从API获取
    stats.value = {
      researchCount: 156,
      competitionCount: 89,
      studentCount: 423,
      teacherCount: 67,
      researchPercent: 63.7,
      competitionPercent: 36.3,
      participationRate: 78.5,
      teacherActivity: 85.2
    }

    departmentStats.value = [
      {
        department: '计算机学院',
        totalProjects: 89,
        researchProjects: 67,
        competitionProjects: 22,
        pendingProjects: 12,
        approvedProjects: 65,
        rejectedProjects: 12,
        completionRate: 73
      },
      {
        department: '机械学院',
        totalProjects: 67,
        researchProjects: 45,
        competitionProjects: 22,
        pendingProjects: 8,
        approvedProjects: 52,
        rejectedProjects: 7,
        completionRate: 78
      },
      {
        department: '经管学院',
        totalProjects: 45,
        researchProjects: 23,
        competitionProjects: 22,
        pendingProjects: 6,
        approvedProjects: 35,
        rejectedProjects: 4,
        completionRate: 78
      },
      {
        department: '其他学院',
        totalProjects: 44,
        researchProjects: 21,
        competitionProjects: 23,
        pendingProjects: 5,
        approvedProjects: 36,
        rejectedProjects: 3,
        completionRate: 82
      }
    ]

    await nextTick()
    initCharts()
  } catch (error) {
    console.error('加载统计数据失败:', error)
    ElMessage.error('加载统计数据失败')
  }
}

// 初始化图表
const initCharts = () => {
  initStatusChart()
  initDepartmentChart()
  initTrendChart()
  initTypeChart()
}

// 项目状态分布图
const initStatusChart = () => {
  if (!statusChart.value) return
  
  statusChartInstance = echarts.init(statusChart.value)
  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      left: 'left'
    },
    series: [
      {
        name: '项目状态',
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: '18',
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: [
          { value: 156, name: '草稿' },
          { value: 89, name: '待审核' },
          { value: 234, name: '已通过' },
          { value: 45, name: '已驳回' },
          { value: 12, name: '已删除' }
        ]
      }
    ]
  }
  statusChartInstance.setOption(option)
}

// 院系项目分布图
const initDepartmentChart = () => {
  if (!departmentChart.value) return
  
  departmentChartInstance = echarts.init(departmentChart.value)
  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    legend: {
      data: ['科研项目', '竞赛项目']
    },
    xAxis: {
      type: 'category',
      data: ['计算机学院', '机械学院', '经管学院', '其他学院']
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        name: '科研项目',
        type: 'bar',
        data: [67, 45, 23, 21]
      },
      {
        name: '竞赛项目',
        type: 'bar',
        data: [22, 22, 22, 23]
      }
    ]
  }
  departmentChartInstance.setOption(option)
}

// 月度项目趋势图
const initTrendChart = () => {
  if (!trendChart.value) return
  
  trendChartInstance = echarts.init(trendChart.value)
  const option = {
    tooltip: {
      trigger: 'axis'
    },
    xAxis: {
      type: 'category',
      data: ['1月', '2月', '3月', '4月', '5月', '6月', '7月', '8月', '9月', '10月', '11月', '12月']
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        name: '项目数量',
        type: 'line',
        smooth: true,
        data: [12, 19, 15, 25, 32, 28, 35, 42, 38, 45, 52, 48]
      }
    ]
  }
  trendChartInstance.setOption(option)
}

// 项目类型对比图
const initTypeChart = () => {
  if (!typeChart.value) return
  
  typeChartInstance = echarts.init(typeChart.value)
  const option = {
    tooltip: {
      trigger: 'item'
    },
    series: [
      {
        name: '项目类型',
        type: 'pie',
        radius: '50%',
        data: [
          { value: 156, name: '科研项目' },
          { value: 89, name: '竞赛项目' }
        ],
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          }
        }
      }
    ]
  }
  typeChartInstance.setOption(option)
}

// 更新统计
const updateStatistics = () => {
  loadStatistics()
}

// 重置时间筛选
const resetTimeFilter = () => {
  timeFilter.startDate = null
  timeFilter.endDate = null
  updateStatistics()
}

// 获取进度条颜色
const getProgressColor = (percentage) => {
  if (percentage >= 80) return '#67C23A'
  if (percentage >= 60) return '#E6A23C'
  return '#F56C6C'
}

// 窗口大小变化时重绘图表
const handleResize = () => {
  if (statusChartInstance) statusChartInstance.resize()
  if (departmentChartInstance) departmentChartInstance.resize()
  if (trendChartInstance) trendChartInstance.resize()
  if (typeChartInstance) typeChartInstance.resize()
}

// 组件挂载
onMounted(() => {
  loadStatistics()
  window.addEventListener('resize', handleResize)
})

// 组件卸载
onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  if (statusChartInstance) statusChartInstance.dispose()
  if (departmentChartInstance) departmentChartInstance.dispose()
  if (trendChartInstance) trendChartInstance.dispose()
  if (typeChartInstance) typeChartInstance.dispose()
})
</script>

<style scoped>
.project-statistics {
  padding: 20px;
}

.stats-overview {
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

.stat-icon.research {
  background: linear-gradient(135deg, #3498db, #2980b9);
}

.stat-icon.competition {
  background: linear-gradient(135deg, #e74c3c, #c0392b);
}

.stat-icon.participation {
  background: linear-gradient(135deg, #27ae60, #2ecc71);
}

.stat-icon.teacher {
  background: linear-gradient(135deg, #9b59b6, #8e44ad);
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

.stat-percent {
  margin: 0;
  color: #7f8c8d;
  font-size: 12px;
}

.charts-row {
  margin-bottom: 20px;
}

.chart-card {
  border-radius: 8px;
}

.chart-container {
  height: 300px;
}

.chart {
  width: 100%;
  height: 100%;
}

.table-row {
  margin-bottom: 20px;
}

.table-card {
  border-radius: 8px;
}

.filter-row {
  margin-bottom: 20px;
}

.filter-card {
  border-radius: 8px;
}

.el-progress {
  margin-top: 8px;
}
</style> 