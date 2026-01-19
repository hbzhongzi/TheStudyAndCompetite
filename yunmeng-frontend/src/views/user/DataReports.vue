<template>
  <div class="data-reports">
    <div class="page-header">
      <h2>数据统计</h2>
      <div class="header-actions">
        <el-button type="primary" @click="exportReport">
          <i class="el-icon-download"></i>
          导出报表
        </el-button>
        <el-button type="success" @click="refreshData">
          <i class="el-icon-refresh"></i>
          刷新数据
        </el-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stat-cards">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon user-icon">
              <i class="el-icon-user"></i>
            </div>
            <div class="stat-info">
              <h4>总用户数</h4>
              <p class="stat-number">{{ statistics.totalUsers }}</p>
              <p class="stat-change positive">+{{ statistics.userGrowth }}%</p>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon project-icon">
              <i class="el-icon-document"></i>
            </div>
            <div class="stat-info">
              <h4>活跃项目</h4>
              <p class="stat-number">{{ statistics.activeProjects }}</p>
              <p class="stat-change positive">+{{ statistics.projectGrowth }}%</p>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon competition-icon">
              <i class="el-icon-trophy"></i>
            </div>
            <div class="stat-info">
              <h4>进行中竞赛</h4>
              <p class="stat-number">{{ statistics.ongoingCompetitions }}</p>
              <p class="stat-change positive">+{{ statistics.competitionGrowth }}%</p>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon achievement-icon">
              <i class="el-icon-medal"></i>
            </div>
            <div class="stat-info">
              <h4>获奖数量</h4>
              <p class="stat-number">{{ statistics.totalAwards }}</p>
              <p class="stat-change positive">+{{ statistics.awardGrowth }}%</p>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 图表区域 -->
    <el-row :gutter="20" class="chart-section">
      <el-col :span="12">
        <el-card class="chart-card">
          <template #header>
            <span>用户增长趋势</span>
          </template>
          <div class="chart-container">
            <div ref="userChartRef" style="width: 100%; height: 300px;"></div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="12">
        <el-card class="chart-card">
          <template #header>
            <span>项目类型分布</span>
          </template>
          <div class="chart-container">
            <div ref="projectChartRef" style="width: 100%; height: 300px;"></div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="chart-section">
      <el-col :span="12">
        <el-card class="chart-card">
          <template #header>
            <span>竞赛参与情况</span>
          </template>
          <div class="chart-container">
            <div ref="competitionChartRef" style="width: 100%; height: 300px;"></div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="12">
        <el-card class="chart-card">
          <template #header>
            <span>月度活跃度</span>
          </template>
          <div class="chart-container">
            <div ref="activityChartRef" style="width: 100%; height: 300px;"></div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 详细数据表格 -->
    <el-card class="data-table-card">
      <template #header>
        <span>详细数据报表</span>
      </template>
      
      <el-tabs v-model="activeTab" type="border-card">
        <el-tab-pane label="用户统计" name="users">
          <el-table :data="userData" style="width: 100%">
            <el-table-column prop="month" label="月份" width="120" />
            <el-table-column prop="newUsers" label="新增用户" width="120" />
            <el-table-column prop="activeUsers" label="活跃用户" width="120" />
            <el-table-column prop="totalUsers" label="总用户数" width="120" />
            <el-table-column prop="growthRate" label="增长率" width="120">
              <template #default="scope">
                <span :class="scope.row.growthRate >= 0 ? 'positive' : 'negative'">
                  {{ scope.row.growthRate >= 0 ? '+' : '' }}{{ scope.row.growthRate }}%
                </span>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        
        <el-tab-pane label="项目统计" name="projects">
          <el-table :data="projectData" style="width: 100%">
            <el-table-column prop="month" label="月份" width="120" />
            <el-table-column prop="newProjects" label="新增项目" width="120" />
            <el-table-column prop="completedProjects" label="完成项目" width="120" />
            <el-table-column prop="ongoingProjects" label="进行中项目" width="120" />
            <el-table-column prop="completionRate" label="完成率" width="120">
              <template #default="scope">
                {{ scope.row.completionRate }}%
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        
        <el-tab-pane label="竞赛统计" name="competitions">
          <el-table :data="competitionData" style="width: 100%">
            <el-table-column prop="month" label="月份" width="120" />
            <el-table-column prop="newCompetitions" label="新增竞赛" width="120" />
            <el-table-column prop="participants" label="参赛人数" width="120" />
            <el-table-column prop="awards" label="获奖数量" width="120" />
            <el-table-column prop="participationRate" label="参与率" width="120">
              <template #default="scope">
                {{ scope.row.participationRate }}%
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'

// 响应式数据
const activeTab = ref('users')
const userChartRef = ref()
const projectChartRef = ref()
const competitionChartRef = ref()
const activityChartRef = ref()

// 统计数据
const statistics = reactive({
  totalUsers: 1234,
  userGrowth: 12.5,
  activeProjects: 89,
  projectGrowth: 8.3,
  ongoingCompetitions: 15,
  competitionGrowth: 15.7,
  totalAwards: 67,
  awardGrowth: 22.1
})

// 用户数据
const userData = ref([
  { month: '2024-01', newUsers: 45, activeUsers: 890, totalUsers: 1234, growthRate: 12.5 },
  { month: '2024-02', newUsers: 52, activeUsers: 920, totalUsers: 1286, growthRate: 8.9 },
  { month: '2024-03', newUsers: 38, activeUsers: 950, totalUsers: 1324, growthRate: 6.2 },
  { month: '2024-04', newUsers: 61, activeUsers: 980, totalUsers: 1385, growthRate: 15.3 },
  { month: '2024-05', newUsers: 47, activeUsers: 1010, totalUsers: 1432, growthRate: 10.2 },
  { month: '2024-06', newUsers: 55, activeUsers: 1050, totalUsers: 1487, growthRate: 12.8 }
])

// 项目数据
const projectData = ref([
  { month: '2024-01', newProjects: 12, completedProjects: 8, ongoingProjects: 89, completionRate: 85.2 },
  { month: '2024-02', newProjects: 15, completedProjects: 10, ongoingProjects: 94, completionRate: 87.3 },
  { month: '2024-03', newProjects: 18, completedProjects: 12, ongoingProjects: 100, completionRate: 89.1 },
  { month: '2024-04', newProjects: 14, completedProjects: 9, ongoingProjects: 105, completionRate: 86.7 },
  { month: '2024-05', newProjects: 20, completedProjects: 15, ongoingProjects: 110, completionRate: 91.2 },
  { month: '2024-06', newProjects: 16, completedProjects: 11, ongoingProjects: 115, completionRate: 88.9 }
])

// 竞赛数据
const competitionData = ref([
  { month: '2024-01', newCompetitions: 3, participants: 156, awards: 12, participationRate: 78.5 },
  { month: '2024-02', newCompetitions: 4, participants: 189, awards: 15, participationRate: 82.1 },
  { month: '2024-03', newCompetitions: 5, participants: 234, awards: 18, participationRate: 85.7 },
  { month: '2024-04', newCompetitions: 2, participants: 167, awards: 10, participationRate: 76.3 },
  { month: '2024-05', newCompetitions: 6, participants: 298, awards: 22, participationRate: 89.2 },
  { month: '2024-06', newCompetitions: 4, participants: 245, awards: 16, participationRate: 83.4 }
])

// 图表实例
let userChart = null
let projectChart = null
let competitionChart = null
let activityChart = null

// 方法
const initCharts = async () => {
  await nextTick()
  
  // 初始化用户增长趋势图
  if (userChartRef.value) {
    userChart = echarts.init(userChartRef.value)
    const userOption = {
      title: {
        text: '用户增长趋势',
        left: 'center'
      },
      tooltip: {
        trigger: 'axis'
      },
      xAxis: {
        type: 'category',
        data: userData.value.map(item => item.month)
      },
      yAxis: {
        type: 'value'
      },
      series: [
        {
          name: '总用户数',
          type: 'line',
          data: userData.value.map(item => item.totalUsers),
          smooth: true,
          itemStyle: {
            color: '#409EFF'
          }
        },
        {
          name: '新增用户',
          type: 'bar',
          data: userData.value.map(item => item.newUsers),
          itemStyle: {
            color: '#67C23A'
          }
        }
      ]
    }
    userChart.setOption(userOption)
  }
  
  // 初始化项目类型分布图
  if (projectChartRef.value) {
    projectChart = echarts.init(projectChartRef.value)
    const projectOption = {
      title: {
        text: '项目类型分布',
        left: 'center'
      },
      tooltip: {
        trigger: 'item'
      },
      series: [
        {
          name: '项目类型',
          type: 'pie',
          radius: '50%',
          data: [
            { value: 45, name: '科研项目' },
            { value: 28, name: '竞赛项目' },
            { value: 16, name: '创新项目' }
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
    projectChart.setOption(projectOption)
  }
  
  // 初始化竞赛参与情况图
  if (competitionChartRef.value) {
    competitionChart = echarts.init(competitionChartRef.value)
    const competitionOption = {
      title: {
        text: '竞赛参与情况',
        left: 'center'
      },
      tooltip: {
        trigger: 'axis'
      },
      xAxis: {
        type: 'category',
        data: competitionData.value.map(item => item.month)
      },
      yAxis: {
        type: 'value'
      },
      series: [
        {
          name: '参赛人数',
          type: 'bar',
          data: competitionData.value.map(item => item.participants),
          itemStyle: {
            color: '#E6A23C'
          }
        },
        {
          name: '获奖数量',
          type: 'line',
          data: competitionData.value.map(item => item.awards),
          smooth: true,
          itemStyle: {
            color: '#F56C6C'
          }
        }
      ]
    }
    competitionChart.setOption(competitionOption)
  }
  
  // 初始化月度活跃度图
  if (activityChartRef.value) {
    activityChart = echarts.init(activityChartRef.value)
    const activityOption = {
      title: {
        text: '月度活跃度',
        left: 'center'
      },
      tooltip: {
        trigger: 'axis'
      },
      xAxis: {
        type: 'category',
        data: ['1月', '2月', '3月', '4月', '5月', '6月']
      },
      yAxis: {
        type: 'value'
      },
      series: [
        {
          name: '活跃用户',
          type: 'line',
          data: [890, 920, 950, 980, 1010, 1050],
          areaStyle: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [
                { offset: 0, color: 'rgba(64, 158, 255, 0.3)' },
                { offset: 1, color: 'rgba(64, 158, 255, 0.1)' }
              ]
            }
          },
          itemStyle: {
            color: '#409EFF'
          }
        }
      ]
    }
    activityChart.setOption(activityOption)
  }
}

const exportReport = () => {
  ElMessage.success('报表导出成功')
  // 这里可以实现实际的导出功能
}

const refreshData = () => {
  ElMessage.success('数据刷新成功')
  // 这里可以重新加载数据
}

// 生命周期
onMounted(() => {
  initCharts()
  
  // 监听窗口大小变化，重新调整图表大小
  window.addEventListener('resize', () => {
    userChart?.resize()
    projectChart?.resize()
    competitionChart?.resize()
    activityChart?.resize()
  })
})
</script>

<style scoped>
.data-reports {
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

.stat-cards {
  margin-bottom: 20px;
}

.stat-card {
  border-radius: 10px;
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

.stat-icon i {
  font-size: 24px;
  color: white;
}

.user-icon {
  background: linear-gradient(135deg, #409EFF 0%, #36A3F7 100%);
}

.project-icon {
  background: linear-gradient(135deg, #67C23A 0%, #85CE61 100%);
}

.competition-icon {
  background: linear-gradient(135deg, #E6A23C 0%, #F0AD4E 100%);
}

.achievement-icon {
  background: linear-gradient(135deg, #F56C6C 0%, #F78989 100%);
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

.stat-change {
  margin: 0;
  font-size: 12px;
  font-weight: 500;
}

.stat-change.positive {
  color: #67C23A;
}

.stat-change.negative {
  color: #F56C6C;
}

.chart-section {
  margin-bottom: 20px;
}

.chart-card {
  border-radius: 10px;
}

.chart-container {
  padding: 10px;
}

.data-table-card {
  margin-bottom: 20px;
  border-radius: 10px;
}

.positive {
  color: #67C23A;
}

.negative {
  color: #F56C6C;
}
</style> 