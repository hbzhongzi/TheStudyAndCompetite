<template>
  <div class="report-view">
    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon><Document /></el-icon>
            </div>
            <div class="stat-info">
              <h4>指导项目</h4>
              <p class="stat-number">{{ stats.projectCount }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon><User /></el-icon>
            </div>
            <div class="stat-info">
              <h4>指导学生</h4>
              <p class="stat-number">{{ stats.studentCount }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon><Trophy /></el-icon>
            </div>
            <div class="stat-info">
              <h4>竞赛指导</h4>
              <p class="stat-number">{{ stats.competitionCount }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon><Clock /></el-icon>
            </div>
            <div class="stat-info">
              <h4>指导时长</h4>
              <p class="stat-number">{{ stats.totalHours }}h</p>
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
            <span>项目类型分布</span>
          </template>
          <div class="chart-container">
            <div ref="projectTypeChart" class="chart"></div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="12">
        <el-card class="chart-card">
          <template #header>
            <span>月度指导时长趋势</span>
          </template>
          <div class="chart-container">
            <div ref="guidanceTrendChart" class="chart"></div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 详细数据表格 -->
    <el-row :gutter="20" class="tables-row">
      <el-col :span="12">
        <el-card class="table-card">
          <template #header>
            <span>学生成绩统计</span>
          </template>
          <el-table :data="studentStats" style="width: 100%">
            <el-table-column prop="name" label="学生姓名" />
            <el-table-column prop="projectCount" label="项目数" />
            <el-table-column prop="competitionCount" label="竞赛数" />
            <el-table-column prop="averageScore" label="平均成绩" />
            <el-table-column prop="achievement" label="主要成就" />
          </el-table>
        </el-card>
      </el-col>
      
      <el-col :span="12">
        <el-card class="table-card">
          <template #header>
            <span>指导记录统计</span>
          </template>
          <el-table :data="guidanceStats" style="width: 100%">
            <el-table-column prop="type" label="指导类型" />
            <el-table-column prop="count" label="次数" />
            <el-table-column prop="totalHours" label="总时长(h)" />
            <el-table-column prop="averageHours" label="平均时长(h)" />
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <!-- 导出功能 -->
    <el-card class="export-card">
      <template #header>
        <span>数据导出</span>
      </template>
      <el-row :gutter="20">
        <el-col :span="6">
          <el-button type="primary" @click="exportStudentReport">
            <el-icon><Download /></el-icon>
            导出学生报告
          </el-button>
        </el-col>
        <el-col :span="6">
          <el-button type="success" @click="exportProjectReport">
            <el-icon><Download /></el-icon>
            导出项目报告
          </el-button>
        </el-col>
        <el-col :span="6">
          <el-button type="warning" @click="exportCompetitionReport">
            <el-icon><Download /></el-icon>
            导出竞赛报告
          </el-button>
        </el-col>
        <el-col :span="6">
          <el-button type="info" @click="exportGuidanceReport">
            <el-icon><Download /></el-icon>
            导出指导报告
          </el-button>
        </el-col>
      </el-row>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Document, User, Trophy, Clock, Download } from '@element-plus/icons-vue'
import * as echarts from 'echarts'

// 响应式数据
const stats = ref({
  projectCount: 15,
  studentCount: 25,
  competitionCount: 8,
  totalHours: 120
})

const studentStats = ref([
  {
    name: '张三',
    projectCount: 3,
    competitionCount: 2,
    averageScore: 85.5,
    achievement: '全国大学生程序设计竞赛二等奖'
  },
  {
    name: '李四',
    projectCount: 2,
    competitionCount: 1,
    averageScore: 82.0,
    achievement: '蓝桥杯程序设计大赛三等奖'
  },
  {
    name: '王五',
    projectCount: 4,
    competitionCount: 3,
    averageScore: 88.5,
    achievement: '全国大学生信息安全竞赛一等奖'
  },
  {
    name: '赵六',
    projectCount: 1,
    competitionCount: 1,
    averageScore: 79.0,
    achievement: '校级创新项目优秀奖'
  }
])

const guidanceStats = ref([
  {
    type: '项目指导',
    count: 45,
    totalHours: 60,
    averageHours: 1.33
  },
  {
    type: '竞赛指导',
    count: 32,
    totalHours: 40,
    averageHours: 1.25
  },
  {
    type: '学术指导',
    count: 18,
    totalHours: 15,
    averageHours: 0.83
  },
  {
    type: '其他指导',
    count: 12,
    totalHours: 5,
    averageHours: 0.42
  }
])

// 图表引用
const projectTypeChart = ref(null)
const guidanceTrendChart = ref(null)

// 初始化项目类型分布图表
const initProjectTypeChart = () => {
  const chart = echarts.init(projectTypeChart.value)
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
        name: '项目类型',
        type: 'pie',
        radius: '50%',
        data: [
          { value: 8, name: '科研项目' },
          { value: 5, name: '创新项目' },
          { value: 2, name: '竞赛项目' }
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
  chart.setOption(option)
}

// 初始化指导时长趋势图表
const initGuidanceTrendChart = () => {
  const chart = echarts.init(guidanceTrendChart.value)
  const option = {
    tooltip: {
      trigger: 'axis'
    },
    xAxis: {
      type: 'category',
      data: ['1月', '2月', '3月', '4月', '5月', '6月', '7月', '8月', '9月', '10月', '11月', '12月']
    },
    yAxis: {
      type: 'value',
      name: '指导时长(h)'
    },
    series: [
      {
        name: '指导时长',
        type: 'line',
        data: [12, 15, 18, 20, 16, 14, 10, 8, 12, 15, 18, 22],
        smooth: true,
        areaStyle: {
          opacity: 0.3
        }
      }
    ]
  }
  chart.setOption(option)
}

// 导出功能
const exportStudentReport = () => {
  ElMessage.success('学生报告导出成功')
}

const exportProjectReport = () => {
  ElMessage.success('项目报告导出成功')
}

const exportCompetitionReport = () => {
  ElMessage.success('竞赛报告导出成功')
}

const exportGuidanceReport = () => {
  ElMessage.success('指导报告导出成功')
}

// 组件挂载时初始化图表
onMounted(() => {
  // 延迟初始化图表，确保DOM已渲染
  setTimeout(() => {
    initProjectTypeChart()
    initGuidanceTrendChart()
  }, 100)
})
</script>

<style scoped>
.report-view {
  padding: 20px;
}

.stats-row {
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
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  display: flex;
  align-items: center;
  justify-content: center;
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

.charts-row {
  margin-bottom: 20px;
}

.chart-card {
  border-radius: 10px;
}

.chart-container {
  height: 300px;
}

.chart {
  width: 100%;
  height: 100%;
}

.tables-row {
  margin-bottom: 20px;
}

.table-card {
  border-radius: 10px;
}

.export-card {
  border-radius: 10px;
}

.export-card .el-button {
  width: 100%;
}
</style> 