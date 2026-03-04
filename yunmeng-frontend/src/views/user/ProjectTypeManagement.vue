<template>
  <div class="statistics-container">

    <!-- ================= 统计卡片 ================= -->
    <el-row :gutter="20" class="card-row">
      <el-col :span="4" v-for="card in cards" :key="card.label">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-value">{{ card.value }}</div>
          <div class="stat-label">{{ card.label }}</div>
        </el-card>
      </el-col>
    </el-row>

    <!-- ================= 图表区域 ================= -->
    <el-row :gutter="20" class="chart-row">

      <el-col :span="12">
        <el-card shadow="never">
          <div class="chart-title">项目状态分布</div>
          <div id="pieChart" class="chart"></div>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card shadow="never">
          <div class="chart-title">月度项目趋势</div>
          <div id="lineChart" class="chart"></div>
        </el-card>
      </el-col>

    </el-row>

  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import * as echarts from 'echarts'
import adminService from '@/services/adminService'

// ================= 统计卡片数据 =================
const cards = ref([
  { label: '项目总数', value: 0 },
  { label: '草稿', value: 0 },
  { label: '已提交', value: 0 },
  { label: '审核中', value: 0 },
  { label: '已通过', value: 0 },
  { label: '已拒绝', value: 0 }
])

// ================= 获取总体统计 =================
const loadOverview = async () => {
  const res = await adminService.getOverview()
  const data = res.data

  cards.value[0].value = data.total
  cards.value[1].value = data.draft
  cards.value[2].value = data.submitted
  cards.value[3].value = data.reviewing
  cards.value[4].value = data.approved
  cards.value[5].value = data.rejected
}

// ================= 状态饼图 =================
const initPieChart = (data) => {
  const chart = echarts.init(document.getElementById('pieChart'))

  chart.setOption({
    tooltip: { trigger: 'item' },
    legend: { bottom: 0 },
    series: [
      {
        type: 'pie',
        radius: '65%',
        data: data.map(item => ({
          name: item.Status,
          value: item.Count
        }))
      }
    ]
  })
}

const loadStatusStats = async () => {
  const res = await adminService.getStatusStats()
  initPieChart(res.data)
}

// ================= 折线图 =================
const initLineChart = (data) => {
  const chart = echarts.init(document.getElementById('lineChart'))

  chart.setOption({
    tooltip: { trigger: 'axis' },
    xAxis: {
      type: 'category',
      data: data.map(i => i.Month)
    },
    yAxis: { type: 'value' },
    series: [
      {
        data: data.map(i => i.Count),
        type: 'line',
        smooth: true
      }
    ]
  })
}

const loadMonthlyStats = async () => {
  const res = await adminService.getmonthlyOverview()
  initLineChart(res.data)
}

// ================= 生命周期 =================
onMounted(async () => {
  await loadOverview()
  await nextTick()
  loadStatusStats()
  loadMonthlyStats()
})
</script>

<style scoped>
.statistics-container {
  padding: 20px;
}

.card-row {
  margin-bottom: 30px;
}

.stat-card {
  text-align: center;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #409eff;
}

.stat-label {
  margin-top: 10px;
  color: #666;
}

.chart-row {
  margin-top: 20px;
}

.chart {
  height: 400px;
}

.chart-title {
  font-weight: bold;
  margin-bottom: 15px;
}
</style>