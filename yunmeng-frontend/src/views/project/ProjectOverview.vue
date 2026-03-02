<template>
  <div class="project-overview">
    <el-row :gutter="20">
      <!-- 项目统计卡片 -->
      <el-col :span="6" v-for="stat in projectStats" :key="stat.key">
        <el-card class="stat-card" :class="stat.type">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon :size="40">
                <component :is="stat.icon" />
              </el-icon>
            </div>
            <div class="stat-info">
              <h3 class="stat-number">{{ stat.value }}</h3>
              <p class="stat-label">{{ stat.label }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 项目状态分布图表 -->
    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>项目状态分布</span>
          </template>
          <div class="chart-container">
            <div ref="statusChartRef" style="width: 100%; height: 300px;"></div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>项目类型分布</span>
          </template>
          <div class="chart-container">
            <div ref="typeChartRef" style="width: 100%; height: 300px;"></div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 最近项目列表 -->
    <el-row style="margin-top: 20px;">
      <el-col :span="24">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>最近项目</span>
              <el-button type="primary" @click="viewAllProjects">查看全部</el-button>
            </div>
          </template>
          
          <el-table :data="recentProjects" style="width: 100%">
            <el-table-column prop="title" label="项目标题" min-width="200" />
            <el-table-column prop="type" label="项目类型" width="100">
              <template #default="{ row }">
                <el-tag :type="getTypeTagType(row.type)">{{ row.type }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="项目状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getStatusTagType(row.status)">{{ getStatusText(row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="progress" label="进度" width="150">
              <template #default="{ row }">
                <el-progress 
                  :percentage="row.progress || 0" 
                  :status="getProgressStatus(row.progress)"
                />
              </template>
            </el-table-column>
            <el-table-column prop="createdAt" label="创建时间" width="160">
              <template #default="{ row }">
                {{ formatDate(row.createdAt) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="120" fixed="right">
              <template #default="{ row }">
                <el-button size="small" type="primary" @click="viewProject(row.id)">
                  查看详情
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <!-- 项目进度时间线 -->
    <el-row style="margin-top: 20px;">
      <el-col :span="24">
        <el-card>
          <template #header>
            <span>项目进度时间线</span>
          </template>
          
          <el-timeline>
            <el-timeline-item
              v-for="activity in projectActivities"
              :key="activity.id"
              :timestamp="formatDate(activity.time)"
              :type="getActivityType(activity.type)"
            >
              <el-card class="activity-item">
                <div class="activity-header">
                  <h4>{{ activity.title }}</h4>
                  <el-tag :type="getActivityTagType(activity.type)" size="small">
                    {{ getActivityTypeText(activity.type) }}
                  </el-tag>
                </div>
                <p class="activity-content">{{ activity.content }}</p>
                <div class="activity-meta">
                  <span class="project-name">项目: {{ activity.projectName }}</span>
                  <span class="operator">操作人: {{ activity.operator }}</span>
                </div>
              </el-card>
            </el-timeline-item>
          </el-timeline>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'
import  projectService from '@/services/projectService'

export default {
  name: 'ProjectOverview',
  setup() {
    const router = useRouter()
    
    // 响应式数据
    const loading = ref(false)
    const recentProjects = ref([])
    const projectActivities = ref([])
    
    // 图表引用
    const statusChartRef = ref(null)
    const typeChartRef = ref(null)
    
    // 项目统计数据
    const projectStats = reactive([
      {
        key: 'total',
        label: '总项目数',
        value: 0,
        type: 'primary',
        icon: 'Document'
      },
      {
        key: 'inProgress',
        label: '进行中',
        value: 0,
        type: 'success',
        icon: 'Loading'
      },
      {
        key: 'pending',
        label: '待审核',
        value: 0,
        type: 'warning',
        icon: 'Clock'
      },
      {
        key: 'completed',
        label: '已完成',
        value: 0,
        type: 'info',
        icon: 'Check'
      }
    ])
    
    // 方法
    const loadProjectStats = async () => {
      try {
        const response = await projectService.getProjectStats()
        const stats = response.data
        
        // 更新统计数据
        projectStats[0].value = stats.totalProjects || 0
        projectStats[1].value = stats.inProgressProjects || 0
        projectStats[2].value = stats.pendingProjects || 0
        projectStats[3].value = stats.completedProjects || 0
        
        // 渲染图表
        renderStatusChart(stats)
        renderTypeChart(stats)
      } catch (error) {
        ElMessage.error('加载项目统计失败')
      }
    }
    
    const loadRecentProjects = async () => {
      try {
        const response = await projectService.getMyProjects()
        recentProjects.value = response.data.slice(0, 10) // 只显示最近10个项目
      } catch (error) {
        ElMessage.error('加载最近项目失败')
      }
    }
    
    const loadProjectActivities = async () => {
      try {
        // 这里需要后端提供项目活动记录的API
        // 暂时使用模拟数据
        projectActivities.value = [
          {
            id: 1,
            title: '项目创建',
            content: '学生张三创建了新项目"智能校园管理系统"',
            type: 'create',
            projectName: '智能校园管理系统',
            operator: '张三',
            time: new Date()
          },
          {
            id: 2,
            title: '项目提交',
            content: '项目"智能校园管理系统"已提交审核',
            type: 'submit',
            projectName: '智能校园管理系统',
            operator: '张三',
            time: new Date(Date.now() - 86400000)
          },
          {
            id: 3,
            title: '项目审核通过',
            content: '项目"智能校园管理系统"审核通过，可以开始执行',
            type: 'approve',
            projectName: '智能校园管理系统',
            operator: '李老师',
            time: new Date(Date.now() - 172800000)
          }
        ]
      } catch (error) {
        console.error('加载项目活动失败:', error)
      }
    }
    
    const renderStatusChart = (stats) => {
      if (!statusChartRef.value) return
      
      const chart = echarts.init(statusChartRef.value)
      const option = {
        title: {
          text: '项目状态分布',
          left: 'center'
        },
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
            radius: '50%',
            data: [
              { value: stats.draftProjects || 0, name: '草稿' },
              { value: stats.pendingProjects || 0, name: '待审核' },
              { value: stats.approvedProjects || 0, name: '已通过' },
              { value: stats.inProgressProjects || 0, name: '进行中' },
              { value: stats.completedProjects || 0, name: '已完成' },
              { value: stats.rejectedProjects || 0, name: '已驳回' }
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
    
    const renderTypeChart = (stats) => {
      if (!typeChartRef.value) return
      
      const chart = echarts.init(typeChartRef.value)
      const option = {
        title: {
          text: '项目类型分布',
          left: 'center'
        },
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            type: 'shadow'
          }
        },
        xAxis: {
          type: 'category',
          data: ['科研项目', '竞赛项目', '创新项目', '实践项目']
        },
        yAxis: {
          type: 'value'
        },
        series: [
          {
            name: '项目数量',
            type: 'bar',
            data: [
              stats.researchProjects || 0,
              stats.competitionProjects || 0,
              stats.innovationProjects || 0,
              stats.practiceProjects || 0
            ],
            itemStyle: {
              color: '#409eff'
            }
          }
        ]
      }
      chart.setOption(option)
    }
    
    const viewAllProjects = () => {
      router.push('/projects')
    }
    
    const viewProject = (projectId) => {
      router.push(`/project/detail/${projectId}`)
    }
    
    // 工具方法
    const formatDate = (date) => {
      if (!date) return ''
      return new Date(date).toLocaleString('zh-CN')
    }
    
    const getTypeTagType = (type) => {
      const typeMap = {
        '科研': 'primary',
        '竞赛': 'success',
        '创新': 'warning',
        '实践': 'info'
      }
      return typeMap[type] || 'default'
    }
    
    const getStatusTagType = (status) => {
      const statusMap = {
        draft: 'info',
        pending: 'warning',
        approved: 'success',
        rejected: 'danger',
        in_progress: 'primary',
        completed: 'success'
      }
      return statusMap[status] || 'info'
    }
    
    const getStatusText = (status) => {
      const statusMap = {
        draft: '草稿',
        pending: '待审核',
        approved: '已通过',
        rejected: '已驳回',
        in_progress: '进行中',
        completed: '已完成'
      }
      return statusMap[status] || status
    }
    
    const getProgressStatus = (progress) => {
      if (progress >= 100) return 'success'
      if (progress >= 80) return 'warning'
      if (progress >= 50) return ''
      return 'exception'
    }
    
    const getActivityType = (type) => {
      const typeMap = {
        create: 'primary',
        submit: 'warning',
        approve: 'success',
        reject: 'danger',
        update: 'info'
      }
      return typeMap[type] || 'info'
    }
    
    const getActivityTagType = (type) => {
      const typeMap = {
        create: 'primary',
        submit: 'warning',
        approve: 'success',
        reject: 'danger',
        update: 'info'
      }
      return typeMap[type] || 'info'
    }
    
    const getActivityTypeText = (type) => {
      const typeMap = {
        create: '创建',
        submit: '提交',
        approve: '审核通过',
        reject: '审核驳回',
        update: '更新'
      }
      return typeMap[type] || type
    }
    
    // 生命周期
    onMounted(() => {
      loadProjectStats()
      loadRecentProjects()
      loadProjectActivities()
    })
    
    return {
      // 响应式数据
      loading,
      recentProjects,
      projectActivities,
      projectStats,
      
      // 图表引用
      statusChartRef,
      typeChartRef,
      
      // 方法
      viewAllProjects,
      viewProject,
      
      // 工具方法
      formatDate,
      getTypeTagType,
      getStatusTagType,
      getStatusText,
      getProgressStatus,
      getActivityType,
      getActivityTagType,
      getActivityTypeText
    }
  }
}
</script>

<style scoped>
.project-overview {
  padding: 20px;
}

.stat-card {
  margin-bottom: 20px;
  transition: transform 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-5px);
}

.stat-card.primary {
  border-left: 4px solid #409eff;
}

.stat-card.success {
  border-left: 4px solid #67c23a;
}

.stat-card.warning {
  border-left: 4px solid #e6a23c;
}

.stat-card.info {
  border-left: 4px solid #909399;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 20px;
}

.stat-icon {
  color: #409eff;
}

.stat-info {
  flex: 1;
}

.stat-number {
  font-size: 28px;
  font-weight: bold;
  margin: 0 0 5px 0;
  color: #303133;
}

.stat-label {
  margin: 0;
  color: #606266;
  font-size: 14px;
}

.chart-container {
  padding: 20px 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.activity-item {
  margin-bottom: 10px;
}

.activity-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.activity-header h4 {
  margin: 0;
  color: #303133;
}

.activity-content {
  margin: 10px 0;
  color: #606266;
}

.activity-meta {
  display: flex;
  justify-content: space-between;
  margin-top: 15px;
  font-size: 12px;
  color: #909399;
}

.project-name,
.operator {
  background-color: #f5f7fa;
  padding: 4px 8px;
  border-radius: 4px;
}

.el-timeline-item {
  padding-bottom: 20px;
}

.el-timeline-item:last-child {
  padding-bottom: 0;
}
</style> 